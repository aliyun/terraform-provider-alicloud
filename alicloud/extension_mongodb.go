package alicloud

type MongoDBStorageEngine string

const (
	WiredTiger = MongoDBStorageEngine("WiredTiger")
	RocksDB    = MongoDBStorageEngine("RocksDB")
)

type MongoDBShardingNodeType string

const (
	MongoDBShardingNodeMongos = MongoDBShardingNodeType("mongos")
	MongoDBShardingNodeShard  = MongoDBShardingNodeType("shard")
)

type MongoDBInstanceType string

const (
	MongoDBSharding  = MongoDBInstanceType("sharding")
	MongoDBReplicate = MongoDBInstanceType("replicate")
)
