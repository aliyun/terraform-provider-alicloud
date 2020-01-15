package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = defaultRegionToTest

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		if v.ServiceStatus == 10 {
			log.Printf("[INFO] Skipping alikafka instance id: %s ", v.InstanceId)
			continue
		}

		// Control the sasl user list request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		request := alikafka.CreateDescribeSaslUsersRequest()
		request.InstanceId = v.InstanceId
		request.RegionId = defaultRegionToTest

		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka sasl users on instance (%s): %s", v.InstanceId, err)
			continue
		}

		saslUserListResp, _ := raw.(*alikafka.DescribeSaslUsersResponse)
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

			deleteUserReq := alikafka.CreateDeleteSaslUserRequest()
			deleteUserReq.InstanceId = v.InstanceId
			deleteUserReq.Username = saslUser.Username
			deleteUserReq.RegionId = defaultRegionToTest

			_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.DeleteSaslUser(deleteUserReq)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete alikafka sasl username (%s): %s", name, err)
			}
		}
	}

	return nil
}

func TestAccAlicloudAlikafkaSaslUser_basic(t *testing.T) {

	var v *alikafka.SaslUserVO
	resourceId := "alicloud_alikafka_sasl_user.default"
	ra := resourceAttrInit(resourceId, alikafkaSaslUserBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslUserConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAlikafkaAclEnable(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_alikafka_instance.default.id}",
					"username":    "${var.name}",
					"password":    "password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand),
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"username": "newSaslUserName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": "newSaslUserName"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"password": "newPassword",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "newPassword"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"username": "${var.name}",
					"password": "password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand),
						"password": "password"}),
				),
			},
		},
	})

}

func TestAccAlicloudAlikafkaSaslUser_multi(t *testing.T) {

	var v *alikafka.SaslUserVO
	resourceId := "alicloud_alikafka_sasl_user.default.1"
	ra := resourceAttrInit(resourceId, alikafkaSaslUserBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslUserConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAlikafkaAclEnable(t)
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "2",
					"instance_id": "${alicloud_alikafka_instance.default.id}",
					"username":    "${var.name}-${count.index}",
					"password":    "password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": fmt.Sprintf("tf-testacc-alikafkasasluserbasic%v-1", rand),
						"password": "password",
					}),
				),
			},
		},
	})

}

func resourceAlikafkaSaslUserConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
 			default = "%v"
		}

		data "alicloud_zones" "default" {
			available_resource_creation= "VSwitch"
		}
		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}
		
		resource "alicloud_vswitch" "default" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.0.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name       = "${var.name}"
		}

		resource "alicloud_alikafka_instance" "default" {
          name = "${var.name}"
		  topic_quota = "50"
		  disk_type = "1"
		  disk_size = "500"
		  deploy_type = "5"
		  io_max = "20"
          vswitch_id = "${alicloud_vswitch.default.id}"
		}
		`, name)
}

var alikafkaSaslUserBasicMap = map[string]string{
	"username": "${var.name}",
	"password": "password",
}
