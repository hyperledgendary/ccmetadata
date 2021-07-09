package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: ccmetadata -cert=<path> -key=<path> -channel=<name> <chaincode>\n\nGet metadata for the specified chaincode name\n\n")
		flag.PrintDefaults()
	}

	cert := flag.String("cert", "", "certificate file")
	key := flag.String("key", "", "private key file")
	channel := flag.String("channel", "", "channel name")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Incorrect Usage: chaincode name required")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// chaincode := flag.Args()[0]
	chaincode := flag.Arg(0)

	// -cert=testcert -key=testkey -channel=testchannel testchaincode
	fmt.Printf("Certificate file: %s\n", *cert)
	fmt.Printf("Private key file: %s\n", *key)
	fmt.Printf("Channel name: %s\n", *channel)
	fmt.Printf("Chaincode name: %s\n", chaincode)
}
