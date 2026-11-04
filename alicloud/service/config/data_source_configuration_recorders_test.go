package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// alicloud_config_configuration_recorder used to open CloudConfig service.
// There can stop the service manually first before running the testcase.
func SkipTestAccAliCloudConfigConfigurationRecordersDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_config_configuration_recorders.example"

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigConfigurationRecordersSourceConfig(rand, map[string]string{}),
	}

	var ExistConfigConfigurationRecordersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"recorders.#":                            "1",
			"recorders.0.id":                         acctest.CHECKSET,
			"recorders.0.account_id":                 acctest.CHECKSET,
			"recorders.0.organization_enable_status": acctest.CHECKSET,
			"recorders.0.organization_master_id":     acctest.CHECKSET,
			"recorders.0.resource_types.#":           "2",
			"recorders.0.status":                     "REGISTERED",
		}
	}

	var FakeConfigConfigurationRecordersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"recorders.#": "0",
		}
	}

	var configConfigurationRecordersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: ExistConfigConfigurationRecordersMapFunc,
		FakeMapFunc:  FakeConfigConfigurationRecordersMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
	}

	configConfigurationRecordersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)

}

// TestAccAliCloudConfigConfigurationRecordersDataSource exercises the read path of the
// alicloud_config_configuration_recorders data source against the account's existing
// configuration recorder. Creating an alicloud_config_configuration_recorder resource
// requires an enterprise account (PreCheckEnterpriseAccountEnabled), which is out of
// scope for this migration's test account. DescribeConfigurationRecorder returns the
// account's single recorder once Cloud Config is enabled, so a read-only query (no
// resource creation, no enterprise account, no extra dependencies) is a faithful,
// self-contained representative ACC for this migrated data source.
func TestAccAliCloudConfigConfigurationRecordersDataSource(t *testing.T) {
	resourceId := "data.alicloud_config_configuration_recorders.example"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `data "alicloud_config_configuration_recorders" "example" { /* read-only representative ACC against the account's existing recorder */ }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "recorders.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "recorders.0.id"),
					resource.TestCheckResourceAttrSet(resourceId, "recorders.0.account_id"),
				),
			},
		},
	})
}

func testAccCheckAlicloudConfigConfigurationRecordersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccConfigConfigurationRecorders%d"
}

resource "alicloud_config_configuration_recorder" "example" {
 resource_types = ["ACS::ECS::Disk","ACS::ECS::Instance"]
}

data "alicloud_config_configuration_recorders" "example"{
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
