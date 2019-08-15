package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"character_set": {
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(CHARACTER_SET_NAME),
				Optional:     true,
				Default:      "utf8",
				ForceNew:     true,
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
	rdsService := RdsService{client}
	request := rds.CreateCreateDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Get("instance_id").(string)
	request.DBName = d.Get("name").(string)
	request.CharacterSetName = d.Get("character_set").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	if inst, err := rdsService.DescribeDBInstance(request.DBInstanceId); err != nil {
		return WrapError(err)
	} else if inst.Engine == string(PostgreSQL) || inst.Engine == string(PPAS) {
		return WrapError(Error("At present, it does not support creating 'PostgreSQL' and 'PPAS' database. Please login DB instance to create."))
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.CreateDatabase(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, request.DBName))

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

	d.Set("instance_id", object.DBInstanceId)
	d.Set("name", object.DBName)
	d.Set("character_set", object.CharacterSetName)
	d.Set("description", object.DBDescription)

	return nil
}

func resourceAlicloudDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("description") && !d.IsNewResource() {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		request := rds.CreateModifyDBDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = parts[0]
		request.DBName = parts[1]
		request.DBDescription = d.Get("description").(string)
		var raw interface{}
		raw, err = client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBDescription(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
	request := rds.CreateDeleteDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = parts[0]
	request.DBName = parts[1]
	// wait instance status is running before deleting database
	if err := rdsService.WaitForDBInstance(parts[0], Running, 1800); err != nil {
		return WrapError(err)
	}
	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DeleteDatabase(request)
	})
	if err != nil {
		if rdsService.NotFoundDBInstance(err) || IsExceptedError(err, InvalidDBNameNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(rdsService.WaitForDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
