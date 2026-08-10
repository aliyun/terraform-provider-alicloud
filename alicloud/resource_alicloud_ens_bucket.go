package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEnsBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEnsBucketCreate,
		Read:   resourceAlicloudEnsBucketRead,
		Update: resourceAlicloudEnsBucketUpdate,
		Delete: resourceAlicloudEnsBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket_acl": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "private",
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
			},
			"logical_bucket_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"sink", "standard"}, false),
			},
			"ens_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dispatch_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"domestic", "oversea"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEnsBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "PutBucket"
	var response map[string]interface{}
	request := make(map[string]interface{})
	var err error
	request["BucketName"] = d.Get("bucket_name")
	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	if v, ok := d.GetOk("bucket_acl"); ok {
		request["BucketAcl"] = v
	}
	if v, ok := d.GetOk("logical_bucket_type"); ok {
		request["LogicalBucketType"] = v
	}
	if v, ok := d.GetOk("ens_region_id"); ok {
		request["EnsRegionId"] = v
	}
	if v, ok := d.GetOk("dispatch_scope"); ok {
		request["DispatchScope"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_bucket", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["BucketName"]))

	// After PutBucket, the bucket may not be queryable immediately. Poll
	// GetBucketInfo until it returns the bucket (NoSuchBucket during this
	// window means not ready yet, not deleted).
	ensServiceV2 := EnsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ensServiceV2.EnsBucketStateRefreshFunc(d.Id(), "", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEnsBucketRead(d, meta)
}

func resourceAlicloudEnsBucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}
	object, err := ensServiceV2.DescribeEnsBucket(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_bucket DescribeEnsBucket Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	// ENS GetBucketInfo is eventually-consistent across replicas after
	// UpdateBucketInfo: a stale replica may regress both Comment and
	// ModifyTime to pre-update values. The Update poll returns on the
	// first fresh read, but a subsequent Read (the Read at the end of
	// Update, a Refresh, or a testAccCheck-triggered refresh) can still
	// hit a stale replica and write the stale comment into state, which
	// the legacy plugin SDK tolerates and which then causes flaky
	// assertion failures. Compare the API response to the state baseline
	// and retry GetBucketInfo until the replica converges or a short
	// deadline expires, covering every Read path that has a baseline.
	stateComment := d.Get("comment").(string)
	stateModifyTime := d.Get("modify_time").(string)
	if ensBucketResponseStale(object, stateComment, stateModifyTime) {
		staleDeadline := time.Now().Add(30 * time.Second)
		for {
			time.Sleep(2 * time.Second)
			freshObject, e := ensServiceV2.DescribeEnsBucket(d.Id())
			if e != nil {
				// A transient error during stale-replica retry is
				// non-fatal: keep the last good response and let the
				// next Read reconcile. A real NotFound was already
				// handled by the initial DescribeEnsBucket above.
				log.Printf("[WARN] Resource alicloud_ens_bucket %s: GetBucketInfo error during stale-replica retry: %s; accepting previous response", d.Id(), e)
				break
			}
			object = freshObject
			if !ensBucketResponseStale(object, stateComment, stateModifyTime) {
				break
			}
			if time.Now().After(staleDeadline) {
				log.Printf("[WARN] Resource alicloud_ens_bucket %s: GetBucketInfo still returning stale replica after 30s; accepting current response", d.Id())
				break
			}
		}
	}
	if object["BucketName"] != nil {
		d.Set("bucket_name", object["BucketName"])
	}
	if object["Comment"] != nil {
		d.Set("comment", object["Comment"])
	}
	if object["BucketAcl"] != nil {
		d.Set("bucket_acl", object["BucketAcl"])
	}
	if object["LogicalBucketType"] != nil {
		d.Set("logical_bucket_type", object["LogicalBucketType"])
	}
	if object["CreateTime"] != nil {
		d.Set("create_time", object["CreateTime"])
	}
	if object["ModifyTime"] != nil {
		d.Set("modify_time", object["ModifyTime"])
	}
	return nil
}

// ensBucketResponseStale reports whether a GetBucketInfo response likely
// came from a stale ENS replica by comparing it to the state baseline.
//   - Comment regression: the state has a non-empty comment and the API
//     returns a different value (the stale replica serves the
//     pre-update comment).
//   - ModifyTime regression: the state has a non-empty modify_time and
//     the API ModifyTime is strictly older (a stale replica serves
//     older data with an older ModifyTime — a deterministic signal).
//
// When the state has no baseline (e.g. the first Read after Create or
// an Import, where comment and modify_time are empty), no comparison is
// performed and the response is treated as fresh.
func ensBucketResponseStale(object map[string]interface{}, stateComment, stateModifyTime string) bool {
	if stateComment == "" && stateModifyTime == "" {
		return false
	}
	if stateComment != "" {
		apiComment := ""
		if object["Comment"] != nil {
			apiComment = fmt.Sprint(object["Comment"])
		}
		if apiComment != stateComment {
			return true
		}
	}
	if stateModifyTime != "" {
		apiModifyTime := ""
		if object["ModifyTime"] != nil {
			apiModifyTime = fmt.Sprint(object["ModifyTime"])
		}
		if apiModifyTime != "" && ensBucketModifyTimeBefore(apiModifyTime, stateModifyTime) {
			return true
		}
	}
	return false
}

// ensBucketModifyTimeBefore reports whether ENS ModifyTime a is strictly
// before b. ENS ModifyTime is an ISO 8601 UTC timestamp (e.g.
// "2026-08-10T18:50:20.000Z"); time comparison is preferred and
// lexicographic comparison is the fallback for same-format timestamps
// when parsing fails.
func ensBucketModifyTimeBefore(a, b string) bool {
	ta, errA := time.Parse(time.RFC3339Nano, a)
	tb, errB := time.Parse(time.RFC3339Nano, b)
	if errA == nil && errB == nil {
		return ta.Before(tb)
	}
	return a < b
}

func resourceAlicloudEnsBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "UpdateBucketInfo"
	var response map[string]interface{}
	var request map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["BucketName"] = d.Id()
	if d.HasChange("comment") {
		request["Comment"] = d.Get("comment")
	}
	if len(request) > 0 {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		// ENS GetBucketInfo is eventually-consistent after UpdateBucketInfo —
		// the comment field may return the old value (or a transient null
		// response) for several seconds. Poll until the updated comment is
		// visible so that the subsequent Read and any refresh/import steps see
		// consistent state, avoiding flaky ImportStateVerify failures.
		if d.HasChange("comment") {
			expectedComment := fmt.Sprint(d.Get("comment"))
			ensServiceV2 := EnsServiceV2{client}
			propWait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(2*time.Minute, func() *resource.RetryError {
				object, e := ensServiceV2.DescribeEnsBucket(d.Id())
				if e != nil {
					// During propagation GetBucketInfo may transiently return
					// an empty (null) response surfaced as NotFound; the
					// bucket still exists after an update, so retry.
					if NeedRetry(e) || NotFoundError(e) {
						propWait()
						return resource.RetryableError(e)
					}
					return resource.NonRetryableError(e)
				}
				if fmt.Sprint(object["Comment"]) != expectedComment {
					propWait()
					return resource.RetryableError(fmt.Errorf("comment update for %s not yet propagated", d.Id()))
				}
				return nil
			})
			if err != nil {
				return WrapError(err)
			}
		}
	}
	return resourceAlicloudEnsBucketRead(d, meta)
}

func resourceAlicloudEnsBucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBucket"
	var response map[string]interface{}
	var request map[string]interface{}
	var err error
	query := make(map[string]interface{})
	query["BucketName"] = d.Id()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	// Poll GetBucketInfo until the bucket is confirmed gone. ENS
	// GetBucketInfo may return an empty (null) response — not a NoSuchBucket
	// error — while a bucket is still deleting; treating that null as
	// NotFound would make Delete return prematurely and leave a dangling
	// resource that CheckDestroy then trips over. Only a real NoSuchBucket
	// SDK error is accepted as definitive deletion; a null or data-bearing
	// response keeps polling until the delete timeout.
	getAction := "GetBucketInfo"
	wait = incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", getAction, query, request, true)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(getAction, response, query)
		wait()
		return resource.RetryableError(fmt.Errorf("bucket %s is still being deleted", d.Id()))
	})
	if err != nil && !NotFoundError(err) {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// ENS GetBucketInfo is eventually-consistent across replicas: the
	// poll above confirmed NoSuchBucket from one replica, but other
	// replicas may still serve stale bucket data for a short window
	// after DeleteBucket. CheckDestroy (which calls DescribeEnsBucket)
	// can hit any replica, so a single NoSuchBucket is not enough to
	// guarantee all replicas have converged. Require two consecutive
	// NoSuchBucket responses (which may hit different replicas) with a
	// short delay between them before returning, so CheckDestroy is
	// less likely to trip over a stale response. A data-bearing
	// response resets the counter and keeps waiting up to the deadline.
	convergenceDeadline := time.Now().Add(30 * time.Second)
	consecutiveNotFound := 0
	for {
		time.Sleep(5 * time.Second)
		resp, e := client.RpcPost("Ens", "2017-11-10", getAction, query, request, true)
		if e != nil {
			consecutiveNotFound++
			if consecutiveNotFound >= 2 {
				break
			}
			continue
		}
		consecutiveNotFound = 0
		addDebug(getAction, resp, query)
		if time.Now().After(convergenceDeadline) {
			log.Printf("[WARN] Resource alicloud_ens_bucket %s: GetBucketInfo still returning stale bucket data after 30s post-delete; proceeding", d.Id())
			break
		}
	}
	return nil
}
