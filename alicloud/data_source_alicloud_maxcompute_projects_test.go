package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMaxComputeProjectDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_project.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_project.default.id}_fake"]`,
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_project.default.id}"]`,
			"name_regex": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_project.default.id}_fake"]`,
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	MaxComputeProjectCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, allConf)
}

var existMaxComputeProjectMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#":                       "1",
		"projects.0.comment":               CHECKSET,
		"projects.0.default_quota":         CHECKSET,
		"projects.0.owner":                 CHECKSET,
		"projects.0.project_name":          CHECKSET,
		"projects.0.properties.#":          "1",
		"projects.0.security_properties.#": "1",
		"projects.0.status":                CHECKSET,
		"projects.0.type":                  CHECKSET,
	}
}

var fakeMaxComputeProjectMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#": "0",
	}
}

var MaxComputeProjectCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_maxcompute_projects.default",
	existMapFunc: existMaxComputeProjectMapFunc,
	fakeMapFunc:  fakeMaxComputeProjectMapFunc,
}

func testAccCheckAlicloudMaxComputeProjectSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf_testaccmaxcp%d"
}

resource "alicloud_maxcompute_project" "default" {
  status = "AVAILABLE"
  ip_white_list {
    ip_list     = "10.0.0.0/8"
    vpc_ip_list = "10.0.0.0/8"
  }

  security_properties {
    project_protection {
      protected        = "true"
      exception_policy = "{\"Version\":\"1\",\"Statement\":[{\"Action\":[\"odps:*\"],\"Resource\":[\"acs:odps:*:projects/ludong/tables/*\"],\"Effect\":\"Allow\",\"Principal\":[\"ALIYUN$ludong@aliyun.com\"]}]}"
    }

    using_acl                            = "false"
    using_policy                         = "false"
    object_creator_has_access_permission = "false"
    object_creator_has_grant_permission  = "false"
    label_security                       = "false"
    enable_download_privilege            = "false"
  }

  tags = {
    For     = "Test"
    Created = "TF-CI"
  }
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = "terraform测试项目"
  properties {
    type_system      = "2"
    sql_metering_max = "10240"
    encryption {
      key       = "f58d854d-7bc0-4a6e-9205-160e10ffedec"
      enable    = "true"
      algorithm = "AESCTR"
    }

  }

}

data "alicloud_maxcompute_projects" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
