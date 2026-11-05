package cen

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// DescribeCenInstanceGrant and WaitForCenInstanceGrant were historically
// defined on the root-package VpcService, but they describe CEN grant rules
// (a CEN concern) and are consumed only by the cen_instance_grant resource.
// They are re-homed onto *CbnService here so the migrated cen package does
// not import the root alicloud package. Behavior is unchanged.

func (s *CbnService) DescribeCenInstanceGrant(id string) (rule vpc.CbnGrantRule, err error) {
	request := vpc.CreateDescribeGrantRulesToCenRequest()
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return rule, errs.WrapError(err)
	}
	cenId := parts[0]
	instanceId := parts[1]
	instanceType, err := GetCenChildInstanceType(instanceId)
	if err != nil {
		return rule, errs.WrapError(err)
	}

	request.RegionId = s.client.RegionId
	request.InstanceId = instanceId
	request.InstanceType = instanceType

	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeGrantRulesToCen(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeGrantRulesToCenResponse)
		ruleList := response.CenGrantRules.CbnGrantRule
		if len(ruleList) <= 0 {
			return errs.WrapErrorf(errs.NotFoundErr("GrantRules", id), errs.NotFoundMsg, errs.ProviderERROR)
		}

		for ruleNum := 0; ruleNum <= len(response.CenGrantRules.CbnGrantRule)-1; ruleNum++ {
			if ruleList[ruleNum].CenInstanceId == cenId {
				rule = ruleList[ruleNum]
				return nil
			}
		}

		return errs.WrapErrorf(errs.NotFoundErr("GrantRules", id), errs.NotFoundMsg, errs.ProviderERROR)
	})
	return
}

func (s *CbnService) WaitForCenInstanceGrant(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	instanceId := parts[1]
	ownerId := parts[2]
	for {
		object, err := s.DescribeCenInstanceGrant(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.CenInstanceId == instanceId && fmt.Sprint(object.CenOwnerId) == ownerId && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.CenInstanceId, instanceId, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}
