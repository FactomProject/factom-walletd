// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FactomProject/factom/wallet/wsapi"
)

func main() {
	// configure the server
	var pflag = flag.Int("p", 8089, "set the port to host the wsapi")
	flag.Parse()
	port := *pflag
	
	// setup handling for os signals and stop the server gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("Captured %v, stopping web server and exiting", sig)
			wsapi.Stop()
			os.Exit(1)
		}
	}()
	
	// start the wsapi server
	wsapi.Start(fmt.Sprintf(":%d", port))
}
