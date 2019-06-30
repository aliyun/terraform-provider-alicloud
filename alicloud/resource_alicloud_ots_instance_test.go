package alicloud

import (
	"fmt"
	"testing"

	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ots_instance", &resource.Sweeper{
		Name: "alicloud_ots_instance",
		F:    testSweepOtsInstances,
	})
}

func testSweepOtsInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var insts []ots.InstanceInfo
	req := ots.CreateListInstanceRequest()
	req.RegionId = client.RegionId
	req.Method = "GET"
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNum = requests.NewInteger(1)
	for {
		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.ListInstance(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving OTS Instances: %s", err)
		}
		resp, _ := raw.(*ots.ListInstanceResponse)
		if resp == nil || len(resp.InstanceInfos.InstanceInfo) < 1 {
			break
		}
		insts = append(insts, resp.InstanceInfos.InstanceInfo...)

		if len(resp.InstanceInfos.InstanceInfo) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNum); err != nil {
			return err
		} else {
			req.PageNum = page
		}
	}
	sweeped := false

	for _, v := range insts {
		name := v.InstanceName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping OTS Instance: %s", name)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting OTS Instance %s table stores.", name)
		raw, err := otsService.client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			return tableStoreClient.ListTable()
		})
		if err != nil {
			log.Printf("[ERROR] List OTS Instance %s table stores got an error: %#v.", name, err)
		}
		tables, _ := raw.(*tablestore.ListTableResponse)
		if tables != nil && len(tables.TableNames) > 0 {
			for _, t := range tables.TableNames {
				req := new(tablestore.DeleteTableRequest)
				req.TableName = t
				if _, err := otsService.client.WithTableStoreClient(name, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
					return tableStoreClient.DeleteTable(req)
				}); err != nil {
					log.Printf("[ERROR] Delete OTS Instance %s table store %s got an error: %#v.", name, t, err)
				}
			}
			time.Sleep(30 * time.Second)
		}
		log.Printf("[INFO] Deleting OTS Instance: %s", name)
		req := ots.CreateDeleteInstanceRequest()
		req.InstanceName = name
		_, err = client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete OTS Instance (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(3 * time.Minute)
	}
	return nil
}

func TestAccAlicloudOtsInstanceCapacity_basic(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsCapacity)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccAlicloudOtsInstanceCapacity_updateAccessBy(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsCapacity)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOtsInstanceUpdateAccessBy(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Vpc"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsCapacity)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
		},
	})

}

func TestAccAlicloudOtsInstanceCapacity_updateTags(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsCapacity)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOtsInstanceUpdateTags(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsCapacity)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Updated", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.From", "TF"),
				),
			},
		},
	})

}
func TestAccAlicloudOtsInstanceCapacity_updateAll(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsCapacity)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOtsInstanceUpdateAll(string(OtsCapacity), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "ConsoleOrVpc"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsCapacity)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Updated", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.From", "TF"),
				),
			},
		},
	})

}

// Test HighPerformance instance
func TestAccAlicloudOtsInstanceHighPerformance_basic(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsHighPerformance)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccAlicloudOtsInstanceHighPerformance_updateAccessBy(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsHighPerformance)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOtsInstanceUpdateAccessBy(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Vpc"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsHighPerformance)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
		},
	})

}

func TestAccAlicloudOtsInstanceHighPerformance_updateTags(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsHighPerformance)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOtsInstanceUpdateTags(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsHighPerformance)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Updated", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.From", "TF"),
				),
			},
		},
	})

}
func TestAccAlicloudOtsInstanceHighPerformance_updateAll(t *testing.T) {
	var instance ots.InstanceInfo
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_ots_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsInstance(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "Any"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsHighPerformance)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Created", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOtsInstanceUpdateAll(string(OtsHighPerformance), rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOtsInstanceExist("alicloud_ots_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "name", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "accessed_by", "ConsoleOrVpc"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "instance_type", string(OtsHighPerformance)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "description", fmt.Sprintf("tf-testAcc%d", rand)),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.Updated", "TF"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.For", "acceptance test"),
					resource.TestCheckResourceAttr("alicloud_ots_instance.foo", "tags.From", "TF"),
				),
			},
		},
	})

}

func testAccCheckOtsInstanceExist(n string, instance *ots.InstanceInfo) resource.TestCheckFunc {
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

		response, err := otsService.DescribeOtsInstance(rs.Primary.ID)

		if err != nil {
			return err
		}
		instance = &response
		return nil
	}
}

func testAccCheckOtsInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ots_instance" {
			continue
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		otsService := OtsService{client}

		if _, err := otsService.DescribeOtsInstance(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Ots instance %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccOtsInstance(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "%s"
	  tags = {
		Created = "TF"
		For = "acceptance test"
	  }
	}
	`, rand, instanceType)
}

func testAccOtsInstanceUpdateAccessBy(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "%s"
	  tags = {
		Created = "TF"
		For = "acceptance test"
	  }
	}
	`, rand, instanceType)
}

func testAccOtsInstanceUpdateTags(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "%s"
	  tags = {
		Updated = "TF"
		For = "acceptance test"
		From = "TF"
	  }
	}
	`, rand, instanceType)
}

func testAccOtsInstanceUpdateAll(instanceType string, rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}
	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "ConsoleOrVpc"
	  instance_type = "%s"
	  tags = {
		Updated = "TF"
		For = "acceptance test"
		From = "TF"
	  }
	}
	`, rand, instanceType)
}
