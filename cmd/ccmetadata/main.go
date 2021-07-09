package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = usage

	cert := flag.String("cert", "", "certificate file")
	key := flag.String("key", "", "private key file")
	mspid := flag.String("mspid", "", "private key file, e.g. Org1MSP")
	channel := flag.String("channel", "", "channel name, e.g. mychannel")
	aslocalhost := flag.Bool("aslocalhost", false, "use discovery service as localhost")
	verbose := flag.Bool("verbose", false, "enable verbose logging")

	flag.Parse()

	required := []string{"cert", "key", "mspid", "channel"}
	provided := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { provided[f.Name] = true })
	for _, r := range required {
		if !provided[r] {
			fmt.Fprintf(os.Stderr, "flag required but not provided: -%s\n", r)
			usage()
			os.Exit(2)
		}
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "argument required but not provided: chaincode")
		usage()
		os.Exit(2)
	}
	chaincode := flag.Arg(0)

	if *verbose {
		fmt.Printf("Certificate file: %s\n", *cert)
		fmt.Printf("Private key file: %s\n", *key)
		fmt.Printf("MSP ID: %s\n", *mspid)
		fmt.Printf("Channel name: %s\n", *channel)
		fmt.Printf("As localhost option: %t\n", *aslocalhost)
		fmt.Printf("Chaincode name: %s\n", chaincode)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: ccmetadata -cert=<path> -key=<path> -mspid=<name> -channel=<name> [-aslocalhost] <chaincode>\n\nGet metadata for the specified chaincode name\n\n")
	flag.PrintDefaults()
}
