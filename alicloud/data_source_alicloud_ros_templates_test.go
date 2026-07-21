package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudROSTemplatesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_templates.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosTemplates%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosTemplatesDependence)

	sharTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template.default.id}"},
			"share_type":     "Private",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template.default.id}"},
			"share_type":     "Shared",
			"enable_details": "true",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_template.default.template_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "{alicloud_ros_template.default.template_name}-fake",
			"enable_details": "true",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ROS",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF-fake",
				"For":     "update test fake",
			},
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_template.default.template_name}",
			"ids":            []string{"${alicloud_ros_template.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_template.default.template_name}-fake",
			"ids":            []string{"${alicloud_ros_template.default.id}"},
			"enable_details": "true",
		}),
	}
	var existRosTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"names.#":                      "1",
			"names.0":                      name,
			"templates.#":                  "1",
			"templates.0.change_set_id":    "",
			"templates.0.description":      "test for ros templates",
			"templates.0.share_type":       CHECKSET,
			"templates.0.stack_group_name": "",
			"templates.0.stack_id":         "",
			"templates.0.tags.%":           "2",
			"templates.0.template_body":    CHECKSET,
			"templates.0.id":               CHECKSET,
			"templates.0.template_id":      CHECKSET,
			"templates.0.template_name":    name,
			"templates.0.template_version": CHECKSET,
		}
	}

	var fakeRosTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"templates.#": "0",
		}
	}

	var rosTemplatesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosTemplatesMapFunc,
		fakeMapFunc:  fakeRosTemplatesMapFunc,
	}

	rosTemplatesInfo.dataSourceTestCheck(t, 0, sharTypeConf, nameRegexConf, tagsConf, idsConf, allConf)
}

func dataSourceRosTemplatesDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ros_template" "default" {
	  template_name=   "%s"
	  template_body=   "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}"
	   tags = {
   			"Created": "TF",
			"For":     "ROS",
 		}
      description= "test for ros templates"
	}
	`, name)
}
