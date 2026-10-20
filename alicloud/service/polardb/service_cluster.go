package polardb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeVSwitch resolves the VPC and zone of a vSwitch used to fill in
// cluster creation defaults.
func (s *Service) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response == nil || response.VSwitchId != id {
			return errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

func (s *Service) CheckKMSAuthorized(id string, tdeRegion string) (map[string]interface{}, error) {
	action := "CheckKMSAuthorized"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	if tdeRegion != "" {
		request["tdeRegion"] = tdeRegion
	}
	var response map[string]interface{}
	var err error
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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
	return response, nil
}

func (s *Service) CreateClusterParamsModifyParameters(d *schema.ResourceData) error {
	request := polardb.CreateModifyDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = d.Id()
	config := make(map[string]interface{})
	allConfig := make(map[string]string)
	changeParams := []string{"loose_polar_log_bin", "default_time_zone", "loose_xengine", "loose_xengine_use_memory_pct"}
	for _, i := range changeParams {
		if v, ok := d.GetOk(i); ok {
			if d.HasChange(i) {
				if i == "loose_xengine" {
					if v == "ON" {
						v = "1"
					} else if v == "OFF" {
						v = "0"
					}
				}
				config[i] = v
				allConfig[i] = fmt.Sprint(v)
			}
		}
	}
	cfg, _ := json.Marshal(config)
	request.Parameters = string(cfg)
	// wait instance status is Normal before modifying
	if err := s.WaitForCluster(d.Id(), Running, DefaultLongTimeout); err != nil {
		return errs.WrapError(err)
	}
	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.ModifyDBClusterParameters(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	// wait instance parameter expect after modifying
	if err := s.WaitForParameter(d.Id(), 1200, allConfig); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *Service) DescribeDBAuditLogCollectorStatus(id string) (collectorStatus string, err error) {
	request := polardb.CreateDescribeDBClusterAuditLogCollectorRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	var response *polardb.DescribeDBClusterAuditLogCollectorResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusterAuditLogCollector(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response, _ = raw.(*polardb.DescribeDBClusterAuditLogCollectorResponse)
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return "", errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return collectorStatus, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return response.CollectorStatus, nil
}

func (s *Service) DescribeDBClusterAccessWhitelist(id string) (instance *polardb.DescribeDBClusterAccessWhitelistResponse, err error) {
	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	var response *polardb.DescribeDBClusterAccessWhitelistResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterAccessWhitelist(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)
		return nil
	})

	if err != nil {
		return nil, errs.WrapError(err)
	}
	return response, nil
}

func (s *Service) DescribeDBClusterAttribute(id string) (object map[string]interface{}, err error) {
	action := "DescribeDBClusterAttribute"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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

func (s *Service) DescribeDBClusterAvailableVersion(id string) (instance *polardb.DescribeDBClusterVersionResponse, err error) {
	request := polardb.CreateDescribeDBClusterVersionRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	request.DescribeType = "AVAILABLE_VERSION"

	var response *polardb.DescribeDBClusterVersionResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterVersion(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterVersionResponse)
		return nil
	})

	if err != nil {
		return nil, errs.WrapError(err)
	}
	return response, nil
}

func (s *Service) DescribeDBClusterServerlessConfig(id string) (object map[string]interface{}, err error) {
	action := "DescribeDBClusterServerlessConf"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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

func (s *Service) DescribeDBClusterStandbyAz(d *schema.ResourceData, timeout int) (zoneId string, err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	standbyAz := ""
	for {
		clusterAttribute, err := s.DescribeClusterAttribute(d.Id())
		if err != nil {
			return "", errs.WrapError(err)
		}

		exist := true
		if clusterAttribute.HotStandbyCluster != "StandbyClusterOFF" {
			exist = false
			zoneIds := strings.Split(clusterAttribute.ZoneIds, ",")
			if len(zoneIds) >= 2 {
				zoneId = zoneIds[0]
				standbyAz = zoneIds[1]
				// standbyAz info for new cluster may be incorrectly set to be equal to zoneId at first, need retry for real value
				if standbyAz != "" && (standbyAz != zoneId || !d.IsNewResource()) {
					exist = true
				}
			}
		}

		if exist {
			return standbyAz, nil
		}
		if time.Now().After(deadline) {
			// adapt for new cluster with exactly the same zoneId and standbyAz
			if standbyAz != "" {
				return standbyAz, nil
			}
			return standbyAz, errs.WrapErrorf(err, errs.RequiredWhenMsg, "standby_az", "hot_standby_cluster", clusterAttribute.HotStandbyCluster)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *Service) DescribeDBClusterTDE(id string) (map[string]interface{}, error) {
	action := "DescribeDBClusterTDE"
	request := map[string]interface{}{
		"DBClusterId": id,
	}
	client := s.client
	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	return response, nil
}

func (s *Service) DescribeDBClusterVersion(id string) (object map[string]interface{}, err error) {
	action := "DescribeDBClusterVersion"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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

func (s *Service) DescribeDBSecurityGroups(clusterId string) ([]string, error) {

	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAccessWhitelist(request)
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)

	groups := make([]string, 0)
	dbClusterSecurityGroups := resp.DBClusterSecurityGroups.DBClusterSecurityGroup
	for _, group := range dbClusterSecurityGroups {

		groups = append(groups, group.SecurityGroupId)
	}
	return groups, nil
}

func (s *Service) DescribeParameters(id string) (ds *polardb.DescribeDBClusterParametersResponse, err error) {
	request := polardb.CreateDescribeDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.RetryContext(context.Background(), 5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterParameters(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR))
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*polardb.DescribeDBClusterParametersResponse)
	return response, err
}

func (s *Service) DescribeAutoRenewAttribute(id string) (instance *polardb.AutoRenewAttribute, err error) {
	request := polardb.CreateDescribeAutoRenewAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterIds = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeAutoRenewAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeAutoRenewAttributeResponse)
	if len(response.Items.AutoRenewAttribute) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items.AutoRenewAttribute[0], nil
}

func (s *Service) DescribeCluster(id string) (instance *polardb.DBCluster, err error) {
	request := polardb.CreateDescribeDBClustersRequest()
	request.RegionId = s.client.RegionId
	dbClusterIds := []string{}
	dbClusterIds = append(dbClusterIds, id)
	request.DBClusterIds = id
	var response *polardb.DescribeDBClustersResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusters(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClustersResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	if len(response.Items.DBCluster) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *Service) DescribeClusterAttribute(id string) (instance *polardb.DescribeDBClusterAttributeResponse, err error) {
	request := polardb.CreateDescribeDBClusterAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	var response *polardb.DescribeDBClusterAttributeResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterAttribute(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterAttributeResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidDBCluster.NotFound"}) {
			return instance, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *Service) DescribeGlobalSecurityGroupIds(id string) (globalSecurityGroupIds []string, err error) {
	action := "DescribeGlobalSecurityIPGroupRelation"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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
		return globalSecurityGroupIds, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.GlobalSecurityIPGroupRel", response)
	if err != nil {
		return globalSecurityGroupIds, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.GlobalSecurityIPGroupRel", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return globalSecurityGroupIds, nil
	}

	for _, v := range resp.([]interface{}) {
		globalSecurityGroupIds = append(globalSecurityGroupIds, fmt.Sprint(v.(map[string]interface{})["GlobalSecurityGroupId"]))
	}

	return globalSecurityGroupIds, nil
}

func (s *Service) DescribeInstanceNetInfo(id string) ([]polardb.DBEndpoint, error) {

	request := polardb.CreateDescribeDBClusterEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	var response *polardb.DescribeDBClusterEndpointsResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterEndpoints(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterEndpointsResponse)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(response.Items) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBInstanceNetInfo", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return response.Items, nil
}

func (s *Service) DescribeTags(resourceId string, resourceType tags.TagResourceType) (tags []polardb.TagResource, err error) {
	request := polardb.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

// WaitForCluster waits for cluster to given status

func (s *Service) ModifyDBClusterAccessWhitelist(d *schema.ResourceData) error {
	if _, ok := d.GetOk("db_cluster_ip_array"); ok {
		removed, added := d.GetChange("db_cluster_ip_array")

		targetIpArrayNames := make(map[string]bool)
		for _, e := range added.(*schema.Set).List() {
			pack := e.(map[string]interface{})
			ipArrayName := pack["db_cluster_ip_array_name"].(string)
			targetIpArrayNames[ipArrayName] = true
		}

		for _, e := range removed.(*schema.Set).List() {
			pack := e.(map[string]interface{})
			if fmt.Sprint(pack["db_cluster_ip_array_name"]) == "default" || targetIpArrayNames[pack["db_cluster_ip_array_name"].(string)] {
				continue
			}
			//ips expand string list
			ipList := convert.ExpandStringList(pack["security_ips"].(*schema.Set).List())
			ipstr := strings.Join(ipList[:], helper.CommaSeparated)
			request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
			request.RegionId = s.client.RegionId
			request.DBClusterId = d.Id()
			request.SecurityIps = ipstr
			request.DBClusterIPArrayName = pack["db_cluster_ip_array_name"].(string)
			request.ModifyMode = "Delete"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(5*time.Minute, func() *retry.RetryError {
				raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
					return polarDBClient.ModifyDBClusterAccessWhitelist(request)
				})
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			if err := s.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
				return errs.WrapError(err)
			}
		}
		for _, e := range added.(*schema.Set).List() {
			pack := e.(map[string]interface{})
			//ips expand string list
			ipList := convert.ExpandStringList(pack["security_ips"].(*schema.Set).List())
			ipstr := strings.Join(ipList[:], helper.CommaSeparated)
			// default disable connect from outside
			if ipstr == "" {
				ipstr = LOCAL_HOST_IP
			}
			request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
			request.RegionId = s.client.RegionId
			request.DBClusterId = d.Id()
			request.SecurityIps = ipstr
			request.DBClusterIPArrayName = pack["db_cluster_ip_array_name"].(string)
			request.ModifyMode = pack["modify_mode"].(string)
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(5*time.Minute, func() *retry.RetryError {
				raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
					return polarDBClient.ModifyDBClusterAccessWhitelist(request)
				})
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			if err := s.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	return nil
}

func (s *Service) ModifyDBSecurityIps(clusterId, ips string) error {

	request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.SecurityIps = ips

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.ModifyDBClusterAccessWhitelist(request)
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

func (s *Service) ModifyParameters(d *schema.ResourceData) error {
	request := polardb.CreateModifyDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = d.Id()
	config := make(map[string]string)
	allConfig := make(map[string]string)
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request.Parameters = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForCluster(d.Id(), Running, DefaultLongTimeout); err != nil {
			return errs.WrapError(err)
		}
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterParameters(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance parameter expect after modifying
		for _, i := range ns.List() {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			allConfig[key] = value
		}
		if err := s.WaitForParameter(d.Id(), DefaultLongTimeout, allConfig); err != nil {
			return errs.WrapError(err)
		}
	}
	return nil
}

func (s *Service) ClusterCategoryRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Category == failState {
				return object, object.Category, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Category))
			}
		}
		return object, object.Category, nil
	}
}

func (s *Service) ClusterProxyStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.ProxyStatus == failState {
				return object, object.ProxyStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.ProxyStatus))
			}
		}
		return object, object.ProxyStatus, nil
	}
}

func (s *Service) ClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeClusterAttribute(id)
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

func (s *Service) ClusterTDEStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		object, err := s.DescribeDBClusterTDE(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["TDEStatus"].(string) == failState {
				return object, object["TDEStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["TDEStatus"].(string)))
			}
		}
		return object, object["TDEStatus"].(string), nil
	}
}

func (s *Service) RefreshParameters(d *schema.ResourceData) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk("parameters")
	object, err := s.DescribeParameters(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	var parameters = make(map[string]interface{})
	for _, i := range object.RunningParameters.Parameter {
		// 创建集群传入参数模板参数
		changeParams := []string{"loose_polar_log_bin", "lower_case_table_names", "default_time_zone", "loose_xengine", "loose_xengine_use_memory_pct"}
		if helper.InArray(i.ParameterName, changeParams) {
			if i.ParameterName == "lower_case_table_names" {
				if _parameterValue, err := strconv.Atoi(i.ParameterValue); err == nil {
					d.Set("lower_case_table_names", _parameterValue)
				}
			} else if d.Get("db_type").(string) == "MySQL" && d.Get("db_version") == "5.6" && i.ParameterName == "loose_polar_log_bin" {
				// mysql 5.6 loose_polar_log_bin values: ON_WITH_GTID、OFF
				if i.ParameterValue == "ON_WITH_GTID" {
					d.Set("loose_polar_log_bin", "ON")
				}
			} else if i.ParameterName == "loose_xengine" {
				if i.ParameterValue == "1" {
					d.Set("loose_xengine", "ON")
				} else {
					d.Set("loose_xengine", "OFF")
				}
			} else if i.ParameterName == "loose_xengine_use_memory_pct" {
				if parameterValue, err := strconv.Atoi(i.ParameterValue); err == nil {
					d.Set("loose_xengine_use_memory_pct", parameterValue)
				}
			} else {
				switch i.ParameterName {
				case "default_time_zone":
					d.Set("default_time_zone", i.ParameterValue)
				case "loose_polar_log_bin":
					d.Set("loose_polar_log_bin", i.ParameterValue)
				}
			}
		}
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}
	if !ok {
		d.Set("parameters", param)
		return nil
	}
	if err := d.Set("parameters", param); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *Service) WaitForCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeClusterAttribute(id)
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

func (s *Service) WaitForNodeClass(id string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		clusters, err := s.DescribeCluster(id)
		if err != nil {
			return errs.WrapError(err)
		}
		clusterAttribute, err := s.DescribeClusterAttribute(id)
		if err != nil {
			return errs.WrapError(err)
		}
		if len(clusters.DBNodes.DBNode) == len(clusterAttribute.DBNodes) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *Service) WaitForParameter(clusterId string, timeout int, expects map[string]string) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	// get engine and engineVersion
	attrbuteInfo, err := s.DescribeClusterAttribute(clusterId)
	if err != nil {
		return errs.WrapError(err)
	}
	dbType := fmt.Sprint(attrbuteInfo.DBType)
	db_version := fmt.Sprint(attrbuteInfo.DBVersion)

	for {
		object, err := s.DescribeParameters(clusterId)
		if err != nil {
			return errs.WrapError(err)
		}

		var actuals = make(map[string]string)
		for _, i := range object.RunningParameters.Parameter {
			// mysql 5.6 loose_polar_log_bin values: ON_WITH_GTID、OFF
			if dbType == "MySQL" && db_version == "5.6" && i.ParameterName == "loose_polar_log_bin" {
				if i.ParameterValue == "ON_WITH_GTID" {
					i.ParameterValue = "ON"
				}
			}
			actuals[i.ParameterName] = i.ParameterValue
		}

		match := true

		got_value := ""
		expected_value := ""

		for name, expect := range expects {
			if actual, ok := actuals[name]; ok {
				if expect != actual {
					match = false
					got_value = actual
					expected_value = expect
					break
				}
			} else {
				match = false
			}
		}

		if match {
			break
		}

		time.Sleep(helper.DefaultIntervalShort * time.Second)

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, clusterId, helper.GetFunc(1), timeout, got_value, expected_value, errs.ProviderERROR)
		}
	}
	return nil
}

func (s *Service) WaitForPayType(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		clusters, err := s.DescribeCluster(id)
		if err != nil {
			return errs.WrapError(err)
		}
		if clusters.PayType == status {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, clusters, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *Service) diffTags(oldTags, newTags []polardb.TagResourcesTag) ([]polardb.TagResourcesTag, []polardb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []polardb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *Service) ignoreTag(t polardb.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *Service) setClusterTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				if !tags.IgnoredTags(v.Key, v.Value) {
					tagKey = append(tagKey, v.Key)
				}
			}
			request := polardb.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = "cluster"
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := polardb.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = "cluster"
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
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

func (s *Service) tagsFromMap(m map[string]interface{}) []polardb.TagResourcesTag {
	result := make([]polardb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, polardb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *Service) tagsToMap(tags []polardb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

// DescribeAvailableInstance reads the billing record of a cluster so the
// subscription release path can reuse the commodity ProductType reported by
// the account instead of hard coding one per station.
func (s *Service) DescribeAvailableInstance(id string) (map[string]interface{}, error) {
	action := "QueryAvailableInstances"
	// The API is keyed by InstanceIDs and is resolved by the account's own
	// station, so Region is deliberately not sent.
	request := map[string]interface{}{
		"ProductCode": "polardb",
		"InstanceIDs": id,
	}
	var response map[string]interface{}
	var endpoint string
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		response, err = s.client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !s.client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList", response)
	}
	list, _ := v.([]interface{})
	if len(list) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("PolarDBInstance", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	instance, _ := list[0].(map[string]interface{})
	if instance == nil || fmt.Sprint(instance["InstanceID"]) != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("PolarDBInstance", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return instance, nil
}

// RefundPrePaidCluster releases a subscription cluster through the billing
// RefundInstance API and returns the refund order id. A cluster that is
// already refunded reports ResourceNotExists, which the caller treats as an
// idempotent success.
func (s *Service) RefundPrePaidCluster(id string, timeout time.Duration) (int64, error) {
	instance, err := s.DescribeAvailableInstance(id)
	if err != nil {
		return 0, err
	}
	productType := fmt.Sprint(instance["ProductType"])
	if productType == "<nil>" {
		productType = ""
	}

	action := "RefundInstance"
	request := map[string]interface{}{
		"ProductCode":        "polardb",
		"ProductType":        productType,
		"InstanceId":         id,
		"ImmediatelyRelease": "1",
		"ClientToken":        helper.BuildClientToken(action),
	}
	var response map[string]interface{}
	var endpoint string
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		var err error
		response, err = s.client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !s.client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return 0, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	// A business rejection can also arrive with HTTP 200 and Success=false;
	// keep the response Code so the caller can classify it like an error code.
	// A response that does not confirm success is a rejection, not a refund.
	success, ok := response["Success"].(bool)
	if !ok || !success {
		code := fmt.Sprint(response["Code"])
		if code == "<nil>" {
			code = ""
		}
		message := fmt.Sprint(response["Message"])
		if message == "<nil>" || message == "" {
			message = "the refund response did not confirm success"
		}
		return 0, errs.WrapErrorf(errs.NewProviderError(code, message), errs.DefaultErrorMsg, id, action, errs.ProviderERROR)
	}

	var orderId int64
	if v, err := jsonpath.Get("$.Data.OrderId", response); err == nil {
		orderId, _ = strconv.ParseInt(fmt.Sprint(v), 10, 64)
	}
	log.Printf("[INFO] Refunded PolarDB cluster %s, refund order id: %d", id, orderId)
	return orderId, nil
}
