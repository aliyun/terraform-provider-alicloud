package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"tags": tagsSchema(),
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

	return resourceAlicloudSlbCACertificateUpdate(d, meta)
}

func resourceAlicloudSlbCACertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	tags, err := slbService.DescribeTags(d.Id(), nil, TagResourceCertificate)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", slbService.tagsToMap(tags))

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
	slbService := SlbService{client}
	if err := slbService.setInstanceTags(d, TagResourceCertificate); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudSlbCACertificateRead(d, meta)
	}
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
			if IsExpectedErrors(err, SlbIsBusy) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"CACertificateId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(slbService.WaitForSlbCACertificate(d.Id(), Deleted, DefaultTimeoutMedium))
}
