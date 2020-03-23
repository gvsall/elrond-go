package bootstrap

import (
	"github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-go/sharding"
)

// StartOfEpochNodesConfigHandler defines the methods to process nodesConfig from epoch start metablocks
type StartOfEpochNodesConfigHandler interface {
	NodesConfigFromMetaBlock(
		currMetaBlock *block.MetaBlock,
		prevMetaBlock *block.MetaBlock,
		publicKey []byte,
	) (*sharding.NodesCoordinatorRegistry, uint32, error)
	IsInterfaceNil() bool
}
