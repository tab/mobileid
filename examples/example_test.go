//go:build example
// +build example

package examples

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tab/mobileid"
)

func Example_CreateSession() {
	client := mobileid.NewClient().
		WithRelyingPartyName("DEMO").
		WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
		WithHashType("SHA512").
		WithText("Enter PIN1").
		WithTextFormat("GSM-7").
		WithLanguage("ENG").
		WithURL("https://tsp.demo.sk.ee/mid-api").
		WithTimeout(60 * time.Second)

	ctx := context.Background()

	phoneNumber := "+37269930366"
	identity := "51307149560"

	session, err := client.CreateSession(ctx, phoneNumber, identity)
	if err != nil {
		fmt.Println("Error creating session:", err)
	}
	fmt.Println("Session created:", session)
}

func Example_FetchSession() {
	client := mobileid.NewClient().
		WithRelyingPartyName("DEMO").
		WithRelyingPartyUUID("00000000-0000-0000-0000-000000000000").
		WithHashType("SHA512").
		WithText("Enter PIN1").
		WithTextFormat("GSM-7").
		WithLanguage("ENG").
		WithURL("https://tsp.demo.sk.ee/mid-api").
		WithTimeout(60 * time.Second)

	ctx := context.Background()

	sessionId := "d3b8f7c3-7e0c-4a4e-9e6b-4b0e6b8e4f4c"

	person, err := client.FetchSession(ctx, sessionId)
	if err != nil {
		fmt.Println("Error fetching session:", err)
		return
	}

	fmt.Println("Person:", person)
}

func Example_ProcessMultipleIdentitiesInBackground() {
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
			continue
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

func Example_ProcessMultipleIdentitiesInBackground_WithTLS() {
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
