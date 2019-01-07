package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNetworkInterfaceAttachmentCreate,
		Read:   resourceAliyunNetworkInterfaceAttachmentRead,
		Delete: resourceAliyunNetworkInterfaceAttachmentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	eniId := d.Get("network_interface_id").(string)
	instanceId := d.Get("instance_id").(string)

	args := ecs.CreateAttachNetworkInterfaceRequest()
	args.InstanceId = instanceId
	args.NetworkInterfaceId = eniId

	_, err := ecsService.DescribeNetworkInterfaceById(instanceId, eniId)
	if err != nil {
		if NotFoundError(err) {
			err := resource.Retry(DefaultTimeout*time.Minute, func() *resource.RetryError {
				_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.AttachNetworkInterface(args)
				})
				if err != nil {
					if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
						return resource.RetryableError(fmt.Errorf("Attach NetworkInterface(%s) failed, %#v", eniId, err))
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Attach NetworkInterface (%s) failed, %#v", eniId, err)
			}

			if err := ecsService.WaitForEcsNetworkInterface(eniId, InUse, DefaultTimeout); err != nil {
				return fmt.Errorf("Wait for the status of NetworkInterface (%s) changes to InUse failed, %#v", eniId, err)
			}
		} else {
			return fmt.Errorf("Describe NetworkInterface (%s) failed when create attachment, %#v", eniId, err)
		}
	}

	d.SetId(eniId + COLON_SEPARATED + instanceId)

	return resourceAliyunNetworkInterfaceAttachmentRead(d, meta)
}

func resourceAliyunNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	eniId, instanceId, err := getEniIDAndInstanceID(d, meta)
	if err != nil {
		return err
	}

	eni, err := ecsService.DescribeNetworkInterfaceById(instanceId, eniId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeNetworkInterfaceAttribute failed, %#v", err)
	}

	d.Set("instance_id", eni.InstanceId)
	d.Set("network_interface_id", eni.NetworkInterfaceId)

	return nil
}

func resourceAliyunNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	eniId, instanceId, err := getEniIDAndInstanceID(d, meta)
	if err != nil {
		return err
	}

	_, err = ecsService.DescribeNetworkInterfaceById(instanceId, eniId)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("Describe NetworkInterface (%s) failed when detach it from Instance (%s), %#v", eniId, instanceId, err)
	}

	args := ecs.CreateDetachNetworkInterfaceRequest()
	args.InstanceId = instanceId
	args.NetworkInterfaceId = eniId

	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachNetworkInterface(args)
		})
		if err != nil {
			if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(fmt.Errorf("Detach NetworkInterface (%s) from Instance (%s) failed, %#v", eniId, instanceId, err))
			}
			return resource.NonRetryableError(fmt.Errorf("Detach NetworkInterface (%s) from Instance (%s) failed, %#v", eniId, instanceId, err))
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := ecsService.WaitForEcsNetworkInterface(eniId, Available, DefaultTimeout); err != nil {
		return fmt.Errorf("Wait for the NetworkInterface(%s) status changes failed, %#v", eniId, err)
	}

	return nil
}

func getEniIDAndInstanceID(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid resource id (%s)", d.Id())
	}

	return parts[0], parts[1], nil
}
