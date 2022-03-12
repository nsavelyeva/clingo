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

Run unit tests:
```
go test ./...
```

#### Weather
Request information about current weather in a certain city, use:
```
./clingo weather --city Amsterdam --token $WEATHER_API_TOKEN
```

#### Currency
Request information about currency rate for the given currency using specified base, execute:
```
./clingo currency --from EUR --to USD,BYN,RUB,PLN --token $CURRENCY_API_TOKEN
```

#### Jokes
Print a short joke, run:
```
./clingo jokes --token $JOKES_API_TOKEN
```

#### News
Read top news, execute:
```
./clingo news --language=nl --sources=rtl-nieuws --from=2022-03-10 --limit=10 --token $NEWS_API_TOKEN
```
_Note._
<br>It seems like providing `pageSize` and `page` parameters in the HTTP GET request
does not have any effect (10 is always a limit).
<br>Hence, it makes sense to keep parameters in the query in case news API will start working
as described at https://newsapi.org/docs/endpoints/top-headlines even for free accounts.
<br>Thus, the limit has the maximum value of 10.

#### Events
See if there is an [upcoming] event according to the provided JSON file, run:
```
./clingo --events events.json
```
The content of `events.json` is as follows:
```
{
  <Format (lines are sorted ascending in calendar year)>
  "<MM(month)>-<DD(day)>": {"year": YYYY, "remind": <integer N or 0>, "type": "<anniversary|birthday|holiday>", "event": "<Description>"},
  ... <Examples> ...
  "01-05": {"year": 2010, "remind": 1, "type": "anniversary", "event": "Someone's anniversary"},
  "02-15": {"year": 2000, "remind": 3, "type": "birthday", "event": "Someone's birthday"},
  "12-25": {"year":    1, "remind": 0, "type": "holiday", "event": "Catholic Christmas Day"},
}
```

## TODO
- Cover functionality with unit tests
- Refactor code:
  - replace `fmt.Printf` with `fmt.Fprint`
- Separate what comes to `stdout` and `stderr`
- Handle exit codes
- Enable logging
