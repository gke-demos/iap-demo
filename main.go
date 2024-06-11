package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/idtoken"
)

func main() {

	args:= os.Args

	if len(args) < 3 {
		log.Fatal("Usage: iap-demo URL")
	}

	url := args[1]

	audience := args[2]
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer

	err = makeIAPRequest(&buf, request, audience)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(&buf)

}

// makeIAPRequest makes a request to an application protected by Identity-Aware
// Proxy with the given audience.
func makeIAPRequest(w io.Writer, request *http.Request, audience string) error {
	ctx := context.Background()

	// self-signed cert
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return fmt.Errorf("idtoken.NewClient: %w", err)
	}


	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer response.Body.Close()
	if _, err := io.Copy(w, response.Body); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	return nil
}
