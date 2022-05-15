package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ddoscoo_instance", &resource.Sweeper{
		Name: "alicloud_ddoscoo_instance",
		F:    testSweepDdoscooInstances,
	})
}

func testSweepDdoscooInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		log.Printf("error getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var insts []map[string]interface{}

	for {

		action := "DescribeInstances"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			log.Printf("[ERROR] %s get an error %#v", action, err)
		}
		resp, err := jsonpath.Get("$.Instances", response)
		if resp == nil || len(resp.([]interface{})) < 1 || err != nil {
			break
		}

		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			insts = append(insts, item)
		}

		if len(resp.([]interface{})) < PageSizeLarge {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, v := range insts {

		name := fmt.Sprint(v["Remark"])
		instanceId := fmt.Sprint(v["InstanceId"])
		skip := true
		for _, prefix := range prefixes {
			if name != "" && strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// Delete the domain resources
		action := "DescribeDomainResource"
		request := map[string]interface{}{
			"InstanceIds": []string{instanceId},
			"PageSize":    PageSizeSmall,
			"PageNumber":  1,
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
		}
		resp, err := jsonpath.Get("$.WebRules", response)
		if err != nil {
			log.Println(err)
		} else {
			result, _ := resp.([]interface{})
			for _, v := range result {
				domain := fmt.Sprint(v.(map[string]interface{})["Domain"])
				action := "DeleteDomainResource"
				request := map[string]interface{}{
					"Domain": domain,
				}

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] %s got an error: %s", action, err)
				}
			}
		}

		// Delete the ports
		action = "DescribePort"
		request = map[string]interface{}{
			"InstanceId": instanceId,
			"PageSize":   PageSizeLarge,
			"PageNumber": 1,
		}

		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
		}
		resp, err = jsonpath.Get("$.NetworkRules", response)
		if err != nil {
			log.Println(err)
		} else {
			result, _ := resp.([]interface{})
			for _, v := range result {
				action := "DeletePort"
				conn, err := client.NewDdoscooClient()
				if err != nil {
					return WrapError(err)
				}
				request := map[string]interface{}{
					"FrontendPort":     v.(map[string]interface{})["FrontendPort"],
					"FrontendProtocol": v.(map[string]interface{})["FrontendProtocol"],
					"InstanceId":       instanceId,
				}

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] %s got an error: %s", action, err)
				}
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ddoscoo Instance: %s", name)
			continue
		}

		log.Printf("[INFO] Deleting Ddoscoo Instance %s .", fmt.Sprint(v["InstanceId"]))
		action = "ReleaseInstance"
		request = map[string]interface{}{
			"InstanceId": fmt.Sprint(v["InstanceId"]),
			"RegionId":   "cn-hangzhou",
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", fmt.Sprint(v["InstanceId"]), err)
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
		CheckDestroy:  nil,
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
				ImportStateVerifyIgnore: []string{"period", "product_type"},
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
func resourceDdoscooInstanceDependence(name string) string {
	return ""
}

var ddoscooInstanceBasicMap = map[string]string{
	"name":              CHECKSET,
	"bandwidth":         "30",
	"base_bandwidth":    "30",
	"service_bandwidth": "100",
	"port_count":        "50",
	"domain_count":      "50",
}

func TestAccAlicloudDdoscooInstance_intl(t *testing.T) {
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
			testAccPreCheckWithAccountSiteType(t, IntlSite)
			testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"bandwidth":         "30",
					"base_bandwidth":    "30",
					"service_bandwidth": "100",
					"port_count":        "50",
					"domain_count":      "50",
					"product_type":      "ddoscoo_intl",
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
				ImportStateVerifyIgnore: []string{"period", "product_type"},
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
					"domain_count": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_count": "60",
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
					"product_type":      "ddoscoo_intl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"bandwidth":         "60",
						"base_bandwidth":    "60",
						"service_bandwidth": "300",
						"port_count":        "60",
						"domain_count":      "60",
						"product_type":      "ddoscoo_intl",
					}),
				),
			},
		},
	})
}
