package dts_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const dtsCheckJobResourceId = "alicloud_dts_check_job.default"

func TestAccAliCloudDTSCheckJob_basic0(t *testing.T) {
	name := fmt.Sprintf("tf-testaccdtscheckjob%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The check job and its two RDS instances are created through the provider region; pin
			// it so the test does not depend on the region the harness exports.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            dtsCheckJobResourceId,
		IDRefreshIgnore:          []string{"status"},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDtsCheckJobDestroy(dtsCheckJobResourceId),
		Steps: []resource.TestStep{
			{
				// instance_class is omitted on purpose: the CHECK instance is created with the API
				// default and the value in state can only come from Read.
				Config: dtsCheckJobDependenceWithoutSSL(name, "Test12345") + dtsCheckJobBasicBody(name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "dts_job_name", name),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "source_endpoint_instance_type", "RDS"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "source_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "source_endpoint_database_name", "test_database"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "destination_endpoint_instance_type", "RDS"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "destination_endpoint_engine_name", "MySQL"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "destination_endpoint_database_name", "test_database"),
					resource.TestCheckResourceAttrSet(dtsCheckJobResourceId, "db_list"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "data_initialization", "false"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "data_synchronization", "false"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "structure_initialization", "false"),
					resource.TestCheckResourceAttrSet(dtsCheckJobResourceId, "dts_instance_id"),
					resource.TestCheckResourceAttrSet(dtsCheckJobResourceId, "instance_class"),
					testAccCheckDtsCheckJobInstanceClassFromAPI(dtsCheckJobResourceId),
					resource.TestCheckResourceAttrPair(dtsCheckJobResourceId, "source_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
					resource.TestCheckResourceAttrPair(dtsCheckJobResourceId, "destination_endpoint_region", "data.alicloud_regions.default", "regions.0.id"),
				),
			},
			{
				// A rename goes through ModifyDtsJobName and must be visible after the refresh.
				Config: dtsCheckJobDependenceWithoutSSL(name, "Test12345") + dtsCheckJobBasicBody(name+"1", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "dts_job_name", name+"1"),
					testAccCheckDtsCheckJobInstanceClassFromAPI(dtsCheckJobResourceId),
				),
			},
			{
				// Starting a one-shot check either leaves it running or completes it before the
				// apply reads the status back; both are legal, an in-flight phase is not.
				Config: dtsCheckJobDependenceWithoutSSL(name, "Test12345") + dtsCheckJobBasicBody(name+"1", `
  status = "Migrating"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "dts_job_name", name+"1"),
					testAccCheckDtsCheckJobStatusIn(dtsCheckJobResourceId, "Migrating", "Finished"),
				),
			},
			{
				ResourceName:      dtsCheckJobResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// data_check_configure and the endpoint passwords are write-only in
				// ConfigureDtsJob and never returned by DescribeDtsJobDetail, so an imported state
				// cannot carry them. The status is server-driven and may move between the last
				// apply and the import.
				ImportStateVerifyIgnore: []string{"data_check_configure", "source_endpoint_password", "destination_endpoint_password", "status"},
			},
		},
	})
}

func TestAccAliCloudDTSCheckJob_allAttributes(t *testing.T) {
	name := fmt.Sprintf("tf-testaccdtscheckjoball%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		IDRefreshName:            dtsCheckJobResourceId,
		IDRefreshIgnore:          []string{"status"},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDtsCheckJobDestroy(dtsCheckJobResourceId),
		Steps: []resource.TestStep{
			{
				Config: dtsCheckJobDependence(name, "Test12345") + dtsCheckJobFullBody(name, "small", "0", "0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "dts_job_name", name),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "instance_class", "small"),
					testAccCheckDtsCheckJobInstanceClassFromAPI(dtsCheckJobResourceId),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "source_endpoint_ssl", "0"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "destination_endpoint_ssl", "0"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "reserve", "{}"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "dts_bis_label", "normal"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "delay_notice", "true"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "delay_phone", "13800000000"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "delay_rule_time", "10"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "error_notice", "true"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "error_phone", "13800000000"),
					resource.TestCheckResourceAttrSet(dtsCheckJobResourceId, "db_list"),
					resource.TestCheckResourceAttrSet(dtsCheckJobResourceId, "data_check_configure"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "data_initialization", "false"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "data_synchronization", "false"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "structure_initialization", "false"),
					resource.TestCheckResourceAttrSet(dtsCheckJobResourceId, "dts_instance_id"),
				),
			},
			{
				// instance_class is ForceNew: the CHECK job is replaced with the requested class.
				// TransferInstanceClass answers Success but DescribeDtsJobDetail keeps reporting the
				// creation class (see evidence in the implementation report), so the new job is the
				// only way to land on a different class, and the assertion reads the class back from
				// the API rather than trusting the configuration.
				Config: dtsCheckJobDependence(name, "Test12345") + dtsCheckJobFullBody(name, "medium", "0", "0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "instance_class", "medium"),
					testAccCheckDtsCheckJobInstanceClassFromAPI(dtsCheckJobResourceId),
				),
			},
			{
				// SSL is ForceNew: ModifyDtsJob with ModifyTypeEnum=UPDATE_RESERVED answers
				// Forbidden.InstanceNotFound for CHECK instances, so the only way to land on a
				// different SSL setting is a new job. The values are read back from
				// SslSolutionEnum after the replacement.
				Config: dtsCheckJobDependence(name, "Test12345") + dtsCheckJobFullBody(name, "medium", "1", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "instance_class", "medium"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "source_endpoint_ssl", "1"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "destination_endpoint_ssl", "1"),
				),
			},
			{
				// Rotating the database account password also rotates the passwords stored on the
				// job through ModifyDtsJobPassword. The job depends on the account privileges, so
				// the databases accept the new password before DTS is updated.
				Config: dtsCheckJobDependence(name, "Test123456") + dtsCheckJobFullBody(name, "medium", "1", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "source_endpoint_password", "Test123456"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "destination_endpoint_password", "Test123456"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "source_endpoint_ssl", "1"),
					resource.TestCheckResourceAttr(dtsCheckJobResourceId, "destination_endpoint_ssl", "1"),
				),
			},
			{
				ResourceName:      dtsCheckJobResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// The create-only settings below are written to ConfigureDtsJob and are not part of
				// the DescribeDtsJobDetail response, so an imported state cannot carry them:
				// data_check_configure, both passwords, reserve, the bis label and the delay/error
				// notification settings. The status is server-driven.
				ImportStateVerifyIgnore: []string{
					"data_check_configure",
					"source_endpoint_password",
					"destination_endpoint_password",
					"reserve",
					"dts_bis_label",
					"delay_notice",
					"delay_phone",
					"delay_rule_time",
					"error_notice",
					"error_phone",
					"status",
				},
			},
		},
	})
}

// TestAccAliCloudDTSCheckJob_missingJobName keeps a focused negative case: a plan with an
// otherwise valid check job and no dts_job_name must be rejected locally. The create-only
// variants that cannot be exercised in a single-account test environment are not simulated
// here; they are recorded as explicit coverage gaps in the implementation report:
// cross-account source/destination owner id + role, Oracle SID, dedicated cluster, the VPC NAT
// vswitches and the self-managed ip/port pair, and checkpoint.
func TestAccAliCloudDTSCheckJob_missingJobName(t *testing.T) {
	name := fmt.Sprintf("tf-testaccdtscheckjobmissing%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDtsCheckJobDestroy(dtsCheckJobResourceId),
		Steps: []resource.TestStep{
			{
				Config:      dtsCheckJobDependenceWithoutSSL(name, "Test12345") + dtsCheckJobMissingNameBody(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// DeleteDtsJob releases the purchased CHECK instance, which is a separate object from the job
// record. A terminal job status ("Finished" or "MigrationFailed") only says the verification run
// ended, so release is verified through the instance id instead: DescribeMigrationInstance is the
// DtsInstanceID form of DescribeDtsJobDetail and reports NotFound once the instance is gone. The
// job lookup is kept for the diagnostic message only.
func testAccCheckDtsCheckJobDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		dtsService := dts.NewService(client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_check_job" {
				continue
			}
			instanceId := rs.Primary.Attributes["dts_instance_id"]
			if instanceId == "" {
				return fmt.Errorf("the state of %s %s has no dts_instance_id, the release of the purchased instance cannot be verified", resourceId, rs.Primary.ID)
			}
			var lastStatus string
			for i := 0; i < 90; i++ {
				if _, err := dtsService.DescribeMigrationInstance(instanceId); err != nil {
					if errs.NotFoundError(err) {
						return nil
					}
					return errs.WrapError(err)
				}
				if object, err := dtsService.DescribeCheckJob(rs.Primary.ID); err == nil {
					lastStatus = fmt.Sprint(object["Status"])
				}
				time.Sleep(10 * time.Second)
			}
			return fmt.Errorf("the CHECK instance %s of %s %s was not released, last job status: %s", instanceId, resourceId, rs.Primary.ID, lastStatus)
		}
		return nil
	}
}

// testAccCheckDtsCheckJobInstanceClassFromAPI compares the state value with DtsJobClass, the only
// API field carrying the class. It proves that a configured value was accepted and that an
// omitted value was populated from the server instead of defaulting in the provider.
func testAccCheckDtsCheckJobInstanceClassFromAPI(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := dts.NewService(client).DescribeCheckJob(rs.Primary.ID)
		if err != nil {
			return errs.WrapError(err)
		}
		apiClass, _ := object["DtsJobClass"].(string)
		if apiClass == "" {
			return fmt.Errorf("DescribeDtsJobDetail returned no DtsJobClass for %s", rs.Primary.ID)
		}
		if got := rs.Primary.Attributes["instance_class"]; got != apiClass {
			return fmt.Errorf("instance_class %q does not match DtsJobClass %q", got, apiClass)
		}
		return nil
	}
}

func testAccCheckDtsCheckJobStatusIn(resourceId string, allowed ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		status := rs.Primary.Attributes["status"]
		for _, value := range allowed {
			if status == value {
				return nil
			}
		}
		return fmt.Errorf("status %q is not one of %v", status, allowed)
	}
}

// dtsCheckJobDependence enables SSL on the RDS fixtures because DTS refuses to switch an
// endpoint to SSL when the database does not accept it, and the SSL update steps of
// TestAccAliCloudDTSCheckJob_allAttributes need that. Scenarios that never exercise SSL use
// dtsCheckJobDependenceWithoutSSL so they do not depend on the SSL modification at all.
func dtsCheckJobDependence(name, password string) string {
	return dtsCheckJobDependenceWithSSL(name, password, true)
}

func dtsCheckJobDependenceWithoutSSL(name, password string) string {
	return dtsCheckJobDependenceWithSSL(name, password, false)
}

func dtsCheckJobDependenceWithSSL(name, password string, enableSSL bool) string {
	sslAction := ""
	if enableSSL {
		sslAction = "\n  ssl_action               = \"Open\""
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "%s"
}

variable "database_name" {
  default = "test_database"
}

// The region has to come from the provider rather than from ALICLOUD_REGION: the test config is
// rendered before PreCheck runs, so an unset variable would interpolate an empty region and
// CreateDtsInstance would fail with MissingDestinationRegion.
data "alicloud_regions" "default" {
  current = true
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

// A busy account can hold an RDS instance in SSL_MODIFYING for a long time, so the fixtures wait
// longer than the 30 minute provider default. ssl_action is only rendered for the scenarios that
// exercise SSL, see dtsCheckJobDependence.
resource "alicloud_db_instance" "default" {
  count                    = 2
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  db_instance_storage_type = "cloud_essd"
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = join("", [var.name, count.index])%s

  timeouts {
    create = "60m"
    update = "60m"
  }
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = "test_mysql"
  account_password = var.password
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}
`, name, password, sslAction)
}

func dtsCheckJobBasicBody(jobName, extra string) string {
	return fmt.Sprintf(`
resource "alicloud_dts_check_job" "default" {
  dts_job_name                       = "%s"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = "${alicloud_db_instance.default.0.id}"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = "${data.alicloud_regions.default.regions.0.id}"
  source_endpoint_database_name      = "test_database"
  source_endpoint_user_name          = "${alicloud_rds_account.default.0.name}"
  source_endpoint_password           = "${var.password}"
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = "${alicloud_db_instance.default.1.id}"
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = "${data.alicloud_regions.default.regions.0.id}"
  destination_endpoint_database_name = "test_database"
  destination_endpoint_user_name     = "${alicloud_rds_account.default.1.name}"
  destination_endpoint_password      = "${var.password}"
  db_list                            = "{\"test_database\":{\"name\":\"test_database\",\"all\":true}}"
  data_check_configure               = "{\"fullDataCheck\":true,\"incrementalDataCheck\":false,\"fullCheckModel\":2,\"dataCheckDbList\":{\"test_database\":{\"name\":\"test_database\",\"all\":true}}}"
  depends_on                         = [alicloud_db_account_privilege.default]
%s}
`, jobName, extra)
}

func dtsCheckJobFullBody(jobName, instanceClass, sourceSSL, destinationSSL string) string {
	return fmt.Sprintf(`
resource "alicloud_dts_check_job" "default" {
  dts_job_name                       = "%s"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = "${alicloud_db_instance.default.0.id}"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = "${data.alicloud_regions.default.regions.0.id}"
  source_endpoint_database_name      = "test_database"
  source_endpoint_user_name          = "${alicloud_rds_account.default.0.name}"
  source_endpoint_password           = "${var.password}"
  source_endpoint_ssl                = "%s"
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = "${alicloud_db_instance.default.1.id}"
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = "${data.alicloud_regions.default.regions.0.id}"
  destination_endpoint_database_name = "test_database"
  destination_endpoint_user_name     = "${alicloud_rds_account.default.1.name}"
  destination_endpoint_password      = "${var.password}"
  destination_endpoint_ssl           = "%s"
  instance_class                     = "%s"
  reserve                            = "{}"
  dts_bis_label                      = "normal"
  delay_notice                       = true
  delay_phone                        = "13800000000"
  delay_rule_time                    = "10"
  error_notice                       = true
  error_phone                        = "13800000000"
  db_list                            = "{\"test_database\":{\"name\":\"test_database\",\"all\":true}}"
  data_check_configure               = "{\"fullDataCheck\":true,\"incrementalDataCheck\":false,\"fullCheckModel\":2,\"dataCheckDbList\":{\"test_database\":{\"name\":\"test_database\",\"all\":true}}}"
  depends_on                         = [alicloud_db_account_privilege.default]
}
`, jobName, sourceSSL, destinationSSL, instanceClass)
}

func dtsCheckJobMissingNameBody() string {
	return `
resource "alicloud_dts_check_job" "default" {
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = "${alicloud_db_instance.default.0.id}"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = "${data.alicloud_regions.default.regions.0.id}"
  source_endpoint_database_name      = "test_database"
  source_endpoint_user_name          = "${alicloud_rds_account.default.0.name}"
  source_endpoint_password           = "${var.password}"
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = "${alicloud_db_instance.default.1.id}"
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = "${data.alicloud_regions.default.regions.0.id}"
  destination_endpoint_database_name = "test_database"
  destination_endpoint_user_name     = "${alicloud_rds_account.default.1.name}"
  destination_endpoint_password      = "${var.password}"
  db_list                            = "{\"test_database\":{\"name\":\"test_database\",\"all\":true}}"
  data_check_configure               = "{\"fullDataCheck\":true,\"incrementalDataCheck\":false,\"fullCheckModel\":2,\"dataCheckDbList\":{\"test_database\":{\"name\":\"test_database\",\"all\":true}}}"
  depends_on                         = [alicloud_db_account_privilege.default]
}
`
}
