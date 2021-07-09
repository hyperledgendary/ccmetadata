# ccmetadata

Get Hyperledger Fabric chaincode metadata

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
ccmetadata -cert=organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem \
-key=organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk \
-mspid=Org1MSP \
-connection=organizations/peerOrganizations/org1.example.com/connection-org1.yaml \
-channel=mychannel \
-aslocalhost \
-verbose \
basic
```
