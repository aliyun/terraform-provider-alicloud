package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_slb_acl", &resource.Sweeper{
		Name: "alicloud_slb_acl",
		F:    testSweepSlbAcl,
	})
}

func testSweepSlbAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	req := slb.CreateDescribeAccessControlListsRequest()
	req.RegionId = client.RegionId
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlLists(req)
	})
	if err != nil {
		return err
	}
	resp, _ := raw.(*slb.DescribeAccessControlListsResponse)

	for _, acl := range resp.Acls.Acl {
		name := acl.AclName
		id := acl.AclId

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Slb Acl: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting Slb Acl : %s (%s)", name, id)
		req := slb.CreateDeleteAccessControlListRequest()
		req.AclId = id
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteAccessControlList(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Slb Acl (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSlbAcl_basic(t *testing.T) {
	var acl *slb.DescribeAccessControlListAttributeResponse

	resourceId := "alicloud_slb_acl.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &acl, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbAclBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSLBAclDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       "${var.name}",
					"ip_version": "${var.ip_version}",
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.basic_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbAcl_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
						{
							"entry":   "172.10.10.0/24",
							"comment": "third",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"entry_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl",
						"entry_list.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbAcl_muilt(t *testing.T) {
	var acl *slb.DescribeAccessControlListAttributeResponse

	resourceId := "alicloud_slb_acl.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &acl, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbAclMuilt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSLBAclDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       "${var.name}-${count.index}",
					"count":      "${var.number}",
					"ip_version": "${var.ip_version}",
					"entry_list": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": "80",
						},
						{
							"entry":   "168.10.10.0/24",
							"comment": "second",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf-testAccSlbAcl-9",
						"ip_version":   "ipv4",
						"entry_list.#": "2",
					}),
				),
			},
		},
	})
}

func resourceSLBAclDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbAcl"
}
variable "basic_name" {
  default = "tf-testAccSlbAcl_name"
}
variable "number" {
  default = "10"
}
variable "ip_version" {
  default = "ipv4"
}
`)
}
