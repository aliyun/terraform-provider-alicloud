package maxcompute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMaxComputeProjectDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_project.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_project.default.id}_fake"]`,
		}),
	}

	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_project.default.id}"]`,
			"name_regex": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_project.default.id}_fake"]`,
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	MaxComputeProjectCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameConf, allConf)
}

var existMaxComputeProjectMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#":                       "1",
		"projects.0.comment":               acctest.CHECKSET,
		"projects.0.default_quota":         acctest.CHECKSET,
		"projects.0.owner":                 acctest.CHECKSET,
		"projects.0.project_name":          acctest.CHECKSET,
		"projects.0.properties.#":          "1",
		"projects.0.security_properties.#": "1",
		"projects.0.status":                acctest.CHECKSET,
		"projects.0.type":                  acctest.CHECKSET,
	}
}

var fakeMaxComputeProjectMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#": "0",
	}
}

var MaxComputeProjectCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_maxcompute_projects.default",
	ExistMapFunc: existMaxComputeProjectMapFunc,
	FakeMapFunc:  fakeMaxComputeProjectMapFunc,
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
