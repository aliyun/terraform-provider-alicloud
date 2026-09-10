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

func resourceAliCloudMongodbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbAccountCreate,
		Read:   resourceAliCloudMongodbAccountRead,
		Update: resourceAliCloudMongodbAccountUpdate,
		Delete: resourceAliCloudMongodbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"character_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"db", "normal"}, false),
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

func resourceAliCloudMongodbAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	if v, ok := d.GetOk("character_type"); ok && InArray(fmt.Sprint(v), []string{"db"}) {
		action := "CreateAccount"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("account_name"); ok {
			request["AccountName"] = v
		}
		if v, ok := d.GetOk("instance_id"); ok {
			request["DBInstanceId"] = v
		}
		request["RegionId"] = client.RegionId

		request["AccountPassword"] = d.Get("account_password")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_account", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], request["AccountName"]))

		mongodbServiceV2 := MongodbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 4*time.Minute, mongodbServiceV2.MongodbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	if v, ok := d.GetOk("character_type"); !ok || InArray(fmt.Sprint(v), []string{"normal"}) {
		action := "ResetAccountPassword"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("account_name"); ok {
			request["AccountName"] = v
		}
		if v, ok := d.GetOk("instance_id"); ok {
			request["DBInstanceId"] = v
		}
		request["RegionId"] = client.RegionId

		if v, ok := d.GetOk("character_type"); ok {
			request["CharacterType"] = v
		}
		request["AccountPassword"] = d.Get("account_password")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_account", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], request["AccountName"]))

		mongodbServiceV2 := MongodbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	return resourceAliCloudMongodbAccountUpdate(d, meta)
}

func resourceAliCloudMongodbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}

	objectRaw, err := mongodbServiceV2.DescribeMongodbAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_account DescribeMongodbAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_description", objectRaw["AccountDescription"])
	d.Set("character_type", objectRaw["CharacterType"])
	d.Set("status", objectRaw["AccountStatus"])
	d.Set("account_name", objectRaw["AccountName"])
	d.Set("instance_id", objectRaw["DBInstanceId"])

	return nil
}

func resourceAliCloudMongodbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	mongodbServiceV2 := MongodbServiceV2{client}
	objectRaw, _ := mongodbServiceV2.DescribeMongodbAccount(d.Id())

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ResetAccountPassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("character_type") {
		update = true
		request["CharacterType"] = d.Get("character_type")
	}

	if !d.IsNewResource() && d.HasChange("account_password") {
		update = true
	}
	request["AccountPassword"] = d.Get("account_password")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
		mongodbServiceV2 := MongodbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	objectRaw, _ = mongodbServiceV2.DescribeMongodbAccount(d.Id())
	enableModifyAccountDescription1 := false
	checkValue00 := objectRaw["CharacterType"]
	if checkValue00 == "normal" {
		enableModifyAccountDescription1 = true
	}
	parts = strings.Split(d.Id(), ":")
	action = "ModifyAccountDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("account_description") {
		update = true
	}
	request["AccountDescription"] = d.Get("account_description")
	if !d.IsNewResource() && d.HasChange("character_type") {
		update = true
		request["CharacterType"] = d.Get("character_type")
	}

	if update && enableModifyAccountDescription1 {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
		mongodbServiceV2 := MongodbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAliCloudMongodbAccountRead(d, meta)
}

func resourceAliCloudMongodbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Account. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
