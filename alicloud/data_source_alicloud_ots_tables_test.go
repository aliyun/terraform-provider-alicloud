package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudOtsTablesDataSource_basic(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsTablesDataSource_basic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_tables.tables"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.id", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.instance_name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.table_name", fmt.Sprintf("testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.name", "pk2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.type", "String"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.time_to_live", "-1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.max_version", "1"),

					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.0", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.0", fmt.Sprintf("testAcc%d", randInt)),
				),
			},
		},
	})
}

func TestAccAlicloudOtsTablesDataSource_ids(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsTablesDataSource_ids_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_tables.tables"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.id", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.instance_name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.table_name", fmt.Sprintf("testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.name", "pk2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.type", "String"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.time_to_live", "-1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.max_version", "1"),

					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.0", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.0", fmt.Sprintf("testAcc%d", randInt)),
				),
			},
			{
				Config: testAccCheckAlicloudOtsTablesDataSource_ids_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_tables.tables"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsTablesDataSource_name_regex(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsTablesDataSource_name_regex_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_tables.tables"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.id", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.instance_name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.table_name", fmt.Sprintf("testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.name", "pk2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.type", "String"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.time_to_live", "-1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.max_version", "1"),

					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.0", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.0", fmt.Sprintf("testAcc%d", randInt)),
				),
			},
			{
				Config: testAccCheckAlicloudOtsTablesDataSource_name_regex_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_tables.tables"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsTablesDataSource_all(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsTablesDataSource_all_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_tables.tables"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.id", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.instance_name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.table_name", fmt.Sprintf("testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.name", "pk2"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.primary_key.1.type", "String"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.time_to_live", "-1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.0.max_version", "1"),

					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.0", fmt.Sprintf("tf-testAcc%d:testAcc%d", randInt, randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.0", fmt.Sprintf("testAcc%d", randInt)),
				),
			},
			{
				Config: testAccCheckAlicloudOtsTablesDataSource_all_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_tables.tables"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "tables.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_tables.tables", "names.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudOtsTablesDataSource_basic(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = [
		{
          name = "pk1"
	      type = "Integer"
	    },
		{
          name = "pk2"
          type = "String"
        },
      ]
	  time_to_live = -1
	  max_version = 1
	}

	data "alicloud_ots_tables" "tables" {
	  instance_name = "${alicloud_ots_table.basic.instance_name}"
	}
	`, randInt)
}

func testAccCheckAlicloudOtsTablesDataSource_name_regex_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = [
		{
          name = "pk1"
	      type = "Integer"
	    },
		{
          name = "pk2"
          type = "String"
        },
      ]
	  time_to_live = -1
	  max_version = 1
	}

	data "alicloud_ots_tables" "tables" {
	  instance_name = "${alicloud_ots_table.basic.instance_name}"
	  name_regex = "${alicloud_ots_table.basic.table_name}"
	}
	`, randInt)
}

func testAccCheckAlicloudOtsTablesDataSource_name_regex_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = [
		{
          name = "pk1"
	      type = "Integer"
	    },
		{
          name = "pk2"
          type = "String"
        },
      ]
	  time_to_live = -1
	  max_version = 1
	}

	data "alicloud_ots_tables" "tables" {
	  instance_name = "${alicloud_ots_table.basic.instance_name}"
	  name_regex = "${alicloud_ots_table.basic.table_name}-fake"
	}
	`, randInt)
}

func testAccCheckAlicloudOtsTablesDataSource_ids_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = [
		{
          name = "pk1"
	      type = "Integer"
	    },
		{
          name = "pk2"
          type = "String"
        },
      ]
	  time_to_live = -1
	  max_version = 1
	}

	data "alicloud_ots_tables" "tables" {
	  instance_name = "${alicloud_ots_table.basic.instance_name}"
      ids = [ "${alicloud_ots_table.basic.id}" ]
	}
	`, randInt)
}

func testAccCheckAlicloudOtsTablesDataSource_ids_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = [
		{
          name = "pk1"
	      type = "Integer"
	    },
		{
          name = "pk2"
          type = "String"
        },
      ]
	  time_to_live = -1
	  max_version = 1
	}

	data "alicloud_ots_tables" "tables" {
	  instance_name = "${alicloud_ots_table.basic.instance_name}"
	  ids = [ "${alicloud_ots_table.basic.id}-fake" ]
	}
	`, randInt)
}

func testAccCheckAlicloudOtsTablesDataSource_all_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = [
		{
          name = "pk1"
	      type = "Integer"
	    },
		{
          name = "pk2"
          type = "String"
        },
      ]
	  time_to_live = -1
	  max_version = 1
	}

	data "alicloud_ots_tables" "tables" {
	  instance_name = "${alicloud_ots_table.basic.instance_name}"
      ids = [ "${alicloud_ots_table.basic.id}" ]
	  name_regex = "${alicloud_ots_table.basic.table_name}"
	}
	`, randInt)
}

func testAccCheckAlicloudOtsTablesDataSource_all_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = [
		{
          name = "pk1"
	      type = "Integer"
	    },
		{
          name = "pk2"
          type = "String"
        },
      ]
	  time_to_live = -1
	  max_version = 1
	}

	data "alicloud_ots_tables" "tables" {
	  instance_name = "${alicloud_ots_table.basic.instance_name}"
	  ids = [ "${alicloud_ots_table.basic.id}-fake" ]
	  name_regex = "${alicloud_ots_table.basic.table_name}-fake"
	}
	`, randInt)
}
