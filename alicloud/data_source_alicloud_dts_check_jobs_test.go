package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDTSCheckJobsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdtscheckjobs%d", rand)
	dependence := AliCloudDTSSynchronizationJobBasicDependence0(name)
	checkJobConfig := dependence + fmt.Sprintf(`
resource "alicloud_dts_check_job" "default" {
  dts_instance_id                    = alicloud_dts_synchronization_instance.default.id
  dts_job_name                       = "%[1]s"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.source.id
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = data.alicloud_regions.default.regions.0.id
  source_endpoint_database_name      = "test_database"
  source_endpoint_user_name          = alicloud_rds_account.source_account.account_name
  source_endpoint_password           = alicloud_rds_account.source_account.account_password
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = alicloud_db_instance.target.id
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = data.alicloud_regions.default.regions.0.id
  destination_endpoint_database_name = "test_database"
  destination_endpoint_user_name     = alicloud_rds_account.target_account.account_name
  destination_endpoint_password      = alicloud_rds_account.target_account.account_password
  db_list = "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"
  data_check_configure = "{\"fullDataCheck\":true,\"incrementalDataCheck\":false}"
}

data "alicloud_dts_check_jobs" "default" {
  ids          = [alicloud_dts_check_job.default.id]
  enable_details = true
}
`, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckAlicloudDtsCheckJobDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: checkJobConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dts_check_jobs.default", "jobs.#", "1"),
					resource.TestCheckResourceAttrPair("data.alicloud_dts_check_jobs.default", "jobs.0.id", "alicloud_dts_check_job.default", "id"),
					resource.TestCheckResourceAttrPair("data.alicloud_dts_check_jobs.default", "jobs.0.dts_job_name", "alicloud_dts_check_job.default", "dts_job_name"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDtsCheckJobDestroy(s *resource.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dts_check_job" {
			continue
		}
		_, err := dtsService.DescribeDtsCheckJob(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("DTS check job %s still exists", rs.Primary.ID)
	}
	return nil
}
