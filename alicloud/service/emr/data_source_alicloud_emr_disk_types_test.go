package emr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEmrDiskTypesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceID := "data.alicloud_emr_disk_types.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%semr-disk-types%v.abc", acctest.DefaultRegionToTest, rand),
		dataSourceEmrDiskTypesConfigDependence)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"destination_resource": "DataDisk",
			"cluster_type":         "HADOOP",
			"instance_charge_type": "PostPaid",
			"instance_type":        "ecs.g5.xlarge",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"destination_resource": "Zone",
			"cluster_type":         "HADOOP",
			"instance_charge_type": "PostPaid",
			"instance_type":        "ecs.g5.xlarge",
		}),
	}

	var (
		existEmrDiskTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":         acctest.CHECKSET,
				"ids.0":         acctest.CHECKSET,
				"types.#":       acctest.CHECKSET,
				"types.0.min":   acctest.CHECKSET,
				"types.0.max":   acctest.CHECKSET,
				"types.0.value": acctest.CHECKSET,
			}
		}

		fakeEmrDiskTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":   "0",
				"types.#": "0",
			}
		}

		emrDiskTypesCheckInfo = acctest.DataSourceAttr{
			ResourceId:   resourceID,
			ExistMapFunc: existEmrDiskTypesMapFunc,
			FakeMapFunc:  fakeEmrDiskTypesMapFunc,
		}
	)

	emrDiskTypesCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

func dataSourceEmrDiskTypesConfigDependence(name string) string {
	return ``
}
