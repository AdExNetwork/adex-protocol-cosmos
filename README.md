# AdEx Protocol on Cosmos

This is the Cosmos SDK implementation of the [AdEx Protocol](https://github.com/AdExNetwork/adex-protocol).

This is built with github.com/cosmos/cosmos-sdk, commit 6cbac799125414d225530241cc8ee44be36a76f7

## How to build:

### Prepare the SDK:

```
mkdir -p $GOPATH/src/github.com/cosmos
cd $GOPATH/src/github.com/cosmos
git clone https://github.com/cosmos/cosmos-sdk
cd cosmos-sdk && git checkout 6cbac799125414d225530241cc8ee44be36a76f7
make get_tools && make get_vendor_deps && make install
```

### Get this repo, build adexd/adexcli:

```
git clone https://github.com/AdExNetwork/adex-protocol-cosmos adex
go build -o build/adexd adex/cmd/adexd/main.go && go build -o build/adexcli adex/cmd/adexcli/main.go
```