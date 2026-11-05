package emr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEmrInstanceTypesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceID := "data.alicloud_emr_instance_types.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%semr-instance-types%v.abc", acctest.DefaultRegionToTest, rand),
		dataSourceEmrInstanceTypesConfigDependence)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"destination_resource":  "InstanceType",
			"cluster_type":          "HADOOP",
			"support_local_storage": "false",
			"instance_charge_type":  "PostPaid",
			"support_node_type":     []string{"CORE"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"destination_resource":  "Zone",
			"cluster_type":          "HADOOP",
			"support_local_storage": "true",
			"instance_charge_type":  "PostPaid",
		}),
	}

	var (
		existEmrInstanceTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":           acctest.CHECKSET,
				"ids.0":           acctest.CHECKSET,
				"types.#":         acctest.CHECKSET,
				"types.0.zone_id": acctest.CHECKSET,
				"types.0.id":      acctest.CHECKSET,
			}
		}

		fakeEmrInstanceTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":   "0",
				"types.#": "0",
			}
		}

		emrInstanceTypesCheckInfo = acctest.DataSourceAttr{
			ResourceId:   resourceID,
			ExistMapFunc: existEmrInstanceTypesMapFunc,
			FakeMapFunc:  fakeEmrInstanceTypesMapFunc,
		}
	)

	emrInstanceTypesCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

func dataSourceEmrInstanceTypesConfigDependence(name string) string {
	return ``
}
