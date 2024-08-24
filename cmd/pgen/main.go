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

const FINALIZED_ROOT_GINDEX = 105

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	url := flag.String("url", "https://docs-demo.quiknode.pro/", "The URL to download the BeaconState data from")
	path := flag.String("input", "", "The full path to the BeaconState input file")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for the HTTP request")

	flag.Usage = func() {
		// Program title and description
		fmt.Fprintf(os.Stderr, "proof generator - generates proofs for ethereum BeaconState data.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults() // Print default flag usage information
	}
	// Parse the command line flags
	flag.Parse()

	if *url == "" && *path == "" {
		log.Fatal().Msg("You need to pass either a download URL or a path to a file with BeaconState data.")
	}
	if *path != "" {
		*url = ""
	}
	if *url != "" && *path != "" {
		log.Fatal().Msg("Please specify either a download URL XOR a path to a file with the BeaconState data, not both.")
	}

	var (
		data []byte
		err  error
	)
	if *path != "" {
		data, err = os.ReadFile(*path)
	} else {
		data, err = fetchBeaconState(*url, *timeout)
	}
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	bst := deneb.BeaconState{}
	err = bst.UnmarshalSSZ(data)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msgf("slot: %d", bst.LatestBlockHeader.Slot)
	parentRoot := hex.EncodeToString(bst.LatestBlockHeader.ParentRoot[:])
	log.Info().Msgf("ParentRoot: %v", parentRoot)

	root, err := bst.GetTree()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	stateRoot := hex.EncodeToString(root.Hash())
	log.Info().Msgf("state root: '%s'", stateRoot)
	// This is from https://github.com/ethereum/consensus-specs
	// ./tests/core/pyspec/eth2spec/deneb/mainnet.py:FINALIZED_ROOT_GINDEX = GeneralizedIndex(105)
	proof, err := root.Prove(FINALIZED_ROOT_GINDEX)
	if err != nil || proof == nil {
		log.Fatal().Msg(err.Error())
	}

	json, err := toJSON(*proof, uint64(bst.LatestBlockHeader.Slot), parentRoot, stateRoot)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	fmt.Println(string(json))
}
