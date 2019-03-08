package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunCommonBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunCommonBandwidthPackageAttachmentCreate,
		Read:   resourceAliyunCommonBandwidthPackageAttachmentRead,
		Delete: resourceAliyunCommonBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}
	args := vpc.CreateAddCommonBandwidthPackageIpRequest()
	args.BandwidthPackageId = Trim(d.Get("bandwidth_package_id").(string))
	args.IpInstanceId = Trim(d.Get("instance_id").(string))
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		ar := args
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AddCommonBandwidthPackageIp(ar)
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Add bandwidth package ip got an error: %#v", err))
		}
		return nil
	}); err != nil {
		return err
	}
	//check the common bandwidth package attachment
	if err := commonBandwidthPackageService.WaitForCommonBandwidthPackageAttachment(args.BandwidthPackageId, args.IpInstanceId, 5*DefaultTimeout); err != nil {
		return fmt.Errorf("Wait for common bandwidth package attachment got error: %#v", err)
	}

	d.SetId(args.BandwidthPackageId + COLON_SEPARATED + args.IpInstanceId)

	return resourceAliyunCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceAliyunCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}

	bandwidthPackageId, ipInstanceId, err := GetBandwidthPackageIdAndIpInstanceId(d, meta)
	err = commonBandwidthPackageService.DescribeCommonBandwidthPackageAttachment(bandwidthPackageId, ipInstanceId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Common Bandwidth Package Attribute: %#v", err)
	}

	d.Set("bandwidth_package_id", bandwidthPackageId)
	d.Set("instance_id", ipInstanceId)
	return nil
}

func resourceAliyunCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	commonBandwidthPackageService := CommonBandwidthPackageService{client}

	bandwidthPackageId, ipInstanceId, err := GetBandwidthPackageIdAndIpInstanceId(d, meta)
	if err != nil {
		return err
	}

	request := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
	request.BandwidthPackageId = bandwidthPackageId
	request.IpInstanceId = ipInstanceId

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.RemoveCommonBandwidthPackageIp(request)
		})
		//Waiting for unassociate the common bandwidth package
		if err != nil {
			if IsExceptedError(err, TaskConflict) {
				return resource.RetryableError(fmt.Errorf("Unassociate Common Bandwidth Package timeout and got an error:%#v.", err))
			}
		}
		//Eusure the instance has been unassociated truly.
		err = commonBandwidthPackageService.DescribeCommonBandwidthPackageAttachment(bandwidthPackageId, ipInstanceId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("Unassociate Common Bandwidth Package timeout."))
	})
}
