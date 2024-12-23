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
	if testSweepPreCheckWithRegions(region, true, connectivity.DdoscooSupportedRegions) {
		log.Printf("[INFO] Skipping ddoscoo unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		log.Printf("error getting AliCloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	var response map[string]interface{}
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
					_, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, false)
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
				request := map[string]interface{}{
					"FrontendPort":     v.(map[string]interface{})["FrontendPort"],
					"FrontendProtocol": v.(map[string]interface{})["FrontendProtocol"],
					"InstanceId":       instanceId,
				}

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					_, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
		if !sweepAll() {
			if skip {
				log.Printf("[INFO] Skipping Ddoscoo Instance: %s", name)
				continue
			}
		}

		log.Printf("[INFO] Deleting Ddoscoo Instance %s .", fmt.Sprint(v["InstanceId"]))
		action = "ReleaseInstance"
		request = map[string]interface{}{
			"InstanceId": fmt.Sprint(v["InstanceId"]),
			"RegionId":   "cn-hangzhou",
		}
		_, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Deleting Instance %s got an error: %#v.", fmt.Sprint(v["InstanceId"]), err)
		}
	}
	return nil
}

func TestAccAliCloudDdosCooInstance_basic0(t *testing.T) {
	var v ddoscoo.Instance
	testAccPreCheckWithRegions(t, true, connectivity.DdoscooInstanceSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := resourceAttrInit(resourceId, ddosCooInstanceBasicMap)
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
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"base_bandwidth":    "30",
					"bandwidth":         "30",
					"service_bandwidth": "100",
					"port_count":        "50",
					"domain_count":      "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"base_bandwidth":    "30",
						"bandwidth":         "30",
						"service_bandwidth": "100",
						"port_count":        "50",
						"domain_count":      "50",
					}),
				),
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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period"},
			},
		},
	})

}

func TestAccAliCloudDdosCooInstance_basic0_twin(t *testing.T) {
	var v ddoscoo.Instance
	testAccPreCheckWithRegions(t, true, connectivity.DdoscooInstanceSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := resourceAttrInit(resourceId, ddosCooInstanceBasicMap)
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
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"base_bandwidth":    "30",
					"bandwidth":         "30",
					"service_bandwidth": "100",
					"port_count":        "50",
					"domain_count":      "50",
					"edition_sale":      "coop",
					"address_type":      "Ipv6",
					"bandwidth_mode":    "2",
					"product_type":      "ddoscoo",
					"period":            "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"base_bandwidth":    "30",
						"bandwidth":         "30",
						"service_bandwidth": "100",
						"port_count":        "50",
						"domain_count":      "50",
						"edition_sale":      "coop",
						"address_type":      "Ipv6",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period"},
			},
		},
	})

}

func TestAccAliCloudDdosCooInstance_basic0_intl(t *testing.T) {
	var v ddoscoo.Instance
	testAccPreCheckWithRegions(t, true, connectivity.DdoscooInstanceSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := resourceAttrInit(resourceId, ddosCooInstanceBasicMap)
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
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"base_bandwidth":    "30",
					"bandwidth":         "30",
					"service_bandwidth": "100",
					"port_count":        "50",
					"domain_count":      "50",
					"edition_sale":      "coop",
					"product_type":      "ddoscoo_intl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"base_bandwidth":    "30",
						"bandwidth":         "30",
						"service_bandwidth": "100",
						"port_count":        "50",
						"domain_count":      "50",
						"edition_sale":      "coop",
					}),
				),
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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period"},
			},
		},
	})
}

func TestAccAliCloudDdosCooInstance_basic1_dip(t *testing.T) {
	var v ddoscoo.Instance
	testAccPreCheckWithRegions(t, true, connectivity.DdoscooInstanceIntlSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := resourceAttrInit(resourceId, ddoscooInstanceBasicDipMap)
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
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"port_count":       "50",
					"domain_count":     "50",
					"normal_bandwidth": "100",
					"normal_qps":       "500",
					"product_plan":     "0",
					"function_version": "0",
					"product_type":     "ddosDip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"port_count":       "50",
						"domain_count":     "50",
						"normal_bandwidth": "100",
						"normal_qps":       "500",
						"product_plan":     "0",
						"function_version": "0",
					}),
				),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period"},
			},
		},
	})

}

func TestAccAliCloudDdosCooInstance_basic1_dip_intl(t *testing.T) {
	var v ddoscoo.Instance
	testAccPreCheckWithRegions(t, true, connectivity.DdoscooInstanceIntlSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := resourceAttrInit(resourceId, ddoscooInstanceBasicDipMap)
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
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"port_count":       "5",
					"domain_count":     "10",
					"normal_bandwidth": "100",
					"normal_qps":       "500",
					"product_plan":     "3",
					"function_version": "0",
					"product_type":     "ddosDip",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"port_count":       "5",
						"domain_count":     "10",
						"normal_bandwidth": "100",
						"normal_qps":       "500",
						"product_plan":     "3",
						"function_version": "0",
					}),
				),
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period"},
			},
		},
	})

}

var ddosCooInstanceBasicMap = map[string]string{
	"normal_bandwidth": CHECKSET,
	"normal_qps":       CHECKSET,
	"edition_sale":     CHECKSET,
	"product_plan":     CHECKSET,
	"address_type":     CHECKSET,
	"function_version": CHECKSET,
	"ip":               CHECKSET,
}

var ddoscooInstanceBasicDipMap = map[string]string{
	"normal_bandwidth": CHECKSET,
	"normal_qps":       CHECKSET,
	"product_plan":     CHECKSET,
	"address_type":     CHECKSET,
	"function_version": CHECKSET,
	"ip":               CHECKSET,
}

func resourceDdoscooInstanceDependence(name string) string {
	return ""
}
