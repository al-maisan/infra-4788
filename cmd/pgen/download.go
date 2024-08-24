package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

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
