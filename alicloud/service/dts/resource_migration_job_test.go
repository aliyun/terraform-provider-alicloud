package dts_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

var dtsTestEndpoints sync.Map
var dtsTestSignVersion sync.Map

// dtsTestClientForRegion rebuilds the package-level sweeper client on the
// exported connectivity API: the flat-suite sharedClientForRegion is not
// reachable from this package.
func dtsTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &dtsTestEndpoints,
		SignVersion: &dtsTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	return conf.Client()
}

func init() {
	resource.AddTestSweepers(
		"alicloud_dts_migration_job",
		&resource.Sweeper{
			Name: "alicloud_dts_migration_job",
			F:    testSweepDTSMigrationJob,
		})
}

func testSweepDTSMigrationJob(region string) error {
	client, err := dtsTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDtsJobs"
	request := map[string]interface{}{}
	request["JobType"] = "MIGRATION"
	request["PageNumber"] = 1
	request["MaxResults"] = helper.PageSizeXLarge

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.DtsJobList", response)
		totalCount, formatErr := convert.FormatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DtsJobList", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["DtsJobName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DtsJobName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DTS Migration Job: %s", item["DtsJobName"].(string))
				continue
			}
			action := "DeleteDtsJob"
			request := map[string]interface{}{
				"DtsJobId": item["DtsJobId"],
			}
			_, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DTS Migration Job (%s): %s", item["DtsJobId"].(string), err)
			}
			log.Printf("[INFO] Delete DTS Migration Job success: %s ", item["DtsJobId"].(string))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

const migrationJobResourceId = "alicloud_dts_migration_job.default"

func TestAccAliCloudDTSMigrationJob_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationjob%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat-suite gates (checkoutSupportedRegions + testAccPreCheckWithRegions)
			// converged to cn-beijing for every input region, so pin it here.
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationJobDestroy(migrationJobResourceId),
		Steps: []resource.TestStep{
			{
				Config: migrationJobBasicDependence(name) + migrationJobBody(name, "MySQL", `
  data_check_configure = "{\"fullDataCheck\":true,\"incrementalDataCheck\":false,\"fullCheckModel\":2,\"dataCheckDbList\":{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}}"
`),
				Check: migrationJobBasicChecks(name, "MySQL",
					resource.TestCheckResourceAttr(migrationJobResourceId, "data_check_configure", `{"fullDataCheck":true,"incrementalDataCheck":false,"fullCheckModel":2,"dataCheckDbList":{"tftestdatabase":{"name":"tftestdatabase","all":true}}}`),
				),
			},
			{
				// data_check_configure is ForceNew; dropping it here would replace the job, and
				// DeleteDtsJob releases the shared dts_instance_id before the new job is configured.
				// Keep the same value so this step only exercises the status transition in place.
				Config: migrationJobBasicDependence(name) + migrationJobBody(name, "MySQL", `
  data_check_configure = "{\"fullDataCheck\":true,\"incrementalDataCheck\":false,\"fullCheckModel\":2,\"dataCheckDbList\":{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}}"
  status = "Migrating"
`),
				Check: migrationJobBasicChecks(name, "MySQL",
					resource.TestCheckResourceAttr(migrationJobResourceId, "status", "Migrating"),
				),
			},
			{
				ResourceName:            migrationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_endpoint_password", "source_endpoint_password", "data_check_configure"},
			},
		},
	})
}

func TestAccAliCloudDTSMigrationJob_mariadb(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationjobmariadb%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationJobDestroy(migrationJobResourceId),
		Steps: []resource.TestStep{
			{
				Config: migrationJobMariadbDependence(name) + migrationJobBody(name, "MARIADB", ""),
				Check:  migrationJobBasicChecks(name, ""),
			},
			{
				Config: migrationJobMariadbDependence(name) + migrationJobBody(name, "MARIADB", `
  status = "Migrating"
`),
				Check: migrationJobBasicChecks(name, "",
					resource.TestCheckResourceAttr(migrationJobResourceId, "status", "Migrating"),
				),
			},
			{
				ResourceName:            migrationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_endpoint_password", "source_endpoint_password"},
			},
		},
	})
}

func TestAccAliCloudDTSMigrationJob_ssl(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationjobssl%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The SSL steps ran in cn-hangzhou in the flat suite.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationJobDestroy(migrationJobResourceId),
		Steps: []resource.TestStep{
			{
				Config: migrationJobSslDependence(name) + migrationJobBody(name, "MySQL", `
  source_endpoint_ssl      = "1"
  destination_endpoint_ssl = "1"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(migrationJobResourceId, "source_endpoint_ssl", "1"),
					resource.TestCheckResourceAttr(migrationJobResourceId, "destination_endpoint_ssl", "1"),
				),
			},
			// Turn SSL off on both endpoints in place. The fields are not ForceNew, so this must
			// go through ModifyDtsJob with ModifyTypeEnum=UPDATE_RESERVED rather than recreating
			// the job.
			{
				Config: migrationJobSslDependence(name) + migrationJobBody(name, "MySQL", `
  source_endpoint_ssl      = "0"
  destination_endpoint_ssl = "0"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(migrationJobResourceId, "source_endpoint_ssl", "0"),
					resource.TestCheckResourceAttr(migrationJobResourceId, "destination_endpoint_ssl", "0"),
				),
			},
			// Turn SSL back on, covering the enable-by-update direction as well as disable.
			{
				Config: migrationJobSslDependence(name) + migrationJobBody(name, "MySQL", `
  source_endpoint_ssl      = "1"
  destination_endpoint_ssl = "1"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(migrationJobResourceId, "source_endpoint_ssl", "1"),
					resource.TestCheckResourceAttr(migrationJobResourceId, "destination_endpoint_ssl", "1"),
				),
			},
			{
				ResourceName:            migrationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "destination_endpoint_password", "source_endpoint_password"},
			},
		},
	})
}

func TestAccAliCloudDTSMigrationJob_vpcNat(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsmigrationjobvpcnat%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The VPC NAT steps ran in cn-hangzhou in the flat suite.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkMigrationJobDestroy(migrationJobResourceId),
		Steps: []resource.TestStep{
			// The four VPC NAT vswitch attributes only exist on ConfigureDtsJob, so this step is
			// what proves the request is accepted with them set — a rejected or misspelled
			// parameter fails the create outright. The primary and secondary vswitches are in
			// different zones, which is the arrangement the paired parameters exist for.
			{
				Config: migrationJobVpcNatDependence(name) + migrationJobVpcNatBody(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(migrationJobResourceId, "dts_job_name", name),
					// Compared against the vswitch resources rather than asserted with CHECKSET,
					// so that a primary/secondary or src/dest swap in the create request mapping
					// is caught instead of passing on any non-empty value.
					resource.TestCheckResourceAttrPair(migrationJobResourceId, "src_primary_vswitch_id", "alicloud_vswitch.primary", "id"),
					resource.TestCheckResourceAttrPair(migrationJobResourceId, "src_secondary_vswitch_id", "alicloud_vswitch.secondary", "id"),
					resource.TestCheckResourceAttrPair(migrationJobResourceId, "dest_primary_vswitch_id", "alicloud_vswitch.primary", "id"),
					resource.TestCheckResourceAttrPair(migrationJobResourceId, "dest_secondary_vswitch_id", "alicloud_vswitch.secondary", "id"),
				),
			},
			// DTS does not return the vswitch ids on DescribeDtsJobDetail, so Read cannot set them
			// and an imported job has them empty. They are ignored here for that reason, not
			// because the values are expected to drift.
			{
				ResourceName:            migrationJobResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "destination_endpoint_password", "source_endpoint_password", "src_primary_vswitch_id", "src_secondary_vswitch_id", "dest_primary_vswitch_id", "dest_secondary_vswitch_id"},
			},
		},
	})
}

func migrationJobBasicChecks(name string, engine string, extras ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(migrationJobResourceId, "dts_instance_id"),
		resource.TestCheckResourceAttr(migrationJobResourceId, "dts_job_name", name),
		resource.TestCheckResourceAttr(migrationJobResourceId, "source_endpoint_instance_type", "RDS"),
		resource.TestCheckResourceAttr(migrationJobResourceId, "destination_endpoint_instance_type", "RDS"),
		resource.TestCheckResourceAttrSet(migrationJobResourceId, "source_endpoint_region"),
		resource.TestCheckResourceAttrSet(migrationJobResourceId, "source_endpoint_user_name"),
		resource.TestCheckResourceAttrSet(migrationJobResourceId, "destination_endpoint_region"),
		resource.TestCheckResourceAttrSet(migrationJobResourceId, "destination_endpoint_user_name"),
		resource.TestCheckResourceAttrSet(migrationJobResourceId, "db_list"),
		resource.TestCheckResourceAttr(migrationJobResourceId, "structure_initialization", "true"),
		resource.TestCheckResourceAttr(migrationJobResourceId, "data_initialization", "true"),
		resource.TestCheckResourceAttr(migrationJobResourceId, "data_synchronization", "true"),
	}
	if engine == "" {
		checks = append(checks,
			resource.TestCheckResourceAttrSet(migrationJobResourceId, "source_endpoint_engine_name"),
			resource.TestCheckResourceAttrSet(migrationJobResourceId, "destination_endpoint_engine_name"),
		)
	} else {
		checks = append(checks,
			resource.TestCheckResourceAttr(migrationJobResourceId, "source_endpoint_engine_name", engine),
			resource.TestCheckResourceAttr(migrationJobResourceId, "destination_endpoint_engine_name", engine),
		)
	}
	return resource.ComposeTestCheckFunc(append(checks, extras...)...)
}

func checkMigrationJobDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_dts_migration_job" {
				continue
			}
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if _, err := dts.NewService(client).DescribeMigrationJob(rs.Primary.ID); err != nil {
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

func migrationJobBody(name, engine, extra string) string {
	return fmt.Sprintf(`
resource "alicloud_dts_migration_job" "default" {
  dts_instance_id                    = "${alicloud_dts_migration_instance.default.id}"
  dts_job_name                       = "%s"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = "${alicloud_db_instance.default.0.id}"
  source_endpoint_engine_name        = "%s"
  source_endpoint_region             = "${data.alicloud_regions.default.regions.0.id}"
  source_endpoint_user_name          = "${alicloud_rds_account.default.0.name}"
  source_endpoint_password           = "${var.password}"
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = "${alicloud_db_instance.default.1.id}"
  destination_endpoint_engine_name   = "%s"
  destination_endpoint_region        = "${data.alicloud_regions.default.regions.0.id}"
  destination_endpoint_user_name     = "${alicloud_rds_account.default.1.name}"
  destination_endpoint_password      = "${var.password}"
  db_list                            = "{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}"
  structure_initialization           = "true"
  data_initialization                = "true"
  data_synchronization               = "true"
  depends_on                         = [alicloud_db_account_privilege.default]
%s}
`, name, engine, engine, extra)
}

func migrationJobVpcNatBody(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dts_migration_job" "default" {
  dts_instance_id                    = "${alicloud_dts_migration_instance.default.id}"
  dts_job_name                       = "%s"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = "${alicloud_db_instance.default.0.id}"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = "${data.alicloud_regions.default.regions.0.id}"
  source_endpoint_user_name          = "${alicloud_rds_account.default.0.name}"
  source_endpoint_password           = "${var.password}"
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = "${alicloud_db_instance.default.1.id}"
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = "${data.alicloud_regions.default.regions.0.id}"
  destination_endpoint_user_name     = "${alicloud_rds_account.default.1.name}"
  destination_endpoint_password      = "${var.password}"
  db_list                            = "{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}"
  src_primary_vswitch_id             = "${alicloud_vswitch.primary.id}"
  src_secondary_vswitch_id           = "${alicloud_vswitch.secondary.id}"
  dest_primary_vswitch_id            = "${alicloud_vswitch.primary.id}"
  dest_secondary_vswitch_id          = "${alicloud_vswitch.secondary.id}"
  structure_initialization           = "true"
  data_initialization                = "true"
  data_synchronization               = "true"
  depends_on                         = [alicloud_db_account_privilege.default]
}
`, name)
}

func migrationJobBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "cn-beijing"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_regions" "default" {
  current = true
}

// Filter the zones by the engine actually used below, and pick a zone that the fixture VPC has a
// vswitch in. An unfiltered lookup returns an arbitrary zone, which left
// data.alicloud_vswitches.default.ids empty and failed the plan with "Invalid index".
data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones[0].id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
  count                    = 2
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  db_instance_storage_type = "cloud_essd"
  vswitch_id               = data.alicloud_vswitches.default.ids[0]
  instance_name    = join("", [var.name, count.index])
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.database_name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}
`, name)
}

func migrationJobMariadbDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_regions" "default" {
  current = true
}

// Filter the zones by the engine actually used below, and pick a zone that the fixture VPC has a
// vswitch in; fall back to creating a vswitch since MariaDB zones may have none in that VPC.
data "alicloud_db_zones" "default" {
  engine                   = "MariaDB"
  engine_version           = "10.3"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones[0].id
  engine                   = "MariaDB"
  engine_version           = "10.3"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vswitch_name = var.name
  vpc_id       = data.alicloud_vpcs.default.ids.0
  zone_id      = data.alicloud_db_zones.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  count                    = 2
  engine                   = "MariaDB"
  engine_version           = "10.3"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  db_instance_storage_type = "cloud_essd"
  vswitch_id               = local.vswitch_id
  instance_name            = join("", [var.name, count.index])
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.database_name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "mariadb"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "mariadb"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}
`, name)
}

func migrationJobSslDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

// The current region is read from the provider rather than from ALICLOUD_REGION via
// defaultRegionToTest: test Config strings are rendered before PreCheck runs, so an unset
// ALICLOUD_REGION would interpolate an empty region and DTS would reject the instance.
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
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
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

resource "alicloud_db_instance" "default" {
  count                    = 2
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  db_instance_storage_type = "cloud_essd"
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = join("", [var.name, count.index])
  ssl_action               = "Open"

  timeouts {
    create = "60m"
    update = "60m"
  }
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.database_name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}
`, name)
}

func migrationJobVpcNatDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

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
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

// The RDS instances live in the primary vswitch; the secondary is created only to be handed to
// DTS as the standby side of the VPC NAT link, so it holds no resources of its own.
resource "alicloud_vswitch" "primary" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = "${var.name}-primary"
}

resource "alicloud_vswitch" "secondary" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.1.id
  vswitch_name = "${var.name}-secondary"
}

resource "alicloud_db_instance" "default" {
  count                    = 2
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  db_instance_storage_type = "cloud_essd"
  vswitch_id               = alicloud_vswitch.primary.id
  instance_name            = join("", [var.name, count.index])
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.database_name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDTSMigrationJob(t *testing.T) {
	resourceSchema := dts.ResourceMigrationJob()
	d, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"dts_instance_id":                    "dts_instance_id",
		"dts_job_name":                       "dts_job_name",
		"source_endpoint_password":           "source_endpoint_password",
		"source_endpoint_database_name":      "source_endpoint_database_name",
		"source_endpoint_engine_name":        "MySQL",
		"source_endpoint_ip":                 "source_endpoint_ip",
		"source_endpoint_instance_id":        "source_endpoint_instance_id",
		"source_endpoint_instance_type":      "RDS",
		"source_endpoint_oracle_sid":         "source_endpoint_oracle_sid",
		"source_endpoint_port":               "source_endpoint_port",
		"source_endpoint_region":             "source_endpoint_region",
		"source_endpoint_user_name":          "source_endpoint_user_name",
		"source_endpoint_role":               "source_endpoint_role",
		"source_endpoint_owner_id":           "source_endpoint_owner_id",
		"destination_endpoint_database_name": "destination_endpoint_database_name",
		"destination_endpoint_engine_name":   "MySQL",
		"destination_endpoint_ip":            "destination_endpoint_ip",
		"destination_endpoint_instance_id":   "destination_endpoint_instance_id",
		"destination_endpoint_instance_type": "RDS",
		"destination_endpoint_oracle_sid":    "destination_endpoint_oracle_sid",
		"destination_endpoint_port":          "destination_endpoint_port",
		"destination_endpoint_region":        "destination_endpoint_region",
		"destination_endpoint_user_name":     "destination_endpoint_user_name",
		"destination_endpoint_password":      "destination_endpoint_password",
		"db_list":                            "db_list",
		"structure_initialization":           true,
		"data_initialization":                true,
		"data_synchronization":               true,
		"status":                             "Migrating",
		"checkpoint":                         "checkpoint",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := dtsTestRegion()
	rawClient, err := dtsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}
	ReadMockResponse := map[string]interface{}{
		"DtsJobId": "dts_job_id",
		"Status":   "Migrating",
		"DtsJobList": []interface{}{
			map[string]interface{}{
				"Checkpoint": "checkpoint",
				"MigrationMode": map[string]interface{}{
					"DataInitialization":       true,
					"DataSynchronization":      true,
					"structure_initialization": true,
				},
				"DbObject": "db_list",
				"DestinationEndpoint": map[string]interface{}{
					"DatabaseName": "destination_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "destination_endpoint_ip",
					"InstanceID":   "destination_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "destination_endpoint_oracle_sid",
					"Port":         "destination_endpoint_port",
					"Region":       "destination_endpoint_region",
					"UserName":     "destination_endpoint_user_name",
				},
				"DtsInstanceID": "dts_instance_id",
				"DtsJobName":    "dts_job_name",
				"PayType":       "PostPaid",
				"SourceEndpoint": map[string]interface{}{
					"DatabaseName": "source_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "source_endpoint_ip",
					"InstanceID":   "source_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "source_endpoint_oracle_sid",
					"Port":         "source_endpoint_port",
					"Region":       "source_endpoint_region",
					"UserName":     "source_endpoint_user_name",
					"RoleName":     "source_endpoint_role",
					"AliyunUid":    "source_endpoint_owner_id",
				},
				"Status":   "Migrating",
				"DtsJobId": "dts_job_id",
			},
		},
	}

	SuspendingReadMockResponse := map[string]interface{}{
		"DtsJobId": "dts_job_id",
		"Status":   "Suspending",
		"DtsJobList": []interface{}{
			map[string]interface{}{
				"Checkpoint": "checkpoint",
				"MigrationMode": map[string]interface{}{
					"DataInitialization":       true,
					"DataSynchronization":      true,
					"structure_initialization": true,
				},
				"DbObject": "db_list",
				"DestinationEndpoint": map[string]interface{}{
					"DatabaseName": "destination_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "destination_endpoint_ip",
					"InstanceID":   "destination_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "destination_endpoint_oracle_sid",
					"Port":         "destination_endpoint_port",
					"Region":       "destination_endpoint_region",
					"UserName":     "destination_endpoint_user_name",
				},
				"DtsInstanceID": "dts_instance_id",
				"DtsJobName":    "dts_job_name",
				"PayType":       "PostPaid",
				"SourceEndpoint": map[string]interface{}{
					"DatabaseName": "source_endpoint_database_name",
					"EngineName":   "MySQL",
					"Ip":           "source_endpoint_ip",
					"InstanceID":   "source_endpoint_instance_id",
					"InstanceType": "RDS",
					"OracleSID":    "source_endpoint_oracle_sid",
					"Port":         "source_endpoint_port",
					"Region":       "source_endpoint_region",
					"UserName":     "source_endpoint_user_name",
					"RoleName":     "source_endpoint_role",
					"AliyunUid":    "source_endpoint_owner_id",
				},
				"Status":   "Suspending",
				"DtsJobId": "dts_job_id",
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dts_migration_job", "dts_instance_id"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateSuspendingNormal": func(errorCode string) (map[string]interface{}, error) {
			result := SuspendingReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceSchema.Create(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceSchema.Create(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceSchema.Create(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("CreateRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceSchema.Create(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceSchema.Update(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("dts_job_id")

	t.Run("UpdateModifyDtsMigrationJobSuspendingStatusAbnormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "Migrating", New: "Suspending"}
		resourceData1, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateSuspendingNormal"]("")
		})
		err := resourceSchema.Update(resourceData1, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobSuspendingStatusNoRetryErrorAbnormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "Migrating", New: "Suspending"}
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceSchema.Update(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobSuspendingStatusNoRetryErrorNormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "Migrating", New: "Suspending"}
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		err := resourceSchema.Update(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyDtsMigrationJobMigratingStatusAbnormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "Suspending", New: "Migrating"}
		resourceData1, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceSchema.Update(resourceData1, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobMigratingStatusNoRetryErrorAbnormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "Suspending", New: "Migrating"}
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceSchema.Update(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyDTSMigrationJobMigratingStatusNoRetryErrorNormal", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "Suspending", New: "Migrating"}
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["UpdateSuspendingNormal"]("")
		})
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceSchema.Update(resourceData, rawClient)
		patcheDescribe.Reset()
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				// retry until the timeout comes
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteDescribeVpcDTSMigrationJobExpectedError", func(t *testing.T) {
		patchRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("Forbidden.InstanceNotFound")
		})
		err := resourceSchema.Delete(d, rawClient)
		patchRequest.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeVpcDTSMigrationJobNotFound", func(t *testing.T) {
		diff := sdkterraform.NewInstanceDiff()
		resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
			return responseMock["NotFoundError"]("NotFoundError")
		})
		err := resourceSchema.Read(resourceData, rawClient)
		patcheDescribe.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeDTSMigrationJobAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceSchema.Delete(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}

// lintignore: R001
func TestUnitAlicloudDTSMigrationJobResumeAcceptsFinished(t *testing.T) {
	resourceSchema := dts.ResourceMigrationJob()
	d, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	d.SetId("dts_job_id")
	assert.Nil(t, d.Set("dts_instance_id", "dts_instance_id"))
	assert.Nil(t, d.Set("dts_job_name", "dts_job_name"))

	rawClient, err := dtsTestClientForRegion(dtsTestRegion())
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}

	// A job without incremental synchronization can leave Suspending and reach
	// Finished without ever reporting Migrating. The resume waiter must accept
	// Finished rather than time out on a status the job will never show again.
	var mu sync.Mutex
	describeCalls := 0
	startCalls := 0
	describePatch := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
		mu.Lock()
		describeCalls++
		status := "Finished"
		if describeCalls == 1 {
			status = "Suspending"
		}
		mu.Unlock()
		return map[string]interface{}{
			"DtsJobId": "dts_job_id",
			"Status":   status,
			"DtsJobList": []interface{}{
				map[string]interface{}{
					"DtsJobId":      "dts_job_id",
					"DtsInstanceID": "dts_instance_id",
					"DtsJobName":    "dts_job_name",
					"Status":        status,
					"MigrationMode": map[string]interface{}{
						"DataInitialization":      false,
						"DataSynchronization":     false,
						"StructureInitialization": false,
					},
					"SourceEndpoint":      map[string]interface{}{},
					"DestinationEndpoint": map[string]interface{}{},
				},
			},
		}, nil
	})
	defer describePatch.Reset()

	requestPatch := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		if action != nil && *action == "StartDtsJob" {
			mu.Lock()
			startCalls++
			mu.Unlock()
		}
		return map[string]interface{}{"Success": true}, nil
	})
	defer requestPatch.Reset()

	diff := sdkterraform.NewInstanceDiff()
	diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "Suspending", New: "Migrating"}
	resourceData, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
	resourceData.SetId(d.Id())

	err = resourceSchema.Update(resourceData, rawClient)
	assert.Nil(t, err)
	assert.Equal(t, "Finished", resourceData.Get("status"))
	mu.Lock()
	assert.Equal(t, 1, startCalls)
	mu.Unlock()
}

// A migration that has only finished its precheck is not running yet. Read must keep
// PreCheckPass in state: folding it into Migrating makes a later explicit status = "Migrating"
// match the state, Terraform schedules no update, and StartDtsJob is never called.
func TestUnitAlicloudDTSMigrationJobReadPreservesPreCheckPass(t *testing.T) {
	resourceSchema := dts.ResourceMigrationJob()
	// A fully built client keeps every call on the patched DoRequest; a zero value would fail in
	// endpoint resolution before the fake is reached.
	meta := dtsUnitTestAPIClient(t)

	d, err := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	assert.Nil(t, err)
	d.SetId("dts-job-precheckpass")

	readPatch := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
		return map[string]interface{}{"Status": "PreCheckPass", "DtsInstanceID": "dts-instance-precheckpass"}, nil
	})
	err = resourceSchema.Read(d, meta)
	readPatch.Reset()
	assert.Nil(t, err)
	assert.Equal(t, "PreCheckPass", d.Get("status"), "an unstarted migration must not be reported as Migrating")

	// Setting status = "Migrating" afterwards has to produce an update that starts the migration.
	diff := sdkterraform.NewInstanceDiff()
	diff.Attributes["status"] = &sdkterraform.ResourceAttrDiff{Old: "PreCheckPass", New: "Migrating"}
	resourceData, err := schema.InternalMap(resourceSchema.Schema).Data(nil, diff)
	assert.Nil(t, err)
	resourceData.SetId("dts-job-precheckpass")

	var mu sync.Mutex
	startCalls := 0
	currentStatus := "PreCheckPass"
	describePatch := gomonkey.ApplyMethod(reflect.TypeOf(&dts.Service{}), "DescribeMigrationJob", func(*dts.Service, string) (map[string]interface{}, error) {
		mu.Lock()
		defer mu.Unlock()
		return map[string]interface{}{"Status": currentStatus, "DtsInstanceID": "dts-instance-precheckpass"}, nil
	})
	defer describePatch.Reset()
	requestPatch := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		if *action == "StartDtsJob" {
			mu.Lock()
			startCalls++
			currentStatus = "Migrating"
			mu.Unlock()
		}
		return map[string]interface{}{"Success": true}, nil
	})
	defer requestPatch.Reset()

	err = resourceSchema.Update(resourceData, meta)
	assert.Nil(t, err)
	mu.Lock()
	assert.Equal(t, 1, startCalls, "an explicit status change must start the migration")
	mu.Unlock()
	assert.Equal(t, "Migrating", resourceData.Get("status"))
}
