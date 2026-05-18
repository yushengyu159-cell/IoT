package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ChannelContract struct {
	Network  *client.Network
	Contract *client.Contract
}

type MultiChannelService struct {
	gateway  *client.Gateway
	channels map[string]*ChannelContract
	mu       sync.RWMutex
}

var MultiChannel = new(MultiChannelService)

var channelChaincodeMap = map[string]string{
	"mychannel":        "basic",
	"access-channel":   "access_cc",
	"billing-channel":  "billing_cc",
	"maintain-channel": "maintain_cc",
	"esg-channel":      "esg_cc",
}

func (m *MultiChannelService) Init(ctx context.Context) error {
	g.Log().Info(ctx, "Initializing MultiChannel service...")
	peerURL := os.Getenv("FABRIC_PEER_URL")
	if peerURL == "" {
		peerURL = "peer0.org1.example.com:7051"
	}
	tlsCAPath := "/app/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
	tlsCABytes, err := os.ReadFile(tlsCAPath)
	if err != nil {
		g.Log().Warning(ctx, "Failed to load TLS CA cert, falling back to insecure:", err)
		return m.initInsecure(ctx, peerURL)
	}

	cp := x509.NewCertPool()
	cp.AppendCertsFromPEM(tlsCABytes)
	creds := credentials.NewTLS(&tls.Config{RootCAs: cp, ServerName: "peer0.org1.example.com", InsecureSkipVerify: true})

	g.Log().Info(ctx, "Connecting to Fabric Peer with TLS:", peerURL)
	grpcCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(grpcCtx, peerURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("gRPC connect failed: %v", err)
	}

	gw, err := m.createGateway(conn)
	if err != nil {
		return err
	}

	return m.setupChannels(ctx, gw)
}

func (m *MultiChannelService) initInsecure(ctx context.Context, peerURL string) error {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})

	grpcCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(grpcCtx, peerURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("gRPC connect failed: %v", err)
	}

	gw, err := m.createGateway(conn)
	if err != nil {
		return err
	}

	return m.setupChannels(ctx, gw)
}

func (m *MultiChannelService) createGateway(conn *grpc.ClientConn) (*client.Gateway, error) {
	certPath := "/app/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
	keyPath := "/app/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"

	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("read cert failed: %v", err)
	}
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read key failed: %v", err)
	}

	block, _ := pem.Decode(certBytes)
	if block == nil {
		return nil, fmt.Errorf("decode cert failed")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse cert failed: %v", err)
	}

	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return nil, fmt.Errorf("decode key failed")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse key failed: %v", err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("create sign failed: %v", err)
	}
	id, err := identity.NewX509Identity("Org1MSP", cert)
	if err != nil {
		return nil, fmt.Errorf("create identity failed: %v", err)
	}

	gw, err := client.Connect(id, client.WithClientConnection(conn), client.WithSign(sign))
	if err != nil {
		return nil, fmt.Errorf("gateway connect failed: %v", err)
	}

	return gw, nil
}

func (m *MultiChannelService) setupChannels(ctx context.Context, gw *client.Gateway) error {
	m.gateway = gw
	m.channels = make(map[string]*ChannelContract)

	for channelName, chaincodeName := range channelChaincodeMap {
		net := gw.GetNetwork(channelName)
		contract := net.GetContract(chaincodeName)
		m.channels[channelName] = &ChannelContract{
			Network:  net,
			Contract: contract,
		}
		g.Log().Infof(ctx, "Channel [%s] -> Chaincode [%s] initialized", channelName, chaincodeName)
	}

	g.Log().Info(ctx, "MultiChannel service initialized successfully")
	return nil
}

func (m *MultiChannelService) GetContract(channelName string) (*client.Contract, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cc, ok := m.channels[channelName]
	if !ok {
		return nil, fmt.Errorf("channel %s not found", channelName)
	}
	return cc.Contract, nil
}

func (m *MultiChannelService) SubmitTransaction(ctx context.Context, channelName, fn string, args ...string) ([]byte, error) {
	contract, err := m.GetContract(channelName)
	if err != nil {
		return nil, err
	}
	return contract.SubmitTransaction(fn, args...)
}

func (m *MultiChannelService) EvaluateTransaction(ctx context.Context, channelName, fn string, args ...string) ([]byte, error) {
	contract, err := m.GetContract(channelName)
	if err != nil {
		return nil, err
	}
	return contract.EvaluateTransaction(fn, args...)
}

func (m *MultiChannelService) ListChannels() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var channels []string
	for name := range m.channels {
		channels = append(channels, name)
	}
	return channels
}

func (m *MultiChannelService) Close() {
	if m.gateway != nil {
		m.gateway.Close()
	}
}
