// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrInstanceCustomizedDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrInstanceCustomizedDomainCreate,
		Read:   resourceAliCloudCrInstanceCustomizedDomainRead,
		Update: resourceAliCloudCrInstanceCustomizedDomainUpdate,
		Delete: resourceAliCloudCrInstanceCustomizedDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"module_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cert_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cert_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCrInstanceCustomizedDomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstanceCustomizedDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("module_name"); ok {
		request["ModuleName"] = v
	}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}
	if v, ok := d.GetOk("cert_id"); ok {
		request["CertId"] = v
	}
	if v, ok := d.GetOk("cert_region_id"); ok {
		request["CertRegionId"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_instance_customized_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["InstanceId"], request["ModuleName"], request["Domain"]))

	return resourceAliCloudCrInstanceCustomizedDomainRead(d, meta)
}

func resourceAliCloudCrInstanceCustomizedDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	objectRaw, err := crServiceV2.DescribeCrInstanceCustomizedDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cr_instance_customized_domain DescribeCrInstanceCustomizedDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("region_id", fmt.Sprint(objectRaw["RegionId"]))
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("module_name", objectRaw["ModuleName"])
	d.Set("domain", objectRaw["Domain"])
	d.Set("cert_id", objectRaw["CertId"])

	createTime, err := strconv.ParseInt(objectRaw["CreateTime"].(json.Number).String(), 10, 64)
	if err != nil {
		return WrapError(err)
	}
	d.Set("create_time", createTime)

	modifiedTime, err := strconv.ParseInt(objectRaw["ModifiedTime"].(json.Number).String(), 10, 64)
	if err != nil {
		return WrapError(err)
	}
	d.Set("modified_time", modifiedTime)

	return nil
}

func resourceAliCloudCrInstanceCustomizedDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateInstanceCustomizedDomain"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ModuleName"] = parts[1]
	request["Domain"] = parts[2]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	if d.HasChange("cert_id") {
		update = true
	}
	request["CertId"] = d.Get("cert_id")

	if v, ok := d.GetOk("cert_region_id"); ok {
		request["CertRegionId"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCrInstanceCustomizedDomainRead(d, meta)
}

func resourceAliCloudCrInstanceCustomizedDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteInstanceCustomizedDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ModuleName"] = parts[1]
	request["Domain"] = parts[2]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

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

	return nil
}
