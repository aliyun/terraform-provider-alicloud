package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlikafkaConsumerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaConsumerGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consumer_id_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"consumer_ids": {
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
						"consumer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaConsumerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	request := alikafka.CreateGetConsumerListRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.RegionId = client.RegionId

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request.CurrentPage = requests.NewInteger(v.(int))
	} else {
		request.CurrentPage = requests.NewInteger(1)
	}
	pageSize := PageSizeLarge
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		pageSize = v.(int)
	}
	request.PageSize = requests.NewInteger(pageSize)

	var r *regexp.Regexp
	if v, ok := d.GetOk("consumer_id_regex"); ok && v.(string) != "" {
		var compileErr error
		r, compileErr = regexp.Compile(v.(string))
		if compileErr != nil {
			return WrapError(compileErr)
		}
	}

	var filteredConsumerGroups []alikafka.ConsumerVO
	for {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetConsumerList(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_consumer_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alikafka.GetConsumerListResponse)

		for _, consumer := range response.ConsumerList.ConsumerVO {
			if r != nil && !r.MatchString(consumer.ConsumerId) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(consumer.InstanceId, ":", consumer.ConsumerId)]; !ok {
					continue
				}
			}

			filteredConsumerGroups = append(filteredConsumerGroups, consumer)
		}

		if isPagingRequest(d) {
			break
		}
		if len(response.ConsumerList.ConsumerVO) < pageSize {
			break
		}
		page, err := getNextpageNumber(request.CurrentPage)
		if err != nil {
			return WrapError(err)
		}
		request.CurrentPage = page
	}

	return alikafkaConsumerGroupsDecriptionAttributes(d, filteredConsumerGroups, meta)
}

func alikafkaConsumerGroupsDecriptionAttributes(d *schema.ResourceData, consumerGroupsInfo []alikafka.ConsumerVO, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range consumerGroupsInfo {
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(item.InstanceId, ":", item.ConsumerId),
			"instance_id": item.InstanceId,
			"consumer_id": item.ConsumerId,
			"remark":      item.Remark,
			"tags":        alikafkaService.tagVOTagsToMap(item.Tags.TagVO),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, item.ConsumerId)
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

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), ids)
	}
	return nil

}
