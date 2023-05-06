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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func SkipTestAccAlicloudKMSSecret_DKMS(t *testing.T) {
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
		},
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
					"dkms_instance_id":              os.Getenv("DKMS_INSTANCE_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data":      name,
						"secret_name":      name,
						"version_id":       "00001",
						"dkms_instance_id": CHECKSET,
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

func TestAccAlicloudKMSSecret_WithSecretTypeGeneric(t *testing.T) {
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
					"secret_type":                   "Generic",
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
						"secret_type":  "Generic",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.usage":   "acceptanceTest",
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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func TestAccAlicloudKMSSecret_WithSecretTypeRds(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsSecretWithSecretTypeRdsConfigDependence)

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
					"secret_data":                   `{\"Accounts\":[{\"AccountName\":\"` + "${alicloud_db_account.default.account_name}" + `\",\"AccountPassword\":\"` + "${alicloud_db_account.default.password}" + `\"}]}`,
					"secret_data_type":              "text",
					"secret_name":                   name,
					"version_id":                    "00001",
					"force_delete_without_recovery": "true",
					"secret_type":                   "Rds",
					"extended_config":               `{\"CustomData\":{\"test\":\"test\"},\"DBInstanceId\":\"` + "${alicloud_db_account.default.db_instance_id}" + `\",\"SecretSubType\":\"SingleUser\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"usage":   "acceptanceTest",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data":     CHECKSET,
						"secret_name":     name,
						"version_id":      "00001",
						"secret_type":     "Rds",
						"extended_config": CHECKSET,
						"tags.%":          "2",
						"tags.Created":    "TF",
						"tags.usage":      "acceptanceTest",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
			},
		},
	})
}

func TestAccAlicloudKMSSecret_WithSecretTypeECS(t *testing.T) {
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsSecretWithSecretTypeECSConfigDependence)

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
					"secret_data":                   `{\"UserName\":\"` + "${alicloud_instance.default.instance_name}" + `\",\"Password\":\"` + "${alicloud_instance.default.password}" + `\"}`,
					"secret_data_type":              "text",
					"secret_name":                   "acs/ecs/" + name,
					"version_id":                    "00001",
					"force_delete_without_recovery": "true",
					"secret_type":                   "ECS",
					"extended_config":               `{\"CommandId\":\"\",\"CustomData\":{\"test\":\"test\"},\"InstanceId\":\"` + "${alicloud_instance.default.id}" + `\",\"RegionId\":\"` + defaultRegionToTest + `\",\"SecretSubType\":\"Password\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"usage":   "acceptanceTest",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_data":     CHECKSET,
						"secret_name":     "acs/ecs/" + name,
						"version_id":      "00001",
						"secret_type":     "ECS",
						"extended_config": CHECKSET,
						"tags.%":          "2",
						"tags.Created":    "TF",
						"tags.usage":      "acceptanceTest",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete_without_recovery", "recovery_window_in_days"},
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

func resourceKmsSecretWithSecretTypeRdsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "Rds"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones[0].id
  		vswitch_name = var.name
	}

	resource "alicloud_db_instance" "default" {
  		engine           = "MySQL"
  		engine_version   = "5.6"
  		instance_type    = "rds.mysql.s1.small"
  		instance_storage = "10"
  		vswitch_id       = alicloud_vswitch.default.id
  		instance_name    = var.name
	}

	resource "alicloud_db_account" "default" {
  		db_instance_id = alicloud_db_instance.default.id
  		account_name   = "tftestnormal"
  		password       = "YourPassword12345!"
	}
`, name)
}

func resourceKmsSecretWithSecretTypeECSConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
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

func TestUnitAlicloudKMSSecret(t *testing.T) {
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
	err = resourceAlicloudKmsSecretCreate(dInit, rawClient)
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
		err := resourceAlicloudKmsSecretCreate(dInit, rawClient)
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
	err = resourceAlicloudMseClusterUpdate(dExisted, rawClient)
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
		err := resourceAlicloudMseClusterUpdate(dExisted, rawClient)
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
		err := resourceAlicloudMseClusterUpdate(dExisted, rawClient)
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
		err := resourceAlicloudKmsSecretRead(dExisted, rawClient)
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
	err = resourceAlicloudKmsSecretDelete(dExisted, rawClient)
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
		err := resourceAlicloudKmsSecretDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "Forbidden.ResourceNotFound":
			assert.Nil(t, err)
		}
	}
}
