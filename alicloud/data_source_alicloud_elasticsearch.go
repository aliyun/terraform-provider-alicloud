package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudElasticsearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudElasticsearchRead,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"es_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"5.5.3_with_X-Pack",
					"6.3.2_with_X-Pack",
				}),
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
						"data_node_disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_node_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"es_version": {
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
	request.Description = d.Get("description").(string)
	request.EsVersion = d.Get("es_version").(string)
	request.Size = requests.NewInteger(PageSizeLarge)

	var instances []elasticsearch.Instance

	for {
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.ListInstance(request)
		})
		if err != nil {
			return err
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
			return err
		} else {
			request.Page = page
		}
	}

	return extractInstance(d, instances)
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
			"data_node_disk":       item.NodeSpec.Disk,
			"data_node_disk_type":  item.NodeSpec.DiskType,
			"status":               item.Status,
			"es_version":           item.EsVersion,
			"created_at":           item.CreatedAt,
			"updated_at":           item.UpdatedAt,
			"vswitch_id":           item.NetworkConfig.VswitchId,
		}

		ids = append(ids, item.InstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
