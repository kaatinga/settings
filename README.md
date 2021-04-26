# env_loader

The package looks up necessary environment variables and use them to set settings for application.

The settings must be formed as struct with byte and string fields.

## Contents

1. [Installation](#installation)
2. [Example](#example)
3. [Limitations](#limits)

<a name="installation"></a>

## Installation

Use go get.

	go get github.com/kaatinga/env_loader

Then import the validator package into your own code.

	import "github.com/kaatinga/env_loader"

<a name="example"></a>

## Example

```go
...
type EnvironmentSettings struct {
Port       string `cfg:"PORT validate:"numeric"`
Database   string `cfg:"DATABASE"`
CacheSize  byte `cfg:"CACHE_SIZE"`
LaunchMode string `cfg:"LAUNCH_MODE"`
}

err := LoadUsingReflect(&EnvironmentSettings)
if err != nil {
    return err
}

...
```

<a name="limits"></a>

## Limitations

The configuration model has some limitations in the way how it is arranged.

First of all, the the wrapped structs must be pointed out via pointer. 

```go
...

type Model2 struct {
    CacheSize   byte `cfg:"CACHE_SIZE"`
}

type Model1 struct {
    Port        string `cfg:"PORT validate:"numeric"`
    Database    string `cfg:"DATABASE"`
    Model2      *Model2
}

...
```

The root model must be also added to the LoadUsingReflect() signature as pointer:

```go

...

err := LoadUsingReflect(&EnvironmentSettings)
if err != nil {
    return err
}

...
```

Otherwise, the function return an error.