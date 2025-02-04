package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

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
		return fmt.Errorf("error getting AliCloud client: %s", err)
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
	swept := false

	for {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
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
			_, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
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

func TestAccAliCloudKmsSecret_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_secret.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsSecretMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsSecret")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsSecretBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_name":                   name,
					"secret_data":                   name,
					"version_id":                    "v1",
					"force_delete_without_recovery": "false",
					"recovery_window_in_days":       "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_name": name,
						"secret_data": name,
						"version_id":  "v1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_data": name + "-update",
					"version_id":  "v2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data": name + "-update",
						"version_id":  "v2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_data_type": "binary",
					"version_id":       "v3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data_type": "binary",
						"version_id":       "v3",
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
					"version_stages": []string{"ACSCurrent", "ACSNext"},
					"version_id":     "v5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_stages.#": "2",
						"version_id":       "v5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Secret",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Secret",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsSecret_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_secret.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsSecretMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsSecret")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsSecretBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_name":                   name,
					"secret_data":                   name,
					"version_id":                    "v1",
					"secret_type":                   "Generic",
					"secret_data_type":              "binary",
					"description":                   name,
					"version_stages":                []string{"ACSCurrent"},
					"force_delete_without_recovery": "false",
					"recovery_window_in_days":       "7",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Secret",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_name":      name,
						"secret_data":      name,
						"version_id":       "v1",
						"secret_type":      "Generic",
						"secret_data_type": "binary",
						"description":      name,
						"version_stages.#": "1",
						"tags.%":           "2",
						"tags.Created":     "TF",
						"tags.For":         "Secret",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsSecret_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_secret.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsSecret")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsSecretBasicDependence1)
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
					"secret_name":                   name,
					"secret_data":                   `{\"Accounts\":[{\"AccountName\":\"` + "tf-testAcc" + `\",\"AccountPassword\":\"` + "YourPassword12345!" + `\"}]}`,
					"version_id":                    "v1",
					"secret_type":                   "Rds",
					"encryption_key_id":             "${alicloud_kms_key.default.id}",
					"dkms_instance_id":              "${alicloud_kms_instance.default.id}",
					"extended_config":               `{\"CustomData\":{\"tf-testAcc\":\"tf-testAcc\"},\"DBInstanceId\":\"` + "${alicloud_db_instance.default.id}" + `\",\"SecretSubType\":\"SingleUser\"}`,
					"force_delete_without_recovery": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_name":       name,
						"secret_data":       CHECKSET,
						"version_id":        "v1",
						"secret_type":       "Rds",
						"encryption_key_id": CHECKSET,
						"dkms_instance_id":  CHECKSET,
						"extended_config":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_automatic_rotation": "true",
					"rotation_interval":         "605800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_automatic_rotation": "true",
						"rotation_interval":         "605800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_automatic_rotation": "false",
					"rotation_interval":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_automatic_rotation": "false",
						"rotation_interval":         REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Secret\"}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Secret",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Secret",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsSecret_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_secret.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsSecret")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsSecretBasicDependence1)
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
					"secret_name":               name,
					"secret_data":               `{\"Accounts\":[{\"AccountName\":\"` + "tf-testAcc" + `\",\"AccountPassword\":\"` + "YourPassword12345!" + `\"}]}`,
					"version_id":                "v1",
					"secret_type":               "Rds",
					"secret_data_type":          "text",
					"encryption_key_id":         "${alicloud_kms_key.default.id}",
					"dkms_instance_id":          "${alicloud_kms_instance.default.id}",
					"extended_config":           `{\"CustomData\":{\"tf-testAcc\":\"tf-testAcc\"},\"DBInstanceId\":\"` + "${alicloud_db_instance.default.id}" + `\",\"SecretSubType\":\"SingleUser\"}`,
					"enable_automatic_rotation": "true",
					"rotation_interval":         "605800s",
					"policy":                    `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Secret\"}],\"Version\": \"1\"}`,
					"description":               name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Secret",
					},
					"force_delete_without_recovery": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_name":               name,
						"secret_data":               CHECKSET,
						"version_id":                "v1",
						"secret_type":               "Rds",
						"secret_data_type":          "text",
						"encryption_key_id":         CHECKSET,
						"dkms_instance_id":          CHECKSET,
						"extended_config":           CHECKSET,
						"enable_automatic_rotation": "true",
						"rotation_interval":         "605800s",
						"policy":                    CHECKSET,
						"description":               name,
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Secret",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsSecret_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_secret.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsSecret")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsSecretBasicDependence1)
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
					"secret_name":                   "acs/ecs/" + name,
					"secret_data":                   `{\"UserName\":\"` + "${alicloud_instance.default.instance_name}" + `\",\"Password\":\"` + "${alicloud_instance.default.password}" + `\"}`,
					"version_id":                    "v1",
					"secret_type":                   "ECS",
					"encryption_key_id":             "${alicloud_kms_key.default.id}",
					"dkms_instance_id":              "${alicloud_kms_instance.default.id}",
					"extended_config":               `{\"CommandId\":\"\",\"CustomData\":{\"tf-testAcc\":\"tf-testAcc\"},\"InstanceId\":\"` + "${alicloud_instance.default.id}" + `\",\"RegionId\":\"` + defaultRegionToTest + `\",\"SecretSubType\":\"Password\"}`,
					"force_delete_without_recovery": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_name":       CHECKSET,
						"secret_data":       CHECKSET,
						"version_id":        "v1",
						"secret_type":       "ECS",
						"encryption_key_id": CHECKSET,
						"dkms_instance_id":  CHECKSET,
						"extended_config":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_automatic_rotation": "true",
					"rotation_interval":         "605800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_automatic_rotation": "true",
						"rotation_interval":         "605800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_automatic_rotation": "false",
					"rotation_interval":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_automatic_rotation": "false",
						"rotation_interval":         REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Secret\"}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": CHECKSET,
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Secret",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Secret",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsSecret_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_secret.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsSecret")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsSecretBasicDependence1)
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
					"secret_name":               "acs/ecs/" + name,
					"secret_data":               `{\"UserName\":\"` + "${alicloud_instance.default.instance_name}" + `\",\"Password\":\"` + "${alicloud_instance.default.password}" + `\"}`,
					"version_id":                "v1",
					"secret_type":               "ECS",
					"secret_data_type":          "text",
					"encryption_key_id":         "${alicloud_kms_key.default.id}",
					"dkms_instance_id":          "${alicloud_kms_instance.default.id}",
					"extended_config":           `{\"CommandId\":\"\",\"CustomData\":{\"tf-testAcc\":\"tf-testAcc\"},\"InstanceId\":\"` + "${alicloud_instance.default.id}" + `\",\"RegionId\":\"` + defaultRegionToTest + `\",\"SecretSubType\":\"Password\"}`,
					"enable_automatic_rotation": "true",
					"rotation_interval":         "605800s",
					"policy":                    `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Secret\"}],\"Version\": \"1\"}`,
					"description":               name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Secret",
					},
					"force_delete_without_recovery": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_name":               CHECKSET,
						"secret_data":               CHECKSET,
						"version_id":                "v1",
						"secret_type":               "ECS",
						"secret_data_type":          "text",
						"encryption_key_id":         CHECKSET,
						"dkms_instance_id":          CHECKSET,
						"extended_config":           CHECKSET,
						"enable_automatic_rotation": "true",
						"rotation_interval":         "605800s",
						"policy":                    CHECKSET,
						"description":               name,
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Secret",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

var AliCloudKmsSecretMap0 = map[string]string{
	"secret_type":      CHECKSET,
	"version_stages.#": CHECKSET,
	"arn":              CHECKSET,
	"create_time":      CHECKSET,
}

var AliCloudKmsSecretMap1 = map[string]string{
	"secret_type":      CHECKSET,
	"policy":           CHECKSET,
	"version_stages.#": CHECKSET,
	"arn":              CHECKSET,
	"create_time":      CHECKSET,
}

func AliCloudKmsSecretBasicDependence0(name string) string {
	return ""
}

func AliCloudKmsSecretBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_zones" "default" {
	}

	data "alicloud_images" "default" {
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone = data.alicloud_zones.default.zones.0.id
  		image_id          = data.alicloud_images.default.images.0.id
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_kms_instance" "default" {
  		product_version = "3"
  		vpc_num         = "1"
  		key_num         = "1000"
  		secret_num      = "1000"
  		spec            = "1000"
  		vpc_id          = data.alicloud_vpcs.default.ids.0
  		vswitch_ids = [
    		data.alicloud_vswitches.default.ids.0
  		]
  		zone_ids = [
    		data.alicloud_zones.default.zones.0.id,
    		data.alicloud_zones.default.zones.1.id
  		]
	}

	resource "alicloud_kms_key" "default" {
  		dkms_instance_id       = alicloud_kms_instance.default.id
  		pending_window_in_days = 7
	}

	resource "alicloud_db_instance" "default" {
  		engine           = "MySQL"
  		engine_version   = "5.6"
  		instance_type    = "rds.mysql.s1.small"
  		instance_storage = "10"
  		vswitch_id       = data.alicloud_vswitches.default.ids.0
  		instance_name    = var.name
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		instance_name              = var.name
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_zones.default.zones.0.id
  		instance_charge_type       = "PostPaid"
  		password                   = "YourPassword12345!"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

func TestUnitAliCloudKmsSecret(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_kms_secret"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_kms_secret"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"description":               "CreateSecretValue",
		"enable_automatic_rotation": false,
		"secret_data":               "CreateSecretValue",
		"secret_data_type":          "CreateSecretValue",
		"secret_name":               "CreateSecretValue",
		"version_id":                "CreateSecretValue",
		"encryption_key_id":         "CreateSecretValue",
		"rotation_interval":         "CreateSecretValue",
		"tags": map[string]string{
			"Created": "TF",
		},
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeSecret
		"Arn":               "CreateSecretValue",
		"Description":       "CreateSecretValue",
		"EncryptionKeyId":   "CreateSecretValue",
		"PlannedDeleteTime": "CreateSecretValue",
		"Tags": map[string]interface{}{
			"Tag": []interface{}{
				map[string]interface{}{
					"Key":   "Created",
					"Value": "TF",
				},
			},
		},
		"SecretData":     "CreateSecretValue",
		"SecretDataType": "CreateSecretValue",
		"VersionId":      "CreateSecretValue",
		"VersionStages": map[string]interface{}{
			"VersionStage": []string{"CreateSecretValue"},
		},
		"SecretName": "CreateSecretValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateSecret
		"SecretName": "CreateSecretValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_kms_secret", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudKmsSecretCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeSecret Response
		"SecretName": "CreateSecretValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateSecret" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudKmsSecretCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_kms_secret"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudKmsSecretUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateSecret
	attributesDiff := map[string]interface{}{
		"description": "UpdateSecretValue",
	}
	diff, err := newInstanceDiff("alicloud_kms_secret", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_kms_secret"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeSecret
		"Description": "UpdateSecretValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateSecret" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudKmsSecretUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_kms_secret"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// PutSecretValue
	attributesDiff = map[string]interface{}{
		"secret_data":      "PutSecretValue",
		"version_id":       "PutSecretValue",
		"secret_data_type": "PutSecretValue",
		"version_stages":   []string{"PutSecretValue", "PutSecretValue"},
	}
	diff, err = newInstanceDiff("alicloud_kms_secret", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_kms_secret"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeSecret
		"SecretData":     "PutSecretValue",
		"SecretDataType": "PutSecretValue",
		"VersionId":      "PutSecretValue",
		"VersionStages": map[string]interface{}{
			"VersionStage": []string{"PutSecretValue"},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "PutSecretValue" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudKmsSecretUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_kms_secret"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeSecret" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudKmsSecretRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudKmsSecretDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "Forbidden.ResourceNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteSecret" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudKmsSecretDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "Forbidden.ResourceNotFound":
			assert.Nil(t, err)
		}
	}
}
