// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factom/wallet/wsapi"
	"github.com/FactomProject/factomd/util"
)

func main() {
	// configure the server
	var (
		pflag = flag.Int("p", 8089, "set the port to host the wsapi")
		wflag = flag.String(
			"w",
			"",
			"set the default wallet location",
		)
		iflag = flag.String("i", "", "Import a version 1 wallet. Set as path to factoid_wallet_bolt.db")
		mflag = flag.String("m", "", "import a wallet from 12 word mnemonic")
		eflag = flag.Bool("e", false, "export a wallet for backup")
		lflag = flag.Bool("l", false, "Create or use an LDB database")

		// Use TLS for the wallet "factom-walletd -wallettls=true"
		walletTLSflag = flag.Bool("wallettls", false, "Set to true to require encrypted connections to the wallet")
		walletTLSKey  = flag.String("walletkey", "", "This file is the PRIVATE TLS key encrypting connections to the wallet. (default ~/.factom/walletAPIpriv.key)")
		walletTLSCert = flag.String("walletcert", "", "This file is the PUBLIC TLS certificate wallet API users will need to connect. (default ~/.factom/walletAPIpub.cert)")

		factomdTLSflag = flag.Bool("factomdtls", false, "Set to true when the factomd API is encrypted")
		factomdTLSCert = flag.String("factomdcert", "", "This file is the TLS certificate provided by the factomd API. (default ~/.factom/m2/factomdAPIpub.cert)")

		walletRpcUser      = flag.String("walletuser", "", "Username to expect before allowing connections")
		walletRpcPassword  = flag.String("walletpassword", "", "Password to expect before allowing connections")
		factomdRpcUser     = flag.String("factomduser", "", "Username for API connections to factomd")
		factomdRpcPassword = flag.String("factomdpassword", "", "Password for API connections to factomd")

		factomdLocation = flag.String("s", "", "IPAddr:port# of factomd API to use to access blockchain (default localhost:8088)")
		walletdLocation = flag.String("selfaddr", "", "comma seperated IPAddresses and DNS names of this factom-walletd to use when creating a cert file")
	)
	flag.Parse()

	// set the wallet path to the wflag or to the default
	walletPath := util.GetHomeDir() + "/.factom/wallet/factom_wallet.db"
	if *lflag {
		walletPath = util.GetHomeDir() + "/.factom/wallet/factom_wallet.ldb"
	}
	if *wflag != "" {
		walletPath = *wflag
	}

	//see if the config file has values which should be used instead of null strings
	filename := util.ConfigFilename() //file name and path to factomd.conf file
	cfg := util.ReadConfig(filename)

	if *walletRpcUser == "" {
		if cfg.Walletd.WalletRpcUser != "" {
			fmt.Printf("using factom-walletd API user and password specified in \"%s\" at WalletRpcUser & WalletRpcPass\n", filename)
			*walletRpcUser = cfg.Walletd.WalletRpcUser
			*walletRpcPassword = cfg.Walletd.WalletRpcPass
		} else {
			fmt.Printf("Warning, factom-walletd API is not password protected. Factoids can be stolen remotely.\n")
		}
	} else {
		fmt.Printf("wallet access protected by user/password specified on command line\n")
	}

	if *factomdRpcUser == "" {
		if cfg.App.FactomdRpcUser != "" {
			fmt.Printf("using factomd API user and password specified in \"%s\" at FactomdRpcUser & FactomdRpcPass\n", filename)
			*factomdRpcUser = cfg.App.FactomdRpcUser
			*factomdRpcPassword = cfg.App.FactomdRpcPass
		}
	}

	if *factomdLocation == "" {
		if cfg.Walletd.FactomdLocation != "localhost:8088" {
			fmt.Printf("using factomd location specified in \"%s\" at FactomdLocation = \"%s\"\n", filename, cfg.Walletd.FactomdLocation)
			*factomdLocation = cfg.Walletd.FactomdLocation
		} else {
			*factomdLocation = "localhost:8088"
		}
	}

	if cfg.Walletd.WalletTlsEnabled == true {
		*walletTLSflag = true
	}
	if *walletTLSflag == true {
		if *walletTLSKey == "" { //if specified, instead use what was on the command line
			if cfg.Walletd.WalletTlsPrivateKey != "/full/path/to/walletAPIpriv.key" { //otherwise check if the the config file has something new
				fmt.Printf("using wallet TLS key file specified in \"%s\" at WalletTlsPrivateKey = \"%s\"\n", filename, cfg.Walletd.WalletTlsPrivateKey)
				*walletTLSKey = cfg.Walletd.WalletTlsPrivateKey
			} else { //if none were specified, use the default file
				*walletTLSKey = fmt.Sprint(util.GetHomeDir(), "/.factom/walletAPIpriv.key")
				fmt.Printf("using default wallet TLS key file \"%s\"\n", *walletTLSKey)
			}
		} else {
			fmt.Printf("using specified wallet TLS key file \"%s\"\n", *walletTLSKey)
		}
		if *walletTLSCert == "" { //if specified, instead use what was on the command line
			if cfg.Walletd.WalletTlsPublicCert != "/full/path/to/walletAPIpub.cert" { //otherwise check if the the config file has something new
				fmt.Printf("using wallet TLS certificate file specified in \"%s\" at WalletTlsPublicCert = \"%s\"\n", filename, cfg.Walletd.WalletTlsPublicCert)
				*walletTLSCert = cfg.Walletd.WalletTlsPublicCert
			} else { //if none were specified, use the default file
				*walletTLSCert = fmt.Sprint(util.GetHomeDir(), "/.factom/walletAPIpub.cert")
				fmt.Printf("using default wallet TLS certificate file \"%s\"\n", *walletTLSCert)
			}
		} else {
			fmt.Printf("using specified wallet TLS certificate file \"%s\"\n", *walletTLSCert)
		}
	} else {
		fmt.Printf("Warning, factom-walletd API connection is unencrypted. Password is unprotected over the network.\n")
	}

	if cfg.Walletd.WalletdLocation != "localhost:8089" {
		//fmt.Printf("using factom-walletd location specified in \"%s\" as WalletdLocation = \"%s\"\n", filename, cfg.Walletd.WalletdLocation)
		var externalIP string
		externalIP += strings.Split(cfg.Walletd.WalletdLocation, ":")[0]
		if *walletdLocation != "" {
			*walletdLocation += ","
		}
		*walletdLocation += externalIP
	}
	if *walletdLocation != "" {
		fmt.Printf("external IP and DNS name to use if making a new TLS keypair = %s\n", *walletdLocation)
	}

	if cfg.App.FactomdTlsEnabled == true {
		*factomdTLSflag = true
	}
	if *factomdTLSflag == true {
		if *factomdTLSCert == "" { //if specified on the command line, don't use the config file
			if cfg.App.FactomdTlsPublicCert != "/full/path/to/factomdAPIpub.cert" { //otherwise check if the the config file has something new
				fmt.Printf("using wallet TLS certificate file specified in \"%s\" at FactomdTlsPublicCert = \"%s\"\n", filename, cfg.App.FactomdTlsPublicCert)
				*factomdTLSCert = cfg.App.FactomdTlsPublicCert
			} else { //if none were specified, use the default file
				*factomdTLSCert = fmt.Sprint(util.GetHomeDir(), "/.factom/m2/factomdAPIpub.cert")
				fmt.Printf("using default factomd TLS certificate file \"%s\"\n", *factomdTLSCert)
			}
		} else {
			fmt.Printf("using specified factomd TLS certificate file \"%s\"\n", *factomdTLSCert)
		}
	}

	port := *pflag
	RPCConfig := factom.RPCConfig{
		WalletTLSEnable:   *walletTLSflag,
		WalletTLSKeyFile:  *walletTLSKey,
		WalletTLSCertFile: *walletTLSCert,
		WalletRPCUser:     *walletRpcUser,
		WalletRPCPassword: *walletRpcPassword,
		WalletServer:      *walletdLocation,
	}
	factom.SetFactomdRpcConfig(*factomdRpcUser, *factomdRpcPassword)
	factom.SetFactomdServer(*factomdLocation)
	factom.SetFactomdEncryption(*factomdTLSflag, *factomdTLSCert)

	if *mflag != "" {
		log.Printf("Creating new wallet with mnemonic")
		w, err := func() (*wallet.Wallet, error) {
			if *lflag {
				return wallet.ImportLDBWalletFromMnemonic(*mflag, walletPath)
			}
			return wallet.ImportWalletFromMnemonic(*mflag, walletPath)
		}()
		if err != nil {
			log.Fatal(err)
		}

		w.Close()
		os.Exit(0)
	}

	if *iflag != "" {
		log.Printf("Importing version 1 wallet %s into %s", *iflag, walletPath)
		w, err := func() (*wallet.Wallet, error) {
			if *lflag {
				return wallet.ImportV1WalletToLDB(*iflag, walletPath)
			}
			return wallet.ImportV1Wallet(*iflag, walletPath)
		}()
		if err != nil {
			log.Fatal(err)
		}
		w.Close()
		os.Exit(0)
	}

	if *eflag {
		m, fs, es, err := func() (string, []*factom.FactoidAddress, []*factom.ECAddress, error) {
			if *lflag {
				return wallet.ExportLDBWallet(walletPath)
			}
			return wallet.ExportWallet(walletPath)
		}()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%q\n", m)
		for _, f := range fs {
			fmt.Println(f.SecString())
		}
		for _, e := range es {
			fmt.Println(e.SecString())
		}
		os.Exit(0)
	}

	// open or create a new wallet file
	fctWallet, err := func() (*wallet.Wallet, error) {
		if *lflag {
			return wallet.NewOrOpenLevelDBWallet(walletPath)
		}
		return wallet.NewOrOpenBoltDBWallet(walletPath)
	}()
	if err != nil {
		log.Fatal(err)
	}

	// open and add a transaction database to the wallet object.
	txdb, err := wallet.NewTXBoltDB(walletPath + ".cache")
	if err != nil {
		log.Println("Could not add transaction database to wallet:", err)
	} else {
		txdb.Update()
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
