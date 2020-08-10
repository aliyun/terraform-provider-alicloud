package alicloud

import (
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	Hb_LAUNCHING            = "LAUNCHING"
	Hb_CREATING             = "CREATING"
	Hb_ACTIVATION           = "ACTIVATION"
	Hb_DELETING             = "DELETING"
	Hb_CREATE_FAILED        = "CREATE_FAILED"
	Hb_NODE_RESIZING        = "HBASE_SCALE_OUT"
	Hb_NODE_RESIZING_FAILED = "NODE_RESIZE_FAILED"
	Hb_DISK_RESIZING        = "HBASE_EXPANDING"
	Hb_DISK_RESIZE_FAILED   = "DISK_RESIZING_FAILED"
)

type HBaseService struct {
	client *connectivity.AliyunClient
}

func (s *HBaseService) setInstanceTags(d *schema.ResourceData) error {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")
	return nil
}

func (s *HBaseService) diffTags(oldTags, newTags []hbase.TagResourcesTag) ([]hbase.TagResourcesTag, []hbase.TagResourcesTag) {
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

func (s *HBaseService) tagsFromMap(m map[string]interface{}) []hbase.TagResourcesTag {
	result := make([]hbase.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, hbase.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *HBaseService) tagsToMap(tags []hbase.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *HBaseService) ignoreTag(t hbase.Tag) bool {
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

func (s *HBaseService) DescribeHBaseInstance(id string) (instance hbase.DescribeInstanceResponse, err error) {
	request := hbase.CreateDescribeInstanceRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeInstance(request)
	})
	response, _ := raw.(*hbase.DescribeInstanceResponse)
	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response == nil || response.InstanceId == "" {
		return instance, WrapErrorf(Error(GetNotFoundMessage("HBase Instance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	instance = *response
	return instance, nil
}

//pop has limit, support next.
func (s *HBaseService) DescribeIpWhitelist(id string) (instance hbase.DescribeIpWhitelistResponse, err error) {
	request := hbase.CreateDescribeIpWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeIpWhitelist(request)
	})
	response, _ := raw.(*hbase.DescribeIpWhitelistResponse)
	if err != nil {
		if IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}

func (s *HBaseService) HBaseClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHBaseInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *HBaseService) ModifyClusterDeletionProtection(clusterId string, protection bool) error {
	request := hbase.CreateModifyClusterDeletionProtectionRequest()
	request.ClusterId = clusterId
	request.Protection = requests.NewBoolean(protection)
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.ModifyClusterDeletionProtection(request)
	})
	if err != nil {
		return WrapErrorf(err, clusterId+" modifyClusterDeletionProtection failed")
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
