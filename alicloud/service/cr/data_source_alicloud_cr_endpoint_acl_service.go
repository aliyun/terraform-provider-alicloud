package cr

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudCrEndpointAclService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCrEndpointAclServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    func(v interface{}) string { return strings.ToLower(v.(string)) },
				ValidateFunc: validation.StringInSlice([]string{"internet", "Internet"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Registry", "Chart"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCrEndpointAclServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "UpdateInstanceEndpointStatus"
	request := make(map[string]interface{})
	request["Enable"] = d.Get("enable")
	request["EndpointType"] = ConvertCrEndpointType(d.Get("endpoint_type").(string))
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cr_instance_endpoint_acl_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", d.Get("endpoint_type").(string)))
	d.Set("endpoint_type", d.Get("endpoint_type").(string))
	d.Set("instance_id", request["InstanceId"])
	d.Set("enable", request["Enable"])

	if v, ok := d.GetOk("enable"); !ok || !v.(bool) {
		d.Set("status", "")
		return nil
	}

	crService := CrService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, crService.CrEndpointAclServiceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	d.Set("status", "RUNNING")
	return nil
}
