package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ddoscoo_instance", &resource.Sweeper{
		Name: "alicloud_ddoscoo_instance",
		F:    testSweepDdoscooInstances,
	})
}

func testSweepDdoscooInstances(region string) error {
	if testSweepPreCheckWithRegions(region, true, []connectivity.Region{connectivity.Hangzhou}) {
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
			log.Printf("[ERROR] %s get an error %#v", req.GetActionName(), err)
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

		log.Printf("[INFO] Deleting Ddoscoo Instance %s .", v.InstanceId)

		releaseReq := ddoscoo.CreateReleaseInstanceRequest()
		releaseReq.InstanceId = v.InstanceId

		_, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.ReleaseInstance(releaseReq)
		})
		if err != nil {
			log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", v.InstanceId, err)
		}
	}
	return nil
}

func TestAccAlicloudDdoscooInstance_basic(t *testing.T) {
	var v ddoscoo.Instance

	resourceId := "alicloud_ddoscoo_instance.default"
	ra := resourceAttrInit(resourceId, ddoscooInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdoscooInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"bandwidth":         "30",
					"base_bandwidth":    "30",
					"service_bandwidth": "100",
					"port_count":        "50",
					"domain_count":      "50",
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
					"bandwidth": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"base_bandwidth": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_bandwidth": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_bandwidth": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_bandwidth": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_count": "55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_count": "55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_count": "55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_count": "55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"bandwidth":         "60",
					"base_bandwidth":    "60",
					"service_bandwidth": "300",
					"port_count":        "60",
					"domain_count":      "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"bandwidth":         "60",
						"base_bandwidth":    "60",
						"service_bandwidth": "300",
						"port_count":        "60",
						"domain_count":      "60",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudDdoscooInstance_multi(t *testing.T) {
	var v ddoscoo.Instance

	resourceId := "alicloud_ddoscoo_instance.default.1"
	ra := resourceAttrInit(resourceId, ddoscooInstanceBasicMap)

	serviceFunc := func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDdoscooInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name + "${count.index}",
					"bandwidth":         "30",
					"base_bandwidth":    "30",
					"service_bandwidth": "100",
					"port_count":        "50",
					"domain_count":      "50",
					"count":             "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func resourceDdoscooInstanceDependence(name string) string {
	return `
    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }
`
}

var ddoscooInstanceBasicMap = map[string]string{
	"name":              CHECKSET,
	"bandwidth":         "30",
	"base_bandwidth":    "30",
	"service_bandwidth": "100",
	"port_count":        "50",
	"domain_count":      "50",
}
