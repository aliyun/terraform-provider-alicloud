package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_slb_tls_cipher_policy",
		&resource.Sweeper{
			Name: "alicloud_slb_tls_cipher_policy",
			F:    testSweepSLBTlsCipherPolicy,
		})
}

func testSweepSLBTlsCipherPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		action := "ListTLSCipherPolicies"
		request := map[string]interface{}{
			"RegionId": region,
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &runtime)
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
			log.Printf("[ERROR] Failed to fetch Slb Tls Cipher Policy: %s", WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_tls_cipher_policies", action, AlibabaCloudSdkGoERROR))
			return nil
		}
		v, err := jsonpath.Get("$.TLSCipherPolicies", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse Slb Tls Cipher Policy: %s", WrapErrorf(err, FailedGetAttributeMsg, action, "$.TLSCipherPolicies", response))
			return nil
		}
		if len(v.([]interface{})) < 1 {
			log.Printf("[ERROR] Failed to fetch Slb Tls Cipher Policy: %s", WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_tls_cipher_policies", action, AlibabaCloudSdkGoERROR))
			return nil
		}

		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Slb Tls Cipher Policy: %s", item["Name"].(string))
				continue
			}

			action := "DeleteTLSCipherPolicy"
			request := map[string]interface{}{
				"RegionId":          region,
				"TLSCipherPolicyId": item["InstanceId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Slb Tls Cipher Policy (%s): %s", item["Name"].(string), err)
			} else {
				log.Printf("[INFO] Delete Slb Tls Cipher Policy success: %s ", item["Name"].(string))
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil

}

func TestAccAlicloudSLBTlsCipherPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_tls_cipher_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudSLBTlsCipherPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbTlsCipherPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbtlscipherpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSLBTlsCipherPolicyBasicDependence0)
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
					"tls_cipher_policy_name": "Tf-testAccCase",
					"tls_versions":           []string{"TLSv1.2"},
					"ciphers":                []string{"AES256-SHA256"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_cipher_policy_name": "Tf-testAccCase",
						"tls_versions.#":         "1",
						"ciphers.#":              "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ciphers": []string{"AES256-SHA256", "AES128-GCM-SHA256"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_versions": []string{"TLSv1.2", "TLSv1.3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_versions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_cipher_policy_name": "Tf-testAccCase1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_cipher_policy_name": "Tf-testAccCase1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tls_cipher_policy_name": "Tf-testAccCase2",
					"tls_versions":           []string{"TLSv1.1", "TLSv1.2"},
					"ciphers":                []string{"AES256-SHA256"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tls_cipher_policy_name": "Tf-testAccCase2",
						"tls_versions.#":         "2",
						"ciphers.#":              "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"access_key_id"},
			},
		},
	})
}

var AlicloudSLBTlsCipherPolicyMap0 = map[string]string{
	"tls_versions.#": CHECKSET,
	"access_key_id":  NOSET,
	"status":         CHECKSET,
	"ciphers.#":      CHECKSET,
}

func AlicloudSLBTlsCipherPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
