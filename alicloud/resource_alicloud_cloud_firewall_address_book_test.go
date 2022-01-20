package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudFirewallAddressBook_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallAddressBookBasicDependence0)
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
					"group_name":       "tf-testAcc-KO98t4",
					"description":      "tf-testAcc-jOgZg",
					"group_type":       "ip",
					"address_list":     []string{"10.21.0.0/16", "10.168.0.0/16"},
					"auto_add_tag_ecs": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       "tf-testAcc-KO98t4",
						"description":      "tf-testAcc-jOgZg",
						"group_type":       "ip",
						"address_list.#":   "2",
						"auto_add_tag_ecs": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "source_ip"},
			},
		},
	})
}

func TestAccAlicloudCloudFirewallAddressBook_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallAddressBookBasicDependence0)
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
					"description":  "tf-testAcc-jOgZg",
					"group_type":   "tag",
					"tag_relation": "and",
					"ecs_tags": []map[string]interface{}{{
						"tag_key":   "created",
						"tag_value": "tfTestAcc0",
					}, {
						"tag_key":   "for",
						"tag_value": "Tftestacc1",
					}},
					"auto_add_tag_ecs": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       name,
						"description":      "tf-testAcc-jOgZg",
						"tag_relation":     "and",
						"group_type":       "tag",
						"ecs_tags.#":       "2",
						"address_list.#":   "0",
						"auto_add_tag_ecs": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "source_ip"},
			},
		},
	})
}

var AlicloudCloudFirewallAddressBookMap0 = map[string]string{
	"lang":           NOSET,
	"address_list.#": CHECKSET,
	"source_ip":      NOSET,
}

func AlicloudCloudFirewallAddressBookBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
func TestAccAlicloudCloudFirewallAddressBook_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallAddressBookMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallAddressBookBasicDependence1)
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
					"description":      "tf-testAcc-jOgZg",
					"group_type":       "ip",
					"address_list":     []string{"10.21.0.0/16", "10.168.0.0/16"},
					"auto_add_tag_ecs": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       name,
						"description":      "tf-testAcc-jOgZg",
						"group_type":       "ip",
						"address_list.#":   "2",
						"auto_add_tag_ecs": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": "tf-testAcc-Apese",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": "tf-testAcc-Apese",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAcc-Sy73h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAcc-Sy73h",
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
				Config: testAccConfig(map[string]interface{}{
					"group_name":       "tf-testAcc-zMGh5",
					"description":      "tf-testAcc-SR8cM",
					"group_type":       "ip",
					"address_list":     []string{"10.21.0.0/16", "10.168.0.0/16"},
					"auto_add_tag_ecs": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       "tf-testAcc-zMGh5",
						"description":      "tf-testAcc-SR8cM",
						"address_list.#":   "2",
						"auto_add_tag_ecs": "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAlicloudCloudFirewallAddressBook_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallAddressBookMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallAddressBookBasicDependence1)
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
					"description":  "tf-testAcc-jOgZg",
					"group_type":   "tag",
					"tag_relation": "and",
					"ecs_tags": []map[string]interface{}{{
						"tag_key":   "created",
						"tag_value": "tfTestAcc0",
					}, {
						"tag_key":   "for",
						"tag_value": "Tftestacc1",
					}},
					"auto_add_tag_ecs": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       name,
						"description":      "tf-testAcc-jOgZg",
						"tag_relation":     "and",
						"group_type":       "tag",
						"ecs_tags.#":       "2",
						"address_list.#":   "0",
						"auto_add_tag_ecs": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": "tf-testAcc-Apese",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": "tf-testAcc-Apese",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAcc-Sy73h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAcc-Sy73h",
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
					"ecs_tags": []map[string]interface{}{{
						"tag_key":   "created",
						"tag_value": "tfTestAcc2",
					}, {
						"tag_key":   "for",
						"tag_value": "Tftestacc3",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecs_tags.#": "2",
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
					"group_name":   "tf-testAcc-zMGh5",
					"description":  "tf-testAcc-SR8cM",
					"tag_relation": "and",
					"ecs_tags": []map[string]interface{}{{
						"tag_key":   "created",
						"tag_value": "tfTestAcc0",
					}, {
						"tag_key":   "for",
						"tag_value": "Tftestacc1",
					}},
					"auto_add_tag_ecs": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       "tf-testAcc-zMGh5",
						"description":      "tf-testAcc-SR8cM",
						"tag_relation":     "and",
						"ecs_tags.#":       "2",
						"address_list.#":   "0",
						"auto_add_tag_ecs": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_ip", "lang"},
			},
		},
	})
}

func TestAccAlicloudCloudFirewallAddressBook_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_address_book.default"
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallAddressBookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallAddressBook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewalladdressbook%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallAddressBookBasicDependence0)
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
					"description":      "tf-testAcc-jOgZg",
					"group_type":       "port",
					"address_list":     []string{"22", "80", "443"},
					"auto_add_tag_ecs": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       name,
						"description":      "tf-testAcc-jOgZg",
						"group_type":       "port",
						"address_list.#":   "3",
						"auto_add_tag_ecs": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":  "tf-testAcc-Apese",
					"description": "tf-testAcc-Apese-desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":  "tf-testAcc-Apese",
						"description": "tf-testAcc-Apese-desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_list": []string{"22", "80"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_list.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "source_ip"},
			},
		},
	})
}

var AlicloudCloudFirewallAddressBookMap1 = map[string]string{
	"lang":           NOSET,
	"address_list.#": CHECKSET,
	"source_ip":      NOSET,
}

func AlicloudCloudFirewallAddressBookBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
