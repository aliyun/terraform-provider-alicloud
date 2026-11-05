package dts_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const migrationInstanceResourceId = "alicloud_dts_migration_instance.default"

func TestAccAliCloudDTSMigrationInstance_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationinstance%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat-suite gates (checkoutSupportedRegions + testAccPreCheckWithRegions)
			// converged to cn-beijing for every input region, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationInstanceDestroy(migrationInstanceResourceId),
		Steps: []resource.TestStep{
			{
				Config: migrationInstanceBasicConfig(name, ""),
				Check:  migrationInstanceBasicChecks(),
			},
			{
				Config: migrationInstanceBasicConfig(name, `
  tags = {
    Created = "TF"
    For     = "Test"
    From    = "acceptance"
  }
`),
				Check: migrationInstanceBasicChecks(
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "tags.%", "3"),
				),
			},
			{
				Config: migrationInstanceBasicConfig(name, `
  tags = {
    Created = "TF"
    For     = "acceptance-Test"
  }
`),
				Check: migrationInstanceBasicChecks(
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "tags.%", "2"),
				),
			},
			{
				ResourceName:            migrationInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "sync_architecture", "compute_unit"},
			},
		},
	})
}

func TestAccAliCloudDTSMigrationInstance_mariadb(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationinstancemariadb%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationInstanceDestroy(migrationInstanceResourceId),
		Steps: []resource.TestStep{
			{
				Config: migrationInstanceDependence(name) + `
resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "mariadb"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "mariadb"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
  compute_unit                     = 2
  database_count                   = 1
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "status"),
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "payment_type", "PayAsYouGo"),
					resource.TestCheckResourceAttr(migrationInstanceResourceId, "instance_class", "small"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "source_endpoint_engine_name"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "source_endpoint_region"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "destination_endpoint_engine_name"),
					resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "destination_endpoint_region"),
				),
			},
			{
				ResourceName:            migrationInstanceResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_count", "sync_architecture", "compute_unit"},
			},
		},
	})
}

func migrationInstanceBasicChecks(extras ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "status"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "payment_type", "PayAsYouGo"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "instance_class", "small"),
		resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "destination_endpoint_region"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "source_endpoint_engine_name", "MySQL"),
		resource.TestCheckResourceAttrSet(migrationInstanceResourceId, "source_endpoint_region"),
		resource.TestCheckResourceAttr(migrationInstanceResourceId, "destination_endpoint_engine_name", "MySQL"),
	}
	return resource.ComposeTestCheckFunc(append(checks, extras...)...)
}

func checkMigrationInstanceDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_migration_instance" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeMigrationInstance(rs.Primary.ID); err != nil {
				if errs.NotFoundError(err) {
					continue
				}
				return err
			}
			return fmt.Errorf("the resource %s %s was not destroyed", resourceId, rs.Primary.ID)
		}
		return nil
	}
}

func migrationInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}
`, name)
}

func migrationInstanceBasicConfig(name string, extra string) string {
	return migrationInstanceDependence(name) + fmt.Sprintf(`
resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
%s}
`, extra)
}

func dtsTestRegion() string {
	if region := os.Getenv("ALICLOUD_REGION"); region != "" {
		return region
	}
	return "cn-beijing"
}
