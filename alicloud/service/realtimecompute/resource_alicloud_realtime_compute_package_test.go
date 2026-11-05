package realtimecompute_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The foas 2018-11-11 package lifecycle APIs operate on a FOAS project
// that has no Terraform resource: it must be pre-created (with a bound cluster
// and storage) and referenced by name. The test is skipped unless that project
// is supplied through ALICLOUD_REALTIME_COMPUTE_PROJECT. The package artifact
// is a real OSS object uploaded in the dependence block.
func testAccRealtimeComputePackageProjectName(t *testing.T) string {
	v := os.Getenv("ALICLOUD_REALTIME_COMPUTE_PROJECT")
	if v == "" {
		t.Skip("Environment variable ALICLOUD_REALTIME_COMPUTE_PROJECT is not set; a pre-created FOAS project is required for the package acceptance test")
	}
	return v
}

func TestAccAliCloudRealtimeComputePackage_basic0(t *testing.T) {
	resourceId := "alicloud_realtime_compute_package.default"
	projectName := testAccRealtimeComputePackageProjectName(t)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-rtcpkg%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRealtimeComputePackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudRealtimeComputePackageBasic0(name, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "project_name", projectName),
					resource.TestCheckResourceAttr(resourceId, "package_name", name),
					resource.TestCheckResourceAttr(resourceId, "type", "JAR"),
					resource.TestCheckResourceAttr(resourceId, "description", "tf acc test package"),
					resource.TestCheckResourceAttrSet(resourceId, "oss_bucket"),
					resource.TestCheckResourceAttrSet(resourceId, "oss_path"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttrSet(resourceId, "creator"),
				),
			},
			{
				Config: testAccAliCloudRealtimeComputePackageBasic0Update(name, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "package_name", name),
					resource.TestCheckResourceAttr(resourceId, "description", "tf acc test package updated"),
					resource.TestCheckResourceAttr(resourceId, "tag", "v2"),
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

func testAccCheckRealtimeComputePackageDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_realtime_compute_package" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := realtimecompute.NewService(client).DescribeRealtimeComputePackage(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing realtime compute package %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("realtime compute package %s still exists", rs.Primary.ID)
	}
	return nil
}

// testAccAliCloudRealtimeComputePackageBasicDependence0 stands up the OSS
// bucket and a small package artifact object; the FOAS project is referenced
// from the pre-created name passed in as projectName.
func testAccAliCloudRealtimeComputePackageBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "current" {}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "packages/test.jar"
  content = "dummy-jar-content-for-foas-package-acc-test"
}
`, name)
}

func testAccAliCloudRealtimeComputePackageBasic0(name, projectName string) string {
	return fmt.Sprintf(`%s

variable "project_name" {
  default = "%s"
}

resource "alicloud_realtime_compute_package" "default" {
  project_name = var.project_name
  package_name = var.name
  oss_bucket   = alicloud_oss_bucket.default.id
  oss_endpoint = "oss-cn-beijing.aliyuncs.com"
  oss_owner    = data.alicloud_account.current.id
  oss_path     = alicloud_oss_bucket_object.default.key
  type         = "JAR"
  description  = "tf acc test package"
}
`, testAccAliCloudRealtimeComputePackageBasicDependence0(name), projectName)
}

func testAccAliCloudRealtimeComputePackageBasic0Update(name, projectName string) string {
	return fmt.Sprintf(`%s

variable "project_name" {
  default = "%s"
}

resource "alicloud_realtime_compute_package" "default" {
  project_name = var.project_name
  package_name = var.name
  oss_bucket   = alicloud_oss_bucket.default.id
  oss_endpoint = "oss-cn-beijing.aliyuncs.com"
  oss_owner    = data.alicloud_account.current.id
  oss_path     = alicloud_oss_bucket_object.default.key
  type         = "JAR"
  description  = "tf acc test package updated"
  tag          = "v2"
}
`, testAccAliCloudRealtimeComputePackageBasicDependence0(name), projectName)
}
