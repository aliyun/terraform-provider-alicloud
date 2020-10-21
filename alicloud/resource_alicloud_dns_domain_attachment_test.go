package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAlidnsDomainAttachment_basic(t *testing.T) {
	var v alidns.DescribeInstanceDomainsResponse

	resourceId := "alicloud_alidns_domain_attachment.default"
	ra := resourceAttrInit(resourceId, alidnsDomainAttachmnetMap)

	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlidnsDomainAttachmentConfigDependence)

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
					"instance_id":  "${alicloud_alidns_instance.default.id}",
					"domain_names": []string{"${alicloud_alidns_domain.default.domain_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
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
