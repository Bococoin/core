package client

import (
	"github.com/Bococoin/core/x/distribution/client/cli"
	"github.com/Bococoin/core/x/distribution/client/rest"
	govclient "github.com/Bococoin/core/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
