package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDasSwitchDasPro() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDasSwitchDasProCreate,
		Read:   resourceAlicloudDasSwitchDasProRead,
		Delete: resourceAlicloudDasSwitchDasProDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
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

func resourceAlicloudDasSwitchDasProCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasService := DasService{client}
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("DAS", "2020-01-16", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_das_switch_das_pro", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dasService.InstanceDasProStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDasSwitchDasProRead(d, meta)
}

func resourceAlicloudDasSwitchDasProRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasService := DasService{client}

	object, err := dasService.DescribeDasSwitchDasPro(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object["instanceId"])
	d.Set("sql_retention", object["sqlRetention"])
	d.Set("user_id", object["userId"])

	describeInstanceDasProObject, err := dasService.DescribeInstanceDasPro(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("status", describeInstanceDasProObject["Data"])

	return nil
}

func resourceAlicloudDasSwitchDasProDelete(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("DAS", "2020-01-16", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"-1020004"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
