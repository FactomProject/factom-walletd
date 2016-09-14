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

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factom/wallet/wsapi"
	"github.com/FactomProject/factomd/util"
)

var homedir = func() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}()

func main() {
	// configure the server
	var (
		pflag = flag.Int("p", 8089, "set the port to host the wsapi")
		wflag = flag.String("w", fmt.Sprint(homedir, "/.factom/wallet.db"),
			"set the default wallet location")
		iflag       = flag.String("i", "", "import a version 1 wallet")
		TLSflag     = flag.Bool("tls", false, "enable tls") //to get tls, run as "factom-walletd -tls=true"
		TLSKeyflag  = flag.String("key", fmt.Sprint(homedir, "/.factom/tlspub.cert"), "set the default tls key location")
		TLSCertflag = flag.String("cert", fmt.Sprint(homedir, "/.factom/tlspriv.key"), "set the default tls cert location")
		//rpcUserflag     = flag.String("rpcuser", "", "Username for JSON-RPC connections")
		//rpcPasswordflag = flag.String("rpcpassword", "", "Password for JSON-RPC connections")

		walletRpcUser      = flag.String("walletuser", "", "Username to expect before allowing connections")
		walletRpcPassword  = flag.String("walletpassword", "", "Password to expect before allowing connections")
		factomdRpcUser     = flag.String("factomduser", "", "Username for API connections to factomd")
		factomdRpcPassword = flag.String("factomdpassword", "", "Password for API connections to factomd")
	)
	flag.Parse()
	/*filename := util.ConfigFilename()
	cfg := util.ReadConfig(filename)
	if *rpcUserflag == "" {
		*rpcUserflag = cfg.FactomdRPCPassword
	}
	if *rpcPasswordflag == "" {
		*rpcPasswordflag = cfg.FactomdRPCPassword
	}
	if *rpcUserflag == "" || *rpcPasswordflag == "" {
		log.Fatal("Rpc user and password did not set, using -rpcuser and -rpcpassword or config file")
	}*/

	//see if the config file has values which should be used instead of null strings
	filename := util.ConfigFilename() //file name and path to factomd.conf file
	cfg := util.ReadConfig(filename).Rpc
	cfgw := util.ReadConfig(filename).Walletd

	if *walletRpcUser == "" {
		if cfgw.WalletRpcUser != "" {
			fmt.Printf("using factom-walletd API user and password specified in \"%s\" at WalletRpcUser & WalletRpcPass\n", filename)
			*walletRpcUser = cfgw.WalletRpcUser
			*walletRpcPassword = cfgw.WalletRpcPass
		}
	}

	if *factomdRpcUser == "" {
		if cfg.FactomdRpcUser != "" {
			fmt.Printf("using factomd API user and password specified in \"%s\" at FactomdRpcUser & FactomdRpcPass\n", filename)
			*factomdRpcUser = cfg.FactomdRpcUser
			*factomdRpcPassword = cfg.FactomdRpcPass
		}
	}

	port := *pflag
	RPCConfig := factom.RPCConfig{
		TLSEnable:          *TLSflag,
		TLSKeyFile:         *TLSKeyflag,
		TLSCertFile:        *TLSCertflag,
		WalletRPCUser:      *walletRpcUser,
		WalletRPCPassword:  *walletRpcPassword,
		FactomdRPCUser:     *factomdRpcUser,
		FactomdRPCPassword: *factomdRpcPassword,
	}

	if *iflag != "" {
		log.Printf("Importing version 1 wallet %s into %s", *iflag, *wflag)
		w, err := wallet.ImportV1Wallet(*iflag, *wflag)
		if err != nil {
			log.Fatal(err)
		}
		w.Close()
		os.Exit(0)
	}

	// open or create a new wallet file
	fctWallet, err := wallet.NewOrOpenBoltDBWallet(*wflag)
	if err != nil {
		log.Fatal(err)
	}

	// open and add a transaction database to the wallet object.
	txdb, err := wallet.NewTXBoltDB(fmt.Sprint(homedir, "/.factom/txdb.db"))
	if err != nil {
		log.Println("Could not add transaction database to wallet:", err)
	} else {
		fctWallet.AddTXDB(txdb)
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
	wsapi.Start(fctWallet, fmt.Sprintf(":%d", port), RPCConfig)
}
