package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_logtail_config", &resource.Sweeper{
		Name: "alicloud_logtail_config",
		F:    testSweepLogConfigs,
	})
}

func testSweepLogConfigs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.ListProject()
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Log Projects: %s", WrapError(err))
	}
	names, _ := raw.([]string)

	for _, v := range names {
		name := v
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				cfNameList, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					cfNames, _, cfErr := slsClient.ListConfig(name, 0, 100)
					return cfNames, cfErr
				})
				if err != nil {
					log.Printf("[ERROR] Error retrieving Log config: %s", WrapError(err))
				}
				for _, cfName := range cfNameList.([]string) {
					log.Printf("[INFO] Deleting Log config: %s", cfName)
					_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
						return nil, slsClient.DeleteConfig(name, cfName)
					})
					if err != nil {
						log.Printf("[ERROR] Failed to delete Log Config (%s): %s", cfName, err)
					}
				}
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Log Project: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting Log Project: %s", name)
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Log Project (%s): %s", name, err)
		}
	}
	return nil
}

func TestAccAlicloudLogTail_basic(t *testing.T) {
	var v *sls.LogConfig
	resourceId := "alicloud_logtail_config.default"
	ra := resourceAttrInit(resourceId, logTailMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogtailconfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogTailDependence)

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
					"project":      "${alicloud_log_project.default.name}",
					"logstore":     "${alicloud_log_store.default.name}",
					"input_type":   "file",
					"name":         name,
					"output_type":  "LogService",
					"input_detail": `{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"input_detail": "{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_sample": "test-logtail-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_sample": "test-logtail-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_type":   "file",
					"input_detail": `{\"autoExtend\":true,\"discardUnmatch\":true,\"enableRawLog\":true,\"fileEncoding\":\"utf8\",\"filePattern\":\"*\",\"key\":[\"test\",\"test2\"],\"logPath\":\"/logPath\",\"logType\":\"delimiter_log\",\"maxDepth\":999,\"quote\":\"\\\"\",\"separator\":\",\",\"timeKey\":\"test\",\"topicFormat\":\"default\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_type":   "file",
						"input_detail": "{\"autoExtend\":true,\"discardUnmatch\":true,\"enableRawLog\":true,\"fileEncoding\":\"utf8\",\"filePattern\":\"*\",\"key\":[\"test\",\"test2\"],\"logPath\":\"/logPath\",\"logType\":\"delimiter_log\",\"maxDepth\":999,\"quote\":\"\\\"\",\"separator\":\",\",\"timeKey\":\"test\",\"topicFormat\":\"default\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project":      "${alicloud_log_project.default.name}",
					"logstore":     "${alicloud_log_store.default.name}",
					"input_type":   "file",
					"name":         name,
					"output_type":  "LogService",
					"input_detail": `{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}`,
					"log_sample":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_detail": "{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}",
						"log_sample":   REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"input_detail"},
			},
		},
	})
}

func TestAccAlicloudLogTail_plugin(t *testing.T) {
	var v *sls.LogConfig
	resourceId := "alicloud_logtail_config.default"
	ra := resourceAttrInit(resourceId, logTailMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogtailconfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogTailDependence)

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
					"project":      "${alicloud_log_project.default.name}",
					"logstore":     "${alicloud_log_store.default.name}",
					"input_type":   "plugin",
					"name":         name,
					"output_type":  "LogService",
					"input_detail": `{\"plugin\":{\"inputs\":[{\"detail\":{\"ExcludeEnv\":null,\"ExcludeLabel\":null,\"IncludeEnv\":null,\"IncludeLabel\":null,\"Stderr\":true,\"Stdout\":true},\"type\":\"service_docker_stdout\"}]}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"input_type":   "plugin",
						"input_detail": "{\"plugin\":{\"inputs\":[{\"detail\":{\"ExcludeEnv\":null,\"ExcludeLabel\":null,\"IncludeEnv\":null,\"IncludeLabel\":null,\"Stderr\":true,\"Stdout\":true},\"type\":\"service_docker_stdout\"}]}}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_sample": "test-logtail-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_sample": "test-logtail-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project":      "${alicloud_log_project.default.name}",
					"logstore":     "${alicloud_log_store.default.name}",
					"input_type":   "plugin",
					"name":         name,
					"output_type":  "LogService",
					"input_detail": `{\"plugin\":{\"inputs\":[{\"detail\":{\"ExcludeEnv\":null,\"ExcludeLabel\":null,\"IncludeEnv\":null,\"IncludeLabel\":null,\"Stderr\":true,\"Stdout\":true},\"type\":\"service_docker_stdout\"}]}}`,
					"log_sample":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_detail": "{\"plugin\":{\"inputs\":[{\"detail\":{\"ExcludeEnv\":null,\"ExcludeLabel\":null,\"IncludeEnv\":null,\"IncludeLabel\":null,\"Stderr\":true,\"Stdout\":true},\"type\":\"service_docker_stdout\"}]}}",
						"log_sample":   REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"input_detail"},
			},
		},
	})
}

func TestAccAlicloudLogTail_multi(t *testing.T) {
	var v *sls.LogConfig
	resourceId := "alicloud_logtail_config.default.4"
	ra := resourceAttrInit(resourceId, logTailMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogtailconfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogTailDependence)

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
					"project":      "${alicloud_log_project.default.name}",
					"logstore":     "${alicloud_log_store.default.name}",
					"input_type":   "file",
					"name":         name + "${count.index}",
					"output_type":  "LogService",
					"input_detail": `{\"discardUnmatch\":false,\"enableRawLog\":true,\"fileEncoding\":\"gbk\",\"filePattern\":\"access.log\",\"logPath\":\"/logPath\",\"logType\":\"json_log\",\"maxDepth\":10,\"topicFormat\":\"default\"}`,
					"count":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceLogTailDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_log_project" "default"{
	name = "${var.name}"
	description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  	project = "${alicloud_log_project.default.name}"
  	name = "${var.name}"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
`, name)
}

var logTailMap = map[string]string{
	"name":         CHECKSET,
	"project":      CHECKSET,
	"logstore":     CHECKSET,
	"input_type":   "file",
	"output_type":  "LogService",
	"input_detail": CHECKSET,
}
