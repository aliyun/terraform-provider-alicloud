package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	req.PageNumber = strconv.Itoa(page)
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
		req.PageNumber = strconv.Itoa(page)
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
func TestAccAlicloudDdoscooInstance_multi(t *testing.T) {
	var v ddoscoo.InstanceSpec

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
					"remark":            name + "${count.index}",
					"bandwidth":         "30",
					"base_bandwidth":    "30",
					"service_bandwidth": "100",
					"port_count":        "50",
					"domain_count":      "50",
					"renewal_status":    "AutoRenewal",
					"renew_period":      "12",
					"normal_qps":        "3000",
					"period":            "1",
					"edition":           "coop",
					"function_version":  "0",
					"service_partner":   "coop-line-001",
					"count":             "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func TestAccAlicloudDdosCooInstance_basic(t *testing.T) {
	var v ddoscoo.InstanceSpec
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := resourceAttrInit(resourceId, DdoscooInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdoscooInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccDdoscooInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DdoscooInstanceBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":         "50",
					"renewal_status":    "AutoRenewal",
					"renew_period":      "12",
					"base_bandwidth":    "30",
					"domain_count":      "60",
					"edition":           "coop",
					"function_version":  "0",
					"port_count":        "50",
					"service_bandwidth": "200",
					"service_partner":   "coop-line-001",
					"period":            "1",
					"normal_qps":        "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":         "50",
						"renewal_status":    "AutoRenewal",
						"renew_period":      "12",
						"base_bandwidth":    "30",
						"domain_count":      "60",
						"edition":           "coop",
						"function_version":  "0",
						"port_count":        "50",
						"service_bandwidth": "200",
						"service_partner":   "coop-line-001",
						"period":            "1",
						"normal_qps":        "3000",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"all", "bandwidth", "domain_count", "modify_type", "normal_qps", "period", "port_count", "renew_period", "renewal_status", "resource_group_id", "service_bandwidth", "service_partner"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "updateremark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "updateremark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"normal_qps":  "3100",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"normal_qps":  "3100",
						"modify_type": "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_bandwidth": "100",
					"modify_type":       "Downgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_bandwidth": "100",
						"modify_type":       "Downgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_count":  "60",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_count":  "60",
						"modify_type": "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_count": "50",
					"modify_type":  "Downgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_count": "50",
						"modify_type":  "Downgrade",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"normal_qps":        "3400",
					"service_bandwidth": "400",
					"port_count":        "70",
					"domain_count":      "70",
					"function_version":  "1",
					"remark":            "lastupdate",
					"modify_type":       "Upgrade",
					"edition":           "coop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"normal_qps":        "3400",
						"service_bandwidth": "400",
						"port_count":        "70",
						"domain_count":      "70",
						"function_version":  "1",
						"remark":            "lastupdate",
						"modify_type":       "Upgrade",
						"edition":           "coop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
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
    }`
}

var ddoscooInstanceBasicMap = map[string]string{
	"bandwidth":         "30",
	"base_bandwidth":    "30",
	"service_bandwidth": "100",
	"port_count":        "50",
	"domain_count":      "50",
}
var DdoscooInstanceMap = map[string]string{
	"status": CHECKSET,
}

func DdoscooInstanceBasicdependence(name string) string {
	return ""
}
