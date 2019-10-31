package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type SagService struct {
	client *connectivity.AliyunClient
}

func (s *SagService) DescribeCloudConnectNetwork(id string) (c smartag.CloudConnectNetwork, err error) {
	request := smartag.CreateDescribeCloudConnectNetworksRequest()
	request.RegionId = s.client.RegionId
	request.CcnId = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.DescribeCloudConnectNetworks(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "CcnNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeCloudConnectNetworksResponse)
	if len(response.CloudConnectNetworks.CloudConnectNetwork) <= 0 || response.CloudConnectNetworks.CloudConnectNetwork[0].CcnId != id {
		return c, WrapErrorf(Error(GetNotFoundMessage("CloudConnectNetwork ", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.CloudConnectNetworks.CloudConnectNetwork[0]
	return c, nil
}

func (s *SagService) DescribeCloudConnectNetworkGrant(id string) (c smartag.GrantRule, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeGrantRulesRequest()
	request.RegionId = s.client.RegionId
	request.AssociatedCcnId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.DescribeGrantRules(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "CcnNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeGrantRulesResponse)
	if len(response.GrantRules.GrantRule) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("GrantRule", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.GrantRules.GrantRule[0]
	return c, nil
}

func (s *SagService) DescribeSagAcl(id string) (c smartag.Acl, err error) {
	request := smartag.CreateDescribeACLsRequest()
	request.RegionId = s.client.RegionId
	request.AclIds = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeACLs(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagAclNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeACLsResponse)
	if len(response.Acls.Acl) <= 0 || response.Acls.Acl[0].AclId != id {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Acl", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.Acls.Acl[0]
	return c, nil
}

func (s *SagService) DescribeSagAclRule(id string) (c smartag.Acr, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeACLAttributeRequest()
	request.RegionId = s.client.RegionId
	request.AclId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeACLAttribute(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagAclRuleNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeACLAttributeResponse)
	if len(response.Acrs.Acr) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Acl Rule", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.Acrs.Acr[0]
	return c, nil
}

func (s *SagService) DescribeSagNetworkopt(id string) (c smartag.NetworkOptimization, err error) {
	request := smartag.CreateDescribeNetworkOptimizationsRequest()
	request.RegionId = s.client.RegionId
	request.NetworkOptId = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeNetworkOptimizations(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagNetworkoptNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeNetworkOptimizationsResponse)
	if len(response.NetworkOptimizations.NetworkOptimization) <= 0 || response.NetworkOptimizations.NetworkOptimization[0].InstanceId != id {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Networkopt", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.NetworkOptimizations.NetworkOptimization[0]
	return c, nil
}

func (s *SagService) DescribeSagNetworkoptSetting(id string) (c smartag.Setting, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeNetworkOptimizationSettingsRequest()
	request.RegionId = s.client.RegionId
	request.NetworkOptId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeNetworkOptimizationSettings(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagNetworkoptSettingNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeNetworkOptimizationSettingsResponse)
	if len(response.Settings.Setting) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Networkopt Setting", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.Settings.Setting[0]
	return c, nil
}

func (s *SagService) DescribeNetworkoptSags(id string) (c smartag.SmartAccessGateway, err error) {
	request := smartag.CreateDescribeNetworkOptimizationSagsRequest()
	request.RegionId = s.client.RegionId
	request.NetworkOptId = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeNetworkOptimizationSags(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagNetworkoptSagsNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeNetworkOptimizationSagsResponse)
	if len(response.SmartAccessGateways.SmartAccessGateway) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Networkopt Sags", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.SmartAccessGateways.SmartAccessGateway[0]
	return c, nil
}

func (s *SagService) DescribeSagClientUser(id string) (c smartag.User, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeSmartAccessGatewayClientUsersRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeSmartAccessGatewayClientUsers(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagClientUserNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeSmartAccessGatewayClientUsersResponse)
	if len(response.Users.User) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Client User", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.Users.User[0]
	return c, nil
}

func (s *SagService) DescribeSagSnatEntry(id string) (c smartag.SnatEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeSnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeSnatEntries(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SnatEntryiesNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeSnatEntriesResponse)
	if len(response.SnatEntries.SnatEntry) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag snat entry", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.SnatEntries.SnatEntry[0]
	return c, nil
}

func (s *SagService) DescribeSagDnatEntry(id string) (c smartag.DnatEntry, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeDnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SagId = parts[0]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeDnatEntries(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "DnatEntryiesNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeDnatEntriesResponse)
	if len(response.DnatEntries.DnatEntry) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag dnat entry", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.DnatEntries.DnatEntry[0]
	return c, nil
}

func (s *SagService) DescribeSagQos(id string) (c smartag.Qos, err error) {
	request := smartag.CreateDescribeQosesRequest()
	request.RegionId = s.client.RegionId
	request.QosIds = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeQoses(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagQosNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeQosesResponse)
	if len(response.Qoses.Qos) <= 0 || response.Qoses.Qos[0].QosId != id {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Qos", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.Qoses.Qos[0]
	return c, nil
}

func (s *SagService) DescribeSagQosPolicy(id string) (c smartag.QosPolicy, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return c, WrapError(err)
	}

	request := smartag.CreateDescribeQosPoliciesRequest()
	request.RegionId = s.client.RegionId
	request.QosId = parts[0]
	request.QosPolicyId = parts[1]

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeQosPolicies(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling, "Throttling.User"}) {
				time.Sleep(DefaultIntervalShort * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, "SagQosPolicyNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeQosPoliciesResponse)
	if len(response.QosPolicies.QosPolicy) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Sag Qos Policy", id)), NotFoundMsg, ProviderERROR)
	}
	c = response.QosPolicies.QosPolicy[0]
	return c, nil
}

func (s *SagService) WaitForSagAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSagAcl(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.AclId == id && status != Deleted {
			break
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AclId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForSagAclRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}

	for {
		object, err := s.DescribeSagAclRule(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.AcrId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AcrId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SagService) WaitForCloudConnectNetworkGrant(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	for {
		object, err := s.DescribeCloudConnectNetworkGrant(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if object.CenInstanceId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.CenInstanceId, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
