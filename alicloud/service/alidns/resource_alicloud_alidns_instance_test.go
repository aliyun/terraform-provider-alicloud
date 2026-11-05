package alidns_test

import (
	"testing"

	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alidns"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAlidnsInstance_basic(t *testing.T) {
	var v alidnssdk.DescribeDnsProductInstanceResponse

	resourceId := "alicloud_alidns_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AlidnsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return alidns.NewAlidnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}

	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", resourceAlidnsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			//acctest.PreCheckWithTime(t, []int{1})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":   "Subscription",
					"dns_security":   "basic",
					"domain_numbers": "2",
					"period":         "1",
					"renew_period":   "1",
					"renewal_status": "ManualRenewal",
					"version_code":   "version_personal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dns_security":   "basic",
						"domain_numbers": "2",
						"period":         "1",
						"renew_period":   "0",
						"renewal_status": "ManualRenewal",
						"version_code":   "version_personal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func resourceAlidnsInstanceConfigDependence(name string) string {
	return ""
}

var AlidnsInstanceBasicMap = map[string]string{}

// lintignore: R001
