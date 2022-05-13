package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			},
			"account_privilege": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ReadWrite", "ReadOnly", "DMLOnly", "DDLOnly"}, false),
				Optional:     true,
				Computed:     true,
			},
			"collate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ctype": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudPolarDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	request := polardb.CreateCreateDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.DBName = d.Get("db_name").(string)
	request.CharacterSetName = d.Get("character_set_name").(string)

	if v, ok := d.GetOk("db_description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}
	request, errors := buildPolarDBDatabaseRequest(d, meta, request)
	if errors != nil {
		return WrapError(errors)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
	//wait for database creation to complete
	polarDBService := PolarDBService{client}
	stateConf := BuildStateConf([]string{""}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, polarDBService.PolarDBDatabaeStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
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
	d.Set("db_cluster_id", parts[0])
	d.Set("db_name", object.DBName)
	d.Set("character_set_name", strings.ToLower(object.CharacterSetName))
	d.Set("db_description", object.DBDescription)
	if object.Engine != "MySQL" {
		if len(object.Accounts.Account) > 0 {
			d.Set("account_name", object.Accounts.Account[0].AccountName)
		}
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

func buildPolarDBDatabaseRequest(d *schema.ResourceData, meta interface{}, request *polardb.CreateDatabaseRequest) (*polardb.CreateDatabaseRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	db_cluster_id := d.Get("db_cluster_id").(string)
	clusterAttribute, error := polarDBService.DescribePolarDBClusterAttribute(db_cluster_id)
	if error != nil {
		return nil, WrapError(error)
	}
	if clusterAttribute.DBType != "MySQL" {
		if v, ok := d.GetOk("account_name"); ok && v.(string) != "" {
			request.AccountName = v.(string)
		}
		if v, ok := d.GetOk("account_privilege"); ok && v.(string) != "" {
			request.AccountPrivilege = v.(string)
		}
		if v, ok := d.GetOk("collate"); ok && v.(string) != "" {
			request.Collate = v.(string)
		}
		if v, ok := d.GetOk("ctype"); ok && v.(string) != "" {
			request.Ctype = v.(string)
		}
	}
	return request, nil
}
