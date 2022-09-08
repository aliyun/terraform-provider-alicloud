# JSONata in Go

Package jsonata is a query and transformation language for JSON.
It's a Go port of the JavaScript library [JSONata](http://jsonata.org/).

It currently has feature parity with jsonata-js 1.5.4. As well as a most of the functions added in newer versions. You can see potentially missing functions by looking at the [jsonata-js changelog](https://github.com/jsonata-js/jsonata/blob/master/CHANGELOG.md).

## Install

    go get github.com/blues/jsonata-go

## Usage

```Go
import (
	"encoding/json"
	"fmt"
	"log"

	jsonata "github.com/blues/jsonata-go"
)

const jsonString = `
    {
        "orders": [
            {"price": 10, "quantity": 3},
            {"price": 0.5, "quantity": 10},
            {"price": 100, "quantity": 1}
        ]
    }
`

func main() {

	var data interface{}

	// Decode JSON.
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Fatal(err)
	}

	// Create expression.
	e := jsonata.MustCompile("$sum(orders.(price*quantity))")

	// Evaluate.
	res, err := e.Eval(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	// Output: 135
}
```

## JSONata Server
A locally hosted version of [JSONata Exerciser](http://try.jsonata.org/)
for testing is [available here](https://github.com/blues/jsonata-go/jsonata-server).

## JSONata tests
A CLI tool for running jsonata-go against the [JSONata test suite](https://github.com/jsonata-js/jsonata/tree/master/test/test-suite) is [available here](https://github.com/blues/jsonata-go/jsonata-test).



## Contributing

We love issues, fixes, and pull requests from everyone. Please run the
unit-tests, staticcheck, and goimports prior to submitting your PR. By participating in this project, you agree to abide by
the Blues Inc [code of conduct](https://blues.github.io/opensource/code-of-conduct).

For details on contributions we accept and the process for contributing, see our
[contribution guide](CONTRIBUTING.md).

In addition to the Go unit tests there is also a test runner that will run against the jsonata-js test
suite in the [jsonata-test](https://github.com/blues/jsonata-go/jsonata-test) directory. A number of these tests currently fail, but we're working towards feature parity with the jsonata-js reference implementation. Pull requests welcome!

If you would like to contribute to this library a good first issue would be to run the jsonata-test suite,
and fix any of the tests not passing.
