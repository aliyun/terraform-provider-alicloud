package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsEniSgAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsEniSgAttachmentCreate,
		Read:   resourceAliCloudEcsEniSgAttachmentRead,
		Delete: resourceAliCloudEcsEniSgAttachmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attach_security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"original_security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effective_security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudEcsEniSgAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	networkInterfaceId := d.Get("network_interface_id").(string)
	object, err := ecsService.DescribeEcsNetworkInterface(networkInterfaceId)
	if err != nil {
		return WrapError(err)
	}

	originalSgIds := flattenEniSecurityGroupIds(object)
	attachSgIds := expandStringList(d.Get("attach_security_group_ids").([]interface{}))
	targetSgIds := mergeSecurityGroupsWithAttach(attachSgIds, originalSgIds)

	if err := modifyEniSecurityGroups(client, d.Timeout(schema.TimeoutCreate), networkInterfaceId, targetSgIds); err != nil {
		return WrapError(err)
	}
	if err := waitForEniSecurityGroups(ecsService, d.Timeout(schema.TimeoutCreate), networkInterfaceId, targetSgIds); err != nil {
		return WrapError(err)
	}

	d.SetId(networkInterfaceId)
	d.Set("original_security_group_ids", originalSgIds)
	d.Set("effective_security_group_ids", targetSgIds)

	return resourceAliCloudEcsEniSgAttachmentRead(d, meta)
}

func resourceAliCloudEcsEniSgAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsNetworkInterface(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_eni_sg_attachment ecsService.DescribeEcsNetworkInterface Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("network_interface_id", d.Id())
	d.Set("effective_security_group_ids", flattenEniSecurityGroupIds(object))

	return nil
}

func resourceAliCloudEcsEniSgAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	networkInterfaceId := d.Get("network_interface_id").(string)
	originalSgIds := expandStringList(d.Get("original_security_group_ids").([]interface{}))
	if len(originalSgIds) == 0 {
		return nil
	}

	if err := modifyEniSecurityGroups(client, d.Timeout(schema.TimeoutDelete), networkInterfaceId, originalSgIds); err != nil {
		return WrapError(err)
	}
	if err := waitForEniSecurityGroups(ecsService, d.Timeout(schema.TimeoutDelete), networkInterfaceId, originalSgIds); err != nil {
		return WrapError(err)
	}

	return nil
}

func flattenEniSecurityGroupIds(object map[string]interface{}) []string {
	result := make([]string, 0)
	ids, ok := object["SecurityGroupIds"].(map[string]interface{})
	if !ok {
		return result
	}
	sgs, ok := ids["SecurityGroupId"].([]interface{})
	if !ok {
		return result
	}
	for _, sg := range sgs {
		result = append(result, fmt.Sprint(sg))
	}
	return result
}

func mergeSecurityGroupsWithAttach(attach []string, existing []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(attach)+len(existing))
	for _, sg := range attach {
		if _, ok := seen[sg]; ok {
			continue
		}
		seen[sg] = struct{}{}
		result = append(result, sg)
	}
	for _, sg := range existing {
		if _, ok := seen[sg]; ok {
			continue
		}
		seen[sg] = struct{}{}
		result = append(result, sg)
	}
	return result
}

func modifyEniSecurityGroups(client *connectivity.AliyunClient, timeout time.Duration, networkInterfaceId string, securityGroupIds []string) error {
	action := "ModifyNetworkInterfaceAttribute"
	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"NetworkInterfaceId": networkInterfaceId,
		"SecurityGroupId":    stringSliceToInterfaceSlice(securityGroupIds),
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(timeout, func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationConflict", "InvalidOperation.InvalidEniState", "ServiceUnavailable"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, networkInterfaceId, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func waitForEniSecurityGroups(ecsService EcsService, timeout time.Duration, networkInterfaceId string, expected []string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Applied"},
		Refresh: func() (interface{}, string, error) {
			object, err := ecsService.DescribeEcsNetworkInterface(networkInterfaceId)
			if err != nil {
				return nil, "", err
			}
			current := flattenEniSecurityGroupIds(object)
			if equalStringSet(current, expected) {
				return object, "Applied", nil
			}
			return object, "Pending", nil
		},
		Timeout:      timeout,
		MinTimeout:   3 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForState()
	return err
}

func equalStringSet(a []string, b []string) bool {
	uniqA := uniqueStrings(a)
	uniqB := uniqueStrings(b)
	if len(uniqA) != len(uniqB) {
		return false
	}
	set := make(map[string]struct{}, len(uniqA))
	for _, item := range uniqA {
		set[item] = struct{}{}
	}
	for _, item := range uniqB {
		if _, ok := set[item]; !ok {
			return false
		}
	}
	return true
}

func uniqueStrings(items []string) []string {
	result := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func stringSliceToInterfaceSlice(v []string) []interface{} {
	result := make([]interface{}, 0, len(v))
	for _, item := range v {
		result = append(result, item)
	}
	return result
}
