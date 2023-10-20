# ccmetadata

Get Hyperledger Fabric chaincode metadata

## Install

Download the latest prebuilt binary for your system and place it in your PATH:

- [ccmetadata-Linux-X64.tgz](https://github.com/hyperledgendary/ccmetadata/releases/latest/download/ccmetadata-Linux-X64.tgz)
- [ccmetadata-macOS-X64.tgz](https://github.com/hyperledgendary/ccmetadata/releases/latest/download/ccmetadata-macOS-X64.tgz)
- [ccmetadata-Windows-X64.tgz](https://github.com/hyperledgendary/ccmetadata/releases/latest/download/ccmetadata-Windows-X64.tgz)

Alternatively if there is not a suitable prebuilt binary, clone this repository and build it as follows:

```
go install ./...
```

## Usage

```
Usage: ccmetadata --gateway=<address> --cert=<path> --key=<path> --mspid=<name> --channel=<name> [--tlscert=<path>] [--override=<hostname>] [--aslocalhost] [--verbose] <chaincode>

Get metadata for the specified chaincode name

  -g, --gateway string
        gateway peer address
  -c, --cert string
        certificate file
  -k, --key string
        private key file
  -m, --mspid string
        membership service provider name, e.g. Org1MSP
  -n, --channel string
        channel name, e.g. mychannel
  -t, --tlscert string
        tls certificate file
  -o, --override string
        server name override
  -l, --aslocalhost
        use discovery service as localhost
  -v, --verbose
        enable verbose logging
```
