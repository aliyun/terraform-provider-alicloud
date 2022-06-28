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
		"alicloud_ga_acl",
		&resource.Sweeper{
			Name: "alicloud_ga_acl",
			F:    testSweepGaAcl,
		})
}

func testSweepGaAcl(region string) error {
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

	request["MaxResults"] = PageSizeLarge
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func TestAccAlicloudGaAcl_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_acl.default"
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGaAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAclBasicDependence0)
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
					"acl_entries": []map[string]interface{}{
						{
							"entry":             "192.168.1.0/24",
							"entry_description": "tf-test1",
						},
						{
							"entry":             "192.168.2.0/24",
							"entry_description": "tf-test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name":      name,
						"acl_entries.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name": name + "update",
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
					"acl_name": name,
					"acl_entries": []map[string]interface{}{
						{
							"entry":             "192.168.1.0/24",
							"entry_description": "tf-test1",
						},
						{
							"entry":             "192.168.2.0/24",
							"entry_description": "tf-test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name":      name,
						"acl_entries.#": "2",
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

func TestAccAlicloudGaAcl_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_acl.default"
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGaAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAclBasicDependence0)
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
					"acl_name":           name,
					"address_ip_version": "IPv4",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_name":           name,
						"address_ip_version": "IPv4",
						"acl_entries.#":      "2",
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

var AlicloudGaAclMap0 = map[string]string{}

func AlicloudGaAclBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
