package ecp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEcpKeyPairsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecp_key_pairs.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcpKeyPairsTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEcpKeyPairsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecp_key_pair.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecp_key_pair.default.id}-fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecp_key_pair.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_ecp_key_pair.default.id}-fake"},
		}),
	}
	var existEcpKeyPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         acctest.CHECKSET,
			"names.#":                       "1",
			"names.0":                       name,
			"pairs.#":                       "1",
			"pairs.0.id":                    acctest.CHECKSET,
			"pairs.0.key_pair_finger_print": acctest.CHECKSET,
			"pairs.0.key_pair_name":         name,
		}
	}

	var fakeEcpKeyPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"pairs.#": "0",
		}
	}

	var EcpKeyPairsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existEcpKeyPairsMapFunc,
		FakeMapFunc:  fakeEcpKeyPairsMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EcpSupportRegions)
	}

	EcpKeyPairsInfo.DataSourceTestCheckWithPreCheck(t, 0, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceEcpKeyPairsDependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_ecp_key_pair" "default" {
			key_pair_name              = "%s"
			public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
		}`, name)
}
