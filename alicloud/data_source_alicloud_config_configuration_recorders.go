package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudConfigConfigurationRecorders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigConfigurationRecordersRead,
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

func dataSourceAlicloudConfigConfigurationRecordersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeConfigurationRecorder"
	request := make(map[string]interface{})
	response, err := client.RpcGet("Config", "2019-01-08", action, request, nil)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_configuration_recorders", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	s := make([]map[string]interface{}, 0)

	v70, err70 := formatInt(response["ConfigurationRecorder"].(map[string]interface{})["AccountId"])
	if err70 != nil {
		return err70
	}
	v71, err71 := formatInt(response["ConfigurationRecorder"].(map[string]interface{})["AccountId"])
	if err71 != nil {
		return err71
	}

	organization_master_idVal, err73 := formatInt(response["ConfigurationRecorder"].(map[string]interface{})["OrganizationMasterId"])
	if err73 != nil {
		return err73
	}
	mapping := map[string]interface{}{
		"id":                         fmt.Sprint(v70),
		"account_id":                 fmt.Sprint(v71),
		"organization_enable_status": response["ConfigurationRecorder"].(map[string]interface{})["OrganizationEnableStatus"],
		"organization_master_id":     organization_master_idVal,
		"resource_types":             response["ConfigurationRecorder"].(map[string]interface{})["ResourceTypes"],
		"status":                     response["ConfigurationRecorder"].(map[string]interface{})["ConfigurationRecorderStatus"],
	}
	s = append(s, mapping)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("recorders", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
