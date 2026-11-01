package adb

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	adbSdk "github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AdbService struct {
	client *connectivity.AliyunClient
}

// NewAdbService constructs an AdbService from the shared client. It replaces
// the root-package &AdbService{client} literal referenced by migrated tests,
// which can no longer reach the unexported struct field from package adb_test.
func NewAdbService(client *connectivity.AliyunClient) *AdbService {
	return &AdbService{client: client}
}

func (s *AdbService) DescribeAdbCluster(id string) (instance *adbSdk.DBClusterInDescribeDBClusters, err error) {
	request := adbSdk.CreateDescribeDBClustersRequest()
	request.RegionId = s.client.RegionId
	dbClusterIds := []string{}
	dbClusterIds = append(dbClusterIds, id)
	request.DBClusterIds = id
	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.DescribeDBClusters(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adbSdk.DescribeDBClustersResponse)
	if len(response.Items.DBCluster) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *AdbService) DescribeAdbClusterAttribute(id string) (instance *adbSdk.DBCluster, err error) {
	request := adbSdk.CreateDescribeDBClusterAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adbSdk.DescribeDBClusterAttributeResponse)
	if len(response.Items.DBCluster) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *AdbService) DescribeAdbAutoRenewAttribute(id string) (instance *adbSdk.AutoRenewAttribute, err error) {
	request := adbSdk.CreateDescribeAutoRenewAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterIds = id

	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.DescribeAutoRenewAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adbSdk.DescribeAutoRenewAttributeResponse)
	if len(response.Items.AutoRenewAttribute) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items.AutoRenewAttribute[0], nil
}

func (s *AdbService) WaitForAdbConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbConnection(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted && object != nil && object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.ConnectionString, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *AdbService) DescribeAdbConnection(id string) (*adbSdk.Address, error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(helper.DefaultIntervalLong) * time.Second)
	for {
		object, err := s.DescribeAdbClusterNetInfo(parts[0])

		if err != nil {
			if errs.NotFoundError(err) {
				return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return nil, errs.WrapError(err)
		}

		if object != nil {
			for _, p := range object {
				if p.NetType == "Public" {
					return &p, nil
				}
			}
		}
		time.Sleep(helper.DefaultIntervalMini * time.Second)
		if time.Now().After(deadline) {
			break
		}
	}

	return nil, errs.WrapErrorf(errs.NotFoundErr("DBConnection", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *AdbService) DescribeAdbClusterNetInfo(id string) ([]adbSdk.Address, error) {

	request := adbSdk.CreateDescribeDBClusterNetInfoRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterNetInfo(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*adbSdk.DescribeDBClusterNetInfoResponse)
	if len(response.Items.Address) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBInstanceNetInfo", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return response.Items.Address, nil
}

func (s *AdbService) WaitForAdbAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbAccount(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object != nil && object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.AccountStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) DescribeAdbAccount(id string) (ds *adbSdk.DBAccount, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := adbSdk.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := helper.NewInvoker()
	invoker.AddCatcher(dbInstanceStatusCatcher)
	var response *adbSdk.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
			return adbClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*adbSdk.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(response.AccountList.DBAccount) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBAccount", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return &response.AccountList.DBAccount[0], nil
}

// WaitForInstance waits for instance to given status
func (s *AdbService) WaitForAdbInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.DBClusterStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) setClusterTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := adbSdk.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = "cluster"
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithAdbClient(func(client *adbSdk.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := adbSdk.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = "cluster"
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithAdbClient(func(client *adbSdk.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

	}

	return nil
}

func (s *AdbService) diffTags(oldTags, newTags []adbSdk.TagResourcesTag) ([]adbSdk.TagResourcesTag, []adbSdk.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []adbSdk.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *AdbService) tagsToMap(tags []adbSdk.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *AdbService) tagsFromMap(m map[string]interface{}) []adbSdk.TagResourcesTag {
	result := make([]adbSdk.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, adbSdk.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *AdbService) ignoreTag(t adbSdk.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagValue)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *AdbService) DescribeTags(resourceId string, resourceType tags.TagResourceType) (tags []adbSdk.TagResource, err error) {
	request := adbSdk.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithAdbClient(func(client *adbSdk.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*adbSdk.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

// WaitForCluster waits for cluster to given status
func (s *AdbService) WaitForCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.DBClusterStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) DescribeDBSecurityIps(clusterId string) (ips []string, err error) {

	request := adbSdk.CreateDescribeDBClusterAccessWhiteListRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterAccessWhiteList(request)
	})
	if err != nil {
		return ips, errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*adbSdk.DescribeDBClusterAccessWhiteListResponse)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range resp.Items.IPArray {
		if ip.DBClusterIPArrayAttribute != "hidden" {
			ipstr += separator + ip.SecurityIPList
			separator = helper.CommaSeparated
		}
	}

	for _, ip := range strings.Split(ipstr, helper.CommaSeparated) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *AdbService) ModifyDBSecurityIps(clusterId, ips string) error {

	request := adbSdk.CreateModifyDBClusterAccessWhiteListRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.SecurityIps = ips

	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.ModifyDBClusterAccessWhiteList(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *AdbService) DescribeAdbBackupPolicy(id string) (policy *adbSdk.DescribeBackupPolicyResponse, err error) {

	request := adbSdk.CreateDescribeBackupPolicyRequest()
	request.DBClusterId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.DescribeBackupPolicy(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return policy, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return raw.(*adbSdk.DescribeBackupPolicyResponse), nil
}

func (s *AdbService) ModifyAdbBackupPolicy(clusterId, backupTime, backupPeriod string) error {

	request := adbSdk.CreateModifyBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime

	raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
		return adbClient.ModifyBackupPolicy(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *AdbService) AdbClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBClusterStatus == failState {
				return object, object.DBClusterStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.DBClusterStatus))
			}
		}
		return object, object.DBClusterStatus, nil
	}
}

func (s *AdbService) DescribeTask(id, taskId string) (*adbSdk.DescribeTaskInfoResponse, error) {
	request := adbSdk.CreateDescribeTaskInfoRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	request.TaskId = requests.Integer(taskId)

	var response *adbSdk.DescribeTaskInfoResponse
	wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithAdbClient(func(adbClient *adbSdk.Client) (interface{}, error) {
			return adbClient.DescribeTaskInfo(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)

			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*adbSdk.DescribeTaskInfoResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *AdbService) AdbTaskStateRefreshFunc(id, taskId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTask(id, taskId)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		return object, object.TaskInfo.Status, nil
	}
}

func (s *AdbService) DescribeAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}

	action := "DescribeAutoRenewAttribute"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBClusterIds": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, true)
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.AutoRenewAttribute", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.AutoRenewAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AnalyticDBForMySQL3.0", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("AnalyticDBForMySQL3.0", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) DescribeDBClusterAccessWhiteList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBClusterAccessWhiteList"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("AnalyticdbForMysql3.0DbCluster", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.IPArray", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.IPArray", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AnalyticDBForMySQL3.0", id), errs.NotFoundWithResponse, response)
	} else {
		ipList := ""
		for _, item := range v.([]interface{}) {
			if item.(map[string]interface{})["DBClusterIPArrayAttribute"] == "hidden" {
				continue
			}
			ipList += item.(map[string]interface{})["SecurityIPList"].(string) + ","
		}
		v.([]interface{})[0].(map[string]interface{})["SecurityIPList"] = strings.TrimSuffix(ipList, ",")
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	var response map[string]interface{}
	var err error
	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.TagIgnored(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err = retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err = retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *AdbService) DescribeAdbDbCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBClusterAttribute"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFoundError"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("AdbDbCluster", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.DBCluster", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.DBCluster", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AnalyticDBForMySQL3.0", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("AnalyticDBForMySQL3.0", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) DescribeDBClusters(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBClusters"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBClusterIds": id,
	}
	response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, true)
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Items.DBCluster", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.DBCluster", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("AnalyticDBForMySQL3.0", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("AnalyticDBForMySQL3.0", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) AdbDbClusterStateRefreshFunc(id string, stateField string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbDbCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object[stateField].(string) == failState {
				return object, object[stateField].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object[stateField].(string)))
			}
		}
		return object, object[stateField].(string), nil
	}
}

func (s *AdbService) DescribeAdbDbClusterLakeVersion(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBClusterAttribute"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("adb", "2021-12-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Adb:DbClusterLakeVersion", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Items.DBCluster", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.DBCluster", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Adb:DbClusterLakeVersion", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["DBClusterId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Adb:DbClusterLakeVersion", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AdbService) AdbDbClusterLakeVersionStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbDbClusterLakeVersion(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBClusterStatus"]) == failState {
				return object, fmt.Sprint(object["DBClusterStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["DBClusterStatus"])))
			}
		}

		return object, fmt.Sprint(object["DBClusterStatus"]), nil
	}
}

func (s *AdbService) DescribeClusterAccessWhiteList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeClusterAccessWhiteList"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("adb", "2021-12-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DBClusterLakeVersion", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Items.IPArray", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Items.IPArray", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DBClusterLakeVersion", id), errs.NotFoundWithResponse, response)
	}

	object = resp.([]interface{})[0].(map[string]interface{})

	return object, nil
}

func (s *AdbService) DescribeAdbResourceGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBResourceGroup"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": parts[0],
		"GroupName":   parts[1],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Adb:ResourceGroup", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.GroupsInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.GroupsInfo", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Adb:ResourceGroup", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["GroupName"]) == strings.ToUpper(parts[1]) {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Adb:ResourceGroup", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *AdbService) DescribeAdbDbClusterSSL(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBClusterSSL"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ADBVersionNotSupport"}) || helper.NeedRetry(err) {
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

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *AdbService) DescribeAdbDbClusterKernelVersion(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeKernelVersion"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("adb", "2019-03-15", action, nil, request, true)
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

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *AdbService) DescribeAdbDbClusterLakeVersionDBClusterSSL(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBClusterSSL"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("adb", "2021-12-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DBClusterLakeVersion", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *AdbService) AdbDbClusterLakeVersionDBClusterSSLStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbDbClusterLakeVersionDBClusterSSL(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["SSLEnabled"]) == failState {
				return object, fmt.Sprint(object["SSLEnabled"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["SSLEnabled"])))
			}
		}

		return object, fmt.Sprint(object["SSLEnabled"]), nil
	}
}
