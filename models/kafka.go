package models

import (
	"net/url"

	"fmt"

	"net/http"
)

// documentation from upstash kafka.

func ProduceMessage(address, key, message, user, pass string) error {
	// Encode the message for inclusion in the URL
	messageEncoded := url.QueryEscape(message)
	keyEncoded := url.QueryEscape(key)

	// Construct the URL with properly encoded components
	url := fmt.Sprintf("https://%s/produce/broadcast/%s?key=%s", address, messageEncoded, keyEncoded)

	// Create a new HTTP GET request (since you're using "GET" method)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.SetBasicAuth(user, pass)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return fmt.Errorf("non-OK status code received: %s", resp.Status)
	}
	return nil
}
