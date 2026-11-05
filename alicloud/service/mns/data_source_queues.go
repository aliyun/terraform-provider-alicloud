package mns

import (
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudMNSQueues() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceAlicloudMNSQueueRead,
		DeprecationMessage: "This data source has been deprecated since v1.188.0 and will be removed in the future. Please use 'alicloud_message_service_queues' instead.",
		Schema: map[string]*schema.Schema{
			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"queues": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"maximum_message_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"message_retention_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"visibility_timeouts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"polling_wait_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMNSQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var namePrefix string
	if v, ok := d.GetOk("name_prefix"); ok {
		namePrefix = v.(string)
	}

	var queueAttr []ali_mns.QueueAttribute
	for {
		var nextMaker string
		raw, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
			return queueManager.ListQueueDetail(nextMaker, 1000, namePrefix)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_mns_queues", "ListQueueDetail", errs.AliMnsERROR)
		}
		logs.AddDebug("ListQueueDetail", raw)
		queueDetails, _ := raw.(ali_mns.QueueDetails)
		for _, attr := range queueDetails.Attrs {
			queueAttr = append(queueAttr, attr)
		}
		nextMaker = queueDetails.NextMarker
		if nextMaker == "" {
			break
		}
	}

	return mnsQueueDescription(d, queueAttr)
}

func mnsQueueDescription(d *schema.ResourceData, queueAttr []ali_mns.QueueAttribute) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range queueAttr {
		mapping := map[string]interface{}{
			"id":                       item.QueueName,
			"name":                     item.QueueName,
			"delay_seconds":            item.DelaySeconds,
			"maximum_message_size":     item.MaxMessageSize,
			"message_retention_period": item.MessageRetentionPeriod,
			"visibility_timeouts":      item.VisibilityTimeout,
			"polling_wait_seconds":     item.PollingWaitSeconds,
		}

		names = append(names, item.QueueName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("queues", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil

}
