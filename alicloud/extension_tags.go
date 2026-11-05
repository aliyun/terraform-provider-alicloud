package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
)

// Tag moved to alicloud/tags (tags/tags.go); this alias keeps
// flat-file references and AddTagsArgs/RemoveTagsArgs compiling unchanged.
type Tag = tags.Tag

type JsonTag struct {
	TagKey   string
	TagValue string
}

type AddTagsArgs struct {
	ResourceId   string
	ResourceType ecs.TagResourceType //image, instance, snapshot or disk
	RegionId     common.Region
	Tag          []Tag
}

type RemoveTagsArgs struct {
	ResourceId   string
	ResourceType ecs.TagResourceType //image, instance, snapshot or disk
	RegionId     common.Region
	Tag          []Tag
}

func AddTags(client *ecs.Client, args *AddTagsArgs) error {
	response := ecs.AddTagsResponse{}
	err := client.Invoke("AddTags", args, &response)
	if err != nil {
		return err
	}
	return err
}

func RemoveTags(client *ecs.Client, args *RemoveTagsArgs) error {
	response := ecs.RemoveTagsResponse{}
	err := client.Invoke("RemoveTags", args, &response)
	if err != nil {
		return err
	}
	return err
}
