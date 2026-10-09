package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudEfloInvocation_basic10348(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_invocation.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloInvocationMap10348)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloInvocationBasicDependence10348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "test",
					"content_encoding": "Base64",
					"name":             name,
					"repeat_mode":      "Once",
					"parameters": map[string]interface{}{
						"\"name\"": "Jack",
					},
					"node_id_list": []string{
						"e01-cn-rno46i6rdfn"},
					"timeout":          "68",
					"command_content":  "ZWNobyAxMjM=",
					"working_dir":      "/home/",
					"username":         "root",
					"enable_parameter": "false",
					"termination_mode": "ProcessTree",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "test",
						"content_encoding": "Base64",
						"name":             name,
						"repeat_mode":      "Once",
						"timeout":          "68",
						"command_content":  "ZWNobyAxMjM=",
						"working_dir":      "/home/",
						"username":         "root",
						"enable_parameter": "false",
						"termination_mode": "ProcessTree",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"command_content", "command_id", "content_encoding", "description", "enable_parameter", "frequency", "launcher", "name", "parameters", "repeat_mode", "termination_mode", "timeout", "username", "working_dir", "node_id_list"},
			},
		},
	})
}

var AlicloudEfloInvocationMap10348 = map[string]string{}

func AlicloudEfloInvocationBasicDependence10348(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// lintignore: AT001
func TestAccAliCloudEfloInvocation_sendFile86848(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_invocation.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloInvocationMap86848)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloSendFileResults")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloInvocationBasicDependence86848)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"invocation_type": "SendFile",
					"name":            name,
					"content":         "tf-sendfile-content",
					"content_type":    "PlainText",
					"target_dir":      "/tmp/",
					"description":     "tf test sendfile",
					"node_id_list": []string{
						"e01-cn-rno46i6rdfn",
						"e01-cn-rno46i6rdf2"},
					"file_mode":  "0644",
					"file_owner": "root",
					"file_group": "root",
					"overwrite":  "true",
					"timeout":    "60",
					"command_id": "c-test",
					"frequency":  "rate(5m)",
					"launcher":   "python3 -u test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"invocation_type": "SendFile",
						"name":            name,
						"content_type":    "PlainText",
						"target_dir":      "/tmp/",
						"description":     "tf test sendfile",
						"file_mode":       "0644",
						"file_owner":      "root",
						"file_group":      "root",
						"overwrite":       "true",
						"timeout":         "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"invocation_type": "SendFile",
					"name":            name,
					"content":         "tf-sendfile-content",
					"content_type":    "PlainText",
					"target_dir":      "/tmp/",
					"description":     "tf test sendfile",
					"node_id_list": []string{
						"e01-cn-rno46i6rdf2",
						"e01-cn-rno46i6rdfn"},
					"file_mode":  "0644",
					"file_owner": "root",
					"file_group": "root",
					"overwrite":  "true",
					"timeout":    "60",
					"command_id": "c-test",
					"frequency":  "rate(5m)",
					"launcher":   "python3 -u test",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"invocation_type": "SendFile",
					"name":            name,
					"content":         "tf-sendfile-content",
					"content_type":    "PlainText",
					"target_dir":      "/tmp/",
					"description":     "tf test sendfile",
					"node_id_list": []string{
						"e01-cn-rno46i6rdf2",
						"e01-cn-rno46i6rdfn"},
					"file_mode":  "0644",
					"file_owner": "root",
					"file_group": "root",
					"overwrite":  "true",
					"timeout":    "60",
					"command_id": "c-test",
					"frequency":  "rate(5m)",
					"launcher":   "python3 -u test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"invocation_type": "SendFile",
						"name":            name,
						"content_type":    "PlainText",
						"target_dir":      "/tmp/",
						"description":     "tf test sendfile",
						"file_mode":       "0644",
						"file_owner":      "root",
						"file_group":      "root",
						"overwrite":       "true",
						"timeout":         "60",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content", "content_type", "description", "file_mode", "file_group", "file_owner", "name", "overwrite", "target_dir", "timeout", "node_id_list", "command_id", "frequency", "launcher"},
			},
		},
	})
}

var AlicloudEfloInvocationMap86848 = map[string]string{}

func AlicloudEfloInvocationBasicDependence86848(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
