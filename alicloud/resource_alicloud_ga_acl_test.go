package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ga_acl",
		&resource.Sweeper{
			Name: "alicloud_ga_acl",
			F:    testSweepGaAcl,
		})
}

func testSweepGaAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListAcls"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeLarge
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
				log.Printf("[INFO] Skipping Ga Acl: %s", item["AclName"].(string))
				continue
			}
			action := "DeleteAcl"
			request := map[string]interface{}{
				"AclId": item["AclId"],
			}
			request["ClientToken"] = buildClientToken("DeleteAcl")
			_, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ga Acl (%s): %s", item["AclId"].(string), err)
			}
			log.Printf("[INFO] Delete Ga Acl success: %s ", item["AclId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudGaAcl_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_acl.default"
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudGaAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAclBasicDependence0)
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
					"address_ip_version": "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_ip_version": "IPv4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name": name,
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
					"acl_entries": []map[string]interface{}{
						{
							"entry":             "192.168.1.0/24",
							"entry_description": "tf-test1",
						},
						{
							"entry":             "192.168.3.0/24",
							"entry_description": "tf-test3",
						},
						{
							"entry":             "192.168.4.0/24",
							"entry_description": "tf-test4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_entries.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Acl",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Acl",
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

func TestAccAliCloudGaAcl_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_acl.default"
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudGaAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAclBasicDependence0)
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
					"address_ip_version": "IPv4",
					"acl_name":           name,
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"acl_entries": []map[string]interface{}{
						{
							"entry":             "192.168.1.0/24",
							"entry_description": "tf-test1",
						},
						{
							"entry":             "192.168.2.0/24",
							"entry_description": "tf-test2/",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Acl",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_ip_version": "IPv4",
						"acl_name":           name,
						"resource_group_id":  CHECKSET,
						"acl_entries.#":      "2",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Acl",
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

var AliCloudGaAclMap0 = map[string]string{
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
}

func AliCloudGaAclBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`)
}
