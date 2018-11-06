package tablestore

import (
	"errors"
)

var (
	errMissMustHeader = func(header string) error {
		return errors.New("[tablestore] miss must header: " + header)
	}
	errTableNameTooLong = func(name string) error {
		return errors.New("[tablestore] table name: \"" + name + "\" too long")
	}

	errInvalidPartitionType    = errors.New("[tablestore] invalid partition key")
	errMissPrimaryKey          = errors.New("[tablestore] missing primary key")
	errPrimaryKeyTooMuch       = errors.New("[tablestore] primary key too much")
	errMultiDeleteRowsTooMuch  = errors.New("[tablestore] multi delete rows too much")
	errCreateTableNoPrimaryKey = errors.New("[tablestore] create table no primary key")
	errUnexpectIoEnd           = errors.New("[tablestore] unexpect io end")
	errTag                     = errors.New("[tablestore] unexpect tag")
	errNoChecksum              = errors.New("[tablestore] expect checksum")
	errChecksum                = errors.New("[tablestore] checksum failed")
	errInvalidInput            = errors.New("[tablestore] invalid input")
)

const (
	ROW_OPERATION_CONFLICT   = "OTSRowOperationConflict"
	NOT_ENOUGH_CAPACITY_UNIT = "OTSNotEnoughCapacityUnit"
	TABLE_NOT_READY          = "OTSTableNotReady"
	PARTITION_UNAVAILABLE    = "OTSPartitionUnavailable"
	SERVER_BUSY              = "OTSServerBusy"
	QUOTA_EXHAUSTED          = "OTSQuotaExhausted"

	STORAGE_TIMEOUT       = "OTSTimeout"
	SERVER_UNAVAILABLE    = "OTSServerUnavailable"
	INTERNAL_SERVER_ERROR = "OTSInternalServerError"
)
