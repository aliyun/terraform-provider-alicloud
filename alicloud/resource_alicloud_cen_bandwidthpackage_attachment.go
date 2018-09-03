package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCenBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthPackageAttachmentCreate,
		Read:   resourceAlicloudCenBandwidthPackageAttachmentRead,
		Delete: resourceAlicloudCenBandwidthPackageAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"cen_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_bandwidthpackage_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	cenId := d.Get("cen_id").(string)
	cenBwpId := d.Get("cen_bandwidthpackage_id").(string)

	request := cbn.CreateAssociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.cenconn.AssociateCenBandwidthPackage(request)
		if err != nil {
			if IsExceptedError(err, InvalidBwpInstanceStatus) || IsExceptedError(err, InvalidBwpBusinessStatus) || IsExceptedError(err, InvalidCenInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Associate CEN bandwidth package %s timeout and got an error: %#v", cenBwpId, err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Associate CEN bandwidth package %s to CEN %s got an error: %#v.", cenBwpId, cenId, err)
	}

	if err := client.WaitForCenBandwidthPackageAssociate(cenBwpId, InUse, 60); err != nil {
		return fmt.Errorf("Timeout when WaitForCenBandwidthPackageInUse")
	}

	d.SetId(d.Get("cen_bandwidthpackage_id").(string) + ":" + d.Get("cen_id").(string))

	return resourceAlicloudCenBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	cenBwpId, cenId, err := getCenIdAndAnotherId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DescribeCenBandwidthPackageById(cenBwpId, cenId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cen_id", resp.CenIds.CenId[0])
	d.Set("cen_bandwidthpackage_id", resp.CenBandwidthPackageId)

	return nil
}

func resourceAlicloudCenBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	cenBwpId, cenId, err := getCenIdAndAnotherId(d.Id())
	if err != nil {
		return err
	}

	request := cbn.CreateUnassociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = client.cenconn.UnassociateCenBandwidthPackage(request)
		if err != nil {
			if IsExceptedError(err, InvalidBwpInstanceStatus) || IsExceptedError(err, InvalidBwpBusinessStatus) {
				return resource.RetryableError(fmt.Errorf("Unassociate cen bandwidth package %s timeout and got an error: %#v", cenBwpId, err))
			}

			return resource.NonRetryableError(fmt.Errorf("Unassociate cen bandwidth %s timeout and got an error: %#v", cenBwpId, err))
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Unassociate cen bandwidth %s from CEN %s got an error: %#v.", cenBwpId, cenId, err)
	}

	if err := client.WaitForCenBandwidthPackage(cenBwpId, Idle, 60); err != nil {
		return fmt.Errorf("Timeout when WaitForCenBandwidthPackageUnassociate")
	}

	return nil
}
