package main

import (
	// "bufio"
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
	// "os"
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
	// peer "github.com/libp2p/go-libp2p-peer"
	// pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	// gologging "github.com/whyrusleeping/go-logging"

	"github.com/HAOYUatHZ/libp2p-Y/common"
)

func main() {
	// golog.LevelDebug
	// golog.LevelInfo
	// golog.LevelError
	// golog.LevelFatal
	// golog.LevelPanic
	golog.SetAllLoggers(golog.LevelInfo)

	// listenP := flag.Int("l", 0, "wait for incoming connections")
	// if *listenP == 0 {
	// 	log.Fatal("Please provide a port to bind on with -l")
	// }

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(6000, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening for connections")
	// Set a stream handler on host A. /p2p/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/p2p/1.0.0", common.HandleStream)

	select {}
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
	// log.Printf("Now run \"go run client/main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	log.Printf("Now run \"go run client/main.go %s\" on a different terminal\n", fullAddr)
	return basicHost, nil
}
