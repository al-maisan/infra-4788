package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	ssz "github.com/ferranbt/fastssz"
)

const genesisTime = 1606824023

type proofJSON struct {
	Slot          uint64   `json:"slot"`
	BlockRoot     string   `json:"beacon_block_root"`
	StateRoot     string   `json:"beacon_state_root"`
	FinalizedRoot string   `json:"finalized_root"`
	BlockTime     uint64   `json:"block_time"`
	Index         int      `json:"index"`
	Leaf          string   `json:"leaf"`
	Hashes        []string `json:"hashes"`
}

// ToJSON converts the Proof struct to a JSON object and returns it as a byte slice
func toJSON(p ssz.Proof, slot uint64, blockRoot, stateRoot, finalizedRoot string) ([]byte, error) {
	// Create an intermediate struct for custom JSON serialization
	intermediate := proofJSON{
		Slot:          slot,
		BlockRoot:     blockRoot,
		StateRoot:     stateRoot,
		FinalizedRoot: finalizedRoot,
		BlockTime:     genesisTime + slot*12,
		Index:         p.Index,
		Leaf:          hex.EncodeToString(p.Leaf),
		Hashes:        make([]string, len(p.Hashes)),
	}

	// Convert each hash in the Hashes slice to a hex string
	for i, hash := range p.Hashes {
		intermediate.Hashes[i] = hex.EncodeToString(hash)
	}

	// Marshal the intermediate struct to a JSON byte slice
	jsonData, err := json.Marshal(intermediate)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Proof struct to JSON: %w", err)
	}
	return jsonData, nil
}
