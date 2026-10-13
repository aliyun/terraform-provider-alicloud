package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudThreatDetectionVulAutoConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_vul_auto_config.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudThreatDetectionVulAutoConfigMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionVulAutoConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccThreatDetectionVulAutoConfig-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudThreatDetectionVulAutoConfigBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":              "vul",
					"start_time":        1700000000,
					"all_uuid":          1,
					"need_snapshot":     0,
					"enable":            1,
					"period_unit":       "day",
					"necessity":         "asap",
					"target_start_time": 0,
					"target_end_time":   23,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":              "vul",
						"start_time":        "1700000000",
						"all_uuid":          "1",
						"need_snapshot":     "0",
						"enable":            "1",
						"period_unit":       "day",
						"necessity":         "asap",
						"target_start_time": "0",
						"target_end_time":   "23",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":              "cve",
					"start_time":        1800000000,
					"all_uuid":          0,
					"need_snapshot":     1,
					"enable":            0,
					"period_unit":       "week",
					"necessity":         "serious",
					"target_start_time": 6,
					"target_end_time":   18,
					"snapshot_name":     name,
					"snapshot_time":     24,
					"rules":             "[]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":              "cve",
						"start_time":        "1800000000",
						"all_uuid":          "0",
						"need_snapshot":     "1",
						"enable":            "0",
						"period_unit":       "week",
						"necessity":         "serious",
						"target_start_time": "6",
						"target_end_time":   "18",
						"snapshot_name":     name,
						"snapshot_time":     "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":              "cve",
					"start_time":        1800000000,
					"all_uuid":          0,
					"need_snapshot":     1,
					"enable":            1,
					"period_unit":       "week",
					"necessity":         "serious",
					"target_start_time": 6,
					"target_end_time":   18,
					"snapshot_name":     name,
					"snapshot_time":     24,
					"rules":             "[]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudThreatDetectionVulAutoConfig_enableOnly(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_vul_auto_config.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudThreatDetectionVulAutoConfigMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionVulAutoConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccThreatDetectionVulAutoConfigEnable-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudThreatDetectionVulAutoConfigBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":          "vul",
					"start_time":    1700000000,
					"all_uuid":      1,
					"need_snapshot": 0,
					"enable":        1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":          "vul",
					"start_time":    1700000000,
					"all_uuid":      1,
					"need_snapshot": 0,
					"enable":        0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "0",
					}),
				),
			},
		},
	})
}

var resourceAlicloudThreatDetectionVulAutoConfigMap = map[string]string{
	"region_id": CHECKSET,
}

func resourceAlicloudThreatDetectionVulAutoConfigBasicDependence(name string) string {
	return ""
}
