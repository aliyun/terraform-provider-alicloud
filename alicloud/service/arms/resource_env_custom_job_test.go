package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Arms EnvCustomJob. >>> Resource test cases, automatically generated.
// Case 4554
func TestAccAliCloudArmsEnvCustomJob_basic4554(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_custom_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsEnvCustomJobMap4554)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsEnvCustomJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvcustomjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsEnvCustomJobBasicDependence4554)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id":      "${alicloud_arms_environment.env-cs.id}",
					"env_custom_job_name": name,
					"config_yaml":         `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 30s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id":      acctest.CHECKSET,
						"env_custom_job_name": name,
						"config_yaml":         acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "run",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "run",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 30s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "stop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "stop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "run",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "run",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 31s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":              "run",
					"environment_id":      "${alicloud_arms_environment.env-cs.id}",
					"env_custom_job_name": name + "-update",
					"config_yaml":         `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 30s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":              "run",
						"environment_id":      acctest.CHECKSET,
						"env_custom_job_name": name + "-update",
						"config_yaml":         acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

var AlicloudArmsEnvCustomJobMap4554 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudArmsEnvCustomJobBasicDependence4554(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "^Default"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [data.alicloud_vswitches.default.ids.0]
  new_nat_gateway      = false
  pod_cidr             = "10.124.0.0/16"
  service_cidr         = "192.168.0.0/16"
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}

resource "alicloud_arms_environment" "env-cs" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = local.cluster_id
  environment_sub_type = "ManagedKubernetes"
}


`, name)
}

// Case 4605
func TestAccAliCloudArmsEnvCustomJob_basic4605(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_custom_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsEnvCustomJobMap4605)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsEnvCustomJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvcustomjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsEnvCustomJobBasicDependence4605)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id":      "${alicloud_arms_environment.env-ecs.id}",
					"env_custom_job_name": name,
					"config_yaml":         `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 30s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id":      acctest.CHECKSET,
						"env_custom_job_name": name,
						"config_yaml":         acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "run",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "run",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 30s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "stop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "stop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "run",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "run",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 31s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":              "run",
					"environment_id":      "${alicloud_arms_environment.env-ecs.id}",
					"env_custom_job_name": name + "-update",
					"config_yaml":         `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 30s`,
					"aliyun_lang":         "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":              "run",
						"environment_id":      acctest.CHECKSET,
						"env_custom_job_name": name + "-update",
						"config_yaml":         acctest.CHECKSET,
						"aliyun_lang":         "zh",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

var AlicloudArmsEnvCustomJobMap4605 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudArmsEnvCustomJobBasicDependence4605(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

resource "alicloud_arms_environment" "env-ecs" {
  environment_type = "ECS"
  environment_name = var.name

  bind_resource_id     = data.alicloud_vpcs.default.ids.0
  environment_sub_type = "ECS"
}


`, name)
}

// Case 4554  twin
func TestAccAliCloudArmsEnvCustomJob_basic4554_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_custom_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsEnvCustomJobMap4554)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsEnvCustomJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvcustomjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsEnvCustomJobBasicDependence4554)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":              "stop",
					"environment_id":      "${alicloud_arms_environment.env-cs.id}",
					"env_custom_job_name": name,
					"config_yaml":         `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 31s`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":              "stop",
						"environment_id":      acctest.CHECKSET,
						"env_custom_job_name": name,
						"config_yaml":         acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

// Case 4605  twin
func TestAccAliCloudArmsEnvCustomJob_basic4605_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_custom_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsEnvCustomJobMap4605)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsEnvCustomJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvcustomjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsEnvCustomJobBasicDependence4605)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":              "stop",
					"environment_id":      "${alicloud_arms_environment.env-ecs.id}",
					"env_custom_job_name": name,
					"config_yaml":         `scrape_configs:\n- job_name: job-demo1\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n- job_name: job-demo2\n  honor_timestamps: false\n  honor_labels: false\n  scrape_interval: 30s\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n  http_sd_configs:\n  - url: 127.0.0.1:9090\n    refresh_interval: 31s`,
					"aliyun_lang":         "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":              "stop",
						"environment_id":      acctest.CHECKSET,
						"env_custom_job_name": name,
						"config_yaml":         acctest.CHECKSET,
						"aliyun_lang":         "en",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

// Test Arms EnvCustomJob. <<< Resource test cases, automatically generated.
