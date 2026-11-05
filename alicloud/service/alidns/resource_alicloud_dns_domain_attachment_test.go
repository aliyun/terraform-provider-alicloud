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

// The dependent resource alidns_instance is prepaid so skip current test case.
func SkipTestAccAlicloudAlidnsDomainAttachment_basic(t *testing.T) {
	var v alidnssdk.DescribeInstanceDomainsResponse

	resourceId := "alicloud_alidns_domain_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, alidnsDomainAttachmnetMap)

	serviceFunc := func() interface{} {
		return alidns.NewDnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}

	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlidnsDomainAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_alidns_instance.default.id}",
					"domain_names": []string{"${alicloud_alidns_domain.default.domain_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    acctest.CHECKSET,
						"domain_names.#": "1",
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
					"domain_names": []string{"${alicloud_alidns_domain.default.domain_name}", "${alicloud_alidns_domain.default1.domain_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_names.#": "2",
					}),
				),
			},
		},
	})
}

func resourceAlidnsDomainAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_alidns_instance" "default" {
 	  dns_security    = "basic"
 	  domain_numbers  = 3
 	  version_code    = "version_personal"
 	  period          = 1
	  renewal_status  = "ManualRenewal"
	}

	resource "alicloud_alidns_domain" "default" {
  	  domain_name = "%s.abc"
	}

	resource "alicloud_alidns_domain" "default1" {
  	  domain_name = "%s1.abc"
	}
`, name, name)
}

var alidnsDomainAttachmnetMap = map[string]string{}
