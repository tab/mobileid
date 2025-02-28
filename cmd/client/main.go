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
