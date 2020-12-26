REVISION = $(shell git describe --tags)
$(info    Make factom-walletd $(REVISION))

LDFLAGS := "-s -w -X github.com/FactomProject/wallet/v2.WalletVersion=$(REVISION)"

default: factom-walletd
install: factom-walletd-install
all: factom-walletd-darwin-amd64 factom-walletd-windows-amd64.exe factom-walletd-windows-386.exe factom-walletd-linux-amd64 factom-walletd-linux-arm64 factom-walletd-linux-arm7

BUILD_FOLDER := build

factom-walletd:
	go build -trimpath -ldflags $(LDFLAGS)
factom-walletd-install:
	go install -trimpath -ldflags $(LDFLAGS)

factom-walletd-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-walletd-darwin-amd64-$(REVISION)
factom-walletd-windows-amd64.exe:
	env GOOS=windows GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-walletd-windows-amd64-$(REVISION).exe
factom-walletd-windows-386.exe:
	env GOOS=windows GOARCH=386 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-walletd-windows-386-$(REVISION).exe
factom-walletd-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-walletd-linux-amd64-$(REVISION)
factom-walletd-linux-arm64:
	env GOOS=linux GOARCH=arm64 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-walletd-linux-arm64-$(REVISION)
factom-walletd-linux-arm7:
	env GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags $(LDFLAGS) -o $(BUILD_FOLDER)/factom-walletd-linux-arm7-$(REVISION)

.PHONY: clean

clean:
	rm -f factom-walletd
	rm -rf build
