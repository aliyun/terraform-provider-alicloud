package alicloud

import (
	"fmt"
	"log"
	"strings"

	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform/helper/schema"
)

func String(v string) *string {
	return &v
}

// tagsSchema returns the schema to use for tags.
//
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeMap,
		//Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
	}
}

// setTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tags"
func setTags(client *AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		// Set tags
		if len(remove) > 0 {
			log.Printf("[DEBUG] Removing tags: %#v from %s", remove, d.Id())
			args := ecs.CreateRemoveTagsRequest()
			args.ResourceType = string(resourceType)
			args.ResourceId = d.Id()
			s := reflect.ValueOf(args).Elem()

			for i, t := range remove {
				s.FieldByName(fmt.Sprintf("Tag%dKey", i+1)).Set(reflect.ValueOf(t.Key))
				s.FieldByName(fmt.Sprintf("Tag%dValue", i+1)).Set(reflect.ValueOf(t.Value))
			}
			if _, err := client.ecsconn.RemoveTags(args); err != nil {
				return fmt.Errorf("Remove tags got error: %s", err)
			}
		}

		if len(create) > 0 {
			log.Printf("[DEBUG] Creating tags: %s for %s", create, d.Id())
			args := ecs.CreateAddTagsRequest()
			args.ResourceType = string(resourceType)
			args.ResourceId = d.Id()
			s := reflect.ValueOf(args).Elem()

			for i, t := range create {
				s.FieldByName(fmt.Sprintf("Tag%dKey", i+1)).Set(reflect.ValueOf(t.Key))
				s.FieldByName(fmt.Sprintf("Tag%dValue", i+1)).Set(reflect.ValueOf(t.Value))
			}
			if _, err := client.ecsconn.AddTags(args); err != nil {
				return fmt.Errorf("Creating tags got error: %s", err)
			}
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
		result[t.TagKey] = t.TagValue
	}

	return result
}

func essTagsToMap(tags []ess.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.Key] = t.Value
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
