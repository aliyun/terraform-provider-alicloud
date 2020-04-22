package alicloud

import (
	"strconv"
	"strings"

	"github.com/emicklei/go-restful/log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type EdasService struct {
	client *connectivity.AliyunClient
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
		return Error("sync resource failed for " + rsp.Message)
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

	log.Printf("rsp: %+v", rsp.Instances.Instance)
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
		return "", Error("QueryApplicationStatus failed for " + response.Message)
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
		return "", Error("QueryApplicationStatus failed for " + response.Message)
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
		return cluster, Error("create cluster failed for " + response.Message)
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
		return group, Error("create cluster failed for " + response.Message)
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

func (e *EdasService) DescribeEdasSlbBind(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, WrapError(err)
	}

	return o, nil
}
