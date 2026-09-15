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
	rdsService := RdsService{client}

	objectRaw, err := rdsServiceV2.DescribeRdsDatabase(d.Id())
	if err != nil {
		// NotFound => the database is gone. A 403 OperationDenied(Read)DBInstanceStatus
		// means the parent instance is in a terminal non-permissible state (refunded /
		// unsubscribed out of band and now in the recycle bin). For that 403, confirm
		// via DescribeDBInstance (which maps the same codes to NotFound) that the
		// parent is actually gone before dropping state — so refresh stops hard-failing
		// without mistakenly clearing a live-but-transiently-locked instance.
		if !d.IsNewResource() && (isParentGone(err) || IsExpectedErrors(err, dbInstanceGoneStatusCodes)) {
			if IsExpectedErrors(err, dbInstanceGoneStatusCodes) {
				// A gone 403 means the parent is (likely) in the recycle bin, but
				// confirm via the choke point before dropping state so a live-but-
				// transiently-locked instance is not cleared by mistake. If the id
				// cannot be parsed, surface that instead of silently dropping state.
				parts, perr := ParseResourceId(d.Id(), 2)
				if perr != nil {
					return WrapError(perr)
				}
				if _, e := rdsService.DescribeDBInstance(parts[0]); e == nil || !isParentGone(e) {
					return WrapError(err)
				}
			}
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
	// If the instance has already been removed outside of Terraform (console /
	// CLI, or orphaned via `terraform state rm` and then deleted) or is being
	// deleted, the database no longer exists and DeleteDatabase cannot be
	// invoked. Return nil here so refresh / destroy does not block for up to
	// 30 minutes waiting for the instance to reach Running.
	instance, err := rdsService.DescribeDBInstance(parts[0])
	if err != nil {
		if isParentGone(err) {
			return nil
		}
		return WrapError(err)
	}
	// The instance is describable. If it is already in a terminal deleting flow
	// the database no longer exists and DeleteDatabase cannot be invoked — return
	// nil so destroy is idempotent. (The recycle-bin / refunded state surfaces as
	// a 403 above, now mapped to NotFound and already returned nil.)
	if status := fmt.Sprint(instance["DBInstanceStatus"]); status == "Deleting" {
		return nil
	}
	// The instance exists and is not in a terminal state; wait for it to reach
	// Running so DeleteDatabase is permitted, then call it with retry on the
	// transient instance-status race (consistent with alicloud_db_instance).
	if err := rdsService.WaitForDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
		// The instance went terminal (recycle bin / refunded) during the wait for
		// Running — the database no longer exists, finish delete idempotently.
		if isParentGone(err) {
			return nil
		}
		return WrapError(err)
	}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
				return nil
			}
			// Terminal 403 (parent instance in recycle bin / refunded): confirm via
			// DescribeDBInstance (the choke point that maps the same codes to NotFound)
			// that the parent is gone, then finish delete idempotently instead of
			// blind-retrying an unmanageable instance until the delete timeout. Mirrors
			// RevokeAccountPrivilege's mid-loop gone-check.
			if IsExpectedErrors(err, dbInstanceGoneStatusCodes) {
				if _, e := rdsService.DescribeDBInstance(parts[0]); e != nil && isParentGone(e) {
					return nil
				}
			}
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus", "IncorrectDBInstanceState"}) || NeedRetry(err) {
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
	return WrapError(rdsService.WaitForDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
