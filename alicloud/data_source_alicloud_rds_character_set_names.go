package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRdsCharacterSetNames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsRdsCharacterSetNamesRead,

		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(MySQL), string(SQLServer), string(PostgreSQL), string(MariaDB)}, false),
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudRdsRdsCharacterSetNamesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeCharacterSetName"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"Engine":   d.Get("engine"),
	}
	var response map[string]interface{}
	var names []string
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_character_set_names", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.CharacterSetNameItems.CharacterSetName", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Regions.Region", response)
	}
	for _, r := range resp.([]interface{}) {
		names = append(names, fmt.Sprint(r))
	}
	d.SetId(dataResourceIdHash(names))
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), names)
	}
	return nil
}
