# Snowflake

[![Build Status](https://travis-ci.org/unicok/snowflake-srv.svg?branch=master)](https://travis-ci.org/unicok/snowflake-srv)

[Snowflake](https://github.com/twitter/snowflake) is a network service for generating unique ID numbers at high scale with some simple guarantees. http://twitter.com/
This is a go language version based on [Consul KV](https://github.com/hashicorp/consul)

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

4. Download and start the service

	```shell
	go get github.com/unicok/snowflake-srv
	./snowflake-srv
	```

## The API
- GetUUID
- Next       

UUID is composed of:
	* unused - 1 bit 
	* time - 41 bits (millisecond precision w/ a custom epoch gives us 69 years)
  	* configured machine id - 10 bits - gives us up to 1024 machines
  	* sequence number - 12 bits - rolls over every 4096 per machine (with protection to avoid rollover in the same ms)

    +-------------------------------------------------------------------------------------------------+
    | UNUSED(1BIT) |         TIMESTAMP(41BIT)           |  MACHINE-ID(10BIT)  |   SERIAL-NO(12BIT)    |
    +-------------------------------------------------------------------------------------------------+ 