// Package ebs. This file is generated automatically. Please do not modify it manually, thank you!
package ebs

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDiskReplicaPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiskReplicaPairCreate,
		Read:   resourceDiskReplicaPairRead,
		Update: resourceDiskReplicaPairUpdate,
		Delete: resourceDiskReplicaPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.IntInSlice([]int{0, 10240, 20480, 51200, 102400}),
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_replica_pair_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"pair_name"},
				Computed:      true,
			},
			"one_shot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == new || convertEbsDiskReplicaPairChargeTypeRequest(old) == new || convertEbsDiskReplicaPairChargeTypeRequest(new) == old
				},
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 60),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rpo": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
			"reverse_replicate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"pair_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'pair_name' has been deprecated since provider version 1.245.0. New field 'disk_replica_pair_name' instead.",
			},
		},
	}
}

func resourceDiskReplicaPairCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDiskReplicaPair"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEbsDiskReplicaPairChargeTypeRequest(v.(string))
	}
	request["DestinationRegionId"] = d.Get("destination_region_id")
	request["DestinationZoneId"] = d.Get("destination_zone_id")
	request["SourceZoneId"] = d.Get("source_zone_id")
	request["DestinationDiskId"] = d.Get("destination_disk_id")
	request["DiskId"] = d.Get("disk_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("period"); ok && v.(int) > 0 {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOkExists("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOkExists("rpo"); ok {
		request["RPO"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("pair_name"); ok || d.HasChange("pair_name") {
		request["PairName"] = v
	}

	if v, ok := d.GetOk("disk_replica_pair_name"); ok {
		request["PairName"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ebs_disk_replica_pair", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ReplicaPairId"]))

	ebsServiceV2 := NewEbsServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"created", "normal", "stopped", "syncing", "manual_syncing"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ebsServiceV2.EbsDiskReplicaPairStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDiskReplicaPairUpdate(d, meta)
}

func resourceDiskReplicaPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := NewEbsServiceV2(client)

	objectRaw, err := ebsServiceV2.DescribeEbsDiskReplicaPair(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_disk_replica_pair DescribeEbsDiskReplicaPair Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("destination_disk_id", objectRaw["DestinationDiskId"])
	d.Set("destination_region_id", objectRaw["DestinationRegion"])
	d.Set("destination_zone_id", objectRaw["DestinationZoneId"])
	d.Set("disk_id", objectRaw["SourceDiskId"])
	d.Set("disk_replica_pair_name", objectRaw["PairName"])
	d.Set("payment_type", convertEbsDiskReplicaPairReplicaPairsChargeTypeResponse(objectRaw["ChargeType"]))
	d.Set("rpo", objectRaw["RPO"])
	d.Set("region_id", objectRaw["SourceRegion"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("source_zone_id", objectRaw["SourceZoneId"])
	d.Set("status", objectRaw["Status"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	d.Set("pair_name", d.Get("disk_replica_pair_name"))
	return nil
}

func resourceDiskReplicaPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("status") {
		ebsServiceV2 := NewEbsServiceV2(client)
		object, err := ebsServiceV2.DescribeEbsDiskReplicaPair(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "normal" {
				action := "StartDiskReplicaPair"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaPairId"] = d.Id()
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
				stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ebsServiceV2.EbsDiskReplicaPairStateRefreshFunc(d.Id(), "#$.LastRecoverPoint", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "stopped" {
				action := "ReprotectDiskReplicaPair"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaPairId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = helper.BuildClientToken(action)
				if v, ok := d.GetOkExists("reverse_replicate"); ok {
					request["ReverseReplicate"] = v
				}
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
				stateConf := helper.BuildStateConf([]string{}, []string{"stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ebsServiceV2.EbsDiskReplicaPairStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "failovered" {
				action := "FailoverDiskReplicaPair"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["ReplicaPairId"] = d.Id()
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
				stateConf := helper.BuildStateConf([]string{}, []string{"failovered"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ebsServiceV2.EbsDiskReplicaPairStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "ModifyDiskReplicaPair"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ReplicaPairId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("pair_name") {
		update = true
		request["PairName"] = d.Get("pair_name")
	}

	if !d.IsNewResource() && d.HasChange("disk_replica_pair_name") {
		update = true
		request["PairName"] = d.Get("disk_replica_pair_name")
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
	request["ResourceType"] = "DiskReplicaPair"
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
		if err := ebsServiceV2.SetResourceTags(d, "DiskReplicaPair"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceDiskReplicaPairRead(d, meta)
}

func resourceDiskReplicaPairDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDiskReplicaPair"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ReplicaPairId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	ebsServiceV2 := NewEbsServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ebsServiceV2.EbsDiskReplicaPairStateRefreshFunc(d.Id(), "$.ReplicaPairId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

func convertEbsDiskReplicaPairReplicaPairsChargeTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "PREPAY":
		return "Subscription"
	case "POSTPAY":
		return "PayAsYouGo"
	}
	return source
}
func convertEbsDiskReplicaPairChargeTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "PREPAY"
	case "PayAsYouGo":
		return "POSTPAY"
	}
	return source
}
