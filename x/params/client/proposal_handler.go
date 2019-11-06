package client

import (
	govclient "github.com/pokt-network/posmint/x/gov/client"
	"github.com/pokt-network/posmint/x/params/client/cli"
	"github.com/pokt-network/posmint/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
