package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPvtzZoneRecord_basic(t *testing.T) {
	var v pvtz.Record

	resourceId := "alicloud_pvtz_zone_record.default"
	ra := resourceAttrInit(resourceId, pvtzZoneRecordBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneRecordConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_record": "www",
					"type":            "A",
					"value":           "2.2.2.2",
					"zone_id":         "${alicloud_pvtz_zone.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "2.2.2.2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "pvtz_zone_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "pvtz_zone_remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "TXT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "TXT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "2.2.2.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "2.2.2.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":     "MX",
					"value":    "bbb.test.com",
					"priority": "2",
					"ttl":      REMOVEKEY,
					"remark":   "pvtz_zone_describe",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "MX",
						"value":    "bbb.test.com",
						"priority": "2",
						"ttl":      "60",
						"remark":   "pvtz_zone_describe",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":     "A",
					"value":    "2.2.2.2",
					"priority": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "A",
						"value":    "2.2.2.2",
						"priority": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudPvtzZoneRecord_multi(t *testing.T) {
	var v pvtz.Record

	resourceId := "alicloud_pvtz_zone_record.default.4"
	ra := resourceAttrInit(resourceId, pvtzZoneRecordBasicMap)

	serviceFunc := func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePvtzZoneRecordConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_record": "www",
					"type":            "A",
					"value":           "2.2.2.${count.index}",
					"zone_id":         "${alicloud_pvtz_zone.default.id}",
					"count":           "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func resourcePvtzZoneRecordConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "default" {
		name = "%s"
	}
	`, name)
}

var pvtzZoneRecordBasicMap = map[string]string{
	"resource_record": "www",
	"type":            "A",
	"value":           CHECKSET,
	"zone_id":         CHECKSET,
	"ttl":             "60",
}
