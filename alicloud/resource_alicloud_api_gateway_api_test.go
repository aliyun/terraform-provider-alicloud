package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_api_gateway_api", &resource.Sweeper{
		Name: "alicloud_api_gateway_api",
		F:    testSweepApiGatewayApi,
	})
}

func testSweepApiGatewayApi(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.ApiGatewayNoSupportedRegions) {
		log.Printf("[INFO] Skipping API Gateway unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	req := cloudapi.CreateDescribeApisRequest()
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApis(req)
	})
	if err != nil {
		log.Printf("[ERROR] %s got an error %#v", req.GetActionName(), err)
		return nil
	}
	apis := raw.(*cloudapi.DescribeApisResponse)

	swept := false

	for _, v := range apis.ApiSummarys.ApiSummary {
		name := v.ApiName
		id := v.ApiId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping api: %s", name)
			continue
		}
		swept = true

		log.Printf("[INFO] Deleting Api: %s", name)

		req := cloudapi.CreateDeleteApiRequest()
		req.GroupId = v.GroupId
		req.ApiId = id

		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteApi(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Api (%s): %s", name, err)
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAliCloudApigatewayApi_basic(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "${alicloud_api_gateway_group.default.name}",
					"group_id":          "${alicloud_api_gateway_group.default.id}",
					"description":       "tf_testAcc_api description",
					"auth_type":         "APP",
					"force_nonce_check": "true",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path",
						"mode":     "MAPPING",
					}},
					"service_type": "HTTP",
					"http_service_config": []map[string]string{{
						"address":   "http://apigateway-backend.alicloudapi.com:8080",
						"method":    "GET",
						"path":      "/web/cloudapi",
						"timeout":   "20",
						"aone_name": "cloudapi-openapi",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                            name,
						"force_nonce_check":               "true",
						"http_service_config.0.address":   "http://apigateway-backend.alicloudapi.com:8080",
						"http_service_config.0.method":    "GET",
						"http_service_config.0.path":      "/web/cloudapi",
						"http_service_config.0.timeout":   "20",
						"http_service_config.0.aone_name": "cloudapi-openapi",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf_testAcc_api description_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf_testAcc_api description_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stage_names": []string{
						"RELEASE",
						"PRE",
						"TEST",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stage_names.#": "3",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"request_config": []map[string]string{{
			//			"protocol": "HTTP",
			//			"method":   "GET",
			//			"path":     "/test/path/test",
			//			"mode":     "MAPPING",
			//		}},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"request_config.0.path": "/test/path/test",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_service_config": []map[string]string{{
						"address":   "http://apigateway-backend.alicloudapi.com:8080",
						"method":    "GET",
						"path":      "/web/cloudapi/update",
						"timeout":   "20",
						"aone_name": "cloudapi-openapi",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_service_config.0.path": "/web/cloudapi/update",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"service_type": "MOCK",
			//		"mock_service_config": []map[string]string{{
			//			"result":    "this is a mock test",
			//			"aone_name": "cloudapi-openapi",
			//		}},
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"service_type": "MOCK",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${alicloud_api_gateway_group.default.name}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api description",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path",
						"mode":     "MAPPING",
					}},
					"service_type": "HTTP",
					"http_service_config": []map[string]string{{
						"address":   "http://apigateway-backend.alicloudapi.com:8080",
						"method":    "GET",
						"path":      "/web/cloudapi",
						"timeout":   "20",
						"aone_name": "cloudapi-openapi",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                       name,
						"service_type":               "HTTP",
						"description":                "tf_testAcc_api description",
						"request_config.0.path":      "/test/path",
						"http_service_config.0.path": "/web/cloudapi",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_param(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	//testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "${alicloud_api_gateway_group.default.name}",
					"group_id":          "${alicloud_api_gateway_group.default.id}",
					"description":       "tf_testAcc_api description",
					"auth_type":         "APP",
					"force_nonce_check": "true",
					"request_config": []map[string]string{{
						"protocol":    "HTTP",
						"method":      "POST",
						"path":        "/test/path",
						"mode":        "MAPPING",
						"body_format": "FORM",
					}},
					"service_type": "HTTP",
					"http_service_config": []map[string]string{{
						"address":   "http://apigateway-backend.alicloudapi.com:8080",
						"method":    "GET",
						"path":      "/web/cloudapi",
						"timeout":   "20",
						"aone_name": "cloudapi-openapi",
					}},
					"request_parameters": []map[string]string{{
						"name":          "request_param",
						"type":          "STRING",
						"required":      "OPTIONAL",
						"in":            "QUERY",
						"in_service":    "QUERY",
						"name_service":  "service_param",
						"description":   "request_desc",
						"default_value": "request_default",
					}},
					"constant_parameters": []map[string]string{{
						"name":        "constant_param",
						"value":       "constant_value",
						"in":          "HEAD",
						"description": "constant_desc",
					}},
					"system_parameters": []map[string]string{{
						"name":         "CaClientIp",
						"name_service": "clientIP",
						"in":           "HEAD",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "name", name),
					resource.TestCheckResourceAttr(resourceId, "force_nonce_check", "true"),
					resource.TestCheckResourceAttr(resourceId, "request_config.0.method", "POST"),
					// request_parameters
					resource.TestCheckResourceAttr(resourceId, "request_parameters.#", "1"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.name", "request_param"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.type", "STRING"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.required", "OPTIONAL"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.in", "QUERY"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.in_service", "QUERY"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.name_service", "service_param"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.description", "request_desc"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.default_value", "request_default"),
					// constant_parameters
					resource.TestCheckResourceAttr(resourceId, "constant_parameters.#", "1"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.name", "constant_param"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.value", "constant_value"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.in", "HEAD"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.description", "constant_desc"),
					// system_parameters
					resource.TestCheckResourceAttr(resourceId, "system_parameters.#", "1"),
					checkRequestParametersAttr(resourceId, "system_parameters.0.name", "CaClientIp"),
					checkRequestParametersAttr(resourceId, "system_parameters.0.name_service", "clientIP"),
					checkRequestParametersAttr(resourceId, "system_parameters.0.in", "HEAD"),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"request_parameters": []map[string]string{{
						"name":          "request_arg",
						"type":          "INT",
						"required":      "REQUIRED",
						"in":            "HEAD",
						"in_service":    "HEAD",
						"name_service":  "service_arg",
						"description":   "request_description",
						"default_value": "1",
					}},
					"constant_parameters": []map[string]string{{
						"name":        "constant_arg",
						"value":       "value",
						"in":          "QUERY",
						"description": "constant_description",
					}},
					"system_parameters": []map[string]string{{
						"name":         "CaDomain",
						"name_service": "domain",
						"in":           "QUERY",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					// request_parameters
					resource.TestCheckResourceAttr(resourceId, "request_parameters.#", "1"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.name", "request_arg"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.type", "INT"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.required", "REQUIRED"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.in", "HEAD"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.in_service", "HEAD"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.name_service", "service_arg"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.description", "request_description"),
					checkRequestParametersAttr(resourceId, "request_parameters.0.default_value", "1"),
					// constant_parameters
					resource.TestCheckResourceAttr(resourceId, "constant_parameters.#", "1"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.name", "constant_arg"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.value", "value"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.in", "QUERY"),
					checkRequestParametersAttr(resourceId, "constant_parameters.0.description", "constant_description"),
					// system_parameters
					resource.TestCheckResourceAttr(resourceId, "system_parameters.#", "1"),
					checkRequestParametersAttr(resourceId, "system_parameters.0.name", "CaDomain"),
					checkRequestParametersAttr(resourceId, "system_parameters.0.name_service", "domain"),
					checkRequestParametersAttr(resourceId, "system_parameters.0.in", "QUERY"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_post(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "${alicloud_api_gateway_group.default.name}",
					"group_id":          "${alicloud_api_gateway_group.default.id}",
					"description":       "tf_testAcc_api description",
					"auth_type":         "APP",
					"force_nonce_check": "true",
					"request_config": []map[string]string{{
						"protocol":    "HTTP",
						"method":      "POST",
						"path":        "/test/path",
						"mode":        "MAPPING",
						"body_format": "FORM",
					}},
					"service_type": "HTTP",
					"http_service_config": []map[string]string{{
						"address":               "http://apigateway-backend.alicloudapi.com:8080",
						"method":                "POST",
						"path":                  "/web/cloudapi",
						"timeout":               "20",
						"aone_name":             "cloudapi-openapi",
						"content_type_category": "DEFAULT",
						"content_type_value":    "application/x-www-form-urlencoded; charset=UTF-8",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                        name,
						"force_nonce_check":                           "true",
						"request_config.0.protocol":                   "HTTP",
						"request_config.0.method":                     "POST",
						"request_config.0.path":                       "/test/path",
						"request_config.0.mode":                       "MAPPING",
						"request_config.0.body_format":                "FORM",
						"http_service_config.0.address":               "http://apigateway-backend.alicloudapi.com:8080",
						"http_service_config.0.method":                "POST",
						"http_service_config.0.path":                  "/web/cloudapi",
						"http_service_config.0.timeout":               "20",
						"http_service_config.0.aone_name":             "cloudapi-openapi",
						"http_service_config.0.content_type_category": "DEFAULT",
						"http_service_config.0.content_type_value":    "application/x-www-form-urlencoded; charset=UTF-8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_service_config": []map[string]string{{
						"address":               "https://www.aliyun.com",
						"method":                "PUT",
						"path":                  "/web/cloudapi/update",
						"timeout":               "30",
						"aone_name":             "cloudapi-task",
						"content_type_category": "CUSTOM",
						"content_type_value":    "application/json",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_service_config.0.address":               "https://www.aliyun.com",
						"http_service_config.0.method":                "PUT",
						"http_service_config.0.path":                  "/web/cloudapi/update",
						"http_service_config.0.timeout":               "30",
						"http_service_config.0.aone_name":             "cloudapi-task",
						"http_service_config.0.content_type_category": "CUSTOM",
						"http_service_config.0.content_type_value":    "application/json",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_convertBackend(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence_vpc)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "${alicloud_api_gateway_group.default.name}" + "_http",
					"group_id":          "${alicloud_api_gateway_group.default.id}",
					"description":       "tf_testAcc_api http",
					"auth_type":         "APP",
					"force_nonce_check": "true",
					"request_config": []map[string]string{{
						"protocol":    "HTTP",
						"method":      "POST",
						"path":        "/test/path",
						"mode":        "MAPPING",
						"body_format": "FORM",
					}},
					"service_type": "HTTP",
					"http_service_config": []map[string]string{{
						"address":   "http://apigateway-backend.alicloudapi.com:8080",
						"method":    "POST",
						"path":      "/web/cloudapi",
						"timeout":   "20",
						"aone_name": "cloudapi-openapi",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                            name + "_http",
						"description":                     "tf_testAcc_api http",
						"force_nonce_check":               "true",
						"request_config.0.protocol":       "HTTP",
						"request_config.0.method":         "POST",
						"request_config.0.path":           "/test/path",
						"request_config.0.mode":           "MAPPING",
						"request_config.0.body_format":    "FORM",
						"http_service_config.0.address":   "http://apigateway-backend.alicloudapi.com:8080",
						"http_service_config.0.method":    "POST",
						"http_service_config.0.path":      "/web/cloudapi",
						"http_service_config.0.timeout":   "20",
						"http_service_config.0.aone_name": "cloudapi-openapi",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              "${alicloud_api_gateway_group.default.name}" + "_vpc",
					"group_id":          "${alicloud_api_gateway_group.default.id}",
					"description":       "tf_testAcc_api vpc",
					"auth_type":         "ANONYMOUS",
					"force_nonce_check": "false",
					"request_config": []map[string]string{{
						"protocol":    "HTTPS",
						"method":      "PUT",
						"path":        "/test/path/vpc",
						"mode":        "MAPPING_PASSTHROUGH",
						"body_format": "STREAM",
					}},
					"service_type": "HTTP-VPC",
					"http_vpc_service_config": []map[string]string{{
						"name":      "${alicloud_api_gateway_vpc_access.default.name}",
						"method":    "POST",
						"path":      "/web/cloudapi/vpc",
						"timeout":   "20",
						"aone_name": "cloudapi-openapi",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                name + "_vpc",
						"description":                         "tf_testAcc_api vpc",
						"auth_type":                           "ANONYMOUS",
						"force_nonce_check":                   "false",
						"request_config.0.protocol":           "HTTPS",
						"request_config.0.method":             "PUT",
						"request_config.0.path":               "/test/path/vpc",
						"request_config.0.mode":               "MAPPING_PASSTHROUGH",
						"request_config.0.body_format":        "STREAM",
						"http_vpc_service_config.0.name":      name,
						"http_vpc_service_config.0.method":    "POST",
						"http_vpc_service_config.0.path":      "/web/cloudapi/vpc",
						"http_vpc_service_config.0.timeout":   "20",
						"http_vpc_service_config.0.aone_name": "cloudapi-openapi",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_vpc(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence_vpc)

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
					"name":        "${alicloud_api_gateway_group.default.name}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api description",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/vpc",
						"mode":     "MAPPING",
					}},
					"service_type": "HTTP-VPC",
					"http_vpc_service_config": []map[string]string{{
						"name":                  "${alicloud_api_gateway_vpc_access.default.name}",
						"method":                "GET",
						"path":                  "/web/cloudapi/vpc",
						"timeout":               "20",
						"aone_name":             "cloudapi-openapi",
						"vpc_scheme":            "https",
						"content_type_category": "DEFAULT",
						"content_type_value":    "application/x-www-form-urlencoded; charset=UTF-8",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                 name,
						"http_vpc_service_config.0.name":       name,
						"http_vpc_service_config.0.method":     "GET",
						"http_vpc_service_config.0.path":       "/web/cloudapi/vpc",
						"http_vpc_service_config.0.timeout":    "20",
						"http_vpc_service_config.0.aone_name":  "cloudapi-openapi",
						"http_vpc_service_config.0.vpc_scheme": "https",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"http_vpc_service_config": []map[string]string{{
						"name":                  "${alicloud_api_gateway_vpc_access.update.name}",
						"method":                "HEAD",
						"path":                  "/web/cloudapi/vpc/update",
						"timeout":               "30",
						"aone_name":             "cloudapi-task",
						"vpc_scheme":            "http",
						"content_type_category": "CUSTOM",
						"content_type_value":    "application/json",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_vpc_service_config.0.name":                  name + "_update",
						"http_vpc_service_config.0.method":                "HEAD",
						"http_vpc_service_config.0.path":                  "/web/cloudapi/vpc/update",
						"http_vpc_service_config.0.timeout":               "30",
						"http_vpc_service_config.0.aone_name":             "cloudapi-task",
						"http_vpc_service_config.0.vpc_scheme":            "http",
						"http_vpc_service_config.0.content_type_category": "CUSTOM",
						"http_vpc_service_config.0.content_type_value":    "application/json",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_fc(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${alicloud_api_gateway_group.default.name}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api description",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/vpc",
						"mode":     "MAPPING",
					}},
					"service_type": "FunctionCompute",
					"fc_service_config": []map[string]string{{
						"function_version": "2.0",
						"function_type":    "FCEvent",
						"qualifier":        "LATEST",
						"region":           defaultRegionToTest,
						"function_name":    name + "Func",
						"service_name":     name,
						"timeout":          "20",
						"arn_role":         "cloudapi-openapi",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                 name,
						"fc_service_config.0.function_version": "2.0",
						"fc_service_config.0.function_type":    "FCEvent",
						"fc_service_config.0.qualifier":        "LATEST",
						"fc_service_config.0.region":           defaultRegionToTest,
						"fc_service_config.0.function_name":    name + "Func",
						"fc_service_config.0.service_name":     name,
						"fc_service_config.0.timeout":          "20",
						"fc_service_config.0.arn_role":         "cloudapi-openapi",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fc_service_config": []map[string]string{{
						"function_version": "2.0",
						"function_type":    "FCEvent",
						"qualifier":        "1.0",
						"region":           "ap-southeast-1",
						"function_name":    name + "_update",
						"service_name":     name + "_update",
						"timeout":          "30",
						"arn_role":         "cloudapi-task",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                 name,
						"fc_service_config.0.function_version": "2.0",
						"fc_service_config.0.function_type":    "FCEvent",
						"fc_service_config.0.qualifier":        "1.0",
						"fc_service_config.0.region":           "ap-southeast-1",
						"fc_service_config.0.function_name":    name + "_update",
						"fc_service_config.0.service_name":     name + "_update",
						"fc_service_config.0.timeout":          "30",
						"fc_service_config.0.arn_role":         "cloudapi-task",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_fc2(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${alicloud_api_gateway_group.default.name}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api description",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/fc",
						"mode":     "MAPPING",
					}},
					"service_type": "FunctionCompute",
					"fc_service_config": []map[string]string{{
						"function_version":   "2.0",
						"function_type":      "HttpTrigger",
						"function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"path":               "/test/path/fc",
						"method":             "GET",
						"only_business_path": "false",
						"region":             defaultRegionToTest,
						"timeout":            "20",
						"arn_role":           "cloudapi-openapi",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                   name,
						"fc_service_config.0.function_version":   "2.0",
						"fc_service_config.0.function_type":      "HttpTrigger",
						"fc_service_config.0.function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"fc_service_config.0.path":               "/test/path/fc",
						"fc_service_config.0.method":             "GET",
						"fc_service_config.0.only_business_path": "false",
						"fc_service_config.0.region":             defaultRegionToTest,
						"fc_service_config.0.timeout":            "20",
						"fc_service_config.0.arn_role":           "cloudapi-openapi",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fc_service_config": []map[string]string{{
						"function_version":   "2.0",
						"function_type":      "HttpTrigger",
						"function_base_url":  "http://apigateway-update.alicloudapi.com/fcapp.run/",
						"path":               "/test/path/fc/update",
						"method":             "POST",
						"only_business_path": "true",
						"region":             defaultRegionToTest,
						"timeout":            "30",
						"arn_role":           "cloudapi-task",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                   name,
						"fc_service_config.0.function_version":   "2.0",
						"fc_service_config.0.function_type":      "HttpTrigger",
						"fc_service_config.0.function_base_url":  "http://apigateway-update.alicloudapi.com/fcapp.run/",
						"fc_service_config.0.path":               "/test/path/fc/update",
						"fc_service_config.0.method":             "POST",
						"fc_service_config.0.only_business_path": "true",
						"fc_service_config.0.region":             defaultRegionToTest,
						"fc_service_config.0.timeout":            "30",
						"fc_service_config.0.arn_role":           "cloudapi-task",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_fc3(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${alicloud_api_gateway_group.default.name}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api description",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/fc3",
						"mode":     "MAPPING",
					}},
					"service_type": "FunctionCompute",
					"fc_service_config": []map[string]string{{
						"function_version":   "3.0",
						"function_type":      "HttpTrigger",
						"function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"path":               "/test/path/fc3",
						"method":             "GET",
						"only_business_path": "false",
						"region":             defaultRegionToTest,
						"timeout":            "20",
						"arn_role":           "cloudapi-openapi",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                   name,
						"fc_service_config.0.function_version":   "3.0",
						"fc_service_config.0.function_type":      "HttpTrigger",
						"fc_service_config.0.function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"fc_service_config.0.path":               "/test/path/fc3",
						"fc_service_config.0.method":             "GET",
						"fc_service_config.0.only_business_path": "false",
						"fc_service_config.0.region":             defaultRegionToTest,
						"fc_service_config.0.timeout":            "20",
						"fc_service_config.0.arn_role":           "cloudapi-openapi",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fc_service_config": []map[string]string{{
						"function_version":   "2.0",
						"function_type":      "HttpTrigger",
						"function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"path":               "/test/path/fc2",
						"method":             "GET",
						"only_business_path": "false",
						"region":             defaultRegionToTest,
						"timeout":            "20",
						"arn_role":           "cloudapi-openapi",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                   name,
						"fc_service_config.0.function_version":   "2.0",
						"fc_service_config.0.function_type":      "HttpTrigger",
						"fc_service_config.0.function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"fc_service_config.0.path":               "/test/path/fc2",
						"fc_service_config.0.method":             "GET",
						"fc_service_config.0.only_business_path": "false",
						"fc_service_config.0.region":             defaultRegionToTest,
						"fc_service_config.0.timeout":            "20",
						"fc_service_config.0.arn_role":           "cloudapi-openapi",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_convertFc(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${alicloud_api_gateway_group.default.name}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api description",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/vpc",
						"mode":     "MAPPING",
					}},
					"service_type": "FunctionCompute",
					"fc_service_config": []map[string]string{{
						"function_version": "2.0",
						"function_type":    "FCEvent",
						"qualifier":        "LATEST",
						"region":           defaultRegionToTest,
						"function_name":    name + "Func",
						"service_name":     name,
						"timeout":          "20",
						"arn_role":         "cloudapi-openapi",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                 name,
						"fc_service_config.0.function_version": "2.0",
						"fc_service_config.0.function_type":    "FCEvent",
						"fc_service_config.0.qualifier":        "LATEST",
						"fc_service_config.0.region":           defaultRegionToTest,
						"fc_service_config.0.function_name":    name + "Func",
						"fc_service_config.0.service_name":     name,
						"fc_service_config.0.timeout":          "20",
						"fc_service_config.0.arn_role":         "cloudapi-openapi",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fc_service_config": []map[string]string{{
						"function_version":   "2.0",
						"function_type":      "HttpTrigger",
						"function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"path":               "/test/path/fc",
						"method":             "GET",
						"only_business_path": "false",
						"region":             defaultRegionToTest,
						"timeout":            "20",
						"arn_role":           "cloudapi-openapi",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                   name,
						"fc_service_config.0.function_version":   "2.0",
						"fc_service_config.0.service_name":       "",
						"fc_service_config.0.function_name":      "",
						"fc_service_config.0.qualifier":          "",
						"fc_service_config.0.function_type":      "HttpTrigger",
						"fc_service_config.0.function_base_url":  "http://apigateway-backend.alicloudapi.com/fcapp.run/",
						"fc_service_config.0.path":               "/test/path/fc",
						"fc_service_config.0.method":             "GET",
						"fc_service_config.0.only_business_path": "false",
						"fc_service_config.0.region":             defaultRegionToTest,
						"fc_service_config.0.timeout":            "20",
						"fc_service_config.0.arn_role":           "cloudapi-openapi",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_multi(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default.9"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${alicloud_api_gateway_group.default.name}" + "${count.index}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api description",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/${count.index}",
						"mode":     "MAPPING",
					}},
					"service_type": "HTTP",
					"http_service_config": []map[string]string{{
						"address":   "http://apigateway-backend.alicloudapi.com:8080",
						"method":    "GET",
						"path":      "/web/cloudapi/${count.index}",
						"timeout":   "20",
						"aone_name": "cloudapi-openapi",
					}},
					"request_parameters": []map[string]string{{
						"name":         "testparam",
						"type":         "STRING",
						"required":     "OPTIONAL",
						"in":           "QUERY",
						"in_service":   "QUERY",
						"name_service": "testparams",
					}},
					"count": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAliCloudApigatewayApi_mock(t *testing.T) {
	var api *cloudapi.DescribeApiResponse
	resourceId := "alicloud_api_gateway_api.default"
	ra := resourceAttrInit(resourceId, apiGatewayApiMap)
	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &api, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApiGatewayApi_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayApiConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${alicloud_api_gateway_group.default.name}",
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"description": "tf_testAcc_api mock",
					"auth_type":   "APP",
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/fc",
						"mode":     "MAPPING",
					}},
					"service_type": "MOCK",
					"mock_service_config": []map[string]string{{
						"result":    "success",
						"aone_name": "cloudapi-openapi",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                            name,
						"description":                     "tf_testAcc_api mock",
						"mock_service_config.0.result":    "success",
						"mock_service_config.0.aone_name": "cloudapi-openapi",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mock_service_config": []map[string]string{{
						"result":    "OK",
						"aone_name": "cloudapi-task",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                            name,
						"mock_service_config.0.result":    "OK",
						"mock_service_config.0.aone_name": "cloudapi-task",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func checkRequestParametersAttr(resourceID, key, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceID]
		if !ok {
			return fmt.Errorf("not found: %s", resourceID)
		}

		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, key) && v != expected {
				return fmt.Errorf("%s: Attribute '%s' expected '%s', got '%s'", resourceID, key, expected, v)
			}
		}

		return nil
	}
}

func resourceApigatewayApiConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
	  default = "%s"
	}

	variable "apigateway_group_description_test" {
	  default = "tf_testAcc_api group description"
	}
	
	resource "alicloud_api_gateway_group" "default" {
	  name = "${var.name}"
	  description = "${var.apigateway_group_description_test}"
	}
	`, name)
}

func resourceApigatewayApiConfigDependence_vpc(name string) string {
	return fmt.Sprintf(`

	variable "name" {
	  default = "%s"
	}
	resource "alicloud_api_gateway_group" "default" {
	  name = "${var.name}"
	  description = "tf_testAcc_api group description"
	}

	resource "alicloud_api_gateway_vpc_access" "default" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	  instance_id = "${alicloud_instance.default.id}"
	  port = "8080"
	}

    resource "alicloud_api_gateway_vpc_access" "update" {
	  name = "${var.name}_update"
	  vpc_id = "${alicloud_vpc.default.id}"
	  instance_id = "${alicloud_instance.default.id}"
	  port = "8848"
	}
	%s
	
	`, name, ApigatewayVpcAccessConfigDependence)
}

var apiGatewayApiMap = map[string]string{
	"name":                      CHECKSET,
	"group_id":                  CHECKSET,
	"description":               "tf_testAcc_api description",
	"auth_type":                 "APP",
	"request_config.0.protocol": "HTTP",
	"request_config.0.method":   "GET",
	"request_config.0.path":     CHECKSET,
	"request_config.0.mode":     "MAPPING",
	"service_type":              CHECKSET,
	"api_id":                    CHECKSET,
}
