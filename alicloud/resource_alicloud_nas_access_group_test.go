package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_nas_access_group", &resource.Sweeper{
		Name: "alicloud_nas_access_group",
		F:    testSweepNasAccessGroup,
		Dependencies: []string{
			"alicloud_nas_mount_target",
		},
	})
}

func testSweepNasAccessGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var ag []nas.AccessGroup
	req := nas.CreateDescribeAccessGroupsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DescribeAccessGroups(req)
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving filesystem: %s", err)
		}
		resp, _ := raw.(*nas.DescribeAccessGroupsResponse)
		if resp == nil || len(resp.AccessGroups.AccessGroup) < 1 {
			break
		}
		ag = append(ag, resp.AccessGroups.AccessGroup...)

		if len(resp.AccessGroups.AccessGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, fs := range ag {

		id := fs.AccessGroupName
		AccessGroupType := fs.AccessGroupType
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(id), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping AccessGroup: %s (%s)", AccessGroupType, id)
			continue
		}
		log.Printf("[INFO] Deleting AccessGroup: %s (%s)", AccessGroupType, id)
		req := nas.CreateDeleteAccessGroupRequest()
		req.AccessGroupName = id
		_, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
			return nasClient.DeleteAccessGroup(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete AccessGroup (%s (%s)): %s", AccessGroupType, id, err)
		}
	}
	return nil
}

func TestAccAlicloudNas_AccessGroup_Upgrade(t *testing.T) {
	var v nas.AccessGroup
	rand := acctest.RandIntRange(10000, 999999)
	resourceID := "alicloud_nas_access_group.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		// module name
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasAccessGroupUpgradeConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_type": "Vpc",
						"description":       "tf-testAccNasConfigDescription",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"file_system_type":  "extreme",
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasAccessGroupUpgradeConfigUpdate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_type": "Vpc",
						"description":       "tf-testAccNasConfigDescriptionUpdate",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
						"file_system_type":  "extreme",
					}),
				),
			},
		},
	})

}

func TestAccAlicloudNas_AccessGroup_update(t *testing.T) {
	var v nas.AccessGroup
	rand := acctest.RandIntRange(10000, 999999)
	resourceID := "alicloud_nas_access_group.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		// module name
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasAccessGroupVpcConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "Vpc",
						"description": "tf-testAccNasConfigDescription",
						"name":        fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasAccessGroupConfigUpdate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "Vpc",
						"description": "tf-testAccNasConfigUpdateDescription",
						"name":        fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
					}),
				),
			},
		},
	})

}

func TestAccAlicloudNas_AccessGroup_Classicupdate(t *testing.T) {
	var v nas.AccessGroup
	rand := acctest.RandIntRange(10000, 999999)
	resourceID := "alicloud_nas_access_group.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		// module name
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasAccessGroupConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "Classic",
						"description": "tf-testAccNasConfigDescription",
						"name":        fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasAccessGroupConfigClassicUpdate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "Classic",
						"description": "tf-testAccNasConfigUpdateDescription",
						"name":        fmt.Sprintf("tf-testAccNasConfigName-%d", rand),
					}),
				),
			},
		},
	})

}

func TestAccAlicloudNas_AccessGroup_multi(t *testing.T) {
	var v nas.AccessGroup
	rand1 := acctest.RandIntRange(10000, 499999)
	rand2 := acctest.RandIntRange(50000, 499999)
	resourceID := "alicloud_nas_access_group.default.9"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasAccessGroupVpcConfigMulti(rand1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "Vpc",
						"description": "tf-testAccNasConfigDescription-2",
						"name":        fmt.Sprintf("tf-testAccNasConfigName-%d-9", rand1),
					}),
				),
			},
			{
				Config: testAccNasAccessGroupClassicConfigMulti(rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":        "Classic",
						"description": "tf-testAccNasConfigDescription-1",
						"name":        fmt.Sprintf("tf-testAccNasConfigName-%d-9", rand2),
					}),
				),
			},
		},
	})
}

func testAccCheckAccessGroupExists(n string, nas *nas.AccessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No NAS ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		nasService := NasService{client}
		instance, err := nasService.DescribeNasAccessGroup(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*nas = instance
		return nil
	}
}

func testAccCheckAccessGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_access_group" {
			continue
		}

		// Try to find the NAS
		instance, err := nasService.DescribeNasAccessGroup(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("NAS %s still exist", instance.AccessGroupName))
	}

	return nil
}

func testAccNasAccessGroupConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_access_group" "default" {
		name = "tf-testAccNasConfigName-%d"
		type = "Classic"
		description = "tf-testAccNasConfigDescription"
	}`, rand)
}

func testAccNasAccessGroupUpgradeConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_access_group" "default" {
		access_group_name = "tf-testAccNasConfigName-%d"
		access_group_type = "Vpc"
		description = "tf-testAccNasConfigDescription"
		file_system_type="extreme"
	}`, rand)
}

func testAccNasAccessGroupUpgradeConfigUpdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_access_group" "default" {
		access_group_name = "tf-testAccNasConfigName-%d"
		access_group_type = "Vpc"
		description = "tf-testAccNasConfigDescriptionUpdate"
		file_system_type="extreme"
	}`, rand)
}

func testAccNasAccessGroupVpcConfig(rand int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Vpc"
                description = "tf-testAccNasConfigDescription"
        }`, rand)
}

func testAccNasAccessGroupConfigUpdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_access_group" "default" {
		name = "tf-testAccNasConfigName-%d"
		type = "Vpc"
		description = "tf-testAccNasConfigUpdateDescription"
	}`, rand)
}

func testAccNasAccessGroupConfigClassicUpdate(rand int) string {
	return fmt.Sprintf(`
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d"
                type = "Classic"
                description = "tf-testAccNasConfigUpdateDescription"
        }`, rand)
}

func testAccNasAccessGroupVpcConfigMulti(rand int) string {
	return fmt.Sprintf(`
	variable "description" {
  		default = "tf-testAccNasConfigDescription"
	}
	resource "alicloud_nas_access_group" "default" {
		name = "tf-testAccNasConfigName-%d-${count.index}"
		type = "Vpc"
		description = "${var.description}-2"
		count = 10
	}`, rand)
}

func testAccNasAccessGroupClassicConfigMulti(rand int) string {
	return fmt.Sprintf(`
        variable "description" {
                default = "tf-testAccNasConfigDescription"
        }
        resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfigName-%d-${count.index}"
                type = "Classic"
                description = "${var.description}-1"
		count = 10
        }`, rand)
}
