package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAlidnsRecord_old_basic(t *testing.T) {
	var v *alidns.DescribeDomainRecordInfoResponse

	resourceId := "alicloud_dns_record.default"
	ra := resourceAttrInit(resourceId, basicMap)

	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDnsRecordConfigDependence)

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
					"name":        "${alicloud_dns.default.name}",
					"host_record": "alimail",
					"type":        "CNAME",
					"value":       "mail.mxhichina.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":  fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand),
						"value": "mail.mxhichina.com",
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
					"host_record": "alimailchange",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"host_record": "alimailchange"}),
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
					testAccCheck(map[string]string{"priority": "3"}),
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
					"routing": "telecom",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"routing": "telecom"}),
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
					"name":        "${alicloud_dns.default.name}",
					"host_record": "alimail",
					"type":        "CNAME",
					"value":       "mail.mxhichin.com",
					"ttl":         "600",
					"priority":    "1",
					"routing":     "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})

}

func TestAccAlicloudAlidnsRecord_multi(t *testing.T) {
	var v *alidns.DescribeDomainRecordInfoResponse
	resourceId := "alicloud_dns_record.default.9"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordmulti%v.abc", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDnsRecordConfigDependence)

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
					"name":        "${alicloud_dns.default.name}",
					"host_record": "alimail${count.index}",
					"type":        "CNAME",
					"value":       "mail.mxhichina${count.index}.com",
					"count":       "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_record": "alimail9",
						"value":       "mail.mxhichina9.com",
					}),
				),
			},
		},
	})
}

func resourceDnsRecordConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns" "default" {
  name = "%s"
}
`, name)
}

var basicMap = map[string]string{
	"host_record": "alimail",
	"type":        "CNAME",
	"ttl":         "600",
	"priority":    "0",
	"value":       "mail.mxhichin.com",
	"routing":     "default",
	"status":      "ENABLE",
	"locked":      "false",
}
