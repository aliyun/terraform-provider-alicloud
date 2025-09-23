package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBOnENSAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBOnENSAccountPrivilegeCreate,
		Read:   resourceAlicloudPolarDBOnENSAccountPrivilegeRead,
		Update: resourceAlicloudPolarDBOnENSAccountPrivilegeUpdate,
		Delete: resourceAlicloudPolarDBOnENSAccountPrivilegeDelete,
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

func resourceAlicloudPolarDBOnENSAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	clusterId := d.Get("db_cluster_id").(string)
	account := d.Get("account_name").(string)
	privilege := d.Get("account_privilege").(string)
	dbList := d.Get("db_names").(*schema.Set).List()

	// wait instance running before granting
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polarDBService.PolarDbZonalClusterStateRefreshFunc(clusterId, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", clusterId, COLON_SEPARATED, account, COLON_SEPARATED, privilege))

	if len(dbList) > 0 {
		for _, db := range dbList {
			if err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				if err := polarDBService.GrantPolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					if IsExpectedErrors(err, OperationDeniedDBStatus) {
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

	return resourceAlicloudPolarDBOnENSAccountPrivilegeRead(d, meta)
}

func resourceAlicloudPolarDBOnENSAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := polarDBService.DescribePolarDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	d.Set("db_names", names)

	return nil
}

func resourceAlicloudPolarDBOnENSAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbService := PolarDbServiceV2{client}

	if d.HasChange("db_names") {
		parts := strings.Split(d.Id(), COLON_SEPARATED)

		o, n := d.GetChange("db_names")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			// wait instance running before revoking
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polarDbService.PolarDbZonalClusterStateRefreshFunc(parts[0], []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			for _, db := range remove {
				if err := polarDbService.RevokePolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}

		if len(add) > 0 {
			// wait instance running before granting
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polarDbService.PolarDbZonalClusterStateRefreshFunc(parts[0], []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			for _, db := range add {
				if err := polarDbService.GrantPolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}
	}

	return resourceAlicloudPolarDBOnENSAccountPrivilegeRead(d, meta)
}

func resourceAlicloudPolarDBOnENSAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := polarDBService.DescribePolarDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if object == nil {
		return nil
	}
	var dbName string

	if len(object.DatabasePrivileges) > 0 {
		for _, pri := range object.DatabasePrivileges {
			if pri.AccountPrivilege == parts[2] {
				dbName = pri.DBName
				if err := polarDBService.RevokePolarDBAccountPrivilege(d.Id(), pri.DBName); err != nil {
					return WrapError(err)
				}
			}
		}
	}

	return polarDBService.WaitForPolarDBAccountPrivilege(d.Id(), dbName, Deleted, DefaultTimeoutMedium)
}
