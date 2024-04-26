[![Tests](https://github.com/kaatinga/settings/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kaatinga/settings/actions/workflows/test.yml)
[![GitHub release](https://img.shields.io/github/release/kaatinga/settings.svg)](https://github.com/kaatinga/settings/releases)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/settings/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kaatinga/settings/branch/main/graph/badge.svg)](https://codecov.io/gh/kaatinga/settings)
[![lint workflow](https://github.com/kaatinga/settings/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/settings/actions?query=workflow%3Alinter)
[![help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/settings/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# settings

The package looks up necessary environment variables and uses them to specify settings for your application. In addition, the package validates the final struct using standard `validate` tags.

## Contents

1. [Installation](#1-installation)
2. [Description](#2-description)
3. [Limitations](#3-limitations)

## 1. Installation

Use `go get` to install the package:
```bash
go get github.com/kaatinga/settings
````

Then, import the package into your own code:
```go
import "github.com/kaatinga/settings"
```

## 2. Description

### How to use

Create a settings model where you can use tags `env`, `default` and `validate`. Announce a variable and call `Load()`:

```go
type Settings struct {
    Port       string `env:"PORT" validate:"numeric"`
    Database   string `env:"DATABASE"`
    CacheSize  byte `env:"CACHE_SIZE" default:"50"`
    LaunchMode string `env:"LAUNCH_MODE"`
}

var settings Settings
err := Load(&settings)
if err != nil {
    return err
}
```

The `env` tag must contain the name of the related environment variable.
The `default` tag contains a default value that is used in case the environment variable was not found.
The `validate` tag may contain an optional validation rule fallowing the documentation of the [validator package](https://github.com/go-playground/validator/). 

### Supported types

| Type          | Real type     |
|---------------| ------------- |
| string        | -             | 
| boolean       | -             | 
| ~int          | -             | 
| ~uint         | -             | 
| time.Duration | int64         | 

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
if err := Load(&settings); err != nil {
    return err
}
```

Nonetheless, if you want, you can do it.

```go
var settings = Model1{Model2: new(Model2)}
if err := Load(&settings); err != nil {
    return err
}
```

## 3. Limitations

The configuration model has some limitations in how it is arranged and used.

### Empty structs are not allowed

If you add an empty struct to your configuration model, `Load()` returns error.

### Load() accepts only pointer to your configuration model

The root model must be initialized and added to the `Load()` signature via pointer:

```go
err := Load(&EnvironmentSettings)
if err != nil {
    return err
}
```

Otherwise, the function returns error.
