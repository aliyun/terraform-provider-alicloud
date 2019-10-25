package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudEmrInstanceTypesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceID := "data.alicloud_emr_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%semr-instance-types%v.abc", defaultRegionToTest, rand),
		dataSourceEmrInstanceTypesConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"destination_resource": "InstanceType",
			"cluster_type":         "HADOOP",
			"instance_charge_type": "PostPaid",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"destination_resource": "Zone",
			"cluster_type":         "HADOOP",
			"instance_charge_type": "PostPaid",
		}),
	}

	var (
		existEmrInstanceTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":           CHECKSET,
				"ids.0":           CHECKSET,
				"types.#":         CHECKSET,
				"types.0.zone_id": CHECKSET,
				"types.0.id":      CHECKSET,
			}
		}

		fakeEmrInstanceTypesMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":   "0",
				"types.#": "0",
			}
		}

		emrInstanceTypesCheckInfo = dataSourceAttr{
			resourceId:   resourceID,
			existMapFunc: existEmrInstanceTypesMapFunc,
			fakeMapFunc:  fakeEmrInstanceTypesMapFunc,
		}
	)

	emrInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourceEmrInstanceTypesConfigDependence(name string) string {
	return ``
}
