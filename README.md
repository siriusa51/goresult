# goresult

goresult is a Golang library that provides Result and Option types, similar to those used in functional programming, to handle errors and optional values in a structured and expressive way.

## Introduction

The goresult library introduces Result and Option types, which are commonly used in functional programming, to simplify error handling and deal with optional values more effectively. These types offer a clean and concise approach to manage success and failure cases, as well as situations where a value may or may not exist.

## Features

- **Result Type**: The Result type represents a value that can either be a success or an error. It enables you to handle potential errors without resorting to traditional error handling techniques, improving code readability and maintainability.

- **Option Type**: The Option type represents an optional value that may or may not be present. It allows you to handle scenarios where a value may be absent, eliminating the need for explicit nil checks and reducing the chance of null pointer errors.

## Installation

To use goresult in your Go project, you can import it using the following command:

```shell
go get github.com/siriusa51/goresult
```

## Usage

Here's a simple example demonstrating the usage of Result and Option types:

```go
package main

import "github.com/siriusa51/goresult"

type KV struct {
	Key   string
	Value string
}

func main() {
	r1 := goresult.Ok(KV{Key: "key", Value: "value"})
	println(r1.IsOk())
	// -> true
	println(r1.IsError())
	// -> false

	r2 := goresult.Error[int]("error")
	println(r2.IsOk())
	// -> false
	println(r2.IsError())
	// -> true

	o := goresult.Some(KV{Key: "key", Value: "value"})
	println(o.IsSome())
	// -> true
	println(o.IsNone())
	// -> false

	n := goresult.None[KV]()
	println(n.IsSome())
	// -> false
	println(n.IsNone())
	// -> true
}
```

For more examples, you can refer to the comments within the code. 

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request on the [GitHub repository](https://github.com/siriusa51/goresult).

## License

This project is licensed under the [MIT License](LICENSE).
