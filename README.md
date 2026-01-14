# Go HTTP Client Assignment

A small Go module that performs an HTTP GET request to `http://example.com/cities/Bangalore` and processes the response. This project demonstrates testability, HTTP mocking, and CI/CD setup.

## Features

- HTTP GET request implementation with dependency injection
- Response processing and error handling
- Comprehensive unit tests using Go's testing framework
- HTTP calls mocked using [httpmatter](https://github.com/therewardstore/httpmatter)
- GitHub Actions CI/CD pipeline
- Clean, testable code structure
- JSON Placeholder style response format

## Project Structure

```
go-http-client/
├── go.mod
├── go.sum
├── README.md
├── .gitignore
├── .github/
│   └── workflows/
│       └── ci.yml
├── client.go          # Main implementation
└── client_test.go     # Unit tests
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-http-client
```

2. Install dependencies:
```bash
go mod download
```

## Usage

The module provides a `FetchData` function that fetches data from the endpoint:

```go
package main

import (
    "fmt"
    "github.com/bharath/go-http-client"
)

func main() {
    client := client.NewDefaultClient()
    data, err := client.FetchData(client)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Response: %s\n", string(data))
}
```

## Running Tests Locally

To run all tests:

```bash
go test ./...
```

To run tests with verbose output:

```bash
go test -v ./...
```

To run tests with coverage:

```bash
go test -cover ./...
```

## Test Coverage

The tests cover:
- ✅ Successful HTTP GET requests with JSON response
- ✅ HTTP error responses (404, 500, 503)
- ✅ Network errors
- ✅ Empty response bodies
- ✅ Malformed JSON responses
- ✅ Response processing

All HTTP calls are mocked using `httpmatter`, ensuring no real network calls are made during testing.

## CI/CD

The project includes a GitHub Actions workflow that:
- Runs on every push and pull request
- Executes `go test ./...`
- Fails the build if tests fail
- Shows test coverage

![CI Status](https://github.com/bharath/go-http-client/actions/workflows/ci.yml/badge.svg)

## Requirements Met

- ✅ Go module (no main() required)
- ✅ HTTP GET request to `http://example.com/cities/Bangalore`
- ✅ Response processing
- ✅ Easily testable code structure
- ✅ Unit tests using Go's testing framework
- ✅ HTTP calls mocked using httpmatter
- ✅ Tests avoid real network calls
- ✅ Tests validate successful and failure responses
- ✅ GitHub Actions CI workflow
- ✅ Clean and readable project structure
- ✅ Table-driven tests
- ✅ Meaningful error handling
- ✅ Coverage for edge and failure cases

## Response Format

The endpoint returns JSON data in a format similar to JSON Placeholder API:

```json
{
  "id": 1,
  "name": "Bangalore",
  "country": "India",
  "population": 8443675,
  "state": "Karnataka"
}
```

## License

This project is created for a coding assignment.
