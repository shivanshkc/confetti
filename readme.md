# Confetti
A simple config manager for Go applications.

## Install
Use the following:  
```
go get -u github.com/shivanshkc/confetti/v2
```

## When to use Confetti
Confetti only has a few, but well implemented set of features. It makes it really easy to use and understand.  
If your application fits the following use-case, Confetti is the best config manager you can get.  
1. The configs have to be loaded/unmarshalled into a struct.
2. The configs are loaded once at application startup and do not change for the entire runtime of the application.
3. The configs are loaded either from the environment or command-line flags (No JSON/YAML files or remote servers).

## Beauty of Confetti
If your application agrees with the above restrictions, you can enjoy the following features of Confetti:  
1. ### Concise syntax using struct tags
    ```go
    package main

    type Configs struct {
        Port string `def:"8080" env:"PORT" arg:"port"`
    }
    ```
    This is all you need to make a struct usable with Confetti. Use the following code to load the configs:
    ```go
    import (
        "fmt"
   
        "github.com/shivanshkc/confetti/v2"
    )   

    func main() {
        loader := confetti.NewDefLoader()

        configs := &Configs{}
        if err := loader.Load(configs); err != nil {
            panic(err)
        }

        fmt.Println("Configs:", configs)
    }
    ```
    Here, as you may have guessed:  
    a. The ```Port``` config has the default value of ```8080```.  
    b. If the environment variable ```PORT``` is provided, it will override the default value.  
    c. If the ```-port``` or ```--port``` flag is provided, it will override both the default value and the environment variable.

2. ### Generates documentation for your application
    If you use the code in the last point, and execute ```go run main.go -h```, you will see the following on your console:
    ```
    Usage of configs:
    -port value
            Doc: not provided
            Default: 8080
            Environment: PORT
    panic: failed to parse flags: flag: help requested
    ```
    Confetti auto-generates this help documentation for your application using Go's ```flag``` package.  
    In the output above, notice the ```Doc: not provided``` line. This is because we did not provide any doc on the ```Port``` config. It can be provided as follows:
    ```go
    type Configs struct {
        Port string `def:"8080" env:"PORT" arg:"port,HTTP server port"`
    }
    ```
    Now, the output will read:
    ```
    Usage of configs:
    -port value
        Doc: HTTP server port
        Default: 8080
        Environment: PORT
    panic: failed to parse flags: flag: help requested
    ```
    Next, you must be getting annoyed by the panic message at the bottom. This is because Go's ```flag``` package returns an error when a ```-h``` or ```-help``` flag is provided. To get rid of this, use the following:
    ```go
    import (
        "errors"
        "fmt"
   
        "github.com/shivanshkc/confetti/v2"   
    )

    func main() {
        loader := confetti.NewDefLoader()

        configs := &Configs{}
        if err := loader.Load(configs); err != nil {
            if errors.Is(err, flag.ErrHelp) {
		        return
	        }
            panic(err)
        }

        fmt.Println("Configs:", configs)
    }
    ```  

3. ### Nested structs
    Confetti is built to handle nested structs. Just make sure that the flag names for all fields are always different, otherwise you'll get a panic.
    ```go
    type Configs struct {
        HTTP struct {
            Port string `def:"8080" env:"HTTP_PORT" arg:"http-port"`
        }

        GRPC struct {
            Port string `def:"7070" env:"GRPC_PORT" arg:"grpc-port"`
        }
    }
    ```
    The above struct is a completely valid Confetti target.

4. ### Automatic type assertions
    Confetti is built to handle all types of configs, and not just strings. Consider the following example:
    ```go
    type Configs struct {
        CORS struct {
            TrustedOrigins []string `def:"[\"google.com\"]" env:"TRUSTED_ORIGINS" arg:"trusted-origins"`
        }

        Pagination struct {
            DefaultLimit int `def:"100" env:"DEFAULT_LIMIT" arg:"default-limit"`
        }

        RedisDetails map[string]string `def:"{}" env:"REDIS_DETAILS" arg:"redis-details"`
    }
    ```
    This struct is also a valid Confetti target. Just make sure that the value of the environment variable or flag is a valid JSON string, otherwise Confetti will give you an error.

## Confetti options
Confetti exposes a ```NewLoader``` function and a ```NewDefLoader``` function (as used in the examples above).  
The ```NewDefLoader``` uses the default options, but users can provide their own options by using the ```NewLoader``` function.  
Here is the explanation of all available options:  
| Name       | Description                                              | Default value |
| ---------- | -------------------------------------------------------- | ------------- |
| Title      | The title that shows up on help documentation.           | configs       |
| DefTagName | The name of the tag that controls the default value.     | def           |
| EnvTagName | The name of the tag that controls the env variable name. | env           |
| ArgTagName | The name of the tag that controls the flag name.         | arg           |
| UseDotEnv  | Whether to use the .env file if present.                 | false         |