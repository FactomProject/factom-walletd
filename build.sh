# Version of walletd should be the same as the wallet library
VERSION=$(cat go.mod | grep "github.com/FactomProject/wallet" | cut -d' ' -f2)
go install -ldflags "-X github.com/FactomProject/wallet/v2.WalletVersion=$VERSION" -v

