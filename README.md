# Go Mobile-ID client

Golang client for the Mobile-ID API (https://www.mobile-id.lt/en/).
It is a simple wrapper around the API, which helps easily integrate Mobile-ID authentication into Golang applications.

## Installation

Use `go get` to install the package

```sh
go get github.com/tab/mobileid
```

## Usage

### Creating a Client

Create a new client using `NewClient()` and customize its configuration using chainable methods.

```go
package main

import (
  "context"
  "log"
  "time"

  "github.com/tab/mobileid"
)

func main() {
  client := mobileid.NewClient().
    WithRelyingPartyName("DEMO").
    WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
    WithHashType("SHA512").
    WithText("Enter PIN1").
    WithTextFormat("GSM-7").
    WithLanguage("ENG").
    WithURL("https://tsp.demo.sk.ee/mid-api").
    WithTimeout(60 * time.Second)

  if err := client.Validate(); err != nil {
    log.Fatal("Invalid configuration:", err)
  }

  // Further processing...
}
```

### Start Authentication

Initiate a new authentication session with the `Mobile-ID` provider by calling `CreateSession`.
This function generates a random hash, constructs the session request, and returns a session that includes an identifier and a verification code.

```go
func main() {
  // Create a client...


  phoneNumber := "+37268000769"
  identity := "60001017869"

  session, err := client.CreateSession(context.Background(), phoneNumber, identity)
  if err != nil {
    log.Fatal("Error creating session:", err)
  }

  fmt.Println("Session created:", session)
}
```

### Fetch Session

```go
func main() {
  // Create a client...

  person, err := client.FetchSession(context.Background(), sessionId)
  if err != nil {
    log.Fatal("Error fetching session:", err)
  }

  fmt.Println("Session status:", session.State)
}
```

### Async Example

For applications requiring the processing of multiple authentication sessions simultaneously, `Mobile-ID` provides a worker model.
Create a worker using `NewWorker`, configure its concurrency and queue size, and then start processing.

```go
package main

import (
  "context"
  "fmt"
  "sync"
  "time"

  "github.com/tab/mobileid"
)

func main() {
  client := mobileid.NewClient().
    WithRelyingPartyName("DEMO").
    WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
    WithHashType("SHA512").
    WithText("Enter PIN1").
    WithTextFormat("GSM-7").
    WithLanguage("ENG").
    WithURL("https://tsp.demo.sk.ee/mid-api").
    WithTimeout(60 * time.Second)

  identities := map[string]string{
    "51307149560": "+37269930366",
    "60001017869": "+37268000769",
    "60001018800": "+37200000566",
    "60001019939": "+37200000266",
    "60001019947": "+37207110066",
    "60001019950": "+37201100266",
    "60001019961": "+37200000666",
    "60001019972": "+37201200266",
    "60001019983": "+37213100266",
    "50001018908": "+37266000266",
  }

  worker := mobileid.NewWorker(client).
    WithConcurrency(50).
    WithQueueSize(100)

  ctx := context.Background()

  worker.Start(ctx)
  defer worker.Stop()

  var wg sync.WaitGroup

  for identity, phoneNumber := range identities {
    wg.Add(1)

    session, err := client.CreateSession(ctx, phoneNumber, identity)
    if err != nil {
      fmt.Println("Error creating session:", err)
      wg.Done()
    }
    fmt.Println("Session created:", session)

    resultCh := worker.Process(ctx, session.Id)
    go func() {
      defer wg.Done()
      result := <-resultCh
      if result.Err != nil {
        fmt.Println("Error fetching session:", result.Err)
      } else {
        fmt.Println("Fetched person:", result.Person)
      }
    }()
  }

  wg.Wait()
}
```

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- [SK ID Solutions](https://www.skidsolutions.eu)
- [Mobile-ID](https://www.mobile-id.lt/en)
