package alicloud

type MongoDBStorageEngine string

const (
	WiredTiger = MongoDBStorageEngine("WiredTiger")
	RocksDB    = MongoDBStorageEngine("RocksDB")
)
