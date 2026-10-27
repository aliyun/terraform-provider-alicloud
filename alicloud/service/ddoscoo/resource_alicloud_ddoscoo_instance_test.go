package ddoscoo_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ddoscoo_instance", &resource.Sweeper{
		Name: "alicloud_ddoscoo_instance",
		F:    testSweepDdosCooInstances,
	})
}

func testSweepDdosCooInstances(region string) error {
	if sweep.SweepPreCheckWithRegions(region, true, connectivity.DdoscooSupportedRegions) {
		log.Printf("[INFO] Skipping ddoscoo unsupported region: %s", region)
		return nil
	}
	rawClient, err := sweep.SharedClientForRegion(region)
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
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var insts []map[string]interface{}

	for {
		action := "DescribeInstances"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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

		if len(resp.([]interface{})) < helper.PageSizeLarge {
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
			"PageSize":    helper.PageSizeSmall,
			"PageNumber":  1,
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
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

				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(1*time.Minute, func() *retry.RetryError {
					_, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
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
			"PageSize":   helper.PageSizeLarge,
			"PageNumber": 1,
		}

		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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

				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(1*time.Minute, func() *retry.RetryError {
					_, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] %s got an error: %s", action, err)
				}
			}
		}
		if !sweep.SweepAll() {
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
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, connectivity.DdoscooInstanceSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosCooInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosCooInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period", "modify_type"},
			},
		},
	})

}

func TestAccAliCloudDdosCooInstance_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, connectivity.DdoscooInstanceSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosCooInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosCooInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period", "modify_type"},
			},
		},
	})

}

func TestAccAliCloudDdosCooInstance_basic0_intl(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, connectivity.DdoscooInstanceSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosCooInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosCooInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.IntlSite)
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period", "modify_type"},
			},
		},
	})
}

func TestAccAliCloudDdosCooInstance_basic1_dip(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, connectivity.DdoscooInstanceIntlSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosCooInstanceMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosCooInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
					"normal_bandwidth": "150",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"normal_bandwidth": "150",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"normal_qps": "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"normal_qps": "3000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_plan": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_plan": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_version": "1",
					"modify_type":      "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_version": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period", "modify_type"},
			},
		},
	})

}

func TestAccAliCloudDdosCooInstance_basic1_dip_intl(t *testing.T) {
	var v map[string]interface{}
	acctest.PreCheckWithRegions(t, true, connectivity.DdoscooInstanceIntlSupportedRegions)
	resourceId := "alicloud_ddoscoo_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosCooInstanceMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosCooInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.IntlSite)
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"port_count":       "5",
					"domain_count":     "10",
					"normal_bandwidth": "100",
					"normal_qps":       "500",
					"product_plan":     "0",
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
				Config: testAccConfig(map[string]interface{}{
					"port_count": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_count": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_count": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_count": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"normal_bandwidth": "150",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"normal_bandwidth": "150",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"normal_qps": "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"normal_qps": "3000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_plan": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_plan": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_version": "1",
					"modify_type":      "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_version": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "bandwidth_mode", "period", "modify_type"},
			},
		},
	})

}

var AliCloudDdosCooInstanceMap0 = map[string]string{
	"normal_bandwidth": acctest.CHECKSET,
	"normal_qps":       acctest.CHECKSET,
	"edition_sale":     acctest.CHECKSET,
	"product_plan":     acctest.CHECKSET,
	"address_type":     acctest.CHECKSET,
	"function_version": acctest.CHECKSET,
	"ip":               acctest.CHECKSET,
	"status":           acctest.CHECKSET,
	"create_time":      acctest.CHECKSET,
}

var AliCloudDdosCooInstanceMap1 = map[string]string{
	"normal_bandwidth": acctest.CHECKSET,
	"normal_qps":       acctest.CHECKSET,
	"product_plan":     acctest.CHECKSET,
	"address_type":     acctest.CHECKSET,
	"function_version": acctest.CHECKSET,
	"ip":               acctest.CHECKSET,
	"status":           acctest.CHECKSET,
	"create_time":      acctest.CHECKSET,
}

func AliCloudDdosCooInstanceBasicDependence0(name string) string {
	return ""
}
