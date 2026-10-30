package dts_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

const jobMonitorRuleResourceId = "alicloud_dts_job_monitor_rule.default"

func TestAccAlicloudDTSJobMonitorRule_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsjobmonitorrule%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: jobMonitorRuleBasicDependence0(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_migration_job.default.id}"
  type       = "delay"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "state"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence0(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_migration_job.default.id}"
  type       = "delay"
  state      = "Y"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "Y"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence0(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_migration_job.default.id}"
  type       = "delay"
  state      = "Y"
  delay_rule_time = "233"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "Y"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "delay_rule_time", "233"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence0(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_migration_job.default.id}"
  type            = "delay"
  state           = "N"
  delay_rule_time = "234"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "N"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "delay_rule_time", "234"),
				),
			},
			{
				ResourceName:      jobMonitorRuleResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func jobMonitorRuleBody(extra string) string {
	return fmt.Sprintf(`
resource "alicloud_dts_job_monitor_rule" "default" {
%s}
`, extra)
}

func jobMonitorRuleBasicDependence0(name string) string {
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

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
  count            = 2
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
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

resource "alicloud_dts_migration_job" "default" {
  dts_instance_id                    = alicloud_dts_migration_instance.default.id
  dts_job_name                       = var.name
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.default.0.id
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = data.alicloud_regions.default.regions.0.id
  source_endpoint_user_name          = alicloud_rds_account.default.0.name
  source_endpoint_password           = var.password
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = alicloud_db_instance.default.1.id
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = data.alicloud_regions.default.regions.0.id
  destination_endpoint_user_name     = alicloud_rds_account.default.1.name
  destination_endpoint_password      = var.password
  db_list                            = "{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}"
  structure_initialization           = true
  data_initialization                = true
  data_synchronization               = true
  depends_on                         = [alicloud_db_account_privilege.default]
}
`, name)
}

func TestAccAlicloudDTSJobMonitorRule_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsjobmonitorrule%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: jobMonitorRuleBasicDependence1(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_synchronization_job.default.id}"
  type       = "delay"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "state"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence1(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_synchronization_job.default.id}"
  type       = "delay"
  state      = "Y"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "Y"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence1(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_synchronization_job.default.id}"
  type       = "delay"
  state      = "Y"
  delay_rule_time = "233"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "Y"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "delay_rule_time", "233"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence1(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_synchronization_job.default.id}"
  type            = "delay"
  state           = "N"
  delay_rule_time = "234"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "N"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "delay_rule_time", "234"),
				),
			},
			{
				ResourceName:      jobMonitorRuleResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func jobMonitorRuleBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

## RDS MySQL Source
resource "alicloud_db_instance" "source" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-source"
}

resource "alicloud_db_database" "source_db" {
  instance_id = alicloud_db_instance.source.id
  name        = "test_database"
}

resource "alicloud_rds_account" "source_account" {
  db_instance_id   = alicloud_db_instance.source.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}

resource "alicloud_db_account_privilege" "source_privilege" {
  instance_id  = alicloud_db_instance.source.id
  account_name = alicloud_rds_account.source_account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.source_db.*.name
}

## RDS MySQL Target
resource "alicloud_db_instance" "target" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-target"
}

resource "alicloud_rds_account" "target_account" {
  db_instance_id   = alicloud_db_instance.target.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}

## DTS Data Synchronization
resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_synchronization_job" "default" {
  data_initialization = "true"
  destination_endpoint_instance_id = "${alicloud_db_instance.target.id}"
  dts_job_name = "tf-testAccCase"
  source_endpoint_password = "${alicloud_rds_account.source_account.account_password}"
  destination_endpoint_password = "${alicloud_rds_account.target_account.account_password}"
  source_endpoint_user_name = "${alicloud_rds_account.source_account.account_name}"
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_engine_name = "MySQL"
  source_endpoint_instance_type = "RDS"
  source_endpoint_instance_id = "${alicloud_db_instance.source.id}"
  source_endpoint_engine_name = "MySQL"
  source_endpoint_region = data.alicloud_regions.default.regions.0.id
  source_endpoint_database_name = "test_database"
  dts_instance_id = "${alicloud_dts_synchronization_instance.default.id}"
  data_synchronization = "true"
  destination_endpoint_region = data.alicloud_regions.default.regions.0.id
  destination_endpoint_user_name = "${alicloud_rds_account.target_account.account_name}"
  structure_initialization = "true"
  destination_endpoint_database_name = "test_database"
  db_list = "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"
}
`, name)
}

func TestAccAlicloudDTSJobMonitorRule_basic2(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsjobmonitorrule%d", dtsTestRegion(), rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: jobMonitorRuleBasicDependence2(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_subscription_job.default.id}"
  type       = "delay"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "state"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence2(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_subscription_job.default.id}"
  type       = "delay"
  state      = "Y"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "Y"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence2(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_subscription_job.default.id}"
  type       = "delay"
  state      = "Y"
  delay_rule_time = "233"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "Y"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "delay_rule_time", "233"),
				),
			},
			{
				Config: jobMonitorRuleBasicDependence2(name) + jobMonitorRuleBody(`
  dts_job_id = "${alicloud_dts_subscription_job.default.id}"
  type            = "delay"
  state           = "N"
  delay_rule_time = "234"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(jobMonitorRuleResourceId, "dts_job_id"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "type", "delay"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "state", "N"),
					resource.TestCheckResourceAttr(jobMonitorRuleResourceId, "delay_rule_time", "234"),
				),
			},
			{
				ResourceName:      jobMonitorRuleResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func jobMonitorRuleBasicDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

## RDS MySQL Source
resource "alicloud_db_instance" "source" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-source"
}

resource "alicloud_db_database" "source_db" {
  instance_id = alicloud_db_instance.source.id
  name        = "test_database"
}

resource "alicloud_rds_account" "source_account" {
  db_instance_id   = alicloud_db_instance.source.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}

resource "alicloud_db_account_privilege" "source_privilege" {
  instance_id  = alicloud_db_instance.source.id
  account_name = alicloud_rds_account.source_account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.source_db.*.name
}

## RDS MySQL Target
resource "alicloud_db_instance" "target" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = "rds-mysql-target"
}

resource "alicloud_rds_account" "target_account" {
  db_instance_id   = alicloud_db_instance.target.id
  account_name     = "test_mysql"
  account_password = "N1cetest"
}

resource "alicloud_dts_subscription_job" "default" {
    dts_job_name                        = var.name
    payment_type                        = "PayAsYouGo"
    source_endpoint_engine_name         = "MySQL"
    source_endpoint_region              = data.alicloud_regions.default.regions.0.id
    source_endpoint_instance_type       = "RDS"
    source_endpoint_instance_id         = "${alicloud_db_instance.source.id}"
    source_endpoint_database_name       = "${alicloud_rds_account.source_account.account_password}"
    source_endpoint_user_name           = "${alicloud_rds_account.source_account.account_name}"
    source_endpoint_password            = "${alicloud_rds_account.source_account.account_password}"
    db_list = "{\"test_database\":{\"name\":\"test_database\",\"all\":true,\"state\":\"normal\"}}"
    subscription_instance_network_type  = "vpc"
    subscription_instance_vpc_id        = data.alicloud_vpcs.default.ids[0]
    subscription_instance_vswitch_id    = data.alicloud_vswitches.default.ids[0]
    status                              = "Normal"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDTSJobMonitorRule(t *testing.T) {
	resourceSchema := dts.ResourceJobMonitorRule()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"dts_job_id":      "CreateJobMonitorRuleValue",
		"type":            "CreateJobMonitorRuleValue",
		"delay_rule_time": "CreateJobMonitorRuleValue",
		"phone":           "CreateJobMonitorRuleValue",
		"state":           "CreateJobMonitorRuleValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := dtsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	ReadMockResponse := map[string]interface{}{
		// DescribeJobMonitorRule
		"DtsJobId": "CreateJobMonitorRuleValue",
		"MonitorRules": []interface{}{
			map[string]interface{}{
				"DelayRuleTime": "CreateJobMonitorRuleValue",
				"Phone":         "CreateJobMonitorRuleValue",
				"State":         "CreateJobMonitorRuleValue",
				"Type":          "CreateJobMonitorRuleValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		//CreateJobMonitorRule
		"Success": true,
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dts_job_monitor_rule", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			for k, v := range operationMockResponse {
				ReadMockResponse[k] = v
			}
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Create(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateJobMonitorRule" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceSchema.Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// CreateJobMonitorRule
	attributesDiff := map[string]interface{}{
		"delay_rule_time": "UpdateValue",
		"phone":           "UpdateValue",
		"state":           "UpdateValue",
		"type":            "UpdateValue",
	}
	// newInstanceDiff merged attributesDiff into attributes in-place, so the dCompare loop below
	// compares against the updated values. Replicate that merge here.
	for k, v := range attributesDiff {
		attributes[k] = v
	}
	diff := sdkterraform.NewInstanceDiff()
	for _, key := range []string{"delay_rule_time", "phone", "state", "type"} {
		diff.Attributes[key] = &sdkterraform.ResourceAttrDiff{Old: dInit.Get(key).(string), New: attributesDiff[key].(string)}
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeJobMonitorRule Response
		"MonitorRules": []interface{}{
			map[string]interface{}{
				"DelayRuleTime": "UpdateValue",
				"Phone":         "UpdateValue",
				"State":         "UpdateValue",
				"Type":          "UpdateValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateJobMonitorRule" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceSchema.Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeJobMonitorRule" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceSchema.Read(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceSchema.Delete(dExisted, rawClient)
	assert.Nil(t, err)
}
