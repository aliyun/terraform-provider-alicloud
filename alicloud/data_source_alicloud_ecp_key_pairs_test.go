package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcpKeyPairsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecp_key_pairs.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcpKeyPairsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcpKeyPairsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecp_key_pair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecp_key_pair.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecp_key_pair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_ecp_key_pair.default.id}-fake"},
		}),
	}
	var existEcpKeyPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         CHECKSET,
			"names.#":                       "1",
			"names.0":                       name,
			"pairs.#":                       "1",
			"pairs.0.id":                    CHECKSET,
			"pairs.0.key_pair_finger_print": CHECKSET,
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

	var EcpKeyPairsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcpKeyPairsMapFunc,
		fakeMapFunc:  fakeEcpKeyPairsMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EcpSupportRegions)
	}

	EcpKeyPairsInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceEcpKeyPairsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ecp_key_pair" "default" {
		key_pair_name              = "%s"
		public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
	}`, name)
}
