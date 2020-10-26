package alicloud

import (
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type R_kvstoreService struct {
	client *connectivity.AliyunClient
}

func (s *R_kvstoreService) DescribeInstanceSSL(id string) (object r_kvstore.DescribeInstanceSSLResponse, err error) {
	request := r_kvstore.CreateDescribeInstanceSSLRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeInstanceSSL(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeInstanceSSLResponse)
	return *response, nil
}

func (s *R_kvstoreService) DescribeSecurityIps(id string) (object r_kvstore.SecurityIpGroup, err error) {
	request := r_kvstore.CreateDescribeSecurityIpsRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeSecurityIps(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeSecurityIpsResponse)

	if len(response.SecurityIpGroups.SecurityIpGroup) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return object, err
	}
	return response.SecurityIpGroups.SecurityIpGroup[0], nil
}

func (s *R_kvstoreService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]r_kvstore.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, r_kvstore.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	// 对系统 Tag 进行过滤
	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !ignoredTags(v, "") {
			removedTagKeys = append(removedTagKeys, v)
		}
	}
	if len(removedTagKeys) > 0 {
		request := r_kvstore.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removedTagKeys
		raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDbInstanceId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := r_kvstore.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
			return r_kvstoreClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDbInstanceId.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *R_kvstoreService) DescribeInstanceAutoRenewalAttribute(id string) (object r_kvstore.Item, err error) {
	request := r_kvstore.CreateDescribeInstanceAutoRenewalAttributeRequest()
	request.RegionId = s.client.RegionId

	request.DBInstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeInstanceAutoRenewalAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeInstanceAutoRenewalAttributeResponse)

	if len(response.Items.Item) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return object, err
	}
	return response.Items.Item[0], nil
}

func (s *R_kvstoreService) DescribeSecurityGroupConfiguration(id string) (object r_kvstore.ItemsInDescribeSecurityGroupConfiguration, err error) {
	request := r_kvstore.CreateDescribeSecurityGroupConfigurationRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeSecurityGroupConfiguration(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeSecurityGroupConfigurationResponse)

	if len(response.Items.EcsSecurityGroupRelation) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return object, err
	}
	return response.Items, nil
}

func (s *R_kvstoreService) DescribeKvstoreInstance(id string) (object r_kvstore.DBInstanceAttribute, err error) {
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeInstanceAttributeResponse)

	if len(response.Instances.DBInstanceAttribute) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreInstance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return object, err
	}
	return response.Instances.DBInstanceAttribute[0], nil
}

func (s *R_kvstoreService) KvstoreInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.InstanceStatus == failState {
				return object, object.InstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.InstanceStatus))
			}
		}
		return object, object.InstanceStatus, nil
	}
}

func (s *R_kvstoreService) DescribeKvstoreConnection(id string) (object r_kvstore.InstanceNetInfo, err error) {
	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeDBInstanceNetInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreConnection", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)

	if len(response.NetInfoItems.InstanceNetInfo) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreConnection", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}

	for _, netInfo := range response.NetInfoItems.InstanceNetInfo {
		// DBInstanceNetType is 0 means public network
		if netInfo.DBInstanceNetType == "0" {
			return netInfo, nil

		}
	}
	err = WrapErrorf(Error("The instance is not bound to the public IP"), DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	return
}

func (s *R_kvstoreService) DescribeKvstoreAccount(id string) (object r_kvstore.Account, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := r_kvstore.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.AccountName = parts[1]
	request.InstanceId = parts[0]

	raw, err := s.client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeAccounts(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KvstoreAccount", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*r_kvstore.DescribeAccountsResponse)

	if len(response.Accounts.Account) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("KvstoreAccount", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Accounts.Account[0], nil
}

func (s *R_kvstoreService) KvstoreAccountStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKvstoreAccount(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.AccountStatus == failState {
				return object, object.AccountStatus, WrapError(Error(FailedToReachTargetStatus, object.AccountStatus))
			}
		}
		return object, object.AccountStatus, nil
	}
}
