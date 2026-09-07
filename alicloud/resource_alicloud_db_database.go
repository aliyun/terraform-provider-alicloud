// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudRdsDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsDatabaseCreate,
		Read:   resourceAliCloudRdsDatabaseRead,
		Update: resourceAliCloudRdsDatabaseUpdate,
		Delete: resourceAliCloudRdsDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_base_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_-]*[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and middle lines, and must begin with letters and end with letters or numbers"),
				ConflictsWith: []string{"name"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_-]*[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and middle lines, and must begin with letters and end with letters or numbers"),
				Deprecated:    "Field 'name' has been deprecated from provider version 1.266.0. New field 'data_base_name' instead.",
				ConflictsWith: []string{"data_base_name"},
			},
			"character_set": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					newArray := strings.Split(new, ",")
					oldArray := strings.Split(old, ",")
					if d.Id() != "" && len(oldArray) > 1 && len(newArray) == 1 && strings.ToLower(newArray[0]) == strings.ToLower(oldArray[0]) {
						return true
					}
					/*
					  SQLServer creates a database, when a non native engine character set is passed in, the SDK will assign the default character set.
					*/
					if old == "Chinese_PRC_CI_AS" && new == "utf8" {
						return true
					}
					return false
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudRdsDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDatabase"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["DBInstanceId"] = v
	}

	if v, ok := d.GetOk("data_base_name"); ok {
		request["DBName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["DBName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "data_base_name" must be set one!`))
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("character_set"); ok {
		request["CharacterSetName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["DBDescription"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.OutofUsage", "OperationDenied.DBInstanceStatus", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "InternalError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_db_database", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], request["DBName"]))

	return resourceAliCloudRdsDatabaseRead(d, meta)
}

func resourceAliCloudRdsDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsServiceV2 := RdsServiceV2{client}

	objectRaw, err := rdsServiceV2.DescribeRdsDatabase(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_db_database DescribeRdsDatabase Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["DBDescription"])
	d.Set("status", objectRaw["DBStatus"])
	d.Set("data_base_name", objectRaw["DBName"])
	d.Set("name", objectRaw["DBName"])
	d.Set("instance_id", objectRaw["DBInstanceId"])

	if string(PostgreSQL) == objectRaw["Engine"] {
		var strArray = []string{objectRaw["CharacterSetName"].(string), objectRaw["Collate"].(string), objectRaw["Ctype"].(string)}
		postgreSQLCharacterSet := strings.Join(strArray, ",")
		d.Set("character_set", postgreSQLCharacterSet)
	} else {
		d.Set("character_set", objectRaw["CharacterSetName"])
	}
	return nil
}

func resourceAliCloudRdsDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDBDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = parts[0]
	request["DBName"] = parts[1]
	request["RegionId"] = client.RegionId
	request["SourceIp"] = client.SourceIp

	if d.HasChange("description") {
		update = true
		request["DBDescription"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, query, request, true)
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

	return resourceAliCloudRdsDatabaseRead(d, meta)
}

func resourceAliCloudRdsDatabaseDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteDatabase"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": parts[0],
		"DBName":       parts[1],
		"SourceIp":     client.SourceIp,
	}
	// wait instance status is running before deleting database
	if err := rdsService.WaitForDBInstance(parts[0], Running, 1800); err != nil {
		return WrapError(err)
	}
	response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return WrapError(rdsService.WaitForDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
