package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSagQosPolicy_basic(t *testing.T) {
	var qospy smartag.QosPolicy
	resourceId := "alicloud_sag_qos_policy.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &qospy, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	timeStr := time.Now().AddDate(0, 0, 1)
	timeParts := strings.Split(timeStr.String(), " ")[0]
	name := fmt.Sprintf("tf-testQosPolicy%s", timeParts)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagQosPolicyDependence)
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
					"qos_id":            "${alicloud_sag_qos.default.id}",
					"priority":          "2",
					"ip_protocol":       "ALL",
					"source_cidr":       "192.168.0.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "10.10.0.0/24",
					"dest_port_range":   "-1/-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"priority":          "2",
						"ip_protocol":       "ALL",
						"source_cidr":       "192.168.0.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "10.10.0.0/24",
						"dest_port_range":   "-1/-1",
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
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testSagQosPolicyDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testSagQosPolicyDescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"start_time": fmt.Sprintf("%sT00:00:00+0800", timeParts),
					"end_time":   fmt.Sprintf("%sT12:00:00+0800", timeParts),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"start_time": fmt.Sprintf("%sT00:00:00+0800", timeParts),
						"end_time":   fmt.Sprintf("%sT12:00:00+0800", timeParts)}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":          "2",
					"ip_protocol":       "ICMP",
					"source_cidr":       "192.168.1.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "10.10.1.0/24",
					"dest_port_range":   "-1/-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":          "2",
						"ip_protocol":       "ICMP",
						"source_cidr":       "192.168.1.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "10.10.1.0/24",
						"dest_port_range":   "-1/-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":          "3",
					"ip_protocol":       "TCP",
					"source_cidr":       "192.168.2.0/24",
					"source_port_range": "1/65535",
					"dest_cidr":         "10.10.2.0/24",
					"dest_port_range":   "80/80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":          "3",
						"ip_protocol":       "TCP",
						"source_cidr":       "192.168.2.0/24",
						"source_port_range": "1/65535",
						"dest_cidr":         "10.10.2.0/24",
						"dest_port_range":   "80/80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":          "4",
					"ip_protocol":       "UDP",
					"source_cidr":       "192.168.3.0/24",
					"source_port_range": "20/20",
					"dest_cidr":         "10.10.3.0/24",
					"dest_port_range":   "1/20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":          "4",
						"ip_protocol":       "UDP",
						"source_cidr":       "192.168.3.0/24",
						"source_port_range": "20/20",
						"dest_cidr":         "10.10.3.0/24",
						"dest_port_range":   "1/20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":          "5",
					"ip_protocol":       "SSH(22)",
					"source_cidr":       "192.168.4.0/24",
					"source_port_range": "22/22",
					"dest_cidr":         "10.10.4.0/24",
					"dest_port_range":   "22/22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":          "5",
						"ip_protocol":       "SSH(22)",
						"source_cidr":       "192.168.4.0/24",
						"source_port_range": "22/22",
						"dest_cidr":         "10.10.4.0/24",
						"dest_port_range":   "22/22",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":          "6",
					"ip_protocol":       "telnet(23)",
					"source_cidr":       "192.168.5.0/24",
					"source_port_range": "23/23",
					"dest_cidr":         "10.10.5.0/24",
					"dest_port_range":   "23/23",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":          "6",
						"ip_protocol":       "telnet(23)",
						"source_cidr":       "192.168.5.0/24",
						"source_port_range": "23/23",
						"dest_cidr":         "10.10.5.0/24",
						"dest_port_range":   "23/23",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":          "7",
					"ip_protocol":       "HTTP(80)",
					"source_cidr":       "192.168.6.0/24",
					"source_port_range": "80/80",
					"dest_cidr":         "10.10.6.0/24",
					"dest_port_range":   "80/80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":          "7",
						"ip_protocol":       "HTTP(80)",
						"source_cidr":       "192.168.6.0/24",
						"source_port_range": "80/80",
						"dest_cidr":         "10.10.6.0/24",
						"dest_port_range":   "80/80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority":          "7",
					"ip_protocol":       "HTTPS(443)",
					"source_cidr":       "192.168.7.0/24",
					"source_port_range": "443/443",
					"dest_cidr":         "10.10.7.0/24",
					"dest_port_range":   "443/443",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":          "7",
						"ip_protocol":       "HTTPS(443)",
						"source_cidr":       "192.168.7.0/24",
						"source_port_range": "443/443",
						"dest_cidr":         "10.10.7.0/24",
						"dest_port_range":   "443/443",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "tf-testSagQosPolicyName-update",
					"description":       "tf-testSagQosPolicyDescription-update",
					"priority":          "1",
					"ip_protocol":       "ALL",
					"source_cidr":       "192.168.0.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "10.10.0.0/24",
					"dest_port_range":   "-1/-1",
					"start_time":        fmt.Sprintf("%sT00:00:00+0800", timeParts),
					"end_time":          fmt.Sprintf("%sT12:00:00+0800", timeParts),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              "tf-testSagQosPolicyName-update",
						"description":       "tf-testSagQosPolicyDescription-update",
						"priority":          "1",
						"ip_protocol":       "ALL",
						"source_cidr":       "192.168.0.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "10.10.0.0/24",
						"dest_port_range":   "-1/-1",
						"start_time":        fmt.Sprintf("%sT00:00:00+0800", timeParts),
						"end_time":          fmt.Sprintf("%sT12:00:00+0800", timeParts),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSagQosPolicy_multi(t *testing.T) {
	var qospy smartag.QosPolicy
	resourceId := "alicloud_sag_qos_policy.default.6"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &qospy, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	timeStr := time.Now().AddDate(0, 0, 1)
	timeParts := strings.Split(timeStr.String(), " ")[0]
	name := fmt.Sprintf("tf-testQosPolicy%s", timeParts)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSagQosPolicyDependence)
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
					"qos_id":            "${alicloud_sag_qos.default.id}",
					"description":       "${var.name}-${count.index}",
					"count":             "7",
					"priority":          "${count.index+1}",
					"ip_protocol":       "ALL",
					"source_cidr":       "192.168.0.0/24",
					"source_port_range": "-1/-1",
					"dest_cidr":         "10.10.0.0/24",
					"dest_port_range":   "-1/-1",
					"start_time":        fmt.Sprintf("%sT00:00:00+0800", timeParts),
					"end_time":          fmt.Sprintf("%sT12:00:00+0800", timeParts),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"qos_id":            CHECKSET,
						"description":       fmt.Sprintf("%s-6", name),
						"priority":          "7",
						"ip_protocol":       "ALL",
						"source_cidr":       "192.168.0.0/24",
						"source_port_range": "-1/-1",
						"dest_cidr":         "10.10.0.0/24",
						"dest_port_range":   "-1/-1",
						"start_time":        fmt.Sprintf("%sT00:00:00+0800", timeParts),
						"end_time":          fmt.Sprintf("%sT12:00:00+0800", timeParts),
					}),
				),
			},
		},
	})
}

func resourceSagQosPolicyDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_sag_qos" "default" {
  name = "${var.name}"
}
`, name)
}
