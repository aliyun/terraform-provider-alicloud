package alidns_test

import (
	"fmt"
	"testing"

	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alidns"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAlidnsRecord_basic(t *testing.T) {
	var v alidnssdk.DescribeDomainRecordInfoResponse

	resourceId := "alicloud_alidns_record.default"
	ra := acctest.ResourceAttrInit(resourceId, alidnsBasicMap)

	serviceFunc := func() interface{} {
		return alidns.NewAlidnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandInt()
	name := fmt.Sprintf("tf-testacc%salidnsrecordbasic%v.abc", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlidnsRecordConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${alicloud_dns.default.name}",
					"rr":          "alimail",
					"type":        "CNAME",
					"ttl":         "600",
					"priority":    "0",
					"value":       "mail.mxhichin.com",
					"line":        "default",
					"status":      "ENABLE",
					"remark":      "test new domain record",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": fmt.Sprintf("tf-testacc%salidnsrecordbasic%v.abc", acctest.DefaultRegionToTest, rand),
						"rr":          "alimail",
						"type":        "CNAME",
						"ttl":         "600",
						"priority":    "0",
						"value":       "mail.mxhichin.com",
						"line":        "default",
						"status":      "ENABLE",
						"remark":      "test new domain record",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"remark"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rr": "alimailchange",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"rr": "alimailchange"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":     "MX",
					"priority": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "MX",
						"priority": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"value": "mail.change.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"value": "mail.change.com"}),
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
					"line": "telecom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"line": "telecom"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"ttl": "600"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "test new domain record again",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "test new domain record again",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${alicloud_dns.default.name}",
					"rr":          "alimail",
					"type":        "CNAME",
					"value":       "mail.mxhichin.com",
					"ttl":         "600",
					"priority":    "0",
					"line":        "default",
					"status":      "ENABLE",
					"remark":      "test new domain record",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": fmt.Sprintf("tf-testacc%salidnsrecordbasic%v.abc", acctest.DefaultRegionToTest, rand),
						"rr":          "alimail",
						"type":        "CNAME",
						"value":       "mail.mxhichin.com",
						"ttl":         "600",
						"priority":    "0",
						"line":        "default",
						"status":      "ENABLE",
						"remark":      "test new domain record",
					}),
				),
			},
		},
	})

}

func resourceAlidnsRecordConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "default" {
  name = "%s"
}
`, name)
}

var alidnsBasicMap = map[string]string{}
