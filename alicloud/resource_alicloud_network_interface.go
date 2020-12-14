package alicloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				MinItems: 1,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ips": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				MaxItems:      10,
				ConflictsWith: []string{"private_ips_count"},
			},
			"private_ips_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.IntBetween(0, 10),
				ConflictsWith: []string{"private_ips"},
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
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
	request := ecs.CreateCreateNetworkInterfaceRequest()
	request.RegionId = client.RegionId
	request.VSwitchId = d.Get("vswitch_id").(string)
	groups := d.Get("security_groups").(*schema.Set).List()

	request.SecurityGroupId = groups[0].(string)

	if primaryIpAddress, ok := d.GetOk("private_ip"); ok {
		request.PrimaryIpAddress = primaryIpAddress.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		request.NetworkInterfaceName = name.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateNetworkInterface(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_network_interface", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	object := raw.(*ecs.CreateNetworkInterfaceResponse)
	d.SetId(object.NetworkInterfaceId)

	if err := ecsService.WaitForNetworkInterface(d.Id(), Available, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunNetworkInterfaceUpdate(d, meta)
}

func resourceAliyunNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeNetworkInterface(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_network_interface ecsService.DescribeNetworkInterface Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.NetworkInterfaceName)
	d.Set("description", object.Description)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("private_ip", object.PrivateIpAddress)
	d.Set("security_groups", object.SecurityGroupIds.SecurityGroupId)
	privateIps := make([]string, 0, len(object.PrivateIpSets.PrivateIpSet))
	for i := range object.PrivateIpSets.PrivateIpSet {
		if !object.PrivateIpSets.PrivateIpSet[i].Primary {
			privateIps = append(privateIps, object.PrivateIpSets.PrivateIpSet[i].PrivateIpAddress)
		}
	}
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("private_ips", privateIps)
	d.Set("private_ips_count", len(privateIps))
	d.Set("mac", object.MacAddress)

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceEni)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	if len(tags) > 0 {
		d.Set("tags", ecsTagsToMap(tags))
	}

	return nil
}

func resourceAliyunNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	d.Partial(true)

	attributeUpdate := false
	request := ecs.CreateModifyNetworkInterfaceAttributeRequest()
	request.NetworkInterfaceId = d.Id()
	request.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		request.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	if !d.IsNewResource() && d.HasChange("name") {
		request.NetworkInterfaceName = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("security_groups") {
		securityGroups := expandStringList(d.Get("security_groups").(*schema.Set).List())
		if len(securityGroups) > 1 || !d.IsNewResource() {
			request.SecurityGroupId = &securityGroups
			attributeUpdate = true
		}
	}

	if attributeUpdate {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyNetworkInterfaceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("security_groups")
		d.SetPartial("description")
		d.SetPartial("name")
	}

	if d.HasChange("private_ips") {
		oldIps, newIps := d.GetChange("private_ips")
		oldIpsSet := oldIps.(*schema.Set)
		newIpsSet := newIps.(*schema.Set)

		unAssignIps := oldIpsSet.Difference(newIpsSet)
		if unAssignIps.Len() > 0 {
			unAssignIpList := expandStringList(unAssignIps.List())
			unAssignPrivateIpAddressesRequest := ecs.CreateUnassignPrivateIpAddressesRequest()
			unAssignPrivateIpAddressesRequest.RegionId = client.RegionId
			unAssignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
			unAssignPrivateIpAddressesRequest.PrivateIpAddress = &unAssignIpList
			err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.UnassignPrivateIpAddresses(unAssignPrivateIpAddressesRequest)
				})
				if err != nil {
					if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(unAssignPrivateIpAddressesRequest.GetActionName(), raw, unAssignPrivateIpAddressesRequest.RpcRequest, unAssignPrivateIpAddressesRequest)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}

		assignIps := newIpsSet.Difference(oldIpsSet)
		if assignIps.Len() > 0 {
			assignIpList := expandStringList(assignIps.List())
			assignPrivateIpAddressesRequest := ecs.CreateAssignPrivateIpAddressesRequest()
			assignPrivateIpAddressesRequest.RegionId = client.RegionId
			assignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
			assignPrivateIpAddressesRequest.PrivateIpAddress = &assignIpList
			err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
				raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
					return ecsClient.AssignPrivateIpAddresses(assignPrivateIpAddressesRequest)
				})
				if err != nil {
					if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
						return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
					}
					return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
				}
				addDebug(assignPrivateIpAddressesRequest.GetActionName(), raw, assignPrivateIpAddressesRequest.RpcRequest, assignPrivateIpAddressesRequest)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}

		if err := ecsService.WaitForPrivateIpsListChanged(d.Id(), expandStringList(newIpsSet.List())); err != nil {
			return WrapError(err)
		}

		d.SetPartial("private_ips")
	}

	if d.HasChange("private_ips_count") {
		privateIpList := expandStringList(d.Get("private_ips").(*schema.Set).List())
		oldIpsCount, newIpsCount := d.GetChange("private_ips_count")
		if oldIpsCount != nil && newIpsCount != nil && newIpsCount != len(privateIpList) {
			diff := newIpsCount.(int) - oldIpsCount.(int)
			if diff > 0 {
				assignPrivateIpAddressesRequest := ecs.CreateAssignPrivateIpAddressesRequest()
				assignPrivateIpAddressesRequest.RegionId = client.RegionId
				assignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
				assignPrivateIpAddressesRequest.SecondaryPrivateIpAddressCount = requests.NewInteger(diff)
				err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.AssignPrivateIpAddresses(assignPrivateIpAddressesRequest)
					})

					if err != nil {
						if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
							return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), AlibabaCloudSdkGoERROR))
						}
						return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), AlibabaCloudSdkGoERROR))
					}
					addDebug(assignPrivateIpAddressesRequest.GetActionName(), raw, assignPrivateIpAddressesRequest.RpcRequest, assignPrivateIpAddressesRequest)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), assignPrivateIpAddressesRequest.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}

			if diff < 0 {
				diff *= -1
				unAssignIps := privateIpList[:diff]
				unAssignPrivateIpAddressesRequest := ecs.CreateUnassignPrivateIpAddressesRequest()
				unAssignPrivateIpAddressesRequest.RegionId = client.RegionId
				err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
					unAssignPrivateIpAddressesRequest.NetworkInterfaceId = d.Id()
					unAssignPrivateIpAddressesRequest.PrivateIpAddress = &unAssignIps
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.UnassignPrivateIpAddresses(unAssignPrivateIpAddressesRequest)
					})
					if err != nil {
						if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
							return resource.RetryableError(err)
						}
						return resource.RetryableError(err)
					}
					addDebug(unAssignPrivateIpAddressesRequest.GetActionName(), raw, unAssignPrivateIpAddressesRequest.RpcRequest, unAssignPrivateIpAddressesRequest)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), unAssignPrivateIpAddressesRequest.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}

			err := ecsService.WaitForPrivateIpsCountChanged(d.Id(), newIpsCount.(int))
			if err != nil {
				return WrapError(err)
			}

			d.SetPartial("private_ips_count")
		}
	}

	if err := setTags(client, TagResourceEni, d); err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceAliyunNetworkInterfaceRead(d, meta)
}

func resourceAliyunNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteNetworkInterfaceRequest()
	request.RegionId = client.RegionId
	request.NetworkInterfaceId = d.Id()

	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteNetworkInterface(request)
		})
		if err != nil {
			if IsExpectedErrors(err, NetworkInterfaceInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ecsService.WaitForNetworkInterface(d.Id(), Deleted, DefaultTimeoutMedium))
}
