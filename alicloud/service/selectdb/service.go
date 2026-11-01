package selectdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) RequestProcess(request map[string]interface{}, action string, method string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	wait := helper.IncrementalWait(5*time.Second, 20*time.Second)
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		if method == "GET" {
			response, err = client.RpcGet("selectdb", "2023-05-22", action, request, nil)
		} else {
			response, err = client.RpcPost("selectdb", "2023-05-22", action, nil, request, true)
		}
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) ||
				errs.IsExpectedErrors(err, []string{"InvalidDBInstanceState.NotSupport"}) ||
				errs.IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) ||
				helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, request, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, request, "$", response)
	}
	result := v.(map[string]interface{})
	return result, nil
}

func (s *Service) RequestProcessPageable(request map[string]interface{}, action string, method string, pageItemJsonpath string) (object []map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	for {
		wait := helper.IncrementalWait(5*time.Second, 20*time.Second)
		err = retry.Retry(10*time.Minute, func() *retry.RetryError {
			if method == "GET" {
				response, err = client.RpcGet("selectdb", "2023-05-22", action, request, nil)
			} else {
				response, err = client.RpcPost("selectdb", "2023-05-22", action, nil, request, true)
			}
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) ||
					errs.IsExpectedErrors(err, []string{"InvalidDBInstanceState.NotSupport"}) ||
					errs.IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) ||
					helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, request, action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get(pageItemJsonpath, response)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, pageItemJsonpath, response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return objects, nil
}

func (s *Service) DescribeDbCluster(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceAttribute"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	if parts[1] == "" {
		return nil, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcess(request, action, "GET")
	if err != nil {
		return nil, err
	}
	clusterInfo := make(map[string]interface{})
	clusterIndex := v["DBClusterList"].([]interface{})
	for _, w := range clusterIndex {
		ws := w.(map[string]interface{})
		if ws["DbClusterId"] == parts[1] {
			clusterInfo = ws
		}
	}
	if len(clusterInfo) == 0 {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DbClusterId", v)
	}
	return clusterInfo, nil
}

func (s *Service) DbClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDbCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *Service) DescribeDbInstance(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceAttribute"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcess(request, action, "GET")

}

func (s *Service) DbInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDbInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *Service) DescribeDbInstances(ids string, tags []map[string]interface{}) (objects []map[string]interface{}, err error) {

	action := "DescribeDBInstances"

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	if len(ids) > 0 {
		instance_ids := strings.Replace(ids, ":", ",", -1)
		request["DBInstanceIds"] = instance_ids
	}
	if tags != nil {
		paramTag, _ := json.Marshal(tags)
		param := string(paramTag)
		request["Tag"] = param
	}
	pageItemJsonpath := "$.Items"

	return s.RequestProcessPageable(request, action, "GET", pageItemJsonpath)
}

func (s *Service) DescribeDBInstanceAccessWhiteList(id string) (object map[string]interface{}, err error) {

	action := "DescribeSecurityIPList"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}
	v, err := s.RequestProcess(request, action, "GET")
	if err != nil {
		return nil, err
	}
	clusterInfo := v["GroupItems"].(map[string]interface{})
	return clusterInfo, nil
}

func (s *Service) UpdateDbClusterConfig(id string, config map[string]string) (object map[string]interface{}, err error) {

	action := "ModifyDBClusterConfig"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	configKey := ""
	if parts[0]+"-fe" == parts[1] {
		configKey = "fe.conf"
	} else {
		configKey = "be.conf"
	}
	paramJson, _ := json.Marshal(config)
	param := string(paramJson)
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
		"ConfigKey":    configKey,
		"Parameters":   param,
	}

	v, err := s.RequestProcess(request, action, "POST")
	if err != nil {
		return nil, err
	}
	return v, nil

}

func (s *Service) DescribeDbClusterConfig(id string) (object []interface{}, err error) {

	action := "DescribeDBClusterConfig"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	configKey := ""
	if parts[0]+"-fe" == parts[1] {
		configKey = "fe.conf"
	} else {
		configKey = "be.conf"
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
		"ConfigKey":    configKey,
	}

	response, err := s.RequestProcess(request, action, "GET")
	if err != nil {
		return nil, err
	}
	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return nil, errs.WrapErrorf(err, errs.IdMsg, "DescribeDBClusterConfig")
	} else {
		clusterInfo := resp.(map[string]interface{})["Params"].([]interface{})
		return clusterInfo, nil
	}

}

func (s *Service) DescribeDbClusterConfigChangeLog(id string) (object []interface{}, err error) {

	action := "DescribeDBClusterConfigChangeLogs"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	configKey := ""
	if parts[0]+"-fe" == parts[1] {
		configKey = "fe.conf"
	} else {
		configKey = "be.conf"
	}
	currentTime := time.Now().Format(time.DateTime)

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
		"ConfigKey":    configKey,
		"StartTime":    "1970-01-01 10:00:00",
		"EndTime":      currentTime,
	}
	v, err := s.RequestProcess(request, action, "GET")
	if err != nil {
		return nil, err
	}
	baseInfo := v["Data"]
	if baseInfo == nil {
		return nil, nil
	}
	clusterInfo := baseInfo.(map[string]interface{})["ParamChangeLogs"]
	if clusterInfo == nil {
		return nil, nil
	}

	return clusterInfo.([]interface{}), nil
}

func (s *Service) DeleteCluster(id string) (object map[string]interface{}, err error) {

	action := "DeleteDBCluster"

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) DeleteInstance(id string) (object map[string]interface{}, err error) {

	action := "DeleteDBInstance"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) ModifyClusterDescription(id string, newDescription string) (object map[string]interface{}, err error) {

	action := "ModifyBEClusterAttribute"

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId":          parts[0],
		"DBClusterId":           parts[1],
		"RegionId":              s.client.RegionId,
		"InstanceAttributeType": "DBInstanceDescription",
		"Value":                 newDescription,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) ModifyCluster(id string, newClass string, newCacheSize int) (object map[string]interface{}, err error) {

	action := "ModifyDBCluster"

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
		"Engine":       "selectdb",
	}

	if newClass != "" {
		request["DBClusterClass"] = newClass
	}

	if newCacheSize > 100 {
		request["CacheSize"] = newCacheSize
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) UpdateClusterStatus(id string, targetStatus string) (object map[string]interface{}, err error) {
	action := ""
	switch targetStatus {
	case "STOPPING", "STOPPED":
		action = "StopBECluster"
	case "STARTING":
		action = "StartBECluster"
	case "RESTART", "RESTARTING":
		action = "RestartDBCluster"
	}
	if action == "" {
		return nil, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, targetStatus))
	}

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) DescribeDbInstanceNetInfo(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceNetInfo"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcess(request, action, "GET")
	if err != nil {
		return nil, err
	}
	response := map[string]interface{}{
		"DBClustersNetInfos": v["DBClustersNetInfos"].([]interface{}),
		"DBInstanceNetInfos": v["DBInstanceNetInfos"].([]interface{}),
	}
	return response, nil

}

func (s *Service) DescribeDbInstanceSecurityIPList(id string) (object []interface{}, err error) {

	action := "DescribeSecurityIPList"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcess(request, action, "GET")
	if err != nil {
		return nil, err
	}
	clusterInfo := v["GroupItems"].([]interface{})
	return clusterInfo, nil

}

func (s *Service) ModifyDbInstanceSecurityIPList(id string, groupName string, newIpList string) (object map[string]interface{}, err error) {

	action := "ModifySecurityIPList"

	request := map[string]interface{}{
		"DBInstanceId":   id,
		"RegionId":       s.client.RegionId,
		"GroupName":      groupName,
		"SecurityIPList": newIpList,
		"ModifyMode":     0,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) ModifyInstanceDescription(id string, newDescription string) (object map[string]interface{}, err error) {

	action := "ModifyDBInstanceAttribute"

	request := map[string]interface{}{
		"DBInstanceId":          id,
		"RegionId":              s.client.RegionId,
		"InstanceAttributeType": "DBInstanceDescription",
		"Value":                 newDescription,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) ModifyInstancePaymentType(id string, paymentRequest map[string]string) (object map[string]interface{}, err error) {

	action := "ModifyDBInstancePayType"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
		"ChargeType":   paymentRequest["payment_type"],
	}
	period, exist := paymentRequest["period"]
	if exist {
		request["Period"] = period
	}
	usedTime, exist := paymentRequest["period_time"]
	if exist {
		request["usedTime"] = usedTime
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) UpgradeInstanceEngineVersion(id string, version string, upgradeInMaintainTime bool) (object map[string]interface{}, err error) {

	action := "UpgradeDBInstanceEngineVersion"
	switchTimeMode := 0

	if upgradeInMaintainTime {
		switchTimeMode = 1
	}

	request := map[string]interface{}{
		"DBInstanceId":   id,
		"RegionId":       s.client.RegionId,
		"EngineVersion":  version,
		"SwitchTimeMode": switchTimeMode,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) AllocateInstancePublicConnection(id string) (object map[string]interface{}, err error) {

	action := "AllocateInstancePublicConnection"

	request := map[string]interface{}{
		"DBInstanceId":           id,
		"RegionId":               s.client.RegionId,
		"ConnectionStringPrefix": id + "-public",
		"NetType":                "Public",
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) ReleaseInstancePublicConnection(id string, connectionString string) (object map[string]interface{}, err error) {

	action := "ReleaseInstancePublicConnection"

	request := map[string]interface{}{
		"DBInstanceId":     id,
		"RegionId":         s.client.RegionId,
		"ConnectionString": connectionString,
	}

	return s.RequestProcess(request, action, "POST")

}

func (s *Service) SetResourceTags(id string, added map[string]interface{}, removed []string) error {

	resourceType := "dbinstance"
	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !tags.IgnoredTags(v, "") {
			removedTagKeys = append(removedTagKeys, v)
		}
	}
	if len(removedTagKeys) > 0 {
		action := "UntagResources"
		request := map[string]interface{}{
			"RegionId":     s.client.RegionId,
			"ResourceType": resourceType,
			"ResourceId.1": id,
		}
		for i, key := range removedTagKeys {
			request[fmt.Sprintf("TagKey.%d", i+1)] = key
		}
		_, err := s.RequestProcess(request, action, "POST")
		if err != nil {
			return err
		}
	}
	if len(added) > 0 {
		action := "TagResources"
		request := map[string]interface{}{
			"RegionId":     s.client.RegionId,
			"ResourceType": resourceType,
			"ResourceId.1": id,
		}
		count := 1
		for key, value := range added {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
		_, err := s.RequestProcess(request, action, "POST")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ModifyInstanceAdminPass(id string, adminpass string) (object map[string]interface{}, err error) {

	action := "ResetAccountPassword"

	request := map[string]interface{}{
		"DBInstanceId":    id,
		"RegionId":        s.client.RegionId,
		"AccountName":     "admin",
		"AccountPassword": adminpass,
	}

	return s.RequestProcess(request, action, "GET")
}
