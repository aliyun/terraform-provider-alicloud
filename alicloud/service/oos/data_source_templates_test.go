package oos_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOOSTemplatesDataSource(t *testing.T) {
	resourceId := "data.alicloud_oos_templates.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosTemplate-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceOosTemplatesDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_template.default.template_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_template.default.template_name}-fake",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}-fake"},
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "template Test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
			"tags": map[string]interface{}{
				"Created": "TF_fake",
				"For":     "template Test",
			},
		}),
	}
	shareTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_oos_template.default.template_name}"},
			"share_type": "Private",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_oos_template.default.template_name}"},
			"share_type": "Public",
		}),
	}
	hasTriggerConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_oos_template.default.template_name}"},
			"has_trigger": "false",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_oos_template.default.template_name}"},
			"has_trigger": "true",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "template Test",
			},
			"has_trigger": "false",
			"share_type":  "Private",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "template Test",
			},
			"has_trigger": "false",
			"share_type":  "Public",
		}),
	}
	var existOosTemplateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        acctest.CHECKSET,
			"templates.#":                  "1",
			"templates.0.category":         acctest.CHECKSET,
			"templates.0.created_date":     acctest.CHECKSET,
			"templates.0.description":      acctest.CHECKSET,
			"templates.0.has_trigger":      acctest.CHECKSET,
			"templates.0.created_by":       acctest.CHECKSET,
			"templates.0.share_type":       "Private",
			"templates.0.template_format":  "JSON",
			"templates.0.template_id":      acctest.CHECKSET,
			"templates.0.id":               name,
			"templates.0.template_name":    name,
			"templates.0.template_type":    "Automation",
			"templates.0.template_version": acctest.CHECKSET,
			"templates.0.updated_by":       acctest.CHECKSET,
			"templates.0.updated_date":     acctest.CHECKSET,
		}
	}

	var fakeOosTemplateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"templates.#": "0",
		}
	}

	var oosTemplatesInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existOosTemplateMapFunc,
		FakeMapFunc:  fakeOosTemplateMapFunc,
	}

	oosTemplatesInfo.DataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf, shareTypeConf, hasTriggerConf, allConf)
}

func dataSourceOosTemplatesDependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_oos_template" "default" {
		  content= <<EOF
		  {
			"FormatVersion": "OOS-2019-06-01",
			"Description": "Update Describe instances of given status",
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
	`, name)
}
