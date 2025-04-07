package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA ScheduledPreloadJob. >>> Resource test cases, automatically generated.
// Case scheduledpreloadjob_test_1
func TestAccAliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_scheduled_preload_job.default"
	ra := resourceAttrInit(resourceId, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1Map)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaScheduledPreloadJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAScheduledPreloadJob%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1BasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"insert_way":                 "oss",
					"site_id":                    "${alicloud_esa_site.resource_Site_ScheduledPreloadJob_test_1.id}",
					"oss_url":                    "https://yandanpub.oss-cn-hangzhou.aliyuncs.com/1.txt",
					"scheduled_preload_job_name": "testscheduledpreloadjob",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"oss_url"},
			},
		},
	})
}

var AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1Map = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1BasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_ScheduledPreloadJob_test_1" {
  site_name   = "terraform.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case scheduledpreloadjob_test_2
func TestAccAliCloudESAScheduledPreloadJobscheduledpreloadjob_test_2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_scheduled_preload_job.default"
	ra := resourceAttrInit(resourceId, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_2Map)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaScheduledPreloadJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAScheduledPreloadJob%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_2BasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"insert_way":                 "textBox",
					"site_id":                    "${alicloud_esa_site.resource_Site_ScheduledPreloadJob_test_2.id}",
					"scheduled_preload_job_name": "testscheduledpreloadjob",
					"url_list":                   "http://test1.gositecdn.cn/test/test.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"url_list"},
			},
		},
	})
}

var AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_2Map = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_2BasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_ScheduledPreloadJob_test_2" {
  site_name   = "terraform.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_ScheduledPreloadJob_ScheduledPreloadExecution_test
func TestAccAliCloudESAScheduledPreloadJobresource_ScheduledPreloadJob_ScheduledPreloadExecution_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_scheduled_preload_job.default"
	ra := resourceAttrInit(resourceId, AliCloudESAScheduledPreloadJobresource_ScheduledPreloadJob_ScheduledPreloadExecution_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaScheduledPreloadJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAScheduledPreloadJob%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAScheduledPreloadJobresource_ScheduledPreloadJob_ScheduledPreloadExecution_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"insert_way":                 "textBox",
					"site_id":                    "${alicloud_esa_site.resource_Site_ScheduledPreloadExecution_test.id}",
					"scheduled_preload_job_name": "test_scheduledpreloadexecution_job",
					"url_list":                   "http://test1.gositecdn.cn/test/test.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
				ImportStateVerifyIgnore: []string{"url_list"},
			},
		},
	})
}

var AliCloudESAScheduledPreloadJobresource_ScheduledPreloadJob_ScheduledPreloadExecution_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAScheduledPreloadJobresource_ScheduledPreloadJob_ScheduledPreloadExecution_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_ScheduledPreloadExecution_test" {
  site_name   = "terraform.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA ScheduledPreloadJob. <<< Resource test cases, automatically generated.
