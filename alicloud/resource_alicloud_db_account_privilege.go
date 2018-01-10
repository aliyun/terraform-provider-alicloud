package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"privilege": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(rds.ReadOnly), string(rds.ReadWrite)}),
				Default:      rds.ReadOnly,
			},

			"db_names": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudDBAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(fmt.Sprintf("%s%s%s%s%s", d.Get("instance_id").(string), COLON_SEPARATED, d.Get("account_name").(string), COLON_SEPARATED, d.Get("privilege").(string)))

	return resourceAlicloudDBAccountPrivilegeUpdate(d, meta)
}

func resourceAlicloudDBAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	account, err := meta.(*AliyunClient).DescribeDatabaseAccount(parts[0], parts[1])
	if err != nil {
		if NotFoundError(err) {
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
		if pri.AccountPrivilege == rds.AccountPrivilege(parts[2]) {
			names = append(names, pri.DBName)
		}
	}
	d.Set("db_names", names)

	return nil
}

func resourceAlicloudDBAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	update := false

	if d.HasChange("privilege") {
		update = true
		d.SetPartial("privilege")
	}

	if d.HasChange("db_names") {
		update = true
		d.SetPartial("db_names")
	}

	if update {
		o, n := d.GetChange("db_names")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			for _, db := range remove {
				if err := client.RevokeAccountPrivilege(parts[0], parts[1], db.(string)); err != nil {
					return err
				}
			}
		}

		if len(add) > 0 {
			for _, db := range add {
				if err := client.GrantAccountPrivilege(parts[0], parts[1], db.(string), parts[2]); err != nil {
					return err
				}
			}
		}
	}

	d.Partial(false)
	return resourceAlicloudDBAccountPrivilegeRead(d, meta)
}

func resourceAlicloudDBAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	account, err := client.DescribeDatabaseAccount(parts[0], parts[1])
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, InvalidAccountNameNotFound) {
			return nil
		}
		return fmt.Errorf("Describe db account got an error: %#v", err)
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		if len(account.DatabasePrivileges.DatabasePrivilege) > 0 {
			for _, pri := range account.DatabasePrivileges.DatabasePrivilege {
				if pri.AccountPrivilege == rds.AccountPrivilege(parts[2]) {
					if err := client.RevokeAccountPrivilege(parts[0], parts[1], pri.DBName); err != nil {
						return resource.NonRetryableError(fmt.Errorf("Revoke DB %s account %s privilege got an error: %#v.", pri.DBName, account, err))
					}
				}
			}
		}
		account, err := client.DescribeDatabaseAccount(parts[0], parts[1])
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidAccountNameNotFound) {
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
