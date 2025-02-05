package alicloud

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_sae_application",
		&resource.Sweeper{
			Name: "alicloud_sae_application",
			F:    testSweepSaeApplication,
		})
}

func testSweepSaeApplication(region string) error {
	prefixes := []string{
		"tftestacc",
	}
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.SaeSupportRegions)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := make(map[string]*string)

	request["ContainCustom"] = StringPointer(strconv.FormatBool(true))
	action := "/pop/v1/sam/namespace/describeNamespaceList"
	response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
	if err != nil {
		log.Printf("[ERROR] %s got an error: %s", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	namespace, _ := resp.([]interface{})

	for _, v := range namespace {
		// item namespace
		item := v.(map[string]interface{})

		action := "/pop/v1/sam/app/listApplications"
		request["NamespaceId"] = StringPointer(item["NamespaceId"].(string))
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_application", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.Applications", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Applications", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				app_name := ""
				if val, exist := item["AppName"]; exist {
					app_name = val.(string)
				}
				if strings.Contains(strings.ToLower(app_name), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Sae Application: %s (%s)", item["AppName"], item["AppId"])
				continue
			}
			action := "/pop/v1/sam/app/deleteApplication"
			request = map[string]*string{
				"AppId": StringPointer(item["AppId"].(string)),
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}

					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil && !IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) {
				return WrapError(err)
			}
			log.Printf("[INFO] Delete Sae Application  success: %v ", item["AppId"])
		}
	}
	return nil
}

// package_type = Image
func TestAccAliCloudSAEApplication_basicImage(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":           name,
					"namespace_id":       "${alicloud_sae_namespace.default.namespace_id}",
					"package_type":       "Image",
					"app_description":    name + "desc",
					"vswitch_id":         "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"image_url":          fmt.Sprintf("registry-vpc.%s.aliyuncs.com/sae-demo-image/consumer:1.0", defaultRegionToTest),
					"replicas":           "5",
					"cpu":                "500",
					"memory":             "2048",
					"package_version":    strconv.FormatInt(time.Now().Unix(), 10) + "create",
					"micro_registration": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":           name,
						"namespace_id":       CHECKSET,
						"package_type":       "Image",
						"app_description":    name + "desc",
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"image_url":          CHECKSET,
						"replicas":           "5",
						"cpu":                "500",
						"memory":             "2048",
						"micro_registration": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_url":       fmt.Sprintf("registry-vpc.%s.aliyuncs.com/lxepoo/apache-php5", defaultRegionToTest),
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "image_url",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_url":       fmt.Sprintf("registry-vpc.%s.aliyuncs.com/lxepoo/apache-php5", defaultRegionToTest),
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_enable_application_scaling_rule": "true",
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_enable_application_scaling_rule": "true",
						"package_version":                      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pre_stop":        `{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pre_stop":        "{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_configs":     `[{\"logDir\":\"/root/logs/hsf/hsf.log\"}]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_configs":     "[{\"logDir\":\"/root/logs/hsf/hsf.log\"}]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"readiness":       `{\"exec\":{\"command\":[\"sleep\",\"6s\"]},\"initialDelaySeconds\":15,\"periodSeconds\":30,\"timeoutSeconds\":12}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"readiness":       "{\"exec\":{\"command\":[\"sleep\",\"6s\"]},\"initialDelaySeconds\":15,\"periodSeconds\":30,\"timeoutSeconds\":12}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"liveness":        `{\"exec\":{\"command\":[\"sleep\",\"5s\"]},\"initialDelaySeconds\":10,\"periodSeconds\":30,\"timeoutSeconds\":11}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"liveness":        "{\"exec\":{\"command\":[\"sleep\",\"5s\"]},\"initialDelaySeconds\":10,\"periodSeconds\":30,\"timeoutSeconds\":11}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timezone":        "Asia/Beijing",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timezone":        "Asia/Beijing",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory": "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory": "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"edas_container_version": "3.5.3",
					"package_version":        strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"edas_container_version": "3.5.3",
						"package_version":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"php_config_location": "/usr/local/etc/php/php.ini",
					"php_config":          "k1=v1",
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"php_config_location": "/usr/local/etc/php/php.ini",
						"php_config":          "k1=v1",
						"package_version":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_enable_application_scaling_rule": "true",
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_enable_application_scaling_rule": "true",
						"package_version":                      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
					"package_version":          strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
						"package_version":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"termination_grace_period_seconds": "30",
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"termination_grace_period_seconds": "30",
						"package_version":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"batch_wait_time": "10",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"batch_wait_time": "10",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"post_start":      `{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"post_start":      "{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command":         "sleep",
					"command_args":    `[\"1d\"]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_args":    "[\"1d\"]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"envs":            `[{\"name\":\"envtmp\",\"value\":\"0\"}]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"envs":            "[{\"name\":\"envtmp\",\"value\":\"0\"}]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_host_alias": `[{\"hostName\":\"samplehost\",\"ip\":\"127.0.0.1\"}]`,
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_host_alias": "[{\"hostName\":\"samplehost\",\"ip\":\"127.0.0.1\"}]",
						"package_version":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":    strconv.FormatInt(time.Now().Unix(), 10) + "micro_registration",
					"micro_registration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":    CHECKSET,
						"micro_registration": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10) + "min_ready_instances",
					"min_ready_instances": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":     CHECKSET,
						"min_ready_instances": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "tomcat_config",
					"tomcat_config":   `{\"contextPath\":\"/\",\"maxThreads\":400,\"port\":8080,\"uriEncoding\":\"UTF-8\",\"useBodyEncodingForUri\":\"false\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"tomcat_config":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "update_strategy",
					"update_strategy": `{\"batchUpdate\":{\"batch\":1,\"batchWaitTime\":1,\"releaseType\":\"auto\"},\"type\":\"BatchUpdate\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"update_strategy": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

// package_type = FatJar Skip testing because api does not support STS authentication
func SkipTestAccAliCloudSAEApplication_basicFatJar(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":        name,
					"namespace_id":    "${alicloud_sae_namespace.default.namespace_id}",
					"package_url":     fmt.Sprintf("http://edas-hz.oss-%s.aliyuncs.com/demo/1.0/hello-sae.jar", defaultRegionToTest),
					"package_type":    "FatJar",
					"app_description": name + "desc",
					"jdk":             "Open JDK 8",
					"vswitch_id":      "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":          "${data.alicloud_vpcs.default.ids.0}",
					"replicas":        "5",
					"cpu":             "500",
					"memory":          "2048",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_description": name + "desc",
						"app_name":        name,
						"namespace_id":    CHECKSET,
						"package_type":    "FatJar",
						"vswitch_id":      CHECKSET,
						"package_url":     CHECKSET,
						"vpc_id":          CHECKSET,
						"replicas":        "5",
						"cpu":             "500",
						"memory":          "2048",
						"jdk":             "Open JDK 8",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pre_stop":        `{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "pre_stop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pre_stop":        "{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_configs":     `[{\"logDir\":\"/root/logs/hsf/hsf.log\"}]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "sls_configs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_configs":     "[{\"logDir\":\"/root/logs/hsf/hsf.log\"}]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"readiness":       `{\"exec\":{\"command\":[\"sleep\",\"6s\"]},\"initialDelaySeconds\":15,\"periodSeconds\":30,\"timeoutSeconds\":12}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "readiness",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"readiness":       "{\"exec\":{\"command\":[\"sleep\",\"6s\"]},\"initialDelaySeconds\":15,\"periodSeconds\":30,\"timeoutSeconds\":12}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"liveness":        `{\"exec\":{\"command\":[\"sleep\",\"5s\"]},\"initialDelaySeconds\":10,\"periodSeconds\":30,\"timeoutSeconds\":11}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "liveness",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"liveness":        "{\"exec\":{\"command\":[\"sleep\",\"5s\"]},\"initialDelaySeconds\":10,\"periodSeconds\":30,\"timeoutSeconds\":11}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timezone":        "Asia/Beijing",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "timezone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timezone":        "Asia/Beijing",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu":             "1000",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "cpu",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cpu":             "1000",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"memory":          "4096",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "memory",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"memory":          "4096",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"php_config_location": "/usr/local/etc/php/php.ini",
					"php_config":          "k1=v1",
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10) + "php_config",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"php_config_location": "/usr/local/etc/php/php.ini",
						"php_config":          "k1=v1",
						"package_version":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_ready_instances": "2",
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10) + "min_ready_instances",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_ready_instances": "2",
						"package_version":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_enable_application_scaling_rule": "true",
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10) + "scaling_rule",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_enable_application_scaling_rule": "true",
						"package_version":                      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
					"package_version":          strconv.FormatInt(time.Now().Unix(), 10) + "arms_config",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
						"package_version":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"termination_grace_period_seconds": "30",
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10) + "termination_grace_period",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"termination_grace_period_seconds": "30",
						"package_version":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"post_start":      `{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "post_start",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"post_start":      "{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"batch_wait_time": "10",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "batch_wait_time",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"batch_wait_time": "10",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command":         "sleep",
					"command_args":    `[\"1d\"]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "command",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_args":    "[\"1d\"]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"envs":            `[{\"name\":\"envtmp\",\"value\":\"0\"}]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "envs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"envs":            "[{\"name\":\"envtmp\",\"value\":\"0\"}]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_host_alias": `[{\"hostName\":\"samplehost\",\"ip\":\"127.0.0.1\"}]`,
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "custom_host_alias",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_host_alias": "[{\"hostName\":\"samplehost\",\"ip\":\"127.0.0.1\"}]",
						"package_version":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "oss_ak_id",
					"oss_ak_id":       os.Getenv("ALICLOUD_ACCESS_KEY"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"oss_ak_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "oss_ak_secret",
					"oss_ak_secret":   os.Getenv("ALICLOUD_SECRET_KEY"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"oss_ak_secret":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "oss_mount_descs",
					"oss_mount_descs": `[{\"bucketName\":\"` + "${alicloud_oss_bucket.default.bucket}" + `\",\"bucketPath\":\"/\",\"mountPath\":\"/tmp/oss\",\"readOnly\":false}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"oss_mount_descs": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":       strconv.FormatInt(time.Now().Unix(), 10) + "config_map_mount_desc",
					"config_map_mount_desc": `[{\"configMapId\":` + "${alicloud_sae_config_map.default.id}" + `,\"key\":\"env.home\",\"mountPath\":\"/tmp/configmap\"}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":       CHECKSET,
						"config_map_mount_desc": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

// package_type = War
func TestAccAliCloudSAEApplication_basicWar(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":          name,
					"app_description":   name + "desc",
					"namespace_id":      "${alicloud_sae_namespace.default.namespace_id}",
					"package_url":       fmt.Sprintf("http://edas-hz.oss-%s.aliyuncs.com/demo/1.0/hello-sae.war?spm=5176.12834076.0.0.60326a68Uw5yB4&file=hello-sae.war", defaultRegionToTest),
					"package_type":      "War",
					"web_container":     "apache-tomcat-8.5.42",
					"jdk":               "Open JDK 8",
					"replicas":          "5",
					"vswitch_id":        "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"cpu":               "500",
					"memory":            "2048",
					"war_start_options": "-Dhello=war",
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":          name,
						"namespace_id":      CHECKSET,
						"app_description":   name + "desc",
						"package_url":       CHECKSET,
						"package_type":      "War",
						"web_container":     "apache-tomcat-8.5.42",
						"jdk":               "Open JDK 8",
						"replicas":          "5",
						"vswitch_id":        CHECKSET,
						"vpc_id":            CHECKSET,
						"cpu":               "500",
						"memory":            "2048",
						"war_start_options": "-Dhello=war",
						"package_version":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "web_container",
					"web_container":   "apache-tomcat-8.5.58",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"web_container":   "apache-tomcat-8.5.58",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "war_start_options",
					"war_start_options": "-Dhello=war2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":   CHECKSET,
						"war_start_options": "-Dhello=war2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pre_stop":        `{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "pre_stop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pre_stop":        "{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sls_configs":     `[{\"logDir\":\"/root/logs/hsf/hsf.log\"}]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "sls_configs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sls_configs":     "[{\"logDir\":\"/root/logs/hsf/hsf.log\"}]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"readiness":       `{\"exec\":{\"command\":[\"sleep\",\"6s\"]},\"initialDelaySeconds\":15,\"periodSeconds\":30,\"timeoutSeconds\":12}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "readiness",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"readiness":       "{\"exec\":{\"command\":[\"sleep\",\"6s\"]},\"initialDelaySeconds\":15,\"periodSeconds\":30,\"timeoutSeconds\":12}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"liveness":        `{\"exec\":{\"command\":[\"sleep\",\"5s\"]},\"initialDelaySeconds\":10,\"periodSeconds\":30,\"timeoutSeconds\":11}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "liveness",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"liveness":        "{\"exec\":{\"command\":[\"sleep\",\"5s\"]},\"initialDelaySeconds\":10,\"periodSeconds\":30,\"timeoutSeconds\":11}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timezone":        "Asia/Beijing",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "timezone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timezone":        "Asia/Beijing",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "cpu",
					"cpu":             "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"cpu":             "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "memory",
					"memory":          "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"memory":          "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"php_config_location": "/usr/local/etc/php/php.ini",
					"php_config":          "k1=v1",
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10) + "php_config",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"php_config_location": "/usr/local/etc/php/php.ini",
						"php_config":          "k1=v1",
						"package_version":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_ready_instances": "2",
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10) + "min_ready_instances",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_ready_instances": "2",
						"package_version":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_enable_application_scaling_rule": "true",
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10) + "enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_enable_application_scaling_rule": "true",
						"package_version":                      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"batch_wait_time": "10",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "batch_wait_time",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"batch_wait_time": "10",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
					"package_version":          strconv.FormatInt(time.Now().Unix(), 10) + "arms_config",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
						"package_version":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"termination_grace_period_seconds": "30",
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10) + "termination_grace",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"termination_grace_period_seconds": "30",
						"package_version":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"post_start":      `{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "post_start",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"post_start":      "{\"exec\":{\"command\":[\"cat\",\"/etc/group\"]}}",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command":         "sleep",
					"command_args":    `[\"1d\"]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "command",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_args":    "[\"1d\"]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"envs":            `[{\"name\":\"envtmp\",\"value\":\"0\"}]`,
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "envs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"envs":            "[{\"name\":\"envtmp\",\"value\":\"0\"}]",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_host_alias": `[{\"hostName\":\"samplehost\",\"ip\":\"127.0.0.1\"}]`,
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "custom_host_alias",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_host_alias": "[{\"hostName\":\"samplehost\",\"ip\":\"127.0.0.1\"}]",
						"package_version":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":      "${data.alicloud_vswitches.default.vswitches.1.id}",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "vswitch_id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":      CHECKSET,
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "tomcat_config_v2",
					"tomcat_config_v2": []map[string]interface{}{
						{
							"port":                      "8081",
							"max_threads":               "401",
							"context_path":              "/",
							"uri_encoding":              "UTF-8",
							"use_body_encoding_for_uri": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":    CHECKSET,
						"tomcat_config_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

// package_type = Php
func TestAccAliCloudSAEApplication_basicPhp(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":                         name,
					"package_type":                     "PhpZip",
					"replicas":                         "1",
					"namespace_id":                     "${alicloud_sae_namespace.default.namespace_id}",
					"vpc_id":                           "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                       "${data.alicloud_vswitches.default.vswitches.0.id}",
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10),
					"package_url":                      fmt.Sprintf("http://edas-hz.oss-%s.aliyuncs.com/demo/1.0/hello-sae-php.zip?spm=5176.12834076.0.0.3b016a68nuAvPc&file=hello-sae-php.zip", defaultRegionToTest),
					"cpu":                              "500",
					"memory":                           "2048",
					"app_description":                  name + "desc",
					"php_config":                       "usrer=A",
					"php_config_location":              "/usr/local/etc/php/php.ini",
					"security_group_id":                "${alicloud_security_group.default.id}",
					"termination_grace_period_seconds": "30",
					"timezone":                         "Asia/Shanghai",
					"micro_registration":               "0",
					"php":                              "PHP-FPM 7.2",
					"programming_language":             "php",
					"liveness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "30",
							"period_seconds":        "30",
							"timeout_seconds":       "5",
							"http_get": []map[string]interface{}{
								{
									"path":   "/",
									"port":   "80",
									"scheme": "HTTP",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":                         name,
						"package_type":                     "PhpZip",
						"replicas":                         "1",
						"namespace_id":                     CHECKSET,
						"vpc_id":                           CHECKSET,
						"vswitch_id":                       CHECKSET,
						"package_version":                  CHECKSET,
						"package_url":                      CHECKSET,
						"cpu":                              "500",
						"memory":                           "2048",
						"app_description":                  name + "desc",
						"php_config":                       "usrer=A",
						"php_config_location":              "/usr/local/etc/php/php.ini",
						"security_group_id":                CHECKSET,
						"termination_grace_period_seconds": "30",
						"timezone":                         "Asia/Shanghai",
						"micro_registration":               "0",
						"php":                              "PHP-FPM 7.2",
						"programming_language":             "php",
						"liveness_v2.#":                    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas":        "5",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "replicas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas":        "5",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":      "${data.alicloud_vswitches.default.vswitches.1.id}",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "vswitch_id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":      CHECKSET,
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "package_url",
					"package_url":     fmt.Sprintf("http://edas-hz.oss-%s.aliyuncs.com/demo/1.0/hello-sae-php.zip?spm=5176.12834076.0.0.3b016a68nuAv&file=hello-sae-php.zip", defaultRegionToTest),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"package_url":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "cpu",
					"cpu":             "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"cpu":             "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "memory",
					"memory":          "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"memory":          "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_enable_application_scaling_rule": "true",
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10) + "enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_enable_application_scaling_rule": "true",
						"package_version":                      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10) + "php_config",
					"php_config":          "k1=v1",
					"php_config_location": "/usr/local/etc/php/php.ini",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":     CHECKSET,
						"php_config":          "k1=v1",
						"php_config_location": "/usr/local/etc/php/php.ini",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":          strconv.FormatInt(time.Now().Unix(), 10) + "arms_config",
					"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":          CHECKSET,
						"php_arms_config_location": "/usr/local/etc/php/conf.d/arms.ini",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "security_group_id",
					"security_group_id": "${alicloud_security_group.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":   CHECKSET,
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10) + "termination_grace",
					"termination_grace_period_seconds": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":                  CHECKSET,
						"termination_grace_period_seconds": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "timezone",
					"timezone":        "Asia/Beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"timezone":        "Asia/Beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":    strconv.FormatInt(time.Now().Unix(), 10) + "micro_registration",
					"micro_registration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":    CHECKSET,
						"micro_registration": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "php",
					"php":             "PHP-FPM 7.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"php":             "PHP-FPM 7.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "liveness_v2",
					"liveness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "31",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"http_get": []map[string]interface{}{
								{
									"path":                "/tmp",
									"port":                "8080",
									"scheme":              "HTTPS",
									"key_word":            "SAE",
									"is_contain_key_word": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"liveness_v2.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "readiness_v2",
					"readiness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "31",
							"period_seconds":        "10",
							"timeout_seconds":       "3",
							"http_get": []map[string]interface{}{
								{
									"path":                "/",
									"port":                "80",
									"scheme":              "HTTP",
									"key_word":            "SAE",
									"is_contain_key_word": "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"readiness_v2.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

// Skip testing because api has bug
func SkipTestAccAliCloudSAEApplication_basicNewImage(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":                         name,
					"package_type":                     "Image",
					"replicas":                         "1",
					"namespace_id":                     "${alicloud_sae_namespace.default.namespace_id}",
					"vpc_id":                           "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                       "${data.alicloud_vswitches.default.vswitches.0.id}",
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10),
					"image_url":                        fmt.Sprintf("registry-vpc.%s.aliyuncs.com/sae-demo-image/consumer:1.0", defaultRegionToTest),
					"cpu":                              "500",
					"memory":                           "2048",
					"app_description":                  name + "desc",
					"security_group_id":                "${alicloud_security_group.default.id}",
					"termination_grace_period_seconds": "30",
					"timezone":                         "Asia/Shanghai",
					"micro_registration":               "0",
					"programming_language":             "other",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":                         name,
						"package_type":                     "Image",
						"replicas":                         "1",
						"namespace_id":                     CHECKSET,
						"vpc_id":                           CHECKSET,
						"vswitch_id":                       CHECKSET,
						"package_version":                  CHECKSET,
						"image_url":                        CHECKSET,
						"cpu":                              "500",
						"memory":                           "2048",
						"app_description":                  name + "desc",
						"security_group_id":                CHECKSET,
						"termination_grace_period_seconds": "30",
						"timezone":                         "Asia/Shanghai",
						"micro_registration":               "0",
						"programming_language":             "other",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas":        "5",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "replicas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas":        "5",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":      "${data.alicloud_vswitches.default.vswitches.1.id}",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "vswitch_id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":      CHECKSET,
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "image_url",
					"image_url":       fmt.Sprintf("registry-vpc.%s.aliyuncs.com/lxepoo/apache-php5", defaultRegionToTest),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"image_url":       fmt.Sprintf("registry-vpc.%s.aliyuncs.com/lxepoo/apache-php5", defaultRegionToTest),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "cpu",
					"cpu":             "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"cpu":             "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "memory",
					"memory":          "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"memory":          "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_enable_application_scaling_rule": "true",
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10) + "enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_enable_application_scaling_rule": "true",
						"package_version":                      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "security_group_id",
					"security_group_id": "${alicloud_security_group.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":   CHECKSET,
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10) + "termination_grace",
					"termination_grace_period_seconds": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":                  CHECKSET,
						"termination_grace_period_seconds": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "timezone",
					"timezone":        "Asia/Beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"timezone":        "Asia/Beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":    strconv.FormatInt(time.Now().Unix(), 10) + "micro_registration",
					"micro_registration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":    CHECKSET,
						"micro_registration": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command":         "sleep",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "command_args_v2",
					"command_args_v2": []string{"1d"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command":           "sleep",
						"package_version":   CHECKSET,
						"command_args_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "custom_host_alias_v2",
					"custom_host_alias_v2": []map[string]interface{}{
						{
							"host_name": "samplehost",
							"ip":        "127.0.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":        CHECKSET,
						"custom_host_alias_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "liveness_v2",
					"liveness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "10",
							"period_seconds":        "30",
							"timeout_seconds":       "12",
							"exec": []map[string]interface{}{
								{
									"command": []string{"sleep", "5s"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"liveness_v2.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "readiness_v2",
					"readiness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "15",
							"period_seconds":        "30",
							"timeout_seconds":       "12",
							"exec": []map[string]interface{}{
								{
									"command": []string{"sleep", "6s"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"readiness_v2.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "post_start_v2",
					"post_start_v2": []map[string]interface{}{
						{
							"exec": []map[string]interface{}{
								{
									"command": []string{"cat", "/etc/group"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"post_start_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "pre_stop_v2",
					"pre_stop_v2": []map[string]interface{}{
						{
							"exec": []map[string]interface{}{
								{
									"command": []string{"cat", "/etc/group"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"pre_stop_v2.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "update_strategy_v2",
					"update_strategy_v2": []map[string]interface{}{
						{
							"type": "BatchUpdate",
							"batch_update": []map[string]interface{}{
								{
									"release_type":    "auto",
									"batch":           "1",
									"batch_wait_time": "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":      CHECKSET,
						"update_strategy_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "pvtz_discovery_svc",
					"pvtz_discovery_svc": []map[string]interface{}{
						{
							"service_name": "testpvtz",
							"namespace_id": "${alicloud_sae_namespace.default.namespace_id}",
							"enable":       "true",
							"port_protocols": []map[string]interface{}{
								{
									"port":     "81",
									"protocol": "UDP",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":      CHECKSET,
						"pvtz_discovery_svc.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

func SkipTestAccAliCloudSAEApplication_basicNewFatJar(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence1)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":                             name,
					"package_type":                         "FatJar",
					"replicas":                             "2",
					"namespace_id":                         "${alicloud_sae_namespace.default.namespace_id}",
					"vpc_id":                               "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                           "${data.alicloud_vswitches.default.vswitches.0.id}",
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10),
					"package_url":                          fmt.Sprintf("http://edas-hz.oss-%s.aliyuncs.com/demo/1.0/hello-sae.jar", defaultRegionToTest),
					"cpu":                                  "500",
					"memory":                               "2048",
					"jdk":                                  "Open JDK 8",
					"jar_start_options":                    "-Dtest=tesst",
					"jar_start_args":                       "hello",
					"app_description":                      name + "desc",
					"auto_config":                          "false",
					"auto_enable_application_scaling_rule": "false",
					"batch_wait_time":                      "1",
					"change_order_desc":                    "1",
					"deploy":                               "true",
					"enable_ahas":                          "true",
					"oss_ak_id":                            os.Getenv("ALICLOUD_ACCESS_KEY"),
					"oss_ak_secret":                        os.Getenv("ALICLOUD_SECRET_KEY"),
					"security_group_id":                    "${alicloud_security_group.default.id}",
					"termination_grace_period_seconds":     "30",
					"timezone":                             "Asia/Shanghai",
					"micro_registration":                   "0",
					"custom_host_alias_v2": []map[string]interface{}{
						{
							"host_name": "www.hello.com",
							"ip":        "127.0.0.1",
						},
					},
					"config_map_mount_desc_v2": []map[string]interface{}{
						{
							"config_map_id": "${alicloud_sae_config_map.default.id}",
							"mount_path":    "/tmp/configmap/a.txt",
							"key":           "env.home",
						},
					},
					"liveness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "30",
							"period_seconds":        "5",
							"timeout_seconds":       "2",
							"exec": []map[string]interface{}{
								{
									"command": []string{"sh", "-c", "cat /home/admin/start.sh"},
								},
							},
						},
					},
					"readiness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "30",
							"period_seconds":        "30",
							"timeout_seconds":       "5",
							"exec": []map[string]interface{}{
								{
									"command": []string{"sh", "-c", "cat /home/admin/start.sh"},
								},
							},
						},
					},
					"post_start_v2": []map[string]interface{}{
						{
							"exec": []map[string]interface{}{
								{
									"command": []string{"cat", "/etc/group"},
								},
							},
						},
					},
					"pre_stop_v2": []map[string]interface{}{
						{
							"exec": []map[string]interface{}{
								{
									"command": []string{"cat", "/etc/group"},
								},
							},
						},
					},
					"update_strategy_v2": []map[string]interface{}{
						{
							"type": "BatchUpdate",
							"batch_update": []map[string]interface{}{
								{
									"release_type":    "auto",
									"batch":           "1",
									"batch_wait_time": "1",
								},
							},
						},
					},
					"nas_configs": []map[string]interface{}{
						{
							"nas_id":       "${alicloud_nas_mount_target.default.file_system_id}",
							"nas_path":     "/",
							"mount_path":   "/tmp/nas",
							"mount_domain": "${alicloud_nas_mount_target.default.mount_target_domain}",
							"read_only":    "false",
						},
					},
					"pvtz_discovery_svc": []map[string]interface{}{
						{
							"service_name": "test",
							"namespace_id": "${alicloud_sae_namespace.default.namespace_id}",
							"enable":       "true",
							"port_protocols": []map[string]interface{}{
								{
									"port":     "80",
									"protocol": "TCP",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":                             name,
						"package_type":                         "FatJar",
						"replicas":                             "2",
						"namespace_id":                         CHECKSET,
						"vpc_id":                               CHECKSET,
						"vswitch_id":                           CHECKSET,
						"package_version":                      CHECKSET,
						"package_url":                          CHECKSET,
						"cpu":                                  "500",
						"memory":                               "2048",
						"jdk":                                  "Open JDK 8",
						"jar_start_options":                    "-Dtest=tesst",
						"jar_start_args":                       "hello",
						"app_description":                      name + "desc",
						"auto_config":                          "false",
						"auto_enable_application_scaling_rule": "false",
						"batch_wait_time":                      "1",
						"change_order_desc":                    "1",
						"deploy":                               "true",
						"enable_ahas":                          "true",
						"oss_ak_id":                            CHECKSET,
						"oss_ak_secret":                        CHECKSET,
						"security_group_id":                    CHECKSET,
						"termination_grace_period_seconds":     "30",
						"timezone":                             "Asia/Shanghai",
						"micro_registration":                   "0",
						"custom_host_alias_v2.#":               "1",
						"config_map_mount_desc_v2.#":           "1",
						"liveness_v2.#":                        "1",
						"readiness_v2.#":                       "1",
						"post_start_v2.#":                      "1",
						"pre_stop_v2.#":                        "1",
						"update_strategy_v2.#":                 "1",
						"nas_configs.#":                        "1",
						"pvtz_discovery_svc.#":                 "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas":        "5",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "replicas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas":        "5",
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":      "${data.alicloud_vswitches.default.vswitches.1.id}",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "vswitch_id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":      CHECKSET,
						"package_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "package_url",
					"package_url":     fmt.Sprintf("https://edas-hz.oss-%s.aliyuncs.com/demo/1.0/hello-sae.jar?file=hello-sae.jar", defaultRegionToTest),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"package_url":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "cpu",
					"cpu":             "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"cpu":             "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "memory",
					"memory":          "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"memory":          "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "jdk",
					"jdk":             "Dragonwell 8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"jdk":             "Dragonwell 8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "jar_start_options",
					"jar_start_options": "-Dtest=tesst2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":   CHECKSET,
						"jar_start_options": "-Dtest=tesst2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "jar_start_args",
					"jar_start_args":  "hello2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"jar_start_args":  "hello2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":                      strconv.FormatInt(time.Now().Unix(), 10) + "enable",
					"auto_enable_application_scaling_rule": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":                      CHECKSET,
						"auto_enable_application_scaling_rule": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "change_order_desc",
					"change_order_desc": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":   CHECKSET,
						"change_order_desc": "12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "enable_ahas",
					"enable_ahas":     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"enable_ahas":     "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":       strconv.FormatInt(time.Now().Unix(), 10) + "enable_grey_tag_route",
					"enable_grey_tag_route": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":       CHECKSET,
						"enable_grey_tag_route": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":          strconv.FormatInt(time.Now().Unix(), 10) + "min_ready_instance_ratio",
					"min_ready_instance_ratio": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":          CHECKSET,
						"min_ready_instance_ratio": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":   strconv.FormatInt(time.Now().Unix(), 10) + "security_group_id",
					"security_group_id": "${alicloud_security_group.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":   CHECKSET,
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":                  strconv.FormatInt(time.Now().Unix(), 10) + "termination_grace",
					"termination_grace_period_seconds": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":                  CHECKSET,
						"termination_grace_period_seconds": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "timezone",
					"timezone":        "Asia/Beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"timezone":        "Asia/Beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version":    strconv.FormatInt(time.Now().Unix(), 10) + "micro_registration",
					"micro_registration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":    CHECKSET,
						"micro_registration": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "custom_host_alias_v2",
					"custom_host_alias_v2": []map[string]interface{}{
						{
							"host_name": "www.hello2.com",
							"ip":        "127.0.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":        CHECKSET,
						"custom_host_alias_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "oss_mount_descs_v2",
					"oss_mount_descs_v2": []map[string]interface{}{
						{
							"bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"bucket_path": "/",
							"mount_path":  "/tmp/oss",
							"read_only":   "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":      CHECKSET,
						"oss_mount_descs_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "config_map_mount_desc_v2",
					"config_map_mount_desc_v2": []map[string]interface{}{
						{
							"config_map_id": "${alicloud_sae_config_map.update.id}",
							"mount_path":    "/tmp/configmap",
							"key":           "env.shell",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":            CHECKSET,
						"config_map_mount_desc_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "liveness_v2",
					"liveness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "60",
							"period_seconds":        "5",
							"timeout_seconds":       "6",
							"tcp_socket": []map[string]interface{}{
								{
									"port": "18091",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"liveness_v2.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "readiness_v2",
					"readiness_v2": []map[string]interface{}{
						{
							"initial_delay_seconds": "60",
							"period_seconds":        "5",
							"timeout_seconds":       "6",
							"tcp_socket": []map[string]interface{}{
								{
									"port": "18091",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"readiness_v2.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "post_start_v2",
					"post_start_v2": []map[string]interface{}{
						{
							"exec": []map[string]interface{}{
								{
									"command": []string{"cat", "/home/admin/start.sh"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"post_start_v2.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "pre_stop_v2",
					"pre_stop_v2": []map[string]interface{}{
						{
							"exec": []map[string]interface{}{
								{
									"command": []string{"cat", "/home/admin/start.sh"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"pre_stop_v2.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "nas_configs",
					"nas_configs": []map[string]interface{}{
						{
							"nas_id":       "${alicloud_nas_mount_target.update.file_system_id}",
							"nas_path":     "test",
							"mount_path":   "/tmp/nas12",
							"mount_domain": "${alicloud_nas_mount_target.update.mount_target_domain}",
							"read_only":    "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"nas_configs.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "kafka_configs",
					"kafka_configs": []map[string]interface{}{
						{
							"kafka_instance_id": "${alicloud_alikafka_instance.default.id}",
							"kafka_endpoint":    "${alicloud_alikafka_instance.default.end_point}",
							"kafka_configs": []map[string]interface{}{
								{
									"log_type":    "file_log",
									"log_dir":     "/tmp/kafka/kafka.log",
									"kafka_topic": "${alicloud_alikafka_topic.default.topic}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version": CHECKSET,
						"kafka_configs.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "pvtz_discovery_svc",
					"pvtz_discovery_svc": []map[string]interface{}{
						{
							"service_name": "testpvtz",
							"namespace_id": "${alicloud_sae_namespace.default.namespace_id}",
							"enable":       "true",
							"port_protocols": []map[string]interface{}{
								{
									"port":     "81",
									"protocol": "UDP",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_version":      CHECKSET,
						"pvtz_discovery_svc.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

func TestAccAliCloudSAEApplication_basicTags(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":        name,
					"namespace_id":    "${alicloud_sae_namespace.default.namespace_id}",
					"package_type":    "Image",
					"app_description": name + "desc",
					"vswitch_id":      "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":          "${data.alicloud_vpcs.default.ids.0}",
					"image_url":       fmt.Sprintf("registry-vpc.%s.aliyuncs.com/sae-demo-image/consumer:1.0", defaultRegionToTest),
					"replicas":        "5",
					"cpu":             "500",
					"memory":          "2048",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10) + "create",
					"tags": map[string]string{
						"Created": "tfTestAcc1",
						"For":     "Tftestacc1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":        name,
						"namespace_id":    CHECKSET,
						"package_type":    "Image",
						"app_description": name + "desc",
						"vswitch_id":      CHECKSET,
						"vpc_id":          CHECKSET,
						"image_url":       CHECKSET,
						"replicas":        "5",
						"cpu":             "500",
						"memory":          "2048",
						"tags.%":          "2",
						"tags.Created":    "tfTestAcc1",
						"tags.For":        "Tftestacc1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc7",
						"For":     "Tftestacc7",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc7",
						"tags.For":     "Tftestacc7",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

// Skip testing because code source settings cannot be configured.
func SkipAccAliCloudSAEApplication_basic3(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AliCloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSAEApplicationBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":            name,
					"namespace_id":        "${alicloud_sae_namespace.default.namespace_id}",
					"package_type":        "Image",
					"app_description":     name + "desc",
					"vswitch_id":          "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"image_url":           "${local.image_url}",
					"replicas":            "5",
					"cpu":                 "500",
					"memory":              "2048",
					"acr_instance_id":     "${data.alicloud_cr_ee_instances.default.ids.0}",
					"acr_assume_role_arn": "${data.alicloud_ram_roles.default.roles.0.arn}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":            name,
						"namespace_id":        CHECKSET,
						"package_type":        "Image",
						"app_description":     name + "desc",
						"vswitch_id":          CHECKSET,
						"vpc_id":              CHECKSET,
						"image_url":           CHECKSET,
						"replicas":            "5",
						"cpu":                 "500",
						"memory":              "2048",
						"acr_instance_id":     CHECKSET,
						"acr_assume_role_arn": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_config", "auto_enable_application_scaling_rule", "change_order_desc", "deploy"},
			},
		},
	})
}

func AliCloudSAEApplicationBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "update" {
  		name   = "${var.name}-update"
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_sae_namespace" "default" {
  		namespace_description = var.name
  		namespace_id          = "%s:%s"
  		namespace_name        = var.name
	}

	resource "alicloud_oss_bucket" "default" {
  		bucket = var.name
	}
	
	resource "alicloud_sae_config_map" "default" {
  		namespace_id = alicloud_sae_namespace.default.namespace_id
  		name         = var.name
  		data         = jsonencode({ "env.home" : "/root", "env.shell" : "/bin/sh" })
	}
`, name, defaultRegionToTest, name)
}

func AliCloudSAEApplicationBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "update" {
  		name   = "${var.name}-update"
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_sae_namespace" "default" {
  		namespace_description = var.name
  		namespace_id          = "%s:%s"
  		namespace_name        = var.name
	}

	resource "alicloud_oss_bucket" "default" {
  		bucket = var.name
	}

	resource "alicloud_sae_config_map" "default" {
  		namespace_id = alicloud_sae_namespace.default.namespace_id
  		name         = var.name
  		data         = jsonencode({ "env.home" : "/root", "env.shell" : "/bin/sh" })
	}

	resource "alicloud_sae_config_map" "update" {
  		namespace_id = alicloud_sae_namespace.default.namespace_id
  		name         = "${var.name}update"
  		data         = jsonencode({ "env.home" : "/root", "env.shell" : "/bin/sh" })
	}

	resource "alicloud_nas_file_system" "default" {
  		protocol_type = "NFS"
  		storage_type  = "Performance"
  		description   = var.name
  		encrypt_type  = "1"
  		vpc_id        = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_nas_mount_target" "default" {
  		file_system_id    = alicloud_nas_file_system.default.id
  		access_group_name = "DEFAULT_VPC_GROUP_NAME"
  		vswitch_id        = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_nas_file_system" "update" {
  		protocol_type = "NFS"
  		storage_type  = "Performance"
  		description   = "${var.name}update"
  		encrypt_type  = "1"
  		vpc_id        = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_nas_mount_target" "update" {
		file_system_id    = alicloud_nas_file_system.update.id
		access_group_name = "DEFAULT_VPC_GROUP_NAME"
	  	vswitch_id        = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_alikafka_instance" "default" {
  		name           = var.name
  		partition_num  = "50"
		disk_type      = "1"
		disk_size      = "500"
		deploy_type    = "5"
		io_max         = "20"
		vswitch_id     = data.alicloud_vswitches.default.ids.0
		security_group = alicloud_security_group.default.id
	}

	resource "alicloud_alikafka_topic" "default" {
  		instance_id = alicloud_alikafka_instance.default.id
  		topic       = var.name
  		remark      = var.name
	}
`, name, defaultRegionToTest, name)
}

func AliCloudSAEApplicationBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data alicloud_cr_ee_instances "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = "${data.alicloud_vpcs.default.ids.0}"
	}

	data "alicloud_ram_roles" "default" {
  		name_regex = "^AliyunServiceRoleForSAE$"
	}

	resource "alicloud_cr_ee_namespace" "default" {
  		instance_id        = data.alicloud_cr_ee_instances.default.ids.0
  		name               = var.name
  		auto_create        = true
  		default_visibility = "PRIVATE"
	}

	resource "alicloud_cr_ee_repo" "default" {
  		instance_id = data.alicloud_cr_ee_instances.default.ids.0
  		namespace   = alicloud_cr_ee_namespace.default.name
  		name        = "${var.name}"
  		summary     = "test summary"
  		repo_type   = "PRIVATE"
  		detail      = "test detail"
	}

	resource "alicloud_sae_namespace" "default" {
  		namespace_description = var.name
  		namespace_id          = "%s:%s"
  		namespace_name        = var.name
	}

	locals {
  		image_url = format("%%s-registry-vpc.cn-hangzhou.cr.aliyuncs.com/%%s/%%s", data.alicloud_cr_ee_instances.default.instances.0.name, alicloud_sae_namespace.default.namespace_name, alicloud_cr_ee_repo.default.name)
	}
`, defaultRegionToTest, name, name)
}

var AliCloudSAEApplicationMap0 = map[string]string{
	"package_version":                  CHECKSET,
	"batch_wait_time":                  CHECKSET,
	"enable_ahas":                      CHECKSET,
	"enable_grey_tag_route":            CHECKSET,
	"min_ready_instances":              CHECKSET,
	"min_ready_instance_ratio":         CHECKSET,
	"security_group_id":                CHECKSET,
	"termination_grace_period_seconds": CHECKSET,
	"timezone":                         CHECKSET,
	"envs":                             CHECKSET,
}
