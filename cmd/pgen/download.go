package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Data represents the structure of the JSON data inside "data"
type Data struct {
	Root      string `json:"root"`
	Canonical bool   `json:"canonical"`
	Header    Header `json:"header"`
}

// Header represents the structure of the JSON data inside "header"
type Header struct {
	Message   Message `json:"message"`
	Signature string  `json:"signature"`
}

// Message represents the structure of the JSON data inside "message"
type Message struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

// BeaconHeader represents the full JSON response structure
type BeaconHeader struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                Data `json:"data"`
}

// getBeaconHeader fetches JSON from a URL, parses it, and returns the data
func getBeaconHeader(url string, timeout time.Duration) (*BeaconHeader, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create HTTP request")
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute HTTP request")
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		log.Error().Int("StatusCode", resp.StatusCode).Msg("Received non-OK HTTP status")
		return nil, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var bh BeaconHeader
	err = json.NewDecoder(resp.Body).Decode(&bh)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode JSON response")
		return nil, err
	}

	return &bh, nil
}

func downloadFileWithTimeout(url, filename string, timeout time.Duration) ([]byte, error) {
	var (
		body []byte
		err  error
	)

	if strings.HasPrefix(url, "http") {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Create the HTTP request with context
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create HTTP request")
			return nil, err
		}

		// Set the required header
		req.Header.Set("Accept", "application/octet-stream;q=1.0,application/json;q=0.")

		// Perform the HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to execute HTTP request")
			return nil, err
		}
		defer resp.Body.Close()

		// Check the HTTP status code
		if resp.StatusCode != http.StatusOK {
			log.Error().Int("StatusCode", resp.StatusCode).Msg("Received non-OK HTTP status")
			return nil, err
		}

		// Read the response body
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read response body")
			return nil, err
		}
		// Write the response body to a file
		file, err := os.Create(filename)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create file")
			return nil, err
		}
		defer file.Close()

		_, err = file.Write(body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to write to file")
			return nil, err
		}

		log.Info().Str("Filename", filename).Msg("File successfully written")
	} else {
		body, err = os.ReadFile(url)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read input file")
			return nil, err
		}
	}

	// Return the response body to the caller
	return body, nil
}
