package factory_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/factory"
	"github.com/stretchr/testify/require"
)

const testHasher = "blake2b"
const testMarshalizer = "json"

func TestNewCoreComponentsFactory_NilConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = nil
	ccf, err := factory.NewCoreComponentsFactory(args)

	require.Nil(t, ccf)
	require.Equal(t, factory.ErrNilConfiguration, err)
}

func TestNewCoreComponentsFactory_OkValuesShouldWork(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	ccf, err := factory.NewCoreComponentsFactory(args)

	require.NoError(t, err)
	require.NotNil(t, ccf)
}

func TestCoreComponentsFactory_CreateCoreComponents_NoHasherConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Marshalizer: config.MarshalizerConfig{
			Type:           testMarshalizer,
			SizeCheckDelta: 0,
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrHasherCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_InvalidHasherConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Marshalizer: config.MarshalizerConfig{
			Type:           testMarshalizer,
			SizeCheckDelta: 0,
		},
		Hasher: config.TypeConfig{
			Type: "invalid_type",
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrHasherCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_NoInternalMarshalizerConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Hasher: config.TypeConfig{
			Type: testHasher,
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrMarshalizerCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_InvalidInternalMarshalizerConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Marshalizer: config.MarshalizerConfig{
			Type:           "invalid_marshalizer_type",
			SizeCheckDelta: 0,
		},
		Hasher: config.TypeConfig{
			Type: testHasher,
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrMarshalizerCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_NoVmMarshalizerConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Hasher: config.TypeConfig{
			Type: testHasher,
		},
		Marshalizer: config.MarshalizerConfig{
			Type:           testMarshalizer,
			SizeCheckDelta: 0,
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrMarshalizerCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_InvalidVmMarshalizerConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Marshalizer: config.MarshalizerConfig{
			Type:           testMarshalizer,
			SizeCheckDelta: 0,
		},
		Hasher: config.TypeConfig{
			Type: testHasher,
		},
		VmMarshalizer: config.TypeConfig{
			Type: "invalid",
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrMarshalizerCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_NoTxSignMarshalizerConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Hasher: config.TypeConfig{
			Type: testHasher,
		},
		Marshalizer: config.MarshalizerConfig{
			Type:           testMarshalizer,
			SizeCheckDelta: 0,
		},
		VmMarshalizer: config.TypeConfig{
			Type: testMarshalizer,
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrMarshalizerCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_InvalidTxSignMarshalizerConfigShouldErr(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	args.Config = &config.Config{
		Marshalizer: config.MarshalizerConfig{
			Type:           testMarshalizer,
			SizeCheckDelta: 0,
		},
		Hasher: config.TypeConfig{
			Type: testHasher,
		},
		VmMarshalizer: config.TypeConfig{
			Type: testMarshalizer,
		},
		TxSignMarshalizer: config.TypeConfig{
			Type: "invalid",
		},
	}
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.Nil(t, cc)
	require.True(t, errors.Is(err, factory.ErrMarshalizerCreation))
}

func TestCoreComponentsFactory_CreateCoreComponents_ShouldWork(t *testing.T) {
	t.Parallel()

	args := getCoreArgs()
	ccf, _ := factory.NewCoreComponentsFactory(args)

	cc, err := ccf.Create()
	require.NoError(t, err)
	require.NotNil(t, cc)
}

func getCoreArgs() factory.CoreComponentsFactoryArgs {
	return factory.CoreComponentsFactoryArgs{
		Config: &config.Config{
			Marshalizer: config.MarshalizerConfig{
				Type:           testMarshalizer,
				SizeCheckDelta: 0,
			},
			Hasher: config.TypeConfig{
				Type: testHasher,
			},
			VmMarshalizer: config.TypeConfig{
				Type: testMarshalizer,
			},
			TxSignMarshalizer: config.TypeConfig{
				Type: testMarshalizer,
			},
		},
		ShardId: "0",
		ChainID: []byte("chainID"),
	}
}
