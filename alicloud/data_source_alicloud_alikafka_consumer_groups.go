package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudAlikafkaConsumerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaConsumerGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_id_regex": {
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
			"consumer_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"consumer_id": {
							Type:     schema.TypeString,
							Computed: true,
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

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetConsumerList(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_consumer_groups", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alikafka.GetConsumerListResponse)

	var filteredConsumerGroups []alikafka.ConsumerVO
	nameRegex, ok := d.GetOk("consumer_id_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, consumer := range response.ConsumerList.ConsumerVO {
			if r != nil && !r.MatchString(consumer.ConsumerId) {
				continue
			}

			filteredConsumerGroups = append(filteredConsumerGroups, consumer)
		}
	} else {
		filteredConsumerGroups = response.ConsumerList.ConsumerVO
	}
	return alikafkaConsumerGroupsDecriptionAttributes(d, filteredConsumerGroups, meta)
}

func alikafkaConsumerGroupsDecriptionAttributes(d *schema.ResourceData, consumerGroupsInfo []alikafka.ConsumerVO, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range consumerGroupsInfo {
		mapping := map[string]interface{}{
			"consumer_id": item.ConsumerId,
		}

		ids = append(ids, item.ConsumerId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("consumer_groups", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
