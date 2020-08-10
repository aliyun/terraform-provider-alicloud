package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func SkipTestAccAlicloudDnsDomainAttachment_basic(t *testing.T) {
	var v alidns.DescribeInstanceDomainsResponse

	resourceId := "alicloud_dns_domain_attachment.default"
	ra := resourceAttrInit(resourceId, dnsDomainAttachmnetMap)

	serviceFunc := func() interface{} {
		return &DnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDnsDomainAttachmentConfigDependence)

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
					"instance_id":  "${alicloud_dns_instance.default.id}",
					"domain_names": []string{"${alicloud_dns.default.name}"},
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
					"domain_names": []string{"${alicloud_dns.default.name}", "${alicloud_dns.default1.name}"},
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

func resourceDnsDomainAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_dns_instance" "default" {
 	  dns_security    = "basic"
 	  domain_numbers  = 3
 	  version_code    = "version_personal"
 	  period          = 1
	  renewal_status  = "ManualRenewal"
	}

	resource "alicloud_dns" "default" {
  	  name = "%s.abc"
	}

	resource "alicloud_dns" "default1" {
  	  name = "%s1.abc"
	}
`, name, name)
}

var dnsDomainAttachmnetMap = map[string]string{}
