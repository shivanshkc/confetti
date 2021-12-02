# Confetti
A minimal configuration manager for Go applications.

## Key features
Confetti provides a really easy way to control your configurations using both environment variables and command-line arguments.  
Consider the following struct:
```go
package main

type Configs struct {
	Port string `def:"3000" env:"PORT" arg:"port"`
	MaxLimit int `def:"100" env:"MAX_LIMIT" arg:"max-limit"`
}
```

Here, the only peculiar thing are the struct tags. Their explanation is as follows. 
#### 1. *def* tag
This tag is used to specify the default value of the field.

#### 2. *env* tag
This tag is used to specify the name of the environment variable that will override the default value.  
If the environment variable is not set, it will be ignored. But if the variable is set, even if it is
an empty string, it will be considered.

#### 3. *arg* tag
This tag is used to specify the name of the command-line flag that will override the default value and the environment variable as well. Even empty string values are considered.

The struct given above can be populated using the following code:
```go
package main

import (
	"fmt"
	
	"github.com/shivanshkc/confetti"
)

type Configs struct {
	Port string `def:"3000" env:"PORT" arg:"port"`
	MaxLimit int `def:"100" env:"MAX_LIMIT" arg:"max-limit"`
}

func main() {
	loader := confetti.NewLoader(nil) // The passed 'nil' means we will get a Loader with default settings.
	
	conf := &Configs{}
	if err := loader.Load(conf); err != nil {
		panic(err)
	}
	
	fmt.Printf("Loaded configs: %+v", conf)
}
```

## To be noted
1. Confetti handles all the type conversions itself. So, you can use custom field types as well. 
2. If your config struct is a nested struct, for example:
    ```go
    package main
    
    type Configs struct {
        Nested struct {
            Value string `def:"OK" env:"VALUE" arg:"value"`
        } `def:"{}" env:"NESTED"`
    }
    ```
    Here, Confetti will first use the tags given to the ```Nested``` field to resolve its value. Then, it will go inside
    the ```Nested``` struct and process the tags given to each field to update their value. If a field does not have tags, it will not be updated.

