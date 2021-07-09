# ccmetadata

Get Hyperledger Fabric chaincode metadata

```
Usage: ccmetadata -cert=<path> -key=<path> -mspid=<name> -channel=<name> [-aslocalhost] <chaincode>

Get metadata for the specified chaincode name

  -aslocalhost
        use discovery service as localhost
  -cert string
        certificate file
  -channel string
        channel name, e.g. mychannel
  -key string
        private key file
  -mspid string
        private key file, e.g. Org1MSP
  -verbose
        enable verbose logging
```

For example:

```
ccmetadata -cert=testcert -key=testkey -mspid=testmsp -channel=testchannel -aslocalhost -verbose testchaincode
```
