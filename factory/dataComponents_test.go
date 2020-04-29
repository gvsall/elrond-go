package factory_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/factory"
	"github.com/ElrondNetwork/elrond-go/factory/mock"
	"github.com/ElrondNetwork/elrond-go/process/economics"
	"github.com/ElrondNetwork/elrond-go/storage/storageUnit"
	"github.com/stretchr/testify/require"
)

func TestNewDataComponentsFactory_NilConfigurationShouldErr(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	args.Config = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, factory.ErrNilConfiguration, err)
}

func TestNewDataComponentsFactory_NilEconomicsDataShouldErr(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	args.EconomicsData = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, factory.ErrNilEconomicsData, err)
}

func TestNewDataComponentsFactory_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	args.ShardCoordinator = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, factory.ErrNilShardCoordinator, err)
}

func TestNewDataComponentsFactory_NilCoreComponentsShouldErr(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	args.Core = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, factory.ErrNilCoreComponents, err)
}

func TestNewDataComponentsFactory_NilPathManagerShouldErr(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	args.PathManager = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, factory.ErrNilPathManager, err)
}

func TestNewDataComponentsFactory_NilEpochStartNotifierShouldErr(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	args.EpochStartNotifier = nil

	dcf, err := factory.NewDataComponentsFactory(args)
	require.Nil(t, dcf)
	require.Equal(t, factory.ErrNilEpochStartNotifier, err)
}

func TestNewDataComponentsFactory_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	args := getDataArgs()

	dcf, err := factory.NewDataComponentsFactory(args)
	require.NoError(t, err)
	require.NotNil(t, dcf)
}

func TestDataComponentsFactory_CreateShouldErrDueBadConfig(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	args.Config.ShardHdrNonceHashStorage = config.StorageConfig{}
	dcf, err := factory.NewDataComponentsFactory(args)
	require.NoError(t, err)

	dc, err := dcf.Create()
	require.Error(t, err)
	require.Nil(t, dc)
}

func TestDataComponentsFactory_CreateForShardShouldWork(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	dcf, err := factory.NewDataComponentsFactory(args)

	require.NoError(t, err)
	dc, err := dcf.Create()
	require.NoError(t, err)
	require.NotNil(t, dc)
}

func TestDataComponentsFactory_CreateForMetaShouldWork(t *testing.T) {
	t.Parallel()

	args := getDataArgs()
	multiShrdCoord := mock.NewMultiShardsCoordinatorMock(3)
	multiShrdCoord.CurrentShard = core.MetachainShardId
	args.ShardCoordinator = multiShrdCoord
	dcf, err := factory.NewDataComponentsFactory(args)
	require.NoError(t, err)
	dc, err := dcf.Create()
	require.NoError(t, err)
	require.NotNil(t, dc)
}

func getDataArgs() factory.DataComponentsFactoryArgs {
	return factory.DataComponentsFactoryArgs{
		Config:             getGeneralConfig(),
		EconomicsData:      &economics.EconomicsData{},
		ShardCoordinator:   mock.NewMultiShardsCoordinatorMock(2),
		Core:               getCoreComponents(),
		PathManager:        &mock.PathManagerStub{},
		EpochStartNotifier: &mock.EpochStartNotifierStub{},
		CurrentEpoch:       0,
	}
}

func getGeneralConfig() *config.Config {
	return &config.Config{
		StoragePruning: config.StoragePruningConfig{
			Enabled:             false,
			FullArchive:         true,
			NumEpochsToKeep:     3,
			NumActivePersisters: 3,
		},
		TxDataPool: config.CacheConfig{
			Size: 10000, Type: "LRU", Shards: 1,
		},
		UnsignedTransactionDataPool: config.CacheConfig{
			Size: 10000, Type: "LRU", Shards: 1,
		},
		RewardTransactionDataPool: config.CacheConfig{
			Size: 10000, Type: "LRU", Shards: 1,
		},
		HeadersPoolConfig: config.HeadersPoolConfig{
			MaxHeadersPerShard:            100,
			NumElementsToRemoveOnEviction: 1,
		},
		TxBlockBodyDataPool: config.CacheConfig{
			Size: 10000, Type: "LRU", Shards: 1,
		},
		PeerBlockBodyDataPool: config.CacheConfig{
			Size: 10000, Type: "LRU", Shards: 1,
		},
		TrieNodesDataPool: config.CacheConfig{
			Size: 10000, Type: "LRU", Shards: 1,
		},
		TxStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "Transactions",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		MiniBlocksStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "MiniBlocks",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		ShardHdrNonceHashStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "ShardHdrHashNonce",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		MetaBlockStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "MetaBlock",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		MetaHdrNonceHashStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "MetaHdrHashNonce",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		UnsignedTransactionStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "UnsignedTransactions",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		RewardTxStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "RewardTransactions",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		BlockHeaderStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "BlockHeaders",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		Heartbeat: config.HeartbeatConfig{
			HeartbeatStorage: config.StorageConfig{
				Cache: config.CacheConfig{
					Size: 10000, Type: "LRU", Shards: 1,
				},
				DB: config.DBConfig{
					FilePath:          "HeartbeatStorage",
					Type:              string(storageUnit.MemoryDB),
					BatchDelaySeconds: 30,
					MaxBatchSize:      6,
					MaxOpenFiles:      10,
				},
			},
		},
		StatusMetricsStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "StatusMetricsStorageDB",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		PeerBlockBodyStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "PeerBlocks",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 30,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		BootstrapStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Size: 10000, Type: "LRU", Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "BootstrapData",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 1,
				MaxBatchSize:      6,
				MaxOpenFiles:      10,
			},
		},
		TxLogsStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Type:   "LRU",
				Size:   1000,
				Shards: 1,
			},
			DB: config.DBConfig{
				FilePath:          "Logs",
				Type:              string(storageUnit.MemoryDB),
				BatchDelaySeconds: 2,
				MaxBatchSize:      100,
				MaxOpenFiles:      10,
			},
		},
	}
}
