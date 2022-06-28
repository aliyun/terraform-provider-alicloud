package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudExpressConnectAccessPointsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_express_connect_access_points.default"
	name := fmt.Sprintf("tf-testacc-expressConnectAccessPoints%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectAccessPointsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": getAccessPointNamePrefix(),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{getAccessPointId()},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{getAccessPointId()},
			"status": "recommended",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"fake"},
			"status": "full",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": getAccessPointNamePrefix(),
			"ids":        []string{getAccessPointId()},
			"status":     "recommended",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake",
			"ids":        []string{"fake"},
			"status":     "full",
		}),
	}

	var existExpressConnectAccessPointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"points.#":                               "1",
			"points.0.access_point_feature_models.#": CHECKSET,
			"points.0.access_point_feature_models.0.feature_key":   CHECKSET,
			"points.0.access_point_feature_models.0.feature_value": CHECKSET,
			"points.0.id":                 CHECKSET,
			"points.0.access_point_id":    getAccessPointId(),
			"points.0.access_point_name":  CHECKSET,
			"points.0.attached_region_no": CHECKSET,
			"points.0.description":        CHECKSET,
			"points.0.host_operator":      CHECKSET,
			"points.0.location":           "",
			"points.0.status":             "recommended",
			"points.0.type":               "VPC",
		}
	}

	var fakeExpressConnectAccessPointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"points.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var ExpressConnectAccessPointsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressConnectAccessPointsMapFunc,
		fakeMapFunc:  fakeExpressConnectAccessPointsMapFunc,
	}

	ExpressConnectAccessPointsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, statusConf, allConf)
}

func dataSourceExpressConnectAccessPointsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		`, name)
}

var AccessPointNamePrefixMap = map[string]string{
	string(connectivity.Hangzhou):     "杭州-余杭-B",
	string(connectivity.Beijing):      "北京-朝阳-A",
	string(connectivity.Shanghai):     "上海-宝山-D",
	string(connectivity.APSouthEast1): "新加坡-A",
	string(connectivity.EUCentral1):   "^欧洲-法兰克福-B",
}
var AccessPointIdMap = map[string]string{
	string(connectivity.Hangzhou):     "ap-cn-hangzhou-yh-B",
	string(connectivity.Beijing):      "ap-cn-beijing-cy-A",
	string(connectivity.Shanghai):     "ap-cn-shanghai-pd-D",
	string(connectivity.APSouthEast1): "ap-sg-singpore-A",
	string(connectivity.EUCentral1):   "ap-eu-frankfurt-B",
}

func getAccessPointNamePrefix() string {
	return AccessPointNamePrefixMap[os.Getenv("ALICLOUD_REGION")]
}
func getAccessPointId() string {
	return AccessPointIdMap[os.Getenv("ALICLOUD_REGION")]
}
