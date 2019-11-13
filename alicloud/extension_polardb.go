package alicloud

type PolarDBAccountPrivilege string

const (
	PolarDBReadOnly  = PolarDBAccountPrivilege("ReadOnly")
	PolarDBReadWrite = PolarDBAccountPrivilege("ReadWrite")
)

type PolarDBAccountType string

const (
	PolarDBAccountNormal = PolarDBAccountType("Normal")
	PolarDBAccountSuper  = PolarDBAccountType("Super")
)

var POLAR_DB_ACCOUNT_PRIVILEGE_NAME = []string{
	"ReadOnly", "ReadWrite", "DMLOnly", "DDLOnly",
}

var POLAR_DB_CHARACTER_SET_NAME = []string{
	"utf8", "gbk", "latin1", "utf8mb4", "Chinese_PRC_CI_AS", "Chinese_PRC_CS_AS",
	"SQL_Latin1_General_CP1_CI_AS", "SQL_Latin1_General_CP1_CS_AS", "Chinese_PRC_BIN",
}
