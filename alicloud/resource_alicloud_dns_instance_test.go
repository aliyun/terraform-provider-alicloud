package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudAlidnsInstance_basic(t *testing.T) {
	var v alidns.DescribeDnsProductInstanceResponse

	resourceId := "alicloud_dns_instance.default"
	ra := resourceAttrInit(resourceId, AlidnsInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", resourceAlidnsInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
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
						"renew_period":   "1",
						"renewal_status": "ManualRenewal",
						"version_code":   "version_personal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "renew_period", "renewal_status"},
			},
		},
	})
}

func resourceAlidnsInstanceConfigDependence(name string) string {
	return ""
}

var AlidnsInstanceBasicMap = map[string]string{}
