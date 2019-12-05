package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDBAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBAccountPrivilegeCreate,
		Read:   resourceAlicloudDBAccountPrivilegeRead,
		Update: resourceAlicloudDBAccountPrivilegeUpdate,
		Delete: resourceAlicloudDBAccountPrivilegeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ReadOnly", "ReadWrite", "DDLOnly", "DMLOnly", "DBOwner"}, false),
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

func resourceAlicloudDBAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rsdService := RdsService{client}
	instanceId := d.Get("instance_id").(string)
	account := d.Get("account_name").(string)
	privilege := d.Get("privilege").(string)
	dbList := d.Get("db_names").(*schema.Set).List()
	// wait instance running before granting
	if err := rsdService.WaitForDBInstance(instanceId, Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", instanceId, COLON_SEPARATED, account, COLON_SEPARATED, privilege))

	if len(dbList) > 0 {
		for _, db := range dbList {
			if err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				if err := rsdService.GrantAccountPrivilege(d.Id(), db.(string)); err != nil {
					if IsExceptedErrors(err, OperationDeniedDBStatus) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAlicloudDBAccountPrivilegeRead(d, meta)
}

func resourceAlicloudDBAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rsdService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := rsdService.DescribeDBAccountPrivilege(d.Id())
	if err != nil {
		if rsdService.NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.DBInstanceId)
	d.Set("account_name", object.AccountName)
	d.Set("privilege", parts[2])
	var names []string
	for _, pri := range object.DatabasePrivileges.DatabasePrivilege {
		if pri.AccountPrivilege == parts[2] {
			names = append(names, pri.DBName)
		}
	}
	d.Set("db_names", names)

	return nil
}

func resourceAlicloudDBAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
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
			if err := rdsService.WaitForDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range remove {
				if err := rdsService.RevokeAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}

		if len(add) > 0 {
			// wait instance running before granting
			if err := rdsService.WaitForDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range add {
				if err := rdsService.GrantAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("db_names")
	}

	d.Partial(false)
	return resourceAlicloudDBAccountPrivilegeRead(d, meta)
}

func resourceAlicloudDBAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := rdsService.DescribeDBAccountPrivilege(d.Id())
	if err != nil {
		if rdsService.NotFoundDBInstance(err) {
			return nil
		}
		return WrapError(err)
	}
	var dbName string

	if len(object.DatabasePrivileges.DatabasePrivilege) > 0 {
		for _, pri := range object.DatabasePrivileges.DatabasePrivilege {
			if pri.AccountPrivilege == parts[2] {
				dbName = pri.DBName
				if err := rdsService.RevokeAccountPrivilege(d.Id(), pri.DBName); err != nil {
					return WrapError(err)
				}
			}
		}
	}

	return rdsService.WaitForAccountPrivilege(d.Id(), dbName, Deleted, DefaultTimeoutMedium)
}
