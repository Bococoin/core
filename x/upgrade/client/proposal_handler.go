package client

import (
	govclient "github.com/Bococoin/core/x/gov/client"
	"github.com/Bococoin/core/x/upgrade/client/cli"
	"github.com/Bococoin/core/x/upgrade/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
