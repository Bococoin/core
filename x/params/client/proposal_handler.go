package client

import (
	govclient "github.com/Bococoin/core/x/gov/client"
	"github.com/Bococoin/core/x/params/client/cli"
	"github.com/Bococoin/core/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
