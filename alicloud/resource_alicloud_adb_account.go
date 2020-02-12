package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudAdbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAdbAccountCreate,
		Read:   resourceAlicloudAdbAccountRead,
		Update: resourceAlicloudAdbAccountUpdate,
		Delete: resourceAlicloudAdbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string("Normal"), string("Super")}, false),
				Default:      "Normal",
				ForceNew:     true,
			},

			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudAdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	request := adb.CreateCreateAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.AccountName = d.Get("account_name").(string)

	password := d.Get("account_password").(string)
	request.AccountPassword = password

	if password == "" {
		return WrapError(Error("'password' should be set."))
	}

	// Description will not be set when account type is normal and it is a API bug
	if v, ok := d.GetOk("account_description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.CreateAccount(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				time.Sleep(5 * time.Second)
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

	if err := adbService.WaitForAdbAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudAdbAccountRead(d, meta)
}

func resourceAlicloudAdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbAccount(d.Id())
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

	d.Set("db_cluster_id", parts[0])
	d.Set("account_name", object.AccountName)
	d.Set("account_type", object.AccountType)
	d.Set("account_description", object.AccountDescription)

	return nil
}

func resourceAlicloudAdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("account_description") {
		if err := adbService.WaitForAdbAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := adb.CreateModifyAccountDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = instanceId
		request.AccountName = accountName
		request.AccountDescription = d.Get("account_description").(string)

		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ModifyAccountDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("account_description")
	}

	if d.HasChange("account_password") {
		if err := adbService.WaitForAdbAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := adb.CreateResetAccountPasswordRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = instanceId
		request.AccountName = accountName

		password := d.Get("account_password").(string)
		if password == "" {
			return WrapError(Error("'password' should be set."))
		}

		if password != "" {
			d.SetPartial("account_password")
			request.AccountPassword = password
		}

		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ResetAccountPassword(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("account_password")
	}

	d.Partial(false)
	return resourceAlicloudAdbAccountRead(d, meta)
}

func resourceAlicloudAdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := adb.CreateDeleteAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DeleteAccount(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidAccountName.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return adbService.WaitForAdbAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
