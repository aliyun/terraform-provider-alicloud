package polardb_test

import (
	"testing"

	polardb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
)

func TestUnitPolarDBLakebaseTenantTokenIsSensitive(t *testing.T) {
	dataSource := polardb.DataSourceLakebaseTenantToken()
	if !dataSource.Schema["token"].Sensitive {
		t.Fatal("token must be sensitive")
	}
}
