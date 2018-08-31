package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCenInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInstanceAttachmentCreate,
		Read:   resourceAlicloudCenInstanceAttachmentRead,
		Delete: resourceAlicloudCenInstanceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"cen_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"child_instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "VPC" && value != "VBR" {
						errors = append(errors, fmt.Errorf("%s must be one of VPC and VBR", k))
					}

					return
				},
			},
			"child_instance_region_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	cenId := d.Get("cen_id").(string)
	instanceId := d.Get("child_instance_id").(string)
	instanceType := d.Get("child_instance_type").(string)
	instanceRegionId := d.Get("child_instance_region_id").(string)

	request := cbn.CreateAttachCenChildInstanceRequest()
	request.CenId = cenId
	request.ChildInstanceId = instanceId
	request.ChildInstanceType = instanceType
	request.ChildInstanceRegionId = instanceRegionId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.cenconn.AttachCenChildInstance(request)
		if err != nil {
			if IsExceptedError(err, InvalidCenInstanceStatus) || IsExceptedError(err, InvalidChildInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Attach CEN child instance %s timeout and got an error: %#v", instanceId, err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Attach child instance %s to CEN %s got an error: %#v.", instanceId, cenId, err)
	}

	waitTime := 60
	if instanceType == "VBR" {
		waitTime = 180 //attach take much longer time
	}

	if err := client.WaitForCenChildInstanceAttached(instanceId, cenId, Status("Attached"), waitTime); err != nil {
		return fmt.Errorf("Timeout when WaitForCenChildInstanceAttached")
	}

	d.SetId(d.Get("child_instance_id").(string) + ":" + d.Get("cen_id").(string))

	return resourceAlicloudCenInstanceAttachmentRead(d, meta)
}

func resourceAlicloudCenInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	instanceId, cenId, err := getCenIdAndAnotherId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DescribeCenAttachedChildInstanceById(instanceId, cenId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cen_id", resp.CenId)
	d.Set("child_instance_id", resp.ChildInstanceId)
	d.Set("child_instance_type", resp.ChildInstanceType)
	d.Set("child_instance_region_id", resp.ChildInstanceRegionId)

	return nil
}

func resourceAlicloudCenInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	instanceId, cenId, err := getCenIdAndAnotherId(d.Id())
	instanceType := d.Get("child_instance_type").(string)
	instanceRegionId := d.Get("child_instance_region_id").(string)
	if err != nil {
		return err
	}

	request := cbn.CreateDetachCenChildInstanceRequest()
	request.CenId = cenId
	request.ChildInstanceId = instanceId
	request.ChildInstanceType = instanceType
	request.ChildInstanceRegionId = instanceRegionId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = client.cenconn.DetachCenChildInstance(request)
		if err != nil {
			if IsExceptedError(err, ParameterInstanceIdNotExist) {
				return nil
			} else if IsExceptedError(err, InvalidCenInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Detach CEN child instance %s got an error: %#v", instanceId, err))
			}

			return resource.NonRetryableError(fmt.Errorf("Detach CEN child instancee %s timeout and got an error: %#v", instanceId, err))
		}

		_, err := client.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("While detach CEN child instance %s, describing CEN child instance got an error: %#v.", instanceId, err))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Detach child instance %s to CEN %s got an error: %#v.", instanceId, cenId, err)
	}

	waitTime := 60
	if instanceType == "VBR" {
		waitTime = 180 //attach take much longer time
	}

	if err := client.WaitForCenChildInstanceDetached(instanceId, cenId, waitTime); err != nil {
		return fmt.Errorf("Timeout when WaitForCenChildInstanceDetached")
	}

	return nil
}
