package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestRequestID = "request-id"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

// RequestRandomReq defines the properties of a request rand request's body
type RequestRandomReq struct {
	BaseReq       rest.BaseReq   `json:"base_req" yaml:"base_req"`             // base req
	Consumer      sdk.AccAddress `json:"consumer" yaml:"consumer"`             // request address
	BlockInterval uint64         `json:"block_interval" yaml:"block_interval"` // block interval
}
