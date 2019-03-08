package alicloud

import (
	"fmt"
	"log"
	"strings"

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

// tagsSchema returns the schema to use for tags.
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

func setVolumeTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("volume_tags") {
		resp, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			req := ecs.CreateDescribeDisksRequest()
			req.InstanceId = d.Id()
			return ecsClient.DescribeDisks(req)
		})
		if err != nil {
			return fmt.Errorf("describe disk for %s failed, %#v", d.Id(), err)
		}

		disks := resp.(*ecs.DescribeDisksResponse)
		if len(disks.Disks.Disk) == 0 {
			return fmt.Errorf("no specified system disk")
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
		args := ecs.CreateUntagResourcesRequest()
		args.ResourceType = string(resourceType)
		args.ResourceId = &ids

		var tagsKey []string
		for _, t := range remove {
			tagsKey = append(tagsKey, t.Key)
		}
		args.TagKey = &tagsKey
		args.All = requests.NewBoolean(true)

		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.UntagResources(args)
		})
		if err != nil {
			return fmt.Errorf("Remove tags got error: %s", err)
		}
	}

	if len(create) > 0 {
		log.Printf("[DEBUG] Creating tags: %s for %#v", create, ids)
		args := ecs.CreateTagResourcesRequest()
		args.ResourceType = string(resourceType)
		args.ResourceId = &ids

		var tags []ecs.TagResourcesTag
		for _, t := range create {
			tags = append(tags, ecs.TagResourcesTag{
				Key:   t.Key,
				Value: t.Value,
			})
		}
		args.Tag = &tags

		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.TagResources(args)
		})
		if err != nil {
			return fmt.Errorf("Creating tags got error: %s", err)
		}
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

func tagsToMap(tags []ecs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !ecsTagIgnored(t) {
			result[t.TagKey] = t.TagValue
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
