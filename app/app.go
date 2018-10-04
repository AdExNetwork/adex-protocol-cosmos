package app

import (
	"encoding/json"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	adex "github.com/cosmos/cosmos-sdk/adex/x/adex"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
)

const (
	appName = "AdExProtocolApp"
)

type AdExProtocolApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the multistore
	keyMain    *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey
	keyIBC     *sdk.KVStoreKey

	// manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
}

func NewAdExProtocolApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *AdExProtocolApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create your application type
	var app = &AdExProtocolApp{
		cdc:        cdc,
		BaseApp:    bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		keyMain:    sdk.NewKVStoreKey("main"),
		keyAccount: sdk.NewKVStoreKey("acc"),
		keyIBC:     sdk.NewKVStoreKey("ibc"),
	}

	// define and attach the mappers and keepers
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.keyAccount, // target store
		auth.ProtoBaseAccount,
	)
	app.coinKeeper = bank.NewBaseKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))

	// register message routes: all messages starting with adex/ will be routed to the adex handler
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("adex", adex.NewHandler(app.coinKeeper))

	// perform initialization logic
	app.SetInitChainer(app.InitChainer)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))

	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()

	return app
}

func MakeCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ibc.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)

	cdc.RegisterConcrete(types.ClaimTokenMsg{}, "adex/ClaimTokenMsg", nil)
	// @TODO: register custom types here using cdc.RegisterConcrete(&theType{}, "adex/theType", nil)

	cdc.Seal()

	return cdc
}

func (app *AdExProtocolApp) EndBlocker(_ sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

func (app *AdExProtocolApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
		panic(err)
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
			panic(err)
		}

		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}


	return abci.ResponseInitChain{}
}

func (app *AdExProtocolApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*types.GenesisAccount{}

	app.accountMapper.IterateAccounts(ctx, func(acc auth.Account) bool {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	})

	genState := types.GenesisState{Accounts: accounts}
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return json.RawMessage{}, validators, nil
}