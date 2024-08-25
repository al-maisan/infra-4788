# what is this?

This is demo software that shows the construction of merkle proofs for a single ehereum [`BeaconState`](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconstate) property: "`finalized_checkpoint.root`".

# how to build it

```
go mod tidy
make
```

The following command displays a short help text:
```
$ bin/pgen -h
proof generator - generates proofs for ethereum BeaconState data.
Usage:
  -blockp string
    	The full path to the BeaconBlock json file
  -statep string
    	The full path to the BeaconState SSZ-snappy file
  -timeout duration
    	Timeout for the HTTP request (default 10s)
  -url string
    	The URL to download the BeaconState data from (default "https://docs-demo.quiknode.pro/")
```

# how to run it

## run by downloading the most recent _finalized_ beacon chain data

Just run `bin/pgen` to download the most recent finalized beacon block/state and it will compute the proof for it.


Example:
```
$ bin/pgen
2024-08-25T15:49:12+02:00 INF BeaconBlock, slot: 9814080, parent_root: d105e75cf23641d5fd725e8028a89b683b37b66905566413300591cc1a882a4f, state_root: 7a5a297644b332b413691ec46b46329addfd202151686a4533d72b033f7fa2db
2024-08-25T15:49:16+02:00 INF File successfully written Filename=bstate-9814080.1724593752
2024-08-25T15:49:16+02:00 INF BeaconState, slot: 9814080, parent_root: d105e75cf23641d5fd725e8028a89b683b37b66905566413300591cc1a882a4f
```

Please note: this writes 2 files:
- [BeaconBlock](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconblock) data in `json` format to `bblock.9814080`
- [BeaconState](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconstate) data in `SSZ-snappy` format to `bstate-9814080.1724593752`

You can use these files to re-run the program on the same block.

The output produced looks as follows:
```json
{
  "slot": 9814080,
  "beacon_block_root": "320c55726160c7fad09077dcae4d6df1eb11227fe374c3d305594b5e3e279204",
  "beacon_state_root": "7a5a297644b332b413691ec46b46329addfd202151686a4533d72b033f7fa2db",
  "finalized_root": "5523029ad957d5049a7d2336d3a9cb31bbb20612aa05ffbedd9a6f14e944046f",
  "block_time": 1724592983,
  "index": 745,
  "leaf": "5523029ad957d5049a7d2336d3a9cb31bbb20612aa05ffbedd9a6f14e944046f",
  "hashes": [
    "00ae040000000000000000000000000000000000000000000000000000000000",
    "fd23d5277e46023ed819d9b25713dfb88ada62452e2a42c22d4500176a36901f",
    "21c0b67dcc28f8b0bd2fb1dbef45710f42c1613b9d0db3bc524e4a8e95c660fa",
    "b17b8a68401c41a894a5736755fafebcbf43860f95d838b15df6d283a7cc8e9d",
    "e8144504535c4b0315705cb0d24c1ac200ea0b1f77b60e5f7fcb2a4634d8d00b",
    "6791b76ab2bdd75660c63314765f4f0878b67fce5d0ff3cd3f29a35b6f94f275",
    "d105e75cf23641d5fd725e8028a89b683b37b66905566413300591cc1a882a4f",
    "cda7207fa2a3fdabbe5f0328dc2be26359b7ba23bbdd76ed1e9cb94142058593",
    "ef14734d032b181ae448830440035bddfe9e09618d74a2f4ee3b11ace5dc1bf9"
  ]
}
```

The bottom three properties (index, leaf, hashes) constitute the proof and are needed for validation along with `block_time` and `beacon_block_root`.

## run using block data on the local file system

Example:
```
$ bin/pgen -blockp bblock.9814080 -statep bstate-9814080.1724593752|jq
2024-08-25T15:54:24+02:00 INF BeaconBlock, slot: 9814080, parent_root: d105e75cf23641d5fd725e8028a89b683b37b66905566413300591cc1a882a4f, state_root: 7a5a297644b332b413691ec46b46329addfd202151686a4533d72b033f7fa2db
2024-08-25T15:54:25+02:00 INF BeaconState, slot: 9814080, parent_root: d105e75cf23641d5fd725e8028a89b683b37b66905566413300591cc1a882a4f
```

The generated `json` is the same as shown above.

# proof generation

This demo software generates a proof for the [BeaconState](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconstate) `finalized_checkpoint.root` property.

It does this by
- fetching the most recent finalized [BeaconBlock](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconblock) and the respective [BeaconState](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconstate)
- unmarshaling both the BeaconBlock and the BeaconState
- converting them to merkle trees for the purpose of proof generation
- grafting the BeaconState subtree onto the BeaconBlock tree in order to gain a full merkle tree whose proofs culminate in the BeaconBlock root
- generating the actual proof
- printing it in `json` format

Please note: I have [slightly extended](https://github.com/al-maisan/fastssz/commit/b4be06eccf42bee8badb0b430752aaed26981858) the [github.com/ferranbt/fastssz](https://github.com/ferranbt/fastssz) repository to enable the tree grafting.
