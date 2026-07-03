package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudDasSqlLogConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDasSqlLogConfigsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"request_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hot_retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cold_retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_log_visible_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudDasSqlLogConfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasServiceV2 := DasServiceV2{client}

	instanceId := d.Get("instance_id").(string)

	objectRaw, err := dasServiceV2.DescribeDasSqlLogConfig(instanceId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId(dataResourceIdHash([]string{instanceId}))
			d.Set("ids", []string{})
			d.Set("configs", []map[string]interface{}{})
			return nil
		}
		return WrapError(err)
	}

	dataRawObj, _ := jsonpath.Get("$.Data", objectRaw)
	dataRaw := make(map[string]interface{})
	if dataRawObj != nil {
		dataRaw = dataRawObj.(map[string]interface{})
	}

	toBool := func(v interface{}) bool {
		b, _ := v.(bool)
		return b
	}
	toStr := func(v interface{}) string {
		if v == nil {
			return ""
		}
		s, _ := v.(string)
		return s
	}

	mapping := map[string]interface{}{
		"id":                   instanceId,
		"instance_id":          instanceId,
		"enable":               toBool(dataRaw["SqlLogEnable"]),
		"request_enable":       toBool(dataRaw["RequestEnable"]),
		"retention":            formatInt(dataRaw["Retention"]),
		"hot_retention":        formatInt(dataRaw["HotRetention"]),
		"cold_retention":       formatInt(dataRaw["ColdRetention"]),
		"version":              toStr(dataRaw["Version"]),
		"log_filter":           toStr(dataRaw["LogFilter"]),
		"sql_log_visible_time": formatInt(dataRaw["SqlLogVisibleTime"]),
	}

	d.SetId(dataResourceIdHash([]string{instanceId}))
	if err := d.Set("instance_id", instanceId); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", []string{instanceId}); err != nil {
		return WrapError(err)
	}
	if err := d.Set("configs", []map[string]interface{}{mapping}); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), []map[string]interface{}{mapping}); err != nil {
			return WrapError(fmt.Errorf("writing output file: %w", err))
		}
	}

	return nil
}
