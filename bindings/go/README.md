

### Usage
This Go package provides an interface to execute Slangroom using an embedded binary. To use this package, include it in your Go project by importing it:

```go
import "github.com/dyne/slangroom-exec/bindings/go"
```

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
//Specify your input inside a SlangInput structure
input := SlangInput{
    Conf:    "", // or simply omit this line if you want it to default to ""
    Contract: contract
    Data:    "", // or simply omit this line if you want it to default to ""
    Keys:    "", // or simply omit this line if you want it to default to ""
    Extra:   "", // or simply omit this line if you want it to default to ""
    Context: "", // or simply omit this line if you want it to default to ""
}

    // Execute  Slangroom
    result, err := slangroom.SlangroomExec()
    if err != nil {
        log.Fatalf("Execution Failed: %v", err)
    }

    // Print the execution output
    fmt.Println("Execution Output:", result.Output)
}
```
