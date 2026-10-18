// Package alicloud. OSS (Object Storage Service) sweeper gap-fill.
//
// This file is part of the repo-wide campaign to fill missing inline sweepers.
// OSS is special: the provider's OSS resources use TWO client paths:
//
//  1. OSS SDK (github.com/aliyun/aliyun-oss-go-sdk/oss) via
//     connectivity.AliyunClient.WithOssClient(...) — used by alicloud_oss_bucket
//     and alicloud_oss_bucket_object (ListBuckets, ListObjects, DeleteBucket,
//     DeleteObject).
//
//  2. The OpenAPI-XML path connectivity.AliyunClient.Do("Oss",
//     helper.XmlParam(<method>, "2019-05-17", <opName>, <path>), query, body,
//     headers, hostMap, autoRetry) — used by EVERY alicloud_oss_bucket_* and
//     alicloud_oss_access_point sub-resource. The XML response body is decoded
//     to a map[string]interface{} and navigated with PaesslerAG/jsonpath. The
//     hostMap["bucket"] key scopes a call to a single bucket; an EMPTY hostMap
//     hits the account-level endpoint (see alicloud_oss_account_public_access_block).
//
// Existing inline OSS sweeper found (NOT duplicated here):
//   - alicloud_oss_bucket  — resource_alicloud_oss_bucket_test.go::testSweepOSSBuckets
//     Lists all buckets via the OSS SDK, gates by region (Location suffix) and
//     test-name prefix (or sweepAll), empties each bucket's objects first, then
//     DeleteBucket. This is the PRIMARY OSS cleanup unit: deleting a bucket
//     transitively removes all of its bucket-bound sub-configurations.
//
// Gap sweeper written here:
//   - alicloud_oss_access_point
//     Access points are the ONLY OSS sub-resource that is enumerable at account
//     scope (ListAccessPoints, doc 2391678) AND has a per-item DeleteAccessPoint.
//     They are regional resources bound to a bucket; the resource id is
//     "<bucket>:<access_point_name>". The sweeper enumerates in-region gated
//     buckets (OSS SDK ListBuckets, same region/name gate as the bucket
//     sweeper) and for each bucket calls bucket-scoped ListAccessPoints, then
//     DeleteAccessPoint(hostMap["bucket"], x-oss-access-point-name). Per-bucket
//     enumeration is used INSTEAD of account-level ListAccessPoints so that
//     region-scoping is guaranteed by construction (red-line: never clean
//     resources outside the sweep region) — every bucket walked is confirmed
//     in-region by its Location field.
//
// Intentionally DEFERRED (no sweep written):
//
//	Bucket-bound single-config sub-resources — each is a single Get/Put/Delete
//	configuration ON one bucket (DescribeOss..., no List API). DeleteBucket
//	(doc 31973) requires the bucket be empty of *objects/parts*; bucket-level
//	configs are not "content" and are removed transitively when the existing
//	alicloud_oss_bucket sweeper deletes the bucket. Writing per-bucket Get-only
//	enumeration is impossible (no List) and unnecessary:
//	  - alicloud_oss_bucket_acl
//	  - alicloud_oss_bucket_cors
//	  - alicloud_oss_bucket_logging
//	  - alicloud_oss_bucket_policy
//	  - alicloud_oss_bucket_referer
//	  - alicloud_oss_bucket_versioning
//	  - alicloud_oss_bucket_website
//	  - alicloud_oss_bucket_server_side_encryption
//	  - alicloud_oss_bucket_request_payment
//	  - alicloud_oss_bucket_transfer_acceleration
//	  - alicloud_oss_bucket_access_monitor
//	  - alicloud_oss_bucket_https_config
//	  - alicloud_oss_bucket_public_access_block
//	  - alicloud_oss_bucket_meta_query
//	  - alicloud_oss_bucket_user_defined_log_fields
//	  - alicloud_oss_bucket_response_header
//	  - alicloud_oss_bucket_overwrite_config
//	  - alicloud_oss_bucket_archive_direct_read
//	  - alicloud_oss_bucket_data_redundancy_transition  (per-bucket transition task; no account-level list; transitive)
//
//	Bucket-scoped collection sub-resources — have a per-bucket List + per-item
//	Delete, but are removed transitively when the bucket is deleted. Nested
//	per-bucket enumeration (list all buckets → per-bucket list collection →
//	delete each) is redundant with the existing bucket sweeper's bucket
//	deletion, so it is not written:
//	  - alicloud_oss_bucket_inventory   (ListBucketInventory + DeleteBucketInventory)
//	  - alicloud_oss_bucket_style       (PutStyle/DeleteStyle; no account-level ListStyle)
//	  - alicloud_oss_bucket_cname       (ListCname + DeleteCname)
//	  - alicloud_oss_bucket_replication (Get/DeleteBucketReplication; whole-bucket config)
//
//	WORM / retention-locked (red-line: do NOT force-delete):
//	  - alicloud_oss_bucket_worm                — Delete is a no-op once the
//	    retention policy is "Locked" (the provider itself logs
//	    "Cannot destroy... resources may remain" and returns nil); AbortBucketWorm
//	    only works on an InProgress/Unlocked worm. A locked WORM cannot be
//	    force-deleted.
//	  - alicloud_oss_bucket_object_worm_configuration — Delete is a no-op
//	    (provider logs "Cannot destroy... resources may remain"); object WORM
//	    retention cannot be force-released.
//
//	No-op / transient (Delete is a no-op, no lingering state):
//	  - alicloud_oss_bucket_cname_token — Delete is a no-op; tokens are
//	    auto-managed by OSS for cname verification and hold no cleanable state.
//
//	Account-level singleton (no List API):
//	  - alicloud_oss_account_public_access_block — exactly one config per
//	    account (id = accountId; Get+Put+Delete, no List). Account-wide; no
//	    name to gate on. Not swept to avoid mutating account-wide state.
//
//	Handled by the existing bucket sweeper directly:
//	  - alicloud_oss_bucket_object — the alicloud_oss_bucket sweeper empties
//	    objects (ListObjects + DeleteObject) from each gated bucket BEFORE
//	    DeleteBucket; sweeping objects individually is impractical (too many)
//	    and redundant.
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// ossSweepPrefixes are the standard test-resource name prefixes. OSS bucket
// names (and access-point names) are globally-unique DNS-style strings, so the
// name is the natural gate. When sweepAll() is true (clean-all on the test
// account) the gate is bypassed for in-region resources.
var ossSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// ossSweepNameOk is the fleet-wipe gate shared by OSS sweepers. It mirrors the
// repo convention: protectedSweepVpcName is always active (NODELETING / 21
// markers), sweepAll() short-circuits to true (clean-all), otherwise the name
// must begin with one of ossSweepPrefixes (case-insensitive).
func ossSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	for _, prefix := range ossSweepPrefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// ossCoerceList normalizes a decoded XML repeated element into a slice of maps.
// The OpenAPI-SDK XML decoder may surface a single repeated child either as a
// bare map[string]interface{} or as a []interface{}; handle both so the sweeper
// never silently skips a lone leftover.
func ossCoerceList(v interface{}) []map[string]interface{} {
	out := []map[string]interface{}{}
	switch t := v.(type) {
	case []interface{}:
		for _, x := range t {
			if m, ok := x.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
	case map[string]interface{}:
		out = append(out, t)
	case []map[string]interface{}:
		out = append(out, t...)
	}
	return out
}

// ossXmlCall is the retry wrapper around the OSS OpenAPI-XML client.Do path
// (product "Oss", version "2019-05-17"), mirroring OssServiceV2's Describe*
// helpers. autoRetry is passed through (true for GET/list, false for DELETE —
// matching the provider's resource patterns).
func ossXmlCall(client *connectivity.AliyunClient, method, op, action string, query map[string]*string, body interface{}, hostMap map[string]*string, autoRetry bool) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.Do("Oss", xmlParam(method, "2019-05-17", op, action), query, body, nil, hostMap, autoRetry)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(op, response)
		return nil
	})
	return response, err
}

// ossListInRegionSweepBuckets enumerates OSS bucket names in the sweep region
// that pass the name gate (or all in-region buckets when sweepAll). Uses the
// OSS SDK ListBuckets and the same region filter as the existing bucket
// sweeper (Location suffix == RegionId). Pagination via Marker/NextMarker.
func ossListInRegionSweepBuckets(client *connectivity.AliyunClient) ([]string, error) {
	out := []string{}
	nextMarker := ""
	for {
		var opts []oss.Option
		opts = append(opts, oss.MaxKeys(1000))
		if nextMarker != "" {
			opts = append(opts, oss.Marker(nextMarker))
		}
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.ListBuckets(opts...)
		})
		if err != nil {
			return nil, fmt.Errorf("ListBuckets: %s", err)
		}
		resp, _ := raw.(oss.ListBucketsResult)
		if len(resp.Buckets) == 0 {
			break
		}
		for _, b := range resp.Buckets {
			// Region filter: OSS bucket Location is the region id (e.g.
			// "oss-cn-beijing"); only buckets whose Location matches the sweep
			// region are walked (red-line: never clean outside the sweep region).
			if !strings.HasSuffix(b.Location, client.RegionId) {
				continue
			}
			if !ossSweepNameOk(b.Name) {
				continue
			}
			out = append(out, b.Name)
		}
		nextMarker = resp.NextMarker
		if nextMarker == "" {
			break
		}
	}
	return out, nil
}

// testSweepOssAccessPoints sweeps alicloud_oss_access_point leftovers. Access
// points are regional resources bound to a bucket; the resource id is
// "<bucket>:<access_point_name>". Enumeration is per in-region gated bucket
// (ListAccessPoints with hostMap["bucket"]), which guarantees region-scoping by
// construction, then DeleteAccessPoint per access point.
func testSweepOssAccessPoints(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	buckets, err := ossListInRegionSweepBuckets(client)
	if err != nil {
		// Without the in-region bucket set we cannot guarantee region-scoping,
		// so do nothing this run (safe). The bucket sweeper still runs.
		log.Printf("[ERROR] OSS access-point sweep: listing in-region buckets failed: %s (skipping)", err)
		return nil
	}

	for _, bucket := range buckets {
		nextToken := ""
		for {
			query := map[string]*string{"max-keys": StringPointer("100")}
			if nextToken != "" {
				query["continuation-token"] = StringPointer(nextToken)
			}
			hostMap := map[string]*string{"bucket": StringPointer(bucket)}
			response, err := ossXmlCall(client, "GET", "ListAccessPoints", "/?accessPoint", query, nil, hostMap, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"NoSuchBucket"}) || NotFoundError(err) {
					break // bucket gone — nothing to sweep here
				}
				log.Printf("[ERROR] ListAccessPoints for bucket %s: %s", bucket, err)
				break
			}
			apsRaw, _ := jsonpath.Get("$.ListAccessPointsResult.AccessPoints.AccessPoint", response)
			for _, ap := range ossCoerceList(apsRaw) {
				name, _ := ap["AccessPointName"].(string)
				if name == "" {
					continue
				}
				log.Printf("[INFO] Deleting OSS access point: %s:%s", bucket, name)
				delHost := map[string]*string{"bucket": StringPointer(bucket)}
				delQuery := map[string]*string{"x-oss-access-point-name": StringPointer(name)}
				if _, err := ossXmlCall(client, "DELETE", "DeleteAccessPoint", "/?accessPoint", delQuery, nil, delHost, false); err != nil {
					if NotFoundError(err) || IsExpectedErrors(err, []string{"NoSuchAccessPoint", "NoSuchBucket"}) {
						continue
					}
					log.Printf("[ERROR] Failed to delete OSS access point %s:%s: %s", bucket, name, err)
				}
			}
			tok, _ := jsonpath.Get("$.ListAccessPointsResult.NextContinuationToken", response)
			if t, ok := tok.(string); ok {
				nextToken = t
			} else {
				nextToken = ""
			}
			if nextToken == "" {
				break
			}
		}
	}
	return nil
}

func init() {
	resource.AddTestSweepers("alicloud_oss_access_point", &resource.Sweeper{
		Name: "alicloud_oss_access_point",
		F:    testSweepOssAccessPoints,
	})
}
