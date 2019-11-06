package client

import (
	"github.com/pokt-network/posmint/x/distribution/client/cli"
	"github.com/pokt-network/posmint/x/distribution/client/rest"
	govclient "github.com/pokt-network/posmint/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
