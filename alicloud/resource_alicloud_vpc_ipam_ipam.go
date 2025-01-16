// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcIpamIpam() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpamIpamCreate,
		Read:   resourceAliCloudVpcIpamIpamRead,
		Update: resourceAliCloudVpcIpamIpamUpdate,
		Delete: resourceAliCloudVpcIpamIpamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipam_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operating_region_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"private_default_scope_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudVpcIpamIpamCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateIpam"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("ipam_name"); ok {
		request["IpamName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("operating_region_list"); ok {
		operatingRegionListMapsArray := v.(*schema.Set).List()
		request["OperatingRegionList"] = operatingRegionListMapsArray
	}

	if v, ok := d.GetOk("ipam_description"); ok {
		request["IpamDescription"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipam_ipam", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpamId"]))

	return resourceAliCloudVpcIpamIpamRead(d, meta)
}

func resourceAliCloudVpcIpamIpamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcIpamServiceV2 := VpcIpamServiceV2{client}

	objectRaw, err := vpcIpamServiceV2.DescribeVpcIpamIpam(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipam_ipam DescribeVpcIpamIpam Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["IpamDescription"] != nil {
		d.Set("ipam_description", objectRaw["IpamDescription"])
	}
	if objectRaw["IpamName"] != nil {
		d.Set("ipam_name", objectRaw["IpamName"])
	}
	if objectRaw["PrivateDefaultScopeId"] != nil {
		d.Set("private_default_scope_id", objectRaw["PrivateDefaultScopeId"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["IpamStatus"] != nil {
		d.Set("status", objectRaw["IpamStatus"])
	}

	operatingRegionList1Raw := make([]interface{}, 0)
	if objectRaw["OperatingRegionList"] != nil {
		operatingRegionList1Raw = objectRaw["OperatingRegionList"].([]interface{})
	}

	d.Set("operating_region_list", operatingRegionList1Raw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcIpamIpamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "UpdateIpam"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("ipam_name") {
		update = true
		request["IpamName"] = d.Get("ipam_name")
	}

	if d.HasChange("ipam_description") {
		update = true
		request["IpamDescription"] = d.Get("ipam_description")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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
	update = false
	action = "ChangeResourceGroup"
	conn, err = client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "IPAM"
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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

	if d.HasChange("operating_region_list") {
		oldEntry, newEntry := d.GetChange("operating_region_list")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "UpdateIpam"
			conn, err := client.NewVpcipamClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["IpamId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			removeOperatingRegionMapsArray := localData
			request["RemoveOperatingRegion"] = removeOperatingRegionMapsArray

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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

		if added.Len() > 0 {
			action := "UpdateIpam"
			conn, err := client.NewVpcipamClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["IpamId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			addOperatingRegionMapsArray := localData
			request["AddOperatingRegion"] = addOperatingRegionMapsArray

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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
	if d.HasChange("tags") {
		vpcIpamServiceV2 := VpcIpamServiceV2{client}
		if err := vpcIpamServiceV2.SetResourceTags(d, "IPAM"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudVpcIpamIpamRead(d, meta)
}

func resourceAliCloudVpcIpamIpamDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteIpam"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["IpamId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
