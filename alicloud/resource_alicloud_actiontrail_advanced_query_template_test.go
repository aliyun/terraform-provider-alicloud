// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Actiontrail AdvancedQueryTemplate. >>> Resource test cases, automatically generated.
// Case AdvancedQueryTemplate线上-templateSql变更后测试 10937
func TestAccAliCloudActiontrailAdvancedQueryTemplate_basic10937(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_actiontrail_advanced_query_template.default"
	ra := resourceAttrInit(resourceId, AlicloudActiontrailAdvancedQueryTemplateMap10937)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailAdvancedQueryTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudActiontrailAdvancedQueryTemplateBasicDependence10937)
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
					"simple_query":  "true",
					"template_name": "testTemplateName",
					"template_sql":  "*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"simple_query":  "true",
						"template_name": "testTemplateName",
						"template_sql":  "*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"simple_query":  "false",
					"template_name": "newTemplateName",
					"template_sql":  "* AND (event.userIdentity.accessKeyId: xxxklnkl)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"simple_query":  "false",
						"template_name": "newTemplateName",
						"template_sql":  "* AND (event.userIdentity.accessKeyId: xxxklnkl)",
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

var AlicloudActiontrailAdvancedQueryTemplateMap10937 = map[string]string{}

func AlicloudActiontrailAdvancedQueryTemplateBasicDependence10937(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case AdvancedQueryTemplate线上 10764
func TestAccAliCloudActiontrailAdvancedQueryTemplate_basic10764(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_actiontrail_advanced_query_template.default"
	ra := resourceAttrInit(resourceId, AlicloudActiontrailAdvancedQueryTemplateMap10764)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailAdvancedQueryTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudActiontrailAdvancedQueryTemplateBasicDependence10764)
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
					"simple_query":  "true",
					"template_name": "testTemplateName",
					"template_sql":  "*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"simple_query":  "true",
						"template_name": "testTemplateName",
						"template_sql":  "*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"simple_query":  "false",
					"template_name": "newTemplateName",
					"template_sql":  "* AND (event.userIdentity.accessKeyId: xxxklnkl)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"simple_query":  "false",
						"template_name": "newTemplateName",
						"template_sql":  "* AND (event.userIdentity.accessKeyId: xxxklnkl)",
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

var AlicloudActiontrailAdvancedQueryTemplateMap10764 = map[string]string{}

func AlicloudActiontrailAdvancedQueryTemplateBasicDependence10764(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case AdvancedQueryTemplate测试用例 10763
func TestAccAliCloudActiontrailAdvancedQueryTemplate_basic10763(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_actiontrail_advanced_query_template.default"
	ra := resourceAttrInit(resourceId, AlicloudActiontrailAdvancedQueryTemplateMap10763)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailAdvancedQueryTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudActiontrailAdvancedQueryTemplateBasicDependence10763)
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
					"simple_query":  "true",
					"template_name": "testTemplateName",
					"template_sql":  "*",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"simple_query":  "true",
						"template_name": "testTemplateName",
						"template_sql":  "*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"simple_query":  "false",
					"template_name": "newTemplateName",
					"template_sql":  "* AND (event.userIdentity.accessKeyId: xxxklnkl)",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"simple_query":  "false",
						"template_name": "newTemplateName",
						"template_sql":  "* AND (event.userIdentity.accessKeyId: xxxklnkl)",
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

var AlicloudActiontrailAdvancedQueryTemplateMap10763 = map[string]string{}

func AlicloudActiontrailAdvancedQueryTemplateBasicDependence10763(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case AdvancedQueryTemplate资源测试 5151
func TestAccAliCloudActiontrailAdvancedQueryTemplate_basic5151(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_actiontrail_advanced_query_template.default"
	ra := resourceAttrInit(resourceId, AlicloudActiontrailAdvancedQueryTemplateMap5151)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ActiontrailServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeActiontrailAdvancedQueryTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccactiontrail%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudActiontrailAdvancedQueryTemplateBasicDependence5151)
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
					"template_sql": "serviceName:actionTrail",
					"simple_query": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_sql": "serviceName:actionTrail",
						"simple_query": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_sql": "serviceName:ecs",
					"simple_query": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_sql": "serviceName:ecs",
						"simple_query": "false",
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

var AlicloudActiontrailAdvancedQueryTemplateMap5151 = map[string]string{}

func AlicloudActiontrailAdvancedQueryTemplateBasicDependence5151(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Actiontrail AdvancedQueryTemplate. <<< Resource test cases, automatically generated.
