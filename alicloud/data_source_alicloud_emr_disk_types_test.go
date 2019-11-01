package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudEmrDiskTypesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceID := "data.alicloud_emr_disk_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%semr-disk-types%v.abc", defaultRegionToTest, rand),
		dataSourceEmrDiskTypesConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"destination_resource": "DataDisk",
			"cluster_type":         "HADOOP",
			"instance_charge_type": "PostPaid",
			"instance_type":        "ecs.g5.xlarge",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"destination_resource": "Zone",
			"cluster_type":         "HADOOP",
			"instance_charge_type": "PostPaid",
			"instance_type":        "ecs.g5.xlarge",
		}),
	}

	var (
		existEmrDiskTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":         CHECKSET,
				"ids.0":         CHECKSET,
				"types.#":       CHECKSET,
				"types.0.min":   CHECKSET,
				"types.0.max":   CHECKSET,
				"types.0.value": CHECKSET,
			}
		}

		fakeEmrDiskTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":   "0",
				"types.#": "0",
			}
		}

		emrDiskTypesCheckInfo = dataSourceAttr{
			resourceId:   resourceID,
			existMapFunc: existEmrDiskTypesMapFunc,
			fakeMapFunc:  fakeEmrDiskTypesMapFunc,
		}
	)

	emrDiskTypesCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourceEmrDiskTypesConfigDependence(name string) string {
	return ``
}
