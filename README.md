[![GitHub release](https://img.shields.io/github/release/kaatinga/env_loader.svg)](https://github.com/kaatinga/env_loader/releases)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/env_loader/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kaatinga/env_loader/branch/main/graph/badge.svg)](https://codecov.io/gh/kaatinga/env_loader)
[![lint workflow](https://github.com/kaatinga/env_loader/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/env_loader/actions?query=workflow%3Alinter)
[![help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/env_loader/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# env_loader

The package looks up necessary environment variables and use them to specify settings for your application. In addition, the
package validates the final struct using standard `validate` tags.

## Contents

1. [Installation](#installation)
2. [Example](#example)
3. [Limitations](#limits)

<a name=installation></a>

## 1. Installation

Use go get.

	go get github.com/kaatinga/env_loader

Then import the validator package into your own code.

	import "github.com/kaatinga/env_loader"

<a name=example></a>

## 2. Description

### How to use

Create a settings model where you can use tags `env`, `default` and `validate`. Announce a variable and call `LoadSettings()`:

```go
type Settings struct {
    Port       string `env:"PORT" validate:"numeric"`
    Database   string `env:"DATABASE"`
    CacheSize  byte `env:"CACHE_SIZE" default:"50"`
    LaunchMode string `env:"LAUNCH_MODE"`
}

var settings Settings
err := LoadSettings(&settings)
if err != nil {
    return err
}
```

The `env` tag must contain the name of the related environment variable.
The `default` tag contains a default value that is used in case the environment variable was not found.
The `validate` tag may contain an optional validation rule fallowing the documentation of the [validator package](https://github.com/go-playground/validator/). 

### Supported types

| Type                   | Real type     |
| -------------          | ------------- |
| string                 | -             | 
| boolean                | -             | 
| any uint               | -             | 
| int, int64             | -             | 
| logrus.Level           | uint32        | 
| syslog.Priority        | int           | 
| time.Duration          | int64         | 

### Nested structs

Nested structs can be added via pointer or without pointer. Example:

```go
type Model2 struct {
    CacheSize   byte `env:"CACHE_SIZE"`
}

type Model3 struct {
    Port        string `env:"PORT validate:"numeric"`
}

type Model1 struct {
    Database    string `env:"DATABASE"`
    Model2      *Model2
    Model3      Model3
}
```

The nested structs that added via pointer must not be necessarily initialized:

```go
var settings Model1
err := LoadSettings(&settings)
if err != nil {
    return err
}
```

Nonetheless, if you want, you can do it.

```go
var settings = Model1{Model2: new(Model2)}
err := LoadSettings(&settings)
if err != nil {
    return err
}
```

<a name=limits></a>

## 3. Limitations

The configuration model has some limitations in the way how it is arranged and used.

### Empty structs are not allowed

If you add an empty struct to your configuration model, `LoadSettings()` returns error.

### LoadSettings() accepts only pointer to your configuration model

The root model must be initialized and added to the `LoadSettings()` signature via pointer:

```go
err := LoadSettings(&EnvironmentSettings)
if err != nil {
    return err
}
```

Otherwise, the function returns error.
