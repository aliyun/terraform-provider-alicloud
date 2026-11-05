package polardb

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceZonalAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceZonalAccountCreate,
		Read:   resourceZonalAccountRead,
		Update: resourceZonalAccountUpdate,
		Delete: resourceZonalAccountDelete,
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
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},

			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string("Normal"), string("Super")}, false),
				Computed:     true,
				ForceNew:     true,
			},

			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceZonalAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBServiceV2 := PolarDbServiceV2{client}
	request := polardb.CreateCreateAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.AccountName = d.Get("account_name").(string)

	password := d.Get("account_password").(string)

	if password == "" {
		return errs.WrapError(errs.Error("the 'password' should be set."))
	}
	request.AccountPassword = password

	if accountType, ok := d.GetOk("account_type"); ok && accountType != "" {
		request.AccountType = accountType.(string)
	} else {
		request.AccountType = "Normal"
	}

	// Description will not be set when account type is normal and it is a API bug
	if v, ok := d.GetOk("account_description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}

	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		err := polarDBServiceV2.CreateAccount(request)
		if err != nil {
			if helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_account", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.AccountName))
	d.Set("account_password", request.AccountPassword)

	if err := polarDBServiceV2.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}

	return resourceZonalAccountRead(d, meta)
}

func resourceZonalAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	object, err := polarDBService.DescribePolarDBAccount(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("account_name", object.AccountName)
	d.Set("account_type", object.AccountType)
	d.Set("account_description", object.AccountDescription)
	return nil
}

func resourceZonalAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if !d.IsNewResource() && d.HasChange("account_description") {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errs.WrapError(err)
		}
		accountDescription := d.Get("account_description").(string)

		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			err := polarDBService.modifyAccountDescription(instanceId, accountName, accountDescription)
			if err != nil {
				if helper.NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug("modifyAccountDescription", err, parts[1])
			return nil
		})
		if err != nil {
			return err
		}
	}

	if !d.IsNewResource() && (d.HasChange("account_password")) {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errs.WrapError(err)
		}

		password := d.Get("account_password").(string)
		if password == "" {
			return errs.WrapError(errs.Error("the 'password' should be set."))
		}

		passwordfinal := password

		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			err := polarDBService.modifyAccountPassword(instanceId, accountName, passwordfinal)
			if err != nil {
				if helper.NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug("modifyAccountPassword", err, accountName)
			return nil
		})
		if err != nil {
			return err
		}
		d.Set("account_password", passwordfinal)
	}

	return resourceZonalAccountRead(d, meta)
}

func resourceZonalAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		err := polarDBService.DeleteAccount(parts[0], parts[1])
		if err != nil {
			if helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug("DeleteAccount", err, parts[1])
		return nil
	})
	if err != nil {
		return err
	}

	return polarDBService.WaitForPolarDBAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
