package alidns_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudTxtGuidDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_dns_domain_txt_guid.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceTxtGuidConfigDependence)
	nameAndTypeConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"domain_name": "test111.abc",
			"type":        "ADD_SUB_DOMAIN",
		}),
	}

	var existTxtGuidMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rr":    acctest.CHECKSET,
			"value": acctest.CHECKSET,
		}
	}

	var fakeTxtGuidMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rr":    "",
			"value": "",
		}
	}

	var txtGuidCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existTxtGuidMapFunc,
		FakeMapFunc:  fakeTxtGuidMapFunc,
	}

	txtGuidCheckInfo.DataSourceTestCheck(t, rand, nameAndTypeConfig)
}
func dataSourceTxtGuidConfigDependence(name string) string {
	return ""
}
