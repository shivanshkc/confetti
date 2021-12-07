# Confetti

A minimal configuration manager for Go applications.

## Quickstart

Confetti makes it really easy to read configs from environment variables and command-line flags.

Consider the following snippet:

```go
package main

import (
	"fmt"

	"github.com/shivanshkc/confetti"
)

type configs struct {
	HTTPServer struct {
		Addr string `default:"0.0.0.0:8080" env:"HTTP_SERVER_ADDR" arg:"http-server-addr"`
	}

	GRPCServer struct {
		Addr string `default:"0.0.0.0:9090" env:"GRPC_SERVER_ADDR" arg:"grpc-server-addr"`
	}

	LogLevel string `default:"info" env:"LOG_LEVEL" arg:"log-level"`
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

Try it on Playground: https://play.golang.org/p/4xV4XL8eUol

Here are the main points to be noticed in this code snippet:

- The `configs` struct is the schema for our configs.
- We can bind the fields to a default value by using the `default` tag.
- We can bind the fields to environment variables using the `env` tag.
- We can bind the fields to command-line args using the `arg` tag.
- To actually load the configs, we need the `confetti.Loader` type, which is provided by `confetti.GetLoader` function.
- The configs are loaded into an instance of our struct, by provided it to the `loader.Load` method.

## Examples

With the provided code-snippet in mind:

1.  The following will output all the default config values.

    ```bash
    go run main.go
    ```

2.  The following will output all the default configs, except for the `HTTPServer.Addr` config, whose default value will be overridden by the environment variable.

    ```bash
    HTTP_SERVER_ADDR=localhost:3000 go run main.go
    ```

3.  This time the `HTTPServer.Addr` property will be equal to the command-line flag provided, **as it is prioritized above the environment variable**.

    ```bash
    HTTP_SERVER_ADDR=localhost:3000 go run main.go --http-server-addr=localhost:4000
    ```

    Make sure to provide your command-line args in: `--<name>=<value>` format.
