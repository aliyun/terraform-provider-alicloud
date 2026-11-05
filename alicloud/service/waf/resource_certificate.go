package waf

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Delete: resourceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"certificate_id"},
				ForceNew:      true,
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"certificate_id"},
				ForceNew:      true,
			},
			"certificate_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"certificate", "private_key", "certificate_name"},
			},
			"certificate_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"certificate_id"},
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error
	var response map[string]interface{}
	if sslId, ok := d.GetOk("certificate_id"); ok {
		action := "CreateCertificateByCertificateId"
		request["InstanceId"] = d.Get("instance_id")
		request["CertificateId"] = sslId.(string)
		if v, ok := d.GetOk("domain"); ok {
			request["Domain"] = v
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_waf_certificate", action, errs.AlibabaCloudSdkGoERROR)
		}
		certificateId, err := convert.FormatInt(response["CertificateId"])
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"], ":", certificateId))
		return resourceCertificateRead(d, meta)
	}

	request = make(map[string]interface{})
	action := "CreateCertificate"
	if v, ok := d.GetOk("certificate"); ok {
		request["Certificate"] = v
	} else {
		return errs.WrapErrorf(err, errs.RequiredWhenMsg, "certificate", "certificate_id", "null")
	}
	if v, ok := d.GetOk("private_key"); ok {
		request["PrivateKey"] = v
	} else {
		return errs.WrapErrorf(err, errs.RequiredWhenMsg, "private_key", "certificate_id", "null")
	}
	if v, ok := d.GetOk("certificate_name"); ok {
		request["CertificateName"] = v
	} else {
		return errs.WrapErrorf(err, errs.RequiredWhenMsg, "certificate_name", "certificate_id", "null")
	}
	request["Domain"] = d.Get("domain")
	request["InstanceId"] = d.Get("instance_id")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_waf_certificate", action, errs.AlibabaCloudSdkGoERROR)
	}

	actualCertificateId, err := convert.FormatInt(response["CertificateId"])
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"], ":", actualCertificateId))

	return resourceCertificateRead(d, meta)
}
func resourceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := NewService(client)
	object, err := wafOpenapiService.DescribeWafCertificate(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_waf_certificate wafOpenapiService.DescribeWafCertificate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("certificate_id", object["CertificateId"])
	d.Set("domain", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("certificate_name", object["CertificateName"])
	return nil
}
func resourceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy ResourceCertificate. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
