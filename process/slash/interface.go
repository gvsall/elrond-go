package slash

import (
	"github.com/ElrondNetwork/elrond-go-core/data"
	coreSlash "github.com/ElrondNetwork/elrond-go-core/data/slash"
	"github.com/ElrondNetwork/elrond-go/process"
)

// SlashingDetector - checks for slashable events and generates proofs to be used for slash
type SlashingDetector interface {
	// VerifyData - checks if an intercepted data represents a slashable event and returns a proof if so,
	// otherwise returns nil and error
	VerifyData(data process.InterceptedData) (coreSlash.SlashingProofHandler, error)
	// ValidateProof - checks if a given proof is valid
	ValidateProof(proof coreSlash.SlashingProofHandler) error
}

// SlashingNotifier - creates a transaction from the generated proof of the slash detector and sends it to the network
type SlashingNotifier interface {
	// CreateShardSlashingTransaction - creates a slash transaction from the generated SlashingProofHandler
	CreateShardSlashingTransaction(proof coreSlash.SlashingProofHandler) (data.TransactionHandler, error)
	// CreateMetaSlashingEscalatedTransaction - creates a transaction for the metachain if x rounds passed
	// and no slash transaction has been created by any of the previous x proposers
	CreateMetaSlashingEscalatedTransaction(proof coreSlash.SlashingProofHandler) data.TransactionHandler
}

// SlashingTxProcessor - processes the proofs from the SlashingNotifier inside shards
type SlashingTxProcessor interface {
	// ProcessTx - processes a slash transaction that contains a proof from the SlashingNotifier
	// if the proof is valid, a SCResult with destination metachain is created,
	// where the actual slash actions are taken (jail, inactivate, remove balance etc)
	ProcessTx(transaction data.TransactionHandler) data.TransactionHandler
}

// Slasher - processes the validated slash proof from the shards
// and applies the necessary actions (jail, inactivate, remove balance etc)
type Slasher interface {
	// ExecuteSlash - processes a slash SCResult that contains information about the slashable event
	// validator could be jailed, inactivated, balance can be decreased
	ExecuteSlash(transaction data.TransactionHandler) data.TransactionHandler
}