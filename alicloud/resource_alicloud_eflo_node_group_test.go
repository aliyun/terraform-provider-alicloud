package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEfloNodeGroup_basic10344(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloNodeGroupMap10344)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloNodeGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloNodeGroupBasicDependence10344)
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
					"node_group_description": "resource-test1",
					"node_group_name":        name,
					"cluster_id":             "i118078301742281630607",
					"machine_type":           "efg2.C48cA3sen",
					"az":                     "cn-hangzhou-b",
					"image_id":               "i198448731735114628708",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_group_description": "resource-test1",
						"node_group_name":        name,
						"cluster_id":             "i118078301742281630607",
						"machine_type":           "efg2.C48cA3sen",
						"az":                     "cn-hangzhou-b",
						"image_id":               "i198448731735114628708",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_group_name": name + "_update",
					"nodes": []map[string]interface{}{
						{
							"login_password": "Alibaba@2025",
							"node_id":        "e01-cn-rno46i6rdfn",
							"vpc_id":         "vpc-bp15z92ev9jflxq14c2l0",
							"vswitch_id":     "vsw-bp1lii10s3tl99bpx20mo",
							"hostname":       "jxyhostname",
						},
						//{
						//	"login_password": "Alibaba@2025",
						//	"node_id":        "e01-cn-rno46i6rdfn-mock",
						//	"vpc_id":         "vpc-bp15z92ev9jflxq14c2l0",
						//	"vswitch_id":     "vsw-bp1lii10s3tl99bpx20mo",
						//	"hostname":       "jxyhostname",
						//},
					},
					"ignore_failed_node_tasks": "true",
					"zone_id":                  "cn-hangzhou-b",
					"user_data":                "YWxpLGFsaSxhbGliYWJh",
					"vpd_subnets": []string{
						"test"},
					"vswitch_zone_id": "cn-hangzhou-b",
					"ip_allocation_policy": []map[string]interface{}{
						{
							"bond_policy": []map[string]interface{}{
								{
									"bond_default_subnet": "test",
									"bonds": []map[string]interface{}{
										{
											"name":   "test",
											"subnet": "test",
										},
									},
								},
							},
							"machine_type_policy": []map[string]interface{}{
								{
									"bonds": []map[string]interface{}{
										{
											"name":   "test",
											"subnet": "test",
										},
									},
									"machine_type": "test",
								},
							},
							"node_policy": []map[string]interface{}{
								{
									"bonds": []map[string]interface{}{
										{
											"name":   "test",
											"subnet": "test",
										},
									},
									"node_id": "e01-cn-rno46i6rdfn",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_group_name":          name + "_update",
						"nodes.#":                  "1",
						"ignore_failed_node_tasks": "true",
						"zone_id":                  "cn-hangzhou-b",
						"user_data":                "YWxpLGFsaSxhbGliYWJh",
						"vpd_subnets.#":            "1",
						"vswitch_zone_id":          "cn-hangzhou-b",
						"ip_allocation_policy.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nodes": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nodes.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_failed_node_tasks", "ip_allocation_policy", "user_data", "vswitch_zone_id", "vpd_subnets", "zone_id"},
			},
		},
	})
}

var AlicloudEfloNodeGroupMap10344 = map[string]string{
	"create_time":   CHECKSET,
	"node_group_id": CHECKSET,
}

func AlicloudEfloNodeGroupBasicDependence10344(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
