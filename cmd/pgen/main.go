package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	url := flag.String("url", "", "The URL to download the file from")
	path := flag.String("input", "", "The full path to the BeaconState input file")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for the HTTP request")

	// Parse the command line flags
	flag.Parse()

	if *url == "" && *path == "" {
		log.Fatal().Msg("You need to pass either a download URL or a path to the BeaconState")
	}
	if *url != "" && *path != "" {
		log.Fatal().Msg("Please specify either a download URL XOR a path to the BeaconState, no both")
	}
	filename := fmt.Sprintf("fbs.%d", time.Now().Unix())

	var (
		data []byte
		err  error
	)
	if path != nil {
		data, err = downloadFileWithTimeout(*path, filename, *timeout)
	} else {
		data, err = downloadFileWithTimeout(*url, filename, *timeout)
	}
	if err != nil {
		log.Error().Msg(err.Error())
	}

	bst := deneb.BeaconState{}
	err = bst.UnmarshalSSZ(data)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Info().Msgf("slot: %d", bst.LatestBlockHeader.Slot)
	log.Info().Msgf("ParentRoot: %v", hex.EncodeToString(bst.LatestBlockHeader.ParentRoot[:]))

	root, err := bst.GetTree()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	proof, err := root.Prove(40)
	if err != nil || proof == nil {
		log.Error().Msg(err.Error())
	}

	json, err := toJSON(*proof)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	log.Info().Msg(string(json))
}