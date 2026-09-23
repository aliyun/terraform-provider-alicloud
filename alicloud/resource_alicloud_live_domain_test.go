// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Live Domain. >>> Resource test cases, automatically generated.
// Case 域名测试_Terraform 11928
func TestAccAliCloudLiveDomain_basic11928(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_live_domain.default"
	ra := resourceAttrInit(resourceId, AliCloudLiveDomainMap11928)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LiveServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLiveDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudLiveDomainBasicDependence11928)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_type": "liveVideo",
					"domain_name": "demo.alicloud.com",
					"region":      "cn-shanghai",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_type": "liveVideo",
						"domain_name": "demo.alicloud.com",
						"region":      "cn-shanghai",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope": "global",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope": "global",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "offline",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "offline",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "online",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "online",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"check_url"},
			},
		},
	})
}

func TestAccAliCloudLiveDomain_basic11928_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_live_domain.default"
	ra := resourceAttrInit(resourceId, AliCloudLiveDomainMap11928)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LiveServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLiveDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.com", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudLiveDomainBasicDependence11928)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_type":       "liveVideo",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"scope":             "global",
					"domain_name":       "demo.alicloud.com",
					"region":            "cn-shanghai",
					"check_url":         "http://demo.alicloud.com/test.html",
					"status":            "online",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_type":       "liveVideo",
						"resource_group_id": CHECKSET,
						"scope":             "global",
						"domain_name":       "demo.alicloud.com",
						"region":            "cn-shanghai",
						"status":            "online",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"check_url"},
			},
		},
	})
}

var AliCloudLiveDomainMap11928 = map[string]string{
	"create_time":       CHECKSET,
	"resource_group_id": CHECKSET,
	"scope":             CHECKSET,
	"status":            CHECKSET,
}

func AliCloudLiveDomainBasicDependence11928(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}


`, name)
}

// Test Live Domain. <<< Resource test cases, automatically generated.
