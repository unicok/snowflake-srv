# SNOWFLAKE

[![Build Status](https://travis-ci.org/unicok/snowflake.svg?branch=master)](https://travis-ci.org/unicok/snowflake)

Snowfalke 是分布式uuid发生器，twitter snowflake的go语言版本

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
	go get github.com/unicok/snowflake
	./snowflake
	```

## The API
- GetUUID
- Next       

uuid格式为:

    +-------------------------------------------------------------------------------------------------+
    | UNUSED(1BIT) |         TIMESTAMP(41BIT)           |  MACHINE-ID(10BIT)  |   SERIAL-NO(12BIT)    |
    +-------------------------------------------------------------------------------------------------+ 