// +build integration

package handler

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/micro/go-micro/client"
	proto "github.com/unicok/snowflake-srv/proto/snowflake"
)

const (
	test_key = "test_key"
)

func TestCasDelay(t *testing.T) {
	casDelay()
}

func TestSnowflake(t *testing.T) {
	// Set up a connection to the server.
	req := client.NewRequest("com.unicok.srv.snowflake", "Snowflake.Next", &proto.Key{Name: test_key})

	rsp := &proto.Value{}
	if err := client.Call(context.Background(), req, rsp); err != nil {
		t.Fatalf("could not get next value: %v", err)
	}

	t.Log(rsp.Value)
}

func BenchmarkSnowflake(b *testing.B) {
	// Set up a connection to the server.
	// Set up a connection to the server.
	req := client.NewRequest("com.unicok.srv.snowflake", "Snowflake.Next", &proto.Key{Name: test_key})
	rsp := &proto.Value{}

	for i := 0; i < b.N; i++ {
		// Contact the server and print out its response.
		if err := client.Call(context.Background(), req, rsp); err != nil {
			b.Fatalf("could not get next value: %v", err)
		}
	}
}

func TestSnowflakeUUID(t *testing.T) {
	// Set up a connection to the server.
	req := client.NewRequest("com.unicok.srv.snowflake", "Snowflake.GetUUID", &proto.NullRequest{})

	rsp := &proto.Value{}
	if err := client.Call(context.Background(), req, rsp); err != nil {
		t.Fatalf("could not get uuid value: %v", err)
	}

	t.Log(rsp.Value)
}

func BenchmarkSnowflakeUUID(b *testing.B) {
	// Set up a connection to the server.
	req := client.NewRequest("com.unicok.srv.snowflake", "Snowflake.GetUUID", &proto.NullRequest{})
	rsp := &proto.Value{}

	for i := 0; i < b.N; i++ {
		// Contact the server and print out its response.
		if err := client.Call(context.Background(), req, rsp); err != nil {
			b.Fatalf("could not get uuid: %v", err)
		}
	}
}
