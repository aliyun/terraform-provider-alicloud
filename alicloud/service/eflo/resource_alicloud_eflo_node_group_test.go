package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass ip_allocation_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.bond_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.machine_type_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.machine_type_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.node_policy -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass ip_allocation_policy.node_policy.bonds -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change
// collection-order: bypass vpd_subnets -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEfloNodeGroup_basic10344(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_node_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloNodeGroupMap10344)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloNodeGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloNodeGroupBasicDependence10344)
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
					"nodes": acctest.REMOVEKEY,
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
	"create_time":   acctest.CHECKSET,
	"node_group_id": acctest.CHECKSET,
}

func AlicloudEfloNodeGroupBasicDependence10344(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
