package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type NetworkAclService struct {
	client *connectivity.AliyunClient
}

func (s *NetworkAclService) DescribeNetworkAcl(id string) (networkAcl vpc.NetworkAcl, err error) {

	request := vpc.CreateDescribeNetworkAclsRequest()
	request.NetworkAclId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeNetworkAcls(request)
	})
	if err != nil {
		if IsExceptedError(err, NetworkAclNotFound) {
			return networkAcl, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return networkAcl, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*vpc.DescribeNetworkAclsResponse)
	addDebug(request.GetActionName(), raw)
	if len(response.NetworkAcls.NetworkAcl) <= 0 || response.NetworkAcls.NetworkAcl[0].NetworkAclId != id {
		return networkAcl, WrapErrorf(Error(GetNotFoundMessage("Network Acl", id)), NotFoundMsg, ProviderERROR)
	}
	return response.NetworkAcls.NetworkAcl[0], nil
}

func (s *NetworkAclService) DescribeNetworkAclAttachment(id string, resource []vpc.Resource) (err error) {

	invoker := NewInvoker()
	return invoker.Run(func() error {
		object, err := s.DescribeNetworkAcl(id)
		if err != nil {
			return WrapError(err)
		}
		if len(object.Resources.Resource) < 1 {
			return WrapErrorf(Error(GetNotFoundMessage("Network Acl Attachment", id)), NotFoundMsg, ProviderERROR)
		}
		success := true
		for _, source := range object.Resources.Resource {
			success = false
			for _, res := range resource {
				if source.ResourceId == res.ResourceId {
					success = true
				}
			}
			if success == false {
				return WrapErrorf(Error(GetNotFoundMessage("Network Acl Attachment", id)), NotFoundMsg, ProviderERROR)
			}
		}
		return nil
	})
}

func (s *NetworkAclService) WaitForNetworkAcl(networkAclId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeNetworkAcl(networkAclId)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, networkAclId, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *NetworkAclService) WaitForNetworkAclAttachment(id string, resource []vpc.Resource, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		err := s.DescribeNetworkAclAttachment(id, resource)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		object, err := s.DescribeNetworkAcl(id)
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}
