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

func resourceAlicloudCrInternetEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrInternetEndpointCreate,
		Read:   resourceAlicloudCrInternetEndpointRead,
		Update: resourceAlicloudCrInternetEndpointUpdate,
		Delete: resourceAlicloudCrInternetEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"entry": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCrInternetEndpointCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UpdateInstanceEndpointStatus"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["ModuleName"] = "Registry"
	request["EndpointType"] = "Internet"
	request["Enable"] = "true"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_internet_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"]))

	crService := CrService{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, crService.CrInternetEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCrInternetEndpointUpdate(d, meta)
}

func resourceAlicloudCrInternetEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	objectRaw, err := crService.DescribeCrInternetEndpoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_internet_endpoint DescribeCrInternetEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["Status"])

	aclEntriesRaw := objectRaw["AclEntries"]
	entriesMaps := make([]map[string]interface{}, 0)
	if aclEntriesRaw != nil {
		for _, aclEntriesChildRaw := range convertToInterfaceArray(aclEntriesRaw) {
			aclEntriesChildRaw := aclEntriesChildRaw.(map[string]interface{})
			// GetInstanceEndpoint always returns an auto-added loopback ACL
			// policy (entry 127.0.0.1/32, comment "default") once the endpoint
			// is enabled. It is system-managed (not creatable/deletable via the
			// ACL policy APIs this resource uses), so exclude it from state to
			// avoid a perpetual plan diff against the user's config.
			if fmt.Sprint(aclEntriesChildRaw["Entry"]) == "127.0.0.1/32" && fmt.Sprint(aclEntriesChildRaw["Comment"]) == "default" {
				continue
			}
			entriesMap := make(map[string]interface{})
			entriesMap["comment"] = aclEntriesChildRaw["Comment"]
			entriesMap["entry"] = aclEntriesChildRaw["Entry"]

			entriesMaps = append(entriesMaps, entriesMap)
		}
	}
	if err := d.Set("entries", entriesMaps); err != nil {
		return err
	}

	d.Set("instance_id", d.Id())

	return nil
}

func resourceAlicloudCrInternetEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("entries") {
		var err error
		oldEntry, newEntry := d.GetChange("entries")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			for _, dataLoop := range removed.List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				action := "DeleteInstanceEndpointAclPolicy"
				request := make(map[string]interface{})
				query := make(map[string]interface{})
				request["InstanceId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ModuleName"] = "Registry"
				request["EndpointType"] = "Internet"
				request["Entry"] = dataLoopTmp["entry"]
				wait := incrementalWait(3*time.Second, 5*time.Second)
				var response map[string]interface{}
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("cr", "2018-12-01", action, query, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"SLB_SERVICE_ERROR"}) || NeedRetry(err) {
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
			for _, dataLoop := range added.List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				action := "CreateInstanceEndpointAclPolicy"
				request := make(map[string]interface{})
				query := make(map[string]interface{})
				request["InstanceId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ModuleName"] = "Registry"
				request["EndpointType"] = "Internet"
				request["Entry"] = dataLoopTmp["entry"]
				if dataLoopTmp["comment"] != nil && fmt.Sprint(dataLoopTmp["comment"]) != "" {
					request["Comment"] = dataLoopTmp["comment"]
				}
				wait := incrementalWait(5*time.Second, 5*time.Second)
				var response map[string]interface{}
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("cr", "2018-12-01", action, query, request, false)
					if err != nil {
						if IsExpectedErrors(err, []string{"SLB_SERVICE_ERROR"}) || NeedRetry(err) {
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

	}
	return resourceAlicloudCrInternetEndpointRead(d, meta)
}

func resourceAlicloudCrInternetEndpointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "UpdateInstanceEndpointStatus"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["Enable"] = "false"
	request["ModuleName"] = "Registry"
	request["EndpointType"] = "Internet"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)
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

	crService := CrService{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, crService.CrInternetEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
