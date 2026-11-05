package smartag

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	smartagsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the acceptance tests) build
// values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// Status and the constants below mirror the flat-layout root package
// declarations historically used by the smartag WaitFor* helpers. No leaf
// package exports them today, so they are kept here verbatim to preserve the
// existing wait semantics during the split-layout migration.
type Status string

const (
	Active   = Status("Active")
	Inactive = Status("Inactive")
	Deleted  = Status("Deleted")
)

// DefaultTimeoutMedium mirrors the root package constant (in seconds) used as
// the default wait timeout by the smartag CRUD handlers.
const DefaultTimeoutMedium = 500

func (s *Service) DescribeCloudConnectNetwork(id string) (c smartagsdk.CloudConnectNetwork, err error) {
	request := smartagsdk.CreateDescribeCloudConnectNetworksRequest()
	request.RegionId = s.client.RegionId
	request.CcnId = id

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartagsdk.Client) (interface{}, error) {
			return ccnClient.DescribeCloudConnectNetworks(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"CcnNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeCloudConnectNetworksResponse)
	if len(response.CloudConnectNetworks.CloudConnectNetwork) <= 0 || response.CloudConnectNetworks.CloudConnectNetwork[0].CcnId != id {
		return c, errs.WrapErrorf(errs.NotFoundErr("CloudConnectNetwork ", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.CloudConnectNetworks.CloudConnectNetwork[0]
	return c, nil
}

func (s *Service) DescribeCloudConnectNetworkGrant(id string) (c smartagsdk.GrantRule, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := smartagsdk.CreateDescribeGrantRulesRequest()
	request.RegionId = s.client.RegionId
	request.AssociatedCcnId = parts[0]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartagsdk.Client) (interface{}, error) {
			return ccnClient.DescribeGrantRules(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"CcnNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeGrantRulesResponse)
	for _, value := range response.GrantRules.GrantRule {
		if value.CcnInstanceId == parts[0] && value.CenInstanceId == parts[1] {
			return value, nil
		}
	}
	return c, errs.WrapErrorf(errs.NotFoundErr("CloudConnectNetworkGrant", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *Service) DescribeCloudConnectNetworkAttachment(id string) (c smartagsdk.SmartAccessGateway, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}
	ccnId := parts[0]
	sagId := parts[1]
	request := smartagsdk.CreateDescribeSmartAccessGatewaysRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = sagId

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(ccnClient *smartagsdk.Client) (interface{}, error) {
			return ccnClient.DescribeSmartAccessGateways(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SmartAccessGatewayNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeSmartAccessGatewaysResponse)
	if len(response.SmartAccessGateways.SmartAccessGateway) <= 0 || response.SmartAccessGateways.SmartAccessGateway[0].AssociatedCcnId != ccnId || response.SmartAccessGateways.SmartAccessGateway[0].SmartAGId != sagId {
		return c, errs.WrapErrorf(errs.NotFoundErr("SmartAccessGatewayAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.SmartAccessGateways.SmartAccessGateway[0]
	return c, nil
}

func (s *Service) DescribeSagAcl(id string) (c smartagsdk.Acl, err error) {
	request := smartagsdk.CreateDescribeACLsRequest()
	request.RegionId = s.client.RegionId
	request.AclIds = id

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeACLs(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SagAclNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeACLsResponse)
	if len(response.Acls.Acl) <= 0 || response.Acls.Acl[0].AclId != id {
		return c, errs.WrapErrorf(errs.NotFoundErr("Sag Acl", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.Acls.Acl[0]
	return c, nil
}

func (s *Service) DescribeSagAclRule(id string) (c smartagsdk.Acr, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := smartagsdk.CreateDescribeACLAttributeRequest()
	request.RegionId = s.client.RegionId
	request.AclId = parts[0]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeACLAttribute(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SagAclRuleNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeACLAttributeResponse)
	for _, value := range response.Acrs.Acr {
		if value.AcrId == parts[1] {
			return value, nil
		}
	}
	return c, errs.WrapErrorf(errs.NotFoundErr("Sag Acl Rule", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *Service) DescribeSagClientUser(id string) (c smartagsdk.User, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := smartagsdk.CreateDescribeSmartAccessGatewayClientUsersRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = parts[0]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeSmartAccessGatewayClientUsers(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SagClientUserNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeSmartAccessGatewayClientUsersResponse)
	for _, value := range response.Users.User {
		if value.UserName == parts[1] {
			return value, nil
		}
	}
	return c, errs.WrapErrorf(errs.NotFoundErr("Sag Client User", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *Service) DescribeSagSnatEntry(id string) (c smartagsdk.SnatEntry, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := smartagsdk.CreateDescribeSnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SmartAGId = parts[0]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeSnatEntries(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SnatEntryiesNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeSnatEntriesResponse)
	for _, value := range response.SnatEntries.SnatEntry {
		if value.InstanceId == parts[1] {
			return value, nil
		}
	}
	return c, errs.WrapErrorf(errs.NotFoundErr("sag_snat_entry", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *Service) DescribeSagDnatEntry(id string) (c smartagsdk.DnatEntry, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := smartagsdk.CreateDescribeDnatEntriesRequest()
	request.RegionId = s.client.RegionId
	request.SagId = parts[0]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeDnatEntries(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DnatEntryiesNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeDnatEntriesResponse)
	for _, value := range response.DnatEntries.DnatEntry {
		if value.DnatEntryId == parts[1] {
			return value, nil
		}
	}
	return c, errs.WrapErrorf(errs.NotFoundErr("sag_dnat_entry", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *Service) DescribeSagQos(id string) (c smartagsdk.Qos, err error) {
	request := smartagsdk.CreateDescribeQosesRequest()
	request.RegionId = s.client.RegionId
	request.QosIds = id

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeQoses(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SagQosNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeQosesResponse)
	if len(response.Qoses.Qos) <= 0 || response.Qoses.Qos[0].QosId != id {
		return c, errs.WrapErrorf(errs.NotFoundErr("Sag Qos", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.Qoses.Qos[0]
	return c, nil
}

func (s *Service) DescribeSagQosPolicy(id string) (c smartagsdk.QosPolicy, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := smartagsdk.CreateDescribeQosPoliciesRequest()
	request.RegionId = s.client.RegionId
	request.QosId = parts[0]
	request.QosPolicyId = parts[1]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeQosPolicies(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SagQosPolicyNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeQosPoliciesResponse)
	if len(response.QosPolicies.QosPolicy) <= 0 || response.QosPolicies.QosPolicy[0].QosPolicyId != parts[1] {
		return c, errs.WrapErrorf(errs.NotFoundErr("Sag Qos Policy", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.QosPolicies.QosPolicy[0]
	return c, nil
}

func (s *Service) DescribeSagQosCar(id string) (c smartagsdk.QosCar, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := smartagsdk.CreateDescribeQosCarsRequest()
	request.RegionId = s.client.RegionId
	request.QosId = parts[0]
	request.QosCarId = parts[1]

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSagClient(func(sagClient *smartagsdk.Client) (interface{}, error) {
			return sagClient.DescribeQosCars(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling, "Throttling.User"}) {
				time.Sleep(helper.DefaultIntervalShort * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SagQosCarNotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*smartagsdk.DescribeQosCarsResponse)
	if len(response.QosCars.QosCar) <= 0 || response.QosCars.QosCar[0].QosCarId != parts[1] {
		return c, errs.WrapErrorf(errs.NotFoundErr("Sag Qos Car", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.QosCars.QosCar[0]
	return c, nil
}

func (s *Service) WaitForCloudConnectNetwork(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCloudConnectNetwork(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.CcnId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.CcnId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForCloudConnectNetworkGrant(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeCloudConnectNetworkGrant(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.CenInstanceId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.CenInstanceId, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForCloudConnectNetworkAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	ccnId := parts[0]
	for {
		object, err := s.DescribeCloudConnectNetworkAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.AssociatedCcnId == ccnId && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.AssociatedCcnId, ccnId, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSagAcl(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.AclId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.AclId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagAclRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}

	for {
		object, err := s.DescribeSagAclRule(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.AcrId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.AcrId, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagQos(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeSagQos(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.QosId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.QosId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagQosPolicy(id string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeSagQosPolicy(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.QosPolicyId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.QosPolicyId, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagQosCar(id string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSagQosCar(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.QosCarId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.QosCarId, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagSnatEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeSagSnatEntry(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.InstanceId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.InstanceId, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagDnatEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeSagDnatEntry(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.DnatEntryId == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.DnatEntryId, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSagClientUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	for {
		object, err := s.DescribeSagClientUser(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		if object.UserName == parts[1] && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.UserName, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) DescribeSmartagFlowLog(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeFlowLogs"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"FlowLogId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Smartag", "2018-03-13", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.FlowLogs.FlowLogSetType", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.FlowLogs.FlowLogSetType", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Smartag", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["FlowLogId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("Smartag", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) SmartagFlowLogStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSmartagFlowLog(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
