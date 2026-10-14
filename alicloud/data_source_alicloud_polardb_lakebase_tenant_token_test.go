package alicloud

import "testing"

func TestUnitPolarDBLakebaseTenantTokenIsSensitive(t *testing.T) {
	dataSource := dataSourceAlicloudPolarDBLakebaseTenantToken()
	if !dataSource.Schema["token"].Sensitive {
		t.Fatal("token must be sensitive")
	}
}
