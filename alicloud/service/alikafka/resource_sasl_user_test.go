package alikafka_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	alikafkasdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_alikafka_sasl_user", &resource.Sweeper{
		Name: "alicloud_alikafka_sasl_user",
		F:    testSweepAlikafkaSaslUser,
		Dependencies: []string{
			"alicloud_alikafka_sasl_acl",
		},
	})
}

func testSweepAlikafkaSaslUser(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafkasdk.CreateGetInstanceListRequest()
	instanceListReq.RegionId = acctest.DefaultRegionToTest

	raw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafkasdk.GetInstanceListResponse)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		if v.ServiceStatus == 10 {
			log.Printf("[INFO] Skipping released alikafka instance id: %s ", v.InstanceId)
			continue
		}

		// Control the sasl user list request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		request := alikafkasdk.CreateDescribeSaslUsersRequest()
		request.InstanceId = v.InstanceId
		request.RegionId = acctest.DefaultRegionToTest

		raw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka sasl users on instance (%s): %s", v.InstanceId, err)
			continue
		}

		saslUserListResp, _ := raw.(*alikafkasdk.DescribeSaslUsersResponse)
		saslUsers := saslUserListResp.SaslUserList.SaslUserVO
		for _, saslUser := range saslUsers {
			name := saslUser.Username
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping alikafka sasl username: %s ", name)
				continue
			}
			log.Printf("[INFO] delete alikafka sasl username: %s ", name)

			// Control the sasl username delete rate
			time.Sleep(time.Duration(400) * time.Millisecond)

			deleteUserReq := alikafkasdk.CreateDeleteSaslUserRequest()
			deleteUserReq.InstanceId = v.InstanceId
			deleteUserReq.Username = saslUser.Username
			deleteUserReq.RegionId = acctest.DefaultRegionToTest

			_, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
				return alikafkaClient.DeleteSaslUser(deleteUserReq)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete alikafka sasl username (%s): %s", name, err)
			}
		}
	}

	return nil
}

func TestAccAliCloudAliKafkaSaslUser_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_user.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaSaslUserMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaSaslUser")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkasasluserbasic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaSaslUserBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_alikafka_instance.default.id}",
					"username":    name,
					"password":    "YourPassword1234!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"username":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword1234!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudAliKafkaSaslUser_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_user.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaSaslUserMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaSaslUser")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkasasluserbasic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaSaslUserBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_alikafka_instance.default.id}",
					"username":    name,
					"type":        "plain",
					"password":    "YourPassword1234!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"username":    name,
						"type":        "plain",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudAliKafkaSaslUser_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_user.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaSaslUserMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaSaslUser")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkasasluserbasic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaSaslUserBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":            "${alicloud_alikafka_instance.default.id}",
					"username":               name,
					"type":                   "scram",
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"username":    name,
						"type":        "scram",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.update.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name + "update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_encrypted_password":   acctest.CHECKSET,
						"kms_encryption_context.%": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

func TestAccAliCloudAliKafkaSaslUser_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_user.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaSaslUserMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaSaslUser")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salikafkasasluserbasic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaSaslUserBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":            "${alicloud_alikafka_instance.default.id}",
					"username":               name,
					"type":                   "scram",
					"mechanism":              "SCRAM-SHA-256",
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"username":    name,
						"type":        "scram",
						"mechanism":   "SCRAM-SHA-256",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"kms_encrypted_password", "kms_encryption_context"},
			},
		},
	})
}

var AliCloudAliKafkaSaslUserMap0 = map[string]string{
	"type": acctest.CHECKSET,
}

var AliCloudAliKafkaSaslUserMap1 = map[string]string{
	"password":  acctest.CHECKSET,
	"mechanism": acctest.CHECKSET,
}

func AliCloudAliKafkaSaslUserBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "10.4.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
	  		vswitch_name = var.name
	  		vpc_id       = alicloud_vpc.default.id
	  		cidr_block   = "10.4.0.0/24"
	  		zone_id      = data.alicloud_zones.default.zones.0.id
		}

		resource "alicloud_security_group" "default" {
	  		vpc_id = alicloud_vpc.default.id
		}

		resource "alicloud_alikafka_instance" "default" {
	  		name            = var.name
	  		partition_num   = 50
	  		disk_type       = "1"
	  		disk_size       = "500"
	  		deploy_type     = "5"
	  		io_max          = "20"
	  		spec_type       = "professional"
	  		service_version = "2.2.0"
	  		config          = "{\"enable.acl\":\"true\"}"
	  		vswitch_id      = alicloud_vswitch.default.id
	  		security_group  = alicloud_security_group.default.id
		}
`, name)
}

func AliCloudAliKafkaSaslUserBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "10.4.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
	  		vswitch_name = var.name
	  		vpc_id       = alicloud_vpc.default.id
	  		cidr_block   = "10.4.0.0/24"
	  		zone_id      = data.alicloud_zones.default.zones.0.id
		}

		resource "alicloud_security_group" "default" {
	  		vpc_id = alicloud_vpc.default.id
		}

		resource "alicloud_alikafka_instance" "default" {
	  		deploy_type     = "4"
	  		instance_type   = "alikafka_serverless"
	  		vswitch_id      = alicloud_vswitch.default.id
	  		spec_type       = "normal"
	  		service_version = "3.3.1"
	  		security_group  = alicloud_security_group.default.id
	  		config          = "{\"enable.acl\":\"true\"}"
	  		serverless_config {
	    		reserved_publish_capacity   = 60
	    		reserved_subscribe_capacity = 60
	  		}
		}

		resource "alicloud_kms_key" "default" {
	  		description            = var.name
	  		status                 = "Enabled"
	  		pending_window_in_days = 7
		}

		resource "alicloud_kms_ciphertext" "default" {
	  		key_id    = alicloud_kms_key.default.id
	  		plaintext = "YourPassword1234!"
	  		encryption_context = {
	    		"name" = var.name
	  		}
		}

		resource "alicloud_kms_ciphertext" "update" {
	  		key_id    = alicloud_kms_key.default.id
	  		plaintext = "YourPassword1234!update"
	  		encryption_context = {
	    		"name" = "${var.name}update"
	  		}
		}
`, name)
}
