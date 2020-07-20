package app

import (
	"encoding/json"
	"github.com/Bococoin/core/x/gov"
	paramsclient "github.com/Bococoin/core/x/params/client"
	"github.com/Bococoin/core/x/slashing"
	"github.com/Bococoin/core/x/upgrade"
	upgradeclient "github.com/Bococoin/core/x/upgrade/client"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/Bococoin/core/simapp"
	"github.com/Bococoin/core/types/module"

	"github.com/Bococoin/core/x/auth"
	"github.com/Bococoin/core/x/auth/vesting"
	"github.com/Bococoin/core/x/bank"
	"github.com/Bococoin/core/x/bocomint"
	distr "github.com/Bococoin/core/x/distribution"
	"github.com/Bococoin/core/x/genutil"
	"github.com/Bococoin/core/x/params"
	"github.com/Bococoin/core/x/staking"
	"github.com/Bococoin/core/x/supply"
)

const (
	appName = "bococoin"
)

var (

	// DefaultCLIHome default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.bococoin/cli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.bococoin")

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		slashing.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, upgradeclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		bocomint.AppModuleBasic{},
		upgrade.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		bocomint.ModuleName:       {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc.Seal()
}

// NewApp extended ABCI application
type bocoCoinApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	distrKeeper    distr.Keeper
	supplyKeeper   supply.Keeper
	mintKeeper     bocomint.Keeper
	paramsKeeper   params.Keeper
	govKeeper      gov.Keeper
	slashingKeeper slashing.Keeper
	upgradeKeeper  upgrade.Keeper
	// Module Manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// verify app interface at compile time
var _ simapp.App = (*bocoCoinApp)(nil)

// NewbococoinApp is a constructor function for bococoinApp
func NewBocoCoinApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *bocoCoinApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		bocomint.StoreKey, supply.StoreKey, slashing.StoreKey,
		distr.StoreKey, gov.StoreKey, params.StoreKey, upgrade.StoreKey)

	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &bocoCoinApp{
		BaseApp:   bApp,
		cdc:       cdc,
		keys:      keys,
		tKeys:     tkeys,
		subspaces: make(map[string]params.Subspace),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	// Set specific supspaces
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[bocomint.ModuleName] = app.paramsKeeper.Subspace(bocomint.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		app.keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keys[staking.StoreKey],
		app.supplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keys[distr.StoreKey],
		app.subspaces[distr.ModuleName],
		&stakingKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		app.subspaces[slashing.ModuleName],
	)
	app.upgradeKeeper = upgrade.NewKeeper(
		map[int64]bool{},
		keys[upgrade.StoreKey],
		app.cdc,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper))

	app.govKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.supplyKeeper,
		&stakingKeeper, govRouter,
	)

	app.mintKeeper = bocomint.NewKeeper(
		app.cdc,
		app.keys[bocomint.StoreKey],
		app.subspaces[bocomint.ModuleName],
		app.accountKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks()),
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		bocomint.NewAppModule(app.mintKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
	)

	app.mm.SetOrderBeginBlockers(upgrade.ModuleName, bocomint.ModuleName, distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	app.mm.SetOrderInitGenesis(
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		bocomint.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
	if err != nil {
		tmos.Exit(err.Error())
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

// InitChainer application update at chain initialization
func (app *bocoCoinApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

// BeginBlocker application updates every begin block
func (app *bocoCoinApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *bocoCoinApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// GetKey returns the KVStoreKey for the provided store key
func (app *bocoCoinApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *bocoCoinApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tKeys[storeKey]
}

// LoadHeight loads a particular height
func (app *bocoCoinApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *bocoCoinApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *bocoCoinApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *bocoCoinApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

//_________________________________________________________

func (app *bocoCoinApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)

	return appState, validators, nil
}
