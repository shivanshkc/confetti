# Confetti

A minimal configuration manager for Go applications.

## Quickstart

```go
package main

import (
	"fmt"

	"github.com/shivanshkc/confetti"
)

type configs struct {
	LogLevel string `default:"info" env:"LOG_LEVEL" arg:"log-level"`
	GRPC struct {
		Addr    string `default:"0.0.0.0:9090" env:"GRPC_ADDR" arg:"grpc-addr"`
		Timeout int    `default:"60" env:"GRPC_TIMEOUT" arg:"grpc-timeout"`
	}
}

func main() {
	conf := &configs{}

	loader := confetti.GetLoader()
	if err := loader.Load(conf); err != nil {
		panic(err)
	}

	fmt.Printf("Conf: %+v\n", conf)
}
```
