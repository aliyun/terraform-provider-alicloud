package ecs

// Local compat copies of other-product Service structs and the specific methods
// ECS calls. ECS must not import other product packages (circular); these copies
// preserve behavior using the shared leaf aliases (errs/helper/logs/...) and SDK.

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Role-policy structs copied from the root ram service; ParseRolePolicyDocument
// is used by JudgeRolePolicyPrincipal below.
type Effect string

const (
	Allow Effect = "Allow"
	Deny  Effect = "Deny"
)

type Principal struct {
	Service []string
	RAM     []string
}

type RolePolicyStatement struct {
	Effect    Effect
	Action    string
	Principal Principal
}

type RolePolicy struct {
	Statement []RolePolicyStatement
	Version   string
}

type VpcService struct {
	client *connectivity.AliyunClient
}

func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return WrapErrorf(NotFoundErr("vswitch", id), NotFoundMsg, ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

func (s *VpcService) DescribeVSwitchWithTeadsl(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}

	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	if v, ok := response["VSwitchId"].(string); ok && v != id {
		return nil, WrapErrorf(NotFoundErr("vswitch", id), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

type RamService struct {
	client *connectivity.AliyunClient
}

// NewRamService mirrors NewService/NewServiceV2/NewEbsService: the migrated
// ram_role_attachment test (external package ecs_test) cannot construct
// RamService via a struct literal because the client field is unexported, so
// it uses this constructor, exactly as the same-package root layout used
// RamService{client}.
func NewRamService(client *connectivity.AliyunClient) *RamService {
	return &RamService{client: client}
}

func (s *RamService) DescribeRamRoleAttachment(id string) (*ecs.DescribeInstanceRamRoleResponse, error) {
	response := &ecs.DescribeInstanceRamRoleResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return response, WrapError(err)
	}
	request := ecs.CreateDescribeInstanceRamRoleRequest()
	request.RegionId = s.client.RegionId
	request.InstanceIds = parts[1]
	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
				return retry.RetryableError(WrapError(err))
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRamRole.NotFound"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response = raw.(*ecs.DescribeInstanceRamRoleResponse)
	instRoleSets := response.InstanceRamRoleSets.InstanceRamRoleSet
	if len(instRoleSets) > 0 {
		var instIds []string
		for _, item := range instRoleSets {
			if item.RamRoleName == parts[0] {
				instIds = append(instIds, item.InstanceId)
			}
		}
		ids := strings.Split(strings.TrimRight(strings.TrimLeft(strings.Replace(strings.Split(id, ":")[1], "\"", "", -1), "["), "]"), ",")
		sort.Strings(instIds)
		sort.Strings(ids)
		if reflect.DeepEqual(instIds, ids) {
			return response, nil
		}
	}
	return response, WrapErrorf(err, NotFoundMsg, ProviderERROR)
}

func (s *RamService) WaitForRamRoleAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRoleAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(object.TotalCount), status, ProviderERROR)
		}
	}
}

func (s *RamService) JudgeRolePolicyPrincipal(roleName string) error {
	request := ram.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = roleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, roleName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ram.GetRoleResponse)
	policy, err := s.ParseRolePolicyDocument(resp.Role.AssumeRolePolicyDocument)
	if err != nil {
		return WrapError(err)
	}
	for _, v := range policy.Statement {
		for _, val := range v.Principal.Service {
			if strings.Trim(val, " ") == "ecs.aliyuncs.com" {
				return nil
			}
		}
	}
	return WrapError(fmt.Errorf("Role policy services must contains 'ecs.aliyuncs.com', Now is \n%v.", resp.Role.AssumeRolePolicyDocument))
}

func (s *RamService) ParseRolePolicyDocument(policyDocument string) (RolePolicy, error) {
	var policy RolePolicy
	err := json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		return RolePolicy{}, WrapError(err)
	}
	return policy, nil
}

type OssService struct {
	client *connectivity.AliyunClient
}

func (s *OssService) WaitForOssBucketObject(bucket *oss.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		exist, err := bucket.IsObjectExist(id)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, "IsObjectExist", AliyunOssGoSdk)
		}
		addDebug("IsObjectExist", exist)

		if !exist {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatBool(exist), status, ProviderERROR)
		}
	}
}

type KmsService struct {
	client *connectivity.AliyunClient
}

func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, WrapError(err)
	}

	var response map[string]interface{}
	client := s.client
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"CiphertextBlob":    ciphertextBlob,
		"EncryptionContext": string(context[:]),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return plaintext, WrapErrorf(err, DefaultErrorMsg, context, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, WrapErrorf(err, FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

// BssOpenApiService is a local compat copy of the root bss-open-api service;
// the ecs instance_types data source calls GetInstanceTypePrice for price lookups.
type BssOpenApiService struct {
	client *connectivity.AliyunClient
}

const ModulesSizeLimit = 50

func (s *BssOpenApiService) GetInstanceTypePrice(productCode, productType, paymentType string, modules []map[string]interface{}) (object []float64, err error) {
	client := s.client
	var response map[string]interface{}
	var endpoint string
	var priceList []float64
	action := "GetSubscriptionPrice"
	request := map[string]interface{}{
		"OrderType":        "NewOrder",
		"SubscriptionType": "Subscription",
		"Region":           client.RegionId,
		"ProductCode":      productCode,
		"ProductType":      productType,
	}

	if paymentType == "Subscription" {
		request["ServicePeriodQuantity"] = 1
		request["ServicePeriodUnit"] = "Month"
		request["Quantity"] = 1
	}

	if paymentType == "PayAsYouGo" {
		action = "GetPayAsYouGoPrice"
		request["SubscriptionType"] = "PayAsYouGo"
	}

	moduleLength := len(modules)

	for {
		if len(modules) < ModulesSizeLimit {
			moduleLength = len(modules)

			request["ModuleList"] = modules
		} else {
			tmp := modules[:ModulesSizeLimit]
			modules = modules[ModulesSizeLimit:]
			moduleLength = len(tmp)

			request["ModuleList"] = &tmp
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
					endpoint = connectivity.BssOpenAPIEndpointInternational
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return object, WrapError(err)
		}

		v, err := jsonpath.Get("$.Data.ModuleDetails.ModuleDetail", response)
		if err != nil {
			return object, WrapError(err)
		}

		for _, originalCosts := range v.([]interface{}) {
			originalCostArg := originalCosts.(map[string]interface{})

			if originalCost, ok := originalCostArg["OriginalCost"]; ok {
				priceListItem, err := formatFloat64(originalCost)
				if err != nil {
					return nil, err
				}
				priceList = append(priceList, priceListItem)
			}
		}

		if moduleLength < ModulesSizeLimit {
			break
		}
	}

	return priceList, nil
}

// EbsService is a local compat copy of the root ebs service; the migrated ECS
// instance test constructs EbsService{...} and calls DescribeEbsDedicatedBlockStorageCluster
// to verify a dedicated block storage cluster. Cross-product, behavior-preserving
// copy per stock-method-disposition; uses the shared leaf aliases and SDK.
type EbsService struct {
	client *connectivity.AliyunClient
}

// NewEbsService mirrors NewService/NewServiceV2: the migrated ECS instance
// test (external package ecs_test) cannot construct EbsService via a struct
// literal because the client field is unexported, so it uses this constructor,
// exactly as the same-package root layout used EbsService{client}.
func NewEbsService(client *connectivity.AliyunClient) *EbsService {
	return &EbsService{client: client}
}

func (s *EbsService) DescribeEbsDedicatedBlockStorageCluster(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"DedicatedBlockStorageClusterId": []string{id},
		"RegionId":                       s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDedicatedBlockStorageClusters"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DedicatedBlockStorageClusters", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DedicatedBlockStorageClusters", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("EBS", id), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}
