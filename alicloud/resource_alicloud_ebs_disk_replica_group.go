package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEbsDiskReplicaGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEbsDiskReplicaGroupCreate,
		Read:   resourceAliCloudEbsDiskReplicaGroupRead,
		Update: resourceAliCloudEbsDiskReplicaGroupUpdate,
		Delete: resourceAliCloudEbsDiskReplicaGroupDelete,
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
				ForceNew: true,
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
				ValidateFunc: StringInSlice([]string{"invalid", "creating", "created", "create_failed", "manual_syncing", "syncing", "normal", "stopping", "stopped", "stop_failed", "failovering", "failovered", "failover_failed", "reprotecting", "reprotect_failed", "deleting", "delete_failed", "deleted"}, false),
			},
			"tags": tagsSchema(),
			"group_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'group_name' has been deprecated since provider version 1.245.0. New field 'disk_replica_group_name' instead.",
			},
		},
	}
}

func resourceAliCloudEbsDiskReplicaGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDiskReplicaGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["SourceZoneId"] = d.Get("source_zone_id")
	request["DestinationRegionId"] = d.Get("destination_region_id")
	request["DestinationZoneId"] = d.Get("destination_zone_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("rpo"); ok {
		request["RPO"] = v
	}
	request["RegionId"] = d.Get("source_region_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("group_name"); ok || d.HasChange("group_name") {
		request["GroupName"] = v
	}

	if v, ok := d.GetOk("disk_replica_group_name"); ok {
		request["GroupName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ebs_disk_replica_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ReplicaGroupId"]))

	ebsServiceV2 := EbsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "Status", []string{"create_failed", "invalid"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEbsDiskReplicaGroupUpdate(d, meta)
}

func resourceAliCloudEbsDiskReplicaGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := EbsServiceV2{client}

	objectRaw, err := ebsServiceV2.DescribeEbsDiskReplicaGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_disk_replica_group DescribeEbsDiskReplicaGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("destination_region_id", objectRaw["DestinationRegionId"])
	d.Set("destination_zone_id", objectRaw["DestinationZoneId"])
	d.Set("disk_replica_group_name", objectRaw["GroupName"])
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
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("group_name", d.Get("disk_replica_group_name"))
	return nil
}

func resourceAliCloudEbsDiskReplicaGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		ebsServiceV2 := EbsServiceV2{client}
		object, err := ebsServiceV2.DescribeEbsDiskReplicaGroup(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "normal" {
				action := "StartDiskReplicaGroup"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("one_shot"); ok {
					request["OneShot"] = v
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
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
				ebsServiceV2 := EbsServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "#$.LastRecoverPoint", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "stopped" {
				action := "ReprotectDiskReplicaGroup"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				if v, ok := d.GetOkExists("reverse_replicate"); ok {
					request["ReverseReplicate"] = v
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"InternalError", "OperationDenied.InvalidStatus"}) || NeedRetry(err) {
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
				ebsServiceV2 := EbsServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"stopped"}, d.Timeout(schema.TimeoutUpdate), 0, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "failovered" {
				action := "FailoverDiskReplicaGroup"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationDenied.InvalidStatus", "InternalError"}) || NeedRetry(err) {
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
				ebsServiceV2 := EbsServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"failovered"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
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
	request["ClientToken"] = buildClientToken(action)
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

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
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
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "DiskReplicaGroup"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
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

	if d.HasChange("tags") {
		ebsServiceV2 := EbsServiceV2{client}
		if err := ebsServiceV2.SetResourceTags(d, "DiskReplicaGroup"); err != nil {
			return WrapError(err)
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
				request["ClientToken"] = buildClientToken(action)
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["ReplicaPairId"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationDenied.NoPairInGroup", "InternalError", "OperationDenied.PairStatusCannotRemoveFromGroup", "OperationDenied.GroupStatusCannotRemovePair"}) || NeedRetry(err) {
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

		if added.Len() > 0 {
			pairIds := added.List()

			for _, item := range pairIds {
				action := "AddDiskReplicaPair"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaGroupId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				if v, ok := item.(string); ok {
					jsonPathResult, err := jsonpath.Get("$", v)
					if err != nil {
						return WrapError(err)
					}
					request["ReplicaPairId"] = jsonPathResult
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationDenied.PairStatusCannotAddToGroup", "OperationDenied.GroupStatusCannotAddPair", "InternalError"}) || NeedRetry(err) {
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
				ebsServiceV2 := EbsServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "#PairIds", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}

	}
	d.Partial(false)
	return resourceAliCloudEbsDiskReplicaGroupRead(d, meta)
}

func resourceAliCloudEbsDiskReplicaGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDiskReplicaGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ReplicaGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "OperationDenied.GroupHasPair", "OperationDenied.InvalidStatus"}) || NeedRetry(err) {
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

	ebsServiceV2 := EbsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ebsServiceV2.EbsDiskReplicaGroupStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
