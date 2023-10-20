package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: ccmetadata --gateway=<address> --cert=<path> --key=<path> --mspid=<name> --channel=<name> [--tlscert=<path>] [--override=<hostname>] [--aslocalhost] [--verbose] <chaincode>\n\n")

	fmt.Fprintf(os.Stderr, "Get metadata for the specified chaincode name\n\n")

	fmt.Fprintf(os.Stderr, "  -g, --gateway string\n        gateway peer address\n")
	fmt.Fprintf(os.Stderr, "  -c, --cert string\n        certificate file\n")
	fmt.Fprintf(os.Stderr, "  -k, --key string\n        private key file\n")
	fmt.Fprintf(os.Stderr, "  -m, --mspid string\n        membership service provider name, e.g. Org1MSP\n")
	fmt.Fprintf(os.Stderr, "  -n, --channel string\n        channel name, e.g. mychannel\n")
	fmt.Fprintf(os.Stderr, "  -t, --tlscert string\n        tls certificate file\n")
	fmt.Fprintf(os.Stderr, "  -o, --override string\n        server name override\n")
	fmt.Fprintf(os.Stderr, "  -l, --aslocalhost\n        use discovery service as localhost\n")
	fmt.Fprintf(os.Stderr, "  -v, --verbose\n        enable verbose logging\n")
}

func main() {
	var gatewayAddress string
	var certPath string
	var keyPath string
	var mspid string
	var channelName string
	var tlsCertPath string
	var serverNameOverride string
	var aslocalhost bool
	var verbose bool

	flag.Usage = usage

	flag.StringVar(&gatewayAddress, "gateway", "", "gateway peer address")
	flag.StringVar(&gatewayAddress, "g", "", "gateway peer address")
	flag.StringVar(&certPath, "cert", "", "certificate file")
	flag.StringVar(&certPath, "c", "", "certificate file")
	flag.StringVar(&keyPath, "key", "", "private key file")
	flag.StringVar(&keyPath, "k", "", "private key file")
	flag.StringVar(&mspid, "mspid", "", "membership service provider name, e.g. Org1MSP")
	flag.StringVar(&mspid, "m", "", "membership service provider name, e.g. Org1MSP")
	flag.StringVar(&channelName, "channel", "", "channel name, e.g. mychannel")
	flag.StringVar(&channelName, "n", "", "channel name, e.g. mychannel")
	flag.StringVar(&tlsCertPath, "tlscert", "", "tls certificate file")
	flag.StringVar(&tlsCertPath, "t", "", "tls certificate file")
	flag.StringVar(&serverNameOverride, "override", "", "server name override")
	flag.StringVar(&serverNameOverride, "o", "", "server name override")
	flag.BoolVar(&aslocalhost, "aslocalhost", false, "use discovery service as localhost")
	flag.BoolVar(&aslocalhost, "l", false, "use discovery service as localhost")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")

	flag.Parse()

	// TODO this doesn't handle short options
	// required := []string{"gateway", "cert", "key", "mspid", "channel"}
	// provided := make(map[string]bool)
	// flag.Visit(func(f *flag.Flag) { provided[f.Name] = true })
	// for _, r := range required {
	// 	if !provided[r] {
	// 		fmt.Fprintf(os.Stderr, "flag required but not provided: -%s\n", r)
	// 		usage()
	// 		os.Exit(2)
	// 	}
	// }

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "argument required but not provided: chaincode")
		usage()
		os.Exit(2)
	}
	chaincodeName := flag.Arg(0)

	// TODO validate that the following arguments have been provided...
	// gatewayAddress
	// certPath
	// keyPath
	// mspid
	// channel

	if verbose {
		fmt.Printf("Gateway address: %s\n", gatewayAddress)
		fmt.Printf("Certificate file: %s\n", certPath)
		fmt.Printf("Private key file: %s\n", keyPath)
		fmt.Printf("MSP ID: %s\n", mspid)
		fmt.Printf("Channel name: %s\n", channelName)
		fmt.Printf("TLS certificate file: %s\n", tlsCertPath)
		fmt.Printf("Server name override: %s\n", serverNameOverride)
		fmt.Printf("As localhost option: %t\n", aslocalhost)
		fmt.Printf("Chaincode name: %s\n", chaincodeName)
	}

	if aslocalhost {
		err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
		if err != nil {
			panic(fmt.Errorf("failed to set DISCOVERY_AS_LOCALHOST environment variable: %w", err))
		}
	}

	clientConnection := newGrpcConnection(gatewayAddress, tlsCertPath, serverNameOverride)
	defer clientConnection.Close()

	id := newIdentity(certPath, mspid)
	sign := newSign(keyPath)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(fmt.Errorf("failed to connect to Fabric Gateway: %w", err))
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	evaluateResult, err := contract.EvaluateTransaction("org.hyperledger.fabric:GetMetadata")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)
	fmt.Print(result)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection(gatewayAddress, tlsCertPath, serverNameOverride string) *grpc.ClientConn {
	var connection *grpc.ClientConn
	var err error

	if len(tlsCertPath) == 0 {
		connection, err = grpc.Dial(gatewayAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Errorf("failed to create insecure gRPC connection: %w", err))
		}
	} else {
		certificate := loadCertificate(tlsCertPath)

		certPool := x509.NewCertPool()
		certPool.AddCert(certificate)
		transportCredentials := credentials.NewClientTLSFromCert(certPool, serverNameOverride)

		connection, err = grpc.Dial(gatewayAddress, grpc.WithTransportCredentials(transportCredentials))
		if err != nil {
			panic(fmt.Errorf("failed to create secure gRPC connection: %w", err))
		}
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity(certPath, mspID string) *identity.X509Identity {
	certificate := loadCertificate(certPath)

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(fmt.Errorf("failed to create identity from certificate: %w", err))
	}

	return id
}

// loadCertificate loads a certificate from a PEM file.
func loadCertificate(filename string) *x509.Certificate {
	certificatePEM, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("failed to create certificate from file contents: %w", err))
	}

	return certificate
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign(keyPath string) identity.Sign {
	privateKeyPEM, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(fmt.Errorf("failed to create private key from file contents: %w", err))
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to create signing function from private key: %w", err))
	}

	return sign
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
