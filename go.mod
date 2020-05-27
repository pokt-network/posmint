module github.com/pokt-network/posmint

go 1.13

require (
	github.com/gogo/protobuf v1.3.1
	github.com/magiconair/properties v1.8.1
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/iavl v0.12.4
	github.com/tendermint/tendermint v0.32.10
	github.com/tendermint/tm-db v0.2.0
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/tendermint/tendermint => github.com/pokt-network/tendermint v0.32.11-0.20200416214829-c67ffb7bf00f
