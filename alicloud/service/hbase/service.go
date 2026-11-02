package hbase

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	Hb_LAUNCHING            = "LAUNCHING"
	Hb_CREATING             = "CREATING"
	Hb_ACTIVATION           = "ACTIVATION"
	Hb_DELETING             = "DELETING"
	Hb_CREATE_FAILED        = "CREATE_FAILED"
	Hb_NODE_RESIZING        = "HBASE_SCALE_OUT"
	HBASE_SCALE_IN          = "HBASE_SCALE_IN"
	Hb_NODE_RESIZING_FAILED = "NODE_RESIZE_FAILED"
	Hb_DISK_RESIZING        = "HBASE_EXPANDING"
	Hb_DISK_RESIZE_FAILED   = "DISK_RESIZING_FAILED"
	Hb_LEVEL_MODIFY         = "INSTANCE_LEVEL_MODIFY"
	Hb_LEVEL_MODIFY_FAILED  = "INSTANCE_LEVEL_MODIFY_FAILED"
	Hb_HBASE_COLD_EXPANDING = "HBASE_COLD_EXPANDING"
)

// PayType is the charge type used by the HBase resources and data sources.
type PayType string

const (
	PostPaid PayType = "PostPaid"
	PrePaid  PayType = "PrePaid"
	Prepaid  PayType = "Prepaid"
)

// localHostIp is the default ip white list value returned when no group is
// configured, matching the legacy root-package behaviour.
const localHostIp = "127.0.0.1"

// Service wraps the HBase API calls used by the resources and data sources
// of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the flat layout's acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := hbase.CreateUnTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.UnTagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := hbase.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.TagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *Service) diffTags(oldTags, newTags []hbase.TagResourcesTag) ([]hbase.TagResourcesTag, []hbase.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []hbase.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *Service) tagsFromMap(m map[string]interface{}) []hbase.TagResourcesTag {
	result := make([]hbase.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, hbase.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *Service) tagsToMap(tags []hbase.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *Service) ignoreTag(t hbase.Tag) bool {
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

func (s *Service) DescribeHBaseInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeInstance"

	query := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("HBase", "2019-01-01", action, query, nil, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Hbase:Instance", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
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

// pop has limit, support next.
func (s *Service) DescribeIpWhitelist(id string) (instance hbase.DescribeIpWhitelistResponse, err error) {
	request := hbase.CreateDescribeIpWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeIpWhitelist(request)
	})
	response, _ := raw.(*hbase.DescribeIpWhitelistResponse)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return instance, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}

func (s *Service) DescribeSecurityGroups(id string) (object hbase.DescribeSecurityGroupsResponse, err error) {
	request := hbase.CreateDescribeSecurityGroupsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeSecurityGroups(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.DescribeSecurityGroupsResponse)
	return *response, nil
}

func (s *Service) HBaseClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHBaseInstance(id)
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

func (s *Service) ModifyClusterDeletionProtection(clusterId string, protection bool) error {
	request := hbase.CreateModifyClusterDeletionProtectionRequest()
	request.ClusterId = clusterId
	request.Protection = requests.NewBoolean(protection)
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.ModifyClusterDeletionProtection(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, clusterId+" modifyClusterDeletionProtection failed")
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *Service) DescribeEndpoints(id string) (object hbase.DescribeEndpointsResponse, err error) {
	request := hbase.CreateDescribeEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeEndpoints(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.DescribeEndpointsResponse)
	return *response, nil
}

func (s *Service) DescribeClusterConnection(id string) (object hbase.DescribeClusterConnectionResponse, err error) {
	request := hbase.CreateDescribeClusterConnectionRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeClusterConnection(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*hbase.DescribeClusterConnectionResponse)
	return *response, nil
}

// DescribeVSwitch resolves the zone and VPC of a vswitch, replacing the
// legacy root-package VpcService.DescribeVSwitch call used during instance
// creation. It is kept behaviour-equivalent: NotFound is surfaced as a
// wrapped NotFound error.
func (s *Service) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, e := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if e != nil {
			if errs.IsExpectedErrors(e, []string{"InvalidVswitchID.NotFound"}) {
				return errs.WrapErrorf(e, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapErrorf(e, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
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
