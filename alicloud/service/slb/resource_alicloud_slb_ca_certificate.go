package slb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCaCertificateCreate,
		Read:   resourceCaCertificateRead,
		Update: resourceCaCertificateUpdate,
		Delete: resourceCaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ca_certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ca_certificate_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.123.1. New field 'ca_certificate_name' instead",
				ConflictsWith: []string{"ca_certificate_name"},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceCaCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "UploadCACertificate"
	request := make(map[string]interface{})
	var err error
	request["CACertificate"] = d.Get("ca_certificate")
	if v, ok := d.GetOk("ca_certificate_name"); ok {
		request["CACertificateName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["CACertificateName"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_ca_certificate", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CACertificateId"]))

	return resourceCaCertificateUpdate(d, meta)
}
func resourceCaCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	object, err := slbService.DescribeSlbCaCertificate(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_ca_certificate slbService.DescribeSlbCaCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("ca_certificate_name", object["CACertificateName"])
	d.Set("name", object["CACertificateName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tags.TagsToMap(v["Tag"]))
	}
	return nil
}
func resourceCaCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	var err error
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := slbService.SetResourceTags(d, "certificate"); err != nil {
			return errs.WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"CACertificateId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("ca_certificate_name") {
		update = true
		request["CACertificateName"] = d.Get("ca_certificate_name")
	} else if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["CACertificateName"] = d.Get("name")
	}
	request["RegionId"] = client.RegionId
	if update {
		action := "SetCACertificateName"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceCaCertificateRead(d, meta)
}
func resourceCaCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCACertificate"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"CACertificateId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"BackendServer.configuring", "OperationBusy", "ServiceIsConfiguring", "ServiceIsStopping", "SystemBusy"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"CACertificateId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
