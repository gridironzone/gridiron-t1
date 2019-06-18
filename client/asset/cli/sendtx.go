package cli

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func preSignCmd(cmd *cobra.Command, _ []string) {
	// Conditionally mark the account and sequence numbers required as no RPC
	// query will be done.
	if viper.GetString(FlagSource) == "gateway" {
		cmd.MarkFlagRequired(FlagGateway)
	}
}

// GetCmdIssueAsset implements the issue asset command
func GetCmdIssueAsset(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token",
		Short: "issue a new token",
		Example: "iriscli asset issue-token --family=<family> --source=<source> --gateway=<gateway-moniker>" +
			" --symbol=<symbol> --name=<token-name> --initial-supply=<initial-supply> --from=<key-name> --chain-id=<chain-id> --fee=0.6iris",
		PreRun: preSignCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			family, ok := asset.StringToAssetFamilyMap[viper.GetString(FlagFamily)]
			if !ok {
				return fmt.Errorf("invalid token family type %s", viper.GetString(FlagFamily))
			}

			source, ok := asset.StringToAssetSourceMap[viper.GetString(FlagSource)]
			if !ok {
				return fmt.Errorf("invalid token source type %s", viper.GetString(FlagSource))
			}

			msg := asset.MsgIssueToken{
				Family:         family,
				Source:         source,
				Gateway:        viper.GetString(FlagGateway),
				Symbol:         viper.GetString(FlagSymbol),
				SymbolAtSource: viper.GetString(FlagSymbolAtSource),
				Name:           viper.GetString(FlagName),
				Decimal:        uint8(viper.GetInt(FlagDecimal)),
				SymbolMinAlias: viper.GetString(FlagSymbolMinAlias),
				InitialSupply:  uint64(viper.GetInt(FlagInitialSupply)),
				MaxSupply:      uint64(viper.GetInt(FlagMaxSupply)),
				Mintable:       viper.GetBool(FlagMintable),
				Owner:          owner,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsAssetIssue)
	cmd.MarkFlagRequired(FlagFamily)
	cmd.MarkFlagRequired(FlagSource)
	cmd.MarkFlagRequired(FlagSymbol)
	cmd.MarkFlagRequired(FlagName)
	cmd.MarkFlagRequired(FlagInitialSupply)

	return cmd
}

// GetCmdCreateGateway implements the create gateway command
func GetCmdCreateGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gateway",
		Short: "create a gateway",
		Example: "iriscli asset create-gateway --moniker=<moniker> --identity=<identity> --details=<details> " +
			"--website=<website> --create-fee=<gateway create fee>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			moniker := viper.GetString(FlagMoniker)
			identity := viper.GetString(FlagIdentity)
			details := viper.GetString(FlagDetails)
			website := viper.GetString(FlagWebsite)
			createFee := viper.GetString(FlagCreateFee)

			createFeeCoin, err := sdk.ParseCoin(createFee)
			if err != nil {
				return err
			}

			if createFeeCoin.Denom == sdk.NativeTokenName {
				createFeeCoin = sdk.NewCoin(sdk.NativeTokenMinDenom, sdk.NewIntWithDecimal(createFeeCoin.Amount.Int64(), 18))
			}

			var msg sdk.Msg
			msg = asset.NewMsgCreateGateway(
				owner, moniker, identity, details, website, createFeeCoin,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsGatewayCreate)
	cmd.MarkFlagRequired(FlagMoniker)
	cmd.MarkFlagRequired(FlagCreateFee)

	return cmd
}

// GetCmdEditGateway implements the edit gateway command
func GetCmdEditGateway(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-gateway",
		Short: "edit a gateway",
		Example: "iriscli asset edit-gateway --moniker=<moniker> --identity=<identity> --details=<details> " +
			"--website=<website>",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))
			txCtx := utils.NewTxContextFromCLI().WithCodec(cdc).
				WithCliCtx(cliCtx)

			owner, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			moniker := viper.GetString(FlagMoniker)
			identity := viper.Get(FlagIdentity)
			details := viper.Get(FlagDetails)
			website := viper.Get(FlagWebsite)

			pIdentity := (*string)(nil)
			pDetails := (*string)(nil)
			pWebsite := (*string)(nil)

			if identity != nil {
				sIdentity := identity.(string)
				pIdentity = &sIdentity
			}
			if details != nil {
				sDetails := details.(string)
				pDetails = &sDetails
			}
			if website != nil {
				sWebsite := website.(string)
				pWebsite = &sWebsite
			}

			var msg sdk.Msg
			msg = asset.NewMsgEditGateway(
				owner, moniker, pIdentity, pDetails, pWebsite,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.SendOrPrintTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsGatewayEdit)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}
