package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
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

func TestAccAlicloudApigatewayApi_basic(t *testing.T) {
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
						"name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
			{
				Config: testAccConfig(map[string]interface{}{
					"request_config": []map[string]string{{
						"protocol": "HTTP",
						"method":   "GET",
						"path":     "/test/path/test",
						"mode":     "MAPPING",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"request_config.0.path": "/test/path/test",
					}),
				),
			},
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
						"description":                "tf_testAcc_api description",
						"request_config.0.path":      "/test/path",
						"http_service_config.0.path": "/web/cloudapi",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudApigatewayApi_multi(t *testing.T) {
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

var apiGatewayApiMap = map[string]string{
	"name":                            CHECKSET,
	"group_id":                        CHECKSET,
	"description":                     "tf_testAcc_api description",
	"auth_type":                       "APP",
	"request_config.0.protocol":       "HTTP",
	"request_config.0.method":         "GET",
	"request_config.0.path":           CHECKSET,
	"request_config.0.mode":           "MAPPING",
	"service_type":                    "HTTP",
	"http_service_config.0.address":   "http://apigateway-backend.alicloudapi.com:8080",
	"http_service_config.0.method":    "GET",
	"http_service_config.0.path":      CHECKSET,
	"http_service_config.0.timeout":   "20",
	"http_service_config.0.aone_name": "cloudapi-openapi",
	"api_id":                          CHECKSET,
}
