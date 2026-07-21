// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcPrefixList() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcPrefixListCreate,
		Read:   resourceAlicloudVpcPrefixListRead,
		Update: resourceAlicloudVpcPrefixListUpdate,
		Delete: resourceAlicloudVpcPrefixListDelete,
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
			"entrys": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_.-]{2,256}$"), "The description of the cidr entry. It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`."),
						},
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"IPV4", "IPV6"}, false),
			},
			"max_entries": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"prefix_list_association": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"prefix_list_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_.-]{2,256}$"), "The description of the prefix list.It must be 2 to 256 characters in length and must start with a letter or Chinese, but cannot start with `http://` or `https://`."),
			},
			"prefix_list_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix_list_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_.-]{2,128}$"), "The name of the prefix list. The name must be 2 to 128 characters in length, and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-)."),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
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

func resourceAlicloudVpcPrefixListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpcPrefixList"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}

	if v, ok := d.GetOk("max_entries"); ok {
		request["MaxEntries"] = v
	}

	if v, ok := d.GetOk("prefix_list_description"); ok {
		request["PrefixListDescription"] = v
	}

	if v, ok := d.GetOk("prefix_list_name"); ok {
		request["PrefixListName"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("entrys"); ok {
		localData := v
		prefixListEntriesMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range localData.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Cidr"] = dataLoopTmp["cidr"]
			dataLoopMap["Description"] = dataLoopTmp["description"]
			prefixListEntriesMaps = append(prefixListEntriesMaps, dataLoopMap)
		}
		request["PrefixListEntries"] = prefixListEntriesMaps
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "SystemBusy", "IncorrectStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_prefix_list", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PrefixListId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcPrefixListStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcPrefixListUpdate(d, meta)
}

func resourceAlicloudVpcPrefixListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcPrefixList(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_prefix_list DescribeVpcPrefixList Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("ip_version", objectRaw["IpVersion"])
	d.Set("max_entries", objectRaw["MaxEntries"])
	d.Set("prefix_list_description", objectRaw["PrefixListDescription"])
	d.Set("prefix_list_name", objectRaw["PrefixListName"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("share_type", objectRaw["ShareType"])
	d.Set("status", objectRaw["Status"])
	d.Set("prefix_list_id", objectRaw["PrefixListId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = vpcServiceV2.DescribeGetVpcPrefixListEntries(d.Id())
	if err != nil {
		return WrapError(err)
	}

	prefixListEntry1Raw := objectRaw["PrefixListEntry"]
	entriesMaps := make([]map[string]interface{}, 0)
	if prefixListEntry1Raw != nil {
		for _, prefixListEntryChild1Raw := range prefixListEntry1Raw.([]interface{}) {
			entriesMap := make(map[string]interface{})
			prefixListEntryChild1Raw := prefixListEntryChild1Raw.(map[string]interface{})
			entriesMap["cidr"] = prefixListEntryChild1Raw["Cidr"]
			entriesMap["description"] = prefixListEntryChild1Raw["Description"]
			entriesMaps = append(entriesMaps, entriesMap)
		}
	}
	d.Set("entrys", entriesMaps)

	objectRaw, err = vpcServiceV2.DescribeGetVpcPrefixListAssociations(d.Id())
	if err != nil {
		return WrapError(err)
	}

	prefixListAssociation1Raw := objectRaw["PrefixListAssociation"]
	prefixListAssociationMaps := make([]map[string]interface{}, 0)
	if prefixListAssociation1Raw != nil {
		for _, prefixListAssociationChild1Raw := range prefixListAssociation1Raw.([]interface{}) {
			prefixListAssociationMap := make(map[string]interface{})
			prefixListAssociationChild1Raw := prefixListAssociationChild1Raw.(map[string]interface{})
			prefixListAssociationMap["owner_id"] = prefixListAssociationChild1Raw["OwnerId"]
			prefixListAssociationMap["prefix_list_id"] = prefixListAssociationChild1Raw["PrefixListId"]
			prefixListAssociationMap["reason"] = prefixListAssociationChild1Raw["Reason"]
			prefixListAssociationMap["region_id"] = prefixListAssociationChild1Raw["RegionId"]
			prefixListAssociationMap["resource_id"] = prefixListAssociationChild1Raw["ResourceId"]
			prefixListAssociationMap["resource_type"] = prefixListAssociationChild1Raw["ResourceType"]
			prefixListAssociationMap["resource_uid"] = prefixListAssociationChild1Raw["ResourceUid"]
			prefixListAssociationMap["status"] = prefixListAssociationChild1Raw["Status"]
			prefixListAssociationMaps = append(prefixListAssociationMaps, prefixListAssociationMap)
		}
	}
	d.Set("prefix_list_association", prefixListAssociationMaps)
	return nil
}

func resourceAlicloudVpcPrefixListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyVpcPrefixList"
	var err error
	request = make(map[string]interface{})

	request["PrefixListId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if !d.IsNewResource() && d.HasChange("prefix_list_name") {
		update = true
		if v, ok := d.GetOk("prefix_list_name"); ok {
			request["PrefixListName"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("prefix_list_description") {
		update = true
		if v, ok := d.GetOk("prefix_list_description"); ok {
			request["PrefixListDescription"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("max_entries") {
		update = true
		if v, ok := d.GetOk("max_entries"); ok {
			request["MaxEntries"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("prefix_list_name")
		d.SetPartial("prefix_list_description")
		d.SetPartial("max_entries")
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})

	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["NewResourceGroupId"] = v
		}
	}

	request["ResourceType"] = "PrefixList"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	update = false
	if !d.IsNewResource() && d.HasChange("entrys") {
		update = true
		oldEntry, newEntry := d.GetChange("entrys")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action = "ModifyVpcPrefixList"
			request = make(map[string]interface{})

			request["PrefixListId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)

			localData := removed.List()

			removePrefixListEntryMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Cidr"] = dataLoopTmp["cidr"]
				dataLoopMap["Description"] = dataLoopTmp["description"]
				removePrefixListEntryMaps = append(removePrefixListEntryMaps, dataLoopMap)
			}
			request["RemovePrefixListEntry"] = removePrefixListEntryMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.PrefixList", "IncorrectStatus", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if added.Len() > 0 {
			action = "ModifyVpcPrefixList"
			request = make(map[string]interface{})

			request["PrefixListId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)

			localData := added.List()

			addPrefixListEntryMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Cidr"] = dataLoopTmp["cidr"]
				dataLoopMap["Description"] = dataLoopTmp["description"]
				addPrefixListEntryMaps = append(addPrefixListEntryMaps, dataLoopMap)
			}
			request["AddPrefixListEntry"] = addPrefixListEntryMaps

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectStatus.PrefixList", "IncorrectStatus", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

	}
	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "PrefixList"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAlicloudVpcPrefixListRead(d, meta)
}

func resourceAlicloudVpcPrefixListDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteVpcPrefixList"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["PrefixListId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "SystemBusy", "DependencyViolation.ShareResource", "IncorrectStatus.PrefixList", "IncorrectStatus.SystemPrefixList", "IncorrectStatus", "OperationFailed.LastTokenProcessing", "LastTokenProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcPrefixListStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
