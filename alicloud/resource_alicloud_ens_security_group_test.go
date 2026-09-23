package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens SecurityGroup. >>> Resource test cases, automatically generated.
// Case 5107
func TestAccAliCloudEnsSecurityGroup_basic5107(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_security_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsSecurityGroupMap5107)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senssecuritygroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsSecurityGroupBasicDependence5107)
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
					"security_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "SecurityGroupDescription_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SecurityGroupDescription_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "SecurityGroupDescription_UPDATE_autotest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SecurityGroupDescription_UPDATE_autotest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":         "SecurityGroupDescription_autotest",
					"security_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "SecurityGroupDescription_autotest",
						"security_group_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsSecurityGroupMap5107 = map[string]string{}

func AlicloudEnsSecurityGroupBasicDependence5107(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 5107  twin
func TestAccAliCloudEnsSecurityGroup_basic5107_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_security_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsSecurityGroupMap5107)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senssecuritygroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsSecurityGroupBasicDependence5107)
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
					"description":         "SecurityGroupDescription_UPDATE_autotest",
					"security_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "SecurityGroupDescription_UPDATE_autotest",
						"security_group_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 5108
func TestAccAliCloudEnsSecurityGroup_permissions5108(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_security_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsSecurityGroupMap5108)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%senssecuritygroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsSecurityGroupBasicDependence5107)
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
					"security_group_name": name,
					"permissions": []map[string]interface{}{
						{
							"direction":         "ingress",
							"ip_protocol":       "TCP",
							"port_range":        "80/80",
							"policy":            "Accept",
							"priority":          1,
							"source_cidr_ip":    "10.0.0.0/8",
							"dest_cidr_ip":      "0.0.0.0/0",
							"source_port_range": "80/80",
							"description":       "test ingress v4 rule",
						},
						{
							"direction":           "ingress",
							"ip_protocol":         "TCP",
							"port_range":          "80/80",
							"policy":              "Accept",
							"priority":            1,
							"ipv6_source_cidr_ip": "2001:db8::/32",
							"ipv6_dest_cidr_ip":   "::/0",
							"source_port_range":   "80/80",
							"description":         "test ingress v6 rule",
						},
						{
							"direction":      "egress",
							"ip_protocol":    "TCP",
							"port_range":     "443/443",
							"policy":         "Drop",
							"priority":       2,
							"source_cidr_ip": "0.0.0.0/0",
							"dest_cidr_ip":   "0.0.0.0/0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_name": name,
						"permissions.#":       "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_name": name,
					"permissions": []map[string]interface{}{
						{
							"direction":         "egress",
							"ip_protocol":       "UDP",
							"port_range":        "53/53",
							"policy":            "Accept",
							"priority":          1,
							"source_cidr_ip":    "0.0.0.0/0",
							"dest_cidr_ip":      "192.168.0.0/16",
							"source_port_range": "53/53",
							"description":       "updated egress v4 rule",
						},
						{
							"direction":           "egress",
							"ip_protocol":         "UDP",
							"port_range":          "53/53",
							"policy":              "Accept",
							"priority":            1,
							"ipv6_source_cidr_ip": "::/0",
							"ipv6_dest_cidr_ip":   "2001:db8:abcd::/48",
							"source_port_range":   "53/53",
							"description":         "updated egress v6 rule",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_name": name,
					"permissions":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"permissions.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEnsSecurityGroupMap5108 = map[string]string{}

// Test Ens SecurityGroup. <<< Resource test cases, automatically generated.
