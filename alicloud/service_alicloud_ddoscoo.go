package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DdoscooService struct {
	client *connectivity.AliyunClient
}

func (s *DdoscooService) DescribeDdoscooInstance(id string) (v ddoscoo.Instance, err error) {
	request := ddoscoo.CreateDescribeInstancesRequest()
	request.RegionId = s.client.RegionId
	request.InstanceIds = "[\"" + id + "\"]"
	request.PageNo = "1"
	request.PageSize = "10"

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{DdoscooInstanceNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		resp, _ := raw.(*ddoscoo.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances) == 0 || resp.Instances[0].InstanceId != id {
			return WrapErrorf(Error(GetNotFoundMessage("Ddoscoo Instance", id)), NotFoundMsg, ProviderERROR)
		}

		v = resp.Instances[0]
		return nil
	})

	return v, WrapError(err)
}

func (s *DdoscooService) DescribeDdoscooInstanceSpec(id string) (v ddoscoo.InstanceSpec, err error) {
	request := ddoscoo.CreateDescribeInstanceSpecsRequest()
	request.RegionId = s.client.RegionId
	request.InstanceIds = "[\"" + id + "\"]"

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstanceSpecs(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{DdoscooInstanceNotFound, InvalidDdoscooInstance}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		resp, _ := raw.(*ddoscoo.DescribeInstanceSpecsResponse)
		if resp == nil || len(resp.InstanceSpecs) == 0 || resp.InstanceSpecs[0].InstanceId != id {
			return WrapErrorf(Error(GetNotFoundMessage("Ddoscoo Instance", id)), NotFoundMsg, ProviderERROR)
		}

		v = resp.InstanceSpecs[0]
		return nil
	})

	return v, WrapError(err)
}

func (s *DdoscooService) UpdateDdoscooInstanceName(instanceId string, name string) error {
	request := ddoscoo.CreateModifyInstanceRemarkRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = instanceId
	request.Remark = name

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.ModifyInstanceRemark(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, instanceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *DdoscooService) UpdateInstanceSpec(schemaName string, specName string, d *schema.ResourceData, meta interface{}) error {
	request := bssopenapi.CreateModifyInstanceRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = d.Id()

	request.ProductCode = "ddos"
	request.ProductType = "ddoscoo"
	request.SubscriptionType = "Subscription"

	o, n := d.GetChange(schemaName)
	oi, _ := strconv.Atoi(o.(string))
	ni, _ := strconv.Atoi(n.(string))
	if ni < oi {
		request.ModifyType = "Downgrade"
	} else {
		request.ModifyType = "Upgrade"
	}

	request.Parameter = &[]bssopenapi.ModifyInstanceParameter{
		{
			Code:  specName,
			Value: d.Get(schemaName).(string),
		},
	}

	raw, err := s.client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
		return bssopenapiClient.ModifyInstance(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*bssopenapi.ModifyInstanceResponse)
	if !resp.Success {
		return WrapError(Error(resp.Message))
	}
	return nil
}
