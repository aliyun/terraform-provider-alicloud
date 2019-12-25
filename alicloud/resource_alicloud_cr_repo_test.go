package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCRRepo_Basic(t *testing.T) {
	var v *cr.GetRepoResponse
	resourceId := "alicloud_cr_repo.default"
	ra := resourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": alicloud_cr_namespace.default.name,
					"name":      var.name,
					"summary":   "summary",
					"repo_type": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace": name,
						"name":      name,
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
					"detail": "detail",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": "detail",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary": "summary update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary": "summary update",
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
					"detail": "detail update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": "detail update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary":   "summary",
					"repo_type": "PUBLIC",
					"detail":    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary":   "summary",
						"repo_type": "PUBLIC",
						"detail":    REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCRRepo_Multi(t *testing.T) {
	var v *cr.GetRepoResponse
	resourceId := "alicloud_cr_repo.default.4"
	ra := resourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": alicloud_cr_namespace.default.name,
					"name":      "${var.name}${count.index}",
					"summary":   "summary",
					"repo_type": "PUBLIC",
					"count":     "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceCRRepoConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cr_namespace" "default" {
	name = var.name
	auto_create	= false
	default_visibility = "PRIVATE"
}
`, name)
}

var crRepoMap = map[string]string{
	"namespace": CHECKSET,
	"name":      CHECKSET,
	"summary":   "summary",
	"repo_type": "PUBLIC",
}
