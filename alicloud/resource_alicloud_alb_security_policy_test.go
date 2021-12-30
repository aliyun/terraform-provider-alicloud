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
		"alicloud_alb_security_policy",
		&resource.Sweeper{
			Name: "alicloud_alb_security_policy",
			F:    testSweepAlbSecurityPolicy,
		})
}

func testSweepAlbSecurityPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListSecurityPolicies"
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

		resp, err := jsonpath.Get("$.SecurityPolicies", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.SecurityPolicies", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["SecurityPolicyName"]; !ok {
				continue
			}

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["SecurityPolicyName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ALB Security Policy: %s", item["SecurityPolicyName"].(string))
				continue
			}
			action := "DeleteSecurityPolicy"
			request := map[string]interface{}{
				"SecurityPolicyId": item["SecurityPolicyId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete ALB Security Policy (%s): %s", item["SecurityPolicyId"].(string), err)
			}
			log.Printf("[INFO] Delete ALB Security Policy success: %s ", item["SecurityPolicyId"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudALBSecurityPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_security_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudALBSecurityPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbSecurityPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbsecuritypolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBSecurityPolicyBasicDependence0)
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
					"security_policy_name": "tf-testAcc-test-secrity",
					"tls_versions":         []string{"TLSv1.0"},
					"ciphers":              []string{"ECDHE-ECDSA-AES128-SHA", "AES256-SHA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": "tf-testAcc-test-secrity",
						"ciphers.#":            "2",
						"tls_versions.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "tf-testAcc-test-secrity-new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": "tf-testAcc-test-secrity-new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{"TLSv1.1"},
					"ciphers":      []string{"AES128-SHA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#":      "1",
						"tls_versions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc4",
						"For":     "Tftestacc4",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc4",
						"tags.For":     "Tftestacc4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc42",
						"For":     "Tftestacc42",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc42",
						"tags.For":     "Tftestacc42",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_policy_name": "tf-testAcc-test-secrity",
					"ciphers":              []string{"TLS_AES_128_GCM_SHA256"},
					"tls_versions":         []string{"TLSv1.3"},
					"tags": map[string]string{
						"Created": "tfTestAcc58",
						"For":     "Tftestacc58",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": "tf-testAcc-test-secrity",
						"ciphers.#":            "1",
						"tls_versions.#":       "1",
						"tags.%":               "2",
						"tags.Created":         "tfTestAcc58",
						"tags.For":             "Tftestacc58",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

func TestAccAlicloudALBSecurityPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_security_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudALBSecurityPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbSecurityPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salbsecuritypolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudALBSecurityPolicyBasicDependence1)
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
					"security_policy_name": "tf-testAcc-test-secrity",
					"tls_versions":         []string{"TLSv1.0"},
					"ciphers":              []string{"ECDHE-ECDSA-AES128-SHA", "AES256-SHA"},
					"dry_run":              "false",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_policy_name": "tf-testAcc-test-secrity",
						"ciphers.#":            "2",
						"tls_versions.#":       "1",
						"dry_run":              "false",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudALBSecurityPolicyMap0 = map[string]string{
	"dry_run":              NOSET,
	"ciphers.#":            CHECKSET,
	"status":               CHECKSET,
	"tls_versions.#":       CHECKSET,
	"resource_group_id":    CHECKSET,
	"tags.%":               NOSET,
	"security_policy_name": "tf-testAcc-test-secrity",
}

func AlicloudALBSecurityPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func AlicloudALBSecurityPolicyBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
