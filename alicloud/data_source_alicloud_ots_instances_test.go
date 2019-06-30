package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudOtsInstancesDataSource_ids(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_ids_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.id", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.status", string(Running)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.write_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.read_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.cluster_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.user_id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.network", "NORMAL"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.description", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.entity_quota"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.0", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.0", fmt.Sprintf("tf-testAcc%d", randInt)),
				),
			},
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_ids_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsInstancesDataSource_name_regex(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_name_regex_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.id", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.status", string(Running)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.write_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.read_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.cluster_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.user_id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.network", "NORMAL"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.description", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.entity_quota"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.0", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.0", fmt.Sprintf("tf-testAcc%d", randInt)),
				),
			},
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_name_regex_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsInstancesDataSource_tags(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_tags_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.id", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.status", string(Running)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.write_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.read_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.cluster_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.user_id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.network", "NORMAL"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.description", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.entity_quota"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.0", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.0", fmt.Sprintf("tf-testAcc%d", randInt)),
				),
			},
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_tags_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsInstancesDataSource_All(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_all_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.id", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.status", string(Running)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.write_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.read_capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.cluster_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.user_id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.network", "NORMAL"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.description", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instances.instances", "instances.0.entity_quota"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.0", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.0", fmt.Sprintf("tf-testAcc%d", randInt)),
				),
			},
			{
				Config: testAccCheckAlicloudOtsInstancesDataSource_all_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instances.instances"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instances.instances", "names.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudOtsInstancesDataSource_ids_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
      ids = [ "${alicloud_ots_instance.foo.id}" ]
	}
	`, randInt)
}

func testAccCheckAlicloudOtsInstancesDataSource_ids_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
      ids = [ "${alicloud_ots_instance.foo.id}-fake" ]
	}
	`, randInt)
}

func testAccCheckAlicloudOtsInstancesDataSource_name_regex_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
      name_regex = "${alicloud_ots_instance.foo.name}"	
	}
	`, randInt)
}

func testAccCheckAlicloudOtsInstancesDataSource_name_regex_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
      name_regex = "${alicloud_ots_instance.foo.name}-fake"	
	}
	`, randInt)
}

func testAccCheckAlicloudOtsInstancesDataSource_tags_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF%d"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
      tags = "${alicloud_ots_instance.foo.tags}"
	}
	`, randInt, randInt)
}

func testAccCheckAlicloudOtsInstancesDataSource_tags_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF%d"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
      tags = {
        Created = "${alicloud_ots_instance.foo.tags.Created}-fake"
        For = "${alicloud_ots_instance.foo.tags.For}-fake"
      }
	}
	`, randInt, randInt)
}

func testAccCheckAlicloudOtsInstancesDataSource_all_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF%d"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
	  ids = [ "${alicloud_ots_instance.foo.id}" ]
      name_regex = "${alicloud_ots_instance.foo.name}"
      tags = "${alicloud_ots_instance.foo.tags}"
	}
	`, randInt, randInt)
}

func testAccCheckAlicloudOtsInstancesDataSource_all_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF%d"
		For = "acceptance test"
	  }
	}
	data "alicloud_ots_instances" "instances" {
	  ids = [ "${alicloud_ots_instance.foo.id}-fake" ]
      name_regex = "${alicloud_ots_instance.foo.name}-fake"	
      tags = {
        Created = "${alicloud_ots_instance.foo.tags.Created}-fake"
        For = "${alicloud_ots_instance.foo.tags.For}-fake"
      }
	}
	`, randInt, randInt)
}
