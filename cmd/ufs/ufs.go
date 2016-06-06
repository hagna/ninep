// This is a ufs server.
package main

import (
	"flag"
	"log"
	"net"

	"github.com/rminnich/ninep/filesystem"
	"github.com/rminnich/ninep/stub"
	"os"
	"runtime/pprof"
)

var (
	ntype      = flag.String("ntype", "tcp4", "Default network type")
	naddr      = flag.String("addr", ":5640", "Network address")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func profile(s string, i ...interface{}) {
	switch s {
	case "Starting readNetPackets":
		// start profile

		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	case "Stop readNetPackets":
		//stop
		pprof.StopCPUProfile()
	}
}

func main() {
	flag.Parse()
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
			s.Trace = nil // log.Printf
			if *cpuprofile != "" {
				s.Trace = profile
			}
			return nil
		})
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}

}
