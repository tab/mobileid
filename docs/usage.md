# Usage

## Configure a client

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

Check client default configuration values:

```go
const (
  Text       = "Enter PIN1"
  TextFormat = "GSM-7"
  Language   = "ENG"
  Timeout    = requests.Timeout
  URL        = "https://tsp.demo.sk.ee/mid-api"
)

cfg := &config.Config{
  HashType:   utils.HashTypeSHA512,
  Text:       Text,
  TextFormat: TextFormat,
  Language:   Language,
  URL:        URL,
  Timeout:    Timeout,
}
```

## Start authentication

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

## Fetch authentication session

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

## Async example

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

## Certificate pinning (optional)

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
  manager, err := mobileid.NewCertificateManager("./certs")
  if err != nil {
    fmt.Println("Failed to create certificate manager:", err)
  }
  tlsConfig := manager.TLSConfig()

  client := mobileid.NewClient().
    WithRelyingPartyName("DEMO").
    WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
    WithHashType("SHA512").
    WithText("Enter PIN1").
    WithTextFormat("GSM-7").
    WithLanguage("ENG").
    WithURL("https://tsp.demo.sk.ee/mid-api").
    WithTimeout(60 * time.Second).
    WithTLSConfig(tlsConfig)

  // Further processing...
```

```sh
Session created: &{b2769811-16de-42f1-a06e-c580d07c1298 5960}
Fetched person: &{PNOEE-60001017869 60001017869 EID2016 TESTNUMBER}
```

- **b2769811-16de-42f1-a06e-c580d07c1298** – session id
- **5960** – verification code


- **PNOEE-60001017869** – formatted identity
- **60001017869** – personal identification code
- **EID2016** – person first name
- **TESTNUMBER** – person last name

