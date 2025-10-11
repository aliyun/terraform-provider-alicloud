package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBOnENSDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBOnENSDatabaseCreate,
		Read:   resourceAlicloudPolarDBOnENSDatabaseRead,
		Update: resourceAlicloudPolarDBOnENSDatabaseUpdate,
		Delete: resourceAlicloudPolarDBOnENSDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"db_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"character_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
				ForceNew: true,
			},

			"db_description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPolarDBOnENSDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	request := polardb.CreateCreateDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.DBName = d.Get("db_name").(string)
	request.CharacterSetName = d.Get("character_set_name").(string)

	cluster, err := polarDBService.DescribePolarDbZonalCluster(request.DBClusterId)

	if err != nil {
		return WrapError(err)
	}

	if cluster.DBType == "PostgreSQL" || cluster.DBType == "Oracle" {
		request.AccountName = d.Get("account_name").(string)
		request.Collate = "C"
		request.Ctype = "C"
	}

	if v, ok := d.GetOk("db_description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := polarDBService.CreateDatabaseZonal(request)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateDatabaseZonal", err, request)
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.DBName))
	if err := polarDBService.WaitForPolarDBDatabase(d.Id(), Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudPolarDBOnENSDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBOnENSDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	object, err := polarDBService.DescribePolarDBDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		d.SetId("")
		return nil
	}

	cluster, err := polarDBService.DescribePolarDbZonalCluster(parts[0])
	if err != nil {
		return WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("db_name", object.DBName)
	d.Set("character_set_name", strings.ToLower(object.CharacterSetName))
	d.Set("db_description", object.DBDescription)

	if cluster.DBType == "PostgreSQL" || cluster.DBType == "Oracle" {
		d.Set("account_name", object.Accounts[0].AccountName)
	}

	return nil
}

func resourceAlicloudPolarDBOnENSDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	if d.HasChange("db_description") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		dbDescription := d.Get("db_description").(string)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			err := polarDBService.ModifyDBDescriptionZonal(parts[0], parts[1], dbDescription)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("ModifyDBDescriptionZonal", err, parts[1])
			return nil
		})
		if err != nil {
			return err
		}
	}
	return resourceAlicloudPolarDBOnENSDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBOnENSDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDbServiceV2{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		err := polarDBService.DeleteDatabaseZonal(parts[0], parts[1])
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteDatabaseZonal", err, parts[1])
		return nil
	})
	if err != nil {
		return err
	}

	return WrapError(polarDBService.WaitForPolarDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
