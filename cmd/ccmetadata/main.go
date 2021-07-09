package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: ccmetadata -cert=<path> -key=<path> -mspid=<name> -connection=<path> -channel=<name> [-aslocalhost] [-verbose] <chaincode>\n\nGet metadata for the specified chaincode name\n\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	certPath := flag.String("cert", "", "certificate file")
	keyPath := flag.String("key", "", "private key file")
	mspid := flag.String("mspid", "", "private key file, e.g. Org1MSP")
	ccpPath := flag.String("connection", "", "connection profile file")
	channelName := flag.String("channel", "", "channel name, e.g. mychannel")
	aslocalhost := flag.Bool("aslocalhost", false, "use discovery service as localhost")
	verbose := flag.Bool("verbose", false, "enable verbose logging")

	flag.Parse()

	required := []string{"cert", "key", "mspid", "connection", "channel"}
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
	chaincodeName := flag.Arg(0)

	if *verbose {
		fmt.Printf("Certificate file: %s\n", *certPath)
		fmt.Printf("Private key file: %s\n", *keyPath)
		fmt.Printf("MSP ID: %s\n", *mspid)
		fmt.Printf("Channel name: %s\n", *channelName)
		fmt.Printf("As localhost option: %t\n", *aslocalhost)
		fmt.Printf("Chaincode name: %s\n", chaincodeName)
	} else {
		logging.SetLevel("", logging.ERROR)
	}

	if *aslocalhost {
		err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
		if err != nil {
			log.Fatalf("Failed to set DISCOVERY_AS_LOCALHOST environment variable: %v", err)
		}
	}

	wallet, err := createWallet(*certPath, *keyPath, *mspid)
	if err != nil {
		log.Fatalf("Failed to get credentials: %v", err)
	}

	connectionConfig := config.FromFile(filepath.Clean(*ccpPath))

	gateway, err := gateway.Connect(
		gateway.WithConfig(connectionConfig),
		gateway.WithIdentity(wallet, "identity"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gateway.Close()

	network, err := gateway.GetNetwork(*channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract(chaincodeName)

	result, err := contract.EvaluateTransaction("org.hyperledger.fabric:GetMetadata")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	fmt.Println(string(result))
}

func createWallet(certPath, keyPath, mspid string) (*gateway.Wallet, error) {
	wallet := gateway.NewInMemoryWallet()

	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return wallet, err
	}

	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return wallet, err
	}

	identity := gateway.NewX509Identity(mspid, string(cert), string(key))
	wallet.Put("identity", identity)

	return wallet, nil
}
