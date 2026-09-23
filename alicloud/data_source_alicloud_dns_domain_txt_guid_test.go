package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudTxtGuidDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_dns_domain_txt_guid.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceTxtGuidConfigDependence)
	nameAndTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "test111.abc",
			"type":        "ADD_SUB_DOMAIN",
		}),
	}

	var existTxtGuidMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rr":    CHECKSET,
			"value": CHECKSET,
		}
	}

	var fakeTxtGuidMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rr":    "",
			"value": "",
		}
	}

	var txtGuidCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existTxtGuidMapFunc,
		fakeMapFunc:  fakeTxtGuidMapFunc,
	}

	txtGuidCheckInfo.dataSourceTestCheck(t, rand, nameAndTypeConfig)
}
func dataSourceTxtGuidConfigDependence(name string) string {
	return ""
}
