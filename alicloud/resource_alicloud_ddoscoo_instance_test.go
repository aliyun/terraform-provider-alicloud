package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"

	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ddoscoo_instance", &resource.Sweeper{
		Name: "alicloud_ddoscoo_instance",
		F:    testSweepDdoscooInstances,
	})
}

func testSweepDdoscooInstances(region string) error {
	if !testSweepPreCheckWithRegions(region, true, []connectivity.Region{connectivity.Hangzhou}) {
		log.Printf("[INFO] only supported region: cn-hangzhou")
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var insts []ddoscoo.Instance
	req := ddoscoo.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = strconv.Itoa(PageSizeLarge)

	var page = 1
	req.PageNo = strconv.Itoa(page)
	for {
		raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error describing OTS Instances: %s", err)
		}
		resp, _ := raw.(*ddoscoo.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances) < 1 {
			break
		}
		insts = append(insts, resp.Instances...)

		if len(resp.Instances) < PageSizeLarge {
			break
		}

		page++
		req.PageNo = strconv.Itoa(page)
	}
	sweeped := false

	for _, v := range insts {
		name := v.Remark
		skip := true
		for _, prefix := range prefixes {
			if name != "" && strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ddoscoo Instance: %s", name)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ddoscoo Instance %s .", v.InstanceId)

		releaseReq := ddoscoo.CreateReleaseInstanceRequest()
		releaseReq.InstanceId = v.InstanceId

		_, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.ReleaseInstance(releaseReq)
		})
		if err != nil {
			log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", v.InstanceId, err)
		}

		time.Sleep(30 * time.Second)
	}
	if sweeped {
		time.Sleep(3 * time.Minute)
	}
	return nil
}

func TestAccAlicloudDdoscooInstance_basic(t *testing.T) {
	var v []ddoscoo.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ddoscoo_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDdoscooDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDdoscooInstanceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDdoscooExists("alicloud_ddoscoo_instance.foo", v),
					resource.TestCheckResourceAttrSet("alicloud_ddoscoo_instance.foo", "band_width"),
				),
			},
		},
	})

}

func testAccCheckDdoscooExists(n string, instances []ddoscoo.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Instance ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ddoscoo.CreateDescribeInstancesRequest()
		request.InstanceIds = "[\"" + rs.Primary.ID + "\"]"
		request.PageNo = "1"
		request.PageSize = "20"
		raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(request)
		})

		if err == nil {
			response, _ := raw.(*ddoscoo.DescribeInstancesResponse)
			instances = response.Instances
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckDdoscooDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ddoscoo_instance" {
			continue
		}

		// Try to find the instance
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ddoscoo.CreateDescribeInstancesRequest()
		request.InstanceIds = "[\"" + rs.Primary.ID + "\"]"
		request.PageNo = "1"
		request.PageSize = "20"

		_, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(request)
		})

		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func testAccDdoscooInstanceConfig() string {
	return fmt.Sprintf(`
	resource "alicloud_ddoscoo_instance" "foo" {
	  business_endpoint        = "business.aliyuncs.com"
      band_width               = "30"
      base_band_width          = "30"
      service_band_width       = "200"
      port_count               = "50"
      domain_count             = "50"
	}`)
}
