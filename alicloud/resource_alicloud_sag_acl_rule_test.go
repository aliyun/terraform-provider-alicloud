package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSagAclRule_basic(t *testing.T) {
	var acr smartag.Acr
	resourceId := "alicloud_sag_acl_rule.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &acr, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testSagAclName")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagAclRuleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SmartagSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id":            "${alicloud_sag_acl.default.id}",
					"description":       "tf-testSagAclRule",
					"policy":            "drop",
					"ip_protocol":       "ALL",
					"direction":         "out",
					"source_cidr":       "10.10.10.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "192.168.10.0/24",
					"dest_port_range":   "-1/-1",
					"priority":          "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":            CHECKSET,
						"description":       "tf-testSagAclRule",
						"policy":            "drop",
						"ip_protocol":       "ALL",
						"direction":         "out",
						"source_cidr":       "10.10.10.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "192.168.10.0/24",
						"dest_port_range":   "-1/-1",
						"priority":          "2",
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
					"description": "tf-testSagAclRule-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testSagAclRule-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": "accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": "accept",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"direction": "in",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"direction": "in",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "ALL",
					"source_cidr":       "10.10.11.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "192.168.11.0/24",
					"dest_port_range":   "-1/-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "ALL",
						"source_cidr":       "10.10.11.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "192.168.11.0/24",
						"dest_port_range":   "-1/-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "ICMP",
					"source_cidr":       "10.10.1.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "192.168.1.0/24",
					"dest_port_range":   "-1/-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "ICMP",
						"source_cidr":       "10.10.1.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "192.168.1.0/24",
						"dest_port_range":   "-1/-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "TCP",
					"source_cidr":       "10.10.2.0/24",
					"source_port_range": "1/65535",
					"dest_cidr":         "192.168.2.0/24",
					"dest_port_range":   "80/80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "TCP",
						"source_cidr":       "10.10.2.0/24",
						"source_port_range": "1/65535",
						"dest_cidr":         "192.168.2.0/24",
						"dest_port_range":   "80/80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "UDP",
					"source_cidr":       "10.10.3.0/24",
					"source_port_range": "20/20",
					"dest_cidr":         "192.168.3.0/24",
					"dest_port_range":   "1/20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "UDP",
						"source_cidr":       "10.10.3.0/24",
						"source_port_range": "20/20",
						"dest_cidr":         "192.168.3.0/24",
						"dest_port_range":   "1/20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "SSH(22)",
					"source_cidr":       "10.10.4.0/24",
					"source_port_range": "1/20",
					"dest_cidr":         "192.168.4.0/24",
					"dest_port_range":   "22/22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "SSH(22)",
						"source_cidr":       "10.10.4.0/24",
						"source_port_range": "1/20",
						"dest_cidr":         "192.168.4.0/24",
						"dest_port_range":   "22/22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "telnet(23)",
					"source_cidr":       "10.10.5.0/24",
					"source_port_range": "1/20",
					"dest_cidr":         "192.168.5.0/24",
					"dest_port_range":   "23/23",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "telnet(23)",
						"source_cidr":       "10.10.5.0/24",
						"source_port_range": "1/20",
						"dest_cidr":         "192.168.5.0/24",
						"dest_port_range":   "23/23",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "HTTP(80)",
					"source_cidr":       "10.10.6.0/24",
					"source_port_range": "1/20",
					"dest_cidr":         "192.168.6.0/24",
					"dest_port_range":   "80/80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "HTTP(80)",
						"source_cidr":       "10.10.6.0/24",
						"source_port_range": "1/20",
						"dest_cidr":         "192.168.6.0/24",
						"dest_port_range":   "80/80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol":       "HTTPS(443)",
					"source_cidr":       "10.10.7.0/24",
					"source_port_range": "1/20",
					"dest_cidr":         "192.168.7.0/24",
					"dest_port_range":   "443/443",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol":       "HTTPS(443)",
						"source_cidr":       "10.10.7.0/24",
						"source_port_range": "1/20",
						"dest_cidr":         "192.168.7.0/24",
						"dest_port_range":   "443/443",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id":            "${alicloud_sag_acl.default.id}",
					"description":       "tf-testSagAclRule",
					"policy":            "drop",
					"ip_protocol":       "ALL",
					"direction":         "out",
					"source_cidr":       "10.10.10.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "192.168.10.0/24",
					"dest_port_range":   "-1/-1",
					"priority":          "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":            CHECKSET,
						"description":       "tf-testSagAclRule",
						"policy":            "drop",
						"ip_protocol":       "ALL",
						"direction":         "out",
						"source_cidr":       "10.10.10.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "192.168.10.0/24",
						"dest_port_range":   "-1/-1",
						"priority":          "2",
					}),
				),
			},
		},
	})
}

func resourceSagAclRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_sag_acl" "default" {
  name = "${var.name}"
}

`, name)
}
