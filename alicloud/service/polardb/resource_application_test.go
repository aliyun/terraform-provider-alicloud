// Package polardb_test holds the SDKv2 acceptance tests for alicloud_polardb_application,
// migrated from the flat alicloud package into the polardb product package.
//
// The flat-layout tests leaned on the shared acceptance framework
// (resourceAttrInit/resourceAttrMapUpdateSet/resourceTestAccConfigFunc) which
// is not importable from a product package, so each step is expressed as an
// explicit, complete HCL config (matching the migrated-subpackage convention
// used by service/ecs). The step semantics are preserved: each step's config
// is the create base plus that step's accumulated overrides, so no Required
// field is dropped and no spurious parameter/whitelist removal diffs arise.
package polardb_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass components -- pre-existing order-insensitive TypeList (stock debt predating this migration PR).
// collection-order: bypass parameters -- pre-existing order-insensitive TypeList (stock debt predating this migration PR).
// collection-order: bypass security_ip_list -- pre-existing order-insensitive TypeList (stock debt predating this migration PR).

// TestAccAliCloudPolarDBApplication_Create exercises create, version upgrade,
// parameter modification and whitelist append/delete on a PostPaid PolarDB AI
// application.
func TestAccAliCloudPolarDBApplication_Create(t *testing.T) {
	resourceId := "alicloud_polardb_application.default"
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			// The flat suite pinned this resource to cn-beijing via
			// testAccPreCheckWithRegions; keep the pin here.
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBApplicationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPolardbApplicationConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigUpgradeVersion(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "upgrade_version", "true"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigParameters(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigWhitelist(name, "Append"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "security_ip_array_name", "test"),
					resource.TestCheckResourceAttr(resourceId, "security_ip_list.#", "1"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigWhitelist(name, "Delete"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "security_ip_array_name", "test"),
					resource.TestCheckResourceAttr(resourceId, "security_ip_list.#", "1"),
				),
			},
			{
				// item 65: tags mutation on a live application via
				// TagResources (ResourceType="application"), read back through
				// ListTagResources. Steps from the prior config are preserved so
				// only tags change, exercising the tag update path.
				Config: testAccPolardbApplicationConfigTags(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "allin"),
					resource.TestCheckResourceAttr(resourceId, "tags.Env", "test"),
				),
			},
		},
	})
}

// TestAccAliCloudPolarDBApplication_CreateFull covers the full attribute set
// including prepaid billing fields, model credentials and an attached cluster.
func TestAccAliCloudPolarDBApplication_CreateFull(t *testing.T) {
	resourceId := "alicloud_polardb_application.default"
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-beijing")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBApplicationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPolardbApplicationConfigFull(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigFullWhitelist(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
				),
			},
		},
	})
}

func testAccCheckPolarDBApplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_application" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := polardb.NewService(client).DescribeApplicationAttribute(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing polardb application %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("polardb application %s still exists", rs.Primary.ID)
	}
	return nil
}

// testAccPolardbApplicationDependence is the shared VPC/VSwitch + node-class
// lookup the application depends on; ported verbatim from the flat-layout
// resourcePolarDBApplicationConfigDependence.
func testAccPolardbApplicationDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_vpc" "default" {
	vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
	zone_id  = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	vpc_id   = alicloud_vpc.default.id
	cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

locals {
	vpc_id     = alicloud_vpc.default.id
	vswitch_id = alicloud_vswitch.default.id
}

data "alicloud_polardb_node_classes" "this" {
	db_type    = "MySQL"
	db_version = "8.0"
	pay_type   = "PostPaid"
	category   = "Normal"
}
`, name)
}

func testAccPolardbApplicationConfigBasic(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
}
`, testAccPolardbApplicationDependence(name))
}

func testAccPolardbApplicationConfigUpgradeVersion(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version = true
}
`, testAccPolardbApplicationDependence(name))
}

func testAccPolardbApplicationConfigParameters(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version = true
	parameters {
		parameter_name  = "secret.dashscope.apiKey"
		parameter_value = "abc"
	}
}
`, testAccPolardbApplicationDependence(name))
}

func testAccPolardbApplicationConfigWhitelist(name, modifyMode string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version        = true
	parameters {
		parameter_name  = "secret.dashscope.apiKey"
		parameter_value = "abc"
	}
	modify_mode            = "%[2]s"
	security_ip_array_name = "test"
	security_ip_list       = ["127.0.0.1"]
}
`, testAccPolardbApplicationDependence(name), modifyMode)
}

// testAccPolardbApplicationFullDependence is the full dependence (adds a
// resource-manager group lookup and a PolarDB cluster the application
// attaches to); ported verbatim from resourcePolarDBApplicationFullConfigDependence.
func testAccPolardbApplicationFullDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_vpc" "default" {
	vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
	zone_id  = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	vpc_id   = alicloud_vpc.default.id
	cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

locals {
	vpc_id     = alicloud_vpc.default.id
	vswitch_id = alicloud_vswitch.default.id
	component_id = alicloud_polardb_application.default.component_id
}

data "alicloud_polardb_node_classes" "this" {
	db_type    = "MySQL"
	db_version = "8.0"
	pay_type   = "PostPaid"
	category   = "Normal"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_polardb_cluster" "cluster" {
	db_type = "MySQL"
	db_version = "8.0"
	pay_type = "PostPaid"
	db_node_count = "2"
	db_node_class = "polar.mysql.x4.large"
	vswitch_id = "${local.vswitch_id}"
	description = "${var.name}"

	// encrypt_new_tables is optional without a computed fallback in the stock
	// cluster resource, and DescribeDBClusters intermittently stops returning
	// the TDEStatus sub-object after the first read; the resulting
	// "OFF -> null" plan fails this fixture's second step. Other polardb
	// fixtures ignore the same field for the same reason. It does not belong
	// to alicloud_polardb_application behaviour.
	lifecycle {
		ignore_changes = [encrypt_new_tables]
	}
}
`, name)
}

func testAccPolardbApplicationConfigFull(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	auto_renew        = false
	period            = 1
	model_api_key     = "sk-xxx"
	model_api         = "openai-completions"
	ai_db_cluster_id  = alicloud_polardb_cluster.cluster.id
	db_cluster_id     = alicloud_polardb_cluster.cluster.id
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
`, testAccPolardbApplicationFullDependence(name))
}

func testAccPolardbApplicationConfigFullWhitelist(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	auto_renew        = false
	period            = 1
	model_api_key     = "sk-xxx"
	model_api         = "openai-completions"
	ai_db_cluster_id  = alicloud_polardb_cluster.cluster.id
	db_cluster_id     = alicloud_polardb_cluster.cluster.id
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
	component_id           = ""
	modify_mode            = "Append"
	security_ip_array_name = "test"
	security_ip_list       = ["127.0.0.1"]
}
`, testAccPolardbApplicationFullDependence(name))
}

// testAccPolardbApplicationConfigTags extends the whitelist-delete config
// with two tags, exercising the TagResources update path (item 65) and the
// ListTagResources read-back. It shares the basic dependence so only the
// tags attribute differs from the prior step.
func testAccPolardbApplicationConfigTags(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version        = true
	parameters {
		parameter_name  = "secret.dashscope.apiKey"
		parameter_value = "abc"
	}
	modify_mode            = "Delete"
	security_ip_array_name = "test"
	security_ip_list       = ["127.0.0.1"]
	tags = {
		Created = "allin"
		Env     = "test"
	}
}
`, testAccPolardbApplicationDependence(name))
}
