package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEaisClientInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEaisClientInstanceAttachmentCreate,
		Read:   resourceAliCloudEaisClientInstanceAttachmentRead,
		Update: resourceAliCloudEaisClientInstanceAttachmentUpdate,
		Delete: resourceAliCloudEaisClientInstanceAttachmentDelete,
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

func resourceAliCloudEaisClientInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if InArray(fmt.Sprint(d.Get("category")), []string{"eais", ""}) {
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

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_eais_client_instance_attachment", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", response["ElasticAcceleratedInstanceId"], response["ClientInstanceId"]))

	}

	if v, ok := d.GetOk("category"); ok && InArray(fmt.Sprint(v), []string{"ei"}) {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_eais_client_instance_attachment", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", response["EiInstanceId"], response["ClientInstanceId"]))

	}

	return resourceAliCloudEaisClientInstanceAttachmentUpdate(d, meta)
}

func resourceAliCloudEaisClientInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eaisServiceV2 := EaisServiceV2{client}

	objectRaw, err := eaisServiceV2.DescribeEaisClientInstanceAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eais_client_instance_attachment DescribeEaisClientInstanceAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["StartTime"])
	d.Set("ei_instance_type", objectRaw["InstanceType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("status", objectRaw["Status"])
	d.Set("client_instance_id", objectRaw["ClientInstanceId"])
	d.Set("instance_id", objectRaw["ElasticAcceleratedInstanceId"])

	return nil
}

func resourceAliCloudEaisClientInstanceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	enableAction := false
	if v, ok := d.GetOk("category"); ok && InArray(fmt.Sprint(v), []string{"ei"}) {
		enableAction = true
	}
	if enableAction && d.HasChange("status") {
		eaisServiceV2 := EaisServiceV2{client}
		object, err := eaisServiceV2.DescribeEaisClientInstanceAttachment(d.Id())
		if err != nil {
			return WrapError(err)
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
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "Bound" {
				parts := strings.Split(d.Id(), ":")
				action := "StopEaisEi"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["EiInstanceId"] = parts[0]
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)
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
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	return resourceAliCloudEaisClientInstanceAttachmentRead(d, meta)
}

func resourceAliCloudEaisClientInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	enableDelete := false
	if InArray(fmt.Sprint(d.Get("category")), []string{"eais", ""}) {
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

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)

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
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

	}

	enableDelete = false
	if v, ok := d.GetOk("category"); ok {
		if InArray(fmt.Sprint(v), []string{"ei"}) {
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

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("eais", "2019-06-24", action, query, request, true)

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
			if NotFoundError(err) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

	}
	return nil
}
