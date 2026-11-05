package vpc_test

// VPC sweep helpers that lived as *VpcService methods on the root
// alicloud package's service_alicloud_common_test.go at baseline
// (needSweepVpc / sweepVSwitch / sweepNatGateway). They did not migrate with
// the resource files, and package vpc_test cannot define methods on
// *vpc.VpcService (a type from another package) nor read its unexported client
// field, so they are re-hosted here as free functions taking the
// *connectivity.AliyunClient directly. needSweepVpc still uses the exported
// DescribeVpc / DescribeVswitch methods via a NewVpcService(client) instance.
// Bodies are transcribed verbatim from the baseline originals (only the
// receiver access s.client -> client changed).

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func needSweepVpc(s *VpcService, vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := s.DescribeVswitch(vswitchId)
		if err != nil && !NotFoundError(err) {
			return false, WrapError(err)
		}
		name := strings.ToLower(object["VSwitchName"].(string))
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the vswitch (%v (%v)).", object["VSwitchId"], object["VSwitchName"])
			return true, nil
		}
		vpcId = fmt.Sprint(object["VpcId"])
	}
	if vpcId != "" {
		object, err := s.DescribeVpc(vpcId)
		if err != nil {
			if NotFoundError(err) {
				return false, nil
			}
			return false, WrapError(err)
		}
		name := strings.ToLower(object["VpcName"].(string))
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the VPC (%v (%v)).", object["VpcId"], object["VpcName"])
			return true, nil
		}
	}
	return false, nil
}

func sweepVSwitch(client *connectivity.AliyunClient, id string) error {
	var err error
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Vswitch %s ...", id)
	action := "DeleteVSwitch"
	request := map[string]interface{}{
		"VSwitchId": id,
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(10*time.Second, func() *retry.RetryError {
		_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return WrapError(err)
}

func sweepNatGateway(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}

	log.Printf("[INFO] Deleting Nat Gateway %s ...", id)
	request := vpc.CreateDeleteNatGatewayRequest()
	request.NatGatewayId = id
	request.Force = requests.NewBoolean(true)
	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteNatGateway(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return WrapError(err)
}

// sweepSecurityGroup is transcribed verbatim from the baseline root
// *EcsService.sweepSecurityGroup (service_alicloud_common_test.go). It did not
// migrate with the resource files, and package vpc_test cannot define methods
// on *EcsService nor read its unexported client field, so it is re-hosted here
// as a free function taking the *connectivity.AliyunClient directly (only the
// receiver access s.client -> client changed). It keeps the original ECS SDK
// DeleteSecurityGroup call, the 1s success sleep and the WrapError return, so
// the dependent-resource sweep retains the original retry/error/wait behavior
// rather than the generic void deleteDependentResource path.
func sweepSecurityGroup(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Security Group %s ...", id)
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.SecurityGroupId = id
	_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteSecurityGroup(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return WrapError(err)
}
