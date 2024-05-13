# slangroom-exec

The missing slangroom executor. We are working the .wasm transpile of the
slangroom, but in the meantime this repo could be used similar to the `zencode-exec`
to embed [https://dyne.org/slangroom](slangroom) into other languages.

`slangroom-exec` is a simple utility that reads from STDIN the following content

1. conf
1. slangroom-contract
1. data
1. keys
1. extra
1. context

separated each per new-line and encoded in `base64` and outputs the slangroom execution to stoud.

### Demo

![Slangroom-exec Demo](./docs/slangroom-exec.gif)

## SLangroom-EXec Format Encoder

This script is used to encode the format of the slangroom-exec command into a string that can be used in the slangroom-exec command.

The script accepts the six parameters that are used in the slangroom-exec command and encodes them into a string. The encoded string is then printed to stdout.

### Usage

For each of the parameters, the script also has option flags:

-c or --conf for conf
-s or --slangroom-contract for slangroom-contract
-d or --data for data
-k or --keys for keys
-e or --extra for extra
-x or --context for context
-F or --filename lookup files based on a prefix
-h or --help to print the help message

#### The named convention `-F` option flag

When you have a suite of files if you follow the formal slangroom name convention as such:

conf: `${prefix}.conf`
slangroom-contract: `${prefix}.slang`
data: `${prefix}.data.json`
keys: `${prefix}.keys.json`
extra: `${prefix}.extra.json`
context: `${prefix}.context`

you can just run

#### STDIN

if you just pass something in `/dev/stdin` is interpreted as the contract.
This also overwrites the `--slangroom-contract` option flag if passed as a duplicate.

## Examples

To encode a slangroom-contract and data, you can run:

To encode parameters from a file, you can run:

Where myParameters.txt is a file containing your parameters.
