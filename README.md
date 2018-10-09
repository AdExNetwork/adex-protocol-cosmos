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
cd adex
make
```


## Message types

Please see [OCEAN](https://github.com/AdExNetwork/adex-protocol/blob/master/OCEAN.md) to better understand the way the start/finalize operations work.

Prior to any of the following messages, bids are unknown to the blockchain. If they're cancelled, their hash is being flagged as cancelled so they cannot ever be interacted with. Otherwise, they get wrapped into a commitment using `commitmentStartMsg`, and once the commitment goal has been delivered, we call `commitmentFinalizeMsg`. For more details on the protocol, see [AdEx Protocol](https://github.com/AdExNetwork/adex-protocol).

`bidCancelMsg` - cancel a bid

`commitmentStartMsg` - start an OCEAN commitment for a specific bid

`commitmentFinalizeMsg` - finalize an OCEAN commitment, submit validator votes

Unlike the Ethereum implementation, the `Timeout` step is not needed here, since it will happen automatically on every block (via `EndBlocker`)


## Rationale, discussion

Compared to the existing Ethereum implementation, a Cosmos implementation offers:

* Scalability and interoperability
* More flexibility, allowing us to tweak fees, perform operations at the end of each block, therefore improving UX
* Possible advancements in upgradability because of the governance

See https://github.com/AdExNetwork/adex-core/issues/12


## Bootstrapping

We have not found any official guidelines on bootstrapping a new Cosmos SDK app. What we did is copied `examples/basecoin` from the official repository and stripped it down, most notably removing the custom `Account` type.

Then, we added a simple `Makefile` that just invokes `go build`.

If you have the `cosmos-sdk` repository in your `$GOPATH` and you've ran all it's needed steps (`make get_tools && make get_vendor_deps && make install`), you will be able to compile your app.

## There be dragons

Because the experimental nature of the cosmos-sdk, there are many potential security (and otherwise) risks that might exist in this codebase.

For example:

* the security of the arithmetic operations with `sdk.Coins`
* nil pointers/slices and all the `sdk` types that are pointers/slices
* go-amino: possible serialization/deserialization bugs and inconsistencies
* sdk.Coins can be negative

