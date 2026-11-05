package eais

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEaisClientInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create:             ResourceEaisClientInstanceAttachmentCreate,
		Read:               ResourceEaisClientInstanceAttachmentRead,
		Update:             ResourceEaisClientInstanceAttachmentUpdate,
		Delete:             ResourceEaisClientInstanceAttachmentDelete,
		DeprecationMessage: "This resource has been deprecated since v1.278.0.",
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"client_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ei_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceEaisClientInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if helper.InArray(fmt.Sprint(d.Get("category")), []string{"eais", ""}) {
		action := "AttachEai"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("instance_id"); ok {
			request["ElasticAcceleratedInstanceId"] = v
		}
		if v, ok := d.GetOk("client_instance_id"); ok {
			request["ClientInstanceId"] = v
		}
		request["RegionId"] = client.RegionId

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_eais_client_instance_attachment", action, errs.AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", response["ElasticAcceleratedInstanceId"], response["ClientInstanceId"]))

	}

	if v, ok := d.GetOk("category"); ok && helper.InArray(fmt.Sprint(v), []string{"ei"}) {
		action := "AttachEaisEi"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("instance_id"); ok {
			request["EiInstanceId"] = v
		}
		if v, ok := d.GetOk("client_instance_id"); ok {
			request["ClientInstanceId"] = v
		}
		request["RegionId"] = client.RegionId

		if v, ok := d.GetOk("ei_instance_type"); ok {
			request["EiInstanceType"] = v
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_eais_client_instance_attachment", action, errs.AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", response["EiInstanceId"], response["ClientInstanceId"]))

	}

	return ResourceEaisClientInstanceAttachmentUpdate(d, meta)
}

func ResourceEaisClientInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eaisServiceV2 := EaisServiceV2{client}

	objectRaw, err := eaisServiceV2.DescribeEaisClientInstanceAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eais_client_instance_attachment DescribeEaisClientInstanceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["StartTime"])
	d.Set("ei_instance_type", objectRaw["InstanceType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("status", objectRaw["Status"])
	d.Set("client_instance_id", objectRaw["ClientInstanceId"])
	d.Set("instance_id", objectRaw["ElasticAcceleratedInstanceId"])

	return nil
}

func ResourceEaisClientInstanceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	enableAction := false
	if v, ok := d.GetOk("category"); ok && helper.InArray(fmt.Sprint(v), []string{"ei"}) {
		enableAction = true
	}
	if enableAction && d.HasChange("status") {
		eaisServiceV2 := EaisServiceV2{client}
		object, err := eaisServiceV2.DescribeEaisClientInstanceAttachment(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "InUse" {
				parts := strings.Split(d.Id(), ":")
				action := "StartEaisEi"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["EiInstanceId"] = parts[0]
				request["RegionId"] = client.RegionId
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			if target == "Bound" {
				parts := strings.Split(d.Id(), ":")
				action := "StopEaisEi"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["EiInstanceId"] = parts[0]
				request["RegionId"] = client.RegionId
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
		}
	}

	return ResourceEaisClientInstanceAttachmentRead(d, meta)
}

func ResourceEaisClientInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := false
	if helper.InArray(fmt.Sprint(d.Get("category")), []string{"eais", ""}) {
		enableDelete = true
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		parts := strings.Split(d.Id(), ":")
		action := "DetachEai"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["ElasticAcceleratedInstanceId"] = parts[0]
		request["RegionId"] = client.RegionId

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)

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

	}

	enableDelete = false
	if v, ok := d.GetOk("category"); ok {
		if helper.InArray(fmt.Sprint(v), []string{"ei"}) {
			enableDelete = true
		}
	}
	if enableDelete {
		client := meta.(*connectivity.AliyunClient)
		parts := strings.Split(d.Id(), ":")
		action := "DetachEaisEi"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["EiInstanceId"] = parts[0]
		request["RegionId"] = client.RegionId

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)

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

	}
	return nil
}
