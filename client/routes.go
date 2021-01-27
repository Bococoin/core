package client

import (
	"github.com/gorilla/mux"

	"github.com/Bococoin/core/client/context"
	"github.com/Bococoin/core/client/rpc"
)

// Register routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	rpc.RegisterRPCRoutes(cliCtx, r)
}
