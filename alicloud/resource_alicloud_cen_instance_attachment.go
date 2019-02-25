package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCenInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInstanceAttachmentCreate,
		Read:   resourceAlicloudCenInstanceAttachmentRead,
		Delete: resourceAlicloudCenInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	cenId := d.Get("instance_id").(string)
	instanceId := d.Get("child_instance_id").(string)
	instanceRegionId := d.Get("child_instance_region_id").(string)
	instanceType, err := cenService.GetCenInstanceType(instanceId)
	if err != nil {
		return err
	}

	request := cbn.CreateAttachCenChildInstanceRequest()
	request.CenId = cenId
	request.ChildInstanceId = instanceId
	request.ChildInstanceType = instanceType
	request.ChildInstanceRegionId = instanceRegionId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.AttachCenChildInstance(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidCenInstanceStatus, InvalidChildInstanceStatus}) {
				return resource.RetryableError(fmt.Errorf("Attach CEN child instance timeout and got an error: %#v", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Attach child instance %s to CEN %s and got an error: %#v.", instanceId, cenId, err)
	}

	if err := cenService.WaitForCenChildInstanceAttached(instanceId, cenId, Status("Attached"), DefaultCenTimeoutLong); err != nil {
		return fmt.Errorf("Timeout when WaitForCenChildInstanceAttached, CEN ID %s, child instance ID %s, error info %#v.", cenId, instanceId, err)
	}

	d.SetId(cenId + COLON_SEPARATED + instanceId)

	return resourceAlicloudCenInstanceAttachmentRead(d, meta)
}

func resourceAlicloudCenInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	cenId, instanceId, err := cenService.GetCenIdAndAnotherId(d.Id())
	if err != nil {
		return err
	}

	resp, err := cenService.DescribeCenAttachedChildInstanceById(instanceId, cenId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("instance_id", resp.CenId)
	d.Set("child_instance_id", resp.ChildInstanceId)
	d.Set("child_instance_region_id", resp.ChildInstanceRegionId)

	return nil
}

func resourceAlicloudCenInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	instanceRegionId := d.Get("child_instance_region_id").(string)
	cenId, instanceId, err := cenService.GetCenIdAndAnotherId(d.Id())
	if err != nil {
		return err
	}
	instanceType, err := cenService.GetCenInstanceType(instanceId)
	if err != nil {
		return err
	}

	request := cbn.CreateDetachCenChildInstanceRequest()
	request.CenId = cenId
	request.ChildInstanceId = instanceId
	request.ChildInstanceType = instanceType
	request.ChildInstanceRegionId = instanceRegionId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DetachCenChildInstance(request)
		})
		if err != nil {
			if IsExceptedError(err, ParameterInstanceIdNotExist) {
				return nil
			}
			if IsExceptedError(err, InvalidCenInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Detach CEN child instance timeout and got an error: %#v", err))
			}

			return resource.NonRetryableError(err)
		}

		_, err = cenService.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Detach child instance %s from CEN %s got an error: %#v.", instanceId, cenId, err)
	}

	if err := cenService.WaitForCenChildInstanceDetached(instanceId, cenId, DefaultCenTimeoutLong); err != nil {
		return fmt.Errorf("Timeout when WaitForCenChildInstanceDetached, CEN ID %s, child instance ID %s, error info: %#v", cenId, instanceId, err)
	}

	return nil
}
