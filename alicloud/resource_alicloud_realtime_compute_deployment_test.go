// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test RealtimeCompute Deployment. >>> Resource test cases, automatically generated.
// Case 测试MAP 11896
func TestAccAliCloudRealtimeComputeDeployment_basic11896(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_deployment.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeDeploymentMap11896)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeDeploymentBasicDependence11896)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"logging_profile": "default",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "",
									"logger_level": "INFO",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "true",
									"expiration_days": "7",
								},
							},
						},
					},
					"deployment_name": name,
					"description":     "This is a test deployment.",
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"local_variables": []map[string]interface{}{
						{
							"value": "value",
							"name":  "name",
						},
					},
					"execution_mode": "STREAMING",
					"labels": map[string]interface{}{
						"\"vvp\"": "nb",
					},
					"deployment_target": []map[string]interface{}{
						{
							"mode": "PER_JOB",
							"name": "default-queue",
						},
					},
					"streaming_resource_setting": []map[string]interface{}{
						{
							"basic_resource_setting": []map[string]interface{}{
								{
									"taskmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "1Gi",
											"cpu":    "1",
										},
									},
									"parallelism": "1",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "1Gi",
											"cpu":    "1",
										},
									},
								},
							},
							"resource_setting_mode": "BASIC",
						},
					},
					"namespace": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script": "create temporary table `datagen` ( id varchar, name varchar ) with ( 'connector' = 'datagen' );  create temporary table `blackhole` ( id varchar, name varchar ) with ( 'connector' = 'blackhole' );  insert into blackhole select * from datagen;",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar"},
								},
							},
						},
					},
					"resource_id": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"flink_conf": map[string]interface{}{
						"\"execution.checkpointing.interval\"":      "180s",
						"\"execution.checkpointing.min-pause\"":     "180s",
						"\"restart-strategy\"":                      "fixed-delay",
						"\"restart-strategy.fixed-delay.attempts\"": "2147483647",
						"\"restart-strategy.fixed-delay.delay\"":    "10 s",
						"\"table.exec.state.ttl\"":                  "36 h",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name":   name,
						"description":       "This is a test deployment.",
						"engine_version":    "vvr-8.0.10-flink-1.17",
						"local_variables.#": "1",
						"execution_mode":    "STREAMING",
						"namespace":         CHECKSET,
						"resource_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"log4j2_configuration_template": "test-template",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "StdOut",
									"logger_level": "DEBUG",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "false",
									"expiration_days": "5",
								},
							},
						},
					},
					"deployment_name": name + "_update",
					"description":     "This is a test deployment 2.",
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"local_variables": []map[string]interface{}{
						{
							"value": "value1",
							"name":  "name1",
						},
						{
							"value": "value2",
							"name":  "name2",
						},
					},
					"labels": map[string]interface{}{
						"\"vvp\"": "b",
					},
					"streaming_resource_setting": []map[string]interface{}{
						{
							"basic_resource_setting": []map[string]interface{}{
								{
									"taskmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
									"parallelism": "2",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
								},
							},
							"resource_setting_mode": "BASIC",
						},
					},
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script": "CREATE TABLE result_table (id BIGINT, name STRING) WITH ( 'connector' = 'jdbc', 'url' = 'jdbc:mysql://localhost:3306/test', 'table-name' = 'result_table' );",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar", "oss://bucket-name/b.jar", "oss://bucket-name/c.jar"},
								},
							},
						},
					},
					"flink_conf": map[string]interface{}{
						"\"execution.checkpointing.interval\"":      "180s",
						"\"restart-strategy\"":                      "fixed-delay",
						"\"restart-strategy.fixed-delay.attempts\"": "2147483647",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name":   name + "_update",
						"description":       "This is a test deployment 2.",
						"engine_version":    "vvr-8.0.10-flink-1.17",
						"local_variables.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"streaming_resource_setting": []map[string]interface{}{
						{
							"resource_setting_mode": "EXPERT",
							"expert_resource_setting": []map[string]interface{}{
								{
									"resource_plan": "{\\\\n  \\\\\\\"ssgProfiles\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"name\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n      \\\\\\\"cpu\\\\\\\": 0.26,\\\\n      \\\\\\\"heap\\\\\\\": \\\\\\\"1 gb\\\\\\\",\\\\n      \\\\\\\"offHeap\\\\\\\": \\\\\\\"32 mb\\\\\\\",\\\\n      \\\\\\\"managed\\\\\\\": {},\\\\n      \\\\\\\"extended\\\\\\\": {}\\\\n    }\\\\n  ],\\\\n  \\\\\\\"nodes\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 1,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecTableSourceScan\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Source: datagen_source[7]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 2,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecCalc\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Calc[8]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 3,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"ConstraintEnforcer[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 4,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Sink: vvptest[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    }\\\\n  ],\\\\n  \\\\\\\"edges\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 1,\\\\n      \\\\\\\"target\\\\\\\": 2,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 2,\\\\n      \\\\\\\"target\\\\\\\": 3,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 3,\\\\n      \\\\\\\"target\\\\\\\": 4,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    }\\\\n  ],\\\\n  \\\\\\\"vertices\\\\\\\": {\\\\n    \\\\\\\"717c7b8afebbfb7137f6f0f99beb2a94\\\\\\\": [\\\\n      1,\\\\n      2,\\\\n      3,\\\\n      4\\\\n    ]\\\\n  }\\\\n}",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
								},
							},
						},
					},
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script":              "CREATE TABLE result_table (id BIGINT, name STRING) WITH ( 'connector' = 'jdbc', 'url' = 'jdbc:mysql://localhost:3306/test', 'table-name' = 'result_table' );",
									"additional_dependencies": []string{},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudRealtimeComputeDeploymentMap11896 = map[string]string{
	"deployment_id": CHECKSET,
}

func AlicloudRealtimeComputeDeploymentBasicDependence11896(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc.id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch.zone_id
}


`, name)
}

// Case 已重置副本 11895
func TestAccAliCloudRealtimeComputeDeployment_basic11895(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_deployment.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeDeploymentMap11895)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeDeploymentBasicDependence11895)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"logging_profile": "default",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "",
									"logger_level": "INFO",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "true",
									"expiration_days": "7",
								},
							},
						},
					},
					"deployment_name": name,
					"description":     "This is a test deployment.",
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"local_variables": []map[string]interface{}{
						{
							"value": "value",
							"name":  "name",
						},
					},
					"execution_mode": "STREAMING",
					"labels": map[string]interface{}{
						"\"vvp\"": "nb",
					},
					"deployment_target": []map[string]interface{}{
						{
							"mode": "PER_JOB",
							"name": "default-queue",
						},
					},
					"streaming_resource_setting": []map[string]interface{}{
						{
							"basic_resource_setting": []map[string]interface{}{
								{
									"taskmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "1Gi",
											"cpu":    "1",
										},
									},
									"parallelism": "1",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "1Gi",
											"cpu":    "1",
										},
									},
								},
							},
							"resource_setting_mode": "BASIC",
						},
					},
					"namespace": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script": "create temporary table `datagen` ( id varchar, name varchar ) with ( 'connector' = 'datagen' );  create temporary table `blackhole` ( id varchar, name varchar ) with ( 'connector' = 'blackhole' );  insert into blackhole select * from datagen;",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar"},
								},
							},
						},
					},
					"resource_id": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"flink_conf": map[string]interface{}{
						"\"execution.checkpointing.interval\"":      "180s",
						"\"execution.checkpointing.min-pause\"":     "180s",
						"\"restart-strategy\"":                      "fixed-delay",
						"\"restart-strategy.fixed-delay.attempts\"": "2147483647",
						"\"restart-strategy.fixed-delay.delay\"":    "10 s",
						"\"table.exec.state.ttl\"":                  "36 h",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name":   name,
						"description":       "This is a test deployment.",
						"engine_version":    "vvr-8.0.10-flink-1.17",
						"local_variables.#": "1",
						"execution_mode":    "STREAMING",
						"namespace":         CHECKSET,
						"resource_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "StdOut",
									"logger_level": "DEBUG",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "false",
									"expiration_days": "5",
								},
							},
						},
					},
					"deployment_name": name + "_update",
					"description":     "This is a test deployment 2.",
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"local_variables": []map[string]interface{}{
						{
							"value": "value1",
							"name":  "name1",
						},
						{
							"value": "value2",
							"name":  "name2",
						},
					},
					"labels": map[string]interface{}{
						"\"vvp\"": "b",
					},
					"streaming_resource_setting": []map[string]interface{}{
						{
							"basic_resource_setting": []map[string]interface{}{
								{
									"taskmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
									"parallelism": "2",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
								},
							},
							"resource_setting_mode": "BASIC",
						},
					},
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script": "CREATE TABLE result_table (id BIGINT, name STRING) WITH ( 'connector' = 'jdbc', 'url' = 'jdbc:mysql://localhost:3306/test', 'table-name' = 'result_table' );",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar", "oss://bucket-name/b.jar", "oss://bucket-name/c.jar"},
								},
							},
						},
					},
					"flink_conf": map[string]interface{}{
						"\"execution.checkpointing.interval\"":      "180s",
						"\"restart-strategy\"":                      "fixed-delay",
						"\"restart-strategy.fixed-delay.attempts\"": "2147483647",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name":   name + "_update",
						"description":       "This is a test deployment 2.",
						"engine_version":    "vvr-8.0.10-flink-1.17",
						"local_variables.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"streaming_resource_setting": []map[string]interface{}{
						{
							"resource_setting_mode": "EXPERT",
							"expert_resource_setting": []map[string]interface{}{
								{
									"resource_plan": "{\\\\n  \\\\\\\"ssgProfiles\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"name\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n      \\\\\\\"cpu\\\\\\\": 0.26,\\\\n      \\\\\\\"heap\\\\\\\": \\\\\\\"1 gb\\\\\\\",\\\\n      \\\\\\\"offHeap\\\\\\\": \\\\\\\"32 mb\\\\\\\",\\\\n      \\\\\\\"managed\\\\\\\": {},\\\\n      \\\\\\\"extended\\\\\\\": {}\\\\n    }\\\\n  ],\\\\n  \\\\\\\"nodes\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 1,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecTableSourceScan\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Source: datagen_source[7]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 2,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecCalc\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Calc[8]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 3,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"ConstraintEnforcer[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 4,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Sink: vvptest[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    }\\\\n  ],\\\\n  \\\\\\\"edges\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 1,\\\\n      \\\\\\\"target\\\\\\\": 2,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 2,\\\\n      \\\\\\\"target\\\\\\\": 3,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 3,\\\\n      \\\\\\\"target\\\\\\\": 4,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    }\\\\n  ],\\\\n  \\\\\\\"vertices\\\\\\\": {\\\\n    \\\\\\\"717c7b8afebbfb7137f6f0f99beb2a94\\\\\\\": [\\\\n      1,\\\\n      2,\\\\n      3,\\\\n      4\\\\n    ]\\\\n  }\\\\n}",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
								},
							},
						},
					},
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script":              "CREATE TABLE result_table (id BIGINT, name STRING) WITH ( 'connector' = 'jdbc', 'url' = 'jdbc:mysql://localhost:3306/test', 'table-name' = 'result_table' );",
									"additional_dependencies": []string{},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudRealtimeComputeDeploymentMap11895 = map[string]string{
	"deployment_id": CHECKSET,
}

func AlicloudRealtimeComputeDeploymentBasicDependence11895(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc.id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch.zone_id
}


`, name)
}

// Case test_4 11872
func TestAccAliCloudRealtimeComputeDeployment_basic11872(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_deployment.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeDeploymentMap11872)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeDeploymentBasicDependence11872)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"deployment_name": name,
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"resource_id":     "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"execution_mode":  "STREAMING",
					"deployment_target": []map[string]interface{}{
						{
							"mode": "PER_JOB",
							"name": "default-queue",
						},
					},
					"namespace": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"artifact": []map[string]interface{}{
						{
							"kind": "PYTHON",
							"python_artifact": []map[string]interface{}{
								{
									"entry_module": "test.py",
									"main_args":    "start from main",
									"additional_python_archives": []string{
										"oss://bucket-name/c.jar"},
									"additional_python_libraries": []string{
										"oss://bucket-name/b.jar"},
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar"},
									"python_artifact_uri": "oss://bucket-name/main.py",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name": name,
						"engine_version":  "vvr-8.0.10-flink-1.17",
						"execution_mode":  "STREAMING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"artifact": []map[string]interface{}{
						{
							"kind": "PYTHON",
							"python_artifact": []map[string]interface{}{
								{
									"entry_module": "test1.py",
									"main_args":    "start from main1",
									"additional_python_archives": []string{
										"oss://bucket-name/a2.jar", "oss://bucket-name/b2.jar", "oss://bucket-name/c2.jar"},
									"additional_python_libraries": []string{
										"oss://bucket-name/a1.jar", "oss://bucket-name/b1.jar", "oss://bucket-name/c1.jar"},
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar", "oss://bucket-name/b.jar", "oss://bucket-name/c.jar"},
									"python_artifact_uri": "oss://bucket-name/main1.py",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"artifact": []map[string]interface{}{
						{
							"kind": "PYTHON",
							"python_artifact": []map[string]interface{}{
								{
									"entry_module":                "test1.py",
									"main_args":                   "start from main1",
									"additional_python_archives":  []string{},
									"additional_python_libraries": []string{},
									"additional_dependencies":     []string{},
									"python_artifact_uri":         "oss://bucket-name/main.py",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudRealtimeComputeDeploymentMap11872 = map[string]string{
	"deployment_id": CHECKSET,
}

func AlicloudRealtimeComputeDeploymentBasicDependence11872(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc4" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch4" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc4.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket4" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket4.id
    }
  }
  vpc_id      = alicloud_vswitch.create_Vswitch4.vpc_id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch4.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch4.zone_id
}


`, name)
}

// Case test3 11873
func TestAccAliCloudRealtimeComputeDeployment_basic11873(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_deployment.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeDeploymentMap11873)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeDeploymentBasicDependence11873)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"deployment_name": name,
					"description":     "This is a test deployment.",
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"resource_id":     "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"execution_mode":  "STREAMING",
					"deployment_target": []map[string]interface{}{
						{
							"mode": "PER_JOB",
							"name": "default-queue",
						},
					},
					"streaming_resource_setting": []map[string]interface{}{
						{
							"resource_setting_mode": "EXPERT",
							"expert_resource_setting": []map[string]interface{}{
								{
									"resource_plan": "{\\\\n  \\\\\\\"ssgProfiles\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"name\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n      \\\\\\\"cpu\\\\\\\": 0.26,\\\\n      \\\\\\\"heap\\\\\\\": \\\\\\\"1 gb\\\\\\\",\\\\n      \\\\\\\"offHeap\\\\\\\": \\\\\\\"32 mb\\\\\\\",\\\\n      \\\\\\\"managed\\\\\\\": {},\\\\n      \\\\\\\"extended\\\\\\\": {}\\\\n    }\\\\n  ],\\\\n  \\\\\\\"nodes\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 1,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecTableSourceScan\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Source: datagen_source[7]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 2,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecCalc\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Calc[8]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 3,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"ConstraintEnforcer[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 4,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Sink: vvptest[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    }\\\\n  ],\\\\n  \\\\\\\"edges\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 1,\\\\n      \\\\\\\"target\\\\\\\": 2,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 2,\\\\n      \\\\\\\"target\\\\\\\": 3,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 3,\\\\n      \\\\\\\"target\\\\\\\": 4,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    }\\\\n  ],\\\\n  \\\\\\\"vertices\\\\\\\": {\\\\n    \\\\\\\"717c7b8afebbfb7137f6f0f99beb2a94\\\\\\\": [\\\\n      1,\\\\n      2,\\\\n      3,\\\\n      4\\\\n    ]\\\\n  }\\\\n}",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
								},
							},
						},
					},
					"namespace": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script": "create temporary table `datagen` ( id varchar, name varchar ) with ( 'connector' = 'datagen' );  create temporary table `blackhole` ( id varchar, name varchar ) with ( 'connector' = 'blackhole' );  insert into blackhole select * from datagen;",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name": name,
						"description":     "This is a test deployment.",
						"engine_version":  "vvr-8.0.10-flink-1.17",
						"execution_mode":  "STREAMING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deployment_target": []map[string]interface{}{
						{
							"mode": "SESSION",
							"name": "tf",
						},
					},
					"streaming_resource_setting": []map[string]interface{}{
						{
							"resource_setting_mode": "EXPERT",
							"expert_resource_setting": []map[string]interface{}{
								{
									"resource_plan": "{\\\\n  \\\\\\\"ssgProfiles\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"name\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n      \\\\\\\"cpu\\\\\\\": 0.26,\\\\n      \\\\\\\"heap\\\\\\\": \\\\\\\"1 gb\\\\\\\",\\\\n      \\\\\\\"offHeap\\\\\\\": \\\\\\\"32 mb\\\\\\\",\\\\n      \\\\\\\"managed\\\\\\\": {},\\\\n      \\\\\\\"extended\\\\\\\": {}\\\\n    }\\\\n  ],\\\\n  \\\\\\\"nodes\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 1,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecTableSourceScan\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Source: datagen_source[7]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 2,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecCalc\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Calc[8]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 3,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"ConstraintEnforcer[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 1,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    },\\\\n    {\\\\n      \\\\\\\"id\\\\\\\": 4,\\\\n      \\\\\\\"type\\\\\\\": \\\\\\\"StreamExecSink\\\\\\\",\\\\n      \\\\\\\"desc\\\\\\\": \\\\\\\"Sink: vvptest[9]\\\\\\\",\\\\n      \\\\\\\"profile\\\\\\\": {\\\\n        \\\\\\\"group\\\\\\\": \\\\\\\"default\\\\\\\",\\\\n        \\\\\\\"parallelism\\\\\\\": 2,\\\\n        \\\\\\\"maxParallelism\\\\\\\": 32768,\\\\n        \\\\\\\"minParallelism\\\\\\\": 1\\\\n      }\\\\n    }\\\\n  ],\\\\n  \\\\\\\"edges\\\\\\\": [\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 1,\\\\n      \\\\\\\"target\\\\\\\": 2,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 2,\\\\n      \\\\\\\"target\\\\\\\": 3,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    },\\\\n    {\\\\n      \\\\\\\"source\\\\\\\": 3,\\\\n      \\\\\\\"target\\\\\\\": 4,\\\\n      \\\\\\\"mode\\\\\\\": \\\\\\\"PIPELINED\\\\\\\",\\\\n      \\\\\\\"strategy\\\\\\\": \\\\\\\"FORWARD\\\\\\\"\\\\n    }\\\\n  ],\\\\n  \\\\\\\"vertices\\\\\\\": {\\\\n    \\\\\\\"717c7b8afebbfb7137f6f0f99beb2a94\\\\\\\": [\\\\n      1,\\\\n      2,\\\\n      3,\\\\n      4\\\\n    ]\\\\n  }\\\\n}",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudRealtimeComputeDeploymentMap11873 = map[string]string{
	"deployment_id": CHECKSET,
}

func AlicloudRealtimeComputeDeploymentBasicDependence11873(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc2" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch2" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc2.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket2" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket2.id
    }
  }
  vpc_id      = alicloud_vswitch.create_Vswitch2.vpc_id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch2.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch2.zone_id
}


`, name)
}

// Case test_2 11870
func TestAccAliCloudRealtimeComputeDeployment_basic11870(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_deployment.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeDeploymentMap11870)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeDeploymentBasicDependence11870)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"deployment_name": name,
					"description":     "This is a test deployment.",
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"resource_id":     "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"batch_resource_setting": []map[string]interface{}{
						{
							"basic_resource_setting": []map[string]interface{}{
								{
									"taskmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "1Gi",
											"cpu":    "1",
										},
									},
									"parallelism": "1",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "1Gi",
											"cpu":    "1",
										},
									},
								},
							},
							"max_slot": "1",
						},
					},
					"execution_mode": "BATCH",
					"deployment_target": []map[string]interface{}{
						{
							"mode": "PER_JOB",
							"name": "default-queue",
						},
					},
					"namespace": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"artifact": []map[string]interface{}{
						{
							"kind": "SQLSCRIPT",
							"sql_artifact": []map[string]interface{}{
								{
									"sql_script": "create temporary table `datagen` ( id varchar, name varchar ) with ( 'connector' = 'datagen' );  create temporary table `blackhole` ( id varchar, name varchar ) with ( 'connector' = 'blackhole' );  insert into blackhole select * from datagen;",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name": name,
						"description":     "This is a test deployment.",
						"engine_version":  "vvr-8.0.10-flink-1.17",
						"execution_mode":  "BATCH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"batch_resource_setting": []map[string]interface{}{
						{
							"basic_resource_setting": []map[string]interface{}{
								{
									"taskmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
									"parallelism": "2",
									"jobmanager_resource_setting_spec": []map[string]interface{}{
										{
											"memory": "2Gi",
											"cpu":    "2",
										},
									},
								},
							},
							"max_slot": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudRealtimeComputeDeploymentMap11870 = map[string]string{
	"deployment_id": CHECKSET,
}

func AlicloudRealtimeComputeDeploymentBasicDependence11870(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc1" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch1" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc1.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket1" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket1.id
    }
  }
  vpc_id      = alicloud_vswitch.create_Vswitch1.vpc_id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch1.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch1.zone_id
}


`, name)
}

// Case test_1_副本1763709151937_副本1763709161121 11871
func TestAccAliCloudRealtimeComputeDeployment_basic11871(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_deployment.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeDeploymentMap11871)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeDeploymentBasicDependence11871)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"log4j2_configuration_template": "<?xml version=\\\\\\\"1.0\\\\\\\" encoding=\\\\\\\"UTF-8\\\\\\\" standalone=\\\\\\\"no\\\\\\\"?>\\\\n<Configuration xmlns=\\\\\\\"http://logging.apache.org/log4j/2.0/config\\\\\\\" strict=\\\\\\\"true\\\\\\\" monitorInterval=\\\\\\\"30\\\\\\\">\\\\n    <Appenders>\\\\n        <Appender name=\\\\\\\"StdOut\\\\\\\" type=\\\\\\\"Console\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%d{yyyy-MM-dd HH:mm:ss,SSS} [%-tn] %-5p %-60c %x - %m%n\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"RollingFile\\\\\\\" type=\\\\\\\"RollingFile\\\\\\\" fileName=\\\\\\\"$${sys:log.file}\\\\\\\" filePattern=\\\\\\\"$${sys:log.file}.%i\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%d{yyyy-MM-dd HH:mm:ss,SSS} [%-tn] %-5p %-60c %x - %m%n\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n            <Policies>\\\\n                <SizeBasedTriggeringPolicy size=\\\\\\\"5 MB\\\\\\\"/>\\\\n            </Policies>\\\\n            <DefaultRolloverStrategy max=\\\\\\\"1\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"StdOutErrConsoleAppender\\\\\\\" type=\\\\\\\"Console\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%m\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"StdOutFileAppender\\\\\\\" type=\\\\\\\"RollingFile\\\\\\\" fileName=\\\\\\\"$${sys:stdout.file}\\\\\\\" filePattern=\\\\\\\"$${sys:stdout.file}.%i\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%m\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n            <Policies>\\\\n                <SizeBasedTriggeringPolicy size=\\\\\\\"5 MB\\\\\\\"/>\\\\n            </Policies>\\\\n            <DefaultRolloverStrategy max=\\\\\\\"2\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"StdErrFileAppender\\\\\\\" type=\\\\\\\"RollingFile\\\\\\\" fileName=\\\\\\\"$${sys:stderr.file}\\\\\\\" filePattern=\\\\\\\"$${sys:stderr.file}.%i\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%m\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n            <Policies>\\\\n                <SizeBasedTriggeringPolicy size=\\\\\\\"5 MB\\\\\\\"/>\\\\n            </Policies>\\\\n            <DefaultRolloverStrategy max=\\\\\\\"2\\\\\\\"/>\\\\n        </Appender>\\\\n    </Appenders>\\\\n    <Loggers>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"org.apache.hadoop\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"org.apache.kafka\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"org.apache.zookeeper\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"WARN\\\\\\\" name=\\\\\\\"com.aliyun.jindodata\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"ERROR\\\\\\\" name=\\\\\\\"org.jboss.netty.channel.DefaultChannelPipeline\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"ERROR\\\\\\\" name=\\\\\\\"com.ververica.platform.logging.appender.oss\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"ERROR\\\\\\\" name=\\\\\\\"org.apache.flink.fs.osshadoop.shaded.com.aliyun.oss\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"OFF\\\\\\\" name=\\\\\\\"org.apache.flink.runtime.rest.handler.job.JobDetailsHandler\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"StdOutErrRedirector.StdOut\\\\\\\" additivity=\\\\\\\"false\\\\\\\">\\\\n            <AppenderRef ref=\\\\\\\"StdOutFileAppender\\\\\\\"/>\\\\n            <AppenderRef ref=\\\\\\\"StdOutErrConsoleAppender\\\\\\\"/>\\\\n        </Logger>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"StdOutErrRedirector.StdErr\\\\\\\" additivity=\\\\\\\"false\\\\\\\">\\\\n            <AppenderRef ref=\\\\\\\"StdErrFileAppender\\\\\\\"/>\\\\n            <AppenderRef ref=\\\\\\\"StdOutErrConsoleAppender\\\\\\\"/>\\\\n        </Logger>\\\\n        <!-- User-configured loggers placeholder (Do not modify this line) -->\\\\n        <Root level=\\\\\\\"{{ rootLoggerLogLevel }}\\\\\\\">\\\\n            <AppenderRef ref=\\\\\\\"StdOut\\\\\\\"/>\\\\n            <AppenderRef ref=\\\\\\\"RollingFile\\\\\\\"/>\\\\n        </Root>\\\\n        \\\\n    {%- for name, level in userConfiguredLoggers -%} \\\\n       <Logger level=\\\\\\\"{{ level }}\\\\\\\" name=\\\\\\\"{{ name }}\\\\\\\"/> \\\\n    {%- endfor -%}\\\\n    \\\\n\\\\n    </Loggers>\\\\n</Configuration>",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "StdOut",
									"logger_level": "DEBUG",
								},
								{
									"logger_name":  "",
									"logger_level": "INFO",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "false",
									"expiration_days": "5",
								},
							},
						},
					},
					"deployment_name": name,
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"resource_id":     "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"execution_mode":  "STREAMING",
					"deployment_target": []map[string]interface{}{
						{
							"mode": "PER_JOB",
							"name": "default-queue",
						},
					},
					"namespace": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"artifact": []map[string]interface{}{
						{
							"kind": "JAR",
							"jar_artifact": []map[string]interface{}{
								{
									"jar_uri":     "oss://bucket-name/main.jar",
									"main_args":   "start from main",
									"entry_class": "org.apache.flink.test",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name": name,
						"engine_version":  "vvr-8.0.10-flink-1.17",
						"execution_mode":  "STREAMING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"artifact": []map[string]interface{}{
						{
							"kind": "JAR",
							"jar_artifact": []map[string]interface{}{
								{
									"jar_uri":                 "oss://bucket-name/main.jar",
									"main_args":               "start from main1",
									"entry_class":             "org.apache.flink.test1",
									"additional_dependencies": []string{},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"logging_profile": "",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_level": "INFO",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "true",
									"expiration_days": "7",
								},
							},
						},
					},
					"artifact": []map[string]interface{}{
						{
							"kind": "JAR",
							"jar_artifact": []map[string]interface{}{
								{
									"jar_uri":     "oss://bucket-name/main1.jar",
									"main_args":   "start from main1",
									"entry_class": "org.apache.flink.test1",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar", "oss://bucket-name/b.jar", "oss://bucket-name/c.jar"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudRealtimeComputeDeploymentMap11871 = map[string]string{
	"deployment_id": CHECKSET,
}

func AlicloudRealtimeComputeDeploymentBasicDependence11871(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc3" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch3" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc3.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket3" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket3.id
    }
  }
  vpc_id      = alicloud_vpc.create_Vpc3.id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch3.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch3.zone_id
}


`, name)
}

// Case test_1 11869
func TestAccAliCloudRealtimeComputeDeployment_basic11869(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_deployment.default"
	ra := resourceAttrInit(resourceId, AlicloudRealtimeComputeDeploymentMap11869)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RealtimeComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRealtimeComputeDeployment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRealtimeComputeDeploymentBasicDependence11869)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"logging_profile":               "default",
							"log4j2_configuration_template": "<?xml version=\\\\\\\"1.0\\\\\\\" encoding=\\\\\\\"UTF-8\\\\\\\" standalone=\\\\\\\"no\\\\\\\"?>\\\\n<Configuration xmlns=\\\\\\\"http://logging.apache.org/log4j/2.0/config\\\\\\\" strict=\\\\\\\"true\\\\\\\" monitorInterval=\\\\\\\"30\\\\\\\">\\\\n    <Appenders>\\\\n        <Appender name=\\\\\\\"StdOut\\\\\\\" type=\\\\\\\"Console\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%d{yyyy-MM-dd HH:mm:ss,SSS} [%-tn] %-5p %-60c %x - %m%n\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"RollingFile\\\\\\\" type=\\\\\\\"RollingFile\\\\\\\" fileName=\\\\\\\"$${sys:log.file}\\\\\\\" filePattern=\\\\\\\"$${sys:log.file}.%i\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%d{yyyy-MM-dd HH:mm:ss,SSS} [%-tn] %-5p %-60c %x - %m%n\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n            <Policies>\\\\n                <SizeBasedTriggeringPolicy size=\\\\\\\"5 MB\\\\\\\"/>\\\\n            </Policies>\\\\n            <DefaultRolloverStrategy max=\\\\\\\"1\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"StdOutErrConsoleAppender\\\\\\\" type=\\\\\\\"Console\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%m\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"StdOutFileAppender\\\\\\\" type=\\\\\\\"RollingFile\\\\\\\" fileName=\\\\\\\"$${sys:stdout.file}\\\\\\\" filePattern=\\\\\\\"$${sys:stdout.file}.%i\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%m\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n            <Policies>\\\\n                <SizeBasedTriggeringPolicy size=\\\\\\\"5 MB\\\\\\\"/>\\\\n            </Policies>\\\\n            <DefaultRolloverStrategy max=\\\\\\\"2\\\\\\\"/>\\\\n        </Appender>\\\\n        <Appender name=\\\\\\\"StdErrFileAppender\\\\\\\" type=\\\\\\\"RollingFile\\\\\\\" fileName=\\\\\\\"$${sys:stderr.file}\\\\\\\" filePattern=\\\\\\\"$${sys:stderr.file}.%i\\\\\\\">\\\\n            <Layout pattern=\\\\\\\"%m\\\\\\\" type=\\\\\\\"PatternLayout\\\\\\\" charset=\\\\\\\"UTF-8\\\\\\\"/>\\\\n            <Policies>\\\\n                <SizeBasedTriggeringPolicy size=\\\\\\\"5 MB\\\\\\\"/>\\\\n            </Policies>\\\\n            <DefaultRolloverStrategy max=\\\\\\\"2\\\\\\\"/>\\\\n        </Appender>\\\\n    </Appenders>\\\\n    <Loggers>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"org.apache.hadoop\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"org.apache.kafka\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"org.apache.zookeeper\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"WARN\\\\\\\" name=\\\\\\\"com.aliyun.jindodata\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"ERROR\\\\\\\" name=\\\\\\\"org.jboss.netty.channel.DefaultChannelPipeline\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"ERROR\\\\\\\" name=\\\\\\\"com.ververica.platform.logging.appender.oss\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"ERROR\\\\\\\" name=\\\\\\\"org.apache.flink.fs.osshadoop.shaded.com.aliyun.oss\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"OFF\\\\\\\" name=\\\\\\\"org.apache.flink.runtime.rest.handler.job.JobDetailsHandler\\\\\\\"/>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"StdOutErrRedirector.StdOut\\\\\\\" additivity=\\\\\\\"false\\\\\\\">\\\\n            <AppenderRef ref=\\\\\\\"StdOutFileAppender\\\\\\\"/>\\\\n            <AppenderRef ref=\\\\\\\"StdOutErrConsoleAppender\\\\\\\"/>\\\\n        </Logger>\\\\n        <Logger level=\\\\\\\"INFO\\\\\\\" name=\\\\\\\"StdOutErrRedirector.StdErr\\\\\\\" additivity=\\\\\\\"false\\\\\\\">\\\\n            <AppenderRef ref=\\\\\\\"StdErrFileAppender\\\\\\\"/>\\\\n            <AppenderRef ref=\\\\\\\"StdOutErrConsoleAppender\\\\\\\"/>\\\\n        </Logger>\\\\n        <!-- User-configured loggers placeholder (Do not modify this line) -->\\\\n        <Root level=\\\\\\\"{{ rootLoggerLogLevel }}\\\\\\\">\\\\n            <AppenderRef ref=\\\\\\\"StdOut\\\\\\\"/>\\\\n            <AppenderRef ref=\\\\\\\"RollingFile\\\\\\\"/>\\\\n        </Root>\\\\n        \\\\n    {%- for name, level in userConfiguredLoggers -%} \\\\n       <Logger level=\\\\\\\"{{ level }}\\\\\\\" name=\\\\\\\"{{ name }}\\\\\\\"/> \\\\n    {%- endfor -%}\\\\n    \\\\n\\\\n    </Loggers>\\\\n</Configuration>",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "StdOut",
									"logger_level": "DEBUG",
								},
								{
									"logger_name":  "",
									"logger_level": "INFO",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "false",
									"expiration_days": "5",
								},
							},
						},
					},
					"deployment_name": name,
					"engine_version":  "vvr-8.0.10-flink-1.17",
					"resource_id":     "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.resource_id}",
					"execution_mode":  "STREAMING",
					"deployment_target": []map[string]interface{}{
						{
							"mode": "PER_JOB",
							"name": "default-queue",
						},
					},
					"namespace": "${alicloud_realtime_compute_vvp_instance.create_VvpInstance.vvp_instance_name}-default",
					"artifact": []map[string]interface{}{
						{
							"kind": "JAR",
							"jar_artifact": []map[string]interface{}{
								{
									"jar_uri":     "oss://bucket-name/main.jar",
									"main_args":   "start from main",
									"entry_class": "org.apache.flink.test",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deployment_name": name,
						"engine_version":  "vvr-8.0.10-flink-1.17",
						"execution_mode":  "STREAMING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"logging_profile": "",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "",
									"logger_level": "INFO",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "true",
									"expiration_days": "7",
								},
							},
						},
					},
					"artifact": []map[string]interface{}{
						{
							"kind": "JAR",
							"jar_artifact": []map[string]interface{}{
								{
									"jar_uri":                 "oss://bucket-name/main.jar",
									"main_args":               "start from main1",
									"entry_class":             "org.apache.flink.test1",
									"additional_dependencies": []string{},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging": []map[string]interface{}{
						{
							"logging_profile": "",
							"log4j_loggers": []map[string]interface{}{
								{
									"logger_name":  "",
									"logger_level": "INFO",
								},
							},
							"log_reserve_policy": []map[string]interface{}{
								{
									"open_history":    "true",
									"expiration_days": "5",
								},
							},
						},
					},
					"artifact": []map[string]interface{}{
						{
							"kind": "JAR",
							"jar_artifact": []map[string]interface{}{
								{
									"jar_uri":     "oss://bucket-name/main1.jar",
									"main_args":   "start from main1",
									"entry_class": "org.apache.flink.test1",
									"additional_dependencies": []string{
										"oss://bucket-name/a.jar", "oss://bucket-name/b.jar", "oss://bucket-name/c.jar"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudRealtimeComputeDeploymentMap11869 = map[string]string{
	"deployment_id": CHECKSET,
}

func AlicloudRealtimeComputeDeploymentBasicDependence11869(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "create_Vpc3" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-deployment"
}

resource "alicloud_vswitch" "create_Vswitch3" {
  is_default   = false
  vpc_id       = alicloud_vpc.create_Vpc3.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vSwitch-deployment"
}

resource "alicloud_oss_bucket" "create_bucket3" {
}

resource "alicloud_realtime_compute_vvp_instance" "create_VvpInstance" {
  vvp_instance_name = "code-test-tf-deployment"
  storage {
    oss {
      bucket = alicloud_oss_bucket.create_bucket3.id
    }
  }
  vpc_id      = alicloud_vswitch.create_Vswitch3.vpc_id
  vswitch_ids = ["${alicloud_vswitch.create_Vswitch3.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.create_Vswitch3.zone_id
}


`, name)
}

// Test RealtimeCompute Deployment. <<< Resource test cases, automatically generated.
