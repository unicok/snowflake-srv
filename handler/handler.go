package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unsafe"

	"github.com/unicok/misc/log"
	proto "github.com/unicok/snowflake/proto/snowflake"

	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
)

const (
	backoff       = 100           // max backoff delay millisecond
	uuidQueue     = 1024          // uuid process queue
	tsMask        = 0x1FFFFFFFFFF // 41bit
	snMask        = 0xFFF         // 12bit
	machineIDMask = 0x3FF         // 10bit
)

type snowflake struct {
	*api.Client
	kv *api.KV

	ns        string
	uuidKey   string
	machineID uint64 // 10-bit machine append

	procCh chan chan uint64
}

func NewSnowflake(mcid, ns, uuidKey, addr string) *snowflake {
	p := new(snowflake)

	config := api.DefaultConfig()
	config.Address = addr

	var err error
	p.Client, err = api.NewClient(config)
	if err != nil {
		log.Fatal("create consul kv error:", err)
	}

	p.kv = p.Client.KV()
	p.ns = ns
	p.uuidKey = uuidKey
	p.procCh = make(chan chan uint64, uuidQueue)

	if len(mcid) > 0 {
		id, err := strconv.Atoi(mcid)
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		p.machineID = (uint64(id) & machineIDMask) << 12
		log.Info("machine id from env specified:", p.machineID)
	} else {
		p.initMachineID()
		log.Info("machine id from kv specified:", p.machineID)
	}

	go p.uuidTask()

	return p
}

func (p *snowflake) initMachineID() {
	for {
		// check the key
		ks, _, err := p.kv.Keys(p.uuidKey, "", nil)
		checkErrPanic(err)

		if len(ks) == 0 {
			_, err := p.kv.Put(&api.KVPair{
				Key:   p.uuidKey,
				Value: []byte("0")}, nil)
			checkErrPanic(err)
		}

		// get the key
		kvpair, _, err := p.kv.Get(p.uuidKey, nil)
		checkErrPanic(err)

		// get preValue & preIndex
		prevValue, err := strconv.Atoi(Bytes2Str(kvpair.Value))
		checkErrPanic(err)
		prevIndex := kvpair.ModifyIndex

		// compareAndSwap
		_, _, err = p.kv.CAS(&api.KVPair{
			Key:         p.uuidKey,
			Value:       []byte(fmt.Sprint(prevValue + 1)),
			ModifyIndex: prevIndex}, nil)
		if err != nil {
			casDelay()
			continue
		}

		// record serial number of this service, already shifted
		p.machineID = (uint64(prevValue+1) & machineIDMask) << 12
		return
	}
}

// Next is get next value of a key, like auto-incrememt in mysql
func (p *snowflake) Next(ctx context.Context, in *proto.Key, out *proto.Value) error {
	key := p.ns + in.Name
	for {
		// get the key
		kvpair, _, err := p.kv.Get(key, nil)
		if err != nil || kvpair == nil {
			log.Fatal(err)
			return errors.New("Key not exists, need to create first")
		}

		// get prevValue & prevIndex
		prevValue, err := strconv.Atoi(Bytes2Str(kvpair.Value))
		if err != nil {
			log.Fatal(err)
			return errors.New("marlformed value")
		}
		prevIndex := kvpair.ModifyIndex

		// compareAndSwap
		_, _, err = p.kv.CAS(&api.KVPair{
			Key:         key,
			Value:       []byte(fmt.Sprint(prevValue + 1)),
			ModifyIndex: prevIndex}, nil)
		if err != nil {
			casDelay()
			continue
		}

		out.Value = int64(prevValue + 1)
		return nil
	}
}

// GetUUID is generate an unique uuid
func (p *snowflake) GetUUID(ctx context.Context, req *proto.NullRequest, rsp *proto.UUID) error {
	q := make(chan uint64, 1)
	p.procCh <- q
	rsp.Uuid = <-q
	return nil
}

// uuid generator
func (p *snowflake) uuidTask() {
	var sn uint64    // 12-bit serial no
	var lastts int64 // last timestamp
	for {
		ret := <-p.procCh
		// get a correct serial number
		t := ts()
		// clock shift backward
		if t < lastts {
			log.Fatal("clock shift happened, waiting until the clock moving to the next millisecond.")
			t = p.waitMs(lastts)
		}

		// same millisecond
		if lastts == t {
			sn = (sn + 1) & snMask
			// serial number overflows, wait until next ms
			if sn == 0 {
				t = p.waitMs(lastts)
			}
		} else { // new millsecond, reset serial number to 0
			sn = 0
		}
		// remember last timestamp
		lastts = t

		// generate uuid, format:
		//
		// 0		0.................0		0..............0	0........0
		// 1-bit	41bit timestamp			10bit machine-id	12bit sn
		var uuid uint64
		uuid |= (uint64(t) & tsMask) << 22
		uuid |= p.machineID
		uuid |= sn
		ret <- uuid
	}
}

// waitMs will spin wait till next millisecond
func (p *snowflake) waitMs(lastts int64) int64 {
	t := ts()
	for t <= lastts {
		t = ts()
	}
	return t
}

// random delay
func casDelay() {
	<-time.After(time.Duration(rand.Int63n(backoff)) * time.Millisecond)
}

// get timestamp
func ts() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func checkErrPanic(err error) {
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
