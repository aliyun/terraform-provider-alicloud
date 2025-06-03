package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudBastionhostInstance_basic(t *testing.T) {
	var v yundun_bastionhost.Instance
	resourceId := "alicloud_bastionhost_instance.default"
	ra := resourceAttrInit(resourceId, bastionhostInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBastionhostInstanceDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":        "${var.name}",
					"license_code":       "bhah_ent_50_asset",
					"plan_code":          "cloudbastion",
					"storage":            "5",
					"bandwidth":          "10",
					"period":             "1",
					"vswitch_id":         "${local.vswitch_id}",
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}", "${alicloud_security_group.default.1.id}"},
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          name,
						"period":               "1",
						"security_group_ids.#": "2",
					}),
				),
			},
			// currenly, there is a api bug when moving resource group
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"resource_group_id": CHECKSET,
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"license_code": "bhah_ent_100_asset",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"license_code": "bhah_ent_100_asset",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{"${alicloud_security_group.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance-test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance-test",
						"tags.Updated": "TF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period":        "2",
					"renewal_period_unit": "M",
					"renewal_status":      "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_period":        "2",
						"renewal_period_unit": "M",
						"renewal_status":      "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "NotRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "NotRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ad_auth_server": []map[string]interface{}{
						{
							"server":         "192.168.1.1",
							"standby_server": "192.168.1.3",
							"port":           "80",
							"domain":         "domain",
							"account":        "cn=Manager,dc=test,dc=com",
							"password":       "YouPassword123",
							"filter":         "objectClass=person",
							"name_mapping":   "nameAttr",
							"email_mapping":  "emailAttr",
							"mobile_mapping": "mobileAttr",
							"is_ssl":         "true",
							"base_dn":        "dc=test,dc=com",
						},
					},
					"ldap_auth_server": []map[string]interface{}{
						{
							"server":             "192.168.1.1",
							"standby_server":     "192.168.1.3",
							"port":               "80",
							"login_name_mapping": "uid",
							"account":            "cn=Manager,dc=test,dc=com",
							"password":           "YouPassword123",
							"filter":             "objectClass=person",
							"name_mapping":       "nameAttr",
							"email_mapping":      "emailAttr",
							"mobile_mapping":     "mobileAttr",
							"is_ssl":             "true",
							"base_dn":            "dc=test,dc=com",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ad_auth_server.#":   "1",
						"ldap_auth_server.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ad_auth_server": []map[string]interface{}{
						{
							"server":   "192.168.1.1",
							"port":     "80",
							"is_ssl":   "false",
							"domain":   "domain",
							"account":  "cn=Manager,dc=test,dc=com",
							"password": "YouPassword123",
							"base_dn":  "dc=test,dc=com",
						},
					},
					"ldap_auth_server": []map[string]interface{}{
						{
							"server":   "192.168.1.1",
							"port":     "80",
							"password": "YouPassword123",
							"account":  "cn=Manager,dc=test,dc=com",
							"base_dn":  "dc=test,dc=com",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ad_auth_server.#":   "1",
						"ldap_auth_server.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					//"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":        "${var.name}",
					"license_code":       "bhah_ent_200_asset",
					"security_group_ids": []string{"${alicloud_security_group.default.0.id}", "${alicloud_security_group.default.1.id}"},
					"tags":               REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"resource_group_id":    CHECKSET,
						"description":          name,
						"license_code":         "bhah_ent_200_asset",
						"security_group_ids.#": "2",
						"tags.%":               REMOVEKEY,
						"tags.Created":         REMOVEKEY,
						"tags.For":             REMOVEKEY,
						"tags.Updated":         REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccAliCloudBastionhostInstance_PublicAccess(t *testing.T) {
	var v yundun_bastionhost.Instance
	resourceId := "alicloud_bastionhost_instance.default"
	ra := resourceAttrInit(resourceId, bastionhostInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBastionhostInstanceDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"license_code":         "bhah_ent_50_asset",
					"period":               "1",
					"plan_code":            "cloudbastion",
					"storage":              "5",
					"bandwidth":            "10",
					"description":          "${var.name}",
					"vswitch_id":           "${local.vswitch_id}",
					"security_group_ids":   []string{"${alicloud_security_group.default.0.id}"},
					"enable_public_access": "false",
					"public_white_list":    []string{"192.168.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          name,
						"period":               "1",
						"plan_code":            "cloudbastion",
						"security_group_ids.#": "1",
						"enable_public_access": "false",
						"public_white_list.#":  "1",
						"public_white_list.0":  "192.168.0.0/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_access": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"public_white_list": []string{"192.168.0.0/18"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_white_list.#": "1",
						"public_white_list.0": "192.168.0.0/18",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func resourceBastionhostInstanceDependence(name string) string {
	return fmt.Sprintf(
		`data "alicloud_zones" "default" {
				  available_resource_creation = "VSwitch"
				}

				data "alicloud_resource_manager_resource_groups" "default"{
					status="OK"
				}
				
				variable "name" {
				  default = "%s"
				}

				data "alicloud_vpcs" "default" {
					name_regex = "^default-NODELETING$"
				}
				data "alicloud_vswitches" "default" {
					vpc_id = data.alicloud_vpcs.default.ids.0
					zone_id = data.alicloud_zones.default.zones.0.id
				}
				
				locals {
				  vswitch_id = data.alicloud_vswitches.default.ids[0]
				}
				
				resource "alicloud_security_group" "default" {
				  count  = 2
				  name   = "${var.name}"
				  vpc_id = data.alicloud_vpcs.default.ids.0
				}
`, name)
}

var bastionhostInstanceBasicMap = map[string]string{
	"description":          CHECKSET,
	"license_code":         "bhah_ent_50_asset",
	"period":               "1",
	"vswitch_id":           CHECKSET,
	"security_group_ids.#": "1",
}
