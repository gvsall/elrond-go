package versionswitch_vmquery

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm/arwen/arwenvm"
	"github.com/stretchr/testify/require"
)

func TestSCExecutionWithVMVersionSwitchingEpochRevertAndVMQueries(t *testing.T) {
	t.Skip("work in progress")

	vmConfig := &config.VirtualMachineConfig{
		ArwenVersions: []config.ArwenVersionByEpoch{
			{StartEpoch: 0, Version: "v1.2"},
			{StartEpoch: 1, Version: "v1.2"},
			{StartEpoch: 2, Version: "v1.2"},
			{StartEpoch: 3, Version: "v1.2"},
			{StartEpoch: 4, Version: "v1.3"},
			{StartEpoch: 5, Version: "v1.2"},
			{StartEpoch: 6, Version: "v1.2"},
		},
	}

	gasSchedule, _ := core.LoadGasScheduleConfig("../../../../cmd/node/config/gasSchedules/gasScheduleV2.toml")
	testContext, err := vm.CreateTxProcessorArwenWithVMConfig(
		vm.ArgEnableEpoch{},
		vmConfig,
		gasSchedule,
	)
	require.Nil(t, err)
	defer testContext.Close()

	_ = arwenvm.SetupERC20Test(testContext, "../../testdata/erc20-c-03/wrc20_arwen.wasm")

	err = arwenvm.RunERC20TransactionSet(testContext)
	require.Nil(t, err)

	repeatSwitching := 20
	for i := 0; i < repeatSwitching; i++ {
		epoch := uint32(4)
		testContext.EpochNotifier.CheckEpoch(arwenvm.MakeHeaderHandlerStub(epoch))
		err = arwenvm.RunERC20TransactionSet(testContext)
		require.Nil(t, err)

		epoch = uint32(5)
		testContext.EpochNotifier.CheckEpoch(arwenvm.MakeHeaderHandlerStub(epoch))
		err = arwenvm.RunERC20TransactionSet(testContext)
		require.Nil(t, err)

		epoch = uint32(6)
		testContext.EpochNotifier.CheckEpoch(arwenvm.MakeHeaderHandlerStub(epoch))
		err = arwenvm.RunERC20TransactionSet(testContext)
		require.Nil(t, err)
	}
}
