package main

import (
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/unicok/misc/log"
	"github.com/unicok/snowflake-srv/handler"
	proto "github.com/unicok/snowflake-srv/proto/snowflake"
)

const (
	default_namespace = "seqs/"
	default_uuidKey   = "snowflake-uuid"
)

func main() {
	service := micro.NewService(
		micro.Name("com.unicok.srv.snowflake"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),

		micro.Flags(
			cli.StringFlag{
				Name:   "uuid_key",
				EnvVar: "UUID_KEY",
				Usage:  "Name for the key-value store",
			},
			cli.StringFlag{
				Name:   "consul_address",
				EnvVar: "CONSUL_ADDRESS",
				Usage:  "Comma-separated list of consul addresses for kv",
			},
			cli.StringFlag{
				Name:   "machine_id",
				EnvVar: "MACHINE_ID",
				Usage:  "specific machine id",
			},
		),
	)

	var (
		uuidKey    string
		consulAddr string
		machineID  string
	)

	service.Init(
		micro.Action(func(c *cli.Context) {

			uuidKey = default_namespace + default_uuidKey
			if len(c.String("uuid_key")) > 0 {
				uuidKey = c.String("uuid_key")
			}

			consulAddr = "127.0.0.1:8500"
			if len(c.String("consul_address")) > 0 {
				consulAddr = c.String("consul_address")
			}

			machineID = ""
			if len(c.String("machine_id")) > 0 {
				machineID = c.String("machine_id")
			}
		}),
	)

	sf := handler.NewSnowflake(machineID, default_namespace, uuidKey, consulAddr)
	proto.RegisterSnowflakeHandler(service.Server(), sf)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
