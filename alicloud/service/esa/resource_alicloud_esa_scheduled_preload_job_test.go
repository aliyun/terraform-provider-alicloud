package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA ScheduledPreloadJob. >>> Resource test cases, automatically generated.
// Case scheduledpreloadjob_test_1
func TestAccAliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_scheduled_preload_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1Map)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaScheduledPreloadJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAScheduledPreloadJob%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_1BasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"id": acctest.CHECKSET,
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
  site_name   = "${lower(var.name)}.com"
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_2Map)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaScheduledPreloadJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAScheduledPreloadJob%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAScheduledPreloadJobscheduledpreloadjob_test_2BasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"id": acctest.CHECKSET,
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
  site_name   = "${lower(var.name)}.com"
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAScheduledPreloadJobresource_ScheduledPreloadJob_ScheduledPreloadExecution_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaScheduledPreloadJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAScheduledPreloadJob%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAScheduledPreloadJobresource_ScheduledPreloadJob_ScheduledPreloadExecution_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"id": acctest.CHECKSET,
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
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA ScheduledPreloadJob. <<< Resource test cases, automatically generated.
