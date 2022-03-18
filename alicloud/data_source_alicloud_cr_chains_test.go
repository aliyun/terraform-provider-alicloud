package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCrChainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CRSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_chain.default.chain_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cr_chain.default.chain_id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_chain.default.chain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cr_chain.default.chain_name}_fake"`,
		}),
	}
	repoNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"repo_name":           `"${alicloud_cr_ee_repo.default.name}"`,
			"repo_namespace_name": `"${alicloud_cr_ee_namespace.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"repo_name":           `"${alicloud_cr_ee_repo.default.name}"`,
			"repo_namespace_name": `"${alicloud_cr_ee_namespace.default.name}"`,
			"name_regex":          `"${alicloud_cr_chain.default.chain_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_cr_chain.default.chain_id}"]`,
			"repo_name":           `"${alicloud_cr_ee_repo.default.name}"`,
			"repo_namespace_name": `"${alicloud_cr_ee_namespace.default.name}"`,
			"name_regex":          `"${alicloud_cr_chain.default.chain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCrChainsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_cr_chain.default.chain_id}_fake"]`,
			"repo_name":           `"${alicloud_cr_ee_repo.default.name}"`,
			"repo_namespace_name": `"${alicloud_cr_ee_namespace.default.name}"`,
			"name_regex":          `"${alicloud_cr_chain.default.chain_name}_fake"`,
		}),
	}
	var existAlicloudCrChainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                         "1",
			"names.#":                                       "1",
			"chains.#":                                      "1",
			"chains.0.id":                                   CHECKSET,
			"chains.0.chain_id":                             CHECKSET,
			"chains.0.chain_name":                           CHECKSET,
			"chains.0.create_time":                          CHECKSET,
			"chains.0.description":                          CHECKSET,
			"chains.0.modified_time":                        CHECKSET,
			"chains.0.instance_id":                          CHECKSET,
			"chains.0.scope_id":                             CHECKSET,
			"chains.0.scope_type":                           CHECKSET,
			"chains.0.chain_config.#":                       "1",
			"chains.0.chain_config.0.nodes.#":               CHECKSET,
			"chains.0.chain_config.0.nodes.0.enable":        CHECKSET,
			"chains.0.chain_config.0.nodes.0.node_name":     CHECKSET,
			"chains.0.chain_config.0.nodes.1.node_config.#": "1",
			"chains.0.chain_config.0.nodes.1.node_config.0.deny_policy.#":             "1",
			"chains.0.chain_config.0.nodes.1.node_config.0.deny_policy.0.issue_count": "1",
			"chains.0.chain_config.0.nodes.1.node_config.0.deny_policy.0.issue_level": "MEDIUM",
			"chains.0.chain_config.0.nodes.1.node_config.0.deny_policy.0.logic":       "AND",
			"chains.0.chain_config.0.nodes.1.node_config.0.deny_policy.0.action":      "BLOCK_DELETE_TAG",
			"chains.0.chain_config.0.routers.#":                                       "6",
			"chains.0.chain_config.0.routers.0.from.#":                                "1",
			"chains.0.chain_config.0.routers.0.from.0.node_name":                      CHECKSET,
			"chains.0.chain_config.0.routers.0.to.#":                                  "1",
			"chains.0.chain_config.0.routers.0.to.0.node_name":                        CHECKSET,
		}
	}
	var fakeAlicloudCrChainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCrChainsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cr_chains.default",
		existMapFunc: existAlicloudCrChainsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCrChainsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCrChainsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, repoNameRegexConf, allConf)
}
func testAccCheckAlicloudCrChainsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tftestacc%d"
}

data "alicloud_cr_ee_instances" "default" {
  name_regex = "tf-testacc"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = data.alicloud_cr_ee_instances.default.ids[0]
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_namespace.default.instance_id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = var.name
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}

resource "alicloud_cr_chain" "default" {
  chain_name          = var.name
  description         = "description"
  instance_id         = alicloud_cr_ee_namespace.default.instance_id
  repo_name           = alicloud_cr_ee_repo.default.name
  repo_namespace_name = alicloud_cr_ee_namespace.default.name
  chain_config {
    routers {
      from {
        node_name = "DOCKER_IMAGE_BUILD"
      }
      to {
        node_name = "DOCKER_IMAGE_PUSH"
      }
    }
    routers {
      from {
        node_name = "DOCKER_IMAGE_PUSH"
      }
      to {
        node_name = "VULNERABILITY_SCANNING"
      }
    }
    routers {
      from {
        node_name = "VULNERABILITY_SCANNING"
      }
      to {
        node_name = "ACTIVATE_REPLICATION"
      }
    }
    routers {
      from {
        node_name = "ACTIVATE_REPLICATION"
      }
      to {
        node_name = "TRIGGER"
      }
    }
    routers {
      from {
        node_name = "VULNERABILITY_SCANNING"
      }
      to {
        node_name = "SNAPSHOT"
      }
    }
    routers {
      from {
        node_name = "SNAPSHOT"
      }
      to {
        node_name = "TRIGGER_SNAPSHOT"
      }
    }

    nodes {
      enable    = true
      node_name = "DOCKER_IMAGE_BUILD"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = true
      node_name = "DOCKER_IMAGE_PUSH"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = true
      node_name = "VULNERABILITY_SCANNING"
      node_config {
        deny_policy {
          issue_level = "MEDIUM"
          issue_count = 1
          action      = "BLOCK_DELETE_TAG"
          logic       = "AND"
        }
      }
    }
    nodes {
      enable    = true
      node_name = "ACTIVATE_REPLICATION"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = true
      node_name = "TRIGGER"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = false
      node_name = "SNAPSHOT"
      node_config {
        deny_policy {}
      }
    }
    nodes {
      enable    = false
      node_name = "TRIGGER_SNAPSHOT"
      node_config {
        deny_policy {}
      }
    }
  }
}

data "alicloud_cr_chains" "default" {	
	instance_id = "${alicloud_cr_chain.default.instance_id}"
	enable_details = true
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
