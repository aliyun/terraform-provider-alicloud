package polardb

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountPrivilegeCreate,
		Read:   resourceAccountPrivilegeRead,
		Update: resourceAccountPrivilegeUpdate,
		Delete: resourceAccountPrivilegeDelete,
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

			"account_privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ReadOnly", "ReadWrite", "DMLOnly", "DDLOnly"}, false),
				Default:      "ReadOnly",
				ForceNew:     true,
			},

			"db_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
		},
	}
}

func resourceAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	clusterId := d.Get("db_cluster_id").(string)
	account := d.Get("account_name").(string)
	privilege := d.Get("account_privilege").(string)
	dbList := d.Get("db_names").(*schema.Set).List()
	// wait instance running before granting
	if err := polarDBService.WaitForPolarDBInstance(clusterId, Running, DefaultLongTimeout); err != nil {
		return errs.WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", clusterId, COLON_SEPARATED, account, COLON_SEPARATED, privilege))

	if len(dbList) > 0 {
		for _, db := range dbList {
			if err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				if err := polarDBService.GrantPolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					if errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				return errs.WrapError(err)
			}
		}
	}

	return resourceAccountPrivilegeRead(d, meta)
}

func resourceAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := polarDBService.DescribePolarDBAccountPrivilege(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("account_name", object.AccountName)
	d.Set("account_privilege", parts[2])
	var names []string
	for _, pri := range object.DatabasePrivileges {
		if pri.AccountPrivilege == parts[2] {
			names = append(names, pri.DBName)
		}
	}
	// Sort the db_names to ensure consistent ordering and avoid plan drift
	sort.Strings(names)
	d.Set("db_names", names)

	return nil
}

func resourceAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	Service := Service{client}
	d.Partial(true)

	if d.HasChange("db_names") {
		parts := strings.Split(d.Id(), COLON_SEPARATED)

		o, n := d.GetChange("db_names")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			// wait instance running before revoking
			if err := Service.WaitForPolarDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return errs.WrapError(err)
			}
			for _, db := range remove {
				if err := Service.RevokePolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					return errs.WrapError(err)
				}
			}
		}

		if len(add) > 0 {
			// wait instance running before granting
			if err := Service.WaitForPolarDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return errs.WrapError(err)
			}
			for _, db := range add {
				if err := Service.GrantPolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					return errs.WrapError(err)
				}
			}
		}
	}

	d.Partial(false)
	return resourceAccountPrivilegeRead(d, meta)
}

func resourceAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	Service := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := Service.DescribePolarDBAccountPrivilege(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}
	var dbName string

	if len(object.DatabasePrivileges) > 0 {
		for _, pri := range object.DatabasePrivileges {
			if pri.AccountPrivilege == parts[2] {
				dbName = pri.DBName
				if err := Service.RevokePolarDBAccountPrivilege(d.Id(), pri.DBName); err != nil {
					return errs.WrapError(err)
				}
			}
		}
	}

	return Service.WaitForPolarDBAccountPrivilege(d.Id(), dbName, Deleted, DefaultTimeoutMedium)
}
