package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDBAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBAccountCreate,
		Read:   resourceAlicloudDBAccountRead,
		Update: resourceAlicloudDBAccountUpdate,
		Delete: resourceAlicloudDBAccountDelete,
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

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(DBAccountNormal), string(DBAccountSuper)}),
				ForceNew:     true,
				Computed:     true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDBAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	request := rds.CreateCreateAccountRequest()
	request.DBInstanceId = d.Get("instance_id").(string)
	request.AccountName = d.Get("name").(string)
	request.AccountPassword = d.Get("password").(string)
	request.AccountType = d.Get("type").(string)

	// Description will not be set when account type is normal and it is a API bug
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}
	// wait instance running before modifying
	if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.CreateAccount(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_db_account", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, request.AccountName))

	if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudDBAccountRead(d, meta)
}

func resourceAlicloudDBAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBAccount(d.Id())
	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.DBInstanceId)
	d.Set("name", object.AccountName)
	d.Set("type", object.AccountType)
	d.Set("description", object.AccountDescription)

	return nil
}

func resourceAlicloudDBAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("description") {
		if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := rds.CreateModifyAccountDescriptionRequest()
		request.DBInstanceId = instanceId
		request.AccountName = accountName
		request.AccountDescription = d.Get("description").(string)

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyAccountDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		d.SetPartial("description")
	}

	if d.HasChange("password") {
		if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := rds.CreateResetAccountPasswordRequest()
		request.DBInstanceId = instanceId
		request.AccountName = accountName
		request.AccountPassword = d.Get("password").(string)

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ResetAccountPassword(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlicloudDBAccountRead(d, meta)
}

func resourceAlicloudDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := rds.CreateDeleteAccountRequest()
	request.DBInstanceId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DeleteAccount(request)
	})
	if err != nil && !IsExceptedError(err, InvalidAccountNameNotFound) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return rdsService.WaitForAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
