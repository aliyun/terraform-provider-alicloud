package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"fmt"
	"strings"
	"testing"
	"time"

	"log"
)

func init() {
	resource.AddTestSweepers("alicloud_kms_secret", &resource.Sweeper{
		Name: "alicloud_kms_secret",
		F:    testSweepKmsSecret,
	})
}

func testSweepKmsSecret(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "ListSecrets"

	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	swept := false

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, SweepDefaultErrorMsg, "alicloud_kms_secret", action)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.SecretList.Secret", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SecretList.Secret", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["SecretName"]; !ok {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["SecretName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Kms Secret: %s", item["SecretName"].(string))
				continue
			}
			swept = true
			action = "DeleteSecret"
			request := map[string]interface{}{
				"SecretName":                 item["SecretName"],
				"ForceDeleteWithoutRecovery": true,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Kms Secret (%s): %s", item["SecretName"].(string), err)
			}
			log.Printf("[INFO] Delete Kms Secret success: %s ", item["SecretName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if swept {
		time.Sleep(3 * time.Second)
	}
	return nil
}

func TestAccAlicloudKMSSecret_Basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_kms_secret.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testaccKmsSecret_%d", rand)
	ra := resourceAttrInit(resourceId, map[string]string{
		"arn":              CHECKSET,
		"description":      "",
		"secret_data_type": "text",
		"version_stages.#": "1",
	})

	serviceFunc := func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeKmsSecret")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsSecretConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.KmsSkippedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_data":                   name,
					"secret_data_type":              "text",
					"secret_name":                   name,
					"version_id":                    "00001",
					"force_delete_without_recovery": "true",
					//"recovery_window_in_days": "7",
					"tags": map[string]string{
						"Created": "TF",
						"usage":   "acceptanceTest",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data":  name,
						"secret_name":  name,
						"version_id":   "00001",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.usage":   "acceptanceTest",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
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
					"tags": map[string]string{
						"Created": "TF",
						"Name":    name,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.usage":   REMOVEKEY,
						"tags.Created": "TF",
						"tags.Name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_data": name + "update",
					"version_id":  "00002",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data": name + "update",
						"version_id":  "00002",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    name + "update",
					"secret_data":    name,
					"version_id":     "00003",
					"version_stages": []string{"ACSCurrent", "UStage1"},
					"tags": map[string]string{
						"Description": name,
						"usage":       "acceptanceTest",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      name + "update",
						"secret_data":      name,
						"version_id":       "00003",
						"version_stages.#": "2",
						"tags.%":           "2",
						"tags.Description": name,
						"tags.usage":       "acceptanceTest",
						"tags.Created":     REMOVEKEY,
						"tags.Name":        REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudKMSSecret_WithKey(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alicloud_kms_secret.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testaccKmsSecretWithKey_%d", rand)
	ra := resourceAttrInit(resourceId, map[string]string{
		"arn":               CHECKSET,
		"description":       "",
		"encryption_key_id": CHECKSET,
		"version_stages.#":  "1",
	})

	serviceFunc := func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeKmsSecret")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsSecretWithKeyConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.KmsSkippedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_data":                   name,
					"secret_name":                   name,
					"version_id":                    "00001",
					"force_delete_without_recovery": "true",
					"encryption_key_id":             "${alicloud_kms_key.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data": name,
						"secret_name": name,
						"version_id":  "00001",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
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
					"secret_data": name + "update",
					"version_id":  "00002",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data": name + "update",
						"version_id":  "00002",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    name + "update",
					"secret_data":    name,
					"version_id":     "00003",
					"version_stages": []string{"ACSCurrent", "UStage1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      name + "update",
						"secret_data":      name,
						"version_id":       "00003",
						"version_stages.#": "2",
					}),
				),
			},
		},
	})
}

func resourceKmsSecretConfigDependence(name string) string {
	return ""
}

func resourceKmsSecretWithKeyConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}
		resource "alicloud_kms_key" "default" {
			description = var.name
			pending_window_in_days = 7
		}
`, name)
}
