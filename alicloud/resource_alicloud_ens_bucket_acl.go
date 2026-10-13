// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ensBucketAclWaitConsistent polls GetBucketAcl until the returned ACL matches
// the value just submitted via PutBucketAcl AND two consecutive reads return
// the same value. The ENS backend is eventually consistent and non-deterministic
// after a write: within the same second GetBucketAcl can oscillate between the
// old and new ACL several times before settling. A single matching read is not
// enough to guarantee the next Read will also see the correct value, so we
// require two consecutive polls that both match the expected value before
// returning. NotFound is retried (the resource may be briefly invisible after
// the write).
func ensBucketAclWaitConsistent(client *connectivity.AliyunClient, id, expected string, timeout time.Duration) error {
	ensServiceV2 := EnsServiceV2{client}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var prevValue string
	consecutive := 0
	return resource.Retry(timeout, func() *resource.RetryError {
		objectRaw, err := ensServiceV2.DescribeEnsBucketAcl(id)
		if err != nil {
			if NotFoundError(err) {
				consecutive = 0
				prevValue = ""
				wait()
				return resource.RetryableError(fmt.Errorf("ENS bucket acl %s not converged yet (still not found)", id))
			}
			return resource.NonRetryableError(err)
		}
		current := fmt.Sprint(objectRaw["BucketAcl"])
		if current != expected {
			consecutive = 0
			prevValue = ""
			wait()
			return resource.RetryableError(fmt.Errorf("ENS bucket acl %s not converged yet (expected %s, got %s)", id, expected, current))
		}
		if current == prevValue {
			consecutive++
		} else {
			consecutive = 1
		}
		prevValue = current
		if consecutive >= 2 {
			return nil
		}
		wait()
		return resource.RetryableError(fmt.Errorf("ENS bucket acl %s matched %s once, waiting for a second consecutive read to confirm stability", id, current))
	})
}

// ensBucketAclReadStable polls GetBucketAcl until two consecutive reads return
// the same ACL value, absorbing the ENS backend's non-deterministic post-write
// oscillation. Unlike ensBucketAclWaitConsistent it has no expected value to
// compare against, so it is used by Read (including import) to ensure the value
// adopted into state is stable rather than a transient stale read. NotFound is
// not retried here: in import/refresh semantics a missing resource should be
// surfaced so the caller can drop it from state.
func ensBucketAclReadStable(client *connectivity.AliyunClient, id string, timeout time.Duration) (map[string]interface{}, error) {
	ensServiceV2 := EnsServiceV2{client}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var prevValue string
	consecutive := 0
	var stableObject map[string]interface{}
	err := resource.Retry(timeout, func() *resource.RetryError {
		objectRaw, err := ensServiceV2.DescribeEnsBucketAcl(id)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		current := fmt.Sprint(objectRaw["BucketAcl"])
		if current == prevValue {
			consecutive++
		} else {
			consecutive = 1
		}
		prevValue = current
		stableObject = objectRaw
		if consecutive >= 2 {
			return nil
		}
		wait()
		return resource.RetryableError(fmt.Errorf("ENS bucket acl %s read not stable yet (got %s), waiting for a second consecutive read", id, current))
	})
	return stableObject, err
}

func resourceAliCloudEnsBucketAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsBucketAclCreate,
		Read:   resourceAliCloudEnsBucketAclRead,
		Update: resourceAliCloudEnsBucketAclUpdate,
		Delete: resourceAliCloudEnsBucketAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket_acl": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"public-read-write", "public-read", "private"}, false),
			},
		},
	}
}

func resourceAliCloudEnsBucketAclCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PutBucketAcl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["BucketName"] = d.Get("bucket_name")
	request["BucketAcl"] = d.Get("bucket_acl")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_bucket_acl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(d.Get("bucket_name")))

	// Wait for the write to be readable before Read persists it to state
	// (GetBucketAcl is eventually consistent after a PutBucketAcl write).
	if err := ensBucketAclWaitConsistent(client, d.Id(), d.Get("bucket_acl").(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudEnsBucketAclRead(d, meta)
}

func resourceAliCloudEnsBucketAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	objectRaw, err := ensBucketAclReadStable(client, d.Id(), d.Timeout(schema.TimeoutRead))
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_bucket_acl DescribeEnsBucketAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bucket_name", d.Id())
	d.Set("bucket_acl", objectRaw["BucketAcl"])

	return nil
}

func resourceAliCloudEnsBucketAclUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PutBucketAcl"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BucketName"] = d.Get("bucket_name")
	request["BucketAcl"] = d.Get("bucket_acl")

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

	// Wait for the write to be readable before Read persists it to state
	// (GetBucketAcl is eventually consistent after a PutBucketAcl write).
	if err := ensBucketAclWaitConsistent(client, d.Id(), d.Get("bucket_acl").(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudEnsBucketAclRead(d, meta)
}

func resourceAliCloudEnsBucketAclDelete(d *schema.ResourceData, meta interface{}) error {

	// BucketAcl is a singleton attribute of an ENS bucket and cannot be
	// removed (there is no DeleteBucketAcl API). On destroy we reset the
	// ACL to the service default "private" so the cloud object returns to
	// its pre-managed state, then drop the resource from Terraform state.
	client := meta.(*connectivity.AliyunClient)
	action := "PutBucketAcl"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BucketName"] = d.Get("bucket_name")
	request["BucketAcl"] = "private"

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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"NoSuchBucket"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
