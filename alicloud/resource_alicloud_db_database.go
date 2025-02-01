package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDBDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBDatabaseCreate,
		Read:   resourceAlicloudDBDatabaseRead,
		Update: resourceAlicloudDBDatabaseUpdate,
		Delete: resourceAlicloudDBDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_-]*[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and middle lines, and must begin with letters and end with letters or numbers"),
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

func resourceAlicloudDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "CreateDatabase"
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"DBInstanceId":     d.Get("instance_id"),
		"DBName":           d.Get("name"),
		"CharacterSetName": d.Get("character_set"),
		"SourceIp":         client.SourceIp,
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request["DBDescription"] = v
	}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
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

	d.SetId(fmt.Sprintf("%v%s%v", request["DBInstanceId"], COLON_SEPARATED, request["DBName"]))

	return resourceAlicloudDBDatabaseRead(d, meta)
}

func resourceAlicloudDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rsdService := RdsService{client}
	object, err := rsdService.DescribeDBDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object["DBInstanceId"])
	d.Set("name", object["DBName"])
	if string(PostgreSQL) == object["Engine"] {
		var strArray = []string{object["CharacterSetName"].(string), object["Collate"].(string), object["Ctype"].(string)}
		postgreSQLCharacterSet := strings.Join(strArray, ",")
		d.Set("character_set", postgreSQLCharacterSet)
	} else {
		d.Set("character_set", object["CharacterSetName"])
	}
	d.Set("description", object["DBDescription"])

	return nil
}

func resourceAlicloudDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("description") && !d.IsNewResource() {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		action := "ModifyDBDescription"
		request := map[string]interface{}{
			"RegionId":      client.RegionId,
			"DBInstanceId":  parts[0],
			"DBName":        parts[1],
			"DBDescription": d.Get("description"),
			"SourceIp":      client.SourceIp,
		}

		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
	}
	return resourceAlicloudDBDatabaseRead(d, meta)
}

func resourceAlicloudDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
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
