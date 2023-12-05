package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudFirewallAddressBook_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallAddressBookBasicDependence0)
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
					"group_name":   name,
					"group_type":   "ip",
					"description":  name,
					"address_list": []string{"10.21.0.0/16", "10.22.0.0/16", "10.168.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":     name,
						"group_type":     "ip",
						"description":    name,
						"address_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_list": []string{"10.21.0.0/16"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallAddressBook_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallAddressBookBasicDependence0)
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
					"group_name":   name,
					"group_type":   "ipv6",
					"description":  name,
					"address_list": []string{"::1/128", "::2/128", "::3/128"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":     name,
						"group_type":     "ipv6",
						"description":    name,
						"address_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_list": []string{"::1/128"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallAddressBook_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallAddressBookBasicDependence0)
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
					"group_name":   name,
					"group_type":   "domain",
					"description":  name,
					"address_list": []string{"alibaba.com", "aliyun.com", "alicloud.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":     name,
						"group_type":     "domain",
						"description":    name,
						"address_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_list": []string{"alibaba.com"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallAddressBook_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallAddressBookBasicDependence0)
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
					"group_name":   name,
					"group_type":   "port",
					"description":  name,
					"address_list": []string{"1/1", "22/22", "88/88"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":     name,
						"group_type":     "port",
						"description":    name,
						"address_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_list": []string{"1/1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallAddressBook_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallAddressBookBasicDependence0)
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
					"group_name":  name,
					"group_type":  "tag",
					"description": name,
					"ecs_tags": []map[string]interface{}{
						{
							"tag_key":   "created",
							"tag_value": "tfTestAcc0",
						},
						{
							"tag_key":   "for",
							"tag_value": "tfTestAcc1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":  name,
						"group_type":  "tag",
						"description": name,
						"ecs_tags.#":  "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_add_tag_ecs": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_tag_ecs": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_relation": "or",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_relation": "or",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_tags": []map[string]interface{}{
						{
							"tag_key":   "created",
							"tag_value": "tfTestAcc0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_tags.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallAddressBook_basic5(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallAddressBookBasicDependence0)
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
					"group_name":       name,
					"group_type":       "tag",
					"description":      name,
					"auto_add_tag_ecs": "1",
					"tag_relation":     "or",
					"lang":             "en",
					"ecs_tags": []map[string]interface{}{
						{
							"tag_key":   "created",
							"tag_value": "tfTestAcc0",
						},
						{
							"tag_key":   "for",
							"tag_value": "tfTestAcc1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       name,
						"group_type":       "tag",
						"description":      name,
						"auto_add_tag_ecs": "1",
						"tag_relation":     "or",
						"lang":             "en",
						"ecs_tags.#":       "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_add_tag_ecs": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_add_tag_ecs": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_relation": "and",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_relation": "and",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecs_tags": []map[string]interface{}{
						{
							"tag_key":   "created",
							"tag_value": "tfTestAcc0",
						},
						{
							"tag_key":   "for",
							"tag_value": "tfTestAcc1",
						},
						{
							"tag_key":   "by",
							"tag_value": "tfTestAcc2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_tags.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AliCloudCloudFirewallAddressBookMap0 = map[string]string{}

func AliCloudCloudFirewallAddressBookBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}
`, name)
}
