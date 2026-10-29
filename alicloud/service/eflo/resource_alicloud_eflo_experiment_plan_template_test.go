package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass template_pipeline -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Eflo ExperimentPlanTemplate. >>> Resource test cases, automatically generated.
// Case 实验计划模版用例_V2_预发 10580
func TestAccAliCloudEfloExperimentPlanTemplate_basic10580(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_experiment_plan_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEfloExperimentPlanTemplateMap10580)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloExperimentPlanTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEfloExperimentPlanTemplateBasicDependence10580)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EfloExperimentPlanTemplateSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"template_pipeline": []map[string]interface{}{
						{
							"workload_id":   "2",
							"workload_name": "MatMul",
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":    "90",
									"gpu_per_worker":    "8",
									"memory_per_worker": "500",
									"share_memory":      "500",
									"worker_num":        "1",
								},
							},
							"pipeline_order": "1",
							"scene":          "baseline",
						},
					},
					"privacy_level": "private",
					"template_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_pipeline.#": "1",
						"privacy_level":       "private",
						"template_name":       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_pipeline": []map[string]interface{}{
						{
							"workload_id":   "1",
							"workload_name": "Bert-base",
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":    "200",
									"gpu_per_worker":    "200",
									"memory_per_worker": "200",
									"share_memory":      "200",
									"worker_num":        "2",
								},
							},
							"pipeline_order": "2",
							"scene":          "baseline2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_pipeline.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_pipeline": []map[string]interface{}{
						{
							"workload_id":   "1",
							"workload_name": "Bert-base",
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":     "200",
									"gpu_per_worker":     "200",
									"memory_per_worker":  "200",
									"share_memory":       "200",
									"worker_num":         "2",
									"py_torch_version":   "1",
									"gpu_driver_version": "1",
									"cuda_version":       "1",
									"nccl_version":       "1",
								},
							},
							"pipeline_order": "2",
							"scene":          "baseline2",
							"setting_params": map[string]interface{}{
								"ITERATION":         "200",
								"NUMERICAL_FORMATS": "FP16",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_pipeline.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_pipeline": []map[string]interface{}{
						{
							"workload_id":   "2",
							"workload_name": "MatMul",
							"scene":         "baseline",
							"setting_params": map[string]interface{}{
								"ITERATION":         "100",
								"NUMERICAL_FORMATS": "FP16",
							},
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":     "100",
									"gpu_per_worker":     "100",
									"memory_per_worker":  "100",
									"share_memory":       "100",
									"worker_num":         "1",
									"py_torch_version":   "1",
									"gpu_driver_version": "1",
									"cuda_version":       "1",
									"nccl_version":       "1",
								},
							},
							"pipeline_order": "1",
						},
						{
							"workload_id":   "1",
							"workload_name": "Bert-base",
							"scene":         "baseline2",
							"setting_params": map[string]interface{}{
								"ITERATION":         "200",
								"NUMERICAL_FORMATS": "FP16",
							},
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":     "200",
									"gpu_per_worker":     "200",
									"memory_per_worker":  "200",
									"share_memory":       "200",
									"worker_num":         "2",
									"py_torch_version":   "1",
									"gpu_driver_version": "2",
									"cuda_version":       "2",
									"nccl_version":       "2",
								},
							},
							"pipeline_order": "2",
						},
						{
							"workload_id":   "6",
							"workload_name": "Unet",
							"scene":         "baseline",
							"setting_params": map[string]interface{}{
								"ITERATION":         "300",
								"NUMERICAL_FORMATS": "FP16",
							},
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":     "300",
									"gpu_per_worker":     "300",
									"memory_per_worker":  "300",
									"share_memory":       "300",
									"worker_num":         "3",
									"py_torch_version":   "3",
									"gpu_driver_version": "3",
									"cuda_version":       "3",
									"nccl_version":       "3",
								},
							},
							"pipeline_order": "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_pipeline.#": "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_pipeline"},
			},
		},
	})
}

func TestAccAliCloudEfloExperimentPlanTemplate_basic10580_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_experiment_plan_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEfloExperimentPlanTemplateMap10580)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloExperimentPlanTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEfloExperimentPlanTemplateBasicDependence10580)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.EfloExperimentPlanTemplateSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"template_pipeline": []map[string]interface{}{
						{
							"workload_id":   "2",
							"workload_name": "MatMul",
							"scene":         "baseline",
							"setting_params": map[string]interface{}{
								"ITERATION":         "100",
								"NUMERICAL_FORMATS": "FP16",
							},
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":     "100",
									"gpu_per_worker":     "100",
									"memory_per_worker":  "100",
									"share_memory":       "100",
									"worker_num":         "1",
									"py_torch_version":   "1",
									"gpu_driver_version": "1",
									"cuda_version":       "1",
									"nccl_version":       "1",
								},
							},
							"pipeline_order": "1",
						},
						{
							"workload_id":   "1",
							"workload_name": "Bert-base",
							"scene":         "baseline2",
							"setting_params": map[string]interface{}{
								"ITERATION":         "200",
								"NUMERICAL_FORMATS": "FP16",
							},
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":     "200",
									"gpu_per_worker":     "200",
									"memory_per_worker":  "200",
									"share_memory":       "200",
									"worker_num":         "2",
									"py_torch_version":   "1",
									"gpu_driver_version": "2",
									"cuda_version":       "2",
									"nccl_version":       "2",
								},
							},
							"pipeline_order": "2",
						},
						{
							"workload_id":   "6",
							"workload_name": "Unet",
							"scene":         "baseline",
							"setting_params": map[string]interface{}{
								"ITERATION":         "300",
								"NUMERICAL_FORMATS": "FP16",
							},
							"env_params": []map[string]interface{}{
								{
									"cpu_per_worker":     "300",
									"gpu_per_worker":     "300",
									"memory_per_worker":  "300",
									"share_memory":       "300",
									"worker_num":         "3",
									"py_torch_version":   "3",
									"gpu_driver_version": "3",
									"cuda_version":       "3",
									"nccl_version":       "3",
								},
							},
							"pipeline_order": "3",
						},
					},
					"privacy_level":        "private",
					"template_name":        name,
					"template_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_pipeline.#":  "3",
						"privacy_level":        "private",
						"template_name":        name,
						"template_description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_pipeline"},
			},
		},
	})
}

var AliCloudEfloExperimentPlanTemplateMap10580 = map[string]string{
	"create_time": acctest.CHECKSET,
	"template_id": acctest.CHECKSET,
}

func AliCloudEfloExperimentPlanTemplateBasicDependence10580(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test Eflo ExperimentPlanTemplate. <<< Resource test cases, automatically generated.
