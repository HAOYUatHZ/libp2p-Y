package main

import (
	"io"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
)

func main(
	// Make a host that listens on the given multiaddress
    ha, err := makeBasicHost(*listenF )
    if err != nil {
        log.Fatal(err)
    }


    log.Println(ha)
)


// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int) (host.Host, error) {

    // If the seed is zero, use real cryptographic randomness. Otherwise, use a
    // deterministic randomness source to make generated keys stay the same
    // across multiple runs
    var r io.Reader
    // if randseed == 0 {
    //     r = rand.Reader
    // } else {
    r = mrand.New(mrand.NewSource(randseed))
    // }

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
    if secio {
        log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
    } else {
        log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
    }

    return basicHost, nil
}

