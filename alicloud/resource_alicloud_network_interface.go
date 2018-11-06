package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNetworkInterfaceCreate,
		Read:   resourceAliyunNetworkInterfaceRead,
		Update: resourceAliyunNetworkInterfaceUpdate,
		Delete: resourceAliyunNetworkInterfaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				MinItems: 1,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ips": &schema.Schema{
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				MaxItems:      10,
				ConflictsWith: []string{"private_ips_count"},
			},
			"private_ips_count": &schema.Schema{
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validateIntegerInRange(1, 10),
				ConflictsWith: []string{"private_ips"},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	args := ecs.CreateCreateNetworkInterfaceRequest()

	args.VSwitchId = d.Get("vswitch_id").(string)
	groups := d.Get("security_groups").(*schema.Set).List()

	args.SecurityGroupId = groups[0].(string)

	if primaryIpAddress, ok := d.GetOk("private_ip"); ok {
		args.PrimaryIpAddress = primaryIpAddress.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		args.NetworkInterfaceName = name.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		args.Description = description.(string)
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateNetworkInterface(args)
	})
	if err != nil {
		return fmt.Errorf("Create NetworkInterface faield, %#v", err)
	}

	resp := raw.(*ecs.CreateNetworkInterfaceResponse)
	d.SetId(resp.NetworkInterfaceId)

	if err := ecsService.WaitForEcsNetworkInterface(d.Id(), Available, DefaultTimeout); err != nil {
		return fmt.Errorf("Wait NetwortInterface(%s) to be avaialbe failed, %#v", d.Id(), err)
	}

	return resourceAliyunNetworkInterfaceUpdate(d, meta)
}

func resourceAliyunNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	eni, err := ecsService.DescribeNetworkInterfaceById("", d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Describe NetworkInterface(%s) failed, %#v", d.Id(), err)
	}

	d.Set("name", eni.NetworkInterfaceName)
	d.Set("description", eni.Description)
	d.Set("vswitch_id", eni.VSwitchId)
	d.Set("private_ip", eni.PrivateIpAddress)
	d.Set("security_groups", eni.SecurityGroupIds.SecurityGroupId)
	privateIps := make([]string, 0, len(eni.PrivateIpSets.PrivateIpSet))
	for i := range eni.PrivateIpSets.PrivateIpSet {
		if !eni.PrivateIpSets.PrivateIpSet[i].Primary {
			privateIps = append(privateIps, eni.PrivateIpSets.PrivateIpSet[i].PrivateIpAddress)
		}
	}
	d.Set("private_ips", privateIps)

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceEni)
	if err != nil && !NotFoundError(err) {
		return fmt.Errorf("DescribeTags of NetworkInterface(%s) failed, %#v", d.Id(), err)
	}

	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	d.Partial(true)

	attributeUpdate := false
	args := ecs.CreateModifyNetworkInterfaceAttributeRequest()
	args.NetworkInterfaceId = d.Id()

	if !d.IsNewResource() && d.HasChange("description") {
		args.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if !d.IsNewResource() && d.HasChange("name") {
		args.NetworkInterfaceName = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("security_groups") {
		securityGroups := expandStringList(d.Get("security_groups").(*schema.Set).List())
		if len(securityGroups) > 1 || !d.IsNewResource() {
			args.SecurityGroupId = &securityGroups
			attributeUpdate = true
		}
	}

	if attributeUpdate {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyNetworkInterfaceAttribute(args)
		})
		if err != nil {
			return fmt.Errorf("Modify NetworkInterface(%s) attributes failed, %#v", d.Id(), err)
		}
		d.SetPartial("security_groups")
		d.SetPartial("description")
		d.SetPartial("name")
	}

	if d.HasChange("private_ips") {
		o, n := d.GetChange("private_ips")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		unassignIps := os.Difference(ns)
		if unassignIps.Len() > 0 {
			unassignIpList := expandStringList(unassignIps.List())
			args := ecs.CreateUnassignPrivateIpAddressesRequest()
			args.NetworkInterfaceId = d.Id()
			args.PrivateIpAddress = &unassignIpList
			err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
				_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.UnassignPrivateIpAddresses(args)
				})
				if err != nil {
					if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
						return resource.RetryableError(fmt.Errorf("Unassign private IP address (%#v) failed, %#v", unassignIpList, err))
					}
					return resource.NonRetryableError(fmt.Errorf("Unassign private IP address (%#v) failed, %#v", unassignIpList, err))
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		assignIps := ns.Difference(os)
		if assignIps.Len() > 0 {
			assignIpList := expandStringList(assignIps.List())
			args := ecs.CreateAssignPrivateIpAddressesRequest()
			args.NetworkInterfaceId = d.Id()
			args.PrivateIpAddress = &assignIpList
			err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
				_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.AssignPrivateIpAddresses(args)
				})
				if err != nil {
					if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
						return resource.RetryableError(fmt.Errorf("Assign private IP address (%#v) failed, %#v", assignIpList, err))
					}
					return resource.NonRetryableError(fmt.Errorf("Assign private IP address (%#v) failed, %#v", assignIpList, err))
				}

				return nil
			})
			if err != nil {
				return err
			}
		}

		if err := ecsService.WaitForPrivateIpsListChanged(d.Id(), expandStringList(ns.List())); err != nil {
			return err
		}

		d.SetPartial("private_ips")
	}

	if d.HasChange("private_ips_count") {
		privateIpList := expandStringList(d.Get("private_ips").(*schema.Set).List())
		o, n := d.GetChange("private_ips_count")
		if o != nil && n != nil && n != len(privateIpList) {
			diff := n.(int) - o.(int)
			if diff > 0 {
				args := ecs.CreateAssignPrivateIpAddressesRequest()
				args.NetworkInterfaceId = d.Id()
				args.SecondaryPrivateIpAddressCount = requests.NewInteger(diff)
				err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
					_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.AssignPrivateIpAddresses(args)
					})

					if err != nil {
						if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
							return resource.RetryableError(fmt.Errorf("Assign private IP address (%s->%d) failed, %#v", o, n, err))
						}
						return resource.NonRetryableError(fmt.Errorf("Assign private address (%d->%d) failed, %#v", o, n, err))
					}
					return nil
				})
				if err != nil {
					return err
				}
			}

			if diff < 0 {
				diff *= -1
				unassignIps := privateIpList[:diff]
				err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
					args := ecs.CreateUnassignPrivateIpAddressesRequest()
					args.NetworkInterfaceId = d.Id()
					args.PrivateIpAddress = &unassignIps
					_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.UnassignPrivateIpAddresses(args)
					})
					if err != nil {
						if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
							return resource.RetryableError(fmt.Errorf("Unassign private IP address (%d->%d) failed, %#v", o, n, err))
						}
						return resource.RetryableError(fmt.Errorf("Unassign private IP address (%d->%d) failed, %#v", o, n, err))
					}
					return nil
				})
				if err != nil {
					return err
				}
			}

			if err := ecsService.WaitForPrivateIpsCountChanged(d.Id(), n.(int)); err != nil {
				return err
			}

			d.SetPartial("private_ips_count")
		}
	}

	if err := setTags(client, TagResourceEni, d); err != nil {
		return fmt.Errorf("SetTags of NetworkInterface(%s) failed, %#v", d.Id(), err)
	} else {
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceAliyunNetworkInterfaceRead(d, meta)
}

func resourceAliyunNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	args := ecs.CreateDeleteNetworkInterfaceRequest()
	args.NetworkInterfaceId = d.Id()

	return resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteNetworkInterface(args)
		})
		if err != nil {
			if IsExceptedErrors(err, NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(fmt.Errorf("Delete NetworkInterface(%s) failed, %#v", d.Id(), err))
			}
			return resource.NonRetryableError(fmt.Errorf("Delete NetworkInterface(%s) failed: %#v", d.Id(), err))
		}

		_, err = ecsService.DescribeNetworkInterfaceById("", d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Describe NetwortInterface(%s) failed, %#v", d.Id(), err))
		}

		return nil
	})
}
