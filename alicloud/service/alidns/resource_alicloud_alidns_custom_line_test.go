package alidns_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alidns"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAlidnsCustomLine_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_custom_line.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlidnsCustomLineMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alidns.NewAlidnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlidnsCustomLine")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlidnsCustomLineBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_line_name": name,
					"domain_name":      "${var.domain_name}",
					"ip_segment_list": []map[string]interface{}{
						{
							"start_ip": "192.0.2.123",
							"end_ip":   "192.0.2.125",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name":       acctest.CHECKSET,
						"custom_line_name":  name,
						"ip_segment_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_line_name": name + "new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_line_name": name + "new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_segment_list": []map[string]interface{}{
						{
							"start_ip": "192.0.2.137",
							"end_ip":   "192.0.2.139",
						},
						{
							"start_ip": "192.0.2.156",
							"end_ip":   "192.0.2.159",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_segment_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_segment_list": []map[string]interface{}{
						{
							"start_ip": "192.0.2.123",
							"end_ip":   "192.0.2.125",
						},
					},
					"custom_line_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_segment_list.#": "1",
						"custom_line_name":  name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AlicloudAlidnsCustomLineMap0 = map[string]string{}

func AlicloudAlidnsCustomLineBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "domain_name" {
  default = "%s"
}
`, name, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"))
}

// lintignore: R001
