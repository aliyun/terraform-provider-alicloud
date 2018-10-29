package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCenBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthPackageAttachmentCreate,
		Read:   resourceAlicloudCenBandwidthPackageAttachmentRead,
		Delete: resourceAlicloudCenBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth_package_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	cenId := d.Get("instance_id").(string)
	cenBwpId := d.Get("bandwidth_package_id").(string)

	request := cbn.CreateAssociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.AssociateCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidBwpInstanceStatus, InvalidBwpBusinessStatus, InvalidCenInstanceStatus}) {
				return resource.RetryableError(fmt.Errorf("Associate bandwidth package %s to CEN %s timeout and got an error: %#v", cenBwpId, cenId, err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Associate bandwidth package %s to CEN %s got an error: %#v.", cenBwpId, cenId, err)
	}

	if err := cenService.WaitForCenBandwidthPackageAttachment(cenBwpId, InUse, DefaultCenTimeout); err != nil {
		return fmt.Errorf("Timeout when WaitForCenBandwidthPackageAttachment, CEN ID %s, bandwidth package ID %s, error info %#v.", cenId, cenBwpId, err)
	}

	d.SetId(cenBwpId)

	return resourceAlicloudCenBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	resp, err := cenService.DescribeCenBandwidthPackageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("instance_id", resp.CenIds.CenId[0])
	d.Set("bandwidth_package_id", resp.CenBandwidthPackageId)

	return nil
}

func resourceAlicloudCenBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	cenId := d.Get("instance_id").(string)
	cenBwpId := d.Get("bandwidth_package_id").(string)

	request := cbn.CreateUnassociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.UnassociateCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidBwpInstanceStatus, InvalidBwpBusinessStatus, InvalidCenInstanceStatus}) {
				return resource.RetryableError(fmt.Errorf("Unassociate bandwidth package %s from CEN %s timeout and got an error: %#v", cenBwpId, cenId, err))
			}

			return resource.NonRetryableError(fmt.Errorf("Unassociate bandwidth %s from CEN %s timeout and got an error: %#v", cenBwpId, cenId, err))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Unassociate bandwidth %s from CEN %s got an error: %#v.", cenBwpId, cenId, err)
	}

	if err := cenService.WaitForCenBandwidthPackageAttachment(cenBwpId, Idle, DefaultCenTimeout); err != nil {
		return fmt.Errorf("Timeout when WaitForCenBandwidthPackageAttachment, CEN ID %s, bandwidth package ID %s, error %#v", cenId, cenBwpId, err)
	}

	return nil
}
