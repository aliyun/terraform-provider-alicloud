package ga

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAlicloudGaDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaDomainCreate,
		Read:   resourceAlicloudGaDomainRead,
		Delete: resourceAlicloudGaDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"domain": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudGaDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("accelerator_id"); ok {
		request["AcceleratorIds.1"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}

	var response map[string]interface{}
	action := "CreateDomain"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ga_domain", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AcceleratorIds.1"], ":", request["Domain"]))

	return resourceAlicloudGaDomainRead(d, meta)
}

func resourceAlicloudGaDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := Service{client}

	object, err := gaService.DescribeGaDomain(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_domain gaService.DescribeGaDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("accelerator_id", parts[0])
	d.Set("domain", object["Domain"])
	d.Set("status", object["State"])

	return nil
}

func resourceAlicloudGaDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"AcceleratorIds.1": parts[0],
		"Domain":           parts[1],
	}

	action := "DeleteDomainAcceleratorRelation"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
