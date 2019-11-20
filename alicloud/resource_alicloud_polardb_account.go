package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudPolarDBAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBAccountCreate,
		Read:   resourceAlicloudPolarDBAccountRead,
		Update: resourceAlicloudPolarDBAccountUpdate,
		Delete: resourceAlicloudPolarDBAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
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
				ValidateFunc: validation.StringInSlice([]string{string(PolarDBAccountNormal), string(PolarDBAccountSuper)}, false),
				Default:      PolarDBAccountNormal,
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

func resourceAlicloudPolarDBAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request := polardb.CreateCreateAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("cluster_id").(string)
	request.AccountName = d.Get("name").(string)
	request.AccountPassword = d.Get("password").(string)
	request.AccountType = d.Get("type").(string)

	// Description will not be set when account type is normal and it is a API bug
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}
	// wait instance running before modifying
	if err := polarDBService.WaitForPolarDBInstance(request.DBClusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateAccount(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_db_account", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.AccountName))

	if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudPolarDBAccountRead(d, meta)
}

func resourceAlicloudPolarDBAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("cluster_id", parts[0])
	d.Set("name", object.AccountName)
	d.Set("type", object.AccountType)
	d.Set("description", object.AccountDescription)

	return nil
}

func resourceAlicloudPolarDBAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("description") {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := polardb.CreateModifyAccountDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = instanceId
		request.AccountName = accountName
		request.AccountDescription = d.Get("description").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyAccountDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("description")
	}

	if d.HasChange("password") {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := polardb.CreateModifyAccountPasswordRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = instanceId
		request.AccountName = accountName
		request.NewAccountPassword = d.Get("password").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyAccountPassword(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBAccountRead(d, meta)
}

func resourceAlicloudPolarDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateDeleteAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DeleteAccount(request)
	})
	if err != nil && !IsExceptedError(err, InvalidAccountNameNotFound) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return polarDBService.WaitForPolarDBAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
