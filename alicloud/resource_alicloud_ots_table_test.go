package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudOtsTableStoreCapatity(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86400"),
				),
			},
		},
	})

}
func TestAccAlicloudOtsTableStoreCapatity_updateMaxVersion(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateMaxVersion(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "2"),
				),
			},
		},
	})

}
func TestAccAlicloudOtsTableStoreCapatity_updateTimeToLive(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateTimeToLive(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "86401"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudOtsTableStoreCapatity_updateDeviationCellVersion(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateDeviationCellVersion(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86401"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsTableStoreCapatity_updateAll(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86400"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateAll(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "86401"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86401"),
				),
			},
		},
	})

}
func TestAccAlicloudOtsTableStoreHighPerformance(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86400"),
				),
			},
		},
	})

}
func TestAccAlicloudOtsTableStoreHighPerformance_updateMaxVersion(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateMaxVersion(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "2"),
				),
			},
		},
	})

}
func TestAccAlicloudOtsTableStoreHighPerformance_updateTimeToLive(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateTimeToLive(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "86401"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsTableStoreHighPerformance_updateDeviationCellVersion(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateDeviationCellVersion(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86401"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsTableStoreHighPerformance_updateAll(t *testing.T) {
	var table tablestore.DescribeTableResponse
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: "alicloud_ots_table.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "-1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86400"),
				),
			},
			{
				Config: testAccOtsTableStoreUpdateAll(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					testAccCheckOtsTableExist("alicloud_ots_table.basic", &table),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "instance_name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "table_name", fmt.Sprintf("testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.name", "pk1"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "primary_key.0.type", "Integer"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "time_to_live", "86401"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "max_version", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_table.basic", "deviation_cell_version_in_sec", "86401"),
				),
			},
		},
	})

}

func testAccCheckOtsTableExist(n string, table *tablestore.DescribeTableResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found OTS table: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no OTS table ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		otsService := OtsService{client}
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		response, err := otsService.DescribeOtsTable(split[0], split[1])

		if err != nil {
			return fmt.Errorf("Error finding OTS table %s: %#v", rs.Primary.ID, err)
		}

		table = response
		return nil
	}
}

func testAccCheckOtsTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ots_table" {
			continue
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		otsService := OtsService{client}
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if _, err := otsService.DescribeOtsTable(split[0], split[1]); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("error! Ots table still exists")
	}

	return nil
}

func testAccOtsTableStore(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = {
	    name = "pk1"
	    type = "Integer"
	  }
	  time_to_live = -1
	  max_version = 1
      deviation_cell_version_in_sec = "86400"
	}
	`, rand, instanceType)
}

func testAccOtsTableStoreUpdateMaxVersion(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = {
	    name = "pk1"
	    type = "Integer"
	  }
	  time_to_live = -1
	  max_version = 2
	}
	`, rand, instanceType)
}
func testAccOtsTableStoreUpdateTimeToLive(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = {
	    name = "pk1"
	    type = "Integer"
	  }
	  time_to_live = 86401
	  max_version = 1
	}
	`, rand, instanceType)
}
func testAccOtsTableStoreUpdateDeviationCellVersion(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = {
	    name = "pk1"
	    type = "Integer"
	  }
      time_to_live = -1
	  max_version = 1
      deviation_cell_version_in_sec = "86401"
	}
	`, rand, instanceType)
}
func testAccOtsTableStoreUpdateAll(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "basic" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  table_name = "${var.name}"
	  primary_key = {
	    name = "pk1"
	    type = "Integer"
	  }
	  time_to_live = 86401
	  max_version = 2
      deviation_cell_version_in_sec = "86401"
	}
	`, rand, instanceType)
}
