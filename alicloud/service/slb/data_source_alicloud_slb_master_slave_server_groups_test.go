package slb_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAlicloudSLBMasterSlaveServerGroupsDataSource_basic(t *testing.T) {
	basicConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_load_balancer.default.id}"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_load_balancer.default.id}"`,
			"name_regex":       `"${alicloud_slb_master_slave_server_group.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_load_balancer.default.id}"`,
			"name_regex":       `"${alicloud_slb_master_slave_server_group.default.name}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_load_balancer.default.id}"`,
			"ids":              `["${alicloud_slb_master_slave_server_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_load_balancer.default.id}"`,
			"ids":              `["${alicloud_slb_master_slave_server_group.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_load_balancer.default.id}"`,
			"ids":              `["${alicloud_slb_master_slave_server_group.default.id}"]`,
			"name_regex":       `"${alicloud_slb_master_slave_server_group.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_load_balancer.default.id}"`,
			"ids":              `["${alicloud_slb_master_slave_server_group.default.id}_fake"]`,
			"name_regex":       `"${alicloud_slb_master_slave_server_group.default.name}"`,
		}),
	}

	var existSLBMasterSlaveServerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                       "1",
			"ids.#":                          "1",
			"names.#":                        "1",
			"groups.0.id":                    CHECKSET,
			"groups.0.name":                  "tf-testAccslbmasterslaveservergroupsdatasourcebasic",
			"groups.0.servers.#":             "2",
			"groups.0.servers.0.weight":      "100",
			"groups.0.servers.0.instance_id": CHECKSET,
			"groups.0.servers.1.instance_id": CHECKSET,
		}
	}

	var fakeSLBMasterSlaveServerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	const dsResourceId = "data.alicloud_slb_master_slave_server_groups.default"
	existCheck := acctest.ResourceAttrInit(dsResourceId, existSLBMasterSlaveServerGroupsMapFunc(-1)).ResourceAttrMapUpdateSet()
	fakeCheck := acctest.ResourceAttrInit(dsResourceId, fakeSLBMasterSlaveServerGroupsMapFunc(-1)).ResourceAttrMapUpdateSet()

	existStep := func(config string) resource.TestStep {
		return resource.TestStep{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				existCheck(nil),
				testAccCheckMasterSlaveServerGroupDescriptions(),
			),
		}
	}
	fakeStep := func(config string) resource.TestStep {
		return resource.TestStep{
			Config: config,
			Check:  resource.ComposeTestCheckFunc(fakeCheck(nil)),
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			existStep(allConf.ExistConfig),
			fakeStep(allConf.FakeConfig),
			existStep(basicConf.ExistConfig),
			existStep(nameRegexConf.ExistConfig),
			fakeStep(nameRegexConf.FakeConfig),
			existStep(idsConf.ExistConfig),
			fakeStep(idsConf.FakeConfig),
		},
	})
}

func testAccCheckMasterSlaveServerGroupDescriptions() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		group, ok := s.RootModule().Resources["alicloud_slb_master_slave_server_group.default"]
		if !ok {
			return fmt.Errorf("master slave server group not found in state")
		}
		ds, ok := s.RootModule().Resources["data.alicloud_slb_master_slave_server_groups.default"]
		if !ok {
			return fmt.Errorf("master slave server groups data source not found in state")
		}

		type expectedServer struct {
			serverType  string
			description string
		}
		expectedDescriptions := map[string]string{
			"Master": "ds-master-desc",
			"Slave":  "ds-slave-desc",
		}
		expected := make(map[string]expectedServer, 2)
		serverKeyPattern := regexp.MustCompile(`^servers\.([^.]+)\.server_id$`)
		for key, serverId := range group.Primary.Attributes {
			match := serverKeyPattern.FindStringSubmatch(key)
			if match == nil {
				continue
			}
			prefix := "servers." + match[1]
			serverType := group.Primary.Attributes[prefix+".server_type"]
			want, ok := expectedDescriptions[serverType]
			if !ok {
				return fmt.Errorf("unexpected configured server_type %q", serverType)
			}
			expected[serverId] = expectedServer{serverType: serverType, description: want}
		}
		if len(expected) != 2 {
			return fmt.Errorf("expected 2 configured backend servers, got %d", len(expected))
		}

		for i := 0; i < 2; i++ {
			serverId := ds.Primary.Attributes[fmt.Sprintf("groups.0.servers.%d.instance_id", i)]
			gotType := ds.Primary.Attributes[fmt.Sprintf("groups.0.servers.%d.server_type", i)]
			gotDescription := ds.Primary.Attributes[fmt.Sprintf("groups.0.servers.%d.description", i)]
			expect, ok := expected[serverId]
			if !ok {
				return fmt.Errorf("data source returned unexpected backend server %q", serverId)
			}
			if gotType != expect.serverType || gotDescription != expect.description {
				return fmt.Errorf("data source backend server %s: got server_type=%q description=%q, want server_type=%q description=%q",
					serverId, gotType, gotDescription, expect.serverType, expect.description)
			}
			delete(expected, serverId)
		}
		if len(expected) != 0 {
			missingServerIds := make([]string, 0, len(expected))
			for serverId := range expected {
				missingServerIds = append(missingServerIds, serverId)
			}
			return fmt.Errorf("data source did not return expected backend servers: %v", missingServerIds)
		}
		return nil
	}
}

func testAccCheckAlicloudSlbMasterSlaveServerGroupsDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccslbmasterslaveservergroupsdatasourcebasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
}

data "alicloud_images" "default" {
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_regex = "^ubuntu"
  most_recent = true
  owners = "system"
}

data "alicloud_instance_types" "default" {
  cpu_core_count = 2
  memory_size = 4
  system_disk_category = "cloud_essd"
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  load_balancer_spec = "slb.s2.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id             = "${alicloud_slb_load_balancer.default.id}"
  master_slave_server_group_id = "${alicloud_slb_master_slave_server_group.default.id}"
  frontend_port                = "22"
  protocol                     = "tcp"
  bandwidth                    = "10"
  health_check_type            = "tcp"
  persistence_timeout          = 3600
  healthy_threshold            = 8
  unhealthy_threshold          = 8
  health_check_timeout         = 8
  health_check_interval        = 5
  health_check_http_code       = "http_2xx"
  health_check_connect_port    = 20
  health_check_uri             = "/console"
  established_timeout          = 600
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"

  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_essd"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"

  count = "2"
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type       = "PostPaid"
}

resource "alicloud_slb_master_slave_server_group" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alicloud_instance.default.0.id}"
      port = 80
      weight = 100
      server_type = "Master"
      description = "ds-master-desc"
  }
  servers {
      server_id = "${alicloud_instance.default.1.id}"
      port = 80
      weight = 100
      server_type = "Slave"
      description = "ds-slave-desc"
  }
}


data "alicloud_slb_master_slave_server_groups" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
