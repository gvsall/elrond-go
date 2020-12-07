package factory

import (
	"fmt"
	"path"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/outport"
	"github.com/ElrondNetwork/elrond-go/outport/drivers/elastic"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/elastic/go-elasticsearch/v7"
)

const (
	withKibanaFolder = "withKibana"
	noKibanaFolder   = "noKibana"
)

// ArgsElasticDriverFactory holds all dependencies required by the data indexer factory in order to create
// new instances
type ArgsElasticDriverFactory struct {
	Enabled                  bool
	IndexerCacheSize         int
	ShardCoordinator         sharding.Coordinator
	Url                      string
	UserName                 string
	Password                 string
	Marshalizer              marshal.Marshalizer
	Hasher                   hashing.Hasher
	AddressPubkeyConverter   core.PubkeyConverter
	ValidatorPubkeyConverter core.PubkeyConverter
	TemplatesPath            string
	Options                  *elastic.Options
	EnabledIndexes           []string
	Denomination             int
	AccountsDB               state.AccountsAdapter
	FeeConfig                *config.FeeSettings
	IsInImportDBMode         bool
}

// NewElasticClient will create a new instance of elastic client
func NewElasticClient(args *ArgsElasticDriverFactory) (outport.Driver, error) {
	err := checkDataIndexerParams(args)
	if err != nil {
		return nil, err
	}

	elasticProcessor, err := createElasticProcessor(args)
	if err != nil {
		return nil, err
	}

	dispatcher, err := elastic.NewDataDispatcher(args.IndexerCacheSize)
	if err != nil {
		return nil, err
	}

	dispatcher.StartIndexData()

	arguments := elastic.ArgDataIndexer{
		Marshalizer:      args.Marshalizer,
		Options:          args.Options,
		ShardCoordinator: args.ShardCoordinator,
		ElasticProcessor: elasticProcessor,
		DataDispatcher:   dispatcher,
	}

	return elastic.NewDataIndexer(arguments)
}

func createDatabaseClient(url, userName, password string) (elastic.DatabaseClientHandler, error) {
	return elastic.NewElasticClient(elasticsearch.Config{
		Addresses: []string{url},
		Username:  userName,
		Password:  password,
	})
}

func createElasticProcessor(args *ArgsElasticDriverFactory) (elastic.ElasticProcessor, error) {
	databaseClient, err := createDatabaseClient(args.Url, args.UserName, args.Password)
	if err != nil {
		return nil, err
	}

	var templatesPath string
	if args.Options.UseKibana {
		templatesPath = path.Join(args.TemplatesPath, withKibanaFolder)
	} else {
		templatesPath = path.Join(args.TemplatesPath, noKibanaFolder)
	}

	indexTemplates, indexPolicies, err := elastic.GetElasticTemplatesAndPolicies(templatesPath, args.Options.UseKibana)
	if err != nil {
		return nil, err
	}

	enabledIndexesMap := make(map[string]struct{})
	for _, index := range args.EnabledIndexes {
		enabledIndexesMap[index] = struct{}{}
	}
	if len(enabledIndexesMap) == 0 {
		return nil, elastic.ErrEmptyEnabledIndexes
	}

	esIndexerArgs := elastic.ArgElasticProcessor{
		IndexTemplates:           indexTemplates,
		IndexPolicies:            indexPolicies,
		Marshalizer:              args.Marshalizer,
		Hasher:                   args.Hasher,
		AddressPubkeyConverter:   args.AddressPubkeyConverter,
		ValidatorPubkeyConverter: args.ValidatorPubkeyConverter,
		Options:                  args.Options,
		DBClient:                 databaseClient,
		EnabledIndexes:           enabledIndexesMap,
		AccountsDB:               args.AccountsDB,
		Denomination:             args.Denomination,
		FeeConfig:                args.FeeConfig,
		IsInImportDBMode:         args.IsInImportDBMode,
		ShardCoordinator:         args.ShardCoordinator,
	}

	return elastic.NewElasticProcessor(esIndexerArgs)
}

func checkDataIndexerParams(arguments *ArgsElasticDriverFactory) error {
	if arguments == nil {
		return outport.ErrNilArgsElasticDriverFactory
	}
	if arguments.IndexerCacheSize < 0 {
		return elastic.ErrNegativeCacheSize
	}
	if check.IfNil(arguments.AddressPubkeyConverter) {
		return fmt.Errorf("%w when setting AddressPubkeyConverter in indexer", outport.ErrNilPubkeyConverter)
	}
	if check.IfNil(arguments.ValidatorPubkeyConverter) {
		return fmt.Errorf("%w when setting ValidatorPubkeyConverter in indexer", outport.ErrNilPubkeyConverter)
	}
	if arguments.Url == "" {
		return outport.ErrNilUrl
	}
	if check.IfNil(arguments.Marshalizer) {
		return outport.ErrNilMarshalizer
	}
	if check.IfNil(arguments.Hasher) {
		return outport.ErrNilHasher
	}
	if arguments.FeeConfig == nil {
		return outport.ErrNilFeeConfig
	}
	if check.IfNil(arguments.AccountsDB) {
		return outport.ErrNilAccountsDB
	}

	return nil
}