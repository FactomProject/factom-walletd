vendorver=./vendor/github.com/FactomProject/factom/wallet/VERSION
godirver=$GOPATH/src/github.com/FactomProject/factom/wallet/VERSION

if [ -f $vendorver ]; then
    go install -ldflags "-X github.com/FactomProject/factom/wallet.WalletVersion=`cat $vendorver`" -v
else
    go install -ldflags "-X github.com/FactomProject/factom/wallet.WalletVersion=`cat $godirver`" -v
fi
