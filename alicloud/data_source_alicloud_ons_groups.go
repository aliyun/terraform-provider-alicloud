package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOnsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOnsGroupsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"group_id_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"http", "tcp"}, false),
				Default:      "tcp",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOnsGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ons.CreateOnsGroupListRequest()
	if v, ok := d.GetOk("group_type"); ok {
		request.GroupType = v.(string)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]ons.OnsGroupListTag, len(v.(map[string]interface{})))
		i := 0
		for key, value := range v.(map[string]interface{}) {
			tags[i] = ons.OnsGroupListTag{
				Key:   key,
				Value: value.(string),
			}
			i++
		}
		request.Tag = &tags
	}
	var objects []ons.SubscribeInfoDo
	var groupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		groupNameRegex = r
	}
	if v, ok := d.GetOk("group_id_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		groupNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response *ons.OnsGroupListResponse
	raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsGroupList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*ons.OnsGroupListResponse)

	for _, item := range response.Data.SubscribeInfoDo {
		if groupNameRegex != nil {
			if !groupNameRegex.MatchString(item.GroupId) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[item.GroupId]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                 object.GroupId,
			"group_name":         object.GroupId,
			"group_type":         object.GroupType,
			"independent_naming": object.IndependentNaming,
			"instance_id":        object.InstanceId,
			"owner":              object.Owner,
			"remark":             object.Remark,
		}
		ids = append(ids, object.GroupId)
		tags := make(map[string]string)
		for _, t := range object.Tags.Tag {
			tags[t.Key] = t.Value
		}
		mapping["tags"] = tags
		names = append(names, object.GroupId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
