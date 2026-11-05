package arms

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIntegrationExporter() *schema.Resource {
	return &schema.Resource{
		Create: ResourceIntegrationExporterCreate,
		Read:   ResourceIntegrationExporterRead,
		Update: ResourceIntegrationExporterUpdate,
		Delete: ResourceIntegrationExporterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"integration_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"param": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func ResourceIntegrationExporterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddPrometheusIntegration"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClusterId"] = d.Get("cluster_id")
	request["IntegrationType"] = d.Get("integration_type")
	request["Param"] = d.Get("param")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_integration_exporter", action, errs.AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_arms_integration_exporter")
	} else {
		instanceId := resp.(map[string]interface{})["InstanceId"]
		d.SetId(fmt.Sprintf("%v:%v:%v", request["ClusterId"], request["IntegrationType"], instanceId))
	}

	return ResourceIntegrationExporterRead(d, meta)
}

func ResourceIntegrationExporterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := Service{client}

	object, err := armsService.DescribeArmsIntegrationExporter(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cluster_id", object["ClusterId"])
	d.Set("integration_type", object["IntegrationType"])
	d.Set("param", object["Param"])
	instance_idVal, err := convert.FormatInt(object["InstanceId"])
	if err != nil {
		return err
	}
	d.Set("instance_id", instance_idVal)

	return nil
}

func ResourceIntegrationExporterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClusterId":       parts[0],
		"IntegrationType": parts[1],
		"InstanceId":      parts[2],
	}

	if d.HasChange("param") {
		update = true
	}
	request["Param"] = d.Get("param")

	if update {
		action := "UpdatePrometheusIntegration"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return ResourceIntegrationExporterRead(d, meta)
}

func ResourceIntegrationExporterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePrometheusIntegration"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClusterId":       parts[0],
		"IntegrationType": parts[1],
		"InstanceId":      parts[2],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
