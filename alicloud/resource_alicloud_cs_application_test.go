package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func SkipTestAccAlicloudCSApplication_swarm(t *testing.T) {
	var basic, env cs.GetProjectResponse
	var swarm cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSApplication_basic(testJavaTemplate, testMultiTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &swarm),
					testAccCheckContainerApplicationExists("alicloud_cs_application.basic", &basic),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.0.name", "slave-java"),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.0.status", "running"),
					testAccCheckContainerApplicationExists("alicloud_cs_application.env", &env),
					resource.TestCheckResourceAttr("alicloud_cs_application.env", "services.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_application.env", "services.1.name", "web"),
					resource.TestCheckResourceAttr("alicloud_cs_application.env", "environment.%", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_application.env", "environment.USER", "swarm"),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCSApplication_update(t *testing.T) {
	var basic cs.GetProjectResponse
	var swarm cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSApplication_updateBefore(testWebTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &swarm),
					testAccCheckContainerApplicationExists("alicloud_cs_application.basic", &basic),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.0.name", "web"),
				),
			},

			{
				Config: testAccCSApplication_updateBlueGreen(testJavaTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &swarm),
					testAccCheckContainerApplicationExists("alicloud_cs_application.basic", &basic),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.#", "2"),
				),
			},

			{
				Config: testAccCSApplication_updateConfirm(testJavaTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &swarm),
					testAccCheckContainerApplicationExists("alicloud_cs_application.basic", &basic),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.0.name", "slave-java"),
				),
			},
		},
	})
}

func testAccCheckContainerApplicationExists(n string, d *cs.GetProjectResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cluster, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found:%s", n)
		}

		if cluster.Primary.ID == "" {
			return fmt.Errorf("No Container cluster ID is set")
		}
		parts := strings.Split(cluster.Primary.ID, COLON_SEPARATED)
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		csService := CsService{client}
		app, err := csService.DescribeContainerApplication(parts[0], parts[1])

		if err != nil {
			return err
		}

		if app.Name == "" {
			return fmt.Errorf("Container application %s not found.", cluster.Primary.ID)
		}

		*d = app
		return nil
	}
}

func testAccCheckContainerApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	csService := CsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_application" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		app, err := csService.DescribeContainerApplication(parts[0], parts[1])

		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"Not Found", "Unable to reach primary cluster manager", AliyunGoClientFailure}) {
				continue
			}
			return err
		}

		if app.Name != "" {
			return fmt.Errorf("Error container application %s still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCSApplication_basic(basic, env string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCSApplication-basic"
}
data "alicloud_images" main {
  most_recent = true
  name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_hbase_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_hbase_zones.default.ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = local.vswitch_id
}

resource "alicloud_cs_application" "basic" {
  cluster_name = "${alicloud_cs_swarm.cs_vpc.name}"
  name = "${var.name}-app-basic"
  version = "1.0"
  description = "from tf creation"
  template = <<DEFINITION
  %s
  DEFINITION
  latest_image = "true"
}

resource "alicloud_cs_application" "env" {
  cluster_name = "${alicloud_cs_swarm.cs_vpc.name}"
  name = "${var.name}-app-env"
  version = "1.0"
  template = <<DEFINITION
  %s
  DEFINITION
  description = "from tf creation"
  latest_image = "true"
  environment = {
	USER = "swarm"
	PASSWORD = "Yourpassword1234"
  }
}
`, basic, env)
}

func testAccCSApplication_updateBefore(web string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCSApplication-basic"
}
data "alicloud_images" main {
  most_recent = true
  name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_hbase_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_hbase_zones.default.ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = local.vswitch_id
}

resource "alicloud_cs_application" "basic" {
  cluster_name = "${alicloud_cs_swarm.cs_vpc.name}"
  name = "${var.name}"
  version = "1.0"
  description = "from tf creation"
  template = <<DEFINITION
  %s
  DEFINITION
  latest_image = "true"
}
`, web)
}

func testAccCSApplication_updateBlueGreen(java string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCSApplication-basic"
}
data "alicloud_images" main {
	most_recent = true
	name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_hbase_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_hbase_zones.default.ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = local.vswitch_id
}

resource "alicloud_cs_application" "basic" {
  cluster_name = "${alicloud_cs_swarm.cs_vpc.name}"
  name = "${var.name}"
  version = "1.1"
  description = "from tf creation"
  template = <<DEFINITION
  %s
  DEFINITION
  latest_image = "true"
  blue_green = "true"
}
`, java)
}

func testAccCSApplication_updateConfirm(java string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCSApplication-basic"
}
data "alicloud_images" main {
	most_recent = true
	name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_hbase_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_hbase_zones.default.ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = local.vswitch_id
}

resource "alicloud_cs_application" "basic" {
  cluster_name = "${alicloud_cs_swarm.cs_vpc.name}"
  name = "${var.name}"
  version = "1.1"
  description = "from tf creation"
  template = <<DEFINITION
  %s
  DEFINITION
  latest_image = "true"
  blue_green = "true"
  blue_green_confirm = "true"
}

`, java)
}

var testJavaTemplate = `
slave-java:
  image: 'registry.aliyuncs.com/acs-sample/jenkins-slave-dind-java'
  volumes:
      - /var/run/docker.sock:/var/run/docker.sock
  restart: always
  labels:
      aliyun.scale: '1'
`

var testWebTemplate = `
web:
  image: registry.cn-beijing.aliyuncs.com/101datumx/web:v0.3.4
  ports:
    - 8080
`

var testMultiTemplate = `
slave-java:
  image: 'registry.aliyuncs.com/acs-sample/jenkins-slave-dind-java'
  volumes:
      - /var/run/docker.sock:/var/run/docker.sock
  restart: always
  labels:
      aliyun.scale: '1'
  environment:
      USER: "$${USER}"
      PASSWORD: "$${PASSWORD}"

web:
  image: registry.cn-beijing.aliyuncs.com/101datumx/web:v0.3.4
  ports:
    - 8080
`
