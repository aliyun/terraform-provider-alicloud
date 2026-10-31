package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceConfigurationRecorders() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceConfigurationRecordersRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recorders": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_enable_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_master_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceConfigurationRecordersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeConfigurationRecorder"
	request := make(map[string]interface{})
	response, err := client.RpcGet("Config", "2019-01-08", action, request, nil)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_config_configuration_recorders", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	s := make([]map[string]interface{}, 0)

	accountId, err := convert.FormatInt(response["ConfigurationRecorder"].(map[string]interface{})["AccountId"])
	if err != nil {
		return err
	}
	accountIdValue, err := convert.FormatInt(response["ConfigurationRecorder"].(map[string]interface{})["AccountId"])
	if err != nil {
		return err
	}

	organization_master_idVal, err := convert.FormatInt(response["ConfigurationRecorder"].(map[string]interface{})["OrganizationMasterId"])
	if err != nil {
		return err
	}
	mapping := map[string]interface{}{
		"id":                         fmt.Sprint(accountId),
		"account_id":                 fmt.Sprint(accountIdValue),
		"organization_enable_status": response["ConfigurationRecorder"].(map[string]interface{})["OrganizationEnableStatus"],
		"organization_master_id":     organization_master_idVal,
		"resource_types":             response["ConfigurationRecorder"].(map[string]interface{})["ResourceTypes"],
		"status":                     response["ConfigurationRecorder"].(map[string]interface{})["ConfigurationRecorderStatus"],
	}
	s = append(s, mapping)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("recorders", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
