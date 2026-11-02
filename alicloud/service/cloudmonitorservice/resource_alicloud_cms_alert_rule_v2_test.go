// Package cloudmonitorservice. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	testingterraform "github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var checkCmsAlertRuleV2Destroy = func() resource.TestCheckFunc {
	var object map[string]interface{}
	rc := acctest.ResourceCheckInitWithDescribeMethod("alicloud_cms_alert_rule_v2.default", &object, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	return rc.CheckResourceDestroy()
}()

func requireCmsAlertRuleV2ProtoV5Factory(t *testing.T) {
	if factory := acctest.ProtoV5ProviderFactories()["alicloud"]; factory == nil {
		t.Fatal("SDK acceptance tests require the real alicloud protocol 5 factory")
	}
}

// checkCmsAlertRuleV2EscalationsByMetric matches composite escalation members by metric_name so the
// assertion does not depend on the order the API echoes for the TypeList.
func checkCmsAlertRuleV2EscalationsByMetric(prefix string, want map[string]map[string]string) resource.TestCheckFunc {
	return func(s *testingterraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_cms_alert_rule_v2.default"]
		if !ok {
			return fmt.Errorf("resource alicloud_cms_alert_rule_v2.default not found in state")
		}
		count, err := strconv.Atoi(rs.Primary.Attributes[prefix+".#"])
		if err != nil {
			return fmt.Errorf("cannot read %s.#: %v", prefix, err)
		}
		seen := make(map[string]bool, len(want))
		for i := 0; i < count; i++ {
			attrPrefix := fmt.Sprintf("%s.%d", prefix, i)
			metricName := rs.Primary.Attributes[attrPrefix+".metric_name"]
			expected, ok := want[metricName]
			if !ok {
				continue
			}
			for key, value := range expected {
				if actual := rs.Primary.Attributes[attrPrefix+"."+key]; actual != value {
					return fmt.Errorf("%s[%s].%s expected %q, got %q", prefix, metricName, key, value, actual)
				}
			}
			seen[metricName] = true
		}
		for metricName := range want {
			if !seen[metricName] {
				return fmt.Errorf("escalation for metric %q not found in %s", metricName, prefix)
			}
		}
		return nil
	}
}

// Test Cms AlertRuleV2. >>> Resource test cases, automatically generated.
// Case resource_AlertRuleV2_UMODEL_NoData_test 13068
func TestAccAliCloudCmsAlertRuleV2_basic13068(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13068_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-umodel-nodata"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13068_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-umodel-nodata-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic13068_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "umodel no-data alert on $${metric}"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        enabled = false
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "regression-umodel-nodata"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13068_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "umodel no-data alert on $${metric}"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        enabled = false
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "regression-umodel-nodata-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_APM_Composite_test 12919
func TestAccAliCloudCmsAlertRuleV2_basic12919(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12919_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.errorrate"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-svc-composite-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-svc-composite-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-composite-rewrite"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "SUM"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "5"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12919_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.errorrate"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-svc-composite-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-svc-composite-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-composite-rewrite-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "SUM"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "10"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12919_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.errorrate"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        service_id_list = ["arms-svc-composite-1", "arms-svc-composite-2"]
    }

        display_name = "regression-apm-composite-rewrite"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "SUM"
        threshold_list {
            severity = "WARNING"
            threshold = 5
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12919_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.errorrate"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        service_id_list = ["arms-svc-composite-1", "arms-svc-composite-2"]
    }

        display_name = "regression-apm-composite-rewrite-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "SUM"
        threshold_list {
            severity = "CRITICAL"
            threshold = 10
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_NotifyStrategies_test 13054
func TestAccAliCloudCmsAlertRuleV2_basic13054(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13054_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13054_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-id-003"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13054_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-4"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13054_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-5"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13054_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic13054_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "regression-notify-strategies"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-1"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13054_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m1"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-id-003"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13054_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m2"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-4"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13054_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m3"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-5"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13054_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m4"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = []
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_ExpressEscalation_test 13053
func TestAccAliCloudCmsAlertRuleV2_basic13053(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13053_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-prom-express-escalation"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13053_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-prom-express-escalation-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic13053_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.2"
    }

        display_name = "regression-prom-express-escalation"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13053_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.1"
    }

        display_name = "regression-prom-express-escalation-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_Basic_test 13052
func TestAccAliCloudCmsAlertRuleV2_basic13052(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13052_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-basic-prometheus"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13052_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-basic-prometheus-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test-updated"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic13052_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "regression-basic-prometheus"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13052_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-basic-prometheus-updated"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test-updated"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Labels_Verify_test 13051
func TestAccAliCloudCmsAlertRuleV2_basic13051(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13051_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.%", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.summary", "CPU usage high"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.description", "Node CPU is above 90%"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "cloudspec_labels_verify_test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.%", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.severity", "critical"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.env", "test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic13051_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.%", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.summary", "CPU usage updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.description", "Updated CPU alert"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.runbook", "https://example.com"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "cloudspec_labels_verify_test_updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.%", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.severity", "warning"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.team", "dev"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.env", "staging"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic13051_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        annotations = {
        "summary" = "CPU usage high"
        "description" = "Node CPU is above 90%%"
    }

        display_name = "cloudspec_labels_verify_test"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        labels = {
        "severity" = "critical"
        "env" = "test"
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic13051_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        annotations = {
        "summary" = "CPU usage updated"
        "description" = "Updated CPU alert"
        "runbook" = "https://example.com"
    }

        display_name = "cloudspec_labels_verify_test_updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        labels = {
        "severity" = "warning"
        "team" = "dev"
        "env" = "staging"
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_IntegrationConfig_test 12939
func TestAccAliCloudCmsAlertRuleV2_basic12939(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12939_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-id-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.1", "action-id-002"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.2", "action-id-003"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "integration-config-roundtrip-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12939_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-id-004"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "integration-config-roundtrip-test-s1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12939_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-id-005"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.1", "action-id-006"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.2", "action-id-007"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "integration-config-roundtrip-test-s2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12939_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-id-008"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.1", "action-id-009"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.2", "action-id-010"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "integration-config-roundtrip-test-s3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12939_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "0"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "integration-config-roundtrip-test-s4-empty"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12939_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        action_integration_config {
        actions = ["action-id-001", "action-id-002", "action-id-003"]
        enabled = true
    }

        arms_integration_config {
        enabled = true
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "integration-config-roundtrip-test"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12939_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        action_integration_config {
        actions = ["action-id-004"]
        enabled = false
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "integration-config-roundtrip-test-s1"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12939_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        action_integration_config {
        actions = ["action-id-005", "action-id-006", "action-id-007"]
        enabled = true
    }

        arms_integration_config {
        enabled = true
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "integration-config-roundtrip-test-s2"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12939_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        action_integration_config {
        actions = ["action-id-008", "action-id-009", "action-id-010"]
        enabled = true
    }

        arms_integration_config {
        enabled = true
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "integration-config-roundtrip-test-s3"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12939_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        action_integration_config {
        actions = []
        enabled = true
    }

        arms_integration_config {
        enabled = true
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "integration-config-roundtrip-test-s4-empty"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_APM_Aggregate_test 12937
func TestAccAliCloudCmsAlertRuleV2_basic12937(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12937_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-agg-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-agg-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-aggregate"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "SUM"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "1000"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.threshold", "500"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12937_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-agg-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-agg-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-aggregate-max"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "LT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "MAX"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "200"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12937_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-agg-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-agg-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-aggregate-p50"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "P50"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "INFO"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "50"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12937_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-agg-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-agg-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-aggregate-p75"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "P75"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "75"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12937_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-agg-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-agg-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-aggregate-p90"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "LT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "P90"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "90"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12937_step5(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-agg-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-agg-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-aggregate-p99"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "P99"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "99"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12937_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc", "rpcType"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-agg-1", "arms-agg-2"]
    }

        display_name = "regression-apm-aggregate"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "SUM"
        threshold_list {
            severity = "CRITICAL"
            threshold = 1000
        }
        threshold_list {
            severity = "WARNING"
            threshold = 500
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12937_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc", "rpcType"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-agg-1", "arms-agg-2"]
    }

        display_name = "regression-apm-aggregate-max"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "LT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "MAX"
        threshold_list {
            severity = "WARNING"
            threshold = 200
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12937_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-agg-1", "arms-agg-2"]
    }

        display_name = "regression-apm-aggregate-p50"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "P50"
        threshold_list {
            severity = "INFO"
            threshold = 50
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12937_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpcType"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-agg-1", "arms-agg-2"]
    }

        display_name = "regression-apm-aggregate-p75"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "P75"
        threshold_list {
            severity = "ERROR"
            threshold = 75
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12937_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc", "rpcType"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-agg-1", "arms-agg-2"]
    }

        display_name = "regression-apm-aggregate-p90"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "LT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "P90"
        threshold_list {
            severity = "WARNING"
            threshold = 90
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12937_step5(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc", "rpcType"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-agg-1", "arms-agg-2"]
    }

        display_name = "regression-apm-aggregate-p99"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "P99"
        threshold_list {
            severity = "CRITICAL"
            threshold = 99
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_QueryConfig_Advanced_test 12936
func TestAccAliCloudCmsAlertRuleV2_basic12936(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12936_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.0.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.0.field", "regionId"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.0.value", "cn-hangzhou"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.1.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.1.field", "instanceId"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.1.value", "i-bp1abc123def"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.0.field", "name"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.0.value", "my-instance"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-queryconfig-advanced"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12936_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.0.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.0.field", "regionId"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.0.value", "cn-shanghai"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.1.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.1.field", "instanceId"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.1.value", "i-bp2ghi456jkl"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.2.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.2.field", "status"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_filters.2.value", "Running"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "memory_usedutilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.0.field", "name"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.0.value", "my-ecs-instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.1.field", "hostName"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_fields.1.value", "web-server-01"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-queryconfig-advanced-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test-updated"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "95"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12936_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        entity_filters {
            operator = "="
            field = "regionId"
            value = "cn-hangzhou"
        }
        entity_filters {
            operator = "="
            field = "instanceId"
            value = "i-bp1abc123def"
        }
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
        entity_fields {
            field = "name"
            value = "my-instance"
        }
    }

        display_name = "regression-queryconfig-advanced"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12936_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "UMODEL"
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        entity_filters {
            operator = "="
            field = "regionId"
            value = "cn-shanghai"
        }
        entity_filters {
            operator = "="
            field = "instanceId"
            value = "i-bp2ghi456jkl"
        }
        entity_filters {
            operator = "="
            field = "status"
            value = "Running"
        }
        metric = "memory_usedutilization"
        metric_set = "acs_ecs_dashboard"
        entity_fields {
            field = "name"
            value = "my-ecs-instance"
        }
        entity_fields {
            field = "hostName"
            value = "web-server-01"
        }
    }

        display_name = "regression-queryconfig-advanced-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test-updated"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "WARNING"
        duration_secs = 120
        threshold = 95
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_APM_CompareList_test 12935
func TestAccAliCloudCmsAlertRuleV2_basic12935(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12935_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-svc-compare-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-svc-compare-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-comparelist"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.aggregate", "AVG"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.yoy_time_value", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.threshold", "80"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.yoy_time_unit", "hour"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_COMPOSITE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.relation", "OR"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12935_step1(name, "AND", "regression-apm-comparelist-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-svc-compare-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-svc-compare-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-comparelist-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.aggregate", "P99"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.yoy_time_value", "24"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.threshold", "95"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.yoy_time_unit", "hour"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_COMPOSITE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.relation", "AND"),
				),
			},
			{
				// condition_config.relation is updatable on its own: this step keeps the ForceNew
				// compare_list values unchanged and only flips the relation, so it reaches Update.
				Config: testAccAliCloudCmsAlertRuleV2_basic12935_step1(name, "OR", "regression-apm-comparelist-relation-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-comparelist-relation-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.relation", "OR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.aggregate", "P99"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.compare_list.0.threshold", "95"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12935_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        service_id_list = ["arms-svc-compare-1", "arms-svc-compare-2"]
    }

        display_name = "regression-apm-comparelist"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        compare_list {
            operator = "GT"
            aggregate = "AVG"
            yoy_time_value = 1
            threshold = 80
            yoy_time_unit = "hour"
        }
        type = "APM_COMPOSITE_CONDITION"
        relation = "OR"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12935_step1(name, relation, displayName string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        service_id_list = ["arms-svc-compare-1", "arms-svc-compare-2"]
    }

        display_name = "%[3]s"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        compare_list {
            operator = "GT"
            aggregate = "P99"
            yoy_time_value = 24
            threshold = 95
            yoy_time_unit = "hour"
        }
        type = "APM_COMPOSITE_CONDITION"
        relation = "%[2]s"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name, relation, displayName)
}

// Case resource_AlertRuleV2_APM_Operator_test 12934
func TestAccAliCloudCmsAlertRuleV2_basic12934(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12934_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-app-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-app-002"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-operator"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "COUNT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "1000"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.threshold", "500"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12934_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-app-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-app-002"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-operator-gte"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "MIN"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "800"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12934_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-app-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-app-002"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-operator-lt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "LT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "CONTINUES"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "INFO"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "10"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12934_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-app-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-app-002"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-operator-lte"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "LT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "COUNT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "20"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12934_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-app-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-app-002"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-operator-eq"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "MIN"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "0"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12934_step5(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-app-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-app-002"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-operator-ne"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "NE"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "CONTINUES"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "999"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12934_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.count"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-app-001", "arms-app-002"]
    }

        display_name = "regression-apm-operator"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "COUNT"
        threshold_list {
            severity = "CRITICAL"
            threshold = 1000
        }
        threshold_list {
            severity = "WARNING"
            threshold = 500
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12934_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.count"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-app-001", "arms-app-002"]
    }

        display_name = "regression-apm-operator-gte"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "MIN"
        threshold_list {
            severity = "WARNING"
            threshold = 800
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12934_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.count"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-app-001", "arms-app-002"]
    }

        display_name = "regression-apm-operator-lt"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "LT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "CONTINUES"
        threshold_list {
            severity = "INFO"
            threshold = 10
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12934_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.count"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-app-001", "arms-app-002"]
    }

        display_name = "regression-apm-operator-lte"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "LT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "COUNT"
        threshold_list {
            severity = "WARNING"
            threshold = 20
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12934_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.count"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-app-001", "arms-app-002"]
    }

        display_name = "regression-apm-operator-eq"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "EQ"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "MIN"
        threshold_list {
            severity = "ERROR"
            threshold = 0
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12934_step5(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.count"
            group_by = ["rpc"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        service_id_list = ["arms-app-001", "arms-app-002"]
    }

        display_name = "regression-apm-operator-ne"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "NE"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "CONTINUES"
        threshold_list {
            severity = "CRITICAL"
            threshold = 999
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_SimpleEscalation_test 12932
func TestAccAliCloudCmsAlertRuleV2_basic12932(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12932_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "memory_usedutilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-umodel-simple-escalation"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GE"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "20"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12932_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "memory_usedutilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-umodel-simple-escalation-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "LE"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "180"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12932_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "memory_usedutilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "regression-umodel-simple-escalation"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GE"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "WARNING"
        duration_secs = 120
        threshold = 20
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12932_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "memory_usedutilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "regression-umodel-simple-escalation-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "LE"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 180
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_APM_Simple_test 12930
func TestAccAliCloudCmsAlertRuleV2_basic12930(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12930_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.value", "http"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "NE"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.value", "/health"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "svc-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "svc-002"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-simple-8"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "AVG"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "1000"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.threshold", "500"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12930_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.value", "/orders"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "svc-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "svc-003"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.2", "svc-004"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-simple-8-modified"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "AVG"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "50"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.severity", "INFO"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.threshold", "30"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.threshold", "80"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12930_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.value", "/orders"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "svc-009"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-simple-8-modified-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "AVG"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "50"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.severity", "INFO"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.threshold", "30"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.threshold", "80"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12930_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpcType"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.value", "/orders"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "svc-010"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "svc-011"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-simple-8-modified-3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "AVG"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "50"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.severity", "INFO"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.threshold", "30"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.threshold", "80"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12930_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "0"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "EQ"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.value", "/orders"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "svc-012"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-simple-8-modified-4-empty"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "AVG"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "50"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.severity", "INFO"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.1.threshold", "30"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.2.threshold", "80"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12930_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc", "rpcType"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "EQ"
            value = "http"
            key = "rpcType"
        }
        filter_list {
            type = "NE"
            value = "/health"
            key = "rpc"
        }
        service_id_list = ["svc-001", "svc-002"]
    }

        display_name = "regression-apm-simple-8"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "AVG"
        threshold_list {
            severity = "CRITICAL"
            threshold = 1000
        }
        threshold_list {
            severity = "WARNING"
            threshold = 500
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12930_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc", "rpcType"]
            window_secs = 120
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "EQ"
            value = "/orders"
            key = "rpc"
        }
        service_id_list = ["svc-001", "svc-003", "svc-004"]
    }

        display_name = "regression-apm-simple-8-modified"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "AVG"
        threshold_list {
            severity = "WARNING"
            threshold = 50
        }
        threshold_list {
            severity = "INFO"
            threshold = 30
        }
        threshold_list {
            severity = "ERROR"
            threshold = 80
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12930_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpc"]
            window_secs = 120
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "EQ"
            value = "/orders"
            key = "rpc"
        }
        service_id_list = ["svc-009"]
    }

        display_name = "regression-apm-simple-8-modified-2"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "AVG"
        threshold_list {
            severity = "WARNING"
            threshold = 50
        }
        threshold_list {
            severity = "INFO"
            threshold = 30
        }
        threshold_list {
            severity = "ERROR"
            threshold = 80
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12930_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = ["rpcType", "rpc"]
            window_secs = 120
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "EQ"
            value = "/orders"
            key = "rpc"
        }
        service_id_list = ["svc-010", "svc-011"]
    }

        display_name = "regression-apm-simple-8-modified-3"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "AVG"
        threshold_list {
            severity = "WARNING"
            threshold = 50
        }
        threshold_list {
            severity = "INFO"
            threshold = 30
        }
        threshold_list {
            severity = "ERROR"
            threshold = 80
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12930_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by = []
            window_secs = 120
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "EQ"
            value = "/orders"
            key = "rpc"
        }
        service_id_list = ["svc-012"]
    }

        display_name = "regression-apm-simple-8-modified-4-empty"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "AVG"
        threshold_list {
            severity = "WARNING"
            threshold = 50
        }
        threshold_list {
            severity = "INFO"
            threshold = 30
        }
        threshold_list {
            severity = "ERROR"
            threshold = 80
        }
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_UMODEL_test 12929
func TestAccAliCloudCmsAlertRuleV2_basic12929(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12929_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.0.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.0.value", "web-server"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.0.name", "app"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.1.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.1.value", "production"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.1.name", "env"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-umodel-10"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12929_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-integration-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.0.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.0.value", "api-gateway"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.0.name", "app"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.1.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.1.value", "staging"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.1.name", "env"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.2.operator", "="),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.2.value", "cn-hangzhou"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.label_filters.2.name", "region"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-umodel-10-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "95"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12929_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "umodel test alert on $${metric}"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        enabled = false
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        label_filters {
            operator = "="
            value = "web-server"
            name = "app"
        }
        label_filters {
            operator = "="
            value = "production"
            name = "env"
        }
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "regression-umodel-10"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12929_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "umodel test alert updated on $${metric}"

        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        actions = ["action-integration-001"]
        enabled = true
    }

        arms_integration_config {
        enabled = true
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        label_filters {
            operator = "="
            value = "api-gateway"
            name = "app"
        }
        label_filters {
            operator = "="
            value = "staging"
            name = "env"
        }
        label_filters {
            operator = "="
            value = "cn-hangzhou"
            name = "region"
        }
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "regression-umodel-10-updated"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 120
        threshold = 95
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_Rich_test 12928
func TestAccAliCloudCmsAlertRuleV2_basic12928(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12928_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "rate(http_requests_total[5m]) > 100"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-enriched-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.5", "6"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.6", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "23:59"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "86400"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "00:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "group-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "group-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.2", "group-3"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12928_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "rate(http_requests_total[5m]) > 200"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-enriched-1-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "18:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "43200"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "09:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test-updated"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12928_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "rate(http_requests_total[5m]) > 200"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-enriched-1-m2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "18:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "43200"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "09:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "group-4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "group-5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.2", "group-6"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12928_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "rate(http_requests_total[5m]) > 200"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-enriched-1-m3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "18:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "43200"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "09:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "group-4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "group-5"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12928_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "rate(http_requests_total[5m]) > 200"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-enriched-1-m4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "18:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "43200"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "09:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "group-7"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12928_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "告警：$${alertName} 触发"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "rate(http_requests_total[5m]) > 100"
    }

        display_name = "regression-enriched-1"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5, 6, 7]
        active_end_time = "23:59"
        silence_time_secs = 86400
        active_start_time = "00:00"
        channels {
            type = "GROUP"
            identifiers = ["group-1", "group-2", "group-3"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12928_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "更新告警：$${alertName}"

        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "rate(http_requests_total[5m]) > 200"
    }

        display_name = "regression-enriched-1-updated"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5]
        active_end_time = "18:00"
        silence_time_secs = 43200
        active_start_time = "09:00"
        channels {
            type = "GROUP"
            identifiers = ["regression-test-updated"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12928_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "更新告警：$${alertName}"

        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "rate(http_requests_total[5m]) > 200"
    }

        display_name = "regression-enriched-1-m2"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5]
        active_end_time = "18:00"
        silence_time_secs = 43200
        active_start_time = "09:00"
        channels {
            type = "GROUP"
            identifiers = ["group-4", "group-5", "group-6"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12928_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "更新告警：$${alertName}"

        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "rate(http_requests_total[5m]) > 200"
    }

        display_name = "regression-enriched-1-m3"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5]
        active_end_time = "18:00"
        silence_time_secs = 43200
        active_start_time = "09:00"
        channels {
            type = "GROUP"
            identifiers = ["group-4", "group-5"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12928_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "更新告警：$${alertName}"

        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "rate(http_requests_total[5m]) > 200"
    }

        display_name = "regression-enriched-1-m4"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5]
        active_end_time = "18:00"
        silence_time_secs = 43200
        active_start_time = "09:00"
        channels {
            type = "GROUP"
            identifiers = ["group-7"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_Array_test 12924
func TestAccAliCloudCmsAlertRuleV2_basic12924(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12924_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-severity-array"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.5", "6"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.6", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "23:59"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "00:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "86400"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "group-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "group-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.2", "group-3"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "INFO"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12924_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-severity-array-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "6"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "6"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.5", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "20:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "43200"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "group-a"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "group-b"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12924_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "regression-severity-array"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5, 6, 7]
        active_end_time = "23:59"
        active_start_time = "00:00"
        silence_time_secs = 86400
        channels {
            type = "GROUP"
            identifiers = ["group-1", "group-2", "group-3"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "INFO"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12924_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-severity-array-updated"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 5, 6, 7]
        active_end_time = "20:00"
        active_start_time = "08:00"
        silence_time_secs = 43200
        channels {
            type = "GROUP"
            identifiers = ["group-a", "group-b"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "ERROR"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_NotifyChannels_test 12923
func TestAccAliCloudCmsAlertRuleV2_basic12923(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12923_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "cmsxx-prometheus-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up == 0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-channels"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "23:59"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "86400"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "00:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "cms-alarm-group-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "cms-alarm-group-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.type", "DINGTALK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.0", "dingtalk-webhook-url-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.1", "dingtalk-webhook-url-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.type", "WEBHOOK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.0", "https://webhook.example.com/alert"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.1", "https://webhook2.example.com/alert"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12923_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "cmsxx-prometheus-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up == 0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-channels-m1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.5", "6"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.6", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "23:59"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "86400"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "00:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "CONTACT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "contact-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "contact-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.2", "contact-3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.type", "FEISHU"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.0", "feishu-webhook-url-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.1", "feishu-webhook-url-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.type", "SLACK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.0", "slack-channel-url-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.1", "slack-channel-url-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.3.type", "WEIXIN"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.3.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.3.identifiers.0", "weixin-bot-key-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.3.identifiers.1", "weixin-bot-key-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12923_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "cmsxx-prometheus-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up == 0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-channels-m2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "6"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "20:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "43200"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "group-a"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "group-b"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.2", "group-c"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.type", "DINGTALK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.0", "dt-url-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.1", "dt-url-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.2", "dt-url-3"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12923_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "cmsxx-prometheus-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up == 0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-channels-m3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "18:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "86400"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "09:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "CONTACT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "admin-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "admin-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.type", "WEBHOOK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.0", "https://hook1.example.com"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.1", "https://hook2.example.com"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12923_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "cmsxx-prometheus-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up == 0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-channels-m4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.4", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.5", "6"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.6", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "23:59"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "86400"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "00:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "ops-group-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "ops-group-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.type", "DINGTALK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.0", "dt-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.1", "dt-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.type", "FEISHU"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.0", "feishu-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.1", "feishu-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12923_step5(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "90"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "cmsxx-prometheus-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up == 0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-channels-m5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.utc_offset", "+08:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.#", "4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.0", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.1", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.2", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_days.3", "7"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_end_time", "22:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.silence_time_secs", "86400"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.active_start_time", "06:00"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "CONTACT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "dev-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.1", "dev-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.2", "dev-3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.type", "SLACK"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.0", "slack-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.1.identifiers.1", "slack-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.type", "WEIXIN"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.0", "wx-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.1", "wx-2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.2.identifiers.2", "wx-3"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "90"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12923_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "cmsxx-prometheus-1"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up == 0"
    }

        display_name = "regression-notify-channels"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 3, 5]
        active_end_time = "23:59"
        silence_time_secs = 86400
        active_start_time = "00:00"
        channels {
            type = "GROUP"
            identifiers = ["cms-alarm-group-1", "cms-alarm-group-2"]
        }
        channels {
            type = "DINGTALK"
            identifiers = ["dingtalk-webhook-url-1", "dingtalk-webhook-url-2"]
        }
        channels {
            type = "WEBHOOK"
            identifiers = ["https://webhook.example.com/alert", "https://webhook2.example.com/alert"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12923_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "cmsxx-prometheus-1"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up == 0"
    }

        display_name = "regression-notify-channels-m1"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5, 6, 7]
        active_end_time = "23:59"
        silence_time_secs = 86400
        active_start_time = "00:00"
        channels {
            type = "CONTACT"
            identifiers = ["contact-1", "contact-2", "contact-3"]
        }
        channels {
            type = "FEISHU"
            identifiers = ["feishu-webhook-url-1", "feishu-webhook-url-2"]
        }
        channels {
            type = "SLACK"
            identifiers = ["slack-channel-url-1", "slack-channel-url-2"]
        }
        channels {
            type = "WEIXIN"
            identifiers = ["weixin-bot-key-1", "weixin-bot-key-2"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12923_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "cmsxx-prometheus-1"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up == 0"
    }

        display_name = "regression-notify-channels-m2"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [2, 4, 6]
        active_end_time = "20:00"
        silence_time_secs = 43200
        active_start_time = "08:00"
        channels {
            type = "GROUP"
            identifiers = ["group-a", "group-b", "group-c"]
        }
        channels {
            type = "DINGTALK"
            identifiers = ["dt-url-1", "dt-url-2", "dt-url-3"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12923_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "cmsxx-prometheus-1"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up == 0"
    }

        display_name = "regression-notify-channels-m3"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5]
        active_end_time = "18:00"
        silence_time_secs = 86400
        active_start_time = "09:00"
        channels {
            type = "CONTACT"
            identifiers = ["admin-1", "admin-2"]
        }
        channels {
            type = "WEBHOOK"
            identifiers = ["https://hook1.example.com", "https://hook2.example.com"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12923_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "cmsxx-prometheus-1"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up == 0"
    }

        display_name = "regression-notify-channels-m4"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 2, 3, 4, 5, 6, 7]
        active_end_time = "23:59"
        silence_time_secs = 86400
        active_start_time = "00:00"
        channels {
            type = "GROUP"
            identifiers = ["ops-group-1", "ops-group-2"]
        }
        channels {
            type = "DINGTALK"
            identifiers = ["dt-1", "dt-2"]
        }
        channels {
            type = "FEISHU"
            identifiers = ["feishu-1", "feishu-2"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12923_step5(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 90
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "cmsxx-prometheus-1"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up == 0"
    }

        display_name = "regression-notify-channels-m5"

        enabled = true

        notify_config {
        utc_offset = "+08:00"
        type = "DIRECT_NOTIFY"
        active_days = [1, 3, 5, 7]
        active_end_time = "22:00"
        silence_time_secs = 86400
        active_start_time = "06:00"
        channels {
            type = "CONTACT"
            identifiers = ["dev-1", "dev-2", "dev-3"]
        }
        channels {
            type = "SLACK"
            identifiers = ["slack-1", "slack-2"]
        }
        channels {
            type = "WEIXIN"
            identifiers = ["wx-1", "wx-2", "wx-3"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_DataCheck_test 12922
func TestAccAliCloudCmsAlertRuleV2_basic12922(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12922_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.enable_data_complete_check", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-data-check"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12922_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.enable_data_complete_check", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-data-check-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test-updated"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12922_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
        enable_data_complete_check = true
    }

        display_name = "regression-data-check"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12922_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
        enable_data_complete_check = false
    }

        display_name = "regression-data-check-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test-updated"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_Condition_test 12921
func TestAccAliCloudCmsAlertRuleV2_basic12921(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12921_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "rate(http_errors_total[5m]) / rate(http_requests_total[5m]) > 0.05"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-prom-condition-rewrite"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12921_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "rate(http_errors_total[10m]) / rate(http_requests_total[10m]) > 0.1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-prom-condition-rewrite-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "180"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12921_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "rate(http_errors_total[5m]) / rate(http_requests_total[5m]) > 0.05"
    }

        display_name = "regression-prom-condition-rewrite"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12921_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "rate(http_errors_total[10m]) / rate(http_requests_total[10m]) > 0.1"
    }

        display_name = "regression-prom-condition-rewrite-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 180
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_EmptyArray_Repro_test 12920
func TestAccAliCloudCmsAlertRuleV2_basic12920(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12920_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "content_template", "empty array repro test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.1", "action-002"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "empty-array-repro-init"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12920_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "content_template", "empty array repro test modified"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-003"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "empty-array-repro-modify-nonempty"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "95"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12920_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "content_template", "empty array repro test emptied"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "0"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "empty-array-repro-modify-empty-array-repro"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "95"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12920_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "empty array repro test"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        actions = ["action-001", "action-002"]
        enabled = true
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "empty-array-repro-init"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12920_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "empty array repro test modified"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        actions = ["action-003"]
        enabled = true
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "empty-array-repro-modify-nonempty"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 95
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12920_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "empty array repro test emptied"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        actions = []
        enabled = true
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "empty-array-repro-modify-empty-array-repro"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 95
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_APM_YOY_test 12918
func TestAccAliCloudCmsAlertRuleV2_basic12918(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12918_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.status.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "status"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "status"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-svc-yoy-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-svc-yoy-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-yoy"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "YOY_UP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "AVG"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.yoy_time_value", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "30"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.yoy_time_unit", "minute"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12918_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.status.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "status"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "status"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-svc-yoy-3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-svc-yoy-4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.2", "arms-svc-yoy-5"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-yoy-week"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "YOY_DOWN"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "SUM"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.yoy_time_value", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "20"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.yoy_time_unit", "week"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12918_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.status.count"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.0", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.group_by.1", "status"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.0.key", "rpc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.filter_list.1.key", "status"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "arms-svc-yoy-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.1", "arms-svc-yoy-2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-yoy-month"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "YOY_UP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "APM_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.aggregate", "MAX"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.yoy_time_value", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold_list.0.threshold", "50"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.yoy_time_unit", "month"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12918_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.status.count"
            group_by = ["rpc", "status"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "ALL"
            key = "status"
        }
        service_id_list = ["arms-svc-yoy-1", "arms-svc-yoy-2"]
    }

        display_name = "regression-apm-yoy"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "YOY_UP"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "AVG"
        yoy_time_value = 5
        threshold_list {
            severity = "CRITICAL"
            threshold = 30
        }
        yoy_time_unit = "minute"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12918_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.status.count"
            group_by = ["rpc", "status"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "ALL"
            key = "status"
        }
        service_id_list = ["arms-svc-yoy-3", "arms-svc-yoy-4", "arms-svc-yoy-5"]
    }

        display_name = "regression-apm-yoy-week"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "YOY_DOWN"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "SUM"
        yoy_time_value = 2
        threshold_list {
            severity = "WARNING"
            threshold = 20
        }
        yoy_time_unit = "week"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12918_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "APM"
    }

        query_config {
        measure_list {
            measure_code = "appstat.status.count"
            group_by = ["rpc", "status"]
            window_secs = 60
        }
        type = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key = "rpc"
        }
        filter_list {
            type = "ALL"
            key = "status"
        }
        service_id_list = ["arms-svc-yoy-1", "arms-svc-yoy-2"]
    }

        display_name = "regression-apm-yoy-month"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "YOY_UP"
        type = "APM_SIMPLE_CONDITION"
        aggregate = "MAX"
        yoy_time_value = 1
        threshold_list {
            severity = "ERROR"
            threshold = 50
        }
        yoy_time_unit = "month"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_ExpressEscalation_test 12931
func TestAccAliCloudCmsAlertRuleV2_basic12931(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12931_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-prom-express-escalation"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12931_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-prom-express-escalation-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12931_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.2"
    }

        display_name = "regression-prom-express-escalation"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12931_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.1"
    }

        display_name = "regression-prom-express-escalation-updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_Basic_test 12933
func TestAccAliCloudCmsAlertRuleV2_basic12933(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12933_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-basic-prometheus"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12933_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-basic-prometheus-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test-updated"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12933_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "regression-basic-prometheus"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12933_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-basic-prometheus-updated"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test-updated"]
        }
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_lifecycle_test 12927
func TestAccAliCloudCmsAlertRuleV2_basic12927(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12927_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "cspec-lifecycle-test-rule"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12927_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "content_template"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.actions.0", "action-integration-001"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "cspec-lifecycle-test-rule-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "95"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12927_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "alert triggered on $${metric}"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        enabled = false
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "cspec-lifecycle-test-rule"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12927_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "alert updated on $${metric}"

        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        actions = ["action-integration-001"]
        enabled = true
    }

        arms_integration_config {
        enabled = true
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "cspec-lifecycle-test-rule-updated"

        enabled = false

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 120
        threshold = 95
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Prometheus_NotifyStrategies_test 12926
func TestAccAliCloudCmsAlertRuleV2_basic12926(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12926_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12926_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-id-003"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12926_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-4"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12926_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.0", "strategy-5"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12926_step4(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.region_id", "cn-hangzhou"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-notify-strategies-m4"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "NOTIFY_POLICY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.notify_strategies.#", "0"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12926_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        display_name = "regression-notify-strategies"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-1"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12926_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m1"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-id-003"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12926_step2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m2"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-4"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12926_step3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m3"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = ["strategy-5"]
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12926_step4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
        region_id = "cn-hangzhou"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        display_name = "regression-notify-strategies-m4"

        enabled = true

        notify_config {
        type = "NOTIFY_POLICY"
        notify_strategies = []
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case delete_minimal 12925
func TestAccAliCloudCmsAlertRuleV2_basic12925(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12925_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "content_template", "minimal delete test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "action_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "arms_integration_config.0.enabled", "false"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_type", "instance"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_METRICSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.entity_domain", "ecs"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric", "CPUUtilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.metric_set", "acs_ecs_dashboard"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "cspec-delete-minimal-test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_METRICSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.threshold", "90"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12925_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "minimal delete test"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        action_integration_config {
        enabled = false
    }

        arms_integration_config {
        enabled = false
    }

        query_config {
        entity_type = "instance"
        type = "UMODEL_METRICSET_QUERY"
        entity_domain = "ecs"
        window_secs = 60
        aggregate = "AVG"
        metric = "CPUUtilization"
        metric_set = "acs_ecs_dashboard"
    }

        display_name = "cspec-delete-minimal-test"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        operator = "GT"
        type = "UMODEL_METRICSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        threshold = 90
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Case resource_AlertRuleV2_Labels_Verify_test 12938
func TestAccAliCloudCmsAlertRuleV2_basic12938(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12938_step0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.%", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.summary", "CPU usage high"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.description", "Node CPU is above 90%"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "cloudspec_labels_verify_test"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.%", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.severity", "critical"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.env", "test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("datasource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("notify_strategy_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("observe_resource_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("partition_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("severity_levels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_cms_alert_rule_v2.default", tfjsonpath.New("updated_at"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_basic12938_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "120"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "PROMETHEUS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.instance_id", "prom-regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "PROMETHEUS_SINGLE_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "up > 2"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.%", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.summary", "CPU usage updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.description", "Updated CPU alert"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "annotations.runbook", "https://example.com"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "cloudspec_labels_verify_test_updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.%", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.severity", "warning"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.team", "dev"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "labels.env", "staging"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "PROMETHEUS_SIMPLE_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "120"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_basic12938_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 1"
    }

        annotations = {
        "summary" = "CPU usage high"
        "description" = "Node CPU is above 90%%"
    }

        display_name = "cloudspec_labels_verify_test"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        labels = {
        "severity" = "critical"
        "env" = "test"
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func testAccAliCloudCmsAlertRuleV2_basic12938_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        schedule_config {
        type = "FIXED"
        interval_secs = 120
    }

        datasource_config {
        type = "PROMETHEUS"
        instance_id = "prom-regression-test"
    }

        query_config {
        type = "PROMETHEUS_SINGLE_QUERY"
        expr = "up > 2"
    }

        annotations = {
        "summary" = "CPU usage updated"
        "description" = "Updated CPU alert"
        "runbook" = "https://example.com"
    }

        display_name = "cloudspec_labels_verify_test_updated"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        labels = {
        "severity" = "warning"
        "team" = "dev"
        "env" = "staging"
    }

        condition_config {
        type = "PROMETHEUS_SIMPLE_CONDITION"
        severity = "WARNING"
        duration_secs = 120
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

// Test Cms AlertRuleV2 log based alert rules (SLS datasource and UMODEL logset query).
// The CloudSpec regression cases above do not cover the log datasource (project/stores)
// or the logset query (log_set) added for the log based alert rule request.
func TestAccAliCloudCmsAlertRuleV2_slsLogDatasource(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2SlsLogDatasource(name, "sls log alert", "WARNING", 1, "INFO"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "SLS"),
					resource.TestCheckResourceAttrPair("alicloud_cms_alert_rule_v2.default", "datasource_config.0.project", "alicloud_log_project.default", "project_name"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.stores.#", "1"),
					resource.TestCheckResourceAttrPair("alicloud_cms_alert_rule_v2.default", "datasource_config.0.stores.0.store", "alicloud_log_store.default", "logstore_name"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.stores.0.store_type", "LOG"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "SLS_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.group_type", "none"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.0.expr", "* | select count(*) as cnt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.0.time_unit", "minute"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.0.window", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.type", "FIXED"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "schedule_config.0.interval_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "SLS_MULTI_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.0.operator", "HAS_DATA_MATCH_COUNT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.0.condition", "cnt > 0"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.0.count_condition", "__count__ > 0"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.alert_count", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.no_data_alert_level", "INFO"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "sls log alert"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.type", "DIRECT_NOTIFY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.type", "GROUP"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),

					resource.TestCheckResourceAttrSet("alicloud_cms_alert_rule_v2.default", "observe_resource_global_scope"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2SlsLogDatasource(name, "sls log alert updated", "ERROR", 2, "WARNING"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "SLS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.stores.#", "1"),
					resource.TestCheckResourceAttrPair("alicloud_cms_alert_rule_v2.default", "datasource_config.0.stores.0.store", "alicloud_log_store.default", "logstore_name"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "SLS_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.group_type", "none"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.0.operator", "HAS_DATA_MATCH_COUNT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.0.count_condition", "__count__ > 0"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.triggers.0.severity", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.alert_count", "2"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.no_data_alert_level", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "sls log alert updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2SlsLogDatasourceRange(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "SLS_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.0.time_unit", "minute"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.0.start", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.queries.0.end", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "sls log alert range"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2SlsLogDatasource(name, displayName, severity string, alertCount int, noDataAlertLevel string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
    description  = "tf-acc cms alert rule v2 sls"
}

resource "alicloud_log_store" "default" {
    provider = alicloud.hz
    project_name = alicloud_log_project.default.project_name
    logstore_name = var.name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "sls log alert rule"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "SLS"
        project = alicloud_log_project.default.project_name
        stores {
            store = alicloud_log_store.default.logstore_name
            store_type = "LOG"
        }
    }

        query_config {
        type = "SLS_MULTI_QUERY"
        group_type = "none"
        queries {
            expr = "* | select count(*) as cnt"
            time_unit = "minute"
            window = 60
        }
    }

        display_name = "%[2]s"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "SLS_MULTI_CONDITION"
        triggers {
            operator = "HAS_DATA_MATCH_COUNT"
            condition = "cnt > 0"
            count_condition = "__count__ > 0"
            severity = "%[3]s"
        }
        alert_count = %[4]d
        no_data_alert_level = "%[5]s"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name, displayName, severity, alertCount, noDataAlertLevel)
}

func testAccAliCloudCmsAlertRuleV2SlsLogDatasourceRange(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
    description  = "tf-acc cms alert rule v2 sls"
}

resource "alicloud_log_store" "default" {
    provider = alicloud.hz
    project_name = alicloud_log_project.default.project_name
    logstore_name = var.name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "sls log alert rule"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "SLS"
        project = alicloud_log_project.default.project_name
        stores {
            store = alicloud_log_store.default.logstore_name
            store_type = "LOG"
        }
    }

        query_config {
        type = "SLS_MULTI_QUERY"
        group_type = "none"
        queries {
            expr = "* | select count(*) as cnt"
            time_unit = "minute"
            start = 1
            end = 60
        }
    }

        display_name = "sls log alert range"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "SLS_MULTI_CONDITION"
        triggers {
            operator = "HAS_DATA_MATCH_COUNT"
            condition = "cnt > 0"
            count_condition = "__count__ > 0"
            severity = "WARNING"
        }
        alert_count = 1
        no_data_alert_level = "INFO"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name)
}

func TestAccAliCloudCmsAlertRuleV2_umodelLogSet(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2UmodelLogSet(name, "umodel logset alert"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "UMODEL"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "UMODEL_LOGSET_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.log_set", name),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.expr", "* | select count(*) as cnt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.window_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.offset_secs", "60"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.type", "UMODEL_LOGSET_CONDITION"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.duration_secs", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.match_operator", "EQUALS"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.match_field", "level"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.match_value", "ERROR"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.count_threshold", "1"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "umodel logset alert"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "true"),

					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "notify_config.0.channels.0.identifiers.0", "regression-test"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2UmodelLogSet(name, "umodel logset alert updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.log_set", name),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "umodel logset alert updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_rule_v2.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2UmodelLogSet(name, displayName string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
        provider = alicloud.hz
        content_template = "umodel logset alert rule"

        schedule_config {
        type = "FIXED"
        interval_secs = 60
    }

        datasource_config {
        type = "UMODEL"
    }

        query_config {
        type = "UMODEL_LOGSET_QUERY"
        log_set = var.name
        expr = "* | select count(*) as cnt"
        window_secs = 60
        offset_secs = 60
    }

        display_name = "%[2]s"

        enabled = true

        notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type = "GROUP"
            identifiers = ["regression-test"]
        }
    }

        condition_config {
        type = "UMODEL_LOGSET_CONDITION"
        severity = "CRITICAL"
        duration_secs = 60
        count_operator = "GT"
        count_threshold = 1
        match_field = "level"
        match_operator = "EQUALS"
        match_value = "ERROR"
    }

        workspace = alicloud_cms_workspace.default.id
}
`, name, displayName)
}

// Case resource_AlertRuleV2_ObserveResource_Rca_test validates the observe_resource_config and rca_config blocks.
func TestAccAliCloudCmsAlertRuleV2_observeResourceRca(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", r)
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkCmsAlertRuleV2Destroy,
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudCmsAlertRuleV2_observeResourceRca(name, "apm_ns_coverage", "apm", "CLOUD_INSTANCE", "tf-probe-agent"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.type", "APM"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.type", "APM_MULTI_QUERY"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.measure_list.0.measure_code", "appstat.transaction.rt"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.service_id_list.0", "svc-coverage-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.entity_domain", "apm"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.entity_type", "apm.service"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.relation_type", "CLOUD_INSTANCE"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.namespace", "apm_ns_coverage"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.product_category", "apm"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.resources.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.resources.0", "svc-coverage-1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "rca_config.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "rca_config.0.enable_rca", "true"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "rca_config.0.digital_employee_name", "tf-probe-agent"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "display_name", "regression-apm-observe-rca"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "enabled", "false"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "workspace", name),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_observeResourceRca(name, "apm_ns_coverage_updated", "apm_updated", "ALL", "tf-probe-agent"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.relation_type", "ALL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.namespace", "apm_ns_coverage_updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.product_category", "apm_updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.entity_domain", "apm"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.entity_type", "apm.service"),
				),
			},
			{
				Config: testAccAliCloudCmsAlertRuleV2_observeResourceRca(name, "apm_ns_coverage_updated", "apm_updated", "ALL", "tf-probe-agent-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "rca_config.0.enable_rca", "true"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "rca_config.0.digital_employee_name", "tf-probe-agent-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "observe_resource_config.0.relation_type", "ALL"),
				),
			},
			{
				ResourceName:      "alicloud_cms_alert_rule_v2.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudCmsAlertRuleV2_observeResourceRca(name, namespace, productCategory, relationType, rcaName string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%[1]s"
}

provider "alicloud" {
    alias  = "hz"
    region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
    provider = alicloud.hz
    project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
    provider = alicloud.hz
    workspace_name = var.name
    sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_rule_v2" "default" {
    provider = alicloud.hz
    schedule_config {
        type          = "FIXED"
        interval_secs = 60
    }
    datasource_config {
        type = "APM"
    }
    query_config {
        measure_list {
            measure_code = "appstat.transaction.rt"
            group_by     = ["rpc"]
            window_secs  = 60
        }
        type            = "APM_MULTI_QUERY"
        filter_list {
            type = "ALL"
            key  = "rpc"
        }
        service_id_list = ["svc-coverage-1"]
    }
    observe_resource_config {
        entity_domain    = "apm"
        entity_type      = "apm.service"
        namespace        = "%[2]s"
        product_category = "%[3]s"
        relation_type    = "%[4]s"
        resources        = ["svc-coverage-1"]
    }
    rca_config {
        enable_rca           = true
        digital_employee_name = "%[5]s"
    }
    display_name = "regression-apm-observe-rca"
    enabled      = false
    notify_config {
        type = "DIRECT_NOTIFY"
        channels {
            type        = "GROUP"
            identifiers = ["regression-test"]
        }
    }
    condition_config {
        operator  = "GT"
        type      = "APM_SIMPLE_CONDITION"
        aggregate = "AVG"
        threshold_list {
            severity  = "WARNING"
            threshold = 100
        }
    }
    workspace = alicloud_cms_workspace.default.id
}
`, name, namespace, productCategory, relationType, rcaName)
}

// Test Cms AlertRuleV2. <<< Resource test cases, automatically generated.

var AlicloudCmsAlertRuleV2Map12929 = map[string]string{
	"created_at":                    acctest.CHECKSET,
	"datasource_type":               acctest.CHECKSET,
	"status":                        acctest.CHECKSET,
	"observe_resource_global_scope": acctest.CHECKSET,
	"severity_levels":               acctest.CHECKSET,
	"updated_at":                    acctest.CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12929(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// TypeList order coverage for alicloud_cms_alert_rule_v2 query_config.label_filters.
func TestAccAliCloudCmsAlertRuleV2_labelFiltersOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
							"label_filters": []map[string]interface{}{
								{
									"operator": "=",
									"value":    "web-server",
									"name":     "app",
								},
								{
									"operator": "=",
									"value":    "production",
									"name":     "env",
								},
							},
						},
					},
					"display_name": "regression-umodel-labelorder",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.label_filters.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
							"label_filters": []map[string]interface{}{
								{
									"operator": "=",
									"value":    "production",
									"name":     "env",
								},
								{
									"operator": "=",
									"value":    "web-server",
									"name":     "app",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.label_filters.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2 query_config.entity_filters.
func TestAccAliCloudCmsAlertRuleV2_entityFiltersOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
							"entity_filters": []map[string]interface{}{
								{
									"operator": "=",
									"field":    "regionId",
									"value":    "cn-hangzhou",
								},
								{
									"operator": "=",
									"field":    "instanceId",
									"value":    "i-bp1abc123def",
								},
							},
						},
					},
					"display_name": "regression-umodel-entityfilterorder",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.entity_filters.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
							"entity_filters": []map[string]interface{}{
								{
									"operator": "=",
									"field":    "instanceId",
									"value":    "i-bp1abc123def",
								},
								{
									"operator": "=",
									"field":    "regionId",
									"value":    "cn-hangzhou",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.entity_filters.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2 query_config.entity_fields.
func TestAccAliCloudCmsAlertRuleV2_entityFieldsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
							"entity_fields": []map[string]interface{}{
								{
									"field": "name",
									"value": "my-instance",
								},
								{
									"field": "status",
									"value": "running",
								},
							},
						},
					},
					"display_name": "regression-umodel-entityfieldorder",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.entity_fields.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
							"entity_fields": []map[string]interface{}{
								{
									"field": "status",
									"value": "running",
								},
								{
									"field": "name",
									"value": "my-instance",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.entity_fields.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2 query_config.measure_list.
func TestAccAliCloudCmsAlertRuleV2_measureListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1", "arms-agg-2"},
							"filter_list": []map[string]interface{}{
								{"type": "ALL", "key": "rpc"},
								{"type": "ALL", "key": "rpcType"},
							},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
								{
									"measure_code": "appstat.transaction.errorrate",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
							},
						},
					},
					"display_name": "regression-apm-measurelistorder",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":     "APM_COMPOSITE_CONDITION",
							"relation": "OR",
							"severity": "CRITICAL",
							"compare_list": []map[string]interface{}{
								{
									"operator":       "YOY_UP",
									"aggregate":      "AVG",
									"yoy_time_value": "1",
									"yoy_time_unit":  "hour",
									"threshold":      "0",
								},
								{
									"operator":       "YOY_DOWN",
									"aggregate":      "MAX",
									"yoy_time_value": "2",
									"yoy_time_unit":  "day",
									"threshold":      "0",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.measure_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1", "arms-agg-2"},
							"filter_list": []map[string]interface{}{
								{"type": "ALL", "key": "rpc"},
								{"type": "ALL", "key": "rpcType"},
							},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.errorrate",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.measure_list.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2 query_config.filter_list.
func TestAccAliCloudCmsAlertRuleV2_filterListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1"},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
							},
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "rpc",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "rpcType",
								},
							},
						},
					},
					"display_name": "regression-apm-filterlistorder",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.filter_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1"},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
							},
							"filter_list": []map[string]interface{}{
								{
									"type":  "EQ",
									"value": "http",
									"key":   "rpcType",
								},
								{
									"type": "ALL",
									"key":  "rpc",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.filter_list.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2 query_config.service_id_list.
func TestAccAliCloudCmsAlertRuleV2_serviceIdListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "APM_MULTI_QUERY",
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
							},
							"service_id_list": []string{
								"arms-agg-1", "arms-agg-2"},
						},
					},
					"display_name": "regression-apm-serviceidlistorder",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.service_id_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type": "APM_MULTI_QUERY",
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
							},
							"service_id_list": []string{
								"arms-agg-2", "arms-agg-1"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.service_id_list.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2 query_config.measure_list.group_by.
func TestAccAliCloudCmsAlertRuleV2_measureListGroupByOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1"},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpc", "rpcType"},
									"window_secs": "60",
								},
							},
						},
					},
					"display_name": "regression-apm-groupbyorder",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.measure_list.0.group_by.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1"},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by": []string{
										"rpcType", "rpc"},
									"window_secs": "60",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.measure_list.0.group_by.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.action_integration_config.actions.
func TestAccAliCloudCmsAlertRuleV2_actionIntegrationConfigActionsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "UMODEL"},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-actions-order",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-integration-001",
								"action-integration-002",
							},
							"enabled": "true",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_integration_config.0.actions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-integration-002",
								"action-integration-001",
							},
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_integration_config.0.actions.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.condition_config.compare_list.
func TestAccAliCloudCmsAlertRuleV2_conditionConfigCompareListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "APM"},
					},
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1"},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by":     []string{"rpc", "rpcType"},
									"window_secs":  "60",
								},
								{
									"measure_code": "appstat.transaction.errorrate",
									"group_by":     []string{"rpc", "rpcType"},
									"window_secs":  "60",
								},
							},
							"filter_list": []map[string]interface{}{
								{"type": "ALL", "key": "rpc"},
								{"type": "ALL", "key": "rpcType"},
							},
						},
					},
					"display_name": "regression-comparelist-order",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":     "APM_COMPOSITE_CONDITION",
							"relation": "OR",
							"severity": "CRITICAL",
							"compare_list": []map[string]interface{}{
								{
									"operator":       "YOY_UP",
									"aggregate":      "AVG",
									"yoy_time_value": "1",
									"threshold":      "0",
									"yoy_time_unit":  "hour",
								},
								{
									"operator":       "YOY_DOWN",
									"aggregate":      "MAX",
									"yoy_time_value": "2",
									"threshold":      "0",
									"yoy_time_unit":  "day",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.compare_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"type":     "APM_COMPOSITE_CONDITION",
							"relation": "OR",
							"severity": "CRITICAL",
							"compare_list": []map[string]interface{}{
								{
									"operator":       "YOY_DOWN",
									"aggregate":      "MAX",
									"yoy_time_value": "2",
									"threshold":      "0",
									"yoy_time_unit":  "day",
								},
								{
									"operator":       "YOY_UP",
									"aggregate":      "AVG",
									"yoy_time_value": "1",
									"threshold":      "0",
									"yoy_time_unit":  "hour",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.compare_list.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.condition_config.threshold_list.
func TestAccAliCloudCmsAlertRuleV2_conditionConfigThresholdListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "APM"},
					},
					"query_config": []map[string]interface{}{
						{
							"type":            "APM_MULTI_QUERY",
							"service_id_list": []string{"arms-agg-1"},
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "appstat.transaction.rt",
									"group_by":     []string{"rpc", "rpcType"},
									"window_secs":  "60",
								},
							},
							"filter_list": []map[string]interface{}{
								{"type": "ALL", "key": "rpc"},
							},
						},
					},
					"display_name": "regression-thresholdlist-order",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.threshold_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.threshold_list.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.notify_config.active_days.
func TestAccAliCloudCmsAlertRuleV2_notifyConfigActiveDaysOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "UMODEL"},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-activedays-order",
					"enabled":      "true",
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"notify_config": []map[string]interface{}{
						{
							"type":        "DIRECT_NOTIFY",
							"active_days": []string{"1", "2"},
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_config.0.active_days.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_config": []map[string]interface{}{
						{
							"type":        "DIRECT_NOTIFY",
							"active_days": []string{"2", "1"},
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_config.0.active_days.#": "2",
					}),
				),
			},
		},
	})
}

// collection-order: bypass notify_config.channels -- the CMS 2.0 API merges same-type channels into one member and drops a CONTACT channel whose identifier does not resolve, so a two-member reorder apply cannot converge (ACC task 25203).

// TypeList order coverage for alicloud_cms_alert_rule_v2.notify_config.channels.identifiers.
func TestAccAliCloudCmsAlertRuleV2_notifyConfigChannelsIdentifiersOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "UMODEL"},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"window_secs":   "60",
							"aggregate":     "AVG",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-identifiers-order",
					"enabled":      "true",
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"reg-A", "reg-B"}},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_config.0.channels.0.identifiers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"reg-B", "reg-A"}},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_config.0.channels.0.identifiers.#": "2",
					}),
				),
			},
		},
	})
}

// collection-order: bypass notify_config.notify_strategies -- the CMS 2.0 API rejects more than one strategy ("NotifyPolicy.notifyStrategies size must be <= 1", RequestId 01A19A04-B6BF-5C39-A671-0A0588186571), so two distinct members cannot exist.

// TypeList order coverage for alicloud_cms_alert_rule_v2.query_config.dimensions.
func TestAccAliCloudCmsAlertRuleV2_queryConfigDimensionsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"content_template": "dimensions order coverage",
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "CLOUD_MONITORING", "product_category": "ecs"},
					},
					"query_config": []map[string]interface{}{
						{
							"type":      "CLOUD_MONITORING_QUERY",
							"namespace": "acs_ecs_dashboard",
							"dimensions": []map[string]interface{}{
								{"instanceId": "i-abc123"},
								{"device": "/dev/vda"},
							},
						},
					},
					"display_name": name,
					"enabled":      false,
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "SIMPLE",
							"simple_escalation": []map[string]interface{}{
								{
									"metric_name": "cpu_total",
									"period":      "60",
									"escalations": []map[string]interface{}{
										{
											"severity":            "WARNING",
											"comparison_operator": "GT",
											"threshold":           "80",
											"times":               "3",
										},
									},
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.dimensions.#": "2",
					}),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.dimensions.0.instanceId", "i-abc123"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.dimensions.1.device", "/dev/vda"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.namespace", "acs_ecs_dashboard"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "datasource_config.0.product_category", "ecs"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type":      "CLOUD_MONITORING_QUERY",
							"namespace": "acs_ecs_dashboard",
							"dimensions": []map[string]interface{}{
								{"device": "/dev/vda"},
								{"instanceId": "i-abc123"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.dimensions.#": "2",
					}),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.dimensions.0.device", "/dev/vda"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "query_config.0.dimensions.1.instanceId", "i-abc123"),
				),
			},
		},
	})
}

// Covers query_config.relation_type/group_id (type=CLOUD_MONITORING_QUERY relationType=GROUP) and
// the updatable query_config.namespace and datasource_config.product_category attributes.
func TestAccAliCloudCmsAlertRuleV2_queryConfigGroupRelation(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"content_template": "cloud monitoring group relation coverage",
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "CLOUD_MONITORING", "product_category": "ecs"},
					},
					"query_config": []map[string]interface{}{
						{
							"type":          "CLOUD_MONITORING_QUERY",
							"namespace":     "acs_ecs_dashboard",
							"relation_type": "GROUP",
							"group_id":      "og-tfacc16955",
						},
					},
					"display_name": name,
					"enabled":      false,
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "SIMPLE",
							"simple_escalation": []map[string]interface{}{
								{
									"metric_name": "cpu_total",
									"period":      "60",
									"escalations": []map[string]interface{}{
										{"severity": "WARNING", "comparison_operator": "GT", "threshold": "80", "times": "3"},
									},
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.relation_type":         "GROUP",
						"query_config.0.group_id":              "og-tfacc16955",
						"query_config.0.namespace":             "acs_ecs_dashboard",
						"datasource_config.0.product_category": "ecs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"datasource_config": []map[string]interface{}{
						{"type": "CLOUD_MONITORING", "product_category": "rds"},
					},
					"query_config": []map[string]interface{}{
						{
							"type":          "CLOUD_MONITORING_QUERY",
							"namespace":     "acs_rds_dashboard",
							"relation_type": "GROUP",
							"group_id":      "og-tfacc16956",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.relation_type":         "GROUP",
						"query_config.0.group_id":              "og-tfacc16956",
						"query_config.0.namespace":             "acs_rds_dashboard",
						"datasource_config.0.product_category": "rds",
					}),
				),
			},
		},
	})
}

func AlicloudCmsAlertRuleV2LogDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
    project_name = var.name
    description  = "tf-acc cms alert rule v2 order coverage"
}

resource "alicloud_log_store" "first" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = "${var.name}-first"
}

resource "alicloud_log_store" "second" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = "${var.name}-second"
}
`, name)
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.datasource_config.stores.
func TestAccAliCloudCmsAlertRuleV2_datasourceConfigStoresOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2LogDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "${alicloud_log_project.default.project_name}",
							"stores": []map[string]interface{}{
								{"store": "${alicloud_log_store.first.logstore_name}", "store_type": "LOG"},
								{"store": "${alicloud_log_store.second.logstore_name}", "store_type": "LOG"},
							},
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
							},
						},
					},
					"display_name": "regression-stores-order",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type": "SLS_MULTI_CONDITION",
							"triggers": []map[string]interface{}{
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 0", "count_condition": "__count__ > 0", "severity": "WARNING"},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"datasource_config.0.stores.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"datasource_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "${alicloud_log_project.default.project_name}",
							"stores": []map[string]interface{}{
								{"store": "${alicloud_log_store.second.logstore_name}", "store_type": "LOG"},
								{"store": "${alicloud_log_store.first.logstore_name}", "store_type": "LOG"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"datasource_config.0.stores.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.query_config.queries.
func TestAccAliCloudCmsAlertRuleV2_queryConfigQueriesOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2LogDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "${alicloud_log_project.default.project_name}",
							"stores": []map[string]interface{}{
								{"store": "${alicloud_log_store.first.logstore_name}", "store_type": "LOG"},
							},
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt2", "time_unit": "minute", "window": "60"},
							},
						},
					},
					"display_name": "regression-queries-order",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type": "SLS_MULTI_CONDITION",
							"triggers": []map[string]interface{}{
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 0", "count_condition": "__count__ > 0", "severity": "WARNING"},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.queries.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt2", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.queries.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.condition_config.triggers.
func TestAccAliCloudCmsAlertRuleV2_conditionConfigTriggersOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2LogDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "${alicloud_log_project.default.project_name}",
							"stores": []map[string]interface{}{
								{"store": "${alicloud_log_store.first.logstore_name}", "store_type": "LOG"},
							},
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
							},
						},
					},
					"display_name": "regression-triggers-order",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type": "SLS_MULTI_CONDITION",
							"triggers": []map[string]interface{}{
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 0", "count_condition": "__count__ > 0", "severity": "WARNING"},
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 10", "count_condition": "__count__ > 0", "severity": "CRITICAL"},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.triggers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"type": "SLS_MULTI_CONDITION",
							"triggers": []map[string]interface{}{
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 10", "count_condition": "__count__ > 0", "severity": "CRITICAL"},
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 0", "count_condition": "__count__ > 0", "severity": "WARNING"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.triggers.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.observe_resource_config.resources.
func TestAccAliCloudCmsAlertRuleV2_observeResourceConfigResourcesOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "APM"},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "APM_MULTI_QUERY",
							"measure_list": []map[string]interface{}{
								{"measure_code": "appstat.transaction.rt", "group_by": []string{"rpc"}, "window_secs": "60"},
							},
							"filter_list": []map[string]interface{}{
								{"type": "ALL", "key": "rpc"},
							},
							"service_id_list": []string{"svc-coverage-1"},
						},
					},
					"observe_resource_config": []map[string]interface{}{
						{
							"entity_domain": "apm",
							"entity_type":   "apm.service",
							"relation_type": "CLOUD_INSTANCE",
							"resources":     []string{"svc-coverage-1", "svc-coverage-2"},
						},
					},
					"display_name": "regression-observe-resources-order",
					"enabled":      "false",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "AVG",
							"threshold_list": []map[string]interface{}{
								{"severity": "WARNING", "threshold": "100"},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"observe_resource_config.0.resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"observe_resource_config": []map[string]interface{}{
						{
							"entity_domain": "apm",
							"entity_type":   "apm.service",
							"relation_type": "CLOUD_INSTANCE",
							"resources":     []string{"svc-coverage-2", "svc-coverage-1"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"observe_resource_config.0.resources.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.condition_config.simple_escalation.escalations.
func TestAccAliCloudCmsAlertRuleV2_conditionConfigSimpleEscalationEscalationsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "CLOUD_MONITORING", "product_category": "ecs"},
					},
					"query_config": []map[string]interface{}{
						{
							"type":      "CLOUD_MONITORING_QUERY",
							"namespace": "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-simple-escalations-order",
					"enabled":      "false",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "SIMPLE",
							"simple_escalation": []map[string]interface{}{
								{
									"metric_name": "cpu_total",
									"period":      "60",
									"escalations": []map[string]interface{}{
										{"severity": "WARNING", "comparison_operator": "GT", "threshold": "80", "times": "3", "pre_condition": "pre-warning", "statistics": "Average"},
										{"severity": "CRITICAL", "comparison_operator": "GT", "threshold": "90", "times": "3", "pre_condition": "pre-critical", "statistics": "Maximum"},
									},
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.simple_escalation.0.escalations.#": "2",
					}),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.metric_name", "cpu_total"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.period", "60"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.threshold", "80"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.pre_condition", "pre-warning"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.statistics", "Average"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.comparison_operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.times", "3"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.threshold", "90"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.pre_condition", "pre-critical"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.statistics", "Maximum"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.comparison_operator", "GT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.times", "3"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "SIMPLE",
							"simple_escalation": []map[string]interface{}{
								{
									"metric_name": "cpu_total",
									"period":      "60",
									"escalations": []map[string]interface{}{
										{"severity": "CRITICAL", "comparison_operator": "GT", "threshold": "90", "times": "3", "pre_condition": "pre-critical", "statistics": "Maximum"},
										{"severity": "WARNING", "comparison_operator": "GT", "threshold": "80", "times": "3", "pre_condition": "pre-warning", "statistics": "Average"},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.simple_escalation.0.escalations.#": "2",
					}),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.threshold", "90"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.threshold", "80"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "SIMPLE",
							"simple_escalation": []map[string]interface{}{
								{
									"metric_name": "memory_usedutilization",
									"period":      "120",
									"escalations": []map[string]interface{}{
										{"severity": "CRITICAL", "comparison_operator": "LT", "threshold": "95", "times": "5", "pre_condition": "pre-critical-updated", "statistics": "Maximum"},
										{"severity": "WARNING", "comparison_operator": "LT", "threshold": "85", "times": "5", "pre_condition": "pre-warning-updated", "statistics": "Average"},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.simple_escalation.0.escalations.#": "2",
					}),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.metric_name", "memory_usedutilization"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.period", "120"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.threshold", "95"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.pre_condition", "pre-critical-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.statistics", "Maximum"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.comparison_operator", "LT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.0.times", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.severity", "WARNING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.threshold", "85"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.pre_condition", "pre-warning-updated"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.statistics", "Average"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.comparison_operator", "LT"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.simple_escalation.0.escalations.1.times", "5"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "COMPOSITE",
							"composite_escalation": []map[string]interface{}{
								{
									"relation": "and",
									"severity": "CRITICAL",
									"times":    "3",
									"escalations": []map[string]interface{}{
										{"metric_name": "cpu_total", "comparison_operator": "GT", "threshold": "88", "period": "60", "pre_condition": "pre-switch", "statistics": "Average"},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.escalation_type", "COMPOSITE"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.composite_escalation.0.relation", "and"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.composite_escalation.0.escalations.0.metric_name", "cpu_total"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.composite_escalation.0.escalations.0.pre_condition", "pre-switch"),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.condition_config.composite_escalation.escalations.
// The CMS 2.0 API echoes the submitted member order but normalizes relation to lower case, so the
// configuration uses "and" to stay convergent.
func TestAccAliCloudCmsAlertRuleV2_conditionConfigCompositeEscalationEscalationsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{"type": "CLOUD_MONITORING", "product_category": "ecs"},
					},
					"query_config": []map[string]interface{}{
						{
							"type":      "CLOUD_MONITORING_QUERY",
							"namespace": "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-composite-escalations-order",
					"enabled":      "false",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "COMPOSITE",
							"no_data_policy":  "NO_DATA_TO_OK",
							"composite_escalation": []map[string]interface{}{
								{
									"relation": "and",
									"severity": "CRITICAL",
									"times":    "3",
									"escalations": []map[string]interface{}{
										{"metric_name": "cpu_total", "comparison_operator": "GT", "threshold": "80", "period": "60", "pre_condition": "pre-cpu", "statistics": "Average"},
										{"metric_name": "memory_usedutilization", "comparison_operator": "GT", "threshold": "70", "period": "300", "pre_condition": "pre-memory", "statistics": "Maximum"},
									},
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.composite_escalation.0.escalations.#": "2",
						"condition_config.0.composite_escalation.0.relation":      "and",
						"condition_config.0.composite_escalation.0.severity":      "CRITICAL",
						"condition_config.0.composite_escalation.0.times":         "3",
						"condition_config.0.no_data_policy":                       "NO_DATA_TO_OK",
					}),
					checkCmsAlertRuleV2EscalationsByMetric("condition_config.0.composite_escalation.0.escalations", map[string]map[string]string{
						"cpu_total":              {"period": "60", "pre_condition": "pre-cpu", "statistics": "Average", "threshold": "80"},
						"memory_usedutilization": {"period": "300", "pre_condition": "pre-memory", "statistics": "Maximum", "threshold": "70"},
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "COMPOSITE",
							"no_data_policy":  "NO_DATA_TO_OK",
							"composite_escalation": []map[string]interface{}{
								{
									"relation": "and",
									"severity": "CRITICAL",
									"times":    "3",
									"escalations": []map[string]interface{}{
										{"metric_name": "memory_usedutilization", "comparison_operator": "GT", "threshold": "70", "period": "300", "pre_condition": "pre-memory", "statistics": "Maximum"},
										{"metric_name": "cpu_total", "comparison_operator": "GT", "threshold": "80", "period": "60", "pre_condition": "pre-cpu", "statistics": "Average"},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.composite_escalation.0.escalations.#": "2",
					}),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.composite_escalation.0.escalations.#", "2"),
					checkCmsAlertRuleV2EscalationsByMetric("condition_config.0.composite_escalation.0.escalations", map[string]map[string]string{
						"cpu_total":              {"period": "60", "pre_condition": "pre-cpu", "statistics": "Average"},
						"memory_usedutilization": {"period": "300", "pre_condition": "pre-memory", "statistics": "Maximum"},
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition_config": []map[string]interface{}{
						{
							"type":            "CLOUD_MONITORING_CONDITION",
							"escalation_type": "COMPOSITE",
							"no_data_policy":  "KEEP_LAST_STATE",
							"composite_escalation": []map[string]interface{}{
								{
									"relation": "or",
									"severity": "ERROR",
									"times":    "5",
									"escalations": []map[string]interface{}{
										{"metric_name": "memory_usedutilization", "comparison_operator": "LT", "threshold": "75", "period": "120", "pre_condition": "pre-memory-updated", "statistics": "Maximum"},
										{"metric_name": "cpu_total", "comparison_operator": "LT", "threshold": "85", "period": "120", "pre_condition": "pre-cpu-updated", "statistics": "Average"},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition_config.0.composite_escalation.0.escalations.#": "2",
						"condition_config.0.composite_escalation.0.relation":      "or",
						"condition_config.0.composite_escalation.0.severity":      "ERROR",
						"condition_config.0.composite_escalation.0.times":         "5",
						"condition_config.0.no_data_policy":                       "KEEP_LAST_STATE",
					}),
					resource.TestCheckResourceAttr("alicloud_cms_alert_rule_v2.default", "condition_config.0.composite_escalation.0.escalations.#", "2"),
					checkCmsAlertRuleV2EscalationsByMetric("condition_config.0.composite_escalation.0.escalations", map[string]map[string]string{
						"cpu_total":              {"period": "120", "pre_condition": "pre-cpu-updated", "statistics": "Average", "threshold": "85"},
						"memory_usedutilization": {"period": "120", "pre_condition": "pre-memory-updated", "statistics": "Maximum", "threshold": "75"},
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.query_config.joinings.
func TestAccAliCloudCmsAlertRuleV2_queryConfigJoiningsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2LogDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "${alicloud_log_project.default.project_name}",
							"stores": []map[string]interface{}{
								{"store": "${alicloud_log_store.first.logstore_name}", "store_type": "LOG"},
							},
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt2", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt3", "time_unit": "minute", "window": "60"},
							},
							"joinings": []map[string]interface{}{
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$0.cnt", "operator": "==", "rhs_field": "$1.cnt2"},
									},
								},
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$1.cnt2", "operator": "==", "rhs_field": "$2.cnt3"},
									},
								},
							},
						},
					},
					"display_name": "regression-joinings-order",
					"enabled":      "false",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type": "SLS_MULTI_CONDITION",
							"triggers": []map[string]interface{}{
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 0", "count_condition": "__count__ > 0", "severity": "WARNING"},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.joinings.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt2", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt3", "time_unit": "minute", "window": "60"},
							},
							"joinings": []map[string]interface{}{
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$1.cnt2", "operator": "==", "rhs_field": "$2.cnt3"},
									},
								},
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$0.cnt", "operator": "==", "rhs_field": "$1.cnt2"},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.joinings.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_alert_rule_v2.query_config.joinings.conditions.
func TestAccAliCloudCmsAlertRuleV2_queryConfigJoiningsConditionsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsAlertRuleV2")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsAlertRuleV2LogDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			requireCmsAlertRuleV2ProtoV5Factory(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{"type": "FIXED", "interval_secs": "60"},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":    "SLS",
							"project": "${alicloud_log_project.default.project_name}",
							"stores": []map[string]interface{}{
								{"store": "${alicloud_log_store.first.logstore_name}", "store_type": "LOG"},
							},
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt2", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt3", "time_unit": "minute", "window": "60"},
							},
							"joinings": []map[string]interface{}{
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$0.cnt", "operator": "==", "rhs_field": "$1.cnt2"},
										{"lhs_field": "$0.cnt", "operator": ">=", "rhs_field": "$1.cnt2"},
									},
								},
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$1.cnt2", "operator": "==", "rhs_field": "$2.cnt3"},
									},
								},
							},
						},
					},
					"display_name": "regression-joinings-conditions-order",
					"enabled":      "false",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{"type": "GROUP", "identifiers": []string{"regression-test"}},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type": "SLS_MULTI_CONDITION",
							"triggers": []map[string]interface{}{
								{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 0", "count_condition": "__count__ > 0", "severity": "WARNING"},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.joinings.0.conditions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type":       "SLS_MULTI_QUERY",
							"group_type": "none",
							"queries": []map[string]interface{}{
								{"expr": "* | select count(*) as cnt", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt2", "time_unit": "minute", "window": "60"},
								{"expr": "* | select count(*) as cnt3", "time_unit": "minute", "window": "60"},
							},
							"joinings": []map[string]interface{}{
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$0.cnt", "operator": ">=", "rhs_field": "$1.cnt2"},
										{"lhs_field": "$0.cnt", "operator": "==", "rhs_field": "$1.cnt2"},
									},
								},
								{
									"type": "INNER",
									"conditions": []map[string]interface{}{
										{"lhs_field": "$1.cnt2", "operator": "==", "rhs_field": "$2.cnt3"},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_config.0.joinings.0.conditions.#": "2",
					}),
				),
			},
		},
	})
}

func TestUnitCmsAlertRuleV2CustomizeDiffConditionalRequirements(t *testing.T) {
	resourceSchema, ok := acctest.Provider.ResourcesMap["alicloud_cms_alert_rule_v2"]
	if !ok {
		t.Fatal("alicloud_cms_alert_rule_v2 is not registered in the provider")
	}

	validConfig := func() map[string]interface{} {
		return map[string]interface{}{
			"workspace": "tf-testacc-customize-diff",
			"schedule_config": []interface{}{
				map[string]interface{}{"type": "FIXED", "interval_secs": 60},
			},
			"datasource_config": []interface{}{
				map[string]interface{}{
					"type":    "SLS",
					"project": "tf-testacc-customize-diff",
					"stores": []interface{}{
						map[string]interface{}{"store": "tf-testacc-store", "store_type": "LOG"},
					},
				},
			},
			"query_config": []interface{}{
				map[string]interface{}{"type": "SLS_MULTI_QUERY", "group_type": "none"},
			},
			"condition_config": []interface{}{
				map[string]interface{}{
					"type": "SLS_MULTI_CONDITION",
					"triggers": []interface{}{
						map[string]interface{}{"operator": "HAS_DATA_MATCH_COUNT", "condition": "cnt > 0"},
					},
				},
			},
			"notify_config": []interface{}{
				map[string]interface{}{
					"type": "DIRECT_NOTIFY",
					"channels": []interface{}{
						map[string]interface{}{"type": "GROUP", "identifiers": []interface{}{"tf-testacc-group"}},
					},
				},
			},
		}
	}

	cases := []struct {
		name   string
		mutate func(map[string]interface{})
		want   string
	}{
		{
			name:   "missing datasource_config",
			mutate: func(config map[string]interface{}) { delete(config, "datasource_config") },
			want:   "datasource_config is required",
		},
		{
			name: "sls datasource without project",
			mutate: func(config map[string]interface{}) {
				config["datasource_config"] = []interface{}{map[string]interface{}{"type": "SLS"}}
			},
			want: "datasource_config.project is required when datasource_config.type is SLS",
		},
		{
			name: "sls datasource without stores",
			mutate: func(config map[string]interface{}) {
				config["datasource_config"] = []interface{}{
					map[string]interface{}{"type": "SLS", "project": "tf-testacc-customize-diff"},
				}
			},
			want: "datasource_config.stores requires at least one store when datasource_config.type is SLS",
		},
		{
			name: "umodel logset query without log set",
			mutate: func(config map[string]interface{}) {
				config["query_config"] = []interface{}{map[string]interface{}{"type": "UMODEL_LOGSET_QUERY"}}
			},
			want: "query_config.log_set is required when query_config.type is UMODEL_LOGSET_QUERY",
		},
		{
			name: "umodel logset query without expr",
			mutate: func(config map[string]interface{}) {
				config["query_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_LOGSET_QUERY", "log_set": "tf-testacc-logset"},
				}
			},
			want: "query_config.expr is required when query_config.type is UMODEL_LOGSET_QUERY",
		},
		{
			name: "umodel logset query without window secs",
			mutate: func(config map[string]interface{}) {
				config["query_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_LOGSET_QUERY", "log_set": "tf-testacc-logset", "expr": "cnt > 0"},
				}
			},
			want: "query_config.window_secs is required when query_config.type is UMODEL_LOGSET_QUERY",
		},
		{
			name: "umodel metricset query without window secs",
			mutate: func(config map[string]interface{}) {
				config["query_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_METRICSET_QUERY", "aggregate": "AVG"},
				}
			},
			want: "query_config.window_secs is required when query_config.type is UMODEL_METRICSET_QUERY",
		},
		{
			name: "umodel metricset query without aggregate",
			mutate: func(config map[string]interface{}) {
				config["query_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_METRICSET_QUERY", "window_secs": 60},
				}
			},
			want: "query_config.aggregate is required when query_config.type is UMODEL_METRICSET_QUERY",
		},
		{
			name: "umodel logset condition without match operator",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_LOGSET_CONDITION", "duration_secs": 60, "severity": "CRITICAL"},
				}
			},
			want: "condition_config.match_operator is required when condition_config.type is UMODEL_LOGSET_CONDITION",
		},
		{
			name: "umodel logset condition without duration secs",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_LOGSET_CONDITION", "match_operator": "PRESENT", "severity": "CRITICAL"},
				}
			},
			want: "condition_config.duration_secs is required when condition_config.type is UMODEL_LOGSET_CONDITION",
		},
		{
			name: "umodel logset condition without severity",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_LOGSET_CONDITION", "match_operator": "PRESENT", "duration_secs": 60},
				}
			},
			want: "condition_config.severity is required when condition_config.type is UMODEL_LOGSET_CONDITION",
		},
		{
			name: "umodel logset condition without match field",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_LOGSET_CONDITION", "match_operator": "CONTAINS", "duration_secs": 60, "severity": "CRITICAL", "match_value": "ERROR"},
				}
			},
			want: "condition_config.match_field is required when condition_config.match_operator is CONTAINS, EQUALS or REGEX",
		},
		{
			name: "umodel logset condition without match value",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_LOGSET_CONDITION", "match_operator": "REGEX", "duration_secs": 60, "severity": "CRITICAL", "match_field": "level"},
				}
			},
			want: "condition_config.match_value is required when condition_config.match_operator is CONTAINS, EQUALS or REGEX",
		},
		{
			name: "umodel metricset condition no data to alert without severity",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_METRICSET_CONDITION", "no_data_policy": "NO_DATA_TO_ALERT"},
				}
			},
			want: "condition_config.no_data_alert_severity is required when condition_config.no_data_policy is NO_DATA_TO_ALERT",
		},
		{
			name: "umodel metricset condition severity forbidden for other policy",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{
					map[string]interface{}{"type": "UMODEL_METRICSET_CONDITION", "no_data_policy": "NO_DATA_TO_OK", "no_data_alert_severity": "WARNING"},
				}
			},
			want: "condition_config.no_data_alert_severity must not be set when condition_config.no_data_policy is not NO_DATA_TO_ALERT",
		},
		{
			name: "apm query without service id list",
			mutate: func(config map[string]interface{}) {
				config["query_config"] = []interface{}{map[string]interface{}{"type": "APM_MULTI_QUERY"}}
			},
			want: "query_config.service_id_list is required when query_config.type is APM_MULTI_QUERY",
		},
		{
			name: "apm query without measure list",
			mutate: func(config map[string]interface{}) {
				config["query_config"] = []interface{}{
					map[string]interface{}{
						"type":            "APM_MULTI_QUERY",
						"service_id_list": []interface{}{"tf-testacc-service"},
					},
				}
			},
			want: "query_config.measure_list is required when query_config.type is APM_MULTI_QUERY",
		},
		{
			name: "sls multi condition without triggers",
			mutate: func(config map[string]interface{}) {
				config["condition_config"] = []interface{}{map[string]interface{}{"type": "SLS_MULTI_CONDITION"}}
			},
			want: "condition_config.triggers requires at least one trigger when condition_config.type is SLS_MULTI_CONDITION",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			config := validConfig()
			testCase.mutate(config)
			_, err := resourceSchema.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil)
			if err == nil {
				t.Fatalf("expected the customize diff to reject the config with %q", testCase.want)
			}
			if !strings.Contains(err.Error(), testCase.want) {
				t.Fatalf("expected an error containing %q, got %v", testCase.want, err)
			}
		})
	}

	t.Run("legal sls combination", func(t *testing.T) {
		if _, err := resourceSchema.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(validConfig()), nil); err != nil {
			t.Fatalf("the legal SLS configuration must pass the customize diff: %v", err)
		}
	})

	t.Run("legal umodel logset combination", func(t *testing.T) {
		config := validConfig()
		config["datasource_config"] = []interface{}{map[string]interface{}{"type": "UMODEL"}}
		config["query_config"] = []interface{}{
			map[string]interface{}{"type": "UMODEL_LOGSET_QUERY", "log_set": "tf-testacc-logset", "expr": "cnt > 0", "window_secs": 60},
		}
		config["condition_config"] = []interface{}{
			map[string]interface{}{"type": "UMODEL_LOGSET_CONDITION", "match_operator": "PRESENT", "duration_secs": 60, "severity": "CRITICAL"},
		}
		if _, err := resourceSchema.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil); err != nil {
			t.Fatalf("the legal UMODEL logset configuration must pass the customize diff: %v", err)
		}
	})

	t.Run("legal umodel logset match combination", func(t *testing.T) {
		config := validConfig()
		config["datasource_config"] = []interface{}{map[string]interface{}{"type": "UMODEL"}}
		config["query_config"] = []interface{}{
			map[string]interface{}{"type": "UMODEL_LOGSET_QUERY", "log_set": "tf-testacc-logset", "expr": "cnt > 0", "window_secs": 60},
		}
		config["condition_config"] = []interface{}{
			map[string]interface{}{"type": "UMODEL_LOGSET_CONDITION", "match_operator": "EQUALS", "match_field": "level", "match_value": "ERROR", "duration_secs": 60, "severity": "CRITICAL"},
		}
		if _, err := resourceSchema.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil); err != nil {
			t.Fatalf("the legal UMODEL logset match configuration must pass the customize diff: %v", err)
		}
	})

	t.Run("legal umodel metricset no data combination", func(t *testing.T) {
		config := validConfig()
		config["datasource_config"] = []interface{}{map[string]interface{}{"type": "UMODEL"}}
		config["query_config"] = []interface{}{
			map[string]interface{}{"type": "UMODEL_METRICSET_QUERY", "window_secs": 60, "aggregate": "AVG"},
		}
		config["condition_config"] = []interface{}{
			map[string]interface{}{"type": "UMODEL_METRICSET_CONDITION", "operator": "GT", "threshold": 90, "duration_secs": 60, "severity": "CRITICAL", "no_data_policy": "NO_DATA_TO_ALERT", "no_data_alert_severity": "WARNING"},
		}
		if _, err := resourceSchema.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil); err != nil {
			t.Fatalf("the legal UMODEL metricset no-data configuration must pass the customize diff: %v", err)
		}
	})
}

func TestUnitCmsAlertRuleV2ParseQueryRelationsErrorPaths(t *testing.T) {
	cases := []struct {
		name string
		raw  map[string]interface{}
		want []string
	}{
		{
			name: "queries element is not an object",
			raw: map[string]interface{}{
				"queries": []interface{}{"not-an-object"},
			},
			want: []string{"reading alicloud_cms_alert_rule_v2 (rule-uuid-queries)", "query_config.0.queries.0 is not an object"},
		},
		{
			name: "joinings element is not an object",
			raw: map[string]interface{}{
				"joinings": []interface{}{
					map[string]interface{}{"type": "INNER"},
					"not-an-object",
				},
			},
			want: []string{"reading alicloud_cms_alert_rule_v2 (rule-uuid-joinings)", "query_config.0.joinings.1 is not an object"},
		},
		{
			name: "joinings conditions element is not an object",
			raw: map[string]interface{}{
				"joinings": []interface{}{
					map[string]interface{}{
						"type": "INNER",
						"conditions": []interface{}{
							map[string]interface{}{"lhsField": "$0.host", "operator": "==", "rhsField": "$1.host"},
							"not-an-object",
						},
					},
				},
			},
			want: []string{"reading alicloud_cms_alert_rule_v2 (rule-uuid-conditions)", "query_config.0.joinings.0.conditions.1 is not an object"},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			resourceId := ""
			for _, part := range testCase.want {
				if strings.HasPrefix(part, "reading alicloud_cms_alert_rule_v2 (") {
					resourceId = strings.TrimSuffix(strings.TrimPrefix(part, "reading alicloud_cms_alert_rule_v2 ("), ")")
				}
			}
			queries, joinings, err := cloudmonitorservice.CmsAlertRuleV2ParseQueryRelations(resourceId, testCase.raw)
			if err == nil {
				t.Fatalf("expected an error instead of silently dropping the malformed element, got queries=%v joinings=%v", queries, joinings)
			}
			for _, part := range testCase.want {
				if !strings.Contains(err.Error(), part) {
					t.Fatalf("expected the error to contain %q, got %v", part, err)
				}
			}
		})
	}
}

func TestUnitCmsAlertRuleV2ParseQueryRelationsKeepsValidElements(t *testing.T) {
	raw := map[string]interface{}{
		"queries": []interface{}{
			map[string]interface{}{"expr": "cnt > 0", "timeUnit": "minute", "start": int64(1), "end": int64(60)},
		},
		"joinings": []interface{}{
			map[string]interface{}{
				"type": "INNER",
				"conditions": []interface{}{
					map[string]interface{}{"lhsField": "$0.host", "operator": "==", "rhsField": "$1.host"},
				},
			},
		},
	}
	queries, joinings, err := cloudmonitorservice.CmsAlertRuleV2ParseQueryRelations("rule-uuid-valid", raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(queries) != 1 || queries[0]["expr"] != "cnt > 0" || queries[0]["time_unit"] != "minute" {
		t.Fatalf("unexpected queries parse result: %v", queries)
	}
	if len(joinings) != 1 || joinings[0]["type"] != "INNER" {
		t.Fatalf("unexpected joinings parse result: %v", joinings)
	}
	conditions, ok := joinings[0]["conditions"].([]map[string]interface{})
	if !ok || len(conditions) != 1 || conditions[0]["lhs_field"] != "$0.host" {
		t.Fatalf("unexpected joinings conditions parse result: %v", joinings[0]["conditions"])
	}
}
