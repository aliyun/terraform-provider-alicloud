package cassandra

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

// NewService constructs the Cassandra Service from a connectivity client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeCassandraCluster(id string) (object cassandra.Cluster, err error) {
	cluster := cassandra.Cluster{}
	request := cassandra.CreateDescribeClusterRequest()
	request.RegionId = s.client.RegionId

	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeCluster(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return cluster, errs.WrapErrorf(errs.NotFoundErr("Cassandra", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return cluster, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeClusterResponse)
	return response.Cluster, nil
}

func (s *Service) CassandraClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCassandraCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *Service) DescribeCassandraDataCenter(id string) (object cassandra.DescribeDataCenterResponse, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := cassandra.CreateDescribeDataCenterRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = parts[1]
	request.DataCenterId = parts[0]

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeDataCenter(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("CassandraDataCenter", id), errs.NotFoundMsg, errs.ProviderERROR)
			return
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeDataCenterResponse)
	return *response, nil
}

func (s *Service) CassandraDataCenterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCassandraDataCenter(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *Service) DescribeIpWhitelist(id string) (object cassandra.DescribeIpWhitelistResponse, err error) {
	request := cassandra.CreateDescribeIpWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeIpWhitelist(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeIpWhitelistResponse)
	return *response, nil
}

func (s *Service) DescribeClusterDataCenter(id string) (object cassandra.DataCenter, err error) {
	request := cassandra.CreateDescribeDataCentersRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeDataCenters(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeDataCentersResponse)
	for _, dc := range response.DataCenters.DataCenter {
		if dc.CommodityInstance == id {
			return dc, nil
		}
	}
	return
}

func (s *Service) DescribeSecurityGroups(id string) (object cassandra.DescribeSecurityGroupsResponse, err error) {
	request := cassandra.CreateDescribeSecurityGroupsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeSecurityGroups(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeSecurityGroupsResponse)
	return *response, nil
}

func (s *Service) DescribeCassandraEndpoints(id string) (object cassandra.DescribeContactPointsResponse, err error) {
	parts := strings.Split(id, ":")
	request := cassandra.CreateDescribeContactPointsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = parts[len(parts)-1]

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeContactPoints(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeContactPointsResponse)
	return *response, nil
}

func (s *Service) setInstanceTags(d *schema.ResourceData) error {
	if !d.HasChange("tags") {
		return nil
	}
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := cassandra.CreateUnTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.UnTagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := cassandra.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
			return cassandraClient.TagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *Service) diffTags(oldTags, newTags []cassandra.TagResourcesTag) ([]cassandra.TagResourcesTag, []cassandra.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []cassandra.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *Service) tagsFromMap(m map[string]interface{}) []cassandra.TagResourcesTag {
	result := make([]cassandra.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, cassandra.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *Service) tagsToMap(tags []cassandra.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *Service) ignoreTag(t cassandra.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *Service) DescribeAccounts(id string) (object cassandra.DescribeAccountsResponse, err error) {
	request := cassandra.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithCassandraClient(func(cassandraClient *cassandra.Client) (interface{}, error) {
		return cassandraClient.DescribeAccounts(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cassandra.DescribeAccountsResponse)
	return *response, nil
}

func (s *Service) DescribeCassandraBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBackupPlans"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cassandra", "2019-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cassandra:BackupPlan", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.BackupPlans.BackupPlan", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cassandra", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if v.(map[string]interface{})["DataCenterId"].(string) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cassandra", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

// describeVSwitch mirrors alicloud.VpcService.DescribeVSwitch for the cassandra
// package without importing the root alicloud package (cross-product dependency
// kept local; the root VpcService remains for not-yet-migrated callers).
func (s *Service) describeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
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
		if response.VSwitchId != id {
			return errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}
