package notifier_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	coreSlash "github.com/ElrondNetwork/elrond-go-core/data/slash"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	crypto "github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go/consensus/mock"
	mockGenesis "github.com/ElrondNetwork/elrond-go/genesis/mock"
	mockIntegration "github.com/ElrondNetwork/elrond-go/integrationTests/mock"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/slash"
	"github.com/ElrondNetwork/elrond-go/process/slash/notifier"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/ElrondNetwork/elrond-go/testscommon/cryptoMocks"
	"github.com/ElrondNetwork/elrond-go/testscommon/hashingMocks"
	"github.com/ElrondNetwork/elrond-go/testscommon/slashMocks"
	stateMock "github.com/ElrondNetwork/elrond-go/testscommon/state"
	"github.com/ElrondNetwork/elrond-go/update"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewSlashingNotifier(t *testing.T) {
	tests := []struct {
		args        func() *notifier.SlashingNotifierArgs
		expectedErr error
	}{
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.PrivateKey = nil
				return args
			},
			expectedErr: crypto.ErrNilPrivateKey,
		},
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.PublicKey = nil
				return args
			},
			expectedErr: crypto.ErrNilPublicKey,
		},
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.PubKeyConverter = nil
				return args
			},
			expectedErr: update.ErrNilPubKeyConverter,
		},
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.Signer = nil
				return args
			},
			expectedErr: crypto.ErrNilSingleSigner,
		},
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.AccountAdapter = nil
				return args
			},
			expectedErr: state.ErrNilAccountsAdapter,
		},
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.Hasher = nil
				return args
			},
			expectedErr: process.ErrNilHasher,
		},
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.Marshaller = nil
				return args
			},
			expectedErr: process.ErrNilMarshalizer,
		},
		{
			args: func() *notifier.SlashingNotifierArgs {
				args := generateSlashingNotifierArgs()
				args.ProofTxDataExtractor = nil
				return args
			},
			expectedErr: process.ErrNilProofTxDataExtractor,
		},
	}

	for _, currTest := range tests {
		_, err := notifier.NewSlashingNotifier(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidProof_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()
	sn, _ := notifier.NewSlashingNotifier(args)

	tx, err := sn.CreateShardSlashingTransaction(&slashMocks.SlashingProofStub{})
	require.Nil(t, tx)
	require.Equal(t, process.ErrInvalidProof, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidPubKey_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()
	errPubKey := errors.New("pub key error")
	args.PublicKey = &cryptoMocks.PublicKeyStub{
		ToByteArrayStub: func() ([]byte, error) {
			return nil, errPubKey
		},
	}

	sn, _ := notifier.NewSlashingNotifier(args)
	proofStub := &slashMocks.MultipleHeaderProposalProofStub{
		GetHeadersCalled: func() []data.HeaderHandler {
			return slash.HeaderList{&block.HeaderV2{}}
		},
	}

	tx, err := sn.CreateShardSlashingTransaction(proofStub)
	require.Nil(t, tx)
	require.Equal(t, errPubKey, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidAccount_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()
	errAcc := errors.New("accounts adapter error")
	args.AccountAdapter = &stateMock.AccountsStub{
		GetExistingAccountCalled: func([]byte) (vmcommon.AccountHandler, error) {
			return nil, errAcc
		},
	}
	sn, _ := notifier.NewSlashingNotifier(args)
	proofStub := &slashMocks.MultipleHeaderProposalProofStub{
		GetHeadersCalled: func() []data.HeaderHandler {
			return slash.HeaderList{&block.HeaderV2{}}
		},
	}

	tx, err := sn.CreateShardSlashingTransaction(proofStub)
	require.Nil(t, tx)
	require.Equal(t, errAcc, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidMarshaller_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()
	errMarshaller := errors.New("marshaller error")
	args.Marshaller = &testscommon.MarshalizerStub{
		MarshalCalled: func(obj interface{}) ([]byte, error) {
			return nil, errMarshaller
		},
	}

	sn, _ := notifier.NewSlashingNotifier(args)
	proof := &slashMocks.MultipleHeaderProposalProofStub{
		GetHeadersCalled: func() []data.HeaderHandler {
			return slash.HeaderList{&block.HeaderV2{}}
		},
	}

	tx, err := sn.CreateShardSlashingTransaction(proof)
	require.Nil(t, tx)
	require.Equal(t, errMarshaller, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidSlashType_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()

	sn, _ := notifier.NewSlashingNotifier(args)
	proofStub := &slashMocks.MultipleHeaderProposalProofStub{
		GetHeadersCalled: func() []data.HeaderHandler {
			return slash.HeaderList{&block.HeaderV2{}}
		},
		GetTypeCalled: func() coreSlash.SlashingType {
			return 9999999
		},
	}

	tx, err := sn.CreateShardSlashingTransaction(proofStub)

	require.Nil(t, tx)
	require.Equal(t, process.ErrInvalidProof, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidProofTxDataExtractor_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()

	expectedErr := errors.New("invalid tx data extractor")
	args.ProofTxDataExtractor = &slashMocks.ProofTxDataExtractorStub{
		GetProofTxDataCalled: func(proof coreSlash.SlashingProofHandler) (*notifier.ProofTxData, error) {
			return nil, expectedErr
		},
	}
	sn, _ := notifier.NewSlashingNotifier(args)
	tx, err := sn.CreateShardSlashingTransaction(&slashMocks.MultipleHeaderSigningProofStub{})
	require.Nil(t, tx)
	require.Equal(t, expectedErr, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidProofSignature_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()
	errSign := errors.New("signature error")
	args.Signer = &cryptoMocks.SignerStub{
		SignCalled: func(private crypto.PrivateKey, msg []byte) ([]byte, error) {
			return nil, errSign
		},
	}

	sn, _ := notifier.NewSlashingNotifier(args)
	proofStub := &slashMocks.MultipleHeaderProposalProofStub{
		GetHeadersCalled: func() []data.HeaderHandler {
			return slash.HeaderList{&block.HeaderV2{}}
		},
	}

	tx, err := sn.CreateShardSlashingTransaction(proofStub)
	require.Nil(t, tx)
	require.Equal(t, errSign, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_InvalidTxSignature_ExpectError(t *testing.T) {
	args := generateSlashingNotifierArgs()
	errSign := errors.New("signature error")
	flag := false
	args.Signer = &cryptoMocks.SignerStub{
		SignCalled: func(private crypto.PrivateKey, msg []byte) ([]byte, error) {
			if flag {
				return nil, errSign
			}

			flag = true
			return []byte("signature"), nil
		},
	}

	sn, _ := notifier.NewSlashingNotifier(args)
	proofStub := &slashMocks.MultipleHeaderProposalProofStub{
		GetHeadersCalled: func() []data.HeaderHandler {
			return slash.HeaderList{&block.HeaderV2{}}
		},
	}

	tx, err := sn.CreateShardSlashingTransaction(proofStub)
	require.Nil(t, tx)
	require.Equal(t, errSign, err)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_MultipleProposalProof(t *testing.T) {
	round := uint64(100000)
	shardID := uint32(2)

	args := generateSlashingNotifierArgs()
	args.Hasher = &testscommon.HasherStub{
		ComputeCalled: func(string) []byte {
			return []byte{byte('a'), byte('b'), byte('c'), byte('d')}
		},
	}
	args.ProofTxDataExtractor = &slashMocks.ProofTxDataExtractorStub{
		GetProofTxDataCalled: func(proof coreSlash.SlashingProofHandler) (*notifier.ProofTxData, error) {
			return &notifier.ProofTxData{
				Round:     round,
				ShardID:   shardID,
				SlashType: coreSlash.MultipleProposal,
			}, nil
		},
	}

	sn, _ := notifier.NewSlashingNotifier(args)

	h1 := &block.HeaderV2{
		Header: &block.Header{
			Round:        round,
			ShardID:      shardID,
			PrevRandSeed: []byte("seed1"),
		},
	}
	h2 := &block.HeaderV2{
		Header: &block.Header{
			Round:        round,
			ShardID:      shardID,
			PrevRandSeed: []byte("seed2"),
		},
	}

	proof := &slashMocks.MultipleHeaderProposalProofStub{
		GetHeadersCalled: func() []data.HeaderHandler {
			return slash.HeaderList{h1, h2}
		},
	}

	expectedData := []byte(fmt.Sprintf("%s@%s@%d@%d@%s@%s", notifier.BuiltInFunctionSlashCommitmentProof,
		[]byte{slash.MultipleProposalProofID}, shardID, round, []byte{byte('c'), byte('d')}, []byte("signature")))

	expectedTx := &transaction.Transaction{
		Data:      expectedData,
		Nonce:     444,
		SndAddr:   []byte("address"),
		Value:     big.NewInt(notifier.CommitmentProofValue),
		GasPrice:  notifier.CommitmentProofGasPrice,
		GasLimit:  notifier.CommitmentProofGasLimit,
		Signature: []byte("signature"),
	}

	tx, _ := sn.CreateShardSlashingTransaction(proof)
	require.Equal(t, expectedTx, tx)
}

func TestSlashingNotifier_CreateShardSlashingTransaction_MultipleSignProof(t *testing.T) {
	round := uint64(100000)
	shardID := uint32(2)
	pk1 := []byte("pubKey1")

	args := generateSlashingNotifierArgs()
	args.Hasher = &testscommon.HasherStub{
		ComputeCalled: func(string) []byte {
			return []byte{byte('a'), byte('b'), byte('c'), byte('d')}
		},
	}
	args.ProofTxDataExtractor = &slashMocks.ProofTxDataExtractorStub{
		GetProofTxDataCalled: func(proof coreSlash.SlashingProofHandler) (*notifier.ProofTxData, error) {
			return &notifier.ProofTxData{
				Round:     round,
				ShardID:   shardID,
				SlashType: coreSlash.MultipleSigning,
			}, nil
		},
	}

	sn, _ := notifier.NewSlashingNotifier(args)

	h1 := &block.HeaderV2{
		Header: &block.Header{
			Round:        round,
			ShardID:      shardID,
			PrevRandSeed: []byte("seed1"),
		},
	}
	h2 := &block.HeaderV2{
		Header: &block.Header{
			Round:        round,
			ShardID:      shardID,
			PrevRandSeed: []byte("seed2"),
		},
	}

	proof := &slashMocks.MultipleHeaderSigningProofStub{
		GetLevelCalled: func([]byte) coreSlash.ThreatLevel {
			return coreSlash.High
		},
		GetHeadersCalled: func([]byte) []data.HeaderHandler {
			return slash.HeaderList{h1, h2}
		},
		GetPubKeysCalled: func() [][]byte {
			return [][]byte{pk1}
		},
	}

	expectedData := []byte(fmt.Sprintf("%s@%s@%d@%d@%s@%s", notifier.BuiltInFunctionSlashCommitmentProof,
		[]byte{slash.MultipleSigningProofID}, shardID, round, []byte{byte('c'), byte('d')}, []byte("signature")))
	expectedTx := &transaction.Transaction{
		Data:      expectedData,
		Nonce:     444,
		SndAddr:   []byte("address"),
		Value:     big.NewInt(notifier.CommitmentProofValue),
		GasPrice:  notifier.CommitmentProofGasPrice,
		GasLimit:  notifier.CommitmentProofGasLimit,
		Signature: []byte("signature"),
	}

	tx, _ := sn.CreateShardSlashingTransaction(proof)
	require.Equal(t, expectedTx, tx)
}

func generateSlashingNotifierArgs() *notifier.SlashingNotifierArgs {
	accountHandler := &mockGenesis.BaseAccountMock{
		Nonce:             444,
		AddressBytesField: []byte("address"),
	}
	accountsAdapter := &stateMock.AccountsStub{
		GetExistingAccountCalled: func([]byte) (vmcommon.AccountHandler, error) {
			return accountHandler, nil
		},
	}
	marshaller := &testscommon.MarshalizerStub{
		MarshalCalled: func(obj interface{}) ([]byte, error) {
			return nil, nil
		},
	}

	return &notifier.SlashingNotifierArgs{
		PrivateKey:           &mock.PrivateKeyMock{},
		PublicKey:            &mock.PublicKeyMock{},
		PubKeyConverter:      &testscommon.PubkeyConverterMock{},
		Signer:               &mockIntegration.SignerMock{},
		AccountAdapter:       accountsAdapter,
		Hasher:               &hashingMocks.HasherMock{},
		Marshaller:           marshaller,
		ProofTxDataExtractor: &slashMocks.ProofTxDataExtractorStub{},
	}
}