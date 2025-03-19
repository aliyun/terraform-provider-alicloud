package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAlidnsRecordWeight_basic(t *testing.T) {
	var v alidns.DescribeDomainRecordInfoResponse

	resourceId := "alicloud_alidns_record_weight.test1_record_weight_1"
	ra := resourceAttrInit(resourceId, alidnsWeightBasicMap)

	serviceFunc := func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%salidnsrecordweightbasic%v.abc", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlidnsRecordWeightConfigDependence)

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
					"domain_name": "${alicloud_dns.default.name}",
					"rr":          "test1",
					"type":        "A",
					"value":       "1.1.1.1",
					"ttl":         "900",
					"priority":    "5",
					"line":        "default",
					"status":      "ENABLE",
					"remark":      "test1 new record 1",
					"weight":      "60",
					"wrr_status":  "ENABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": fmt.Sprintf("tf-testacc%salidnsrecordweightbasic%v.abc", defaultRegionToTest, rand),
						"rr":          "test1",
						"type":        "A",
						"value":       "1.1.1.1",
						"ttl":         "900",
						"priority":    "5",
						"line":        "default",
						"status":      "ENABLE",
						"remark":      "test1 new record 1",
						"weight":      "60",
						"wrr_status":  "ENABLE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remark", "weight", "wrr_status"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rr": "test1change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"rr": "test1change"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":     "CNAME",
					"priority": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "CNAME",
						"priority": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"weight": "70"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"ttl": "800"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "test updated remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "test updated remark",
					}),
				),
			},
		},
	})

}

func resourceAlidnsRecordWeightConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "default" {
  name = "%s"
}
`, name)
}

var alidnsWeightBasicMap = map[string]string{}
