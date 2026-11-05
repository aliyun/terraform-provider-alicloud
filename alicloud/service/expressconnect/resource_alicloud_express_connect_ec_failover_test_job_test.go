package expressconnect_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnect"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ExpressConnect EcFailoverTestJob. >>> Resource test cases, automatically generated.
// Case 5403
func TestAccAliCloudExpressConnectEcFailoverTestJob_basic5403(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_ec_failover_test_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectEcFailoverTestJobMap5403)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectEcFailoverTestJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectecfailovertestjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectEcFailoverTestJobBasicDependence5403)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"job_type": "StartNow",
					"resource_id": []string{
						"${data.alicloud_express_connect_physical_connections.default.ids.0}", "${data.alicloud_express_connect_physical_connections.default.ids.1}"},
					"job_duration":              "1",
					"resource_type":             "PHYSICALCONNECTION",
					"ec_failover_test_job_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_type":                  "StartNow",
						"resource_id.#":             "2",
						"job_duration":              "1",
						"resource_type":             "PHYSICALCONNECTION",
						"ec_failover_test_job_name": name,
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

var AlicloudExpressConnectEcFailoverTestJobMap5403 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudExpressConnectEcFailoverTestJobBasicDependence5403(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}
`, name)
}

// Case 5369
func TestAccAliCloudExpressConnectEcFailoverTestJob_basic5369(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_ec_failover_test_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectEcFailoverTestJobMap5369)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectEcFailoverTestJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectecfailovertestjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectEcFailoverTestJobBasicDependence5369)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"job_type": "StartLater",
					"resource_id": []string{
						"${data.alicloud_express_connect_physical_connections.default.ids.0}", "${data.alicloud_express_connect_physical_connections.default.ids.1}"},
					"job_duration":              "2",
					"resource_type":             "PHYSICALCONNECTION",
					"ec_failover_test_job_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_type":                  "StartLater",
						"resource_id.#":             "2",
						"job_duration":              "2",
						"resource_type":             "PHYSICALCONNECTION",
						"ec_failover_test_job_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "meijian-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "meijian-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Init",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Init",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_id": []string{
						"${data.alicloud_express_connect_physical_connections.default.ids.1}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_id.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"job_duration": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_duration": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ec_failover_test_job_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ec_failover_test_job_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "meijian-test-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "meijian-test-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_id": []string{
						"${data.alicloud_express_connect_physical_connections.default.ids.0}", "${data.alicloud_express_connect_physical_connections.default.ids.1}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_id.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"job_duration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_duration": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ec_failover_test_job_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ec_failover_test_job_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Testing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Testing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
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

var AlicloudExpressConnectEcFailoverTestJobMap5369 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudExpressConnectEcFailoverTestJobBasicDependence5369(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}
`, name)
}

// Case 5403  twin
func TestAccAliCloudExpressConnectEcFailoverTestJob_basic5403_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_ec_failover_test_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectEcFailoverTestJobMap5403)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectEcFailoverTestJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectecfailovertestjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectEcFailoverTestJobBasicDependence5403)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "meijian-test",
					"job_type":    "StartNow",
					"resource_id": []string{
						"${data.alicloud_express_connect_physical_connections.default.ids.0}", "${data.alicloud_express_connect_physical_connections.default.ids.1}"},
					"job_duration":              "1",
					"resource_type":             "PHYSICALCONNECTION",
					"ec_failover_test_job_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":               "meijian-test",
						"job_type":                  "StartNow",
						"resource_id.#":             "2",
						"job_duration":              "1",
						"resource_type":             "PHYSICALCONNECTION",
						"ec_failover_test_job_name": name,
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

// Case 5369  twin
func TestAccAliCloudExpressConnectEcFailoverTestJob_basic5369_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_ec_failover_test_job.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectEcFailoverTestJobMap5369)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectEcFailoverTestJob")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectecfailovertestjob%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectEcFailoverTestJobBasicDependence5369)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "meijian-test-1",
					"job_type":    "StartLater",
					"resource_id": []string{
						"${data.alicloud_express_connect_physical_connections.default.ids.0}", "${data.alicloud_express_connect_physical_connections.default.ids.1}"},
					"job_duration":              "1",
					"resource_type":             "PHYSICALCONNECTION",
					"ec_failover_test_job_name": name,
					"status":                    "Testing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":               "meijian-test-1",
						"job_type":                  "StartLater",
						"resource_id.#":             "2",
						"job_duration":              "1",
						"resource_type":             "PHYSICALCONNECTION",
						"ec_failover_test_job_name": name,
						"status":                    "Testing",
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

// Test ExpressConnect EcFailoverTestJob. <<< Resource test cases, automatically generated.
