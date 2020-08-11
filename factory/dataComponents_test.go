package factory_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/errors"
	"github.com/ElrondNetwork/elrond-go/factory"
	"github.com/ElrondNetwork/elrond-go/factory/mock"
	"github.com/ElrondNetwork/elrond-go/process/economics"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewDataComponentsFactory_NilEconomicsDataShouldErr(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	args.EconomicsData = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, errors.ErrNilEconomicsData, err)
}

func TestNewDataComponentsFactory_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	args.ShardCoordinator = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, errors.ErrNilShardCoordinator, err)
}

func TestNewDataComponentsFactory_NilCoreComponentsShouldErr(t *testing.T) {
	t.Parallel()

	args := getDataArgs(nil)
	args.Core = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, errors.ErrNilCoreComponents, err)
}

func TestNewDataComponentsFactory_NilEpochStartNotifierShouldErr(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	args.EpochStartNotifier = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, errors.ErrNilEpochStartNotifier, err)
}

func TestNewDataComponentsFactory_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	dcf, err := factory.NewDataComponentsFactory(args)
	require.NoError(t, err)
	require.NotNil(t, dcf)
}

func TestDataComponentsFactory_CreateShouldErrDueBadConfig(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	args.Config.ShardHdrNonceHashStorage = config.StorageConfig{}
	dcf, err := factory.NewDataComponentsFactory(args)
	require.NoError(t, err)

	dc, err := dcf.Create()
	require.Error(t, err)
	require.Nil(t, dc)
}

func TestDataComponentsFactory_CreateForShardShouldWork(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	dcf, err := factory.NewDataComponentsFactory(args)

	require.NoError(t, err)
	dc, err := dcf.Create()
	require.NoError(t, err)
	require.NotNil(t, dc)
}

func TestDataComponentsFactory_CreateForMetaShouldWork(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	multiShrdCoord := mock.NewMultiShardsCoordinatorMock(3)
	multiShrdCoord.CurrentShard = core.MetachainShardId
	args.ShardCoordinator = multiShrdCoord
	dcf, err := factory.NewDataComponentsFactory(args)
	require.NoError(t, err)
	dc, err := dcf.Create()
	require.NoError(t, err)
	require.NotNil(t, dc)
}

// ------------ Test ManagedDataComponents --------------------
func TestManagedDataComponents_CreateWithInvalidArgs_ShouldErr(t *testing.T) {
	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	args.Config.ShardHdrNonceHashStorage = config.StorageConfig{}
	managedDataComponents, err := factory.NewManagedDataComponents(factory.DataComponentsHandlerArgs(args))
	require.NoError(t, err)
	err = managedDataComponents.Create()
	require.Error(t, err)
	require.Nil(t, managedDataComponents.Blockchain())
}

func TestManagedDataComponents_Create_ShouldWork(t *testing.T) {
	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	managedDataComponents, err := factory.NewManagedDataComponents(factory.DataComponentsHandlerArgs(args))
	require.NoError(t, err)
	require.Nil(t, managedDataComponents.Blockchain())
	require.Nil(t, managedDataComponents.StorageService())
	require.Nil(t, managedDataComponents.Datapool())

	err = managedDataComponents.Create()
	require.NoError(t, err)
	require.NotNil(t, managedDataComponents.Blockchain())
	require.NotNil(t, managedDataComponents.StorageService())
	require.NotNil(t, managedDataComponents.Datapool())
}

func TestManagedDataComponents_Close(t *testing.T) {
	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	managedDataComponents, _ := factory.NewManagedDataComponents(factory.DataComponentsHandlerArgs(args))
	err := managedDataComponents.Create()
	require.NoError(t, err)

	err = managedDataComponents.Close()
	require.NoError(t, err)
	require.Nil(t, managedDataComponents.Blockchain())
}

// ------------ Test DataComponents --------------------
func TestManagedDataComponents_Close_ShouldWork(t *testing.T) {
	t.Parallel()

	coreComponents := getCoreComponents()
	args := getDataArgs(coreComponents)
	dcf, _ := factory.NewDataComponentsFactory(args)

	dc, _ := dcf.Create()

	err := dc.Close()
	require.NoError(t, err)
}

func getDataArgs(coreComponents factory.CoreComponentsHolder) factory.DataComponentsFactoryArgs {
	testEconomics := &economics.TestEconomicsData{EconomicsData: &economics.EconomicsData{}}
	testEconomics.SetMinGasPrice(200000000000)

	return factory.DataComponentsFactoryArgs{
		Config:             testscommon.GetGeneralConfig(),
		EconomicsData:      testEconomics.EconomicsData,
		ShardCoordinator:   mock.NewMultiShardsCoordinatorMock(2),
		Core:               coreComponents,
		EpochStartNotifier: &mock.EpochStartNotifierStub{},
		CurrentEpoch:       0,
	}
}
