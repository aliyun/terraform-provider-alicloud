package dts_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
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

const (
	unitTestCheckJobInstanceId = "dts-instance-test"
	unitTestCheckJobJobId      = "dts-job-test"
)

type dtsUnitRequestRecorder struct {
	calls    []string
	requests []map[string]interface{}
}

func (r *dtsUnitRequestRecorder) count(action string) int {
	count := 0
	for _, call := range r.calls {
		if call == action {
			count++
		}
	}
	return count
}

func (r *dtsUnitRequestRecorder) requestFor(action string) map[string]interface{} {
	for i := len(r.calls) - 1; i >= 0; i-- {
		if r.calls[i] == action {
			return r.requests[i]
		}
	}
	return nil
}

// dtsUnitTestAPIClient builds a client that never reaches the network: the DoRequest patched in
// each test intercepts every call, so the credentials below only have to be syntactically valid.
func dtsUnitTestAPIClient(t *testing.T) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("failed to build test credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoints.Store("dts", "dts.cn-hangzhou.aliyuncs.com")
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		Region:               "cn-hangzhou",
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "https",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	apiClient, err := config.Client()
	if err != nil {
		t.Fatalf("failed to build test client: %v", err)
	}
	return apiClient
}

func dtsUnitTestAPIError(code string) error {
	return &tea.SDKError{
		Code:       tea.String(code),
		Data:       tea.String(code),
		Message:    tea.String(code),
		StatusCode: tea.Int(400),
	}
}

func patchDtsUnitAPI(recorder *dtsUnitRequestRecorder, responses map[string]interface{}, failures map[string]error) *gomonkey.Patches {
	return gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		recorder.calls = append(recorder.calls, *action)
		recorder.requests = append(recorder.requests, body)
		if err, ok := failures[*action]; ok {
			return nil, err
		}
		if response, ok := responses[*action]; ok {
			responseBody, _ := response.(map[string]interface{})
			return responseBody, nil
		}
		return map[string]interface{}{"Success": true}, nil
	})
}

// newCheckJobUnitResource builds the resource through the public factory and injects the client
// through the Framework Configure hook, the same path the provider uses.
func newCheckJobUnitResource(t *testing.T, apiClient *connectivity.AliyunClient) fwresource.Resource {
	t.Helper()
	r := dts.NewCheckJobResource()
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigure", r)
	}
	var resp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: apiClient}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", resp.Diagnostics)
	}
	return r
}

func checkJobSchemaAndType(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResp.Schema.Type().TerraformType(ctx))
	}
	return schemaResp.Schema, objType
}

// checkJobPlanValue builds a plan/state value: every attribute is null unless overridden.
func checkJobPlanValue(objType tftypes.Object, overrides map[string]tftypes.Value) tftypes.Value {
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range overrides {
		values[name] = value
	}
	return tftypes.NewValue(objType, values)
}

func checkJobBaseOverrides() map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"dts_job_name":                       tftypes.NewValue(tftypes.String, "tf-test-check-job"),
		"db_list":                            tftypes.NewValue(tftypes.String, `{"test_database":{"name":"test_database","all":true}}`),
		"data_check_configure":               tftypes.NewValue(tftypes.String, `{"fullDataCheck":true,"dataCheckDbList":{"test_database":{"name":"test_database","all":true}}}`),
		"source_endpoint_instance_type":      tftypes.NewValue(tftypes.String, "RDS"),
		"source_endpoint_engine_name":        tftypes.NewValue(tftypes.String, "MySQL"),
		"destination_endpoint_instance_type": tftypes.NewValue(tftypes.String, "RDS"),
		"destination_endpoint_engine_name":   tftypes.NewValue(tftypes.String, "MySQL"),
	}
}

func checkJobCreatePlanValue(t *testing.T, r fwresource.Resource, extra map[string]tftypes.Value) (fwschema.Schema, tftypes.Value) {
	t.Helper()
	resourceSchema, objType := checkJobSchemaAndType(t, r)
	overrides := checkJobBaseOverrides()
	overrides["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	for name, value := range extra {
		overrides[name] = value
	}
	return resourceSchema, checkJobPlanValue(objType, overrides)
}

func checkJobStateString(t *testing.T, state tfsdk.State, name string) string {
	t.Helper()
	var value types.String
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

// An invalid or non-object reserve must fail while the ConfigureDtsJob request is only being
// built; nothing may be purchased and no resource id may be left behind.
func TestUnitAlicloudDTSCheckJobCreateValidatesBeforePurchase(t *testing.T) {
	cases := map[string]struct {
		reserve string
		wantErr string
	}{
		"invalidJson": {reserve: "not-json", wantErr: "invalid character"},
		"jsonNull":    {reserve: "null", wantErr: "reserve must be a JSON object"},
		"jsonArray":   {reserve: "[]", wantErr: "reserve must be a JSON object"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient := dtsUnitTestAPIClient(t)
			r := newCheckJobUnitResource(t, apiClient)
			resourceSchema, planValue := checkJobCreatePlanValue(t, r, map[string]tftypes.Value{
				"reserve":             tftypes.NewValue(tftypes.String, c.reserve),
				"source_endpoint_ssl": tftypes.NewValue(tftypes.String, "1"),
			})
			req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
			resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

			recorder := &dtsUnitRequestRecorder{}
			patches := patchDtsUnitAPI(recorder, nil, nil)
			r.Create(context.Background(), req, resp)
			patches.Reset()

			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected a local validation error for reserve %q", c.reserve)
			}
			assert.Contains(t, fmt.Sprint(resp.Diagnostics), c.wantErr)
			assert.Empty(t, recorder.calls, "no API call may be sent before the local validation passes")
			assert.Empty(t, checkJobStateString(t, resp.State, "id"))
		})
	}
}

// A purchase that cannot be configured must be released again, and the released instance must not
// be remembered as an existing resource.
func TestUnitAlicloudDTSCheckJobCreateReleasesUnconfiguredInstance(t *testing.T) {
	apiClient := dtsUnitTestAPIClient(t)
	r := newCheckJobUnitResource(t, apiClient)
	resourceSchema, planValue := checkJobCreatePlanValue(t, r, nil)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	recorder := &dtsUnitRequestRecorder{}
	patches := patchDtsUnitAPI(recorder, map[string]interface{}{
		"CreateDtsInstance": map[string]interface{}{"Success": true, "InstanceId": unitTestCheckJobInstanceId},
		"DeleteDtsJob":      map[string]interface{}{"Success": true},
	}, map[string]error{
		"ConfigureDtsJob": dtsUnitTestAPIError("InvalidParameter"),
	})
	r.Create(context.Background(), req, resp)
	patches.Reset()

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the failed configuration to surface as an error")
	}
	assert.Contains(t, fmt.Sprint(resp.Diagnostics), "InvalidParameter")
	assert.True(t, resp.State.Raw.IsNull(), "a released instance must not keep a resource id")
	assert.Equal(t, []string{"CreateDtsInstance", "ConfigureDtsJob", "DeleteDtsJob"}, recorder.calls)
	releaseRequest := recorder.requestFor("DeleteDtsJob")
	assert.Equal(t, unitTestCheckJobInstanceId, releaseRequest["DtsInstanceId"])
	assert.NotContains(t, releaseRequest, "DtsJobId")
}

// When the release itself fails the instance id stays as the resource id, refresh keeps it, and
// destroy addresses the instance instead of a job that never existed.
func TestUnitAlicloudDTSCheckJobCreateCompensationFailureStaysTrackable(t *testing.T) {
	apiClient := dtsUnitTestAPIClient(t)
	r := newCheckJobUnitResource(t, apiClient)
	resourceSchema, planValue := checkJobCreatePlanValue(t, r, nil)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	recorder := &dtsUnitRequestRecorder{}
	patches := patchDtsUnitAPI(recorder, map[string]interface{}{
		"CreateDtsInstance": map[string]interface{}{"Success": true, "InstanceId": unitTestCheckJobInstanceId},
	}, map[string]error{
		"ConfigureDtsJob": dtsUnitTestAPIError("InvalidParameter"),
		"DeleteDtsJob":    dtsUnitTestAPIError("InvalidParameter"),
	})
	r.Create(context.Background(), req, resp)
	patches.Reset()

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected the failed configuration to surface as an error")
	}
	assert.Contains(t, fmt.Sprint(resp.Diagnostics), unitTestCheckJobInstanceId)
	assert.Contains(t, fmt.Sprint(resp.Diagnostics), "releasing the purchased DTS instance")
	assert.Equal(t, unitTestCheckJobInstanceId, checkJobStateString(t, resp.State, "id"), "an unreleased instance must stay as the resource id")
	assert.Equal(t, unitTestCheckJobInstanceId, checkJobStateString(t, resp.State, "dts_instance_id"))

	// The instance has no job, so the job lookup reports not found; the retained id must survive
	// that refresh instead of being cleared.
	readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: resp.State.Raw}}
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: resp.State.Raw}}
	readRecorder := &dtsUnitRequestRecorder{}
	readPatches := patchDtsUnitAPI(readRecorder, nil, map[string]error{
		"DescribeDtsJobDetail": dtsUnitTestAPIError("InvalidJobId"),
	})
	r.Read(context.Background(), readReq, readResp)
	readPatches.Reset()
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
	}
	assert.Equal(t, unitTestCheckJobInstanceId, checkJobStateString(t, readResp.State, "id"))
	assert.Equal(t, unitTestCheckJobInstanceId, checkJobStateString(t, readResp.State, "dts_instance_id"))

	// Destroy must release the instance by its instance id alone.
	deleteReq := fwresource.DeleteRequest{State: tfsdk.State{Schema: resourceSchema, Raw: resp.State.Raw}}
	deleteResp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: resourceSchema, Raw: resp.State.Raw}}
	deleteRecorder := &dtsUnitRequestRecorder{}
	deletePatches := patchDtsUnitAPI(deleteRecorder, map[string]interface{}{
		"DeleteDtsJob": map[string]interface{}{"Success": true},
	}, nil)
	r.Delete(context.Background(), deleteReq, deleteResp)
	deletePatches.Reset()
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics: %v", deleteResp.Diagnostics)
	}
	assert.Equal(t, 1, deleteRecorder.count("DeleteDtsJob"))
	deleteRequest := deleteRecorder.requestFor("DeleteDtsJob")
	assert.Equal(t, unitTestCheckJobInstanceId, deleteRequest["DtsInstanceId"])
	assert.NotContains(t, deleteRequest, "DtsJobId")
}

// A configured job must be identified by its job id; the instance id stays an attribute and no
// compensation is issued.
func TestUnitAlicloudDTSCheckJobCreateUsesJobIdAfterConfigure(t *testing.T) {
	apiClient := dtsUnitTestAPIClient(t)
	r := newCheckJobUnitResource(t, apiClient)
	resourceSchema, planValue := checkJobCreatePlanValue(t, r, nil)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	recorder := &dtsUnitRequestRecorder{}
	patches := patchDtsUnitAPI(recorder, map[string]interface{}{
		"CreateDtsInstance": map[string]interface{}{"Success": true, "InstanceId": unitTestCheckJobInstanceId},
		"ConfigureDtsJob":   map[string]interface{}{"Success": true, "DtsJobId": unitTestCheckJobJobId},
		"DescribeDtsJobDetail": map[string]interface{}{
			"Success":             true,
			"Status":              "Finished",
			"DtsJobName":          "tf-test-check-job",
			"DtsInstanceID":       unitTestCheckJobInstanceId,
			"DtsJobClass":         "small",
			"DbObject":            `{"test_database":{"name":"test_database","all":true}}`,
			"SourceEndpoint":      map[string]interface{}{},
			"DestinationEndpoint": map[string]interface{}{},
			"MigrationMode":       map[string]interface{}{},
		},
	}, nil)
	r.Create(context.Background(), req, resp)
	patches.Reset()

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}
	assert.Equal(t, unitTestCheckJobJobId, checkJobStateString(t, resp.State, "id"))
	assert.Equal(t, unitTestCheckJobInstanceId, checkJobStateString(t, resp.State, "dts_instance_id"))
	assert.Equal(t, "small", checkJobStateString(t, resp.State, "instance_class"))
	assert.Equal(t, "Finished", checkJobStateString(t, resp.State, "status"))
	assert.Equal(t, 0, recorder.count("DeleteDtsJob"))
	configureRequest := recorder.requestFor("ConfigureDtsJob")
	assert.Equal(t, unitTestCheckJobInstanceId, configureRequest["DtsInstanceId"])
	assert.Equal(t, "CHECK", configureRequest["JobType"])
}

// Read must expose DtsJobClass, the destination account identity and fold the transient precheck
// phases of a check job into the writable domain.
func TestUnitAlicloudDTSCheckJobRead(t *testing.T) {
	cases := []struct {
		apiStatus string
		expected  string
	}{
		{"Prechecking", "Migrating"},
		{"PreCheckPass", "Migrating"},
		{"Migrating", "Migrating"},
		{"Suspending", "Suspending"},
		{"Finished", "Finished"},
		{"Failed", "Failed"},
		{"Locked", "Locked"},
	}
	apiClient := dtsUnitTestAPIClient(t)
	for _, c := range cases {
		t.Run(c.apiStatus, func(t *testing.T) {
			r := newCheckJobUnitResource(t, apiClient)
			resourceSchema, objType := checkJobSchemaAndType(t, r)
			overrides := checkJobBaseOverrides()
			overrides["id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobJobId)
			overrides["dts_instance_id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobInstanceId)
			stateValue := checkJobPlanValue(objType, overrides)
			readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
			readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			recorder := &dtsUnitRequestRecorder{}
			patches := patchDtsUnitAPI(recorder, map[string]interface{}{
				"DescribeDtsJobDetail": map[string]interface{}{
					"Success":       true,
					"Status":        c.apiStatus,
					"DtsJobClass":   "medium",
					"DtsJobName":    "tf-test-check-job",
					"DtsInstanceID": unitTestCheckJobInstanceId,
					"SourceEndpoint": map[string]interface{}{
						"AliyunUid": "1234567890",
						"RoleName":  "source-role",
					},
					"DestinationEndpoint": map[string]interface{}{
						"AliyunUid": "1234567890",
						"RoleName":  "destination-role",
					},
				},
			}, nil)
			r.Read(context.Background(), readReq, readResp)
			patches.Reset()

			if readResp.Diagnostics.HasError() {
				t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
			}
			assert.Equal(t, c.expected, checkJobStateString(t, readResp.State, "status"))
			assert.Equal(t, "medium", checkJobStateString(t, readResp.State, "instance_class"))
			assert.Equal(t, unitTestCheckJobInstanceId, checkJobStateString(t, readResp.State, "dts_instance_id"))
			// The destination account identity is ForceNew and must survive as a read value.
			assert.Equal(t, "1234567890", checkJobStateString(t, readResp.State, "destination_endpoint_owner_id"))
			assert.Equal(t, "destination-role", checkJobStateString(t, readResp.State, "destination_endpoint_role"))
			assert.Equal(t, "1234567890", checkJobStateString(t, readResp.State, "source_endpoint_owner_id"))
		})
	}
}

// DeleteDtsJob failing with the documented not-found codes must be treated as already released.
func TestUnitAlicloudDTSCheckJobDelete(t *testing.T) {
	apiClient := dtsUnitTestAPIClient(t)
	cases := map[string]struct {
		failures  map[string]error
		wantError bool
	}{
		"DeleteNormal": {failures: nil, wantError: false},
		"DeleteRepeatedInvalidJobId": {failures: map[string]error{
			"DeleteDtsJob": dtsUnitTestAPIError("InvalidJobId"),
		}, wantError: false},
		"DeleteRepeatedInstanceNotFound": {failures: map[string]error{
			"DeleteDtsJob": dtsUnitTestAPIError("Forbidden.InstanceNotFound"),
		}, wantError: false},
		"DeleteUnexpectedError": {failures: map[string]error{
			"DeleteDtsJob": dtsUnitTestAPIError("InvalidParameter"),
		}, wantError: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := newCheckJobUnitResource(t, apiClient)
			resourceSchema, objType := checkJobSchemaAndType(t, r)
			overrides := checkJobBaseOverrides()
			overrides["id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobJobId)
			overrides["dts_instance_id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobInstanceId)
			stateValue := checkJobPlanValue(objType, overrides)
			deleteReq := fwresource.DeleteRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
			deleteResp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			recorder := &dtsUnitRequestRecorder{}
			patches := patchDtsUnitAPI(recorder, nil, c.failures)
			r.Delete(context.Background(), deleteReq, deleteResp)
			patches.Reset()

			if c.wantError {
				assert.True(t, deleteResp.Diagnostics.HasError())
				return
			}
			assert.False(t, deleteResp.Diagnostics.HasError())
		})
	}
}

// instance_class is replacement-only. TransferInstanceClass answers Success for a CHECK job but
// DescribeDtsJobDetail keeps reporting the creation class, so an in-place resize can never
// converge and would leave a permanent diff once the class is read back.
func TestUnitAlicloudDTSCheckJobInstanceClassRequiresReplacement(t *testing.T) {
	r := dts.NewCheckJobResource()
	resourceSchema, _ := checkJobSchemaAndType(t, r)
	attr, ok := resourceSchema.Attributes["instance_class"].(fwschema.StringAttribute)
	if !ok {
		t.Fatalf("instance_class is %T, want schema.StringAttribute", resourceSchema.Attributes["instance_class"])
	}
	assert.True(t, attr.Optional)
	assert.True(t, attr.Computed)
	requiresReplace := false
	for _, modifier := range attr.PlanModifiers {
		if strings.Contains(modifier.Description(context.Background()), "destroy and recreate") {
			requiresReplace = true
		}
	}
	assert.True(t, requiresReplace, "instance_class cannot be changed in place")
}

// Suspending must not start the job: a running job is suspended directly, and a finished one-shot
// check job already satisfies the requested state without any call at all.
func TestUnitAlicloudDTSCheckJobUpdateSuspendDoesNotStart(t *testing.T) {
	cases := map[string]struct {
		current      string
		wantSuspends int
	}{
		"running":  {current: "Migrating", wantSuspends: 1},
		"finished": {current: "Finished", wantSuspends: 0},
	}
	apiClient := dtsUnitTestAPIClient(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := newCheckJobUnitResource(t, apiClient)
			resourceSchema, objType := checkJobSchemaAndType(t, r)

			stateOverrides := checkJobBaseOverrides()
			stateOverrides["id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobJobId)
			stateOverrides["dts_instance_id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobInstanceId)
			stateOverrides["status"] = tftypes.NewValue(tftypes.String, c.current)
			stateValue := checkJobPlanValue(objType, stateOverrides)

			planOverrides := checkJobBaseOverrides()
			planOverrides["id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobJobId)
			planOverrides["dts_instance_id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobInstanceId)
			planOverrides["status"] = tftypes.NewValue(tftypes.String, "Suspending")
			planValue := checkJobPlanValue(objType, planOverrides)

			updateReq := fwresource.UpdateRequest{
				Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
				State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
			}
			updateResp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			recorder := &dtsUnitRequestRecorder{}
			suspended := false
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				recorder.calls = append(recorder.calls, *action)
				recorder.requests = append(recorder.requests, body)
				if *action == "SuspendDtsJob" {
					suspended = true
				}
				if *action == "DescribeDtsJobDetail" {
					status := c.current
					if suspended {
						status = "Suspending"
					}
					return map[string]interface{}{"Success": true, "Status": status, "DtsJobClass": "small", "DtsJobName": "tf-test-check-job", "DtsInstanceID": unitTestCheckJobInstanceId}, nil
				}
				return map[string]interface{}{"Success": true}, nil
			})
			r.Update(context.Background(), updateReq, updateResp)
			patches.Reset()

			if updateResp.Diagnostics.HasError() {
				t.Fatalf("Update returned diagnostics: %v", updateResp.Diagnostics)
			}
			assert.Equal(t, 0, recorder.count("StartDtsJob"), "suspending must not start the job")
			assert.Equal(t, c.wantSuspends, recorder.count("SuspendDtsJob"))
		})
	}
}

// TestUnitAlicloudDTSCheckJobUpdateKeepsPlannedStatus covers the apply-time status drift of a
// one-shot check: the job reaches Finished while an update is in flight, while the plan still
// carries the previous Migrating value. Terraform rejects a known planned value that changes
// during apply, so the update keeps the planned status and the next refresh reports the server
// value.
func TestUnitAlicloudDTSCheckJobUpdateKeepsPlannedStatus(t *testing.T) {
	apiClient := dtsUnitTestAPIClient(t)
	r := newCheckJobUnitResource(t, apiClient)
	resourceSchema, objType := checkJobSchemaAndType(t, r)

	stateOverrides := checkJobBaseOverrides()
	stateOverrides["id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobJobId)
	stateOverrides["dts_instance_id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobInstanceId)
	stateOverrides["status"] = tftypes.NewValue(tftypes.String, "Migrating")
	stateValue := checkJobPlanValue(objType, stateOverrides)

	planOverrides := checkJobBaseOverrides()
	planOverrides["id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobJobId)
	planOverrides["dts_instance_id"] = tftypes.NewValue(tftypes.String, unitTestCheckJobInstanceId)
	planOverrides["status"] = tftypes.NewValue(tftypes.String, "Migrating")
	planOverrides["dts_job_name"] = tftypes.NewValue(tftypes.String, "tf-test-check-job-renamed")
	planValue := checkJobPlanValue(objType, planOverrides)

	updateReq := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	updateResp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	recorder := &dtsUnitRequestRecorder{}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		recorder.calls = append(recorder.calls, *action)
		recorder.requests = append(recorder.requests, body)
		if *action == "DescribeDtsJobDetail" {
			// The check finished after the plan was taken.
			return map[string]interface{}{
				"Success":       true,
				"Status":        "Finished",
				"DtsJobClass":   "small",
				"DtsJobName":    "tf-test-check-job-renamed",
				"DtsInstanceID": unitTestCheckJobInstanceId,
			}, nil
		}
		return map[string]interface{}{"Success": true}, nil
	})
	r.Update(context.Background(), updateReq, updateResp)
	patches.Reset()

	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", updateResp.Diagnostics)
	}
	assert.Equal(t, 1, recorder.count("ModifyDtsJobName"))
	assert.Equal(t, "Migrating", checkJobStateString(t, updateResp.State, "status"))
	assert.Equal(t, "tf-test-check-job-renamed", checkJobStateString(t, updateResp.State, "dts_job_name"))
}
