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


## Message types

Please see [OCEAN](https://github.com/AdExNetwork/adex-protocol/blob/master/OCEAN.md) to better understand the way the start/finalize operations work.

`commitmentStartMsg` - start an OCEAN commitment for a specific bid

`commitmentFinalizeMsg` - finalize an OCEAN commitment, submit validator votes

Unlike the Ethereum implementation, the `Timeout` step is not needed here, since it will happen automatically on every block (via `EndBlocker`)


## Rationale, discussion

Compared to the existing Ethereum implementation, a Cosmos implementation offers:

* Scalability and interoperability
* More flexibility, allowing us to tweak fees, perform operations at the end of each block, therefore improving UX
* Possible advancements in upgradability because of the governance

See https://github.com/AdExNetwork/adex-core/issues/12


