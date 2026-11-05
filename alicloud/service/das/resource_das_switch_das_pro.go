package das

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDasSwitchDasPro() *schema.Resource {
	return &schema.Resource{
		Create: resourceDasSwitchDasProCreate,
		Read:   resourceDasSwitchDasProRead,
		Delete: resourceDasSwitchDasProDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		DeprecationMessage: "This resource has been deprecated since v1.249.0 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sql_retention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{30, 180, 365, 1095, 1825}),
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDasSwitchDasProCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasService := NewService(client)
	var response map[string]interface{}
	action := "EnableDasPro"
	request := make(map[string]interface{})
	var err error

	request["InstanceId"] = d.Get("instance_id")

	if v, ok := d.GetOk("sql_retention"); ok {
		request["SqlRetention"] = v
	}
	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("DAS", "2020-01-16", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_das_switch_das_pro", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dasService.InstanceDasProStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDasSwitchDasProRead(d, meta)
}

func resourceDasSwitchDasProRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasService := NewService(client)

	object, err := dasService.DescribeDasSwitchDasPro(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("instance_id", object["instanceId"])
	d.Set("sql_retention", object["sqlRetention"])
	d.Set("user_id", object["userId"])

	describeInstanceDasProObject, err := dasService.DescribeInstanceDasPro(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("status", describeInstanceDasProObject["Data"])

	return nil
}

func resourceDasSwitchDasProDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DisableDasPro"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("DAS", "2020-01-16", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"-1020004"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
