package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	restName = "name"
)

// RegisterRoutes registers nameservice-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, querierRoute string) {
	registerQueryRoutes(cliCtx, r, querierRoute)
	registerTxRoutes(cliCtx, r, querierRoute)
}
