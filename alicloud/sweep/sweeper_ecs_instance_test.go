package sweep_test

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfconvert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_instance", &resource.Sweeper{
		Name: "alicloud_instance",
		F:    testSweepInstances,
		// When implemented, these should be removed firstly
		// Now, the resource alicloud_havip_attachment has been published.
		//Dependencies: []string{
		//	"alicloud_havip_attachment",
		//},
	})
}

func testSweepInstances(region string) error {
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-example",
		"tf_example",
		"terraform-example",
	}

	var insts []ecs.Instance
	req := ecs.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(tfhelper.PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Instances: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 1 {
			break
		}
		insts = append(insts, resp.Instances.Instance...)

		if len(resp.Instances.Instance) < tfhelper.PageSizeLarge {
			break
		}

		page, err := tfhelper.GetNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	vpcService := newVpcSweepService(client)
	for _, v := range insts {
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		if !tfsweep.SweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
			if skip {
				if need, err := vpcService.needSweepVpc(v.VpcAttributes.VpcId, v.VpcAttributes.VSwitchId); err == nil {
					skip = !need
				}

			}
			if skip {
				log.Printf("[INFO] Skipping Instance: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Deleting Instance: %s (%s)", name, id)
		if v.DeletionProtection {
			request := ecs.CreateModifyInstanceAttributeRequest()
			request.InstanceId = id
			request.DeletionProtection = requests.NewBoolean(false)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceAttribute(request)
			})
			if err != nil {
				log.Printf("[ERROR] %#v", tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, request.GetActionName(), tferrs.AlibabaCloudSdkGoERROR))
				continue
			}
		}
		if v.InstanceChargeType == string(PrePaid) {
			request := ecs.CreateModifyInstanceChargeTypeRequest()
			request.InstanceIds = tfconvert.ConvertListToJsonString(append(make([]interface{}, 0, 1), id))
			request.InstanceChargeType = string(PostPaid)
			request.IncludeDataDisks = requests.NewBoolean(true)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceChargeType(request)
			})
			if err != nil {
				log.Printf("[ERROR] %#v", tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, request.GetActionName(), tferrs.AlibabaCloudSdkGoERROR))
				continue
			}
			time.Sleep(3 * time.Second)
		}

		req := ecs.CreateDeleteInstanceRequest()
		req.InstanceId = id
		req.Force = requests.NewBoolean(true)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 20 seconds to eusure these instances have been deleted.
		time.Sleep(20 * time.Second)
	}
	return nil
}
