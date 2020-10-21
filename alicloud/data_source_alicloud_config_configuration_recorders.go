package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	request := config.CreateDescribeConfigurationRecorderRequest()

	var response *config.DescribeConfigurationRecorderResponse
	raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.DescribeConfigurationRecorder(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_configuration_recorders", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*config.DescribeConfigurationRecorderResponse)

	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"id":                         fmt.Sprintf("%v", response.ConfigurationRecorder.AccountId),
		"account_id":                 fmt.Sprintf("%v", response.ConfigurationRecorder.AccountId),
		"organization_enable_status": response.ConfigurationRecorder.OrganizationEnableStatus,
		"organization_master_id":     response.ConfigurationRecorder.OrganizationMasterId,
		"resource_types":             response.ConfigurationRecorder.ResourceTypes,
		"status":                     response.ConfigurationRecorder.ConfigurationRecorderStatus,
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
