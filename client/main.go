package main

import (
	"bufio"
	"context"
	"crypto/rand"
	// "crypto/sha256"
	// "encoding/hex"
	// "encoding/json"
	// "flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	// "strconv"
	// "strings"
	// "sync"
	// "time"

	// "github.com/davecgh/go-spew/spew"
	golog "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	// net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	// gologging "github.com/whyrusleeping/go-logging"

	"github.com/HAOYUatHZ/libp2p-Y/common"
)

func main() {
	golog.SetAllLoggers(golog.LevelInfo)

	// listenP := flag.Int("l", 0, "wait for incoming connections")
	// if *listenP == 0 {
	// 	log.Fatal("Please provide a port to bind on with -l")
	// }

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(6001, 0)
	if err != nil {
		log.Fatal(err)
	}

	ha.SetStreamHandler("/p2p/1.0.0", common.HandleStream)

	// The following code extracts target's peer ID from the
	// given multiaddress
	target := os.Args[1]
	ipfsaddr, err := ma.NewMultiaddr(target)
	if err != nil {
		log.Fatalln(err)
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /p2p/1.0.0 protocol
	s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
	if err != nil {
		log.Fatalln(err)
	}
	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	go common.WriteData(rw)
	go common.ReadData(rw)

	select {} // hang forever
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	return basicHost, nil
}
