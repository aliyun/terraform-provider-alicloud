package ebs

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDiskReplicaGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiskReplicaGroupCreate,
		Read:   resourceDiskReplicaGroupRead,
		Update: resourceDiskReplicaGroupUpdate,
		Delete: resourceDiskReplicaGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_replica_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"group_name"},
				Computed:      true,
			},
			"enable_rtc": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"one_shot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"pair_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rpo": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reverse_replicate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"invalid", "creating", "created", "create_failed", "manual_syncing", "syncing", "normal", "stopping", "stopped", "stop_failed", "failovering", "failovered", "failover_failed", "reprotecting", "reprotect_failed", "deleting", "delete_failed", "deleted"}, false),
			},
			"tags": tags.TagsSchema(),
			"group_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'group_name' has been deprecated since provider version 1.245.0. New field 'disk_replica_group_name' instead.",
			},
		},
	}
}

func resourceDiskReplicaGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDiskReplicaGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["SourceZoneId"] = d.Get("source_zone_id")
	request["DestinationRegionId"] = d.Get("destination_region_id")
	request["DestinationZoneId"] = d.Get("destination_zone_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("rpo"); ok {
		request["RPO"] = v
	}
	if v, ok := d.GetOkExists("enable_rtc"); ok {
		request["EnableRtc"] = v
	}
	request["RegionId"] = d.Get("source_region_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("group_name"); ok || d.HasChange("group_name") {
		request["GroupName"] = v
	}

	if v, ok := d.GetOk("disk_replica_group_name"); ok {
		request["GroupName"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ebs_disk_replica_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ReplicaGroupId"]))

	ebsServiceV2 := NewEbsServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "Status", []string{"create_failed", "invalid"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDiskReplicaGroupUpdate(d, meta)
}

func resourceDiskReplicaGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := NewEbsServiceV2(client)

	objectRaw, err := ebsServiceV2.DescribeEbsDiskReplicaGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_disk_replica_group DescribeEbsDiskReplicaGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("destination_region_id", objectRaw["DestinationRegionId"])
	d.Set("destination_zone_id", objectRaw["DestinationZoneId"])
	d.Set("disk_replica_group_name", objectRaw["GroupName"])
	d.Set("enable_rtc", objectRaw["EnableRtc"])
	d.Set("rpo", objectRaw["RPO"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("source_region_id", objectRaw["SourceRegionId"])
	d.Set("source_zone_id", objectRaw["SourceZoneId"])
	d.Set("status", objectRaw["Status"])

	pairIdsRaw := make([]interface{}, 0)
	if objectRaw["PairIds"] != nil {
		pairIdsRaw = objectRaw["PairIds"].([]interface{})
	}

	d.Set("pair_ids", pairIdsRaw)

	objectRaw, err = ebsServiceV2.DescribeDiskReplicaGroupListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	d.Set("group_name", d.Get("disk_replica_group_name"))
	return nil
}

func resourceDiskReplicaGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		ebsServiceV2 := NewEbsServiceV2(client)
		object, err := ebsServiceV2.DescribeEbsDiskReplicaGroup(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "normal" {
				action := "StartDiskReplicaGroup"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				if v, ok := d.GetOkExists("one_shot"); ok {
					request["OneShot"] = v
				}
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				ebsServiceV2 := NewEbsServiceV2(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "#$.LastRecoverPoint", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "stopped" {
				action := "ReprotectDiskReplicaGroup"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				if v, ok := d.GetOkExists("reverse_replicate"); ok {
					request["ReverseReplicate"] = v
				}
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.InvalidStatus"}) || helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				ebsServiceV2 := NewEbsServiceV2(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"stopped"}, d.Timeout(schema.TimeoutUpdate), 0, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "failovered" {
				action := "FailoverDiskReplicaGroup"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"OperationDenied.InvalidStatus", "InternalError"}) || helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				ebsServiceV2 := NewEbsServiceV2(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"failovered"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "ModifyDiskReplicaGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ReplicaGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("group_name") {
		update = true
		request["GroupName"] = d.Get("group_name")
	}

	if !d.IsNewResource() && d.HasChange("disk_replica_group_name") {
		update = true
		request["GroupName"] = d.Get("disk_replica_group_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("enable_rtc") {
		update = true
		request["EnableRtc"] = d.Get("enable_rtc")
	}

	if !d.IsNewResource() && d.HasChange("rpo") {
		update = true
		request["RPO"] = d.Get("rpo")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "DiskReplicaGroup"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		ebsServiceV2 := NewEbsServiceV2(client)
		if err := ebsServiceV2.SetResourceTags(d, "DiskReplicaGroup"); err != nil {
			return errs.WrapError(err)
		}
	}
	if d.HasChange("pair_ids") {
		oldEntry, newEntry := d.GetChange("pair_ids")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			pairIds := removed.List()

			for _, item := range pairIds {
				action := "RemoveDiskReplicaPair"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return errs.WrapError(err)
					}
					request["ReplicaPairId"] = jsonPathResult
				}
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"OperationDenied.NoPairInGroup", "InternalError", "OperationDenied.PairStatusCannotRemoveFromGroup", "OperationDenied.GroupStatusCannotRemovePair"}) || helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}

			}
		}

		if added.Len() > 0 {
			pairIds := added.List()

			for _, item := range pairIds {
				action := "AddDiskReplicaPair"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return errs.WrapError(err)
					}
					request["ReplicaPairId"] = jsonPathResult
				}
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"OperationDenied.PairStatusCannotAddToGroup", "OperationDenied.GroupStatusCannotAddPair", "InternalError"}) || helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				ebsServiceV2 := NewEbsServiceV2(client)
				stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "#PairIds", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}

	}
	d.Partial(false)
	return resourceDiskReplicaGroupRead(d, meta)
}

func resourceDiskReplicaGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDiskReplicaGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ReplicaGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.GroupHasPair", "OperationDenied.InvalidStatus"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	ebsServiceV2 := NewEbsServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
