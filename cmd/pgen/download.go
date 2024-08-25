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

func downloadWithTimeout(url, filename string, timeout time.Duration) ([]byte, error) {
	var (
		body []byte
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
		return nil, fmt.Errorf("invalid url: '%v'", url)
	}

	// Return the response body to the caller
	return body, nil
}

func fetchBeaconState(url, stateRoot string, slot uint64, timeout time.Duration) ([]byte, error) {
	filename := fmt.Sprintf("bstate-%d.%d", slot, time.Now().Unix())
	u2 := url + "/eth/v2/debug/beacon/states/0x" + stateRoot
	data, err := downloadWithTimeout(u2, filename, timeout)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func fetchBeaconBlock(url string, timeout time.Duration) ([]byte, error) {
	u1 := url + "/eth/v2/beacon/blocks/finalized"
	return getBeaconBlock(u1, timeout)
}

// Structure representing the JSON hierarchy up to the "body" field
type BBMessage struct {
	Data struct {
		Message json.RawMessage `json:"message"`
	} `json:"data"`
}

func getBeaconBlock(url string, timeout time.Duration) ([]byte, error) {
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

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return nil, err
	}
	var msg BBMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, err
	}
	return msg.Data.Message, nil
}
