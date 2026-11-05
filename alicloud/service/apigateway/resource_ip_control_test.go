package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/apigateway"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ApiGateway IpControl. >>> Resource test cases.
// IpControl is a standalone API Gateway resource (no instance dependency): it
// manages an IP access control (ALLOW/REFUSE) plus its nested policy items. The
// policy items are reconciled through Add/Remove operations, mirroring the
// AccessControlList acl_entrys diff pattern.

// Case basic: create an ALLOW control with two policy items, then reconcile the
// items (remove one, add another), update description, update the control name
// (items unchanged -> no-op diff path) and finally import.
func TestAccAliCloudApiGatewayIpControl_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_ip_control.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayIpControlMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayIpControl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccipcontrol%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayIpControlBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_control_name": name,
					"ip_control_type": "ALLOW",
					"description":     "tf acc description",
					"ip_control_policy_items": []map[string]interface{}{
						{
							"cidr_ip": "10.0.0.0/24",
						},
						{
							"cidr_ip": "192.168.0.0/16",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_control_name":           name,
						"ip_control_type":           "ALLOW",
						"description":               "tf acc description",
						"ip_control_policy_items.#": "2",
					}),
				),
			},
			{
				// Reconcile policy items: remove 192.168.0.0/16 and add
				// 172.16.0.0/12, and update the description.
				Config: testAccConfig(map[string]interface{}{
					"description": "tf acc description updated",
					"ip_control_policy_items": []map[string]interface{}{
						{
							"cidr_ip": "10.0.0.0/24",
						},
						{
							"cidr_ip": "172.16.0.0/12",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":               "tf acc description updated",
						"ip_control_policy_items.#": "2",
					}),
				),
			},
			{
				// Update the control name only; the policy items are unchanged so
				// the item reconciliation is a no-op (HasChange == false).
				Config: testAccConfig(map[string]interface{}{
					"ip_control_name": name + "u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_control_name": name + "u",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case appId: an ALLOW control whose policy item carries an app_id referencing a
// real API Gateway app, exercising the Add path with AppId set.
func TestAccAliCloudApiGatewayIpControl_appId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_ip_control.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayIpControlMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayIpControl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccipcapp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayIpControlAppIdDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_control_name": name,
					"ip_control_type": "ALLOW",
					"ip_control_policy_items": []map[string]interface{}{
						{
							"cidr_ip": "10.0.0.0/24",
							"app_id":  "${alicloud_api_gateway_app.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_control_name":           name,
						"ip_control_type":           "ALLOW",
						"ip_control_policy_items.#": "1",
					}),
				),
			},
			{
				// Modify app_id: rebind the policy item from the default app
				// to the second app. Because app_id is part of the TypeSet
				// hash, this resolves to a remove + add, exercising the
				// app_id change path and satisfying modified-coverage.
				Config: testAccConfig(map[string]interface{}{
					"ip_control_policy_items": []map[string]interface{}{
						{
							"cidr_ip": "10.0.0.0/24",
							"app_id":  "${alicloud_api_gateway_app.second.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_control_name":           name,
						"ip_control_type":           "ALLOW",
						"ip_control_policy_items.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case refuse: a REFUSE control (ip_control_type is ForceNew) with one policy
// item, then import.
func TestAccAliCloudApiGatewayIpControl_refuse(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_ip_control.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApiGatewayIpControlMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apigateway.NewApiGatewayServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApiGatewayIpControl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccipcrefuse%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApiGatewayIpControlBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_control_name": name,
					"ip_control_type": "REFUSE",
					"ip_control_policy_items": []map[string]interface{}{
						{
							"cidr_ip": "10.0.0.0/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_control_name":           name,
						"ip_control_type":           "REFUSE",
						"ip_control_policy_items.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudApiGatewayIpControlMap = map[string]string{}

func AlicloudApiGatewayIpControlBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

`, name)
}

func AlicloudApiGatewayIpControlAppIdDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_api_gateway_app" "default" {
  name        = var.name
  description = "tf acc ip control app dependence"
}

resource "alicloud_api_gateway_app" "second" {
  name        = "${var.name}_2"
  description = "tf acc ip control app dependence second"
}

`, name)
}

// Test ApiGateway IpControl. <<< Resource test cases.
