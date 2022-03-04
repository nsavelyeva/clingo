# Clingo

## Description
This project is for learning purpose to study how CLI tools can be written in Go language
using [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper).

Inspired by [this article](https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/)
and [this fork](https://github.com/corverroos/stingoftheviper).

## How to...
To compile the project and generate the executable file for direct use,
simply run `go build -o clingo`
but to make it work in Alpine docker image you will have to execute the following command:
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o clingo main.go
```

To run unit tests:
```
go test ./...
```

To request information about current weather in a certain city, use:
```
./clingo weather --city Amsterdam --token $WEATHER_API_TOKEN
```

To request information about currency rate for the given currency using specified base, execute:
```
./clingo currency --from EUR --to USD,BYR,RUB,PLN --token $CURRENCY_API_TOKEN
```

To print a short joke, run:
```
./clingo jokes --token $JOKES_API_TOKEN
```

To see if there is an [upcoming] event according to the provided JSON file, run:
```
./clingo --events events.json
```

## TODO
- Cover functionality with unit tests
- Refactor code:
  - replace `fmt.Printf` with `fmt.Fprint`
- Separate what comes to `stdout` and `stderr`
- Handle exit codes
- Enable logging
