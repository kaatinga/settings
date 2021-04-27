# env_loader

The package looks up necessary environment variables and use them to set settings for application.

The settings must be formed as struct with byte and string fields.

## Contents

1. [Installation](#installation)
2. [Example](#example)
3. [Limitations](#limits)

<a name="installation"></a>

## 1. Installation

Use go get.

	go get github.com/kaatinga/env_loader

Then import the validator package into your own code.

	import "github.com/kaatinga/env_loader"

<a name="example"></a>

## 2. Example

```go
...
type EnvironmentSettings struct {
    Port       string `env:"PORT validate:"numeric"`
    Database   string `env:"DATABASE"`
    CacheSize  byte `env:"CACHE_SIZE"`
    LaunchMode string `env:"LAUNCH_MODE"`
}

err := LoadUsingReflect(&EnvironmentSettings)
if err != nil {
    return err
}

...
```

<a name="limits"></a>

## 3. Limitations

The configuration model has some limitations in the way how it is arranged.

First of all, the the nested structs must be pointed out via pointer. 

```go
...

type Model2 struct {
    CacheSize   byte `env:"CACHE_SIZE"`
}

type Model1 struct {
    Port        string `env:"PORT validate:"numeric"`
    Database    string `env:"DATABASE"`
    Model2      *Model2
}

...
```

The root model must be also added to the LoadUsingReflect() signature via pointer:

```go

...

err := LoadUsingReflect(&EnvironmentSettings)
if err != nil {
    return err
}

...
```

Otherwise, the function returns error.