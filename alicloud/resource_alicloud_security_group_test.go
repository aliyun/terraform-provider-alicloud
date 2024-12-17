package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_security_group", &resource.Sweeper{
		Name: "alicloud_security_group",
		F:    testSweepSecurityGroups,
		//When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
			"alicloud_ecs_network_interface",
			"alicloud_bastionhost_instance",
			"alicloud_cs_kubernetes",
		},
	})
}

func testSweepSecurityGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ecs.SecurityGroup
	req := ecs.CreateDescribeSecurityGroupsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroups(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Security Groups: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
		if resp == nil || len(resp.SecurityGroups.SecurityGroup) < 1 {
			break
		}
		groups = append(groups, resp.SecurityGroups.SecurityGroup...)

		if len(resp.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	vpcService := VpcService{client}
	ecsService := EcsService{client}
	for _, v := range groups {
		name := v.SecurityGroupName
		id := v.SecurityGroupId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a Security Group created by other service, it should be fetched by vpc name and deleted.
			if skip {
				if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
					skip = !need
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Security Group: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Deleting Security Group: %s (%s)", name, id)
		if err := ecsService.sweepSecurityGroup(id); err != nil {
			log.Printf("[ERROR] Failed to delete Security Group (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

// Test Ecs SecurityGroup. >>> Resource test cases, automatically generated.
// Case 全生命周期 8588
func TestAccAliCloudEcsSecurityGroup_basic8588(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSecurityGroupMap8588)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSecurityGroupBasicDependence8588)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inner_access_policy": "Drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inner_access_policy": "Drop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inner_access_policy": "Accept",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inner_access_policy": "Accept",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inner_access_policy": "Drop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inner_access_policy": "Drop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudEcsSecurityGroup_basic8588_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSecurityGroupMap8588)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSecurityGroupBasicDependence8588)
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
					"description":         name,
					"inner_access_policy": "Drop",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"security_group_name": name,
					"security_group_type": "normal",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         name,
						"inner_access_policy": "Drop",
						"resource_group_id":   CHECKSET,
						"security_group_name": name,
						"security_group_type": "normal",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Test",
						"vpc_id":              CHECKSET,
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

// Case 全生命周期, 适配废弃字段name, inner_access 8600
func TestAccAliCloudEcsSecurityGroup_basic8600(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSecurityGroupMap8588)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSecurityGroupBasicDependence8588)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
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
					"inner_access": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inner_access": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inner_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inner_access": "true",
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

func TestAccAliCloudEcsSecurityGroup_basic8600_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSecurityGroupMap8588)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSecurityGroupBasicDependence8588)
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
					"description":         name,
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"security_group_type": "normal",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"vpc_id":       "${data.alicloud_vpcs.default.ids.0}",
					"name":         name,
					"inner_access": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         name,
						"resource_group_id":   CHECKSET,
						"security_group_type": "normal",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Test",
						"vpc_id":              CHECKSET,
						"name":                name,
						"inner_access":        "false",
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

func TestAccAliCloudEcsSecurityGroup_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_security_group.default.5"
	ra := resourceAttrInit(resourceId, AliCloudEcsSecurityGroupMap8588)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSecurityGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSecurityGroupBasicDependence8588)
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
					"count":               "6",
					"description":         name + "-${count.index}",
					"inner_access_policy": "Drop",
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"security_group_name": name + "-${count.index}",
					"security_group_type": "normal",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"vpc_id": "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         name + fmt.Sprint(-5),
						"inner_access_policy": "Drop",
						"resource_group_id":   CHECKSET,
						"security_group_name": name + fmt.Sprint(-5),
						"security_group_type": "normal",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Test",
						"vpc_id":              CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudEcsSecurityGroupMap8588 = map[string]string{
	"create_time":         CHECKSET,
	"inner_access_policy": CHECKSET,
	"security_group_type": CHECKSET,
	"vpc_id":              CHECKSET,
}

func AliCloudEcsSecurityGroupBasicDependence8588(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
`, name)
}

// Test Ecs SecurityGroup. <<< Resource test cases, automatically generated.
