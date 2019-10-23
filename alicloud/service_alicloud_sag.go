package alicloud

import (
	"strings"
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

func (s *SagService) DescribeCloudConnectNetworkAttachment(id string) (c smartag.GrantRule, err error) {
	parts, _ := ParseResourceId(id, 2)
	ccn_id := parts[0]
	request := smartag.CreateDescribeGrantRulesRequest()
	request.RegionId = s.client.RegionId
	request.AssociatedCcnId = ccn_id

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

func (s *SagService) DescribeSagGrantRules(id string) (c smartag.GrantRule, err error) {
	request := smartag.CreateDescribeGrantSagRulesRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = id

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
			return sagClient.DescribeGrantSagRules(request)
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
		if IsExceptedError(err, "GrantSagRuleNotExist") {
			return c, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return c, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartag.DescribeGrantSagRulesResponse)
	if len(response.GrantRules.GrantRule) <= 0 {
		return c, WrapErrorf(Error(GetNotFoundMessage("Grant sag rule", id)), NotFoundMsg, ProviderERROR)
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
	parts_ := strings.Split(id, ":")
	acl_id := ""
	if len(parts_) != 2 {
		acl_id = id
	} else {
		parts, _ := ParseResourceId(id, 2)
		acl_id = parts[0]
	}
	request := smartag.CreateDescribeACLAttributeRequest()
	request.RegionId = s.client.RegionId
	request.AclId = acl_id

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
	parts, _ := ParseResourceId(id, 3)
	opt_id := parts[0]

	request := smartag.CreateDescribeNetworkOptimizationSettingsRequest()
	request.RegionId = s.client.RegionId
	request.NetworkOptId = opt_id

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
	parts, _ := ParseResourceId(id, 2)
	sag_id := parts[0]

	request := smartag.CreateDescribeSmartAccessGatewayClientUsersRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = sag_id

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
	parts, _ := ParseResourceId(id, 2)
	sag_id := parts[0]

	request := smartag.CreateDescribeSnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = sag_id

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
	parts, _ := ParseResourceId(id, 2)
	sag_id := parts[0]

	request := smartag.CreateDescribeDnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SagId = sag_id

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
	parts_ := strings.Split(id, ":")
	qos_id := ""
	qospy_id := ""
	if len(parts_) != 2 {
		qos_id = id
	} else {
		parts, _ := ParseResourceId(id, 2)
		qos_id = parts[0]
		qospy_id = parts[1]
	}
	request := smartag.CreateDescribeQosPoliciesRequest()
	request.RegionId = s.client.RegionId
	request.QosId = qos_id
	request.QosPolicyId = qospy_id

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
