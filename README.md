# Go HTTP Client Assignment

A small Go module that performs an HTTP GET request to `https://jsonplaceholder.typicode.com/posts` (JSON Placeholder API) and processes the response. This project demonstrates testability, HTTP mocking, and CI/CD setup.

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
- ✅ Successful HTTP GET requests with JSON response (array of posts)
- ✅ HTTP error responses (400, 401, 403, 404, 500, 502, 503)
- ✅ Network errors (connection refused, timeout, DNS failure)
- ✅ Empty response bodies
- ✅ Malformed JSON responses
- ✅ Response body read errors
- ✅ Response processing

All HTTP calls are mocked using `httpmatter`, ensuring no real network calls are made during testing.

## CI/CD

The project includes a GitHub Actions workflow that:
- Runs on every push and pull request
- Executes `go test ./...`
- Fails the build if tests fail
- Shows test coverage

![CI Status](https://github.com/beerappabvgp/vananam-assignment/actions/workflows/ci.yml/badge.svg)

## Requirements Met

- ✅ Go module (no main() required)
- ✅ HTTP GET request to `https://jsonplaceholder.typicode.com/posts`
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

The endpoint returns JSON data in JSON Placeholder API format (array of posts):

```json
[
  {
    "userId": 1,
    "id": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
  },
  {
    "userId": 1,
    "id": 2,
    "title": "qui est esse",
    "body": "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque"
  }
]
```

## License

This project is created for a coding assignment.
