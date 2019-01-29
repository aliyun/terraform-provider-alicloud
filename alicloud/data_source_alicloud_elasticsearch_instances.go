package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudElasticsearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudElasticsearchRead,

		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"5.5.3_with_X-Pack",
					"6.3.2_with_X-Pack",
				}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_node_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_node_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_node_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_node_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudElasticsearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := elasticsearch.CreateListInstanceRequest()
	request.RegionId = client.RegionId
	request.EsVersion = d.Get("version").(string)
	request.Size = requests.NewInteger(PageSizeLarge)
	request.Page = requests.NewInteger(1)

	var instances []elasticsearch.Instance

	for {
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.ListInstance(request)
		})
		if err != nil {
			WrapErrorf(err, DataDefaultErrorMsg, "elasticsearch_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*elasticsearch.ListInstanceResponse)
		if resp == nil || len(resp.Result) < 1 {
			break
		}

		for _, item := range resp.Result {
			instances = append(instances, item)
		}

		if len(resp.Result) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.Page); err != nil {
			return WrapError(err)
		} else {
			request.Page = page
		}
	}

	var filteredInstances []elasticsearch.Instance

	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			descriptionRegex = r
		}
	}

	if len(instances) > 0 {
		for _, instance := range instances {
			if descriptionRegex != nil && !descriptionRegex.MatchString(instance.Description) {
				continue
			}

			filteredInstances = append(filteredInstances, instance)
		}
	}

	return WrapError(extractInstance(d, filteredInstances))
}

func extractInstance(d *schema.ResourceData, instances []elasticsearch.Instance) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range instances {
		mapping := map[string]interface{}{
			"id":                   item.InstanceId,
			"description":          item.Description,
			"instance_charge_type": getChargeType(item.PaymentType),
			"data_node_amount":     item.NodeAmount,
			"data_node_spec":       item.NodeSpec.Spec,
			"data_node_disk_size":  item.NodeSpec.Disk,
			"data_node_disk_type":  item.NodeSpec.DiskType,
			"status":               item.Status,
			"version":              item.EsVersion,
			"created_at":           item.CreatedAt,
			"updated_at":           item.UpdatedAt,
			"vswitch_id":           item.NetworkConfig.VswitchId,
		}

		ids = append(ids, item.InstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
