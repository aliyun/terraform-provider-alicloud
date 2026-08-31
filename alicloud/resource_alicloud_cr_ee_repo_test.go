package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCREERepo_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_repo.default"
	ra := resourceAttrInit(resourceId, AliCloudCREERepoMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEERepo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREERepoBasicDependence0)
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
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
					"namespace":   "${alicloud_cr_ee_namespace.default.name}",
					"name":        name,
					"repo_type":   "PUBLIC",
					"summary":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"repo_type":   "PUBLIC",
						"summary":     name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_type": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_immutability": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_immutability": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_immutability": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_immutability": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCREERepo_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_repo.default"
	ra := resourceAttrInit(resourceId, AliCloudCREERepoMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEERepo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREERepoBasicDependence0)
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
					"instance_id":      "${alicloud_cr_ee_instance.default.id}",
					"namespace":        "${alicloud_cr_ee_namespace.default.name}",
					"name":             name,
					"repo_type":        "PUBLIC",
					"summary":          name,
					"detail":           name,
					"tag_immutability": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":      CHECKSET,
						"namespace":        CHECKSET,
						"name":             name,
						"repo_type":        "PUBLIC",
						"summary":          name,
						"detail":           name,
						"tag_immutability": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCREERepo_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_repo.default.5"
	ra := resourceAttrInit(resourceId, AliCloudCREERepoMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEERepo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREERepoBasicDependence0)
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
					"count":       "6",
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
					"namespace":   "${alicloud_cr_ee_namespace.default.name}",
					"name":        name + "-${count.index}",
					"repo_type":   "PUBLIC",
					"summary":     name,
					"detail":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"namespace":   CHECKSET,
						"name":        name + fmt.Sprint(-5),
						"repo_type":   "PUBLIC",
						"summary":     name,
						"detail":      name,
					}),
				),
			},
		},
	})
}

var AliCloudCREERepoMap0 = map[string]string{
	"repo_id": CHECKSET,
}

func AliCloudCREERepoBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cr_ee_instance" "default" {
		payment_type   = "Subscription"
		period         = 1
		renewal_status = "ManualRenewal"
		instance_type  = "Basic"
		instance_name  = var.name
	}

	resource "alicloud_cr_ee_namespace" "default" {
  		instance_id        = alicloud_cr_ee_instance.default.id
  		name               = var.name
  		auto_create        = false
  		default_visibility = "PRIVATE"
	}
`, name)
}
