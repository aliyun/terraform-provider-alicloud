package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCSApplication_swarm(t *testing.T) {
	var basic, env cs.GetProjectResponse
	var swarm cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerApplicationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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

func TestAccAlicloudCSApplication_update(t *testing.T) {
	var basic cs.GetProjectResponse
	var swarm cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerApplicationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSApplication_updateBefore(testWebTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &swarm),
					testAccCheckContainerApplicationExists("alicloud_cs_application.basic", &basic),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.#", "1"),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.0.name", "web"),
				),
			},

			resource.TestStep{
				Config: testAccCSApplication_updateBlueGreen(testJavaTemplate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &swarm),
					testAccCheckContainerApplicationExists("alicloud_cs_application.basic", &basic),
					resource.TestCheckResourceAttr("alicloud_cs_application.basic", "services.#", "2"),
				),
			},

			resource.TestStep{
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
		client := testAccProvider.Meta().(*AliyunClient)
		app, err := client.DescribeContainerApplication(parts[0], parts[1])

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
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_application" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		app, err := client.DescribeContainerApplication(parts[0], parts[1])

		if err != nil {
			if NotFoundError(err) ||
				IsExceptedError(err, ApplicationNotFound) ||
				IsExceptedError(err, ApplicationErrorIgnore) ||
				IsExceptedError(err, AliyunGoClientFailure) {
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
	default = "testAccCSApplication-basic"
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

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
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
	PASSWORD = "Test12345"
  }
}
`, basic, env)
}

func testAccCSApplication_updateBefore(web string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "testAccCSApplication-update"
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

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
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
	default = "testAccCSApplication-update"
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

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
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
	default = "testAccCSApplication-update"
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

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
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
