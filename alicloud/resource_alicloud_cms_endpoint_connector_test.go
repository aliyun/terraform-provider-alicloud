package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms EndpointConnector. >>> Resource test cases, hand-written.
func TestAccAliCloudCmsEndpointConnector_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_endpoint_connector.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEndpointConnectorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEndpointConnector")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEndpointConnectorBasicDependence)
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
					"workspace":   "${alicloud_cms_workspace.default.workspace_name}",
					"type":        "model_service",
					"name":        name,
					"endpoint":    "https://example.com/api",
					"alias":       name,
					"description": name,
					"credential":  map[string]interface{}{"token": fmt.Sprintf("secret-%s", name)},
					"headers": []map[string]interface{}{
						{"key": "X-Custom-Header", "value": "foo"},
					},
					"properties": map[string]interface{}{"modelId": "model-001"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace":    CHECKSET,
						"type":         "model_service",
						"name":         name,
						"endpoint":     "https://example.com/api",
						"alias":        name,
						"description":  name,
						"connector_id": CHECKSET,
						"region_id":    CHECKSET,
						"created_at":   CHECKSET,
						"updated_at":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":   "${alicloud_cms_workspace.default.workspace_name}",
					"type":        "model_service",
					"name":        fmt.Sprintf("%s-v2", name),
					"endpoint":    "https://example.com/api/v2",
					"alias":       fmt.Sprintf("%s-updated", name),
					"description": fmt.Sprintf("%s-updated", name),
					"credential":  map[string]interface{}{"token": fmt.Sprintf("secret-%s-v2", name)},
					"headers": []map[string]interface{}{
						{"key": "X-Other-Header", "value": "bar"},
					},
					"properties": map[string]interface{}{"modelId": "model-002"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint":    "https://example.com/api/v2",
						"alias":       fmt.Sprintf("%s-updated", name),
						"description": fmt.Sprintf("%s-updated", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":   "${alicloud_cms_workspace.default.workspace_name}",
					"type":        "model_service",
					"name":        name,
					"endpoint":    "https://example.com/api/v2",
					"description": "",
					"alias":       "",
					"credential":  map[string]interface{}{"token": fmt.Sprintf("secret-%s-v2", name)},
					"properties":  map[string]interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "",
						"alias":       "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential"},
			},
		},
	})
}

func TestAccAliCloudCmsEndpointConnector_agentApp(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_endpoint_connector.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEndpointConnectorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEndpointConnector")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecagent%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEndpointConnectorBasicDependence)
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
					"workspace":  "${alicloud_cms_workspace.default.workspace_name}",
					"type":       "agent_app",
					"name":       name,
					"endpoint":   "https://agent.example.com",
					"credential": map[string]interface{}{"apiKey": fmt.Sprintf("ak-%s", name)},
					"properties": map[string]interface{}{"appId": "app-001"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "agent_app",
						"name":     name,
						"endpoint": "https://agent.example.com",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential"},
			},
		},
	})
}

func TestAccAliCloudCmsEndpointConnector_disappear(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_endpoint_connector.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsEndpointConnectorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEndpointConnector")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecdisp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsEndpointConnectorBasicDependence)
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
					"workspace":  "${alicloud_cms_workspace.default.workspace_name}",
					"type":       "model_service",
					"name":       name,
					"endpoint":   "https://example.com/api",
					"credential": map[string]interface{}{"token": fmt.Sprintf("secret-%s", name)},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":     "model_service",
						"name":     name,
						"endpoint": "https://example.com/api",
					}),
				),
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				ResourceName:       resourceId,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

var AliCloudCmsEndpointConnectorMap = map[string]string{
	"connector_id": CHECKSET,
	"region_id":    CHECKSET,
	"created_at":   CHECKSET,
	"updated_at":   CHECKSET,
}

func AliCloudCmsEndpointConnectorBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}
`, name)
}
