// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudResourceManagerMessageContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerMessageContactCreate,
		Read:   resourceAliCloudResourceManagerMessageContactRead,
		Update: resourceAliCloudResourceManagerMessageContactUpdate,
		Delete: resourceAliCloudResourceManagerMessageContactDelete,
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
			"email_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message_contact_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"phone_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"TechnicalDirector", "MaintenanceDirector", "ProjectDirector", "CEO", "Other"}, false),
			},
		},
	}
}

func resourceAliCloudResourceManagerMessageContactCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddMessageContact"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["EmailAddress"] = d.Get("email_address")
	request["Name"] = d.Get("message_contact_name")
	if v, ok := d.GetOk("phone_number"); ok {
		request["PhoneNumber"] = v
	}
	if v, ok := d.GetOk("message_types"); ok {
		messageTypesMapsArray := v.([]interface{})
		request["MessageTypes"] = messageTypesMapsArray
	}

	request["Title"] = d.Get("title")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceDirectoryMaster", "2022-04-19", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_message_contact", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Contact.ContactId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudResourceManagerMessageContactUpdate(d, meta)
}

func resourceAliCloudResourceManagerMessageContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerMessageContact(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_message_contact DescribeResourceManagerMessageContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("email_address", objectRaw["EmailAddress"])
	d.Set("message_contact_name", objectRaw["Name"])
	d.Set("phone_number", objectRaw["PhoneNumber"])
	d.Set("status", objectRaw["Status"])
	d.Set("title", objectRaw["Title"])

	messageTypesRaw := make([]interface{}, 0)
	if objectRaw["MessageTypes"] != nil {
		messageTypesRaw = objectRaw["MessageTypes"].([]interface{})
	}

	d.Set("message_types", messageTypesRaw)

	return nil
}

func resourceAliCloudResourceManagerMessageContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateMessageContact"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ContactId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("email_address") {
		update = true
	}
	request["EmailAddress"] = d.Get("email_address")
	if !d.IsNewResource() && d.HasChange("message_contact_name") {
		update = true
	}
	request["Name"] = d.Get("message_contact_name")
	if !d.IsNewResource() && d.HasChange("phone_number") {
		update = true
		request["PhoneNumber"] = d.Get("phone_number")
	}

	if !d.IsNewResource() && d.HasChange("message_types") {
		update = true
	}
	if v, ok := d.GetOk("message_types"); ok || d.HasChange("message_types") {
		messageTypesMapsArray := v.([]interface{})
		request["MessageTypes"] = messageTypesMapsArray
	}

	if !d.IsNewResource() && d.HasChange("title") {
		update = true
	}
	request["Title"] = d.Get("title")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceDirectoryMaster", "2022-04-19", action, query, request, true)
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

	return resourceAliCloudResourceManagerMessageContactRead(d, meta)
}

func resourceAliCloudResourceManagerMessageContactDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMessageContact"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ContactId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceDirectoryMaster", "2022-04-19", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExists.Contact"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
