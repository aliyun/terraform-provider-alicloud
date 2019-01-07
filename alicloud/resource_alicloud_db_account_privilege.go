package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateAllowedStringValue([]string{string(ReadOnly), string(ReadWrite)}),
				Default:      ReadOnly,
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
	if err := rsdService.WaitForDBInstance(instanceId, Running, 1800); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}
	if len(dbList) > 0 {
		for _, db := range dbList {
			if err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				if e := rsdService.GrantAccountPrivilege(instanceId, account, db.(string), privilege); e != nil {
					if IsExceptedErrors(e, OperationDeniedDBStatus) {
						return resource.RetryableError(fmt.Errorf("Grant Account %s Privilege %s fot DB %s timeout and got an error: %#v", account, privilege, db.(string), e))
					}
					return resource.NonRetryableError(fmt.Errorf("Grant Account %s Privilege %s fot DB %s got an error: %#v", account, privilege, db.(string), e))
				}
				return nil
			}); err != nil {
				return fmt.Errorf("Grant Account %s Privilege %s got an error: %#v", account, privilege, err)
			}
		}
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", instanceId, COLON_SEPARATED, account, COLON_SEPARATED, privilege))

	return resourceAlicloudDBAccountPrivilegeUpdate(d, meta)
}

func resourceAlicloudDBAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rsdService := RdsService{client}
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	account, err := rsdService.DescribeDatabaseAccount(parts[0], parts[1])
	if err != nil {
		if rsdService.NotFoundDBInstance(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Describe db account got an error: %#v", err)
	}

	d.Set("instance_id", account.DBInstanceId)
	d.Set("account_name", account.AccountName)
	d.Set("privilege", parts[2])
	var names []string
	for _, pri := range account.DatabasePrivileges.DatabasePrivilege {
		if pri.AccountPrivilege == parts[2] {
			names = append(names, pri.DBName)
		}
	}
	d.Set("db_names", names)

	return nil
}

func resourceAlicloudDBAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rsdService := RdsService{client}
	d.Partial(true)

	if d.HasChange("db_names") && !d.IsNewResource() {
		parts := strings.Split(d.Id(), COLON_SEPARATED)

		o, n := d.GetChange("db_names")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			// wait instance running before revoking
			if err := rsdService.WaitForDBInstance(parts[0], Running, 500); err != nil {
				return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
			}
			for _, db := range remove {
				if err := rsdService.RevokeAccountPrivilege(parts[0], parts[1], db.(string)); err != nil {
					return err
				}
			}
		}

		if len(add) > 0 {
			// wait instance running before granting
			if err := rsdService.WaitForDBInstance(parts[0], Running, 500); err != nil {
				return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
			}
			for _, db := range add {
				if err := rsdService.GrantAccountPrivilege(parts[0], parts[1], db.(string), parts[2]); err != nil {
					return err
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
	rsdService := RdsService{client}
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	account, err := rsdService.DescribeDatabaseAccount(parts[0], parts[1])
	if err != nil {
		if rsdService.NotFoundDBInstance(err) {
			return nil
		}
		return fmt.Errorf("Describe db account got an error: %#v", err)
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		if len(account.DatabasePrivileges.DatabasePrivilege) > 0 {
			for _, pri := range account.DatabasePrivileges.DatabasePrivilege {
				if pri.AccountPrivilege == parts[2] {
					if err := rsdService.RevokeAccountPrivilege(parts[0], parts[1], pri.DBName); err != nil {
						return resource.NonRetryableError(fmt.Errorf("Revoke DB %s account %s privilege got an error: %#v.", pri.DBName, account, err))
					}
				}
			}
		}
		account, err := rsdService.DescribeDatabaseAccount(parts[0], parts[1])
		if err != nil {
			if rsdService.NotFoundDBInstance(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe db account got an error: %#v", err))
		}
		if len(account.DatabasePrivileges.DatabasePrivilege) > 0 {
			return resource.RetryableError(fmt.Errorf("Revoke account %s privilege timeout and got an error: %#v.", account, err))
		}
		return nil
	})
}
