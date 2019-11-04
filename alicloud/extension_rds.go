package alicloud

type Engine string

const (
	MySQL      = Engine("MySQL")
	SQLServer  = Engine("SQLServer")
	PPAS       = Engine("PPAS")
	PostgreSQL = Engine("PostgreSQL")
)

type DBAccountPrivilege string

const (
	ReadOnly  = DBAccountPrivilege("ReadOnly")
	ReadWrite = DBAccountPrivilege("ReadWrite")
)

type DBAccountType string

const (
	DBAccountNormal = DBAccountType("Normal")
	DBAccountSuper  = DBAccountType("Super")
)

var WEEK_ENUM = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

var BACKUP_TIME = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}

var CHARACTER_SET_NAME = []string{
	"utf8", "gbk", "latin1", "utf8mb4", "Mohawk_100_BIN",
	"Chinese_PRC_CI_AS", "Chinese_PRC_CS_AS", "SQL_Latin1_General_CP1_CI_AS", "SQL_Latin1_General_CP1_CS_AS", "Chinese_PRC_BIN",
}

type KVStoreInstanceType string

const (
	KVStoreRedis    = KVStoreInstanceType("Redis")
	KVStoreMemcache = KVStoreInstanceType("Memcache")
)

type KVStoreEngineVersion string

const (
	KVStore2Dot8 = KVStoreEngineVersion("2.8")
	KVStore4Dot0 = KVStoreEngineVersion("4.0")
)
