package alicloud

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudOnsTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOnsTopicsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"relation": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"relation_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
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

func dataSourceAlicloudOnsTopicsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}

	request := ons.CreateOnsTopicListRequest()
	request.RegionId = client.RegionId
	request.PreventCache = onsService.GetPreventCache()
	request.InstanceId = d.Get("instance_id").(string)

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsTopicList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_topics", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ons.OnsTopicListResponse)

	var filteredTopics []ons.PublishInfoDo
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, topic := range response.Data.PublishInfoDo {
			if r != nil && !r.MatchString(topic.Topic) {
				continue
			}

			filteredTopics = append(filteredTopics, topic)
		}
	} else {
		filteredTopics = response.Data.PublishInfoDo
	}
	return onsTopicsDecriptionAttributes(d, filteredTopics, meta)
}

func onsTopicsDecriptionAttributes(d *schema.ResourceData, topicsInfo []ons.PublishInfoDo, meta interface{}) error {
	var names []string
	var s []map[string]interface{}

	for _, item := range topicsInfo {
		mapping := map[string]interface{}{
			"topic":              item.Topic,
			"owner":              item.Owner,
			"relation":           item.Relation,
			"relation_name":      item.RelationName,
			"message_type":       item.MessageType,
			"independent_naming": item.IndependentNaming,
			"create_time":        time.Unix(int64(item.CreateTime)/1000, 0).Format("2006-01-02 03:04:05"),
			"remark":             item.Remark,
		}

		names = append(names, item.Topic)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("topics", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
