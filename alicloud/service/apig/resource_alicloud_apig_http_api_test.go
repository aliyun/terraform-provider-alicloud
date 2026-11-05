package apig_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	apig "github.com/aliyun/terraform-provider-alicloud/alicloud/service/apig"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Apig HttpApi. >>> Resource test cases, automatically generated.
// Case http_api_crud_test 12898
func TestAccAliCloudApigHttpApi_basic12898(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigHttpApiMap12898)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigHttpApi")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigHttpApiBasicDependence12898)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_name": name,
					"protocols": []string{
						"HTTP"},
					"type":              "Rest",
					"description":       "test description for cspec",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"ai_protocols": []string{
						"OpenAI"},
					"base_path":      "/cspec-test",
					"enable_auth":    "false",
					"model_category": "LLM",
					"deploy_configs": []string{
						"{\\\"autoDeploy\\\":false}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name":     name,
						"protocols.#":       "1",
						"type":              "Rest",
						"description":       "test description for cspec",
						"resource_group_id": acctest.CHECKSET,
						"ai_protocols.#":    "1",
						"base_path":         "/cspec-test",
						"enable_auth":       "false",
						"model_category":    "LLM",
						"deploy_configs.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "updated description for cspec",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "updated description for cspec",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"HTTPS"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"HTTP", "HTTPS"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"HTTP"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ai_protocols", "base_path", "deploy_configs", "model_category"},
			},
		},
	})
}

var AlicloudApigHttpApiMap12898 = map[string]string{}

func AlicloudApigHttpApiBasicDependence12898(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

`, name)
}

// Case HttpApi测试 6880
func TestAccAliCloudApigHttpApi_basic6880(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigHttpApiMap6880)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigHttpApi")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigHttpApiBasicDependence6880)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":   "/v1",
					"description": "zhiwei_pop_testcase",
					"type":        "Rest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name": name,
						"protocols.#":   "1",
						"base_path":     "/v1",
						"description":   "zhiwei_pop_testcase",
						"type":          "Rest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":   "/v2",
					"description": "1783585689",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "2",
						"base_path":   "/v2",
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ai_protocols", "base_path", "deploy_configs", "enable_auth", "model_category"},
			},
		},
	})
}

var AlicloudApigHttpApiMap6880 = map[string]string{}

func AlicloudApigHttpApiBasicDependence6880(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}


`, name)
}

// Case HttpApi测试_CCApi 9287
func TestAccAliCloudApigHttpApi_basic9287(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigHttpApiMap9287)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigHttpApi")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigHttpApiBasicDependence9287)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":         "/v1",
					"description":       "zhiwei_pop_testcase",
					"type":              "Rest",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name":     name,
						"protocols.#":       "1",
						"base_path":         "/v1",
						"description":       "zhiwei_pop_testcase",
						"type":              "Rest",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":         "/v2",
					"description":       "1783585689",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#":       "2",
						"base_path":         "/v2",
						"description":       acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ai_protocols", "base_path", "deploy_configs", "enable_auth", "model_category"},
			},
		},
	})
}

var AlicloudApigHttpApiMap9287 = map[string]string{}

func AlicloudApigHttpApiBasicDependence9287(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case HttpApi测试_CCApi_2 9288
func TestAccAliCloudApigHttpApi_basic9288(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigHttpApiMap9288)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigHttpApi")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigHttpApiBasicDependence9288)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":         "/v1",
					"description":       "zhiwei_pop_testcase",
					"type":              "Rest",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name":     name,
						"protocols.#":       "1",
						"base_path":         "/v1",
						"description":       "zhiwei_pop_testcase",
						"type":              "Rest",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":         "/v2",
					"description":       "1783585689",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#":       "2",
						"base_path":         "/v2",
						"description":       acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ai_protocols", "base_path", "deploy_configs", "enable_auth", "model_category"},
			},
		},
	})
}

var AlicloudApigHttpApiMap9288 = map[string]string{}

func AlicloudApigHttpApiBasicDependence9288(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case HttpApi测试_副本1732026511257 9021
func TestAccAliCloudApigHttpApi_basic9021(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigHttpApiMap9021)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigHttpApi")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigHttpApiBasicDependence9021)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_name": name,
					"protocols": []string{
						"${var.protocol}"},
					"base_path":         "/v1",
					"description":       "zhiwei_pop_testcase",
					"type":              "Rest",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name":     name,
						"protocols.#":       "1",
						"base_path":         "/v1",
						"description":       "zhiwei_pop_testcase",
						"type":              "Rest",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol_https}", "${var.protocol}"},
					"base_path":         "/v2",
					"description":       "1783585689",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#":       "2",
						"base_path":         "/v2",
						"description":       acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{
						"${var.protocol}"},
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"base_path":   "/v1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ai_protocols", "base_path", "deploy_configs", "enable_auth", "model_category"},
			},
		},
	})
}

var AlicloudApigHttpApiMap9021 = map[string]string{}

func AlicloudApigHttpApiBasicDependence9021(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "protocol" {
  default = "HTTP"
}

variable "protocol_https" {
  default = "HTTPS"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Apig HttpApi. <<< Resource test cases, automatically generated.

// Case minimal attribute set: create with only the fields the API requires
// (http_api_name + protocols + type, plus base_path which is required when type=Rest),
// to verify that the optional fields have no server-side default that would produce a
// perpetual diff.
func TestAccAliCloudApigHttpApi_minimal(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigHttpApiMapMinimal)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigHttpApi")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigHttpApiBasicDependenceMinimal)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_name": name,
					"protocols": []string{
						"HTTP"},
					"type":      "Rest",
					"base_path": "/v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"http_api_name": name,
						"protocols.#":   "1",
						"type":          "Rest",
						"base_path":     "/v1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ai_protocols", "base_path", "deploy_configs", "model_category"},
			},
		},
	})
}

var AlicloudApigHttpApiMapMinimal = map[string]string{}

func AlicloudApigHttpApiBasicDependenceMinimal(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Convergent member-order coverage for the configurable TypeList attributes
// (protocols, deploy_configs). Each step reorders exactly one list over the
// same resource; the post-apply plan must stay empty. ai_protocols is NOT
// converged here: resourceHttpApiCreate sends aiProtocols, but
// resourceHttpApiUpdate does not process ai_protocols and resourceHttpApiRead
// does not d.Set it, so server-side member-order round-trip is unverifiable;
// see the collection-order: bypass ai_protocols stock marker below.
func TestAccAliCloudApigHttpApi_memberOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_http_api.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigHttpApiMapMemberOrder)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigHttpApi")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigHttpApiBasicDependenceMinimal)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"http_api_name": name,
					"type":          "Rest",
					"base_path":     "/morder",
					"protocols":     []string{"HTTP", "HTTPS"},
					"ai_protocols":  []string{"OpenAI", "Anthropic"},
					"deploy_configs": []string{
						"{\\\"autoDeploy\\\":false}",
						"{\\\"autoDeploy\\\":true}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#":      "2",
						"ai_protocols.#":   "2",
						"deploy_configs.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{"HTTPS", "HTTP"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deploy_configs": []string{
						"{\\\"autoDeploy\\\":true}",
						"{\\\"autoDeploy\\\":false}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deploy_configs.#": "2",
					}),
				),
			},
		},
	})
}

// collection-order: bypass ai_protocols -- pre-existing attribute outside this change: resourceHttpApiCreate sends aiProtocols (resource_alicloud_apig_http_api.go:116) but resourceHttpApiUpdate does not process ai_protocols/aiProtocols and resourceHttpApiRead does not d.Set it, so member-order round-trip is unverifiable; the previously-added reorder Step only asserted ai_protocols.#=2 without sending the update to the server. Per-property adjudication, not a uniform historical bypass.

var AlicloudApigHttpApiMapMemberOrder = map[string]string{}
