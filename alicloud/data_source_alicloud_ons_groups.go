package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudOnsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOnsGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
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
	onsService := OnsService{client}

	request := ons.CreateOnsGroupListRequest()
	request.PreventCache = onsService.GetPreventCache()
	request.InstanceId = d.Get("instance_id").(string)

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsGroupList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ons.OnsGroupListResponse)

	var filteredGroups []ons.SubscribeInfoDo
	nameRegex, ok := d.GetOk("group_id_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, group := range response.Data.SubscribeInfoDo {
			if r != nil && !r.MatchString(group.GroupId) {
				continue
			}

			filteredGroups = append(filteredGroups, group)
		}
	} else {
		filteredGroups = response.Data.SubscribeInfoDo
	}

	return onsGroupsDecriptionAttributes(d, filteredGroups, meta)
}

func onsGroupsDecriptionAttributes(d *schema.ResourceData, topicsInfo []ons.SubscribeInfoDo, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range topicsInfo {
		mapping := map[string]interface{}{
			"id":                 item.GroupId,
			"owner":              item.Owner,
			"independent_naming": item.IndependentNaming,
			"remark":             item.Remark,
		}

		ids = append(ids, item.GroupId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
