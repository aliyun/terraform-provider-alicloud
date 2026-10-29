package rds

// Engine and its engine constants are the package-rds copies of the root
// alicloud/extension_rds.go symbols. They are root-native (not leaf-package
// re-exports) and are shared across several products (adb, polardb, mongodb,
// kvstore, clickhouse), so per stock-method-disposition rule #4 the root keeps
// its copy and RDS keeps this behavior-preserving local copy rather than
// importing the root package. BACKUP_TIME likewise moves with the RDS backup
// schema.

type Engine string

const (
	MySQL      = Engine("MySQL")
	SQLServer  = Engine("SQLServer")
	PPAS       = Engine("PPAS")
	PostgreSQL = Engine("PostgreSQL")
	MariaDB    = Engine("MariaDB")
)

var BACKUP_TIME = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}
