package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type SelectDBService struct {
	client *connectivity.AliyunClient
}

func (s *SelectDBService) RequestProcessForSelectDB(request map[string]interface{}, action string, method string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	wait := incrementalWait(5*time.Second, 20*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		if method == "GET" {
			response, err = client.RpcGet("selectdb", "2023-05-22", action, request, nil)
		} else {
			response, err = client.RpcPost("selectdb", "2023-05-22", action, nil, request, true)
		}
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) ||
				IsExpectedErrors(err, []string{"InvalidDBInstanceState.NotSupport"}) ||
				IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) ||
				NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, request, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return nil, WrapErrorf(err, FailedGetAttributeMsg, request, "$", response)
	}
	result := v.(map[string]interface{})
	return result, nil
}

func (s *SelectDBService) RequestProcessPageableForSelectDB(request map[string]interface{}, action string, method string, pageItemJsonpath string) (object []map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	for {
		wait := incrementalWait(5*time.Second, 20*time.Second)
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			if method == "GET" {
				response, err = client.RpcGet("selectdb", "2023-05-22", action, request, nil)
			} else {
				response, err = client.RpcPost("selectdb", "2023-05-22", action, nil, request, true)
			}
			if err != nil {
				if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) ||
					IsExpectedErrors(err, []string{"InvalidDBInstanceState.NotSupport"}) ||
					IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) ||
					NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, request, action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get(pageItemJsonpath, response)
		if err != nil {
			return nil, WrapErrorf(err, FailedGetAttributeMsg, action, pageItemJsonpath, response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return objects, nil
}

func (s *SelectDBService) DescribeSelectDBDbCluster(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceAttribute"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	if parts[1] == "" {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "GET")
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
		return nil, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DbClusterId", v)
	}
	return clusterInfo, nil
}

func (s *SelectDBService) SelectDBDbClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSelectDBDbCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *SelectDBService) DescribeSelectDBDbInstance(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceAttribute"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "GET")

}

func (s *SelectDBService) SelectDBDbInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSelectDBDbInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *SelectDBService) DescribeSelectDBDbInstances(ids string, tags []map[string]interface{}) (objects []map[string]interface{}, err error) {

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

	return s.RequestProcessPageableForSelectDB(request, action, "GET", pageItemJsonpath)
}

func (s *SelectDBService) DescribeDBInstanceAccessWhiteList(id string) (object map[string]interface{}, err error) {

	action := "DescribeSecurityIPList"
	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}
	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return nil, err
	}
	clusterInfo := v["GroupItems"].(map[string]interface{})
	return clusterInfo, nil
}

func (s *SelectDBService) UpdateSelectDBDbClusterConfig(id string, config map[string]string) (object map[string]interface{}, err error) {

	action := "ModifyDBClusterConfig"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
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

	v, err := s.RequestProcessForSelectDB(request, action, "POST")
	if err != nil {
		return nil, err
	}
	return v, nil

}

func (s *SelectDBService) DescribeSelectDBDbClusterConfig(id string) (object []interface{}, err error) {

	action := "DescribeDBClusterConfig"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
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

	response, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return nil, err
	}
	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return nil, WrapErrorf(err, IdMsg, "DescribeDBClusterConfig")
	} else {
		clusterInfo := resp.(map[string]interface{})["Params"].([]interface{})
		return clusterInfo, nil
	}

}

func (s *SelectDBService) DescribeSelectDBDbClusterConfigChangeLog(id string) (object []interface{}, err error) {

	action := "DescribeDBClusterConfigChangeLogs"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
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
	v, err := s.RequestProcessForSelectDB(request, action, "GET")
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

func (s *SelectDBService) DeleteSelectDBCluster(id string) (object map[string]interface{}, err error) {

	action := "DeleteDBCluster"

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) DeleteSelectDBInstance(id string) (object map[string]interface{}, err error) {

	action := "DeleteDBInstance"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBClusterDescription(id string, newDescription string) (object map[string]interface{}, err error) {

	action := "ModifyBEClusterAttribute"

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId":          parts[0],
		"DBClusterId":           parts[1],
		"RegionId":              s.client.RegionId,
		"InstanceAttributeType": "DBInstanceDescription",
		"Value":                 newDescription,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBCluster(id string, newClass string, newCacheSize int) (object map[string]interface{}, err error) {

	action := "ModifyDBCluster"

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
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

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) UpdateSelectDBClusterStatus(id string, targetStatus string) (object map[string]interface{}, err error) {
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
		return nil, WrapError(Error(FailedToReachTargetStatus, targetStatus))
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"DBClusterId":  parts[1],
		"RegionId":     s.client.RegionId,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) DescribeSelectDBDbInstanceNetInfo(id string) (object map[string]interface{}, err error) {

	action := "DescribeDBInstanceNetInfo"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return nil, err
	}
	response := map[string]interface{}{
		"DBClustersNetInfos": v["DBClustersNetInfos"].([]interface{}),
		"DBInstanceNetInfos": v["DBInstanceNetInfos"].([]interface{}),
	}
	return response, nil

}

func (s *SelectDBService) DescribeSelectDBDbInstanceSecurityIPList(id string) (object []interface{}, err error) {

	action := "DescribeSecurityIPList"

	request := map[string]interface{}{
		"DBInstanceId": id,
		"RegionId":     s.client.RegionId,
	}

	v, err := s.RequestProcessForSelectDB(request, action, "GET")
	if err != nil {
		return nil, err
	}
	clusterInfo := v["GroupItems"].([]interface{})
	return clusterInfo, nil

}

func (s *SelectDBService) ModifySelectDBDbInstanceSecurityIPList(id string, groupName string, newIpList string) (object map[string]interface{}, err error) {

	action := "ModifySecurityIPList"

	request := map[string]interface{}{
		"DBInstanceId":   id,
		"RegionId":       s.client.RegionId,
		"GroupName":      groupName,
		"SecurityIPList": newIpList,
		"ModifyMode":     0,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBInstanceDescription(id string, newDescription string) (object map[string]interface{}, err error) {

	action := "ModifyDBInstanceAttribute"

	request := map[string]interface{}{
		"DBInstanceId":          id,
		"RegionId":              s.client.RegionId,
		"InstanceAttributeType": "DBInstanceDescription",
		"Value":                 newDescription,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ModifySelectDBInstancePaymentType(id string, paymentRequest map[string]string) (object map[string]interface{}, err error) {

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

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) UpgradeSelectDBInstanceEngineVersion(id string, version string, upgradeInMaintainTime bool) (object map[string]interface{}, err error) {

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

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) AllocateSelectDBInstancePublicConnection(id string) (object map[string]interface{}, err error) {

	action := "AllocateInstancePublicConnection"

	request := map[string]interface{}{
		"DBInstanceId":           id,
		"RegionId":               s.client.RegionId,
		"ConnectionStringPrefix": id + "-public",
		"NetType":                "Public",
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) ReleaseSelectDBInstancePublicConnection(id string, connectionString string) (object map[string]interface{}, err error) {

	action := "ReleaseInstancePublicConnection"

	request := map[string]interface{}{
		"DBInstanceId":     id,
		"RegionId":         s.client.RegionId,
		"ConnectionString": connectionString,
	}

	return s.RequestProcessForSelectDB(request, action, "POST")

}

func (s *SelectDBService) SetResourceTags(id string, added map[string]interface{}, removed []string) error {

	resourceType := "dbinstance"
	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !ignoredTags(v, "") {
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
		_, err := s.RequestProcessForSelectDB(request, action, "POST")
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
		_, err := s.RequestProcessForSelectDB(request, action, "POST")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SelectDBService) ModifySelectDBInstanceAdminPass(id string, adminpass string) (object map[string]interface{}, err error) {

	action := "ResetAccountPassword"

	request := map[string]interface{}{
		"DBInstanceId":    id,
		"RegionId":        s.client.RegionId,
		"AccountName":     "admin",
		"AccountPassword": adminpass,
	}

	return s.RequestProcessForSelectDB(request, action, "GET")
}
