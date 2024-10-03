

### Usage
This Go package provides an interface to execute Slangroom using an embedded binary. To use this package, include it in your Go project by importing it:

```go
import "github.com/dyne/slangroom-exec/bindings/go"


### Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/dyne/slangroom-exec/bindings/go"
)

func main() {
    // Define your contract and inputs
    contract := `Rule unknown ignore
Given I fetch the local timestamp in seconds and output into 'timestamp'
Given I have a 'time' named 'timestamp'
Then print the string '😘 Welcome to the Slangroom World 🌈'
Then print the 'timestamp'`

    // Execute  Slangroom
    result, err := slangroom.SlangroomExec("", contract, "", "", "", "")
    if err != nil {
        log.Fatalf("Execution Failed: %v", err)
    }

    // Print the execution output
    fmt.Println("Execution Output:", result.Output)
}

