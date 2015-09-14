package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/pin/tftp"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
)

func main() {
	address := flag.String("address", "127.0.0.1", "Address of the machine you wish to send requests to")
	port := flag.Int("port", 69, "Port that is listening to TFTP requests")
	tmpFile := flag.String("tmpfile", "/var/tmp/testfile.img", "Absolute path for a temporary file that will be used during the testing")
	oui := flag.String("oui", "0004f2", "Phone OUI used for the file name in the GET request.")
	num := flag.Int("num", 100, "Number of requests you wish to send")
	flag.Parse()

	fullAddress := fmt.Sprintf("%s:%d", *address, *port)

	addr, e := net.ResolveUDPAddr("udp", fullAddress)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Could not resolve UDP address")
	}

	file, e := os.Create(*tmpFile)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Could not create file /var/tmp/testfile.img")
	}

	w := bufio.NewWriter(file)
	log := log.New(os.Stderr, "", log.Ldate|log.Ltime)
	wg := new(sync.WaitGroup)
	wg.Add(*num)
	hex := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

	for loop := 0; loop < *num; loop++ {

		/* Generate random mac address by appending to the OUI */
		mac := bytes.NewBufferString(*oui)
		for i := 0; i < 6; i++ {
			mac.WriteString(hex[rand.Intn(16)])
		}
		filename := mac.String()

		c := tftp.Client{addr, log}

		go c.Get(filename, "netascii", func(reader *io.PipeReader) {
			n, readError := w.ReadFrom(reader)
			if readError != nil {
				fmt.Fprintf(os.Stderr, "Can't get %s: %v\n", filename, readError)
			} else {
				fmt.Fprintf(os.Stderr, "Got %s (%d bytes)\n", filename, n)
			}

			w.Flush()
			file.Close()
			wg.Done()
		})
	}

	wg.Wait()
	os.Remove(*tmpFile)
}
