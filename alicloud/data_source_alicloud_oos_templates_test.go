package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOOSTemplatesDataSource(t *testing.T) {
	resourceId := "data.alicloud_oos_templates.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosTemplate-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOosTemplatesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_template.default.template_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_template.default.template_name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}-fake"},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "template Test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
			"tags": map[string]interface{}{
				"Created": "TF_fake",
				"For":     "template Test",
			},
		}),
	}
	shareTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_oos_template.default.template_name}"},
			"share_type": "Private",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_oos_template.default.template_name}"},
			"share_type": "Public",
		}),
	}
	hasTriggerConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_oos_template.default.template_name}"},
			"has_trigger": "false",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_oos_template.default.template_name}"},
			"has_trigger": "true",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_template.default.template_name}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "template Test",
			},
			"has_trigger": "false",
			"share_type":  "Private",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":                        CHECKSET,
			"templates.#":                  "1",
			"templates.0.category":         CHECKSET,
			"templates.0.created_date":     CHECKSET,
			"templates.0.description":      CHECKSET,
			"templates.0.has_trigger":      CHECKSET,
			"templates.0.created_by":       CHECKSET,
			"templates.0.share_type":       "Private",
			"templates.0.template_format":  "JSON",
			"templates.0.template_id":      CHECKSET,
			"templates.0.id":               name,
			"templates.0.template_name":    name,
			"templates.0.template_type":    "Automation",
			"templates.0.template_version": CHECKSET,
			"templates.0.updated_by":       CHECKSET,
			"templates.0.updated_date":     CHECKSET,
		}
	}

	var fakeOosTemplateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"templates.#": "0",
		}
	}

	var oosTemplatesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOosTemplateMapFunc,
		fakeMapFunc:  fakeOosTemplateMapFunc,
	}

	oosTemplatesInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, tagsConf, shareTypeConf, hasTriggerConf, allConf)
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
