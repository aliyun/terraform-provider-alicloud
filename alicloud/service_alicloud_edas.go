package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EdasService struct {
	client *connectivity.AliyunClient
}

type Hook struct {
	Exec      *Exec      `json:"exec,omitempty"`
	HttpGet   *HttpGet   `json:"httpGet,omitempty"`
	TcpSocket *TcpSocket `json:"tcpSocket,omitempty"`
}

type Exec struct {
	Command []string `json:"command"`
}

type HttpGet struct {
	Path        string       `json:"path"`
	Port        int          `json:"port"`
	Scheme      string       `json:"scheme"`
	HttpHeaders []HttpHeader `json:"httpHeaders"`
}

type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TcpSocket struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Prober struct {
	FailureThreshold    int `json:"failureThreshold"`
	InitialDelaySeconds int `json:"initialDelaySeconds"`
	SuccessThreshold    int `json:"successThreshold"`
	TimeoutSeconds      int `json:"timeoutSeconds"`
	Hook                `json:",inline"`
}

const (
	ChangeOrderStatusReady            = "0"
	ChangeOrderStatusRunning          = "1"
	ChangeOrderStatusSuccess          = "2"
	ChangeOrderStatusFailed           = "3"
	ChangeOrderStatusAbort            = "6"
	ChangeOrderStatusWaitBatchConfirm = "8"
	ChangeOrderStatusAutoBatchWait    = "9"
	ChangeOrderStatusSystemFail       = "10"
)

const (
	ChangeOrderStatusReadyStr            = "Ready"
	ChangeOrderStatusRunningStr          = "Running"
	ChangeOrderStatusSuccessStr          = "Success"
	ChangeOrderStatusFailedStr           = "Failed"
	ChangeOrderStatusAbortStr            = "Abort"
	ChangeOrderStatusWaitBatchConfirmStr = "WaitBatchConfirm"
	ChangeOrderStatusAutoBatchWaitStr    = "AutoBatchWait"
	ChangeOrderStatusSystemFailStr       = "SystemFail"
)

var ChangeOrderStatusMap = map[int]string{
	0:  ChangeOrderStatusReadyStr,
	1:  ChangeOrderStatusRunningStr,
	2:  ChangeOrderStatusSuccessStr,
	3:  ChangeOrderStatusFailedStr,
	6:  ChangeOrderStatusAbortStr,
	8:  ChangeOrderStatusWaitBatchConfirmStr,
	9:  ChangeOrderStatusAutoBatchWaitStr,
	10: ChangeOrderStatusSystemFailStr,
}

func (e *EdasService) GetChangeOrderStatus(id string) (info *edas.ChangeOrderInfo, err error) {
	request := edas.CreateGetChangeOrderInfoRequest()
	request.RegionId = e.client.RegionId
	request.ChangeOrderId = id

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetChangeOrderInfo(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"OperationDenied.InvalidDBClusterIdNotFound", "OperationDenied.InvalidDBClusterNameNotFound"}) {
			return info, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return info, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp := raw.(*edas.GetChangeOrderInfoResponse)
	return &rsp.ChangeOrderInfo, nil

}
func (e *EdasService) GetChangeOrderStatusV2(id string) (int, error) {
	if response, err := e.doPopRequest(map[string]string{
		"action":     "/pop/v5/changeorder/change_order_info",
		"httpMethod": "POST",
		"id":         id,
	}, map[string]string{
		"ChangeOrderId": id,
	}); err != nil {
		return -1, err
	} else {
		if v, err := jsonpath.Get("$.changeOrderInfo.Status", response); err != nil {
			return -1, err
		} else {
			return parseInt(v)
		}
	}
}

func (e *EdasService) GetDeployGroup(appId, groupId string) (groupInfo *edas.DeployGroup, err error) {
	request := edas.CreateListDeployGroupRequest()
	request.RegionId = e.client.RegionId
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	if err != nil {
		return groupInfo, WrapErrorf(err, DefaultErrorMsg, appId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp := raw.(*edas.ListDeployGroupResponse)
	if rsp.Code != 200 {
		return groupInfo, Error("get deploy group failed for " + rsp.Message)
	}
	for _, group := range rsp.DeployGroupList.DeployGroup {
		if group.GroupId == groupId {
			return &group, nil
		}
	}
	return groupInfo, nil
}

func (e *EdasService) EdasChangeOrderStatusRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := e.GetChangeOrderStatus(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if strconv.Itoa(object.Status) == failState {
				return object, strconv.Itoa(object.Status), WrapError(Error(FailedToReachTargetStatus, strconv.Itoa(object.Status)))
			}
		}

		return object, strconv.Itoa(object.Status), nil
	}
}

func (e *EdasService) EdasChangeOrderStatusRefreshFuncV2(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		status, err := e.GetChangeOrderStatusV2(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		statusStr := "not found"
		if v, ok := ChangeOrderStatusMap[status]; ok {
			statusStr = v
		}
		for _, failState := range failStates {
			if statusStr == failState {
				return nil, statusStr, WrapError(Error(FailedToReachTargetStatus, statusStr))
			}
		}
		return status, statusStr, nil
	}
}

func (e *EdasService) WaitForChangeOrderFinished(resourceId string, changeOrderId string, timeout time.Duration) error {
	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf(
			[]string{ChangeOrderStatusReadyStr, ChangeOrderStatusRunningStr},
			[]string{ChangeOrderStatusSuccessStr},
			timeout,
			5*time.Second,
			e.EdasChangeOrderStatusRefreshFuncV2(changeOrderId,
				[]string{
					ChangeOrderStatusFailedStr,
					ChangeOrderStatusAbortStr,
					ChangeOrderStatusSystemFailStr}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, resourceId)
		}
	}
	return nil
}

func (e *EdasService) WaitForChangeOrderFinishedNonRetryable(resourceId string, changeOrderId string, timeout time.Duration) *resource.RetryError {
	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf(
			[]string{ChangeOrderStatusReadyStr, ChangeOrderStatusRunningStr},
			[]string{ChangeOrderStatusSuccessStr},
			timeout,
			5*time.Second,
			e.EdasChangeOrderStatusRefreshFuncV2(changeOrderId,
				[]string{
					ChangeOrderStatusFailedStr,
					ChangeOrderStatusAbortStr,
					ChangeOrderStatusSystemFailStr}))
		if _, err := stateConf.WaitForState(); err != nil {
			return resource.NonRetryableError(WrapErrorf(err, "change order failed, app id %s", resourceId))
		}
	}
	return nil
}

func parseInt(obj interface{}) (int, error) {
	if obj == nil {
		return -1, fmt.Errorf("try to parse int, but got nil")
	}
	switch obj.(type) {
	case string:
		return strconv.Atoi(obj.(string))
	case int:
		return obj.(int), nil
	case int32:
		return int(obj.(int32)), nil
	case int64:
		return int(obj.(int64)), nil
	case json.Number:
		return strconv.Atoi(obj.(json.Number).String())
	default:
		return -1, fmt.Errorf("unknown type of object: %v", reflect.TypeOf(obj))
	}
}

// HasOngoingTasks check if there is an ongoing task, ignore all api error, just return false
func (e *EdasService) HasOngoingTasks(appId string) bool {
	if response, err := e.doPopRequest(map[string]string{
		"action":     "/pop/v5/changeorder/change_order_list",
		"httpMethod": "POST",
		"id":         appId,
	}, map[string]string{
		"AppId": appId,
	}); err != nil {
		return false
	} else if v, err := jsonpath.Get("$.ChangeOrderList.ChangeOrder", response); err != nil {
		return false
	} else {
		if len(v.([]interface{})) < 1 {
			return false
		}
		for _, co := range v.([]interface{}) {
			changeOrder := co.(map[string]interface{})
			status := changeOrder["Status"]
			if v, err := parseInt(status); err != nil {
				return false
			} else if statusStr, exist := ChangeOrderStatusMap[v]; exist {
				if statusStr == ChangeOrderStatusReadyStr ||
					statusStr == ChangeOrderStatusRunningStr ||
					statusStr == ChangeOrderStatusAutoBatchWaitStr ||
					statusStr == ChangeOrderStatusWaitBatchConfirmStr {
					return true
				}
			} else {
				return false
			}
		}
	}
	return false
}

func (e *EdasService) SyncResource(resourceType string) error {
	request := edas.CreateSynchronizeResourceRequest()
	request.RegionId = e.client.RegionId
	request.Type = resourceType

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.SynchronizeResource(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "sync resource", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp := raw.(*edas.SynchronizeResourceResponse)
	if rsp.Code != 200 || !rsp.Success {
		return WrapError(Error("sync resource failed for " + rsp.Message))
	}

	return nil
}

func (e *EdasService) CheckEcsStatus(instanceIds string, count int) error {
	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = e.client.RegionId
	request.Status = "Running"
	request.PageSize = requests.NewInteger(100)
	request.InstanceIds = instanceIds

	raw, err := e.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstances(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"OperationDenied.InvalidDBClusterIdNotFound", "OperationDenied.InvalidDBClusterNameNotFound"}) {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, instanceIds, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	rsp := raw.(*ecs.DescribeInstancesResponse)

	if len(rsp.Instances.Instance) != count {
		return WrapErrorf(Error("not enough instances"), DefaultErrorMsg, instanceIds, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}

func (e *EdasService) GetLastPackgeVersion(appId, groupId string) (string, error) {
	var versionId string
	request := edas.CreateQueryApplicationStatusRequest()
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryApplicationStatus(request)
	})
	if err != nil {
		return "", WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_application_package_version", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	response, _ := raw.(*edas.QueryApplicationStatusResponse)

	if response.Code != 200 {
		return "", WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}

	for _, group := range response.AppInfo.GroupList.Group {
		if group.GroupId == groupId {
			versionId = group.PackageVersionId
		}
	}

	rq := edas.CreateListHistoryDeployVersionRequest()
	rq.AppId = appId

	raw, err = e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListHistoryDeployVersion(rq)
	})
	if err != nil {
		return "", WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_application_package_version_list", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	rsp, _ := raw.(*edas.ListHistoryDeployVersionResponse)

	if rsp.Code != 200 {
		return "", WrapError(Error("QueryApplicationStatus failed for " + response.Message))
	}

	for _, version := range rsp.PackageVersionList.PackageVersion {
		if version.Id == versionId {
			return version.PackageVersion, nil
		}
	}

	return "", nil
}

func (e *EdasService) DescribeEdasApplication(appId string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	regionId := e.client.RegionId

	request := edas.CreateGetApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(request)
	})
	if err != nil {
		return application, WrapError(err)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetApplicationResponse)
	if response.Code != 200 {
		return application, WrapError(Error("get application error :" + response.Message))
	}

	v := response.Applcation

	return &v, nil
}

func (e *EdasService) DescribeEdasApplicationV2(appId string) (object map[string]interface{}, err error) {
	if response, err := e.doPopRequest(map[string]string{
		"action":     "/pop/v5/app/app_info",
		"httpMethod": "POST",
	}, map[string]string{
		"id":    appId,
		"AppId": appId,
	}); err != nil {
		return nil, err
	} else if application, exist := response["Application"].(map[string]interface{}); exist {
		return application, nil
	} else {
		return nil, fmt.Errorf("%s: Edas K8s Applicatio, AppId: %s", ResourceNotfound, appId)
	}
}

func (e *EdasService) doPopRequest(requestConfig map[string]string, params map[string]string) (map[string]interface{}, error) {
	var response map[string]interface{}
	client := e.client
	request := map[string]*string{}
	for k := range params {
		v := params[k]
		request[k] = &v
	}
	action := requestConfig["action"]
	httpMethod := requestConfig["httpMethod"]
	id := requestConfig["id"]

	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		switch httpMethod {
		case "POST":
			response, err = client.RoaPost("Edas", "2017-08-01", action, request, nil, nil, false)
		case "PUT":
			response, err = client.RoaPut("Edas", "2017-08-01", action, request, nil, nil, false)
		case "DELETE":
			response, err = client.RoaDelete("Edas", "2017-08-01", action, request, nil, nil, false)
		case "GET":
			response, err = client.RoaGet("Edas", "2017-08-01", action, request, nil, nil)
		default:
			response, err = nil, fmt.Errorf("httpMethod %s is not support", httpMethod)
		}
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (e *EdasService) DescribeEdasCluster(clusterId string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}
	regionId := e.client.RegionId

	request := edas.CreateGetClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	if err != nil {
		return cluster, WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetClusterResponse)
	if response.Code != 200 {
		return cluster, WrapError(Error("create cluster failed for " + response.Message))
	}

	v := response.Cluster

	return &v, nil
}

func (e *EdasService) DescribeEdasDeployGroup(id string) (*edas.DeployGroup, error) {
	group := &edas.DeployGroup{}
	regionId := e.client.RegionId

	strs := strings.Split(id, ":")

	request := edas.CreateListDeployGroupRequest()
	request.RegionId = regionId
	request.AppId = strs[0]

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	if err != nil {
		return group, WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_deploy_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.ListDeployGroupResponse)
	if response.Code != 200 {
		return group, WrapError(Error("create cluster failed for " + response.Message))
	}

	for _, v := range response.DeployGroupList.DeployGroup {
		if v.ClusterName == strs[1] {
			return &v, nil
		}
	}

	return group, nil
}

func (e *EdasService) DescribeEdasInstanceClusterAttachment(id string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasCluster(v[0])
	if err != nil {
		return cluster, WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasApplicationDeployment(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasApplicationScale(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasSlbAttachment(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, WrapError(err)
	}

	return o, nil
}

type CommandArg struct {
	Argument string `json:"argument" xml:"argument"`
}

func (e *EdasService) GetK8sCommandArgs(args []interface{}) (string, error) {
	aString := make([]CommandArg, 0)
	for _, v := range args {
		aString = append(aString, CommandArg{Argument: v.(string)})
	}
	b, err := json.Marshal(aString)
	if err != nil {
		return "", WrapError(err)
	}
	return string(b), nil
}

func (e *EdasService) GetK8sCommandArgsForDeploy(args []interface{}) (string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return "", WrapError(err)
	}
	return string(b), nil
}

type K8sEnv struct {
	Name  string `json:"name" xml:"name"`
	Value string `json:"value" xml:"value"`
}

func (e *EdasService) GetK8sEnvs(envs map[string]interface{}) (string, error) {
	k8sEnvs := make([]K8sEnv, 0)
	for n, v := range envs {
		k8sEnvs = append(k8sEnvs, K8sEnv{Name: n, Value: v.(string)})
	}

	b, err := json.Marshal(k8sEnvs)
	if err != nil {
		return "", WrapError(err)
	}
	return string(b), nil
}

func (e *EdasService) QueryK8sAppPackageType(appId string) (string, error) {
	request := edas.CreateGetApplicationRequest()
	request.RegionId = e.client.RegionId
	request.AppId = appId
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(request)
	})

	if err != nil {
		return "", WrapError(err)
	}

	response, _ := raw.(*edas.GetApplicationResponse)
	if response.Code != 200 {
		return "", WrapError(Error("get application for appId:" + appId + " failed:" + response.Message))
	}
	if len(response.Applcation.ApplicationType) > 0 {
		return response.Applcation.ApplicationType, nil
	}
	return "", WrapError(Error("not package type for appId:" + appId))
}

func (e *EdasService) DescribeEdasK8sCluster(clusterId string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}
	regionId := e.client.RegionId

	request := edas.CreateGetClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	if err != nil {
		return cluster, WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetClusterResponse)
	if response.Code != 200 {
		if strings.Contains(response.Message, "does not exist") {
			return cluster, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return cluster, WrapError(Error("create k8s cluster failed for " + response.Message))
	}

	v := response.Cluster

	return &v, nil
}

func (e *EdasService) DescribeEdasK8sApplication(appId string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	regionId := e.client.RegionId

	request := edas.CreateGetK8sApplicationRequest()
	request.RegionId = regionId
	request.AppId = appId

	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetK8sApplication(request)
	})
	if err != nil {
		return application, WrapError(err)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetK8sApplicationResponse)
	if response.Code != 200 {
		if strings.Contains(response.Message, "does not exist") {
			return application, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return application, WrapError(Error("get k8s application error :" + response.Message))
	}

	v := response.Applcation

	return &v, nil
}

func (e *EdasService) PreStopEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldHook Hook
	err := json.Unmarshal([]byte(oldStr), &oldHook)
	if err != nil {
		return false
	}
	var newHook Hook
	err = json.Unmarshal([]byte(newStr), &newHook)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldHook, newHook)
}

func (e *EdasService) PostStartEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldHook Hook
	err := json.Unmarshal([]byte(oldStr), &oldHook)
	if err != nil {
		return false
	}
	var newHook Hook
	err = json.Unmarshal([]byte(newStr), &newHook)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldHook, newHook)
}

func (e *EdasService) LivenessEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldProber Prober
	err := json.Unmarshal([]byte(oldStr), &oldProber)
	if err != nil {
		return false
	}
	var newProber Prober
	err = json.Unmarshal([]byte(newStr), &newProber)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldProber, newProber)
}

func (e *EdasService) ReadinessEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldProber Prober
	err := json.Unmarshal([]byte(oldStr), &oldProber)
	if err != nil {
		return false
	}
	var newProber Prober
	err = json.Unmarshal([]byte(newStr), &newProber)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldProber, newProber)
}
func (e *EdasService) IsJsonEqual(v1, v2 interface{}) bool {
	s1 := v1.(string)
	s2 := v2.(string)
	var i1 interface{}
	err := json.Unmarshal([]byte(s1), &i1)
	if err != nil {
		return false
	}
	var i2 interface{}
	err = json.Unmarshal([]byte(s2), &i2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(&i1, &i2)
}

func (e *EdasService) IsJsonArrayEqual(v1, v2 interface{}) bool {
	s1 := v1.(string)
	s2 := v2.(string)
	var i1 []interface{}
	err := json.Unmarshal([]byte(s1), &i1)
	if err != nil {
		return false
	}
	var i2 []interface{}
	err = json.Unmarshal([]byte(s2), &i2)
	if err != nil {
		return false
	}
	if len(i1) != len(i2) {
		return false
	}
	for _, obj := range i1 {
		found := false
		for _, obj2 := range i2 {
			if reflect.DeepEqual(obj, obj2) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func (s *EdasService) DescribeEdasNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v5/user_region_defs"
	request := map[string]*string{}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RoaPost("Edas", "2017-08-01", action, request, nil, nil, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.UserDefineRegionList.UserDefineRegionEntity", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.UserDefineRegionList.UserDefineRegionEntity", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("EDAS", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Id"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("EDAS", id)), NotFoundWithResponse, response)
	}
	return object, nil
}
func (e *EdasService) BindK8sSlb(appId string, slbConfig *map[string]interface{}, timeout time.Duration) *resource.RetryError {
	if e.HasOngoingTasks(appId) {
		return resource.RetryableError(Error("there is an ongoing task"))
	}
	portMappings := (*slbConfig)["port_mappings"].(*schema.Set).List()
	var requestServicePortInfos []map[string]interface{}
	if len(portMappings) > 0 {
		for _, pmInterface := range portMappings {
			pm := pmInterface.(map[string]interface{})
			if v, ok := pm["service_port"]; ok {
				spSlice := v.(*schema.Set).List()
				sp := spSlice[0].(map[string]interface{})
				info := map[string]interface{}{
					"certId":               pm["cert_id"],
					"loadBalancerProtocol": pm["loadbalancer_protocol"],
					"targetPort":           sp["target_port"],
					"port":                 sp["port"],
				}
				requestServicePortInfos = append(requestServicePortInfos, info)
			}
		}
	}
	portInfoBytes, err := json.Marshal(requestServicePortInfos)
	if err != nil {
		return resource.NonRetryableError(WrapErrorf(err, "marshal port mappings failed, value %v", requestServicePortInfos))
	}
	params := map[string]string{
		"AppId":            appId,
		"Type":             (*slbConfig)["type"].(string),
		"SlbId":            (*slbConfig)["slb_id"].(string),
		"Scheduler":        (*slbConfig)["scheduler"].(string),
		"Specification":    (*slbConfig)["specification"].(string),
		"ServicePortInfos": string(portInfoBytes),
	}
	if response, err := e.doPopRequest(map[string]string{
		"action":     "/pop/v5/k8s/acs/k8s_slb_binding",
		"httpMethod": "POST",
		"id":         appId,
	}, params); err != nil {
		return resource.NonRetryableError(err)
	} else {
		return e.WaitForChangeOrderFinishedNonRetryable(appId, response["ChangeOrderId"].(string), timeout)
	}
}

func (e *EdasService) UpdateK8sAppSlbInfos(appId string, oldInfos, newInfos *[]interface{}, timeout time.Duration) *resource.RetryError {
	buildMap := func(infos *[]interface{}) *map[string]interface{} {
		m := map[string]interface{}{}
		if len(*infos) > 0 {
			for _, v := range *infos {
				info := v.(map[string]interface{})
				m[info["name"].(string)] = info
			}
		}
		return &m
	}
	// build a map with key: slbName, value: slbInfo
	oldInfosMap := buildMap(oldInfos)
	newInfosMap := buildMap(newInfos)

	// find removed slb configs and unbind them
	for _, v := range *oldInfos {
		oInfo := v.(map[string]interface{})
		slbName := oInfo["name"].(string)
		slbType := oInfo["type"].(string)
		if _, ok := (*newInfosMap)[slbName]; !ok {
			if err := e.UnbindK8sSlb(appId, slbType, slbName, timeout); err != nil {
				return err
			}
		}
	}
	// find new or modified slb configs
	for _, v := range *newInfos {
		nInfo := v.(map[string]interface{})
		slbName := nInfo["name"].(string)
		if ov, ok := (*oldInfosMap)[slbName]; ok {
			oInfo := ov.(map[string]interface{})
			if pm, ok := nInfo["port_mappings"].([]interface{}); !ok || len(pm) == 0 {
				// in this case, the port_mappings were cleared, so unbind this slb

				// type is required to unbind slb, but it is cleared in nInfo
				slbType := oInfo["type"].(string)
				if err := e.UnbindK8sSlb(appId, slbType, slbName, timeout); err != nil {
					return err
				}
				continue
			}
			if !reflect.DeepEqual(oInfo, nInfo) {
				// modified
				if err := e.updateK8sSlb(appId, &nInfo, timeout); err != nil {
					return err
				}
			}
		} else {
			// new one
			if err := e.BindK8sSlb(appId, &nInfo, timeout); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *EdasService) updateK8sSlb(appId string, slbConfig *map[string]interface{}, timeout time.Duration) *resource.RetryError {
	if e.HasOngoingTasks(appId) {
		return resource.RetryableError(Error("there is an ongoing task"))
	}
	portMappings := (*slbConfig)["port_mappings"].(*schema.Set).List()
	var requestServicePortInfos []map[string]interface{}
	if len(portMappings) > 0 {
		for _, pmInterface := range portMappings {
			pm := pmInterface.(map[string]interface{})
			if v, ok := pm["service_port"]; ok {
				spSlice := v.(*schema.Set).List()
				sp := spSlice[0].(map[string]interface{})
				info := map[string]interface{}{
					"certId":               pm["cert_id"],
					"loadBalancerProtocol": pm["loadbalancer_protocol"],
					"targetPort":           sp["target_port"],
					"port":                 sp["port"],
				}
				requestServicePortInfos = append(requestServicePortInfos, info)
			}
		}
	}
	portInfoBytes, err := json.Marshal(requestServicePortInfos)
	if err != nil {
		return resource.NonRetryableError(WrapErrorf(err, "marshal port mappings failed, value %v", requestServicePortInfos))
	}
	params := map[string]string{
		"AppId":            appId,
		"SlbName":          (*slbConfig)["name"].(string),
		"Type":             (*slbConfig)["type"].(string),
		"Scheduler":        (*slbConfig)["scheduler"].(string),
		"Specification":    (*slbConfig)["specification"].(string),
		"ClusterId":        "ClusterId",
		"ServicePortInfos": string(portInfoBytes),
	}
	if response, err := e.doPopRequest(map[string]string{
		"action":     "/pop/v5/k8s/acs/k8s_slb_binding",
		"httpMethod": "PUT",
		"id":         appId,
	}, params); err != nil {
		return resource.NonRetryableError(err)
	} else {
		return e.WaitForChangeOrderFinishedNonRetryable(appId, response["ChangeOrderId"].(string), timeout)
	}
}

func (e *EdasService) UnbindK8sSlb(appId, slbType, slbName string, timeout time.Duration) *resource.RetryError {
	if e.HasOngoingTasks(appId) {
		return resource.RetryableError(Error("there is an ongoing task"))
	}
	if response, err := e.doPopRequest(map[string]string{
		"action":     "/pop/v5/k8s/acs/k8s_slb_binding",
		"httpMethod": "DELETE",
		"id":         appId,
	}, map[string]string{
		"AppId":   appId,
		"Type":    slbType,
		"SlbName": slbName,
	}); err != nil {
		return resource.NonRetryableError(err)
	} else {
		return e.WaitForChangeOrderFinishedNonRetryable(appId, response["ChangeOrderId"].(string), timeout)
	}
}

func (e *EdasService) DescribeEdasK8sSlbAttachment(appId string) ([]map[string]interface{}, error) {
	application, err := e.DescribeEdasApplicationV2(appId)
	if err != nil {
		return nil, err
	}
	slbInfo := ""
	if v, ok := application["SlbInfo"]; ok {
		slbInfo = v.(string)
	}
	var slbConfigs []map[string]interface{}
	if !jsonEmpty(slbInfo) {
		filteredSlbInfo, err := filterSlbInfo(slbInfo)
		if err != nil {
			return nil, err
		}
		if len(*filteredSlbInfo) > 0 {
			for _, slbInfo := range *filteredSlbInfo {
				slbConfig, err := parseSlbConfig(&slbInfo)
				if err != nil {
					return nil, err
				}
				slbConfigs = append(slbConfigs, *slbConfig)
			}
		}
	}

	return slbConfigs, nil
}
