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
	"os/user"
	"syscall"

	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factom/wallet/wsapi"
)

var homedir = func() string {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return usr.HomeDir
}()

func main() {
	// configure the server
	var (
		pflag = flag.Int("p", 8089, "set the port to host the wsapi")
		wflag = flag.String("w", fmt.Sprint(homedir, "/.factom/wallet"),
			"set the default wallet location")
	)
	flag.Parse()

	port := *pflag

	// open or create a new wallet file
	fctWallet, err := wallet.NewOrOpenWallet(*wflag)
	if err != nil {
		log.Fatal(err)
	}

	// setup handling for os signals and stop the server gracefully
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		for sig := range c {
			log.Printf("Captured %v, stopping web server and exiting", sig)
			wsapi.Stop()
		}
	}()

	// start the wsapi server
	wsapi.Start(fctWallet, fmt.Sprintf(":%d", port))
}
