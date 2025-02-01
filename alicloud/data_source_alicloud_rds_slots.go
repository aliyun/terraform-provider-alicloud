package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRdsSlots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsRdsSlotsRead,

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slot_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slot_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"temporary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slot_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wal_delay": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRdsRdsSlotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeSlots"
	request := map[string]interface{}{
		"DBInstanceId": d.Get("db_instance_id"),
		"SourceIp":     client.SourceIp,
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	var response map[string]interface{}
	var err error
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_slots", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Slots", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Slots", response)
	}
	s := make([]map[string]interface{}, 0)
	names := make([]string, 0)
	for _, r := range resp.([]interface{}) {
		item := r.(map[string]interface{})
		mapping := map[string]interface{}{
			"slot_name":   fmt.Sprint(item["SlotName"]),
			"plugin":      fmt.Sprint(item["Plugin"]),
			"database":    fmt.Sprint(item["Database"]),
			"temporary":   fmt.Sprint(item["Temporary"]),
			"slot_status": fmt.Sprint(item["SlotStatus"]),
			"wal_delay":   fmt.Sprint(item["WalDelay"]),
			"slot_type":   fmt.Sprint(item["SlotType"]),
		}
		names = append(names, fmt.Sprint(item["SlotName"]))
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(names))
	if err := d.Set("slots", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
