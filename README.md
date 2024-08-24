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
  -input string
    	The full path to the BeaconState input file
  -timeout duration
    	Timeout for the HTTP request (default 10s)
  -url string
    	The URL to download the BeaconState data from (default "https://docs-demo.quiknode.pro/")
```

# how to run it

## run by downloading the most recent _finalized_ beacon state

Just run `bin/pgen` to download the most recent finalized beacon state and it will compute the proof for it.


Example:
```
$ bin/pgen
2024-08-24T19:00:27+02:00 INF slot: 9807808, state root: 0x5875ec5bab5d9bf89b6e29962e15e53739ca34a4822e4e0068d82ac2bab211fd
2024-08-24T19:00:31+02:00 INF File successfully written Filename=fbs-9807808.1724518827
2024-08-24T19:00:31+02:00 INF slot: 9807808
2024-08-24T19:00:31+02:00 INF ParentRoot: 073810e1f6f095e0d8f1a78f20c23c063b3c47e12e6138a93032ef144d4ceba4
2024-08-24T19:00:41+02:00 INF state root: '5875ec5bab5d9bf89b6e29962e15e53739ca34a4822e4e0068d82ac2bab211fd'
```

As you can see the beacon state was also written to a file (in SSZ-snappy format) so you can re-run the program for it again.

The output produced looks as follows:
```json
{
  "slot": 9807808,
  "parent_root": "073810e1f6f095e0d8f1a78f20c23c063b3c47e12e6138a93032ef144d4ceba4",
  "state_root": "5875ec5bab5d9bf89b6e29962e15e53739ca34a4822e4e0068d82ac2bab211fd",
  "block_time": 1724517719,
  "index": 105,
  "leaf": "2a001e07d57e5e48da758ef6e815042f97e8565941c58983e1c3a3e270a7884e",
  "hashes": [
    "3cad040000000000000000000000000000000000000000000000000000000000",
    "8c4dd08fd8b5d27c2329b738550630a503150b12c5859ca48987e4ca598b140e",
    "261f2f3c6566c808887f617f1942f39fdb6d179044fb0a1fcf6bb3b8f96b9463",
    "c94ea64c6b9a2808ca0e07f32eb5176ce9c49669ab8f4e8e5e4cf64e696fba85",
    "319781ab3ea888a99904bde940aab3fa4c041e7745a9980810ad3a390cee26c6",
    "a11922ead3e2eee1ccc4b2fa91e7e6bea3a2e3a8e3aa1dbe3d0c4969e86322e8"
  ]
}
```

The bottom three properties (index, leaf, hashes) constitute the proof and are needed for validation.

## run using a beacon state stored in a local (SSZ-snappy) file

Example:
```
$ bin/pgen -input fbs-9807808.1724518827
2024-08-24T19:04:52+02:00 INF slot: 9807808
2024-08-24T19:04:52+02:00 INF ParentRoot: 073810e1f6f095e0d8f1a78f20c23c063b3c47e12e6138a93032ef144d4ceba4
2024-08-24T19:05:03+02:00 INF state root: '5875ec5bab5d9bf89b6e29962e15e53739ca34a4822e4e0068d82ac2bab211fd'
```

The generated `json` is the same as shown above.
