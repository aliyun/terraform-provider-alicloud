package alikafka

import (
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alikafkasdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_alikafka_instance", &resource.Sweeper{
		Name: "alicloud_alikafka_instance",
		F:    testSweepAlikafkaInstance,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_alikafka_consumer_group",
			"alicloud_alikafka_sasl_acl",
			"alicloud_alikafka_topic",
			"alicloud_alikafka_sasl_user",
		},
	})
}

func testSweepAlikafkaInstance(region string) error {
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
	instanceListReq.RegionId = region

	raw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafkasdk.GetInstanceListResponse)
	service := VpcService{client}
	for _, v := range instanceListResp.InstanceList.InstanceVO {

		name := v.Name
		skip := true
		for _, prefix := range prefixes {

			// ServiceStatus equals 5 means the instance is in running status.
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a ES description is not set successfully, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping alikafka instance: %s ", name)
			continue
		}
		if v.ServiceStatus != 10 {
			log.Printf("[INFO] release alikafka instance: %s ", name)

			request := alikafkasdk.CreateReleaseInstanceRequest()
			request.InstanceId = v.InstanceId
			request.ForceDeleteInstance = requests.NewBoolean(true)
			_, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
				return alikafkaClient.ReleaseInstance(request)
			})

			if err != nil {
				log.Printf("[ERROR] Failed to release alikafka instance (%s): %s", name, err)
			}
		}

		log.Printf("[INFO] Delete alikafka instance: %s ", name)
		request2 := alikafkasdk.CreateDeleteInstanceRequest()
		request2.InstanceId = v.InstanceId
		_, err = client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.DeleteInstance(request2)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete alikafka instance (%s): %s", name, err)
		}
	}

	return nil
}
