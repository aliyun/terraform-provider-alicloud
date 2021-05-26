package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
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

	prefix := "tf_testacc"

	req := kms.CreateListSecretsRequest()
	raw, err := client.WithKmsClient(func(kmsclient *kms.Client) (interface{}, error) {
		return kmsclient.ListSecrets(req)
	})
	log.Printf("[ERROR] %s got an error: %v\n.", req.GetActionName(), err)
	secrets := raw.(*kms.ListSecretsResponse)
	swept := false

	for _, v := range secrets.SecretList.Secret {

		if strings.HasPrefix(strings.ToLower(v.SecretName), prefix) {
			req := kms.CreateDeleteSecretRequest()
			req.SecretName = v.SecretName
			req.ForceDeleteWithoutRecovery = "true"
			raw, err = client.WithKmsClient(func(kmsclient *kms.Client) (interface{}, error) {
				return kmsclient.DeleteSecret(req)
			})
			swept = true
			log.Printf("[ERROR] %s got an error: %v\n.", req.GetActionName(), err)
			break
		}
	}
	if swept {
		time.Sleep(3 * time.Second)
	}
	return nil
}

func TestAccAlicloudKmsSecret_Basic(t *testing.T) {
	var v kms.DescribeSecretResponse

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

func TestAccAlicloudKmsSecret_WithKey(t *testing.T) {
	var v kms.DescribeSecretResponse

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
