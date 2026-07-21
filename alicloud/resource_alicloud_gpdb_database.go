// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGpdbDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbDatabaseCreate,
		Read:   resourceAliCloudGpdbDatabaseRead,
		Delete: resourceAliCloudGpdbDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"character_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"collate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ctype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudGpdbDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDatabase"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DatabaseName"] = d.Get("database_name")
	query["DBInstanceId"] = d.Get("db_instance_id")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("ctype"); ok {
		request["Ctype"] = v
	}
	request["Owner"] = d.Get("owner")
	if v, ok := d.GetOk("collate"); ok {
		request["Collate"] = v
	}
	if v, ok := d.GetOk("character_set_name"); ok {
		request["CharacterSetName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_database", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], query["DatabaseName"]))

	return resourceAliCloudGpdbDatabaseRead(d, meta)
}

func resourceAliCloudGpdbDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbDatabase(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_database DescribeGpdbDatabase Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CharacterSetName"] != nil {
		d.Set("character_set_name", objectRaw["CharacterSetName"])
	}
	if objectRaw["Collate"] != nil {
		d.Set("collate", objectRaw["Collate"])
	}
	if objectRaw["Ctype"] != nil {
		d.Set("ctype", objectRaw["Ctype"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["Owner"] != nil {
		d.Set("owner", objectRaw["Owner"])
	}
	if objectRaw["DatabaseName"] != nil {
		d.Set("database_name", objectRaw["DatabaseName"])
	}
	if objectRaw["DBInstanceId"] != nil {
		d.Set("db_instance_id", objectRaw["DBInstanceId"])
	}

	return nil
}

func resourceAliCloudGpdbDatabaseDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDatabase"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DatabaseName"] = parts[1]
	query["DBInstanceId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)

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
