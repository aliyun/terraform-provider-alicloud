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

func resourceAlicloudPolarDBDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBDatabaseCreate,
		Read:   resourceAlicloudPolarDBDatabaseRead,
		Update: resourceAlicloudPolarDBDatabaseUpdate,
		Delete: resourceAlicloudPolarDBDatabaseDelete,
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

			"account_privilege": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudPolarDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request := polardb.CreateCreateDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.DBName = d.Get("db_name").(string)
	request.CharacterSetName = d.Get("character_set_name").(string)

	cluster, err := polarDBService.DescribePolarDBCluster(request.DBClusterId)

	if err != nil {
		return WrapError(err)
	}

	if cluster.DBType == "PostgreSQL" || cluster.DBType == "Oracle" {
		request.AccountName = d.Get("account_name").(string)
		request.AccountPrivilege = d.Get("account_privilege").(string)
		request.Collate = "C"
		request.Ctype = "C"
	}

	if v, ok := d.GetOk("db_description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDatabase(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.DBName))
	if err := polarDBService.WaitForPolarDBDatabase(d.Id(), Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAlicloudPolarDBDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
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

	cluster, err := polarDBService.DescribePolarDBCluster(parts[0])
	if err != nil {
		return WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("db_name", object.DBName)
	d.Set("character_set_name", strings.ToLower(object.CharacterSetName))
	d.Set("db_description", object.DBDescription)

	if cluster.DBType == "PostgreSQL" || cluster.DBType == "Oracle" {
		d.Set("account_name", object.Accounts.Account[0].AccountName)
	}

	return nil
}

func resourceAlicloudPolarDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("db_description") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		request := polardb.CreateModifyDBDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBName = parts[1]
		request.DBDescription = d.Get("db_description").(string)
		var raw interface{}
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBDescription(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudPolarDBDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateDeleteDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBName = parts[1]
	// wait instance status is running before deleting database
	if err := polarDBService.WaitForPolarDBInstance(parts[0], Running, 1800); err != nil {
		return WrapError(err)
	}
	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DeleteDatabase(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(polarDBService.WaitForPolarDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
