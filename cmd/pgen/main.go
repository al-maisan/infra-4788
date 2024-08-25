package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const FINALIZED_ROOT_GINDEX = 745

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	url := flag.String("url", "https://docs-demo.quiknode.pro/", "The URL to download the BeaconState data from")
	blockp := flag.String("blockp", "", "The full path to the BeaconBlock json file")
	statep := flag.String("statep", "", "The full path to the BeaconState SSZ-snappy file")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for the HTTP request")

	flag.Usage = func() {
		// Program title and description
		fmt.Fprintf(os.Stderr, "proof generator - generates proofs for ethereum BeaconState data.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults() // Print default flag usage information
	}
	// Parse the command line flags
	flag.Parse()

	if *url == "" && *blockp == "" {
		log.Fatal().Msg("You need to pass either a download URL or a path to a file with BeaconBlock data.")
	}

	var (
		data []byte
		err  error
	)
	if *blockp != "" {
		data, err = os.ReadFile(*blockp)
	} else {
		data, err = fetchBeaconBlock(*url, *timeout)
	}
	if err != nil {
		log.Fatal().Msgf("failed to fetch data, %v", err.Error())
	}

	// BeaconBlock
	// parse
	bblock := deneb.BeaconBlock{}
	err = bblock.UnmarshalJSON(data)
	if err != nil {
		log.Fatal().Msgf("failed to unmarshal BeaconBlock, %v", err.Error())
	}
	bbParentRoot := hex.EncodeToString(bblock.ParentRoot[:])
	bbStateRoot := hex.EncodeToString(bblock.StateRoot[:])
	log.Info().Msgf("BeaconBlock, slot: %v, parent_root: %v, state_root: %v", bblock.Slot, bbParentRoot, bbStateRoot)
	// persist
	err = writeToFile(fmt.Sprintf("bblock.%v", bblock.Slot), data)
	if err != nil {
		log.Fatal().Msgf("failed to write BeaconBlock data, %v", err.Error())
	}
	// get tree
	bbRoot, err := bblock.GetTree()
	if err != nil {
		log.Fatal().Msgf("failed to get tree for BeaconBlock, %v", err.Error())
	}

	// BeaconState
	// read
	if *statep != "" {
		data, err = os.ReadFile(*statep)
	} else {
		data, err = fetchBeaconState(*url, bbStateRoot, uint64(bblock.Slot), *timeout)
	}
	if err != nil {
		log.Fatal().Msgf("failed to read BeaconState data, %v", err.Error())
	}
	// parse
	bstate := deneb.BeaconState{}
	err = bstate.UnmarshalSSZ(data)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	bsParentRoot := hex.EncodeToString(bstate.LatestBlockHeader.ParentRoot[:])
	log.Info().Msgf("BeaconState, slot: %v, parent_root: %v", bstate.Slot, bsParentRoot)
	if bbParentRoot != bsParentRoot {
		log.Fatal().Msg("parent root mismatch")
	}
	// get tree
	bsRoot, err := bstate.GetTree()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// graft the BeaconState subtree onto the BeaconBlock tree
	bbHashBefore := bbRoot.Hash()
	log.Info().Msgf("BeaconBlock, before grafting, hash: %v", hex.EncodeToString(bbHashBefore))
	ok := graftSubtree(bbRoot, bblock.StateRoot[:], bsRoot)
	if !ok {
		log.Fatal().Msg("BeaconBlock tree grafting failed")
	}
	bbHashAfter := bbRoot.Hash()
	log.Info().Msgf("BeaconBlock, after grafting, hash: %v", hex.EncodeToString(bbHashAfter))
	if !bytes.Equal(bbHashBefore, bbHashAfter) {
		log.Fatal().Msg("BeaconBlock tree grafting caused tree root mismatch")
	}

	// generate proof
	proof, err := bbRoot.Prove(FINALIZED_ROOT_GINDEX)
	if err != nil || proof == nil {
		log.Fatal().Msgf("failed to geneate proof, %v", err.Error())
	}

	// convert proof to json
	finalizedRoot := hex.EncodeToString(bstate.FinalizedCheckpoint.Root[:])
	json, err := toJSON(*proof, uint64(bstate.LatestBlockHeader.Slot), hex.EncodeToString(bbHashBefore), bbStateRoot, finalizedRoot)
	if err != nil {
		log.Fatal().Msgf("failed to convert proof to json, %v", err.Error())
	}

	fmt.Println(string(json))
}

func writeToFile(filename string, data []byte) error {
	// Open the file with the necessary flags:
	// O_WRONLY - open the file for writing only
	// O_CREATE - create the file if it doesn't exist
	// O_TRUNC  - truncate the file if it already exists
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the []byte data to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
