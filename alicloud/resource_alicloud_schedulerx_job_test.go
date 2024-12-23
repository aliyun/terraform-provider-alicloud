package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Schedulerx Job. >>> Resource test cases, automatically generated.
// Case 预发环境_20241220_杭州region_组合用例B 9597
func TestAccAliCloudSchedulerxJob_basic9597(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_schedulerx_job.default"
	ra := resourceAttrInit(resourceId, AlicloudSchedulerxJobMap9597)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SchedulerxServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSchedulerxJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sschedulerxjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSchedulerxJobBasicDependence9597)
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
					"status":                "Enable",
					"max_attempt":           "0",
					"description":           "Job资源用例自动生成的的任务",
					"success_notice_enable": "false",
					"job_name":              name,
					"max_concurrency":       "1",
					"time_config": []map[string]interface{}{
						{
							"time_type": "-1",
							"calendar":  "workday",
						},
					},
					"namespace": "${alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid}",
					"group_id":  "${alicloud_schedulerx_app_group.CreateAppGroup.group_id}",
					"job_type":  "shell",
					"job_monitor_info": []map[string]interface{}{
						{
							"monitor_config": []map[string]interface{}{
								{
									"timeout":      "7200",
									"send_channel": "webhook",
								},
							},
							"contact_info": []map[string]interface{}{
								{
									"user_name": "tangtao-1",
								},
								{
									"user_name": "tangtao-2",
								},
								{
									"user_name": "tangtao-3",
								},
							},
						},
					},
					"content":          "echo 'hello'",
					"namespace_source": "schedulerx",
					"attempt_interval": "30",
					"execute_mode":     "standalone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                "Enable",
						"max_attempt":           "0",
						"description":           "Job资源用例自动生成的的任务",
						"success_notice_enable": "false",
						"job_name":              name,
						"max_concurrency":       "1",
						"namespace":             CHECKSET,
						"group_id":              CHECKSET,
						"job_type":              "shell",
						"content":               "echo 'hello'",
						"namespace_source":      "schedulerx",
						"attempt_interval":      "30",
						"execute_mode":          "standalone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": "echo 'helllo schedulerx'",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content": "echo 'helllo schedulerx'",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fail_times", "namespace_source", "success_notice_enable", "task_dispatch_mode", "template", "timezone"},
			},
		},
	})
}

var AlicloudSchedulerxJobMap9597 = map[string]string{
	"job_id": CHECKSET,
}

func AlicloudSchedulerxJobBasicDependence9597(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_schedulerx_namespace" "CreateNameSpace" {
  namespace_name = var.name
  description    = "namespace 资源用例自动创建的命名空间"
}

resource "alicloud_schedulerx_app_group" "CreateAppGroup" {
  namespace             = alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid
  group_id              = "test-appgroup-pop-autotest"
  description           = "appgroup 资源用例生成"
  monitor_contacts_json = "[{\"userName\":\"张三\",\"userPhone\":\"89756******\"},{\"userName\":\"李四\",\"ding\":\"http://www.example.com\"}]"
  app_name              = "test-appgroup-pop-autotest"
  app_version           = "1"
  namespace_name        = alicloud_schedulerx_namespace.CreateNameSpace.namespace_name
  monitor_config_json   = "{\"sendChannel\":\"sms,ding\"}"
  app_type              = "2"
  max_jobs              = "100"
  namespace_source      = "schedulerx"
}


`, name)
}

// Case 预发环境_20241220_杭州region_组合用例A 9548
func TestAccAliCloudSchedulerxJob_basic9548(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_schedulerx_job.default"
	ra := resourceAttrInit(resourceId, AlicloudSchedulerxJobMap9548)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SchedulerxServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSchedulerxJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sschedulerxjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSchedulerxJobBasicDependence9548)
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
					"timezone":        "GTM+7",
					"status":          "Enable",
					"max_attempt":     "0",
					"description":     "job资源用例自动生成的任务",
					"parameters":      "hello word",
					"job_name":        name,
					"max_concurrency": "1",
					"time_config": []map[string]interface{}{
						{
							"data_offset":     "1",
							"time_expression": "100000",
							"time_type":       "3",
							"calendar":        "workday",
						},
					},
					"map_task_xattrs": []map[string]interface{}{
						{
							"task_max_attempt":      "1",
							"task_attempt_interval": "1",
							"consumer_size":         "5",
							"queue_size":            "10000",
							"dispatcher_size":       "5",
							"page_size":             "100",
						},
					},
					"namespace": "${alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid}",
					"group_id":  "${alicloud_schedulerx_app_group.CreateAppGroup.group_id}",
					"job_type":  "java",
					"job_monitor_info": []map[string]interface{}{
						{
							"contact_info": []map[string]interface{}{
								{
									"user_phone": "18838204961",
									"user_name":  "tangtao-1",
									"ding":       "https://alidocs.dingtalk.com",
									"user_mail":  "12345678@qq.com",
								},
								{
									"user_phone": "18838204961",
									"user_name":  "tangtao-2",
									"ding":       "https://alidocs.dingtalk.com1",
									"user_mail":  "123456789@qq.com",
								},
							},
							"monitor_config": []map[string]interface{}{
								{
									"timeout":             "7200",
									"send_channel":        "sms",
									"timeout_kill_enable": "true",
									"timeout_enable":      "true",
									"fail_enable":         "true",
									"miss_worker_enable":  "true",
								},
							},
						},
					},
					"class_name":            "com.aliyun.schedulerx.example.processor.SimpleJob",
					"namespace_source":      "schedulerx",
					"attempt_interval":      "30",
					"fail_times":            "1",
					"execute_mode":          "batch",
					"x_attrs":               "{\\\"consumerSize\\\":5,\\\"dispatcherSize\\\":5,\\\"taskMaxAttempt\\\":1,\\\"taskAttemptInterval\\\":1,\\\"taskDispatchMode\\\":\\\"push\\\",\\\"failover\\\":true,\\\"execOnMaster\\\":true,\\\"produceInterval\\\":3,\\\"pageSize\\\":100,\\\"queueSize\\\":10000,\\\"globalConsumerSize\\\":1000}",
					"success_notice_enable": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timezone":              "GTM+7",
						"status":                "Enable",
						"max_attempt":           "0",
						"description":           "job资源用例自动生成的任务",
						"parameters":            "hello word",
						"job_name":              name,
						"max_concurrency":       "1",
						"namespace":             CHECKSET,
						"group_id":              CHECKSET,
						"job_type":              "java",
						"class_name":            "com.aliyun.schedulerx.example.processor.SimpleJob",
						"namespace_source":      "schedulerx",
						"attempt_interval":      "30",
						"fail_times":            "1",
						"execute_mode":          "batch",
						"x_attrs":               "{\"consumerSize\":5,\"dispatcherSize\":5,\"taskMaxAttempt\":1,\"taskAttemptInterval\":1,\"taskDispatchMode\":\"push\",\"failover\":true,\"execOnMaster\":true,\"produceInterval\":3,\"pageSize\":100,\"queueSize\":10000,\"globalConsumerSize\":1000}",
						"success_notice_enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timezone":        "GTM+8",
					"max_attempt":     "2",
					"description":     "hahaha",
					"parameters":      "hello word ,schedlerx",
					"job_name":        name + "_update",
					"max_concurrency": "2",
					"time_config": []map[string]interface{}{
						{
							"time_type":       "1",
							"data_offset":     "2",
							"time_expression": "0 0 18 1 */1 ?",
							"calendar":        "finance",
						},
					},
					"map_task_xattrs": []map[string]interface{}{
						{
							"task_max_attempt":      "2",
							"task_attempt_interval": "2",
							"consumer_size":         "10",
							"queue_size":            "5000",
							"dispatcher_size":       "10",
							"page_size":             "200",
						},
					},
					"job_monitor_info": []map[string]interface{}{
						{
							"monitor_config": []map[string]interface{}{
								{
									"timeout_kill_enable": "false",
									"timeout_enable":      "false",
									"fail_enable":         "false",
									"miss_worker_enable":  "false",
									"timeout":             "1000",
									"send_channel":        "ding",
								},
							},
							"contact_info": []map[string]interface{}{
								{
									"user_name":  "tangtao-update",
									"user_mail":  "iamispangpang@163.com",
									"user_phone": "12588888888",
									"ding":       "https://alidocs.dingtalk.com/i/nodes/dQPGYqjpJYZnRbNYCLggRQjP8akx1Z5N?utm_scene=team_space",
								},
								{
									"user_phone": "18888888888",
									"user_name":  "tangtao-3",
									"user_mail":  "12345678@qq.com",
									"ding":       "https://alidocs.dingtalk.com",
								},
								{
									"user_phone": "18888888666",
									"user_name":  "tangtao-4",
									"ding":       "1233",
									"user_mail":  "22345678@qq.com",
								},
							},
						},
					},
					"class_name":            "hello word",
					"attempt_interval":      "100",
					"fail_times":            "3",
					"execute_mode":          "grid",
					"x_attrs":               "{\\\"consumerSize\\\":10,\\\"dispatcherSize\\\":10,\\\"taskMaxAttempt\\\":2,\\\"taskAttemptInterval\\\":2,\\\"taskDispatchMode\\\":\\\"push\\\",\\\"failover\\\":true,\\\"execOnMaster\\\":true,\\\"produceInterval\\\":3,\\\"pageSize\\\":200,\\\"queueSize\\\":5000,\\\"globalConsumerSize\\\":1000}",
					"success_notice_enable": "true",
					"task_dispatch_mode":    "push",
					"template":              "apiVersion: v1 kind: Pod metadata:   name: schedulerx-python-{JOB_ID}11   namespace: {NAMESPACE} spec:   containers:   - name: python-job     image: python     imagePullPolicy: IfNotPresent     volumeMounts:     - name: script-python       mountPath: script/python     command: [\\\"python\\\",\\\"script/python/python-{JOB_ID}.py\\\"]   volumes:   - name: script-python     configMap:       name: schedulerx-configmap       items:       - key: schedulerx-python-{JOB_ID}         path: python-{JOB_ID}.py   restartPolicy: Never",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timezone":              "GTM+8",
						"max_attempt":           "2",
						"description":           "hahaha",
						"parameters":            "hello word ,schedlerx",
						"job_name":              name + "_update",
						"max_concurrency":       "2",
						"class_name":            "hello word",
						"attempt_interval":      "100",
						"fail_times":            "3",
						"execute_mode":          "grid",
						"x_attrs":               "{\"consumerSize\":10,\"dispatcherSize\":10,\"taskMaxAttempt\":2,\"taskAttemptInterval\":2,\"taskDispatchMode\":\"push\",\"failover\":true,\"execOnMaster\":true,\"produceInterval\":3,\"pageSize\":200,\"queueSize\":5000,\"globalConsumerSize\":1000}",
						"success_notice_enable": "true",
						"task_dispatch_mode":    "push",
						"template":              "apiVersion: v1 kind: Pod metadata:   name: schedulerx-python-{JOB_ID}11   namespace: {NAMESPACE} spec:   containers:   - name: python-job     image: python     imagePullPolicy: IfNotPresent     volumeMounts:     - name: script-python       mountPath: script/python     command: [\"python\",\"script/python/python-{JOB_ID}.py\"]   volumes:   - name: script-python     configMap:       name: schedulerx-configmap       items:       - key: schedulerx-python-{JOB_ID}         path: python-{JOB_ID}.py   restartPolicy: Never",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disable",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"fail_times", "namespace_source", "success_notice_enable", "task_dispatch_mode", "template", "timezone"},
			},
		},
	})
}

var AlicloudSchedulerxJobMap9548 = map[string]string{
	"job_id": CHECKSET,
}

func AlicloudSchedulerxJobBasicDependence9548(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_schedulerx_namespace" "CreateNameSpace" {
  namespace_name = var.name
  description    = "pop资源用例自动生成"
}

resource "alicloud_schedulerx_app_group" "CreateAppGroup" {
  description    = "appgroup 资源用例生成"
  namespace      = alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid
  group_id       = "test-appgroup-pop-autotest"
  app_name       = "test-appgroup-pop-autotest"
  app_version    = "2"
  namespace_name = alicloud_schedulerx_namespace.CreateNameSpace.namespace_name
  app_type       = "2"
  max_jobs       = "1000"
}


`, name)
}

// Test Schedulerx Job. <<< Resource test cases, automatically generated.
