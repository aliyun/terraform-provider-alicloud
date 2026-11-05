package oos_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOOSExecutionsDataSource(t *testing.T) {
	resourceId := "data.alicloud_oos_executions.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosTemplate-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceOosExecutionsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_execution.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_execution.default.id}-fake"},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_oos_execution.default.id}"},
			"status": "Success",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_oos_execution.default.id}"},
			"status": "Cancelled",
		}),
	}
	templateNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alicloud_oos_execution.default.id}"},
			"template_name": "${alicloud_oos_template.default.template_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alicloud_oos_execution.default.id}"},
			"template_name": "${alicloud_oos_template.default.template_name}-fake",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alicloud_oos_execution.default.id}"},
			"status":        "Success",
			"template_name": "${alicloud_oos_template.default.template_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alicloud_oos_execution.default.id}"},
			"status":        "Cancelled",
			"template_name": "${alicloud_oos_template.default.template_name}",
		}),
	}
	var existOosExecutionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"ids.0":                            acctest.CHECKSET,
			"executions.#":                     "1",
			"executions.0.category":            "Other",
			"executions.0.counters":            acctest.CHECKSET,
			"executions.0.create_date":         acctest.CHECKSET,
			"executions.0.end_date":            acctest.CHECKSET,
			"executions.0.executed_by":         acctest.CHECKSET,
			"executions.0.id":                  acctest.CHECKSET,
			"executions.0.execution_id":        acctest.CHECKSET,
			"executions.0.is_parent":           "false",
			"executions.0.mode":                "Automatic",
			"executions.0.outputs":             acctest.CHECKSET,
			"executions.0.parameters":          acctest.CHECKSET,
			"executions.0.parent_execution_id": "",
			"executions.0.ram_role":            "",
			"executions.0.start_date":          acctest.CHECKSET,
			"executions.0.status":              "Success",
			"executions.0.status_message":      "",
			"executions.0.status_reason":       "",
			"executions.0.template_id":         acctest.CHECKSET,
			"executions.0.template_name":       name,
			"executions.0.template_version":    acctest.CHECKSET,
			"executions.0.update_date":         acctest.CHECKSET,
		}
	}

	var fakeOosExecutionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"executions.#": "0",
		}
	}

	var oosExecutionsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existOosExecutionMapFunc,
		FakeMapFunc:  fakeOosExecutionMapFunc,
	}

	oosExecutionsInfo.DataSourceTestCheck(t, 0, statusConf, idsConf, templateNameConf, allConf)
}

func dataSourceOosExecutionsDependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_oos_template" "default" {
		  content= <<EOF
		  {
			"FormatVersion": "OOS-2019-06-01",
			"Description": "Describe instances of given status",
			"Parameters":{
			  "Status":{
				"Type": "String",
				"Description": "(Required) The status of the Ecs instance."
			  }
			},
			"Tasks": [
			  {
				"Properties" :{
				  "Parameters":{
					"Status": "{{ Status }}"
				  },
				  "API": "DescribeInstances",
				  "Service": "Ecs"
				},
				"Name": "foo",
				"Action": "ACS::ExecuteApi"
			  }]
		  }
		  EOF
		  template_name = "%s"
		  version_name = "test"
		  tags = {
			"Created" = "TF",
			"For" = "template Test"
		  }
		}
		
		resource "alicloud_oos_execution" "default"{
			template_name = alicloud_oos_template.default.template_name
			description = "From TF Test"
			parameters = <<EOF
				{"Status":"Running"}
		  	EOF
		}
	`, name)
}
