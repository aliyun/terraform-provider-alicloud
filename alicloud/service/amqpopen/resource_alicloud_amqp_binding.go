package amqpopen

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAmqpBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAmqpBindingCreate,
		Read:   resourceAliCloudAmqpBindingRead,
		Delete: resourceAliCloudAmqpBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_exchange": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"binding_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"EXCHANGE", "QUEUE"}, false),
			},
			"binding_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"argument": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"x-match:all", "x-match:any"}, false),
			},
		},
	}
}

func resourceAliCloudAmqpBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBinding"
	request := make(map[string]interface{})
	var err error

	request["InstanceId"] = d.Get("instance_id")
	request["VirtualHost"] = d.Get("virtual_host_name")
	request["SourceExchange"] = d.Get("source_exchange")
	request["DestinationName"] = d.Get("destination_name")
	request["BindingType"] = d.Get("binding_type")
	request["BindingKey"] = d.Get("binding_key")

	if v, ok := d.GetOk("argument"); ok {
		request["Argument"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_amqp_binding", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", request["InstanceId"], request["VirtualHost"], request["SourceExchange"], request["DestinationName"]))

	return resourceAliCloudAmqpBindingRead(d, meta)
}

func resourceAliCloudAmqpBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := NewAmqpOpenService(client)

	object, err := amqpOpenService.DescribeAmqpBinding(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_binding amqpOpenService.DescribeAmqpBinding Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("virtual_host_name", parts[1])
	d.Set("source_exchange", object["SourceExchange"])
	d.Set("destination_name", object["DestinationName"])
	d.Set("binding_type", object["BindingType"])
	d.Set("binding_key", object["BindingKey"])
	d.Set("argument", object["Argument"])

	return nil
}

func resourceAliCloudAmqpBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBinding"
	var response map[string]interface{}
	var err error

	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId":      parts[0],
		"VirtualHost":     parts[1],
		"SourceExchange":  parts[2],
		"DestinationName": parts[3],
		"BindingType":     d.Get("binding_type"),
		"BindingKey":      d.Get("binding_key"),
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("amqp-open", "2019-12-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ExchangeNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
