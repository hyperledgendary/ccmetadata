# ccmetadata

Get Hyperledger Fabric chaincode metadata

## Install

Download the latest prebuilt binary for your system and place it in your PATH:

- [ccmetadata-Darwin-x86_64.tgz](https://github.com/hyperledgendary/ccmetadata/releases/latest/download/ccmetadata-Darwin-x86_64.tgz)
- [ccmetadata-Linux-x86_64.tgz](https://github.com/hyperledgendary/ccmetadata/releases/latest/download/ccmetadata-Linux-x86_64.tgz)

Alternatively if there is not a suitable prebuilt binary, clone this repository and build it as follows:

```
go install ./...
```

## Usage

```
Usage: ccmetadata -cert=<path> -key=<path> -mspid=<name> -connection=<path> -channel=<name> [-aslocalhost] [-verbose] <chaincode>

Get metadata for the specified chaincode name

  -aslocalhost
        use discovery service as localhost
  -cert string
        certificate file
  -channel string
        channel name, e.g. mychannel
  -connection string
        connection profile file
  -key string
        private key file
  -mspid string
        private key file, e.g. Org1MSP
  -verbose
        enable verbose logging
```

For example, in the _fabric-samples/test-network_ directory:

```
ccmetadata -cert=organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem \
-key=organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk \
-mspid=Org1MSP \
-connection=organizations/peerOrganizations/org1.example.com/connection-org1.yaml \
-channel=mychannel \
-aslocalhost \
basic
```
