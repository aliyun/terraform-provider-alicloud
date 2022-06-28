package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_alb_acl",
		&resource.Sweeper{
			Name: "alicloud_alb_acl",
			F:    testSweepAlbAcl,
		})
}

func testSweepAlbAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListAcls"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeXLarge

	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Acls", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Acls", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["AclName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["AclName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Alb Acl: %s", item["AclName"].(string))
				continue
			}
			action := "DeleteAcl"
			request := map[string]interface{}{
				"AclId": item["AclId"],
			}
			request["ClientToken"] = buildClientToken("DeleteAcl")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Alb Acl (%s): %s", item["AclId"].(string), err)
			}
			log.Printf("[INFO] Delete Alb Acl success: %s ", item["AclId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudALBAcl_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudALBAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccalbacl%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBAclBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_entries": []map[string]interface{}{
						{
							"description": name,
							"entry":       "10.0.0.0/24",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/23",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/25",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/26",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/27",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/28",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/29",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/10",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/11",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/12",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/13",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/14",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/15",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/16",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/17",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/18",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/19",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/10",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/11",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/12",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/13",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/14",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/15",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/16",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/17",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/18",
						},
						{
							"description": name,
							"entry":       "11.0.0.0/19",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/10",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/11",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/12",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/13",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/14",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/15",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/16",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/17",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/18",
						},
						{
							"description": name,
							"entry":       "12.0.0.0/19",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entries.#": "37",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_entries": []map[string]interface{}{
						{
							"description": name,
							"entry":       "10.0.0.0/23",
						},
						{
							"description": name,
							"entry":       "10.0.0.0/22",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entries.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
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
						"Created": "tfTestAcc7",
						"For":     "Tftestacc7",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc7",
						"tags.For":     "Tftestacc7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": "${var.name}",
					"tags": map[string]string{
						"Created": "tfTestAcc99",
						"For":     "Tftestacc99",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name":     name,
						"tags.%":       "2",
						"tags.Created": "tfTestAcc99",
						"tags.For":     "Tftestacc99",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudALBAcl_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudALBAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccalbacl%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBAclBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name":          "${var.name}",
					"dry_run":           "false",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name":          name,
						"dry_run":           "false",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudALBAclMap0 = map[string]string{
	"tags.%":            NOSET,
	"dry_run":           NOSET,
	"resource_group_id": CHECKSET,
	"acl_entries.#":     CHECKSET,
	"status":            CHECKSET,
}

func AlicloudALBAclBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
