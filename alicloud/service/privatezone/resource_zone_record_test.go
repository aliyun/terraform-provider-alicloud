package privatezone_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/privatezone"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPvtzZoneRecord_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_pvtz_zone_record.default"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneRecordBasicMap)

	serviceFunc := func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeZoneRecord")

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneRecordConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_client_ip": "172.16.0.0/16",
					"rr":             "www",
					"type":           "A",
					"value":          "2.2.2.2",
					"zone_id":        "${alicloud_pvtz_zone.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "2.2.2.2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_client_ip"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_client_ip": "172.16.0.0/16-upd",
					"status":         "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rr": "aaa",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rr": "aaa",
					}),
				),
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
					"lang": "en",
					"ttl":  "30",
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
					"ttl":      acctest.REMOVEKEY,
					"remark":   "pvtz_zone_describe",
					"rr":       "www",
					"status":   "ENABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "MX",
						"value":    "bbb.test.com",
						"priority": "2",
						"ttl":      "60",
						"remark":   "pvtz_zone_describe",
						"rr":       "www",
						"status":   "ENABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":     "A",
					"value":    "2.2.2.2",
					"priority": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "A",
						"value":    "2.2.2.2",
						"priority": acctest.REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPvtzZoneRecord_multi(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_pvtz_zone_record.default.4"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneRecordBasicMap)

	serviceFunc := func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeZoneRecord")

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneRecordConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rr":      "www",
					"type":    "A",
					"value":   "2.2.2.${count.index}",
					"zone_id": "${alicloud_pvtz_zone.default.id}",
					"count":   "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var pvtzZoneRecordAAAAMap = map[string]string{
	"rr":      "www",
	"type":    "AAAA",
	"value":   acctest.CHECKSET,
	"zone_id": acctest.CHECKSET,
	"ttl":     "60",
}

func TestAccAliCloudPvtzZoneRecord_aaaa(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_pvtz_zone_record.default"
	ra := acctest.ResourceAttrInit(resourceId, pvtzZoneRecordAAAAMap)

	serviceFunc := func() interface{} {
		return privatezone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeZoneRecord")

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%d.test.com", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourcePvtzZoneRecordConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rr":      "www",
					"type":    "AAAA",
					"value":   "2001:db8::1",
					"zone_id": "${alicloud_pvtz_zone.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":  "AAAA",
						"value": "2001:db8::1",
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
					"value": "2001:db8::2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": "2001:db8::2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rr": "ipv6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rr": "ipv6",
					}),
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
	"rr":      "www",
	"type":    "A",
	"value":   acctest.CHECKSET,
	"zone_id": acctest.CHECKSET,
	"ttl":     "60",
}
