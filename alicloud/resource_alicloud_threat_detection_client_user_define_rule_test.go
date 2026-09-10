package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection ClientUserDefineRule. >>> Resource test cases, automatically generated.
// Case 4405
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4405(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4405)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4405)
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
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"type":                         "5",
					"registry_content":             "123",
					"registry_key":                 "123",
					"proc_path":                    "/root/bash",
					"cmdline":                      "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"type":                         "5",
						"registry_content":             "123",
						"registry_key":                 "123",
						"proc_path":                    "/root/bash",
						"cmdline":                      "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"registry_content": "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"registry_content": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"registry_key": "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"registry_key": "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"platform": "windows",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform": "windows",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_user_define_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_user_define_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"registry_content": "123sada",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"registry_content": "123sada",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_proc_path": "/root/bash/dsa",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_proc_path": "/root/bash/dsa",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cmdline": "bash b",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cmdline": "bash b",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/root/bash/d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/root/bash/d",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "bash b",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "bash b",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"registry_key": "123ads",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"registry_key": "123ads",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type":                  "0",
					"platform":                     "windows",
					"registry_content":             "123",
					"client_user_define_rule_name": name + "_update",
					"parent_proc_path":             "/root/bash",
					"type":                         "5",
					"cmdline":                      "bash",
					"proc_path":                    "/root/bash",
					"parent_cmdline":               "bash",
					"registry_key":                 "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "windows",
						"registry_content":             "123",
						"client_user_define_rule_name": name + "_update",
						"parent_proc_path":             "/root/bash",
						"type":                         "5",
						"cmdline":                      "bash",
						"proc_path":                    "/root/bash",
						"parent_cmdline":               "bash",
						"registry_key":                 "123",
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

var AlicloudThreatDetectionClientUserDefineRuleMap4405 = map[string]string{
	"create_time": CHECKSET,
	"port_str":    CHECKSET,
}

func AlicloudThreatDetectionClientUserDefineRuleBasicDependence4405(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4350
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4350(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4350)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4350)
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
					"platform":                     "all",
					"action_type":                  "0",
					"client_user_define_rule_name": name,
					"type":                         "1",
					"hash":                         "dacsda",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform":                     "all",
						"action_type":                  "0",
						"client_user_define_rule_name": name,
						"type":                         "1",
						"hash":                         "dacsda",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hash": "dacsda",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hash": "dacsda",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"platform": "all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform": "all",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_user_define_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_user_define_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hash": "acsdscsdsacd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hash": "acsdscsdsacd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type":                  "0",
					"platform":                     "all",
					"client_user_define_rule_name": name + "_update",
					"hash":                         "dacsda",
					"type":                         "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "all",
						"client_user_define_rule_name": name + "_update",
						"hash":                         "dacsda",
						"type":                         "1",
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

var AlicloudThreatDetectionClientUserDefineRuleMap4350 = map[string]string{
	"create_time": CHECKSET,
	"port_str":    CHECKSET,
}

func AlicloudThreatDetectionClientUserDefineRuleBasicDependence4350(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4404
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4404(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4404)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4404)
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
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"type":                         "3",
					"parent_cmdline":               "bash",
					"cmdline":                      "bash",
					"parent_proc_path":             "/root/bash",
					"proc_path":                    "/root/bash",
					"ip":                           "0.0.0.0",
					"port_str":                     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"type":                         "3",
						"parent_cmdline":               "bash",
						"cmdline":                      "bash",
						"parent_proc_path":             "/root/bash",
						"proc_path":                    "/root/bash",
						"ip":                           "0.0.0.0",
						"port_str":                     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_str": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_str": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip": "0.0.0.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip": "0.0.0.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"platform": "linux",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform": "linux",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_user_define_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_user_define_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_str": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_str": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip": "1.1.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip": "1.1.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name + "_update",
					"parent_proc_path":             "/root/bash",
					"port_str":                     "1",
					"type":                         "3",
					"cmdline":                      "bash",
					"proc_path":                    "/root/bash",
					"parent_cmdline":               "bash",
					"ip":                           "0.0.0.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name + "_update",
						"parent_proc_path":             "/root/bash",
						"port_str":                     "1",
						"type":                         "3",
						"cmdline":                      "bash",
						"proc_path":                    "/root/bash",
						"parent_cmdline":               "bash",
						"ip":                           "0.0.0.0",
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

var AlicloudThreatDetectionClientUserDefineRuleMap4404 = map[string]string{
	"create_time": CHECKSET,
	"port_str":    CHECKSET,
}

func AlicloudThreatDetectionClientUserDefineRuleBasicDependence4404(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4353
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4353(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4353)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4353)
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
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"type":                         "2",
					"cmdline":                      "bash",
					"proc_path":                    "/root/erfds/ef",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"type":                         "2",
						"cmdline":                      "bash",
						"proc_path":                    "/root/erfds/ef",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_proc_path": "/root/fscd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_proc_path": "/root/fscd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_str": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_str": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/root/erfds/ef",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/root/erfds/ef",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "/dcsaasc/cdsas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "/dcsaasc/cdsas",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"platform": "linux",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform": "linux",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_user_define_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_user_define_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_proc_path": "/fvs/fcds/fds",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_proc_path": "/fvs/fcds/fds",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/vsf/vfd/vfd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/vsf/vfd/vfd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name + "_update",
					"parent_proc_path":             "/root/fscd",
					"port_str":                     "1",
					"type":                         "2",
					"cmdline":                      "bash",
					"proc_path":                    "/root/erfds/ef",
					"parent_cmdline":               "/dcsaasc/cdsas",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name + "_update",
						"parent_proc_path":             "/root/fscd",
						"port_str":                     "1",
						"type":                         "2",
						"cmdline":                      "bash",
						"proc_path":                    "/root/erfds/ef",
						"parent_cmdline":               "/dcsaasc/cdsas",
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

var AlicloudThreatDetectionClientUserDefineRuleMap4353 = map[string]string{
	"create_time": CHECKSET,
	"port_str":    CHECKSET,
}

func AlicloudThreatDetectionClientUserDefineRuleBasicDependence4353(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4403
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4403(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4403)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4403)
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
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"type":                         "4",
					"cmdline":                      "bash",
					"proc_path":                    "/root/bash",
					"file_path":                    "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"type":                         "4",
						"cmdline":                      "bash",
						"proc_path":                    "/root/bash",
						"file_path":                    "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_str": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_str": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"platform": "linux",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform": "linux",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_user_define_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_user_define_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_path": "/root/bash2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path": "/root/bash2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_path":                    "/root/bash",
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name + "_update",
					"parent_proc_path":             "/root/bash",
					"port_str":                     "1",
					"type":                         "4",
					"cmdline":                      "bash",
					"proc_path":                    "/root/bash",
					"parent_cmdline":               "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path":                    "/root/bash",
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name + "_update",
						"parent_proc_path":             "/root/bash",
						"port_str":                     "1",
						"type":                         "4",
						"cmdline":                      "bash",
						"proc_path":                    "/root/bash",
						"parent_cmdline":               "bash",
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

var AlicloudThreatDetectionClientUserDefineRuleMap4403 = map[string]string{
	"create_time": CHECKSET,
	"port_str":    CHECKSET,
}

func AlicloudThreatDetectionClientUserDefineRuleBasicDependence4403(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4407
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4407(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4407)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4407)
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
					"action_type":                  "0",
					"client_user_define_rule_name": name,
					"type":                         "5",
					"platform":                     "windows",
					"proc_path":                    "/root/bash",
					"registry_key":                 "123",
					"registry_content":             "123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"client_user_define_rule_name": name,
						"type":                         "5",
						"platform":                     "windows",
						"proc_path":                    "/root/bash",
						"registry_key":                 "123",
						"registry_content":             "123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"new_file_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"new_file_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/root/bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/root/bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "bash",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"platform": "linux",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform": "linux",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_user_define_rule_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_user_define_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_path": "/root/bash1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path": "/root/bash1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"new_file_path": "/root/bash2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"new_file_path": "/root/bash2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cmdline": "bash d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cmdline": "bash d",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proc_path": "/root/bash3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proc_path": "/root/bash3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parent_cmdline": "bash d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parent_cmdline": "bash d",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_path":                    "/root/bash",
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name + "_update",
					"parent_proc_path":             "/root/bash",
					"new_file_path":                "/root/bash",
					"cmdline":                      "bash",
					"proc_path":                    "/root/bash",
					"parent_cmdline":               "bash",
					"type":                         "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path":                    "/root/bash",
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name + "_update",
						"parent_proc_path":             "/root/bash",
						"new_file_path":                "/root/bash",
						"cmdline":                      "bash",
						"proc_path":                    "/root/bash",
						"parent_cmdline":               "bash",
						"type":                         "5",
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

var AlicloudThreatDetectionClientUserDefineRuleMap4407 = map[string]string{
	"create_time": CHECKSET,
	"port_str":    CHECKSET,
}

func AlicloudThreatDetectionClientUserDefineRuleBasicDependence4407(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4405  twin
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4405_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4405)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4405)
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
					"action_type":                  "0",
					"platform":                     "windows",
					"registry_content":             "123sada",
					"client_user_define_rule_name": name,
					"parent_proc_path":             "/root/bash/dsa",
					"type":                         "5",
					"cmdline":                      "bash b",
					"proc_path":                    "/root/bash/d",
					"parent_cmdline":               "bash b",
					"registry_key":                 "123ads",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "windows",
						"registry_content":             "123sada",
						"client_user_define_rule_name": name,
						"parent_proc_path":             "/root/bash/dsa",
						"type":                         "5",
						"cmdline":                      "bash b",
						"proc_path":                    "/root/bash/d",
						"parent_cmdline":               "bash b",
						"registry_key":                 "123ads",
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

// Case 4350  twin
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4350_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4350)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4350)
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
					"action_type":                  "0",
					"platform":                     "all",
					"client_user_define_rule_name": name,
					"hash":                         "acsdscsdsacd",
					"type":                         "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "all",
						"client_user_define_rule_name": name,
						"hash":                         "acsdscsdsacd",
						"type":                         "1",
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

// Case 4404  twin
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4404_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4404)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4404)
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
					"action_type":                  "1",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"parent_proc_path":             "/root/bash",
					"port_str":                     "5",
					"type":                         "3",
					"cmdline":                      "bash",
					"proc_path":                    "/root/bash",
					"parent_cmdline":               "bash",
					"ip":                           "1.1.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "1",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"parent_proc_path":             "/root/bash",
						"port_str":                     "5",
						"type":                         "3",
						"cmdline":                      "bash",
						"proc_path":                    "/root/bash",
						"parent_cmdline":               "bash",
						"ip":                           "1.1.1.1",
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

// Case 4353  twin
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4353_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4353)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4353)
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
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"parent_proc_path":             "/fvs/fcds/fds",
					"port_str":                     "1",
					"type":                         "2",
					"cmdline":                      "bash",
					"proc_path":                    "/vsf/vfd/vfd",
					"parent_cmdline":               "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"parent_proc_path":             "/fvs/fcds/fds",
						"port_str":                     "1",
						"type":                         "2",
						"cmdline":                      "bash",
						"proc_path":                    "/vsf/vfd/vfd",
						"parent_cmdline":               "bash",
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

// Case 4403  twin
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4403_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4403)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4403)
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
					"file_path":                    "/root/bash2",
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"parent_proc_path":             "/root/bash",
					"port_str":                     "1",
					"type":                         "4",
					"cmdline":                      "bash",
					"proc_path":                    "/root/bash",
					"parent_cmdline":               "bash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path":                    "/root/bash2",
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"parent_proc_path":             "/root/bash",
						"port_str":                     "1",
						"type":                         "4",
						"cmdline":                      "bash",
						"proc_path":                    "/root/bash",
						"parent_cmdline":               "bash",
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

// Case 4407  twin
func TestAccAliCloudThreatDetectionClientUserDefineRule_basic4407_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_client_user_define_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionClientUserDefineRuleMap4407)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionClientUserDefineRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionClientUserDefineRuleBasicDependence4407)
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
					"file_path":                    "/root/bash1",
					"action_type":                  "0",
					"platform":                     "linux",
					"client_user_define_rule_name": name,
					"parent_proc_path":             "/root/bash",
					"new_file_path":                "/root/bash2",
					"cmdline":                      "bash d",
					"proc_path":                    "/root/bash3",
					"parent_cmdline":               "bash d",
					"type":                         "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_path":                    "/root/bash1",
						"action_type":                  "0",
						"platform":                     "linux",
						"client_user_define_rule_name": name,
						"parent_proc_path":             "/root/bash",
						"new_file_path":                "/root/bash2",
						"cmdline":                      "bash d",
						"proc_path":                    "/root/bash3",
						"parent_cmdline":               "bash d",
						"type":                         "7",
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

// Test ThreatDetection ClientUserDefineRule. <<< Resource test cases, automatically generated.
