package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBZonalAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBZonalAccountCreate,
		Read:   resourceAlicloudPolarDBZonalAccountRead,
		Update: resourceAlicloudPolarDBZonalAccountUpdate,
		Delete: resourceAlicloudPolarDBZonalAccountDelete,
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

func resourceAlicloudPolarDBZonalAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBServiceV2 := PolarDbServiceV2{client}
	request := polardb.CreateCreateAccountRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.AccountName = d.Get("account_name").(string)

	password := d.Get("account_password").(string)

	if password == "" {
		return WrapError(Error("the 'password' should be set."))
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

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := polarDBServiceV2.CreateAccount(request)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_account", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.AccountName))
	d.Set("account_password", request.AccountPassword)

	if err := polarDBServiceV2.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudPolarDBZonalAccountRead(d, meta)
}

func resourceAlicloudPolarDBZonalAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
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

	d.Set("db_cluster_id", parts[0])
	d.Set("account_name", object.AccountName)
	d.Set("account_type", object.AccountType)
	d.Set("account_description", object.AccountDescription)
	return nil
}

func resourceAlicloudPolarDBZonalAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if !d.IsNewResource() && d.HasChange("account_description") {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		accountDescription := d.Get("account_description").(string)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			err := polarDBService.modifyAccountDescription(instanceId, accountName, accountDescription)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("modifyAccountDescription", err, parts[1])
			return nil
		})
		if err != nil {
			return err
		}
	}

	if !d.IsNewResource() && (d.HasChange("account_password")) {
		if err := polarDBService.WaitForPolarDBAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}

		password := d.Get("account_password").(string)
		if password == "" {
			return WrapError(Error("the 'password' should be set."))
		}

		passwordfinal := password

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			err := polarDBService.modifyAccountPassword(instanceId, accountName, passwordfinal)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("modifyAccountPassword", err, accountName)
			return nil
		})
		if err != nil {
			return err
		}
		d.Set("account_password", passwordfinal)
	}

	return resourceAlicloudPolarDBZonalAccountRead(d, meta)
}

func resourceAlicloudPolarDBZonalAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := polarDBService.DeleteAccount(parts[0], parts[1])
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteAccount", err, parts[1])
		return nil
	})
	if err != nil {
		return err
	}

	return polarDBService.WaitForPolarDBAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
