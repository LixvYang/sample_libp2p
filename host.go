package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	noise "github.com/libp2p/go-libp2p-noise"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	run()
}


func run() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	h, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	defer h.Close()

	fmt.Printf("Hello World, my hosts ID is %s\n", h.ID())

	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		panic(err)
	}
	fmt.Println("My private key is:", priv)

	var idht *dht.IpfsDHT

	conn, _ := connmgr.NewConnManager(
		100,         // Lowwater
		400,         // HighWater,
	)
	h2, err := libp2p.New(
		// 使用我们初始化的密钥对
		libp2p.Identity(priv),
		// 多地址监听
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/9000",
			"/ip4/0.0.0.0/udp/9000/quic",
		),
		// 支持TLS
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// 连接管理
		libp2p.ConnectionManager(conn),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),	
		// 使用dht去寻找其他节点
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),

		libp2p.EnableAutoRelay(),
	)	
	if err != nil {
		panic(err)
	}
	defer h2.Close()

	for _, addr := range dht.DefaultBootstrapPeers {
		pi, _ := peer.AddrInfoFromP2pAddr(addr)
		h2.Connect(ctx, *pi)
	}
	fmt.Printf("Hello World, my second hosts ID is %s\n", h2.ID())
}