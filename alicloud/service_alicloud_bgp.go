package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type BgpService struct {
	client *connectivity.AliyunClient
}

func (s *BgpService) DescribeBgpNetwork(id string) (v vpc.BgpNetwork, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return v, WrapError(err)
	}
	routerId := parts[0]
	dstCidr := parts[1]

	request := vpc.CreateDescribeBgpNetworksRequest()
	request.RegionId = s.client.RegionId
	request.RouterId = routerId

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeBgpNetworks(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{BgpNotFound, BgpInternalError}) {
			return v, WrapErrorf(Error(GetNotFoundMessage("BgpNetwork", id)), NotFoundMsg, ProviderERROR)
		}
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.DescribeBgpNetworksResponse)

	for _, resp := range response.BgpNetworks.BgpNetwork {
		if routerId == resp.RouterId && dstCidr == resp.DstCidrBlock {
			return resp, nil
		}
	}
	return v, WrapErrorf(Error(GetNotFoundMessage("BgpNetwork", id)), NotFoundMsg, ProviderERROR)
}

func (s *BgpService) DescribeBgpGroup(id string) (v vpc.BgpGroup, err error) {
	request := vpc.CreateDescribeBgpGroupsRequest()
	request.BgpGroupId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeBgpGroups(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{BgpNotFound, BgpInternalError}) {
			return v, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*vpc.DescribeBgpGroupsResponse)

	for _, resp := range response.BgpGroups.BgpGroup {
		if id == resp.BgpGroupId {
			return resp, nil
		}
	}
	return v, WrapErrorf(Error(GetNotFoundMessage("BgpGroup", id)), NotFoundMsg, ProviderERROR)
}

func (s *BgpService) DescribeBgpPeer(id string) (v vpc.BgpPeer, err error) {
	request := vpc.CreateDescribeBgpPeersRequest()
	request.BgpPeerId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeBgpPeers(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{BgpNotFound, BgpInternalError}) {
			return v, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*vpc.DescribeBgpPeersResponse)

	for _, resp := range response.BgpPeers.BgpPeer {
		if id == resp.BgpPeerId {
			return resp, nil
		}
	}
	return v, WrapErrorf(Error(GetNotFoundMessage("BgpPeer", id)), NotFoundMsg, ProviderERROR)
}

func (s *BgpService) WaitForBgpNetwork(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeBgpNetwork(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if string(status) == object.Status {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *BgpService) WaitForBgpGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeBgpGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if string(status) == object.Status {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *BgpService) WaitForBgpPeer(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeBgpPeer(id)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(5 * time.Second)
	}
}
