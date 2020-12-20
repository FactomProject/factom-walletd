FROM golang:1.15

# Where factom-walletd sources will live
WORKDIR $GOPATH/src/github.com/FactomProject/factom-walletd

# Populate the rest of the source
COPY . .

ARG GOOS=linux

# Build and install factom-walletd
RUN ./build.sh

ENTRYPOINT ["/go/bin/factom-walletd"]

EXPOSE 8089
