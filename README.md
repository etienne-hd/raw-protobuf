# raw-protobuf

![Go version](https://img.shields.io/github/go-mod/go-version/etienne-hd/raw-protobuf?style=for-the-badge)
[![GitHub license](https://img.shields.io/github/license/etienne-hd/raw-protobuf?style=for-the-badge)](https://github.com/etienne-hd/raw-protobuf/blob/master/LICENSE)

**Working with raw Protobuf has never been easier!**

```go
type Person struct {
    Name string `num:"1"`
    Age  uint32 `num:"2"`
}

person := Person{
    Name: "Etienne",
    Age:  20,
}

data, _ := protobuf.Marshal(person)

fmt.Printf("raw protobuf: %00x\n", data)
```

## Installation

Requires **Go 1.23+**.

```bash
go get github.com/etienne-hd/raw-protobuf
```

## Usage

### Marshal

`Marshal` works similarly to `json.Marshal`: define a struct, add field tags, and you're ready to go.

The `raw-protobuf` field tag is `num`.

```go
type A struct {
    B string `num:"1"`
}

data, _ := protobuf.Marshal(A{
    B: "Hello, World!",
})
```

Raw Protobuf messages can become deeply nested, especially when working with large or complex structures. Instead of defining every nested struct, you can describe the entire field path directly using multiple numbers separated by `.`.

```go
type A struct {
    B string `num:"1.2.3.4.5.6.7.8.9"`
}

data, _ := protobuf.Marshal(A{
    B: "Hello, World!",
})
```

`raw-protobuf` automatically builds the corresponding nested Protobuf structure for you.

### Unmarshal

`Unmarshal` works similarly to `json.Unmarshal`: provide the raw Protobuf data and a destination struct.

```go
type Person struct {
    Name string `num:"1"`
    Age  uint32 `num:"2"`
}

var person Person

err := protobuf.Unmarshal(data, &person)
if err != nil {
    panic(err)
}

fmt.Println(person.Name) // Etienne
fmt.Println(person.Age)  // 20
```

### gRPC

`raw-protobuf` also provides utilities for working with **gRPC frames**, including **Gzip compression**.

The gRPC utilities are available in the `protobuf/grpc` package:

```go
import "github.com/etienne-hd/raw-protobuf/protobuf/grpc"
```

A gRPC message consists of a 5-byte header followed by the payload:

- **1 byte**: compression flag
- **4 bytes**: payload length
- **N bytes**: Protobuf payload

#### Frame

Use `Frame` to create a gRPC frame from a raw Protobuf payload.

```go
data, _ := protobuf.Marshal(Person{
    Name: "Etienne",
    Age:  20,
})

frame, err := grpc.Frame(data, grpc.CompressionNone)
if err != nil {
    panic(err)
}
```

Gzip compression is supported through `CompressionGzip`:

```go
frame, err := grpc.Frame(data, grpc.CompressionGzip)
if err != nil {
    panic(err)
}
```

#### Unframe

Use `Unframe` to extract and decompress the Protobuf payload from a gRPC frame.

```go
data, err := grpc.Unframe(frame, grpc.CompressionGzip)
if err != nil {
    panic(err)
}

var person Person

err = protobuf.Unmarshal(data, &person)
if err != nil {
    panic(err)
}

fmt.Println(person.Name) // Etienne
fmt.Println(person.Age)  // 20
```

Supported compression algorithms:

```go
grpc.CompressionNone
grpc.CompressionGzip
```

## License

This project is licensed under the [MIT License](LICENSE).

## Support

<a href="https://www.buymeacoffee.com/etienneh" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;"></a>

You can contact me via [Telegram](https://t.me/etienne_hd) or [Discord](https://discord.com/users/1153975318990827552) if you need help with scraping services or want to write a library.
