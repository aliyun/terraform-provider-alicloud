package kms_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_kms_secret", &resource.Sweeper{
		Name: "alicloud_kms_secret",
		F:    testSweepKmsSecret,
	})
}

func testSweepKmsSecret(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	request := map[string]interface{}{
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "ListSecrets"

	var response map[string]interface{}
	swept := false

	for {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			return errs.WrapErrorf(err, errs.SweepDefaultErrorMsg, "alicloud_kms_secret", action)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.SecretList.Secret", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.SecretList.Secret", response)
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
		if len(result) < helper.PageSizeLarge {
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsSecretMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsSecretBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsSecretMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsSecretBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsSecretBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"secret_data":       acctest.CHECKSET,
						"version_id":        "v1",
						"secret_type":       "Rds",
						"encryption_key_id": acctest.CHECKSET,
						"dkms_instance_id":  acctest.CHECKSET,
						"extended_config":   acctest.CHECKSET,
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
					"rotation_interval":         acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_automatic_rotation": "false",
						"rotation_interval":         acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Secret\"}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsSecretBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"secret_data":               acctest.CHECKSET,
						"version_id":                "v1",
						"secret_type":               "Rds",
						"secret_data_type":          "text",
						"encryption_key_id":         acctest.CHECKSET,
						"dkms_instance_id":          acctest.CHECKSET,
						"extended_config":           acctest.CHECKSET,
						"enable_automatic_rotation": "true",
						"rotation_interval":         "605800s",
						"policy":                    acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsSecretBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_name":                   "acs/ecs/" + name,
					"secret_data":                   `{\"UserName\":\"` + "${alicloud_instance.default.instance_name}" + `\",\"Password\":\"` + "${alicloud_instance.default.password}" + `\"}`,
					"version_id":                    "v1",
					"secret_type":                   "ECS",
					"encryption_key_id":             "${alicloud_kms_key.default.id}",
					"dkms_instance_id":              "${alicloud_kms_instance.default.id}",
					"extended_config":               `{\"CommandId\":\"\",\"CustomData\":{\"tf-testAcc\":\"tf-testAcc\"},\"InstanceId\":\"` + "${alicloud_instance.default.id}" + `\",\"RegionId\":\"` + acctest.DefaultRegionToTest + `\",\"SecretSubType\":\"Password\"}`,
					"force_delete_without_recovery": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_name":       acctest.CHECKSET,
						"secret_data":       acctest.CHECKSET,
						"version_id":        "v1",
						"secret_type":       "ECS",
						"encryption_key_id": acctest.CHECKSET,
						"dkms_instance_id":  acctest.CHECKSET,
						"extended_config":   acctest.CHECKSET,
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
					"rotation_interval":         acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_automatic_rotation": "false",
						"rotation_interval":         acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Secret\"}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsSecretMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsSecret")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsSecret_%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsSecretBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"extended_config":           `{\"CommandId\":\"\",\"CustomData\":{\"tf-testAcc\":\"tf-testAcc\"},\"InstanceId\":\"` + "${alicloud_instance.default.id}" + `\",\"RegionId\":\"` + acctest.DefaultRegionToTest + `\",\"SecretSubType\":\"Password\"}`,
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
						"secret_name":               acctest.CHECKSET,
						"secret_data":               acctest.CHECKSET,
						"version_id":                "v1",
						"secret_type":               "ECS",
						"secret_data_type":          "text",
						"encryption_key_id":         acctest.CHECKSET,
						"dkms_instance_id":          acctest.CHECKSET,
						"extended_config":           acctest.CHECKSET,
						"enable_automatic_rotation": "true",
						"rotation_interval":         "605800s",
						"policy":                    acctest.CHECKSET,
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
	"secret_type":      acctest.CHECKSET,
	"version_stages.#": acctest.CHECKSET,
	"arn":              acctest.CHECKSET,
	"create_time":      acctest.CHECKSET,
}

var AliCloudKmsSecretMap1 = map[string]string{
	"secret_type":      acctest.CHECKSET,
	"policy":           acctest.CHECKSET,
	"version_stages.#": acctest.CHECKSET,
	"arn":              acctest.CHECKSET,
	"create_time":      acctest.CHECKSET,
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
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone = data.alicloud_zones.default.zones.2.id
  		image_id          = data.alicloud_images.default.images.0.id
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.2.id
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
		payment_type    = "PayAsYouGo"
  		vswitch_ids = [
    		data.alicloud_vswitches.default.ids.0
  		]
  		zone_ids = [
    		data.alicloud_zones.default.zones.2.id,
    		data.alicloud_zones.default.zones.1.id
  		]
	}

	resource "alicloud_kms_key" "default" {
  		dkms_instance_id       = alicloud_kms_instance.default.id
  		pending_window_in_days = 7
	}

	resource "alicloud_db_instance" "default" {
  		engine           = "MySQL"
  		engine_version   = "8.0"
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
  		availability_zone          = data.alicloud_zones.default.zones.2.id
  		instance_charge_type       = "PostPaid"
  		password                   = "YourPassword12345!"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

// lintignore: R001
