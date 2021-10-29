package mandosConverter

import (
	"testing"

	mge "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/elrondgo-exporter"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm"
	"github.com/ElrondNetwork/elrond-go/state"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

// CheckAccounts will verify if mandosAccounts correspond to AccountsAdapter accounts
func CheckAccounts(t *testing.T, accAdapter state.AccountsAdapter, mandosAccounts []*mge.TestAccount) {
	for _, mandosAcc := range mandosAccounts {
		accHandler, err := accAdapter.LoadAccount(mandosAcc.GetAddress())
		require.Nil(t, err)
		account := accHandler.(state.UserAccountHandler)

		require.Equal(t, mandosAcc.GetBalance(), account.GetBalance())
		require.Equal(t, mandosAcc.GetNonce(), account.GetNonce())

		scOwnerAddress := mandosAcc.GetOwner()
		if len(scOwnerAddress) == 0 {
			require.Nil(t, account.GetOwnerAddress())
		} else {
			require.Equal(t, mandosAcc.GetOwner(), account.GetOwnerAddress())
		}

		codeHash := account.GetCodeHash()
		code := accAdapter.GetCode(codeHash)
		require.Equal(t, len(mandosAcc.GetCode()), len(code))

		mandosAccStorage := mandosAcc.GetStorage()
		accStorage := account.DataTrieTracker()
		CheckStorage(t, accStorage, mandosAccStorage)
	}
}

// CheckStorage checks if the dataTrie of an account equals with the storage of the corresponding mandosAccount
func CheckStorage(t *testing.T, dataTrie state.DataTrieTracker, mandosAccStorage map[string][]byte) {
	for key := range mandosAccStorage {
		dataTrieValue, err := dataTrie.RetrieveValue([]byte(key))
		require.Nil(t, err)
		require.Equal(t, mandosAccStorage[key], dataTrieValue)
	}
}

// CheckTransactions checks if the transactions correspond with the mandosTransactions
func CheckTransactions(t *testing.T, transactions []*transaction.Transaction, mandosTransactions []*mge.Transaction) {
	expectedLength := len(mandosTransactions)
	require.Equal(t, expectedLength, len(transactions))
	for i := 0; i < expectedLength; i++ {
		expectedSender := mandosTransactions[i].GetSenderAddress()
		expectedReceiver := mandosTransactions[i].GetReceiverAddress()
		expectedCallValue := mandosTransactions[i].GetCallValue()
		expectedCallFunction := mandosTransactions[i].GetCallFunction()
		expectedCallArguments := mandosTransactions[i].GetCallArguments()
		expectedGasLimit, expectedGasPrice := mandosTransactions[i].GetGasLimitAndPrice()
		expectedNonce := mandosTransactions[i].GetNonce()
		//expectedEsdtTransfers := mandosTransactions[i].GetESDTTransfers()

		require.Equal(t, expectedSender, transactions[i].GetSndAddr())
		require.Equal(t, expectedReceiver, transactions[i].GetRcvAddr())
		require.Equal(t, expectedCallValue, transactions[i].GetValue())
		require.Equal(t, expectedGasLimit, transactions[i].GetGasLimit())
		require.Equal(t, expectedGasPrice, transactions[i].GetGasPrice())
		require.Equal(t, expectedNonce, transactions[i].GetNonce())

		expectedData := createData(expectedCallFunction, expectedCallArguments)
		actualData := transactions[i].GetData()
		require.Equal(t, expectedData, actualData)
	}
}

// SetStateFromMandosTest recieves path to mandosTest, returns a VMTestContext with the specified accounts, an array with the specified transactions and an error
func SetStateFromMandosTest(mandosTestPath string) (testContext *vm.VMTestContext, transactions []*transaction.Transaction, err error) {
	mandosAccounts, deployedMandosAccounts, mandosTransactions, deployMandosTransactions, err := mge.GetAccountsAndTransactionsFromMandos(mandosTestPath)
	if err != nil {
		return nil, nil, err
	}
	testContext, err = vm.CreatePreparedTxProcessorWithVMs(vm.ArgEnableEpoch{})
	if err != nil {
		return nil, nil, err
	}
	err = CreateAccountsFromMandosAccs(testContext, mandosAccounts)
	if err != nil {
		return nil, nil, err
	}
	err = DeploySCsFromMandosDeployTxs(testContext, deployMandosTransactions, mandosTransactions, deployedMandosAccounts)
	if err != nil {
		return nil, nil, err
	}
	transactions = CreateTransactionsFromMandosTxs(mandosTransactions)
	return testContext, transactions, nil
}

// RunSingleTransactionBenchmark receives the VMTestContext (which can be created with SetStateFromMandosTest), a tx and performs a benchmark on that specific tx. If processing transaction fails, it will return error, else will return nil
func RunSingleTransactionBenchmark(b *testing.B, testContext *vm.VMTestContext, tx *transaction.Transaction) (err error) {
	var returnCode vmcommon.ReturnCode
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		returnCode, err = testContext.TxProcessor.ProcessTransaction(tx)
		tx.Nonce++
	}
	b.StopTimer()
	if err != nil {
		return err
	}
	if returnCode != vmcommon.Ok {
		return errReturnCodeNotOk
	}
	return nil
}
