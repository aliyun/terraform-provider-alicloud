package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := make(map[string]*string)

	request["ContainCustom"] = StringPointer(strconv.FormatBool(true))
	action := "/pop/v1/sam/namespace/describeNamespaceList"
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
	if err != nil {
		log.Printf("[ERROR] %s got an error: %s", action, err)
		return nil
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return nil
	}
	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
	}
	result, _ := resp.([]interface{})

	for _, v := range result {
		// item namespace
		item := v.(map[string]interface{})

		action := "/pop/v1/sam/app/listApplications"
		conn, err = client.NewServerlessClient()
		if err != nil {
			return WrapError(err)
		}
		request["NamespaceId"] = StringPointer(item["NamespaceId"].(string))
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_application", action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "AlicloudSaelistApplications", response))
		}
		resp, err := jsonpath.Get("$.Data.Applications", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Applications", response)
		}
		sweeped := false
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
			sweeped = true
			action := "/pop/v1/sam/app/deleteApplication"
			request = map[string]*string{
				"AppId": StringPointer(item["AppId"].(string)),
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) {
						wait()
						return resource.RetryableError(err)
					}

					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return WrapError(err)
			}
			if sweeped {
				time.Sleep(10 * time.Second)
			}
			log.Printf("[INFO] Delete Sae Application  success: %v ", item["AppId"])
		}
	}
	return nil
}

//package_type = Image
func TestAccAlicloudSAEApplication_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEApplicationBasicDependence0)
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
			// todo: The responses of the Update and Get corresponding to mount_desc are inconsistent.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"nas_id":          "${alicloud_nas_file_system.default.id}",
			//		"mount_desc":      `[{mountPath: \"/tmp\", nasPath: \"/\"}]`,
			//		"mount_host":      "${alicloud_nas_mount_target.example.mount_target_domain}",
			//		"package_version": strconv.FormatInt(time.Now().Unix(), 10),
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"nas_id":          CHECKSET,
			//			"mount_desc":      CHECKSET,
			//			"mount_host":      CHECKSET,
			//			"package_version": CHECKSET,
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas": "6",
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
					"min_ready_instances": "2",
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10),
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
					"min_ready_instances": "2",
					"package_version":     strconv.FormatInt(time.Now().Unix(), 10),
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
				ImportStateVerifyIgnore: []string{"auto_enable_application_scaling_rule", "batch_wait_time", "config_map_mount_desc"},
			},
		},
	})
}

//package_type = FatJar
func TestAccAlicloudSAEApplication_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEApplicationBasicDependence0)
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
			// todo: The responses of the Update and Get corresponding to mount_desc are inconsistent.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"nas_id":          "${alicloud_nas_file_system.default.id}",
			//		"package_version": strconv.FormatInt(time.Now().Unix(), 10),
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"nas_id":          CHECKSET,
			//			"package_version": CHECKSET,
			//		}),
			//	),
			//},
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
				ImportStateVerifyIgnore: []string{"auto_enable_application_scaling_rule", "batch_wait_time", "config_map_mount_desc"},
			},
		},
	})
}

//package_type = War
func TestAccAlicloudSAEApplication_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEApplicationBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":        name,
					"app_description": name + "desc",
					"namespace_id":    "${alicloud_sae_namespace.default.namespace_id}",
					"package_url":     fmt.Sprintf("http://edas-hz.oss-%s.aliyuncs.com/demo/1.0/hello-sae.war?spm=5176.12834076.0.0.60326a68Uw5yB4&file=hello-sae.war", defaultRegionToTest),
					"package_type":    "War",
					"web_container":   "apache-tomcat-8.5.42",
					"jdk":             "Open JDK 8",
					"replicas":        "5",
					"vswitch_id":      "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":          "${data.alicloud_vpcs.default.ids.0}",
					"cpu":             "500",
					"memory":          "2048",
					"package_version": strconv.FormatInt(time.Now().Unix(), 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":        name,
						"namespace_id":    CHECKSET,
						"app_description": name + "desc",
						"package_url":     CHECKSET,
						"package_type":    "War",
						"web_container":   "apache-tomcat-8.5.42",
						"jdk":             "Open JDK 8",
						"replicas":        "5",
						"vswitch_id":      CHECKSET,
						"vpc_id":          CHECKSET,
						"cpu":             "500",
						"memory":          "2048",
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
			// todo: The responses of the Update and Get corresponding to mount_desc are inconsistent.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"nas_id":          "${alicloud_nas_file_system.default.id}",
			//		"package_version": strconv.FormatInt(time.Now().Unix(), 10),
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"nas_id":          CHECKSET,
			//			"package_version": CHECKSET,
			//		}),
			//	),
			//},
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
				ImportStateVerifyIgnore: []string{"auto_enable_application_scaling_rule", "batch_wait_time", "config_map_mount_desc"},
			},
		},
	})
}

func TestAccAlicloudSAEApplication_basicTags(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_application.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEApplicationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEApplicationBasicDependence0)
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
				ImportStateVerifyIgnore: []string{"auto_enable_application_scaling_rule", "batch_wait_time", "config_map_mount_desc"},
			},
		},
	})
}

func AlicloudSAEApplicationBasicDependence0(name string) string {
	return fmt.Sprintf(`
data "alicloud_vpcs" "default"	{
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id = "%s:%s"
  namespace_name = var.name
}

#resource "alicloud_nas_file_system" "default" {
 # protocol_type = "NFS"
  #storage_type  = "Performance"
  #description   = var.name
  #encrypt_type = "1"
  #vpc_id        = data.alicloud_vpcs.default.ids.0
#}

#resource "alicloud_nas_mount_target" "example" {
 # file_system_id    = alicloud_nas_file_system.default.id
  #access_group_name = "DEFAULT_VPC_GROUP_NAME"
  #vswitch_id        = data.alicloud_vswitches.default.ids.0
#}

variable "name" {
  default = "%s"
}
`, defaultRegionToTest, name, name)
}

var AlicloudSAEApplicationMap0 = map[string]string{
	"namespace_id": CHECKSET,
}
