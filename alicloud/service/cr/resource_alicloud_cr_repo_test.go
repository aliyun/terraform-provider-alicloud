package cr_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	crsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cr"

	crsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCRRepo_Basic(t *testing.T) {
	var v *crsdk.GetRepoResponse
	resourceId := "alicloud_cr_repo.default"
	ra := acctest.ResourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": "${alicloud_cr_namespace.default.name}",
					"name":      "${var.name}",
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
					"detail":    acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary":   "summary",
						"repo_type": "PUBLIC",
						"detail":    acctest.REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudCRRepo_Multi(t *testing.T) {
	var v *crsdk.GetRepoResponse
	resourceId := "alicloud_cr_repo.default.4"
	ra := acctest.ResourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": "${alicloud_cr_namespace.default.name}",
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
  name               = "${var.name}"
  auto_create        = false
  default_visibility = "PRIVATE"
}
`, name)
}

var crRepoMap = map[string]string{
	"namespace": acctest.CHECKSET,
	"name":      acctest.CHECKSET,
	"summary":   "summary",
	"repo_type": "PUBLIC",
}

func TestAccAliCloudCRRepo_StateMigrationV0ToV1(t *testing.T) {
	if os.Getenv("ALICLOUD_STATE_MIGRATION_V0_V1") == "" {
		t.Skip("ALICLOUD_STATE_MIGRATION_V0_V1 not set. Test for V0 -> V1 state migration")
	}
	var v map[string]interface{}
	resourceId := "alicloud_cr_repo.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{})
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return crsvc.NewCrService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCrRepo")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-statemig-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"alicloud": {
						Source:            "aliyun/alicloud",
						VersionConstraint: "1.282.0",
					},
				},
				Config: testAccCRRepoStateMigrationConfigV0(name),
			},
			{
				ProviderFactories: acctest.SDKProviderFactories,
				Config:            testAccCRRepoStateMigrationConfigV1(name),
			},
		},
	})
}

func testAccCRRepoStateMigrationConfigV0(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cr_namespace" "default" {
  name               = "%[1]s"
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_repo" "default" {
  namespace = alicloud_cr_namespace.default.name
  name      = "%[1]s"
  summary   = "state migration test"
  repo_type = "PUBLIC"
}

output "domain_list" {
  value = alicloud_cr_repo.default.domain_list["vpc"]
}
`, name)
}

func testAccCRRepoStateMigrationConfigV1(name string) string {
	return fmt.Sprintf(`
resource "alicloud_cr_namespace" "default" {
  name               = "%[1]s"
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_repo" "default" {
  namespace = alicloud_cr_namespace.default.name
  name      = "%[1]s"
  summary   = "state migration test"
  repo_type = "PUBLIC"
}

output "domain_list" {
  value = alicloud_cr_repo.default.domain_list.0.vpc
}
`, name)
}
