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
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group" {
			continue
		}

		_, err := ecsService.DescribeSecurityGroup(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return WrapError(Error("Error SecurityGroup still exist"))
	}
	return nil
}

func TestAccAliCloudECSSecurityGroup_basic0(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSSecurityGroupBasicDependence0)
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "SecurityGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "SecurityGroup",
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

func TestAccAliCloudECSSecurityGroup_basic0_twin(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSSecurityGroupBasicDependence0)
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
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"security_group_type": "normal",
					"name":                name,
					"description":         name,
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"inner_access_policy": "Drop",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "SecurityGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":              CHECKSET,
						"security_group_type": "normal",
						"name":                name,
						"description":         name,
						"resource_group_id":   CHECKSET,
						"inner_access_policy": "Drop",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "SecurityGroup",
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

func TestAccAliCloudECSSecurityGroup_basic1(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSSecurityGroupBasicDependence0)
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "SecurityGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "SecurityGroup",
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

func TestAccAliCloudECSSecurityGroup_basic1_twin(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSSecurityGroupBasicDependence0)
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
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"security_group_type": "normal",
					"name":                name,
					"description":         name,
					"resource_group_id":   "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"inner_access":        "false",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "SecurityGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":              CHECKSET,
						"security_group_type": "normal",
						"name":                name,
						"description":         name,
						"resource_group_id":   CHECKSET,
						"inner_access":        "false",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "SecurityGroup",
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

func TestAccAliCloudECSSecurityGroupMulti(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alicloud_security_group.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSecurityGroupName%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityGroupConfigMulti(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupConfigMulti(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = "${var.name}_vpc"
  		cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_security_group" "default" {
  		count             = 10
  		vpc_id            = alicloud_vpc.default.id
  		resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.1.id
  		inner_access      = false
  		name              = var.name
  		description       = "${var.name}_describe"
  		tags = {
    		Created  = "TF"
    		For = "SecurityGroup"
  		}
	}
`, name)
}

var testAccCheckSecurityBasicMap = map[string]string{
	"security_group_type": CHECKSET,
	"inner_access_policy": CHECKSET,
	"inner_access":        CHECKSET,
}

func AliCloudECSSecurityGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}
`, name)
}
