package alicloud

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func getPrimaryKeyType(primaryKeyType string) tablestore.PrimaryKeyType {
	var keyType tablestore.PrimaryKeyType
	t := PrimaryKeyTypeString(primaryKeyType)
	switch t {
	case IntegerType:
		keyType = tablestore.PrimaryKeyType_INTEGER
	case StringType:
		keyType = tablestore.PrimaryKeyType_STRING
	case BinaryType:
		keyType = tablestore.PrimaryKeyType_BINARY
	}
	return keyType
}

func describeOtsTable(tableName string, meta interface{}) (*tablestore.DescribeTableResponse, error) {
	client := meta.(*AliyunClient).otsconn

	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = tableName

	return client.DescribeTable(describeTableReq)
}

func deleteOtsTable(tableName string, meta interface{}) (bool, error) {
	client := meta.(*AliyunClient).otsconn

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = tableName
	_, err := client.DeleteTable(deleteReq)

	describ, _ := describeOtsTable(tableName, meta)

	if describ.TableMeta != nil {
		return false, err
	}

	return true, err
}

// Convert tablestore.PrimaryKeyType to PrimaryKeyTypeString
func convertPrimaryKeyType(t tablestore.PrimaryKeyType) PrimaryKeyTypeString {
	var typeString PrimaryKeyTypeString
	switch t {
	case tablestore.PrimaryKeyType_INTEGER:
		typeString = IntegerType
	case tablestore.PrimaryKeyType_BINARY:
		typeString = BinaryType
	case tablestore.PrimaryKeyType_STRING:
		typeString = StringType
	}
	return typeString
}
