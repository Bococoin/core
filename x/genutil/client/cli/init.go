package cli

import (
	"encoding/json"
	"fmt"
	boco "github.com/Bococoin/core/types"
	"github.com/Bococoin/core/x/staking"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"

	"github.com/Bococoin/core/client/flags"
	"github.com/Bococoin/core/types/module"
	"github.com/Bococoin/core/x/genutil"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	srvCfg "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagOverwrite  = "overwrite"
	flagClientHome = "home-client"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string,
	appMessage json.RawMessage) printInfo {

	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))
	return err
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager,
	defaultNodeHome string) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			chainID := viper.GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			config.Moniker = args[0]
			config.Consensus.TimeoutCommit = time.Second * boco.BlockTime
			config.RPC.MaxSubscriptionsPerClient = 15
			config.P2P.PersistentPeers = "" +
				"ee360fe56121129c93c5f084d1121ab3aeceaf30@95.217.186.21:26656, " +
				"48f14931a5d3bc92f1ea686b927aa0f748ef10ef@95.217.144.116:26656, " +
				"ad05238ba1e342abd2259ba9512d1591eb555681@148.251.6.157:26656, " +
				"010b084161479734e252122cef39f855678fa6e2@95.216.224.245:26656"

			genFile := config.GenesisFile()
			if !viper.GetBool(flagOverwrite) && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			genesis := mbm.DefaultGenesis()

			skGenState := staking.DefaultGenesisState()
			skGenState.Params.BondDenom = boco.DefaultDenom
			genesis[staking.ModuleName], err = staking.ModuleCdc.MarshalJSON(skGenState)
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			appState, err := codec.MarshalJSONIndent(cdc, genesis)
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return errors.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			genDoc.ChainID = chainID
			genDoc.ConsensusParams = types.DefaultConsensusParams()
			genDoc.ConsensusParams.Block.MaxGas = boco.DefaultMaxGas
			genDoc.Validators = nil
			genDoc.AppState = appState
			if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return errors.Wrap(err, "Failed to export genesis file")
			}

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			appConfigFilePath := filepath.Join(config.RootDir, "config/app.toml")

			appConf, _ := srvCfg.ParseConfig()
			appConf.MinGasPrices = boco.DefaultMinGasPrice
			srvCfg.WriteConfigFile(appConfigFilePath, appConf)

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}
