FROM golang

RUN mkdir -p /btcd/log

# Grab and install the latest version of btcd and it's dependencies.
RUN ["go", "get", "-u", "github.com/Masterminds/glide"]
#RUN ["git", "clone", "https://github.com/btcsuite/btcd $GOPATH/src/github.com/btcsuite/btcd"]
RUN git clone https://github.com/btcsuite/btcd $GOPATH/src/github.com/btcsuite/btcd
WORKDIR $GOPATH/src/github.com/btcsuite/btcd
RUN ["glide", "install"]

RUN ["go", "install", ".", "./cmd/.."]

# wallet, p2p, and rpc
EXPOSE 8332 8333 8334

# testnet wallet, p2p, and rpc
EXPOSE 18332 18333 18334

ENTRYPOINT ["btcd"]
#ENTRYPOINT ["/bin/bash"]