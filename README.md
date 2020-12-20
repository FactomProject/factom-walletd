factom-walletd
===
[![CircleCI](https://circleci.com/gh/FactomProject/factom-walletd/tree/develop.svg?style=svg)](https://circleci.com/gh/FactomProject/factom-walletd/tree/develop)

`factom-walletd` serves the [wallet/wsapi](https://github.com/FactomProject/wallet/tree/master/wsapi) interface to the wallet library at [wallet](https://github.com/FactomProject/wallet/tree/master).

Here is the [API documentation](https://docs.factomprotocol.org/start/factom-api-docs/factom-walletd-api).

`factom-walletd` manages a database of Factoid and Entry Credit wallets and
uses them to compose transactions for the Factom blockchain. It can compose
transactions for sending Factoids, purchasing Entry Credits, or creating Chains
or Entries.

## Dependencies
### Build Dependencies
- Go > 1.13

### Optional Dependencies

- [`factom-cli`](https://github.com/FactomProject/factom-cli)
- [`factomd`](https://github.com/FactomProject/factomd)

## Package distribution

`factom-walletd` is installed with the Factom package found
[here](https://github.com/FactomProject/distribution). This also installs
`factomd` and `factom-cli`.

## Build and install

Alternatively you can build and install from source.
```
./build.sh
```
