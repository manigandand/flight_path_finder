## flight path calculator

### Few Assumptions:

- currently we consider the airport has only one destinations. Ex. airport1 -> airport2, airport2 -> airport3, etc...

- Destination Airport: the assumption from the input data is, the destination airport is not a another-source airport from the list.

- Source Airport: is the one reaches the destination airport while going through all the other source airports.

### How To run

```
make run

2022/12/30 21:59:45 Starting server on port: 8080
```

or

```
make run_clean

./golangci_linter.sh ./...
Running golangci-lint run --out-format=colored-line-number --issues-exit-code=2 --path-prefix=api --config=.golangci.yaml ./...
golangci-lint found no issues
go test -v ./... --failfast
?       volumefi        [no test files]
=== RUN   TestFlightPathCalculator
=== RUN   TestFlightPathCalculator/test_find_flight_path_calculator_-_invalid_routes
=== RUN   TestFlightPathCalculator/test_find_flight_path_calculator
--- PASS: TestFlightPathCalculator (0.00s)
    --- PASS: TestFlightPathCalculator/test_find_flight_path_calculator_-_invalid_routes (0.00s)
    --- PASS: TestFlightPathCalculator/test_find_flight_path_calculator (0.00s)
PASS
ok      volumefi/api    (cached)
?       volumefi/config [no test files]
?       volumefi/middleware     [no test files]
go mod tidy
env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o volumefi -ldflags="-X 'main.version='" .
./volumefi
2022/12/30 21:59:45 Starting server on port: 8080
```

### API Examples:

> Valid Request:

```
curl --location --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "paths": [
        ["IND", "EWR"],
        ["SFO", "ATL"],
        ["GSO", "IND"],
        ["ATL", "GSO"]
    ]
}'

200 OK
{
    "result": [
        "SFO",
        "EWR"
    ],
    "flight_path": [
        "SFO",
        "ATL",
        "GSO",
        "IND",
        "EWR"
    ]
}

```

> Invalid Requests

```
curl --location --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "paths": []
}'

400 Bad Request
{
    "error": "flight routes required",
    "status": 400
}
```

```
curl --location --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "paths": [["SFO"]]
}'

400 Bad Request
{
    "error": "invalid flight route [SFO]",
    "status": 400
}
```

```
curl --location --request POST 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "paths": [["ATL", "SFO"], ["SFO", "ATL"]]
}'

400 Bad Request
{
    "error": "can't able to identify the destination airport from the input",
    "status": 400
}
```

---

Problem Statement:

Senior Software Engineer Take-Home Programming Assignment for Golang

Story: There are over 100,000 flights a day, with millions of people and cargo being transferred around the world. With so many people and different carrier/agency groups, it can be hard to track where a person might be. In order to determine the flight path of a person, we must sort through all of their flight records.

Goal: To create a simple microservice API that can help us understand and track how a particular person's flight path may be queried. The API should accept a request that includes a list of flights, which are defined by a source and destination airport code. These flights may not be listed in order and will need to be sorted to find the total flight paths starting and ending airports.

Required JSON structure:
[["SFO", "EWR"]] => ["SFO", "EWR"]
[["ATL", "EWR"], ["SFO", "ATL"]] => ["SFO", "EWR"]
[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]] => ["SFO", "EWR"]

