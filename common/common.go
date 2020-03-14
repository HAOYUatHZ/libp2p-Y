package common

import (
	"bufio"
	// "context"
	// "crypto/rand"
	// "crypto/sha256"
	// "encoding/hex"
	"encoding/json"
	// "flag"
	"fmt"
	// "io"
	"log"
	// mrand "math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	// golog "github.com/ipfs/go-log"
	// libp2p "github.com/libp2p/go-libp2p"
	// crypto "github.com/libp2p/go-libp2p-crypto"
	// host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	// peer "github.com/libp2p/go-libp2p-peer"
	// pstore "github.com/libp2p/go-libp2p-peerstore"
	// ma "github.com/multiformats/go-multiaddr"
	// gologging "github.com/whyrusleeping/go-logging"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       string
	Hash      string
	PrevHash  string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

var mutex = &sync.Mutex{}

func HandleStream(s net.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func ReadData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(chain) > len(Blockchain) {
				Blockchain = chain
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color:     \x1b[32m
				// Reset console color:     \x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			mutex.Unlock()
		}
	}
}

func WriteData(rw *bufio.ReadWriter) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		// newBlock :=Block{} generateBlock(Blockchain[len(Blockchain)-1], bpm)
		newBlock := Block{BPM: sendData}

		// if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		mutex.Unlock()
		// }

		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}

		spew.Dump(Blockchain)

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	}
}
