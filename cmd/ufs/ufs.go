// This is a ufs server.
package main

import (
	"flag"
	"log"
	"net"

	"github.com/harvey-os/ninep/filesystem"
	"github.com/harvey-os/ninep/stub"
	"os"
	"runtime/pprof"
)

var (
	ntype      = flag.String("ntype", "tcp4", "Default network type")
	naddr      = flag.String("addr", ":5640", "Network address")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func profile(c net.Conn) func(s string, i ...interface{}) {
	var ppFile *os.File
	var ident string
	ident = c.RemoteAddr().String()
	ppFile, err := os.Create(*cpuprofile + ident)
	if err != nil {
		log.Fatal(err)
	}
	return func(s string, i ...interface{}) {
		switch s {
		case "Starting readNetPackets":
			// start profile
			log.Println("starting profile for", ident)

			pprof.StartCPUProfile(ppFile)
		case "Stop readNetPackets":
			//stop
			log.Println("Writing profile", ident)
			pprof.StopCPUProfile()
			ppFile.Close()
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	l, err := net.Listen(*ntype, *naddr)
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Accept: %v", err)
		}

		_, err = ufs.NewUFS(func(s *stub.Server) error {
			s.FromNet, s.ToNet = c, c
			s.Trace = log.Printf
			if *cpuprofile != "" {

				s.Trace = profile(c)
			}
			return nil
		})
		if err != nil {
			log.Printf("Error: %v", err)
		}

	}

}
