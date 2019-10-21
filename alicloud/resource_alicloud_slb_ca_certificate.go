package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudSlbCACertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbCACertificateCreate,
		Read:   resourceAlicloudSlbCACertificateRead,
		Update: resourceAlicloudSlbCACertificateUpdate,
		Delete: resourceAlicloudSlbCACertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSlbCACertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	request := slb.CreateUploadCACertificateRequest()
	request.RegionId = client.RegionId

	if val, ok := d.GetOk("name"); ok && val.(string) != "" {
		request.CACertificateName = val.(string)
	}

	if val, ok := d.GetOk("resource_group_id"); ok && val.(string) != "" {
		request.ResourceGroupId = val.(string)
	}

	if val, ok := d.GetOk("ca_certificate"); ok && val.(string) != "" {
		request.CACertificate = val.(string)
	} else {
		return WrapError(Error("UploadCACertificate got an error, ca_certificate should be not null"))
	}

	raw, err := slbService.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.UploadCACertificate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_ca_certificate", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*slb.UploadCACertificateResponse)

	d.SetId(response.CACertificateId)

	return resourceAlicloudSlbCACertificateRead(d, meta)
}

func resourceAlicloudSlbCACertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	object, err := slbService.DescribeSlbCACertificate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	err = d.Set("name", object.CACertificateName)
	err = d.Set("resource_group_id", object.ResourceGroupId)
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudSlbCACertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") {
		request := slb.CreateSetCACertificateNameRequest()
		request.RegionId = client.RegionId
		request.CACertificateId = d.Id()
		request.CACertificateName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetCACertificateName(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudSlbCACertificateRead(d, meta)
}

func resourceAlicloudSlbCACertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteCACertificateRequest()
	request.RegionId = client.RegionId
	request.CACertificateId = d.Id()

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteCACertificate(request)
		})
		if err != nil {
			if IsExceptedErrors(err, SlbIsBusy) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExceptedError(err, SlbCACertificateIdNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(slbService.WaitForSlbCACertificate(d.Id(), Deleted, DefaultTimeoutMedium))
}
