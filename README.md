# Clingo

This project is for learning purpose to study how CLI tools can be written in Go language
using [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper).

Inspired by [this article](https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/)
and [this fork](https://github.com/corverroos/stingoftheviper).

To compile the project and generate the executable file for direct use, simply run `go build -o clingo`
but to make it work in Alpine docker image you will have to execute the following command:
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o clingo main.go
```

To run unit tests:
```
go test ./...
```
