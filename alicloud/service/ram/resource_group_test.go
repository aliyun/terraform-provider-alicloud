package ram_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_group", &resource.Sweeper{
		Name: "alicloud_ram_group",
		F:    testSweepRamGroups,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_ram_user",
		},
	})
}

func testSweepRamGroups(region string) error {
	rawClient, err := ramTestClientForRegion(region)
	if err != nil {
		return errs.WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ramsdk.Group
	request := ramsdk.CreateListGroupsRequest()
	for {
		raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListGroups(request)
		})
		if err != nil {
			return errs.WrapError(err)
		}
		resp, _ := raw.(*ramsdk.ListGroupsResponse)
		if len(resp.Groups.Group) < 1 {
			break
		}
		groups = append(groups, resp.Groups.Group...)

		if !resp.IsTruncated {
			break
		}
		request.Marker = resp.Marker
	}
	sweeped := false

	for _, v := range groups {
		name := v.GroupName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram Group: %s", name)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ram Group: %s", name)
		request := ramsdk.CreateListPoliciesForGroupRequest()
		request.GroupName = name

		raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to list Ram Group (%s): %s", name, err)
		}
		response, _ := raw.(*ramsdk.ListPoliciesForGroupResponse)
		for _, p := range response.Policies.Policy {
			request := ramsdk.CreateDetachPolicyFromGroupRequest()
			request.PolicyType = p.PolicyType
			request.GroupName = name
			request.PolicyName = p.PolicyName
			log.Printf("[INFO] Detaching Ram policy %s from group: %s", p.PolicyName, name)
			_, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.DetachPolicyFromGroup(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to detach policy from Group (%s): %s", name, err)
			}
		}
		_, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			request := ramsdk.CreateDeleteGroupRequest()
			request.GroupName = name
			return ramClient.DeleteGroup(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ram Group (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

// Test Ram Group. >>> Resource test cases, automatically generated.
// Case Group资源测试_副本1737429980161 10096
func TestAccAliCloudRamGroup_basic10096(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRamGroupMap10096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRamGroupBasicDependence10096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comments": "for test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comments": "for test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudRamGroup_basic10096_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRamGroupMap10096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRamGroupBasicDependence10096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name,
					"comments":   "this is a policy test",
					"force":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name,
						"comments":   "this is a policy test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

// Case Group资源测试_副本1737429980161 10098 适配废弃字段name
func TestAccAliCloudRamGroup_basic10098(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRamGroupMap10096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRamGroupBasicDependence10096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comments": "for test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comments": "for test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudRamGroup_basic10098_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRamGroupMap10096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRamGroupBasicDependence10096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     name,
					"comments": "this is a policy test",
					"force":    "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"comments": "this is a policy test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AliCloudRamGroupMap10096 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AliCloudRamGroupBasicDependence10096(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

func TestAccAliCloudRamGroup_basic_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_group.default.2"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRamGroupMap10096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRamGroupBasicDependence10096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":      "3",
					"group_name": name + "-${count.index}",
					"comments":   "this is a policy test",
					"force":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + fmt.Sprint(-2),
						"comments":   "this is a policy test",
					}),
				),
			},
		},
	})
}

// Test Ram Group. <<< Resource test cases, automatically generated.
