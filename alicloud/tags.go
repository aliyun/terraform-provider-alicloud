package alicloud

import (
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func String(v string) *string {
	return &v
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

// setTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tags"
func setTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		return updateTags(client, []string{d.Id()}, resourceType, oraw, nraw)
	}

	return nil
}

func setCdnTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		return updateCdnTags(client, []string{d.Id()}, resourceType, oraw, nraw)
	}

	return nil
}

func setVolumeTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("volume_tags") {
		request := ecs.CreateDescribeDisksRequest()
		request.InstanceId = d.Id()
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		disks := raw.(*ecs.DescribeDisksResponse)
		if len(disks.Disks.Disk) == 0 {
			return WrapError(Error("no specified system disk"))
		}

		var ids []string
		for i := range disks.Disks.Disk {
			ids = append(ids, disks.Disks.Disk[i].DiskId)
		}

		oraw, nraw := d.GetChange("volume_tags")
		return updateTags(client, ids, resourceType, oraw, nraw)
	}

	return nil
}

func updateTags(client *connectivity.AliyunClient, ids []string, resourceType TagResourceType, oraw, nraw interface{}) error {
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

	// Set tags
	if len(remove) > 0 {
		log.Printf("[DEBUG] Removing tags: %#v from %#v", remove, ids)
		request := ecs.CreateUntagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tagsKey []string
		for _, t := range remove {
			tagsKey = append(tagsKey, t.Key)
		}
		request.TagKey = &tagsKey
		request.All = requests.NewBoolean(true)

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		log.Printf("[DEBUG] Creating tags: %s for %#v", create, ids)
		request := ecs.CreateTagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tags []ecs.TagResourcesTag
		for _, t := range create {
			tags = append(tags, ecs.TagResourcesTag{
				Key:   t.Key,
				Value: t.Value,
			})
		}
		request.Tag = &tags

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func updateCdnTags(client *connectivity.AliyunClient, ids []string, resourceType TagResourceType, oraw, nraw interface{}) error {
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

	// Set tags
	if len(remove) > 0 {
		log.Printf("[DEBUG] Removing tags: %#v from %#v", remove, ids)
		request := cdn.CreateUntagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tagsKey []string
		for _, t := range remove {
			tagsKey = append(tagsKey, t.Key)
		}
		request.TagKey = &tagsKey

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		log.Printf("[DEBUG] Creating tags: %s for %#v", create, ids)
		request := cdn.CreateTagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tags []cdn.TagResourcesTag
		for _, t := range create {
			tags = append(tags, cdn.TagResourcesTag{
				Key:   t.Key,
				Value: t.Value,
			})
		}
		request.Tag = &tags

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

// diffTags takes our tags locally and the ones remotely and returns
// the set of tags that must be created, and the set of tags that must
// be destroyed.
func diffTags(oldTags, newTags []Tag) ([]Tag, []Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []Tag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return tagsFromMap(create), remove
}

func diffGpdbTags(oldTags, newTags []gpdb.TagResourcesTag) ([]gpdb.TagResourcesTag, []gpdb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []gpdb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return gpdbTagsFromMap(create), remove
}

// tagsFromMap returns the tags for the given map of data.
func tagsFromMap(m map[string]interface{}) []Tag {
	result := make([]Tag, 0, len(m))
	for k, v := range m {
		result = append(result, Tag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func gpdbTagsFromMap(m map[string]interface{}) []gpdb.TagResourcesTag {
	result := make([]gpdb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, gpdb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func tagsToMap(tags []ecs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !ecsTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func cdnTagsToMap(tags []cdn.TagItem) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !cdnTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func essTagsToMap(tags []ess.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !essTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func otsTagsToMap(tags []ots.TagInfo) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.TagKey] = t.TagValue
	}

	return result
}

func tagsMapEqual(expectMap map[string]interface{}, compareMap map[string]string) bool {
	if len(expectMap) != len(compareMap) {
		return false
	} else {
		for key, eVal := range expectMap {
			if eStr, ok := eVal.(string); !ok {
				// type is mismatch.
				return false
			} else {
				if cStr, ok := compareMap[key]; ok {
					if eStr != cStr {
						return false
					}
				} else {
					return false
				}
			}
		}
	}
	return true
}

func tagsToString(tags []ecs.Tag) string {
	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		ecsTags := ecs.Tag{
			TagKey:   tag.TagKey,
			TagValue: tag.TagValue,
		}
		result = append(result, ecsTags.TagKey+":"+ecsTags.TagValue)
	}

	return strings.Join(result, ",")
}

// tagIgnored compares a tag against a list of strings and checks if it should be ignored or not
func ecsTagIgnored(t ecs.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

// tagIgnored compares a tag against a list of strings and checks if it should be ignored or not
func essTagIgnored(t ess.Tag) bool {
	filter := []string{"^aliyun", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func cdnTagIgnored(t cdn.TagItem) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}
