// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEcdDesktopGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcdDesktopGroupCreate,
		Read:   resourceAliCloudEcdDesktopGroupRead,
		Update: resourceAliCloudEcdDesktopGroupUpdate,
		Delete: resourceAliCloudEcdDesktopGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_auto_setup": {
				Type:      schema.TypeInt,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"allow_buffer_count": {
				Type:      schema.TypeInt,
				Optional:  true,
				Sensitive: true,
			},
			"bundle_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_disk_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_disk_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"desktop_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"directory_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_user_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_count": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"gpu_spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keep_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_desktops_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_desktops_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"office_site_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"office_site_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"own_bundle_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"res_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scale_strategy_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"system_disk_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEcdDesktopGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDesktopGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["ChargeType"] = "PostPaid"

	if v, ok := d.GetOkExists("allow_auto_setup"); ok {
		request["AllowAutoSetup"] = v
	}
	request["OfficeSiteId"] = d.Get("office_site_id")
	if v, ok := d.GetOk("end_user_ids"); ok {
		endUserIdsMapsArray := convertToInterfaceArray(v)

		request["EndUserIds"] = endUserIdsMapsArray
	}

	if v, ok := d.GetOkExists("max_desktops_count"); ok {
		request["MaxDesktopsCount"] = v
	}
	request["BundleId"] = d.Get("bundle_id")
	if v, ok := d.GetOkExists("keep_duration"); ok {
		request["KeepDuration"] = v
	}
	if v, ok := d.GetOk("directory_id"); ok {
		request["DirectoryId"] = v
	}
	if v, ok := d.GetOk("desktop_group_name"); ok {
		request["DesktopGroupName"] = v
	}
	if v, ok := d.GetOk("comments"); ok {
		request["Comments"] = v
	}
	if v, ok := d.GetOkExists("allow_buffer_count"); ok {
		request["AllowBufferCount"] = v
	}
	if v, ok := d.GetOkExists("min_desktops_count"); ok {
		request["MinDesktopsCount"] = v
	}
	if v, ok := d.GetOk("scale_strategy_id"); ok {
		request["ScaleStrategyId"] = v
	}
	request["PolicyGroupId"] = d.Get("policy_group_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_desktop_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DesktopGroupId"]))

	return resourceAliCloudEcdDesktopGroupUpdate(d, meta)
}

func resourceAliCloudEcdDesktopGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdServiceV2 := EcdServiceV2{client}

	objectRaw, err := ecdServiceV2.DescribeEcdDesktopGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_desktop_group DescribeEcdDesktopGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("allow_auto_setup", objectRaw["AllowAutoSetup"])
	d.Set("allow_buffer_count", objectRaw["AllowBufferCount"])
	d.Set("bundle_id", objectRaw["OwnBundleId"])
	d.Set("comments", objectRaw["Comments"])
	d.Set("cpu", objectRaw["Cpu"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("creator", objectRaw["Creator"])
	d.Set("data_disk_category", objectRaw["DataDiskCategory"])
	d.Set("data_disk_size", objectRaw["DataDiskSize"])
	d.Set("desktop_group_name", objectRaw["DesktopGroupName"])
	d.Set("directory_id", objectRaw["DirectoryId"])
	d.Set("directory_type", objectRaw["DirectoryType"])
	d.Set("expired_time", objectRaw["ExpiredTime"])
	d.Set("gpu_count", objectRaw["GpuCount"])
	d.Set("gpu_spec", objectRaw["GpuSpec"])
	d.Set("keep_duration", objectRaw["KeepDuration"])
	d.Set("max_desktops_count", objectRaw["MaxDesktopsCount"])
	d.Set("memory", objectRaw["Memory"])
	d.Set("min_desktops_count", objectRaw["MinDesktopsCount"])
	d.Set("office_site_id", objectRaw["OfficeSiteId"])
	d.Set("office_site_name", objectRaw["OfficeSiteName"])
	d.Set("office_site_type", objectRaw["OfficeSiteType"])
	d.Set("own_bundle_name", objectRaw["OwnBundleName"])
	d.Set("pay_type", objectRaw["PayType"])
	d.Set("policy_group_id", objectRaw["PolicyGroupId"])
	d.Set("policy_group_name", objectRaw["PolicyGroupName"])
	d.Set("res_type", objectRaw["ResType"])
	d.Set("system_disk_category", objectRaw["SystemDiskCategory"])
	d.Set("system_disk_size", objectRaw["SystemDiskSize"])

	usersObject, err := ecdServiceV2.DescribeDesktopGroupDescribeUsersInGroup(d.Id())
	if err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
	} else {
		d.Set("end_user_ids", usersObject["EndUserIds"])
	}

	return nil
}

func resourceAliCloudEcdDesktopGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyDesktopGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DesktopGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("allow_auto_setup") {
		update = true
		request["AllowAutoSetup"] = d.Get("allow_auto_setup")
	}

	if !d.IsNewResource() && d.HasChange("min_desktops_count") {
		update = true
		request["MinDesktopsCount"] = d.Get("min_desktops_count")
	}

	if !d.IsNewResource() && d.HasChange("scale_strategy_id") {
		update = true
		request["ScaleStrategyId"] = d.Get("scale_strategy_id")
	}

	if !d.IsNewResource() && d.HasChange("max_desktops_count") {
		update = true
		request["MaxDesktopsCount"] = d.Get("max_desktops_count")
	}

	if !d.IsNewResource() && d.HasChange("bundle_id") {
		update = true
	}
	request["OwnBundleId"] = d.Get("bundle_id")
	if !d.IsNewResource() && d.HasChange("keep_duration") {
		update = true
		request["KeepDuration"] = d.Get("keep_duration")
	}

	if !d.IsNewResource() && d.HasChange("desktop_group_name") {
		update = true
		request["DesktopGroupName"] = d.Get("desktop_group_name")
	}

	if !d.IsNewResource() && d.HasChange("comments") {
		update = true
		request["Comments"] = d.Get("comments")
	}

	if !d.IsNewResource() && d.HasChange("allow_buffer_count") {
		update = true
		request["AllowBufferCount"] = d.Get("allow_buffer_count")
	}

	if !d.IsNewResource() && d.HasChange("policy_group_id") {
		update = true
	}
	request["PolicyGroupId"] = d.Get("policy_group_id")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if !d.IsNewResource() && d.HasChange("end_user_ids") {
		var err error
		oldEntry, newEntry := d.GetChange("end_user_ids")
		oldUsers := make(map[string]interface{})
		for _, v := range oldEntry.([]interface{}) {
			oldUsers[fmt.Sprint(v)] = v
		}
		newUsers := make(map[string]interface{})
		for _, v := range newEntry.([]interface{}) {
			newUsers[fmt.Sprint(v)] = v
		}
		removed := make([]interface{}, 0)
		for _, v := range oldEntry.([]interface{}) {
			if _, ok := newUsers[fmt.Sprint(v)]; !ok {
				removed = append(removed, v)
			}
		}
		added := make([]interface{}, 0)
		for _, v := range newEntry.([]interface{}) {
			if _, ok := oldUsers[fmt.Sprint(v)]; !ok {
				added = append(added, v)
			}
		}

		if len(removed) > 0 {
			action := "RemoveUserFromDesktopGroup"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["DesktopGroupId"] = d.Id()
			request["RegionId"] = client.RegionId
			endUserIdsMapsArray := removed
			request["EndUserIds"] = endUserIdsMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action := "AddUserToDesktopGroup"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["DesktopGroupId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			endUserIdsMapsArray := added
			request["EndUserIds"] = endUserIdsMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}
	return resourceAliCloudEcdDesktopGroupRead(d, meta)
}

func resourceAliCloudEcdDesktopGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	ecdServiceV2 := EcdServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error

	// DeleteDesktopGroup requires that no end user is authorized to the
	// desktop group, so remove all authorized end users before deletion.
	usersObject, err := ecdServiceV2.DescribeDesktopGroupDescribeUsersInGroup(d.Id())
	if err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
	} else if endUserIds, ok := usersObject["EndUserIds"].([]interface{}); ok && len(endUserIds) > 0 {
		action := "RemoveUserFromDesktopGroup"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["DesktopGroupId"] = d.Id()
		request["RegionId"] = client.RegionId
		request["EndUserIds"] = endUserIds

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)
			if err != nil {
				if NeedRetry(err) || strings.Contains(err.Error(), "NotAllowOperation") {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	action := "DeleteDesktopGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DesktopGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, query, request, true)
		if err != nil {
			// The desktop group cannot be deleted while it is in a transient
			// state (e.g. desktops are being created or released in the group).
			if NeedRetry(err) || strings.Contains(err.Error(), "NotAllowOperation") {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	// The desktops in the desktop group are released asynchronously after the
	// desktop group is deleted, and they keep referencing resources such as the
	// policy group until the release completes. Deleting those resources in the
	// meantime fails with DependencyViolation, so wait until DescribeDesktopsInGroup
	// no longer lists any desktop (releasing desktops are listed with status
	// Deleted until they are purged).
	groupId := d.Id()
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		listAction := "DescribeDesktopsInGroup"
		nextToken := ""
		remaining := 0
		for {
			listRequest := make(map[string]interface{})
			listQuery := make(map[string]interface{})
			listRequest["RegionId"] = client.RegionId
			listRequest["DesktopGroupId"] = groupId
			listRequest["MaxResults"] = PageSizeLarge
			if nextToken != "" {
				listRequest["NextToken"] = nextToken
			}
			listResponse, listErr := client.RpcPost("ecd", "2020-09-30", listAction, listQuery, listRequest, true)
			if listErr != nil {
				if NotFoundError(listErr) {
					return nil
				}
				if NeedRetry(listErr) {
					return resource.RetryableError(listErr)
				}
				return resource.NonRetryableError(listErr)
			}
			// Desktops that are being released stay listed with status Deleted
			// for a long time, but they no longer block the deletion of the
			// resources they reference; only the desktops that are not in the
			// Deleted status are counted here.
			for _, key := range []string{"$.PostPaidDesktops[*]", "$.PaidDesktops[*]"} {
				if v, err := jsonpath.Get(key, listResponse); err == nil {
					if desktops, ok := v.([]interface{}); ok {
						for _, desktop := range desktops {
							if desktopMap, ok := desktop.(map[string]interface{}); ok && fmt.Sprint(desktopMap["DesktopStatus"]) != "Deleted" {
								remaining++
							}
						}
					}
				}
			}
			if token, ok := listResponse["NextToken"].(string); ok && token != "" {
				nextToken = token
			} else {
				break
			}
		}
		if remaining > 0 {
			return resource.RetryableError(fmt.Errorf("there are still %d desktops being released in the desktop group %s", remaining, groupId))
		}
		return nil
	})
	if err != nil {
		return WrapError(err)
	}

	return nil
}
