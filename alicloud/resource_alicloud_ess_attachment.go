package alicloud

import (
	"fmt"
	"time"

	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEssAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssAttachmentCreate,
		Read:   resourceAliyunEssAttachmentRead,
		Update: resourceAliyunEssAttachmentUpdate,
		Delete: resourceAliyunEssAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 20,
				MinItems: 1,
			},

			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunEssAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("scaling_group_id").(string))

	return resourceAliyunEssAttachmentUpdate(d, meta)
}

func resourceAliyunEssAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	d.Partial(true)

	groupId := d.Id()
	if d.HasChange("instance_ids") {
		group, err := essService.DescribeScalingGroup(groupId)
		if err != nil {
			return WrapError(err)
		}
		if group.LifecycleState == string(Inactive) {
			return WrapError(fmt.Errorf("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState))
		} else {
			if err := essService.WaitForScalingGroup(group.ScalingGroupId, Active, DefaultTimeout); err != nil {
				return WrapError(err)
			}
		}
		o, n := d.GetChange("instance_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := convertArrayInterfaceToArrayString(ns.Difference(os).List())

		if len(add) > 0 {
			req := ess.CreateAttachInstancesRequest()
			req.ScalingGroupId = groupId
			s := reflect.ValueOf(req).Elem()

			if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				for i, id := range add {
					s.FieldByName(fmt.Sprintf("InstanceId%d", i+1)).Set(reflect.ValueOf(id))
				}

				_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.AttachInstances(req)
				})
				if err != nil {
					if IsExceptedError(err, IncorrectCapacityMaxSize) {
						instances, err := essService.DescribeScalingInstances(d.Id(), "", make([]string, 0), "")
						if err != nil {
							return resource.NonRetryableError(fmt.Errorf("DescribeScalingInstances got an error: %#v", err))
						}
						var autoAdded, attached []string
						if len(instances) > 0 {
							for _, inst := range instances {
								if inst.CreationType == "Attached" {
									attached = append(attached, inst.InstanceId)
								} else {
									autoAdded = append(autoAdded, inst.InstanceId)
								}
							}
						}
						if len(add) > group.MaxSize {
							return resource.NonRetryableError(fmt.Errorf("To attach %d instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size.", len(add), group.MaxSize))
						}

						if len(autoAdded) > 0 {
							if d.Get("force").(bool) {
								if err := essService.EssRemoveInstances(groupId, autoAdded); err != nil {
									return resource.NonRetryableError(err)
								}
								time.Sleep(5)
								return resource.RetryableError(fmt.Errorf("Autocreated result in attaching instances got an error: %#v", err))
							} else {
								return resource.NonRetryableError(fmt.Errorf("To attach the instances, the total capacity will be greater than the scaling group max size %d."+
									"Please enlarge scaling group max size or set 'force' to true to remove autocreated instances: %#v.", group.MaxSize, autoAdded))
							}
						}

						if len(attached) > 0 {
							return resource.NonRetryableError(fmt.Errorf("To attach the instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size or remove already attached instances: %#v.", group.MaxSize, attached))
						}
					}
					if IsExceptedError(err, ScalingActivityInProgress) {
						time.Sleep(5)
						return resource.RetryableError(fmt.Errorf("Progress results in Attaching instances got an error: %#v", err))
					}
					return resource.NonRetryableError(fmt.Errorf("Attaching instances got an error: %#v", err))
				}
				return nil
			}); err != nil {
				return err
			}

			if err := resource.Retry(3*time.Minute, func() *resource.RetryError {

				instances, err := essService.DescribeScalingInstances(d.Id(), "", add, "")
				if err != nil {
					return resource.NonRetryableError(err)
				}
				if len(instances) < 0 {
					return resource.RetryableError(fmt.Errorf("There are no ECS instances have been attached."))
				}

				for _, inst := range instances {
					if inst.LifecycleState != string(InService) {
						return resource.RetryableError(fmt.Errorf("There are still ECS instances are not %s.", string(InService)))
					}
				}
				return nil
			}); err != nil {
				return err
			}
		}
		if len(remove) > 0 {
			if err := essService.EssRemoveInstances(groupId, convertArrayInterfaceToArrayString(remove)); err != nil {
				return err
			}
		}

		d.SetPartial("instance_ids")
	}

	d.Partial(false)

	return resourceAliyunEssAttachmentRead(d, meta)
}

func resourceAliyunEssAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	instances, err := essService.DescribeScalingInstances(d.Id(), "", make([]string, 0), string(Attached))

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS scaling instances: %#v", err)
	}

	if len(instances) < 1 {
		d.SetId("")
		return nil
	}

	var instanceIds []string
	for _, inst := range instances {
		instanceIds = append(instanceIds, inst.InstanceId)
	}

	d.Set("scaling_group_id", instances[0].ScalingGroupId)
	d.Set("instance_ids", instanceIds)

	return nil
}

func resourceAliyunEssAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	group, err := essService.DescribeScalingGroup(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if group.LifecycleState != string(Active) {
		return WrapError(fmt.Errorf("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState))
	}

	return WrapError(essService.EssRemoveInstances(d.Id(), convertArrayInterfaceToArrayString(d.Get("instance_ids").(*schema.Set).List())))
}

func convertArrayInterfaceToArrayString(elm []interface{}) (arr []string) {
	if len(elm) < 1 {
		return
	}
	for _, e := range elm {
		arr = append(arr, e.(string))
	}
	return
}
