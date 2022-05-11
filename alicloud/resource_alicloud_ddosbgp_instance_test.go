package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"

	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ddosbgp_instance", &resource.Sweeper{
		Name: "alicloud_ddosbgp_instance",
		F:    testSweepDdosbgpInstances,
	})
}

func testSweepDdosbgpInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DdosbgpSupportedRegions) {
		log.Printf("[INFO] Does not support this Region")
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
	}

	var insts []ddosbgp.Instance
	req := ddosbgp.CreateDescribeInstanceListRequest()
	req.RegionId = client.RegionId
	req.DdosRegionId = client.RegionId
	req.PageSize = requests.Integer(fmt.Sprint(PageSizeLarge))

	var page = 1
	req.PageNo = requests.Integer(fmt.Sprint(page))
	for {
		raw, err := client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
			return ddosbgpClient.DescribeInstanceList(req)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error %#v", req.GetActionName(), err)
		}
		resp, _ := raw.(*ddosbgp.DescribeInstanceListResponse)
		if resp == nil || len(resp.InstanceList) < 1 {
			break
		}
		insts = append(insts, resp.InstanceList...)

		if len(resp.InstanceList) < PageSizeLarge {
			break
		}

		page++
		req.PageNo = requests.Integer(fmt.Sprint(page))
	}

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
			log.Printf("[INFO] Skipping Ddosbgp Instance: %s", name)
			continue
		}

		log.Printf("[INFO] Deleting Ddosbgp Instance %s .", v.InstanceId)

		releaseReq := ddosbgp.CreateReleaseInstanceRequest()
		releaseReq.InstanceId = v.InstanceId

		_, err := client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
			return ddosbgpClient.ReleaseInstance(releaseReq)
		})
		if err != nil {
			log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", v.InstanceId, err)
		}
	}
	return nil
}

func TestAccAlicloudDdosbgpInstance_basic(t *testing.T) {
	var v ddosbgp.Instance

	resourceId := "alicloud_ddosbgp_instance.default"
	ra := resourceAttrInit(resourceId, ddosbgpInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &DdosbgpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdosbgpInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdosbgpSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":           name,
					"base_bandwidth": "20",
					"bandwidth":      "201",
					"ip_count":       "100",
					"ip_type":        "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "-update",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name":           name,
					"base_bandwidth": "20",
					"bandwidth":      "201",
					"ip_count":       "100",
					"ip_type":        "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"base_bandwidth": "20",
						"bandwidth":      "201",
						"ip_count":       "100",
						"ip_type":        "IPv4",
					}),
				),
			},
		},
	})
}

func resourceDdosbgpInstanceDependence(name string) string {
	return ``
}

var ddosbgpInstanceBasicMap = map[string]string{
	"name":      CHECKSET,
	"bandwidth": "201",
	"ip_count":  "100",
	"ip_type":   "IPv4",
}
