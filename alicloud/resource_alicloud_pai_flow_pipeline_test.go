package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiFlow Pipeline. >>> Resource test cases, automatically generated.
// Case Case02 10139
func TestAccAliCloudPaiFlowPipeline_basic10139(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_flow_pipeline.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiFlowPipelineMap10139)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiFlowServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiFlowPipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := randIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccpaiflow%d", rand)
	uuid := fmt.Sprintf("terraformuuid%v", rand)
	uuidUpdate := fmt.Sprintf("tfuuidupdatee%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiFlowPipelineBasicDependence10139)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_id": "${alicloud_pai_workspace_workspace.defaultWjQD1e.id}",
					"manifest":     `apiVersion: \"core/v1\"\nmetadata:\n  provider: \"` + "${data.alicloud_account.default.id}" + `\"\n  version: \"v1\"\n  identifier: \"my_pipeline\"\n  name: \"source-transform\"\n  uuid: \"` + uuid + `\"\n  annotations: {}\n  labels: {}\nspec:\n  inputs:\n    artifacts: []\n    parameters:\n    - name: \"execution_maxcompute\"\n      type: \"Map\"\n      value:\n        spec:\n          endpoint: \"http://service.cn.maxcompute.aliyun-inc.com/api\"\n          odpsProject: \"test_i****\"\n  outputs:\n    artifacts: []\n    parameters: []\n  arguments:\n    artifacts: []\n    parameters: []\n  dependencies: []\n  initContainers: []\n  sideCarContainers: []\n  pipelines:\n  - apiVersion: \"core/v1\"\n    metadata:\n      provider: \"pai\"\n      version: \"v1\"\n      identifier: \"data_source\"\n      name: \"data-source\"\n      uuid: \"2ftahdnzcod2rt6u9q\"\n      displayName: \"读数据表-1\"\n      annotations: {}\n      labels: {}\n    spec:\n      inputs:\n        artifacts: []\n        parameters: []\n      outputs:\n        artifacts: []\n        parameters: []\n      arguments:\n        artifacts: []\n        parameters:\n        - name: \"inputTableName\"\n          value: \"pai_online_project.wumai_data\"\n        - name: \"execution\"\n          from: \"{{inputs.parameters.execution_maxcompute}}\"\n      dependencies: []\n      initContainers: []\n      sideCarContainers: []\n      pipelines: []\n      volumes: []\n  - apiVersion: \"core/v1\"\n    metadata:\n      provider: \"pai\"\n      version: \"v1\"\n      identifier: \"type_transform\"\n      name: \"type-transform\"\n      uuid: \"gacnnnl4ksvbabfh6l\"\n      displayName: \"类型转换-1\"\n      annotations: {}\n      labels: {}\n    spec:\n      inputs:\n        artifacts: []\n        parameters: []\n      outputs:\n        artifacts: []\n        parameters: []\n      arguments:\n        artifacts:\n        - name: \"inputTable\"\n          from: \"{{pipelines.data_source.outputs.artifacts.outputTable}}\"\n        parameters:\n        - name: \"cols_to_double\"\n          value: \"time,hour,pm2,pm10,so2,co,no2\"\n        - name: \"execution\"\n          from: \"{{inputs.parameters.execution_maxcompute}}\"\n      dependencies:\n      - \"data_source\"\n      initContainers: []\n      sideCarContainers: []\n      pipelines: []\n      volumes: []\n  volumes: []`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_id": CHECKSET,
						"manifest":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"manifest": `apiVersion: \"core/v1\"\nmetadata:\n  provider: \"` + "${data.alicloud_account.default.id}" + `\"\n  version: \"v1\"\n  identifier: \"my_pipeline\"\n  name: \"source-transform\"\n  uuid: \"` + uuidUpdate + `\"\n  annotations: {}\n  labels: {}\nspec:\n  inputs:\n    artifacts: []\n    parameters:\n    - name: \"execution_maxcompute\"\n      type: \"Map\"\n      value:\n        spec:\n          endpoint: \"http://service.cn.maxcompute.aliyun-inc.com/api/v2\"\n          odpsProject: \"test_i****\"\n  outputs:\n    artifacts: []\n    parameters: []\n  arguments:\n    artifacts: []\n    parameters: []\n  dependencies: []\n  initContainers: []\n  sideCarContainers: []\n  pipelines:\n  - apiVersion: \"core/v1\"\n    metadata:\n      provider: \"pai\"\n      version: \"v1\"\n      identifier: \"data_source\"\n      name: \"data-source\"\n      uuid: \"2ftahdnzcod2rt6u9q\"\n      displayName: \"读数据表-1\"\n      annotations: {}\n      labels: {}\n    spec:\n      inputs:\n        artifacts: []\n        parameters: []\n      outputs:\n        artifacts: []\n        parameters: []\n      arguments:\n        artifacts: []\n        parameters:\n        - name: \"inputTableName\"\n          value: \"pai_online_project.wumai_data\"\n        - name: \"execution\"\n          from: \"{{inputs.parameters.execution_maxcompute}}\"\n      dependencies: []\n      initContainers: []\n      sideCarContainers: []\n      pipelines: []\n      volumes: []\n  - apiVersion: \"core/v1\"\n    metadata:\n      provider: \"pai\"\n      version: \"v1\"\n      identifier: \"type_transform\"\n      name: \"type-transform\"\n      uuid: \"gacnnnl4ksvbabfh6l\"\n      displayName: \"类型转换-1\"\n      annotations: {}\n      labels: {}\n    spec:\n      inputs:\n        artifacts: []\n        parameters: []\n      outputs:\n        artifacts: []\n        parameters: []\n      arguments:\n        artifacts:\n        - name: \"inputTable\"\n          from: \"{{pipelines.data_source.outputs.artifacts.outputTable}}\"\n        parameters:\n        - name: \"cols_to_double\"\n          value: \"time,hour,pm2,pm10,so2,co,no2\"\n        - name: \"execution\"\n          from: \"{{inputs.parameters.execution_maxcompute}}\"\n      dependencies:\n      - \"data_source\"\n      initContainers: []\n      sideCarContainers: []\n      pipelines: []\n      volumes: []\n  volumes: []`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"manifest": CHECKSET,
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

func TestAccAliCloudPaiFlowPipeline_basic10139_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_flow_pipeline.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiFlowPipelineMap10139)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiFlowServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiFlowPipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := randIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccpaiflow%d", rand)
	uuid := fmt.Sprintf("terraformuuid%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiFlowPipelineBasicDependence10139)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace_id": "${alicloud_pai_workspace_workspace.defaultWjQD1e.id}",
					"manifest":     `apiVersion: \"core/v1\"\nmetadata:\n  provider: \"` + "${data.alicloud_account.default.id}" + `\"\n  version: \"v1\"\n  identifier: \"my_pipeline\"\n  name: \"source-transform\"\n  uuid: \"` + uuid + `\"\n  annotations: {}\n  labels: {}\nspec:\n  inputs:\n    artifacts: []\n    parameters:\n    - name: \"execution_maxcompute\"\n      type: \"Map\"\n      value:\n        spec:\n          endpoint: \"http://service.cn.maxcompute.aliyun-inc.com/api\"\n          odpsProject: \"test_i****\"\n  outputs:\n    artifacts: []\n    parameters: []\n  arguments:\n    artifacts: []\n    parameters: []\n  dependencies: []\n  initContainers: []\n  sideCarContainers: []\n  pipelines:\n  - apiVersion: \"core/v1\"\n    metadata:\n      provider: \"pai\"\n      version: \"v1\"\n      identifier: \"data_source\"\n      name: \"data-source\"\n      uuid: \"2ftahdnzcod2rt6u9q\"\n      displayName: \"读数据表-1\"\n      annotations: {}\n      labels: {}\n    spec:\n      inputs:\n        artifacts: []\n        parameters: []\n      outputs:\n        artifacts: []\n        parameters: []\n      arguments:\n        artifacts: []\n        parameters:\n        - name: \"inputTableName\"\n          value: \"pai_online_project.wumai_data\"\n        - name: \"execution\"\n          from: \"{{inputs.parameters.execution_maxcompute}}\"\n      dependencies: []\n      initContainers: []\n      sideCarContainers: []\n      pipelines: []\n      volumes: []\n  - apiVersion: \"core/v1\"\n    metadata:\n      provider: \"pai\"\n      version: \"v1\"\n      identifier: \"type_transform\"\n      name: \"type-transform\"\n      uuid: \"gacnnnl4ksvbabfh6l\"\n      displayName: \"类型转换-1\"\n      annotations: {}\n      labels: {}\n    spec:\n      inputs:\n        artifacts: []\n        parameters: []\n      outputs:\n        artifacts: []\n        parameters: []\n      arguments:\n        artifacts:\n        - name: \"inputTable\"\n          from: \"{{pipelines.data_source.outputs.artifacts.outputTable}}\"\n        parameters:\n        - name: \"cols_to_double\"\n          value: \"time,hour,pm2,pm10,so2,co,no2\"\n        - name: \"execution\"\n          from: \"{{inputs.parameters.execution_maxcompute}}\"\n      dependencies:\n      - \"data_source\"\n      initContainers: []\n      sideCarContainers: []\n      pipelines: []\n      volumes: []\n  volumes: []`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"workspace_id": CHECKSET,
						"manifest":     CHECKSET,
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

var AliCloudPaiFlowPipelineMap10139 = map[string]string{
	"create_time": CHECKSET,
}

func AliCloudPaiFlowPipelineBasicDependence10139(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_pai_workspace_workspace" "defaultWjQD1e" {
  description    = "paiflow resource record test"
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["dev"]
}
`, name)
}

// Test PaiFlow Pipeline. <<< Resource test cases, automatically generated.
