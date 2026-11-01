package slb

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCertificateCreate,
		Read:   resourceServerCertificateRead,
		Update: resourceServerCertificateUpdate,
		Delete: resourceServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_certificate": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateDiffSuppressFunc,
			},
			"private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateDiffSuppressFunc,
			},
			"alicloud_certificate_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alicloud_certifacte_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'alicloud_certifacte_id' has been deprecated from provider version 1.68.0. Use 'alicloud_certificate_id' replaces it.",
			},
			"alicloud_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alicloud_certifacte_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'alicloud_certifacte_name' has been deprecated from provider version 1.68.0. Use 'alicloud_certificate_name' replaces it.",
			},
			"alicloud_certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateUploadServerCertificateRequest()
	request.RegionId = client.RegionId

	if val, ok := d.GetOk("name"); ok && val != "" {
		request.ServerCertificateName = val.(string)
	}

	if val, ok := d.GetOk("server_certificate"); ok && val != "" {
		request.ServerCertificate = val.(string)
	}

	if val, ok := d.GetOk("private_key"); ok && val != "" {
		request.PrivateKey = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_region_id"); ok && val != "" {
		request.AliCloudCertificateRegionId = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_id"); ok && val != "" {
		request.AliCloudCertificateId = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_name"); ok && val != "" {
		request.AliCloudCertificateName = val.(string)
	}

	if val, ok := d.GetOk("resource_group_id"); ok && val != "" {
		request.ResourceGroupId = val.(string)
	}

	// check server_certificate and private_key
	if request.AliCloudCertificateId == "" {
		if val := strings.Trim(request.ServerCertificate, " "); val == "" {
			return errs.WrapError(errs.Error("UploadServerCertificate got an error, as server_certificate should be not null when alicloud_certificate_id is null."))
		}

		if val := strings.Trim(request.PrivateKey, " "); val == "" {
			return errs.WrapError(errs.Error("UploadServerCertificate got an error, as either private_key or private_file  should be not null when alicloud_certificate_id is null."))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.UploadServerCertificate(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.UploadServerCertificateResponse)
	d.SetId(response.ServerCertificateId)

	return resourceServerCertificateUpdate(d, meta)
}

func resourceServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	tagsMap, err := slbService.ListTagResources(d.Id(), "certificate")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(tagsMap))

	serverCertificate, err := slbService.DescribeSlbServerCertificate(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if err := d.Set("name", serverCertificate.ServerCertificateName); err != nil {
		return errs.WrapError(err)
	}

	if serverCertificate.AliCloudCertificateId != "" {
		if err := d.Set("alicloud_certificate_id", serverCertificate.AliCloudCertificateId); err != nil {
			return errs.WrapError(err)
		}
	}

	if serverCertificate.AliCloudCertificateName != "" {
		if err := d.Set("alicloud_certificate_name", serverCertificate.AliCloudCertificateName); err != nil {
			return errs.WrapError(err)
		}
	}

	if serverCertificate.ResourceGroupId != "" {
		if err := d.Set("resource_group_id", serverCertificate.ResourceGroupId); err != nil {
			return errs.WrapError(err)
		}
	}

	return nil
}

func resourceServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	if err := slbService.setInstanceTags(d, tags.TagResourceCertificate); err != nil {
		return errs.WrapError(err)
	}
	if d.IsNewResource() {
		return resourceServerCertificateRead(d, meta)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetServerCertificateNameRequest()
		request.RegionId = client.RegionId
		request.ServerCertificateId = d.Id()
		request.ServerCertificateName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetServerCertificateName(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceServerCertificateRead(d, meta)
}

func resourceServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)

	request := slb.CreateDeleteServerCertificateRequest()
	request.RegionId = client.RegionId
	request.ServerCertificateId = d.Id()
	err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteServerCertificate(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"CertificateAndPrivateKeyIsRefered"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)

		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ServerCertificateId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return errs.WrapError(slbService.WaitForSlbServerCertificate(d.Id(), Deleted, DefaultTimeoutMedium))

}
