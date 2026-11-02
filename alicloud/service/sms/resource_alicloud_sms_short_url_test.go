package sms_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSMSShortUrl_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sms_short_url.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.SMSSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSMSShortUrlMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSmsShortUrl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10, 99)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSMSShortUrlBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"short_url_name": "${var.name}",
					"source_url":     "${var.source_url}",
					"effective_days": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"short_url_name": name,
						"source_url":     acctest.CHECKSET,
						"effective_days": "30",
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

var AlicloudSMSShortUrlMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudSMSShortUrlBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
variable "source_url" {
  default = "%s/%s"
}
`, name, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"), name)
}
