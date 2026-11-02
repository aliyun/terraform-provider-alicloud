package polardb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribePolarDBCluster(id string) (instance *polardb.DBCluster, err error) {
	request := polardb.CreateDescribeDBClustersRequest()
	request.RegionId = s.client.RegionId
	dbClusterIds := []string{}
	dbClusterIds = append(dbClusterIds, id)
	request.DBClusterIds = id
	var response *polardb.DescribeDBClustersResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusters(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClustersResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	if len(response.Items.DBCluster) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *Service) DescribePolarDBClusterAttribute(id string) (instance *polardb.DescribeDBClusterAttributeResponse, err error) {
	request := polardb.CreateDescribeDBClusterAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	var response *polardb.DescribeDBClusterAttributeResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterAttribute(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterAttributeResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *Service) DescribeDBClusterAttribute(id string) (object map[string]interface{}, err error) {
	action := "DescribeDBClusterAttribute"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribePolarDBApplicationAttribute(id string) (object map[string]interface{}, err error) {
	action := "DescribeApplicationAttribute"
	request := map[string]interface{}{
		"ApplicationId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidParameter", "ResourceNotFound", "InvalidApplicationId.NotFound"}) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribePolarDBBatchTaskAttribute(id string) (object map[string]interface{}, err error) {
	action := "DescribeBatchTask"
	request := map[string]interface{}{
		"BatchId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidParameter", "ResourceNotFound", "InvalidBatchId.NotFound"}) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribePolarDBApplications(id string) (object map[string]interface{}, err error) {
	action := "DescribeApplications"
	request := map[string]interface{}{
		"RegionId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DeletePolarDBApplication(id string) (object map[string]interface{}, err error) {
	action := "DeleteApplication"
	request := map[string]interface{}{
		"ApplicationId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribePolarDBAutoRenewAttribute(id string) (instance *polardb.AutoRenewAttribute, err error) {
	request := polardb.CreateDescribeAutoRenewAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterIds = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeAutoRenewAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return instance, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeAutoRenewAttributeResponse)
	if len(response.Items.AutoRenewAttribute) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items.AutoRenewAttribute[0], nil
}

func (s *Service) DescribeParameters(id string) (ds *polardb.DescribeDBClusterParametersResponse, err error) {
	request := polardb.CreateDescribeDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.RetryContext(context.Background(), 5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterParameters(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				return retry.NonRetryableError(errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR))
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*polardb.DescribeDBClusterParametersResponse)
	return response, err
}

func (s *Service) GrantPolarDBAccountPrivilege(id, dbName string) error {
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	request := polardb.CreateGrantAccountPrivilegeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName
	request.AccountPrivilege = parts[2]

	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.GrantAccountPrivilege(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if err := s.WaitForPolarDBAccountPrivilege(id, dbName, Available, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func (s *Service) RevokePolarDBAccountPrivilege(id, dbName string) error {
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	request := polardb.CreateRevokeAccountPrivilegeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName

	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.RevokeAccountPrivilege(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if err := s.WaitForPolarDBAccountPrivilegeRevoked(id, dbName, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func (s *Service) WaitForPolarDBAccountPrivilegeRevoked(id, dbName string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			return errs.WrapError(err)
		}

		exist := false
		if object != nil {
			for _, dp := range object.DatabasePrivileges {
				if dp.DBName == dbName {
					exist = true
					break
				}
			}
		}

		if !exist {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, "", dbName, errs.ProviderERROR)
		}

	}
	return nil
}

func (s *Service) WaitForPolarDBAccountPrivilege(id, dbName string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		ready := false
		if object != nil {
			for _, dp := range object.DatabasePrivileges {
				if dp.DBName == dbName && dp.AccountPrivilege == parts[2] {
					ready = true
					break
				}
			}
		}
		if status == Deleted && !ready {
			break
		}
		if status != Deleted && ready {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, "", dbName, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) DescribePolarDBAccountPrivilege(id string) (account *polardb.DBAccount, err error) {
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := polardb.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := helper.NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *polardb.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeAccounts(request)
		})
		if err != nil {
			return errs.WrapError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBAccountPrivilege", id), NotFoundMsg, errs.ProviderERROR)
	}
	return &response.Accounts[0], nil
}

func (s *Service) WaitForPolarDBConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBConnectionV2(id, "Public")
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if status != Deleted && object != nil && object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			// object stays nil when the connection never became visible within the
			// wait window, so it must not be dereferenced in the timeout message.
			connectionString := ""
			if object != nil {
				connectionString = object.ConnectionString
			}
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, connectionString, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *Service) WaitPolardbEndpointConfigEffect(id string, item map[string]string, timeout int) error {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		effected := true
		object, err := s.DescribePolarDBInstanceNetInfo(parts[0])

		if err != nil {
			if errs.NotFoundError(err) {
				return errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapError(err)
		}

		var endpoint polardb.DBEndpoint
		if object != nil {
			for _, o := range object {
				if o.DBEndpointId == parts[1] {
					endpoint = o
					break
				}
			}
		}
		if value, ok := item["Nodes"]; ok {
			if endpoint.Nodes != value {
				effected = false
			}
		}
		if value, ok := item["ReadWriteMode"]; ok {
			if endpoint.ReadWriteMode != value {
				effected = false
			}
		}
		if value, ok := item["DBEndpointDescription"]; ok {
			if endpoint.DBEndpointDescription != value {
				effected = false
			}
		}
		if value, ok := item["AutoAddNewNodes"]; ok {
			if endpoint.AutoAddNewNodes != value {
				effected = false
			}
		}
		if value, ok := item["SccMode"]; ok {
			if !strings.EqualFold(endpoint.SccMode, value) {
				effected = false
			}
		}
		if value, ok := item["EndpointConfig"]; ok {
			expectConfig := make(map[string]string)
			actualConfig := make(map[string]string)
			err = json.Unmarshal([]byte(value), &expectConfig)
			err = json.Unmarshal([]byte(endpoint.EndpointConfig), &actualConfig)
			for k, v := range expectConfig {
				if subVal, ok := actualConfig[k]; !ok || subVal != v {
					effected = false
					break
				}
			}
		}
		if effected {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, endpoint, item, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *Service) WaitForPolarDBEndpoints(d *schema.ResourceData, status Status, endpointIds *schema.Set, timeout int) (string, error) {
	var dbEndpointId string
	if d.Id() != "" {
		parts, err := helper.ParseResourceId(d.Id(), 2)
		if err != nil {
			return "", errs.WrapError(err)
		}
		dbEndpointId = parts[1]
	}
	dbClusterId := d.Get("db_cluster_id").(string)
	endpointType := d.Get("endpoint_type").(string)

	newEndpoint := make(map[string]string)
	newEndpoint["endpoint_type"] = endpointType

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		endpoints, err := s.DescribePolarDBInstanceNetInfo(dbClusterId)
		if err != nil {
			return "", errs.WrapError(err)
		}

		var deleted bool
		deleted = true
		for _, value := range endpoints {
			if status == Deleted {
				if dbEndpointId == value.DBEndpointId {
					deleted = false
				}
				continue
			}
			if !endpointIds.Contains(value.DBEndpointId) && value.EndpointType == endpointType {
				return value.DBEndpointId, nil
			}
		}
		if status == Deleted && deleted {
			return "", nil
		}

		if time.Now().After(deadline) {
			return "", errs.WrapErrorf(err, WaitTimeoutMsg, dbClusterId, helper.GetFunc(1), timeout, endpoints, newEndpoint, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *Service) DescribePolarDBConnectionV2(id string, netType string) (*polardb.Address, error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(helper.DefaultIntervalLong) * time.Second)
	for {
		object, err := s.DescribePolarDBInstanceNetInfo(parts[0])

		if err != nil {
			if errs.NotFoundError(err) {
				return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return nil, errs.WrapError(err)
		}

		if object != nil {
			for _, o := range object {
				if o.DBEndpointId == parts[1] {
					for _, p := range o.AddressItems {
						if p.NetType == netType {
							return &p, nil
						}
					}
				}
			}
		}
		time.Sleep(helper.DefaultIntervalMini * time.Second)
		if time.Now().After(deadline) {
			break
		}
	}

	return nil, errs.WrapErrorf(errs.NotFoundErr("DBConnection", id), NotFoundMsg, errs.ProviderERROR)
}

func (s *Service) DescribePolarDBConnection(id string) (*polardb.Address, error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(helper.DefaultIntervalLong) * time.Second)
	for {
		object, err := s.DescribePolarDBInstanceNetInfo(parts[0])

		if err != nil {
			if errs.NotFoundError(err) {
				return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return nil, errs.WrapError(err)
		}

		if object != nil {
			for _, o := range object {
				if o.DBEndpointId == parts[1] {
					for _, p := range o.AddressItems {
						if p.NetType == "Public" {
							return &p, nil
						}
					}
				}
			}
		}
		time.Sleep(helper.DefaultIntervalMini * time.Second)
		if time.Now().After(deadline) {
			break
		}
	}

	return nil, errs.WrapErrorf(errs.NotFoundErr("DBConnection", id), NotFoundMsg, errs.ProviderERROR)
}

func (s *Service) DescribePolarDBInstanceNetInfo(id string) ([]polardb.DBEndpoint, error) {

	request := polardb.CreateDescribeDBClusterEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	var response *polardb.DescribeDBClusterEndpointsResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterEndpoints(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterEndpointsResponse)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(response.Items) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBInstanceNetInfo", id), NotFoundMsg, errs.ProviderERROR)
	}

	return response.Items, nil
}

func (s *Service) DescribePolarDBClusterEndpoint(id string) (*polardb.DBEndpoint, error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]

	request := polardb.CreateDescribeDBClusterEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = dbClusterId
	request.DBEndpointId = dbEndpointId

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterEndpoints(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				time.Sleep(10 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	response, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)
	if len(response.Items) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBEndpoint", dbEndpointId), NotFoundMsg, errs.ProviderERROR)
	}

	return &response.Items[0], nil
}

func (s *Service) DescribePolarDBClusterSSL(d *schema.ResourceData) (ssl *polardb.DescribeDBClusterSSLResponse, err error) {
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	dbClusterId := parts[0]

	request := polardb.CreateDescribeDBClusterSSLRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = dbClusterId
	var raw interface{}
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterSSL(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				time.Sleep(10 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	response, _ := raw.(*polardb.DescribeDBClusterSSLResponse)

	return response, nil
}

func (s *Service) DescribePolarDBDatabase(id string) (ds *polardb.Database, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	dbName := parts[1]
	request := polardb.CreateDescribeDatabasesRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.DBName = dbName

	var raw interface{}
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDatabases(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				time.Sleep(10 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDatabasesResponse)
	if len(response.Databases.Database) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBDatabase", dbName), NotFoundMsg, errs.ProviderERROR)
	}
	ds = &response.Databases.Database[0]
	return ds, nil
}

func (s *Service) WaitForPolarDBDatabase(id string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBDatabase(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
				if status == Running {
					continue
				}
			}
			return errs.WrapError(err)
		}
		if status != Deleted && object != nil && object.DBName == parts[1] {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.DBName, parts[1], errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForPolarDBAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccount(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.AccountStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) DescribePolarDBAccount(id string) (ds *polardb.DBAccount, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := polardb.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := helper.NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *polardb.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*polardb.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("DBAccount", id), NotFoundMsg, errs.ProviderERROR)
	}
	return &response.Accounts[0], nil
}

// WaitForInstance waits for instance to given status
func (s *Service) WaitForPolarDBInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.DBClusterStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForPolarDBConnectionPrefix(id, prefix, newPort string, netType string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBConnectionV2(id, netType)
		if err != nil {
			return errs.WrapError(err)
		}
		parts := strings.Split(object.ConnectionString, ".")
		port := object.Port

		if (newPort == "" || newPort == port) && (prefix == "" || prefix == parts[0]) {
			break
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, parts[0], prefix, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) fillingPolarDBEndpointSslCertificateUrl(sslEnabled string, d *schema.ResourceData) {
	if sslEnabled == "Enable" {
		d.Set("ssl_certificate_url", "https://apsaradb-public.oss-ap-southeast-1.aliyuncs.com/ApsaraDB-CA-Chain.zip?file=ApsaraDB-CA-Chain.zip&regionId="+s.client.RegionId)
	} else {
		d.Set("ssl_certificate_url", "")
	}
}

func (s *Service) RefreshEndpointConfig(d *schema.ResourceData) error {
	var config map[string]interface{}
	config = make(map[string]interface{}, 0)
	documented, ok := d.GetOk("endpoint_config")
	if !ok {
		d.Set("endpoint_config", config)
		return nil
	}
	object, err := s.DescribePolarDBClusterEndpoint(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	var endpointConfig = make(map[string]interface{})
	err = json.Unmarshal([]byte(object.EndpointConfig), &endpointConfig)
	if err != nil {
		return errs.WrapError(err)
	}

	for k, v := range documented.(map[string]interface{}) {
		if _, ok := endpointConfig[k]; ok {
			config[k] = v
		}
	}
	if err := d.Set("endpoint_config", config); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

// lintignore: R001
func (s *Service) RefreshParameters(d *schema.ResourceData) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk("parameters")
	object, err := s.DescribeParameters(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	var parameters = make(map[string]interface{})
	for _, i := range object.RunningParameters.Parameter {
		// 创建集群传入参数模板参数
		changeParams := []string{"loose_polar_log_bin", "lower_case_table_names", "default_time_zone", "loose_xengine", "loose_xengine_use_memory_pct"}
		if IsContain(changeParams, i.ParameterName) {
			if i.ParameterName == "lower_case_table_names" {
				if _parameterValue, err := strconv.Atoi(i.ParameterValue); err == nil {
					d.Set(i.ParameterName, _parameterValue)
				}
			} else if d.Get("db_type").(string) == "MySQL" && d.Get("db_version") == "5.6" && i.ParameterName == "loose_polar_log_bin" {
				// mysql 5.6 loose_polar_log_bin values: ON_WITH_GTID、OFF
				if i.ParameterValue == "ON_WITH_GTID" {
					d.Set(i.ParameterName, "ON")
				}
			} else if i.ParameterName == "loose_xengine" {
				if i.ParameterValue == "1" {
					d.Set("loose_xengine", "ON")
				} else {
					d.Set("loose_xengine", "OFF")
				}
			} else if i.ParameterName == "loose_xengine_use_memory_pct" {
				if parameterValue, err := strconv.Atoi(i.ParameterValue); err == nil {
					d.Set(i.ParameterName, parameterValue)
				}
			} else {
				d.Set(i.ParameterName, i.ParameterValue)
			}
		}
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}
	if !ok {
		d.Set("parameters", param)
		return nil
	}
	if err := d.Set("parameters", param); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *Service) ModifyParameters(d *schema.ResourceData) error {
	request := polardb.CreateModifyDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = d.Id()
	config := make(map[string]string)
	allConfig := make(map[string]string)
	o, n := d.GetChange("parameters")
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request.Parameters = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForCluster(d.Id(), Running, DefaultLongTimeout); err != nil {
			return errs.WrapError(err)
		}
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterParameters(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance parameter expect after modifying
		for _, i := range ns.List() {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			allConfig[key] = value
		}
		if err := s.WaitForPolarDBParameter(d.Id(), DefaultLongTimeout, allConfig); err != nil {
			return errs.WrapError(err)
		}
	}
	return nil
}

func (s *Service) CreateClusterParamsModifyParameters(d *schema.ResourceData) error {
	request := polardb.CreateModifyDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = d.Id()
	config := make(map[string]interface{})
	allConfig := make(map[string]string)
	changeParams := []string{"loose_polar_log_bin", "default_time_zone", "loose_xengine", "loose_xengine_use_memory_pct"}
	for _, i := range changeParams {
		if v, ok := d.GetOk(i); ok {
			if d.HasChange(i) {
				if i == "loose_xengine" {
					if v == "ON" {
						v = "1"
					} else if v == "OFF" {
						v = "0"
					}
				}
				config[i] = v
				allConfig[i] = fmt.Sprint(v)
			}
		}
	}
	cfg, _ := json.Marshal(config)
	request.Parameters = string(cfg)
	// wait instance status is Normal before modifying
	if err := s.WaitForCluster(d.Id(), Running, DefaultLongTimeout); err != nil {
		return errs.WrapError(err)
	}
	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.ModifyDBClusterParameters(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	// wait instance parameter expect after modifying
	if err := s.WaitForPolarDBParameter(d.Id(), 1200, allConfig); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *Service) setClusterTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				if !tags.IgnoredTags(v.Key, v.Value) {
					tagKey = append(tagKey, v.Key)
				}
			}
			request := polardb.CreateUntagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = "cluster"
			request.TagKey = &tagKey
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := polardb.CreateTagResourcesRequest()
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = "cluster"
			request.RegionId = s.client.RegionId
			raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

	}

	return nil
}

func (s *Service) diffTags(oldTags, newTags []polardb.TagResourcesTag) ([]polardb.TagResourcesTag, []polardb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []polardb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *Service) tagsToMap(tags []polardb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *Service) tagsFromMap(m map[string]interface{}) []polardb.TagResourcesTag {
	result := make([]polardb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, polardb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *Service) ignoreTag(t polardb.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *Service) DescribeTags(resourceId string, resourceType TagResourceType) (tags []polardb.TagResource, err error) {
	request := polardb.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

// WaitForCluster waits for cluster to given status
func (s *Service) WaitForCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.DBClusterStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) DescribeDBSecurityIps(clusterId string, dbClusterIPArrayName string) (ips []string, err error) {

	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAccessWhitelist(request)
	})
	if err != nil {
		return ips, errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range resp.Items.DBClusterIPArray {
		if ip.DBClusterIPArrayName == dbClusterIPArrayName {
			ipstr += separator + ip.SecurityIps
			separator = COMMA_SEPARATED
		}
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *Service) ModifyDBSecurityIps(clusterId, ips string) error {

	request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.SecurityIps = ips

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.ModifyDBClusterAccessWhitelist(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *Service) DescribeBackupPolicy(id string) (object map[string]interface{}, err error) {

	var response map[string]interface{}
	client := s.client
	action := "DescribeBackupPolicy"
	request := map[string]interface{}{
		"DBClusterId": id,
		"RegionId":    s.client.RegionId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeLogBackupPolicy(id string) (object map[string]interface{}, err error) {

	var response map[string]interface{}
	client := s.client
	action := "DescribeLogBackupPolicy"
	request := map[string]interface{}{
		"DBClusterId": id,
		"RegionId":    s.client.RegionId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) ModifyDBBackupPolicy(clusterId, backupTime, backupPeriod, backupRetentionPolicyOnClusterDeletion string) error {

	request := polardb.CreateModifyBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime
	request.BackupRetentionPolicyOnClusterDeletion = backupRetentionPolicyOnClusterDeletion

	raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.ModifyBackupPolicy(request)
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func (s *Service) DescribeDBAuditLogCollectorStatus(id string) (collectorStatus string, err error) {
	request := polardb.CreateDescribeDBClusterAuditLogCollectorRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	var response *polardb.DescribeDBClusterAuditLogCollectorResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusterAuditLogCollector(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response, _ = raw.(*polardb.DescribeDBClusterAuditLogCollectorResponse)
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return "", errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return collectorStatus, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return response.CollectorStatus, nil
}

func (s *Service) PolarDBClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBClusterStatus == failState {
				return object, object.DBClusterStatus, errs.WrapError(errs.Error(FailedToReachTargetStatus, object.DBClusterStatus))
			}
		}
		return object, object.DBClusterStatus, nil
	}
}

func (s *Service) PolarDBApplicationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBApplicationAttribute(id)
		if err != nil {
			// 当资源被删除后，API 可能返回各种错误，我们应该将这些错误都视为 "资源不存在" 的情况
			// 包括 NotFoundError 和其他可能的错误类型
			if errs.NotFoundError(err) || strings.Contains(err.Error(), "InvalidParameter") || strings.Contains(err.Error(), "ResourceNotFound") {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"] == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *Service) PolarDBBatchTaskStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBBatchTaskAttribute(id)
		if err != nil {
			// 当资源被删除后，API 可能返回各种错误，我们应该将这些错误都视为 "资源不存在" 的情况
			// 包括 NotFoundError 和其他可能的错误类型
			if errs.NotFoundError(err) || strings.Contains(err.Error(), "InvalidParameter") || strings.Contains(err.Error(), "ResourceNotFound") {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"] == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, object["Status"].(string), nil
	}
}

// WaitForDBParameter waits for instance parameter to given value.
// Status of DB instance is Running after ModifyParameters API was
// call, so we can not just wait for instance status become
// Running, we should wait until parameters have expected values.
func (s *Service) WaitForPolarDBParameter(clusterId string, timeout int, expects map[string]string) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	// get engine and engineVersion
	attrbuteInfo, err := s.DescribePolarDBClusterAttribute(clusterId)
	if err != nil {
		return errs.WrapError(err)
	}
	dbType := fmt.Sprint(attrbuteInfo.DBType)
	db_version := fmt.Sprint(attrbuteInfo.DBVersion)

	for {
		object, err := s.DescribeParameters(clusterId)
		if err != nil {
			return errs.WrapError(err)
		}

		var actuals = make(map[string]string)
		for _, i := range object.RunningParameters.Parameter {
			// mysql 5.6 loose_polar_log_bin values: ON_WITH_GTID、OFF
			if dbType == "MySQL" && db_version == "5.6" && i.ParameterName == "loose_polar_log_bin" {
				if i.ParameterValue == "ON_WITH_GTID" {
					i.ParameterValue = "ON"
				}
			}
			actuals[i.ParameterName] = i.ParameterValue
		}

		match := true

		got_value := ""
		expected_value := ""

		for name, expect := range expects {
			if actual, ok := actuals[name]; ok {
				if expect != actual {
					match = false
					got_value = actual
					expected_value = expect
					break
				}
			} else {
				match = false
			}
		}

		if match {
			break
		}

		time.Sleep(helper.DefaultIntervalShort * time.Second)

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, clusterId, helper.GetFunc(1), timeout, got_value, expected_value, errs.ProviderERROR)
		}
	}
	return nil
}

func (s *Service) DescribeDBClusterTDE(id string) (map[string]interface{}, error) {
	action := "DescribeDBClusterTDE"
	request := map[string]interface{}{
		"DBClusterId": id,
	}
	client := s.client
	response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	return response, nil
}

func (s *Service) CheckKMSAuthorized(id string, tdeRegion string) (map[string]interface{}, error) {
	action := "CheckKMSAuthorized"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	if tdeRegion != "" {
		request["tdeRegion"] = tdeRegion
	}
	var response map[string]interface{}
	var err error
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return response, nil
}

func (s *Service) WaitForPolarDBTDEStatus(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Minute)
	for {
		object, err := s.DescribeDBClusterTDE(id)
		if err != nil {
			return errs.WrapError(err)
		}
		if object["TDEStatus"] == status {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *Service) WaitForPolarDBNodeClass(id string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		clusters, err := s.DescribePolarDBCluster(id)
		if err != nil {
			return errs.WrapError(err)
		}
		clusterAttribute, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			return errs.WrapError(err)
		}
		if len(clusters.DBNodes.DBNode) == len(clusterAttribute.DBNodes) {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *Service) WaitForPolarDBPayType(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		clusters, err := s.DescribePolarDBCluster(id)
		if err != nil {
			return errs.WrapError(err)
		}
		if clusters.PayType == status {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, clusters, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalMedium * time.Second)
	}
	return nil
}

func (s *Service) PolarDBClusterTDEStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		object, err := s.DescribeDBClusterTDE(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["TDEStatus"].(string) == failState {
				return object, object["TDEStatus"].(string), errs.WrapError(errs.Error(FailedToReachTargetStatus, object["TDEStatus"].(string)))
			}
		}
		return object, object["TDEStatus"].(string), nil
	}
}

func (s *Service) DescribeDBSecurityGroups(clusterId string) ([]string, error) {

	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAccessWhitelist(request)
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)

	groups := make([]string, 0)
	dbClusterSecurityGroups := resp.DBClusterSecurityGroups.DBClusterSecurityGroup
	for _, group := range dbClusterSecurityGroups {

		groups = append(groups, group.SecurityGroupId)
	}
	return groups, nil
}

func (s *Service) ModifyDBClusterAccessWhitelist(d *schema.ResourceData) error {
	if _, ok := d.GetOk("db_cluster_ip_array"); ok {
		removed, added := d.GetChange("db_cluster_ip_array")

		targetIpArrayNames := make(map[string]bool)
		for _, e := range added.(*schema.Set).List() {
			pack := e.(map[string]interface{})
			ipArrayName := pack["db_cluster_ip_array_name"].(string)
			targetIpArrayNames[ipArrayName] = true
		}

		for _, e := range removed.(*schema.Set).List() {
			pack := e.(map[string]interface{})
			if fmt.Sprint(pack["db_cluster_ip_array_name"]) == "default" || targetIpArrayNames[pack["db_cluster_ip_array_name"].(string)] {
				continue
			}
			//ips expand string list
			ipList := convert.ExpandStringList(pack["security_ips"].(*schema.Set).List())
			ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
			request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
			request.RegionId = s.client.RegionId
			request.DBClusterId = d.Id()
			request.SecurityIps = ipstr
			request.DBClusterIPArrayName = pack["db_cluster_ip_array_name"].(string)
			request.ModifyMode = "Delete"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(5*time.Minute, func() *retry.RetryError {
				raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
					return polarDBClient.ModifyDBClusterAccessWhitelist(request)
				})
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			if err := s.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
				return errs.WrapError(err)
			}
		}
		for _, e := range added.(*schema.Set).List() {
			pack := e.(map[string]interface{})
			//ips expand string list
			ipList := convert.ExpandStringList(pack["security_ips"].(*schema.Set).List())
			ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
			// default disable connect from outside
			if ipstr == "" {
				ipstr = LOCAL_HOST_IP
			}
			request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
			request.RegionId = s.client.RegionId
			request.DBClusterId = d.Id()
			request.SecurityIps = ipstr
			request.DBClusterIPArrayName = pack["db_cluster_ip_array_name"].(string)
			request.ModifyMode = pack["modify_mode"].(string)
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err := retry.Retry(5*time.Minute, func() *retry.RetryError {
				raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
					return polarDBClient.ModifyDBClusterAccessWhitelist(request)
				})
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			if err := s.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	return nil
}

func (s *Service) DescribeDBClusterAccessWhitelist(id string) (instance *polardb.DescribeDBClusterAccessWhitelistResponse, err error) {
	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	var response *polardb.DescribeDBClusterAccessWhitelistResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterAccessWhitelist(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)
		return nil
	})

	if err != nil {
		return nil, errs.WrapError(err)
	}
	return response, nil
}

func convertPolarDBIpsSetListToString(arr1 *schema.Set) []string {
	var ips []string
	for _, v := range arr1.List() {
		ips = append(ips, v.(string))
	}
	return ips
}

func convertPolarDBIpsSetToString(sourceIps string) []string {
	ipsMap := make(map[string]string)

	for _, ip := range strings.Split(sourceIps, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}
	var ips []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			ips = append(ips, key)
		}
	}
	return ips
}

func arrValueEqual(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	sort.Strings(arr2)
	sort.Strings(arr1)
	for i, v := range arr1 {
		if v != arr2[i] {
			return false
		}
	}
	return true
}

func (s *Service) PolarDBClusterCategoryRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Category == failState {
				return object, object.Category, errs.WrapError(errs.Error(FailedToReachTargetStatus, object.Category))
			}
		}
		return object, object.Category, nil
	}
}

func (s *Service) DescribePolarDBGlobalDatabaseNetwork(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeGlobalDatabaseNetwork"
	request := map[string]interface{}{
		"GDNId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"GDN.NotFound"}) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) PolarDBGlobalDatabaseNetworkRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		object, err := s.DescribePolarDBGlobalDatabaseNetwork(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["GDNStatus"]) == failState {
				return object, fmt.Sprint(object["GDNStatus"]), errs.WrapError(errs.Error(FailedToReachTargetStatus, fmt.Sprint(object["GDNStatus"])))
			}
		}
		return object, fmt.Sprint(object["GDNStatus"]), nil
	}
}

func (s *Service) DescribePolarDBParameterGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeParameterGroup"
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"ParameterGroupId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParamGroupsNotExist"}) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ParameterGroup", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ParameterGroup", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("PolarDB", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ParameterGroupId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("PolarDB", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeDBClusterServerlessConfig(id string) (object map[string]interface{}, err error) {
	action := "DescribeDBClusterServerlessConf"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribeDBClusterVersion(id string) (object map[string]interface{}, err error) {
	action := "DescribeDBClusterVersion"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}
func (s *Service) DescribeDBClusterAvailableVersion(id string) (instance *polardb.DescribeDBClusterVersionResponse, err error) {
	request := polardb.CreateDescribeDBClusterVersionRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	request.DescribeType = "AVAILABLE_VERSION"

	var response *polardb.DescribeDBClusterVersionResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusterVersion(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeDBClusterVersionResponse)
		return nil
	})

	if err != nil {
		return nil, errs.WrapError(err)
	}
	return response, nil
}

func (s *Service) PolarDBClusterProxyStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.ProxyStatus == failState {
				return object, object.ProxyStatus, errs.WrapError(errs.Error(FailedToReachTargetStatus, object.ProxyStatus))
			}
		}
		return object, object.ProxyStatus, nil
	}
}
func (s *Service) DescribeDBClusterStandbyAz(d *schema.ResourceData, timeout int) (zoneId string, err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	standbyAz := ""
	for {
		clusterAttribute, err := s.DescribePolarDBClusterAttribute(d.Id())
		if err != nil {
			return "", errs.WrapError(err)
		}

		exist := true
		if clusterAttribute.HotStandbyCluster != "StandbyClusterOFF" {
			exist = false
			zoneIds := strings.Split(clusterAttribute.ZoneIds, ",")
			if len(zoneIds) >= 2 {
				zoneId = zoneIds[0]
				standbyAz = zoneIds[1]
				// standbyAz info for new cluster may be incorrectly set to be equal to zoneId at first, need retry for real value
				if standbyAz != "" && (standbyAz != zoneId || !d.IsNewResource()) {
					exist = true
				}
			}
		}

		if exist {
			return standbyAz, nil
		}
		if time.Now().After(deadline) {
			// adapt for new cluster with exactly the same zoneId and standbyAz
			if standbyAz != "" {
				return standbyAz, nil
			}
			return standbyAz, errs.WrapErrorf(err, RequiredWhenMsg, "standby_az", "hot_standby_cluster", clusterAttribute.HotStandbyCluster)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *Service) DescribePolarDBGlobalSecurityGroupIds(id string) (globalSecurityGroupIds []string, err error) {
	action := "DescribeGlobalSecurityIPGroupRelation"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return globalSecurityGroupIds, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.GlobalSecurityIPGroupRel", response)
	if err != nil {
		return globalSecurityGroupIds, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.GlobalSecurityIPGroupRel", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return globalSecurityGroupIds, nil
	}

	for _, v := range resp.([]interface{}) {
		globalSecurityGroupIds = append(globalSecurityGroupIds, fmt.Sprint(v.(map[string]interface{})["GlobalSecurityGroupId"]))
	}

	return globalSecurityGroupIds, nil
}

func (s *Service) DescribePolarDBAIClusterAttribute(id string) (object map[string]interface{}, err error) {
	action := "DescribeAIDBClusterAttribute"
	request := map[string]interface{}{
		"DBClusterId": id,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidParameter", "ResourceNotFound"}) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribePolarDBGatewayAttribute(id string) (object map[string]interface{}, err error) {
	action := "DescribeGatewayAttribute"
	request := map[string]interface{}{"GwClusterId": id, "RegionId": s.client.RegionId}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.NotFoundError(err) {
			return object, errs.WrapErrorf(err, NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *Service) PolarDBGatewayStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBGatewayAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return "DELETED", "DELETED", nil
			}
			return nil, "", errs.WrapError(err)
		}
		status := fmt.Sprint(object["Status"])
		for _, failState := range failStates {
			if status == failState {
				return object, status, errs.WrapError(errs.Error(FailedToReachTargetStatus, status))
			}
		}
		return object, status, nil
	}
}

func (s *Service) PolarDBAIClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBAIClusterAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		status, _ := object["DBClusterStatus"].(string)

		for _, failState := range failStates {
			if status == failState {
				return object, status, errs.WrapError(errs.Error(FailedToReachTargetStatus, status))
			}
		}
		return object, status, nil
	}
}

func (s *Service) DescribeAIDBClusterApiKeys(regionId string) (object map[string]interface{}, err error) {
	action := "DescribeAIDBClusterApiKeys"
	request := map[string]interface{}{
		"RegionId": regionId,
	}
	var response map[string]interface{}
	client := s.client
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, "", action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "", "$", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}

// DescribePolarDBDynamoEndpointAddress resolves the DynamoDB-compatible endpoint address
// (http://<connection_string>:<port>) of a PolarDB cluster. It prefers the Public address
// and falls back to the first available one. It is mainly used on resource import, where
// only the resource ID (which contains the db_cluster_id) is available.
func (s *Service) DescribePolarDBDynamoEndpointAddress(dbClusterId string) (string, error) {
	endpoints, err := s.DescribePolarDBInstanceNetInfo(dbClusterId)
	if err != nil {
		return "", errs.WrapError(err)
	}
	fallback := ""
	for _, endpoint := range endpoints {
		if endpoint.EndpointType != "DynamoDB" {
			continue
		}
		for _, address := range endpoint.AddressItems {
			if address.ConnectionString == "" {
				continue
			}
			addr := fmt.Sprintf("http://%s:%s", address.ConnectionString, address.Port)
			if address.NetType == "Public" {
				return addr, nil
			}
			if fallback == "" {
				fallback = addr
			}
		}
	}
	if fallback != "" {
		return fallback, nil
	}
	return "", errs.WrapErrorf(errs.NotFoundErr("PolarDBDynamoEndpoint", dbClusterId), NotFoundMsg, errs.ProviderERROR)
}

// DescribePolarDBDynamoAccount finds the DynamoDB-type account of a PolarDB cluster and
// returns its account name and DynamoDB auth password. The DynamoDB-compatible endpoint
// only accepts this account (the provider AK/SK are rejected), so it is used on resource
// import, where only the resource ID (which contains the db_cluster_id) is available.
func (s *Service) DescribePolarDBDynamoAccount(dbClusterId string) (string, string, error) {
	action := "DescribeAccounts"
	request := map[string]interface{}{
		"DBClusterId": dbClusterId,
	}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = s.client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return "", "", errs.WrapErrorf(err, errs.DefaultErrorMsg, dbClusterId, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Accounts[*]", response)
	if err != nil {
		return "", "", errs.WrapErrorf(errs.NotFoundErr("PolarDBDynamoAccount", dbClusterId), NotFoundMsg, response)
	}
	for _, raw := range v.([]interface{}) {
		account := raw.(map[string]interface{})
		if fmt.Sprint(account["AccountType"]) != "DynamoDB" {
			continue
		}
		accountName := fmt.Sprint(account["AccountName"])
		authPassword := ""
		if p, ok := account["DynamoDBAuthPassword"]; ok && p != nil {
			authPassword = fmt.Sprint(p)
		}
		if accountName != "" && authPassword != "" {
			return accountName, authPassword, nil
		}
	}
	return "", "", errs.WrapErrorf(errs.NotFoundErr("PolarDBDynamoAccount", dbClusterId), NotFoundMsg, response)
}

// PolarDBDynamoTableStateRefreshFunc returns a StateRefreshFunc for DynamoDB table status
func (s *Service) PolarDBDynamoTableStateRefreshFunc(dbClusterId, tableName, endpoint, accessKey, secretKey string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		dynamoClient, err := s.client.NewPolarDBDynamoClient(endpoint, accessKey, secretKey)
		if err != nil {
			return nil, "", errs.WrapError(err)
		}

		output, err := dynamoClient.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			if isDynamoNotFoundError(err) {
				return nil, "", nil
			}
			if isDynamoRetryableError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapErrorf(err, errs.DefaultErrorMsg, fmt.Sprintf("%s:%s", dbClusterId, tableName), "DescribeTable", errs.AlibabaCloudSdkGoERROR)
		}

		if output.Table == nil {
			return nil, "", fmt.Errorf("missing Table field in DescribeTable response")
		}

		status := string(output.Table.TableStatus)
		if status == "" {
			status = "ACTIVE"
		}

		return output.Table, status, nil
	}
}

type fwTimeoutDurationValidator struct{}

type fwWireField struct {
	Name   string
	Single bool
	Read   bool
	Input  bool
	Fields map[string]fwWireField
}

type fwMapping struct {
	Source, Target, Encoding, Constant string
	Group, GroupKey, GroupValue        string
	Repeat                             bool
	IndexSources                       []int
	Values                             map[string]string
}

type fwMappedValue struct {
	Value   interface{}
	Indexes []int
}

type fwScopedReadIdentity struct{ Source, Target, Value string }

func fwTimeoutTypes() map[string]attr.Type {
	return map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}
}

func fwTimeoutBlock() fwschema.SingleNestedBlock {
	attributes := map[string]fwschema.Attribute{}
	for _, operation := range []string{"create", "read", "update", "delete"} {
		attributes[operation] = fwschema.StringAttribute{
			Optional:    true,
			Description: "Maximum retry and wait duration for this operation, for example 10m.",
			Validators:  []validator.String{fwTimeoutDurationValidator{}},
		}
	}
	return fwschema.SingleNestedBlock{Attributes: attributes}
}

func (fwTimeoutDurationValidator) Description(context.Context) string {
	return "A positive Go duration, for example 10m or 1h."
}
func (v fwTimeoutDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (fwTimeoutDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	duration, err := time.ParseDuration(req.ConfigValue.ValueString())
	if err != nil || duration <= 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
	}
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if value.IsNull() {
		return fallback, diagnostics
	}
	if value.IsUnknown() {
		diagnostics.AddError("Unknown operation timeout", "Resolve the timeouts block before applying.")
		return 0, diagnostics
	}
	raw, ok := value.Attributes()[operation]
	if !ok || raw.IsNull() {
		return fallback, diagnostics
	}
	configured, ok := raw.(types.String)
	if !ok || configured.IsUnknown() {
		diagnostics.AddError("Invalid operation timeout", "Resolve the operation timeout to a duration string before applying.")
		return 0, diagnostics
	}
	duration, err := time.ParseDuration(configured.ValueString())
	if err != nil || duration <= 0 {
		diagnostics.AddError("Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
		return 0, diagnostics
	}
	return duration, diagnostics
}

// retryRpcPost wraps RpcPost with the shared retry policy: helper.NeedRetry covers
// throttling, 5xx and client timeouts, and backupTransientErrors adds the PolarDB status
// codes that clear on their own.
func retryRpcPost(client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.Retry(timeout, func() *retry.RetryError {
		resp, e := client.RpcPost(product, version, action, nil, request, true)
		response = resp
		if e != nil {
			if helper.NeedRetry(e) || backupTransientError(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		return response, errs.WrapErrorf(err, "calling %s failed", action)
	}
	return response, nil
}

// backupTransientCodes are the PolarDB status codes for which RpcPost can be retried.
var backupTransientCodes = []string{"OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy"}

func backupTransientError(err error) bool {
	return errs.IsExpectedErrors(err, backupTransientCodes)
}

// backupReadNotFoundCodes are the NotFound error codes returned by DescribeBackupPolicy when
// the cluster (and therefore its backup) no longer exists. F5: the empty codes list previously
// never recognized NotFound, so an externally deleted cluster surfaced as a diagnostic instead
// of clearing state.
var backupReadNotFoundCodes = []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}

// backupDeleteNotFoundCodes are the NotFound error codes returned by DeleteBackup when the
// backup set or cluster is already gone. F5: Delete only accepted generic ResourceNotFound
// codes and rejected the official InvalidBackupSetID.NotFound / InvalidDBClusterId.NotFound.
var backupDeleteNotFoundCodes = []string{"InvalidBackupSetID.NotFound", "InvalidDBClusterId.NotFound", "ResourceNotFound.Backup", "ResourceNotFound"}

// backupParseID splits the composite resource id "DBClusterId:BackupId" into its parts.
// F1: the stable identity of a backup set is the BackupId returned by CreateBackup (the
// BackupJobId, which is the backup set id). Pairing it with the owning DBClusterId lets Read
// find the exact backup and Delete target only the owned backup instead of "last wins".
func backupParseID(id string) (dbClusterId, backupId string) {
	if idx := strings.Index(id, ":"); idx >= 0 {
		return id[:idx], id[idx+1:]
	}
	return id, ""
}

// backupValidateID accepts only the composite "DBClusterId:BackupId" form. This resource was
// introduced with the composite id, so a single-segment id has no legitimate user: importing
// one would record a backup without a backup-set identity, which Read would always report as
// present and Delete could never target (DeleteBackup requires BackupId).
func backupValidateID(id string) (dbClusterId, backupId string, err error) {
	dbClusterId, backupId = backupParseID(id)
	if dbClusterId == "" || backupId == "" || strings.Contains(backupId, ":") {
		return "", "", fmt.Errorf("invalid backup id %q: expected <db_cluster_id>:<backup_id>", id)
	}
	return dbClusterId, backupId, nil
}

// errBackupResourceGone marks a resource whose cluster and backup set are both gone. It is
// distinct from a nil object so the create path can keep retrying a cluster that is briefly
// invisible, while Read can stop immediately for a genuinely deleted backup.
var errBackupResourceGone = errors.New("polardb backup resource no longer exists")

// errBackupReadInconclusive marks a read in which the policy API reports the cluster as absent,
// the live list has no match, and the released-cluster backup list cannot answer (throttling,
// 5xx or timeout). The backup set may still exist there, so the caller keeps the current state
// instead of removing or failing the resource.
var errBackupReadInconclusive = errors.New("backup set state cannot be determined: the released-cluster backup list is unavailable")

// backupComposeID builds the composite resource id "DBClusterId:BackupId".
func backupComposeID(dbClusterId, backupId string) string {
	return dbClusterId + ":" + backupId
}

// backupDetachedReadKey marks a read object that came from DescribeDetachedBackups, i.e. a
// backup set that outlived its released cluster. Such an object carries backup-set fields
// only: the cluster-level backup policy no longer exists and is preserved from state instead.
const backupDetachedReadKey = "__detached_backup_set"

// backupGet describes the resource by id. Returns (nil, nil) when the backup set is gone and
// errBackupResourceGone when both the cluster policy and the detached backup list report it
// absent.
func backupGet(client *connectivity.AliyunClient, id string, includeSecondary ...bool) (map[string]interface{}, error) {
	return backupGetWithStart(client, id, "", len(includeSecondary) == 0 || includeSecondary[0])
}

// backupGetWithStart describes the resource by id, using knownStart (the BackupStartTime from
// state) to narrow the backup-set lookup. DescribeBackupPolicy is the primary source because
// it carries both the cluster policy and the identity confirmation; when it reports the cluster
// as absent the backup set may still exist: releasing a cluster with
// BackupRetentionPolicyOnClusterDeletion=ALL keeps its backup sets, and DescribeDetachedBackups
// is the API that reports exactly those sets. Clearing state on the policy NotFound alone would
// drop a live, billed backup from Terraform (P1).
func backupGetWithStart(client *connectivity.AliyunClient, id, knownStart string, completeRead bool) (map[string]interface{}, error) {
	// F1: the resource id is the composite "DBClusterId:BackupId". DescribeBackupPolicy is
	// keyed on DBClusterId, so resolve the cluster part; the stable BackupId is injected into
	// the read object below so DescribeBackups/DescribeBackupTasks filter on the exact backup
	// instead of merging "last wins" over every backup set in the cluster.
	dbClusterId, backupId, idErr := backupValidateID(id)
	if idErr != nil {
		return nil, idErr
	}
	getValues := map[string]interface{}{"RegionId": client.RegionId, "DBClusterId": dbClusterId}
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	if err := fwApplyMappings(request, getValues, []fwMapping{
		{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		return nil, fmt.Errorf("building DescribeBackupPolicy request: %w", err)
	}
	if value, ok := request["DBClusterId"].(string); !ok || value == "" || value != dbClusterId {
		return nil, fmt.Errorf("building DescribeBackupPolicy request: required identity does not match")
	}
	response, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeBackupPolicy", request, 1*time.Minute)
	if err != nil {
		if !errs.IsExpectedErrors(err, backupReadNotFoundCodes) {
			return nil, err
		}
		// P1: the cluster is gone from the policy API, but the backup set may still exist. Both
		// remaining sources are consulted before concluding that it does not:
		//  1. a live cluster can race the policy read (immediately after creation the policy API
		//     reports the cluster as absent), so the live list is asked first;
		//  2. releasing a cluster with BackupRetentionPolicyOnClusterDeletion=ALL keeps its
		//     backup sets, and DescribeDetachedBackups is the API that reports those.
		// Only a released cluster whose detached list has no exact match proves the resource is
		// gone; a still-live cluster keeps retrying through backupReadExisting.
		liveSet, _, liveErr := backupFindLiveBackupSetByID(client, dbClusterId, backupId, knownStart, true)
		if liveErr != nil {
			return nil, liveErr
		}
		if liveSet != nil {
			return backupDetachedReadObject(dbClusterId, backupId, liveSet), nil
		}
		set, sawEntries, detachedErr := backupFindDetachedBackupSetByID(client, dbClusterId, backupId, knownStart)
		if errors.Is(detachedErr, errBackupClusterNotReleased) {
			// DescribeDetachedBackups could not describe the cluster: either it still exists
			// (both lists answered without a match) or its record is gone altogether, which
			// takes the backup set with it. Both are a completed enumeration without a match.
			return nil, errBackupResourceGone
		}
		if detachedErr != nil {
			if backupSourceUnavailable(detachedErr) {
				// The released-cluster list answers with 5xx/throttling right after a release.
				// An unavailable source is not evidence that the backup set is gone, so the
				// caller must keep the current state.
				return nil, fmt.Errorf("%w: %v", errBackupReadInconclusive, detachedErr)
			}
			return nil, detachedErr
		}
		if set == nil {
			if sawEntries {
				// The retained-backup enumeration is populated and does not contain this set.
				return nil, errBackupRetainedNotListed
			}
			// The cluster is released and the retained-backup list answered completely without
			// any entry. The list is populated asynchronously, so an empty answer alone does not
			// prove the set is gone: the caller weighs it against the retention contract.
			return nil, errBackupRetainedAbsent
		}
		return backupDetachedReadObject(dbClusterId, backupId, set), nil
	}
	response, err = fwSelectScopedReadObject(response, request, []fwScopedReadIdentity{
		{Source: "$.DBClusterId", Target: "DBClusterId", Value: dbClusterId},
	}, []string{"$.AdvancedDataPolicies.AdvancedDataPolicy[*].AutoCreated", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].BakType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestRegion", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DumpAction", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterKey", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterValue", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachDay", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachHour", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].PolicyId", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionValue", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcRegion", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].StorageClass", "$.BackupFrequency", "$.BackupPolicyLevel", "$.BackupRetentionPolicyOnClusterDeletion", "$.DataLevel1BackupFrequency", "$.DataLevel1BackupPeriod", "$.DataLevel1BackupRetentionPeriod", "$.DataLevel1BackupTime", "$.DataLevel2BackupAnotherRegionRegion", "$.DataLevel2BackupAnotherRegionRetentionPeriod", "$.DataLevel2BackupPeriod", "$.DataLevel2BackupRetentionPeriod", "$.PreferredBackupPeriod", "$.PreferredBackupTime", "$.PreferredNextBackupTime"})
	if err != nil {
		return nil, fmt.Errorf("selecting DescribeBackupPolicy resource: %w", err)
	}
	if response == nil {
		return nil, nil
	}

	// F1: inject the stable BackupId so the secondary DescribeBackups read filters on the exact
	// owned backup instead of merging every backup set in the cluster with "last wins".
	// BackupJobId is the CreateBackup task id and is preserved from state, not derived here:
	// DescribeBackups never reports it (verified against the live API).
	response["BackupId"] = backupId
	return backupCompleteRead(client, response, dbClusterId, backupId, knownStart, completeRead)
}

// backupDetachedReadObject builds the read object for a backup set that survived its released
// cluster. Only the detached list's backup-set fields are available; the cluster policy fields
// are absent and the caller keeps their state values.
func backupDetachedReadObject(dbClusterId, backupId string, set map[string]interface{}) map[string]interface{} {
	object := map[string]interface{}{
		backupDetachedReadKey: true,
		"DBClusterId":         dbClusterId,
		"BackupId":            backupId,
	}
	for _, field := range backupSetResponseFields {
		if value, ok := set[field.Source]; ok && value != nil {
			object[field.Target] = value
		}
	}
	return object
}

// errBackupRetainedAbsent marks an absence that cannot be proven away: the cluster's policy API
// reports the cluster as missing and the released-cluster backup list answered with a complete,
// empty enumeration. An empty page of that list is not evidence that a retained backup set does
// not exist — the list is populated asynchronously after a release and DescribeDetachedBackups
// accepts no identity filter — so the caller decides by the documented retention contract.
var errBackupRetainedAbsent = errors.New("backup set is not listed among the released cluster's backups")

// errBackupRetainedNotListed marks a backed set that the released cluster's retained-backup list
// answered without, while it did report other backup sets. Under BackupRetentionPolicyOnClusterDeletion
// LATEST only one set is kept, so a populated enumeration that lacks this id proves it was not
// retained.
var errBackupRetainedNotListed = errors.New("backup set is not among the released cluster's retained backups although the list reports other backups")

// errBackupAbsenceUnproven marks a tracked backup set that is not visible right now although the
// resource contract says it may still exist (retention keeps it, or a source could not answer).
// Read keeps the current state instead of dropping a possibly billed backup.
var errBackupAbsenceUnproven = errors.New("backup set absence cannot be proven")

// backupRetentionKeepsBackups reports whether BackupRetentionPolicyOnClusterDeletion requires
// the backup set to survive its cluster. ALL keeps every backup permanently and LATEST keeps the
// most recent one; while either applies, a missing list entry is a visibility problem, never an
// absence proof. NONE (or an unknown value written before the field was tracked) carries no such
// guarantee.
func backupRetentionKeepsBackups(retention string) bool {
	switch strings.ToUpper(strings.TrimSpace(retention)) {
	case "ALL", "LATEST":
		return true
	}
	return false
}

// backupGoneGrace bounds how long Read keeps polling after the sources answer "no match" before
// it acts on the absence: it is the poll budget for the required consecutive confirmations below,
// not an absence proof by itself.
var backupGoneGrace = 3 * time.Minute

// backupAbsenceConfirmations is the number of consecutive complete "no match" answers required
// before Read acts on an absence. A single partial or lagging enumeration must never remove a
// tracked backup set.
const backupAbsenceConfirmations = 3

// backupPollInterval is the delay between list polls. A shorter operation budget scales the
// delay down (see backupPollDelay), so a configured timeout is never outlived by a single sleep.
var backupPollInterval = 5 * time.Second

// backupPollDelay bounds the delay between polls by the operation budget: a short timeout also
// shortens the interval instead of sleeping past the deadline. Production budgets keep the
// five-second interval.
func backupPollDelay(budget time.Duration) time.Duration {
	delay := budget / 10
	if delay > backupPollInterval {
		delay = backupPollInterval
	}
	if delay < time.Millisecond {
		delay = time.Millisecond
	}
	return delay
}

// backupCreateResolveTimeout bounds how long Create waits for the backup set built by the
// CreateBackup task to appear. A configured create timeout shorter than this also shortens it.
var backupCreateResolveTimeout = 5 * time.Minute

// backupReadExisting reads a resource that Terraform already tracks. It tolerates the list-API
// lag of a live cluster and of a just-released cluster, and it only acts on an absence after
// backupAbsenceConfirmations consecutive complete empty answers. For a cluster released with
// BackupRetentionPolicyOnClusterDeletion=ALL/LATEST the backup set cannot be deleted at all, so
// a missing entry is reported as unproven (state kept) instead of being treated as gone.
func backupReadExisting(client *connectivity.AliyunClient, id, knownStart, retention string, timeout time.Duration) (map[string]interface{}, error) {
	deadline := time.Now().Add(timeout)
	interval := backupPollDelay(timeout)
	confirmations := 0
	var lastErr error
	for {
		object, err := backupGetWithStart(client, id, knownStart, true)
		switch {
		case err == nil && object != nil:
			return object, nil
		case err == nil:
			// The cluster answered but a completed enumeration does not list the backup set.
			// This is absence evidence, so it counts towards the same confirmations.
			confirmations++
		case errors.Is(err, errBackupResourceGone):
			confirmations++
		case errors.Is(err, errBackupRetainedAbsent):
			confirmations++
			lastErr = err
			if backupRetentionKeepsBackups(retention) {
				// The retention policy guarantees the set survives its cluster: an empty
				// listing is a visibility delay, never proof that the billed backup is gone.
				return nil, fmt.Errorf("%w: %s is not listed among the released cluster's backups although backup_retention_policy_on_cluster_deletion=%s keeps it; keeping the current state", errBackupAbsenceUnproven, id, retention)
			}
		case errors.Is(err, errBackupRetainedNotListed):
			confirmations++
			lastErr = err
			if strings.EqualFold(strings.TrimSpace(retention), "ALL") {
				// ALL keeps every backup set, so a populated list that does not (yet) include
				// this one may still be filling in: keep the state.
				return nil, fmt.Errorf("%w: %s is not listed among the released cluster's backups although backup_retention_policy_on_cluster_deletion=ALL keeps every set", errBackupAbsenceUnproven, id)
			}
		default:
			return nil, err
		}
		if time.Now().After(deadline) {
			if confirmations >= backupAbsenceConfirmations {
				return nil, nil
			}
			if lastErr != nil {
				return nil, fmt.Errorf("%w: %v", errBackupAbsenceUnproven, lastErr)
			}
			// Inconsistent answers inside the window: the backup set was visible at least once
			// while polling, so nothing proves it is gone.
			return nil, fmt.Errorf("%w: %s was still reported by a source during the confirmation window", errBackupAbsenceUnproven, id)
		}
		time.Sleep(interval)
	}
}

// backupGetWithRetry tolerates the brief window right after cluster creation in which
// DescribeBackupPolicy still reports the cluster as absent: the backup itself exists, so
// clearing or failing on the first empty read would be wrong.
func backupGetWithRetry(client *connectivity.AliyunClient, id string, timeout time.Duration, completeRead ...bool) (map[string]interface{}, error) {
	deadline := time.Now().Add(timeout)
	for {
		object, err := backupGet(client, id, completeRead...)
		if errors.Is(err, errBackupRetainedNotListed) {
			// The retained list is populated and does not contain the set: it does not exist.
			return nil, nil
		}
		if err != nil && !errors.Is(err, errBackupResourceGone) && !errors.Is(err, errBackupRetainedAbsent) {
			return nil, err
		}
		if object != nil {
			return object, nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		time.Sleep(backupPollInterval)
	}
}

// backupGetPolicyWithRetry reads only the cluster backup policy (no backup-set lookup) and
// tolerates the brief window in which DescribeBackupPolicy still reports a just-created
// cluster as absent.
func backupGetPolicyWithRetry(client *connectivity.AliyunClient, id string, timeout time.Duration) (map[string]interface{}, error) {
	return backupGetWithRetry(client, id, timeout, false)
}

// backupMergeResolvedBackupSet overlays the backup set resolved at create time onto the
// policy object, so the state does not depend on a second list call that may lag.
func backupMergeResolvedBackupSet(object, set map[string]interface{}, backupId string) map[string]interface{} {
	if object == nil {
		object = map[string]interface{}{}
	}
	for _, field := range backupSetResponseFields {
		if value, ok := set[field.Source]; ok && value != nil {
			object[field.Target] = value
		}
	}
	object["BackupId"] = backupId
	return object
}

// backupIdentityChunk is the width of one identity-query window. DescribeBackups silently
// returns nothing for wide StartTime/EndTime ranges (observed with 10 years and 90 days in
// ACC, while 24h works), so identity queries walk backwards in proven 7-day chunks instead.
const backupIdentityChunk = 7 * 24 * time.Hour

// backupIdentityScanChunks bounds the backward walk of a lookup that has no known start time
// (imports and first reads of an untracked set). It is a cost bound, not a retention limit: an
// exhausted scan is reported as an error rather than "not found", because a backup set older
// than the scanned period would otherwise be indistinguishable from a deleted one. 520
// seven-day windows cover ten years, beyond any PolarDB backup retention period.
const backupIdentityScanChunks = 520

// errBackupIdentityScanExhausted reports that an identity lookup without a known start time
// used its whole scan budget without a match. It is deliberately an error, never an absence.
var errBackupIdentityScanExhausted = errors.New("backup identity scan exhausted")

// backupFindRetainedBackupSet confirms whether a backup set still exists as a detached backup.
// DeleteBackup cannot remove such a set (it requires a running cluster) and reports the cluster
// as absent, so the detached list is retried briefly before the set is accepted as already gone.
func backupFindRetainedBackupSet(client *connectivity.AliyunClient, dbClusterId, backupId, knownStart string) (map[string]interface{}, error) {
	for attempt := 0; ; attempt++ {
		set, _, err := backupFindDetachedBackupSetByID(client, dbClusterId, backupId, knownStart)
		if errors.Is(err, errBackupClusterNotReleased) {
			return nil, nil
		}
		if err != nil || set != nil {
			return set, err
		}
		if attempt >= 2 {
			return nil, nil
		}
		time.Sleep(backupPollInterval)
	}
}

// backupFindBackupSetByID resolves the exact backup set from the live list first and from the
// detached list second. A known backup start time narrows the query to one window; otherwise
// the lookup walks 7-day windows backwards from now, so a backup is still found long after
// creation instead of being reported as missing by a fixed retention window.
func backupFindBackupSetByID(client *connectivity.AliyunClient, dbClusterId, backupId, knownStart string) (map[string]interface{}, error) {
	live, _, err := backupFindLiveBackupSetByID(client, dbClusterId, backupId, knownStart, false)
	if err != nil || live != nil {
		return live, err
	}
	detached, _, err := backupFindDetachedBackupSetByID(client, dbClusterId, backupId, knownStart)
	if errors.Is(err, errBackupClusterNotReleased) {
		return nil, nil
	}
	return detached, err
}

// backupFindLiveBackupSetByID resolves the exact backup set from DescribeBackups. The
// cluster-level NotFound codes mean this action cannot describe the cluster (a released
// cluster has no live backup sets), which is reported as "no match" rather than an error.
// clampEnd bounds the window end to now; callers that run while the cluster may already be
// released (where a future EndTime can be rejected) pass true, while the regular read keeps
// the small forward margin that tolerates a backup finishing within the current minute.
func backupFindLiveBackupSetByID(client *connectivity.AliyunClient, dbClusterId, backupId, knownStart string, clampEnd bool) (map[string]interface{}, bool, error) {
	const layout = "2006-01-02T15:04Z"
	nowUTC := time.Now().UTC()
	// answered reports whether DescribeBackups could describe the cluster at all; a cluster
	// NotFound means the action cannot see the cluster (released or deleted), not that this
	// particular backup set is missing.
	query := func(startTime, endTime string) (map[string]interface{}, bool, error) {
		entries, _, err := backupQueryBackupSets(client, "DescribeBackups", dbClusterId, backupId, startTime, endTime)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) {
				return nil, false, nil
			}
			return nil, false, err
		}
		if len(entries) > 1 {
			return nil, true, fmt.Errorf("DescribeBackups returned %d entries for backup set %s on cluster %s", len(entries), backupId, dbClusterId)
		}
		if len(entries) == 1 {
			return entries[0], true, nil
		}
		return nil, true, nil
	}
	if knownStart != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", knownStart); err == nil {
			end := t.Add(72 * time.Hour)
			if clampEnd && end.After(nowUTC) {
				end = nowUTC
			}
			return query(t.Add(-24*time.Hour).Format(layout), end.Format(layout))
		}
	}
	// The window end carries a forward margin: the query is formatted to the minute and a set
	// that finished within this minute (or just after the client clock) otherwise falls outside
	// EndTime and the set is reported as missing. Create already uses the same margin.
	now := nowUTC.Add(time.Hour)
	if clampEnd {
		now = nowUTC
	}
	for i := 0; i < backupIdentityScanChunks; i++ {
		end := now.Add(-time.Duration(i) * backupIdentityChunk)
		start := end.Add(-backupIdentityChunk)
		set, answered, err := query(start.Format(layout), end.Format(layout))
		if err != nil {
			return nil, answered, err
		}
		if set != nil {
			return set, answered, nil
		}
		if !answered {
			return nil, false, nil
		}
	}
	// A fixed age limit can never be reported as "not found": an existing backup set older
	// than the scan budget would otherwise be treated exactly like a deleted one.
	return nil, true, fmt.Errorf("%w: backup set %s on cluster %s was not found while scanning %s back in %s windows; its age is unknown, so absence is not proven",
		errBackupIdentityScanExhausted, backupId, dbClusterId, time.Duration(backupIdentityScanChunks)*backupIdentityChunk, backupIdentityChunk)
}

func backupCompleteRead(client *connectivity.AliyunClient, object map[string]interface{}, dbClusterId, backupId, knownStart string, includeSecondary bool) (map[string]interface{}, error) {
	object = fwNormalizeResponse(object, []fwMapping{
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].AutoCreated", Target: "$.AdvancedDataPolicies[*].AutoCreated", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].BakType", Target: "$.AdvancedDataPolicies[*].BakType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestRegion", Target: "$.AdvancedDataPolicies[*].DestRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestType", Target: "$.AdvancedDataPolicies[*].DestType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DumpAction", Target: "$.AdvancedDataPolicies[*].DumpAction", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterKey", Target: "$.AdvancedDataPolicies[*].FilterKey", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterType", Target: "$.AdvancedDataPolicies[*].FilterType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterValue", Target: "$.AdvancedDataPolicies[*].FilterValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachDay", Target: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachDay", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachHour", Target: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachHour", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].PolicyId", Target: "$.AdvancedDataPolicies[*].PolicyId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionType", Target: "$.AdvancedDataPolicies[*].RetentionType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionValue", Target: "$.AdvancedDataPolicies[*].RetentionValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcRegion", Target: "$.AdvancedDataPolicies[*].SrcRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcType", Target: "$.AdvancedDataPolicies[*].SrcType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].StorageClass", Target: "$.AdvancedDataPolicies[*].StorageClass", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupFrequency", Target: "$.BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupPolicyLevel", Target: "$.BackupPolicyLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupRetentionPolicyOnClusterDeletion", Target: "$.BackupRetentionPolicyOnClusterDeletion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupFrequency", Target: "$.DataLevel1BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupPeriod", Target: "$.DataLevel1BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupRetentionPeriod", Target: "$.DataLevel1BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupTime", Target: "$.DataLevel1BackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupAnotherRegionRegion", Target: "$.DataLevel2BackupAnotherRegionRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupAnotherRegionRetentionPeriod", Target: "$.DataLevel2BackupAnotherRegionRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupPeriod", Target: "$.DataLevel2BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupRetentionPeriod", Target: "$.DataLevel2BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PreferredBackupPeriod", Target: "$.PreferredBackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PreferredBackupTime", Target: "$.PreferredBackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PreferredNextBackupTime", Target: "$.PreferredNextBackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupEndTime", Target: "$.BackupEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupId", Target: "$.BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMethod", Target: "$.BackupMethod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMode", Target: "$.BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupSetSize", Target: "$.BackupSetSize", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupStartTime", Target: "$.BackupStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupType", Target: "$.BackupType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupsLevel", Target: "$.BackupsLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Comment", Target: "$.Comment", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ConsistentTime", Target: "$.ConsistentTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.IsAvail", Target: "$.IsAvail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Status", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupJobId", Target: "$.BackupJobId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupProgressStatus", Target: "$.BackupProgressStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.CreateTime", Target: "$.CreateTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.JobMode", Target: "$.JobMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Process", Target: "$.Process", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TaskAction", Target: "$.TaskAction", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupEndTime", Target: "$.BackupEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupId", Target: "$.BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMethod", Target: "$.BackupMethod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMode", Target: "$.BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupSetSize", Target: "$.BackupSetSize", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupStartTime", Target: "$.BackupStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupType", Target: "$.BackupType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupsLevel", Target: "$.BackupsLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ConsistentTime", Target: "$.ConsistentTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.IsAvail", Target: "$.IsAvail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Status", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.StoreStatus", Target: "$.StoreStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	object["DBClusterId"] = dbClusterId
	// F1: keep the injected stable BackupId so the secondary read below filters on the exact
	// owned backup. It is overwritten with the materialized value when the backup set is found;
	// if it is never found the gate at the end reports the resource as gone.
	if _, ok := object["BackupId"]; !ok {
		object["BackupId"] = backupId
	}
	if !includeSecondary {
		return object, nil
	}
	// The composite id names one specific backup set. Read must confirm that exact
	// backup set (exact BackupId + DBClusterId) instead of merging every entry of a list
	// response onto one object ("last wins"), which could bind a sibling backup's attributes
	// or report a foreign backup as this resource.
	//
	// The window is narrowed by the backup start time: the one recorded in state is the
	// authoritative value (DescribeBackupPolicy does not return it), so an already-tracked
	// backup is confirmed with a single exact query instead of a full scan.
	start := backupStringField(object, "BackupStartTime")
	if start == "" {
		start = knownStart
	}
	set, err := backupFindBackupSetByID(client, dbClusterId, backupId, start)
	if err != nil {
		return nil, err
	}
	if set == nil {
		// The cluster policy outlives the backup set: no exact match means the tracked
		// backup is gone, so Read clears state instead of keeping stale or foreign values.
		return nil, nil
	}
	for _, field := range backupSetResponseFields {
		if value, ok := set[field.Source]; ok && value != nil {
			object[field.Target] = value
		}
	}
	// Task progress fields belong to the CreateBackup task recorded in state and are not part
	// of the backup-set response. They are enriched by the caller, which holds the task id
	// from state: DescribeBackups never reports a BackupJobId for the set.
	return object, nil
}

// backupSetResponseFields maps a single backup-set entry onto the resource model. Only
// fields documented on DescribeBackups/DescribeDetachedBackups are listed.
var backupSetResponseFields = []struct{ Target, Source string }{
	{"BackupEndTime", "BackupEndTime"},
	{"BackupId", "BackupId"},
	{"BackupMethod", "BackupMethod"},
	{"BackupMode", "BackupMode"},
	{"BackupSetSize", "BackupSetSize"},
	{"BackupStartTime", "BackupStartTime"},
	{"BackupType", "BackupType"},
	{"BackupsLevel", "BackupsLevel"},
	{"Comment", "Comment"},
	{"ConsistentTime", "ConsistentTime"},
	{"DBClusterId", "DBClusterId"},
	{"IsAvail", "IsAvail"},
	{"Status", "BackupStatus"},
	{"StoreStatus", "StoreStatus"},
}

// backupTaskResponseFields maps a single DescribeBackupTasks entry onto the resource model.
var backupTaskResponseFields = []struct{ Target, Source string }{
	{"BackupJobId", "BackupJobId"},
	{"BackupProgressStatus", "BackupProgressStatus"},
	{"CreateTime", "StartTime"},
	{"JobMode", "JobMode"},
	{"Process", "Process"},
	{"TaskAction", "TaskAction"},
}

// backupStringField returns a response field as a string, treating absence and nil as empty.
func backupStringField(object map[string]interface{}, key string) string {
	value, ok := object[key]
	if !ok || value == nil {
		return ""
	}
	text := fmt.Sprint(value)
	if text == "<nil>" {
		return ""
	}
	return text
}

// backupTaskList extracts the BackupJob task entries from a DescribeBackupTasks response.
func backupTaskList(resp map[string]interface{}) []map[string]interface{} {
	out := []map[string]interface{}{}
	items, _ := resp["Items"].(map[string]interface{})
	jobs, _ := items["BackupJob"].([]interface{})
	for _, j := range jobs {
		if m, ok := j.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// backupBackupsList extracts the backup-set entries from a DescribeBackups response.
func backupBackupsList(resp map[string]interface{}) []map[string]interface{} {
	out := []map[string]interface{}{}
	items, _ := resp["Items"].(map[string]interface{})
	backs, _ := items["Backup"].([]interface{})
	for _, b := range backs {
		if m, ok := b.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// backupPageSize is the largest page size DescribeBackups/DescribeDetachedBackups accept.
const backupPageSize = 100

// backupMaxPages bounds paging so an unfiltered response cannot loop forever.
const backupMaxPages = 50

// backupEntryMatches reports whether a list entry is the requested backup set: the exact
// BackupId, confirmed against DBClusterId whenever the API reports it (a missing
// DBClusterId in the response is not counter-evidence).
func backupEntryMatches(entry map[string]interface{}, dbClusterId, backupId string) bool {
	if backupId == "" || backupStringField(entry, "BackupId") != backupId {
		return false
	}
	if entryClusterId := backupStringField(entry, "DBClusterId"); entryClusterId != "" && entryClusterId != dbClusterId {
		return false
	}
	return true
}

// backupQueryBackupSets pages through DescribeBackups or DescribeDetachedBackups for one
// cluster and returns the entries matching the requested backup set exactly. An empty
// backupId lists every entry in the range.
//
// The server-side BackupId filter is deliberately not sent: once a cluster is switched to
// BackupPolicyLevel=Advanced, DescribeBackups answers a BackupId-filtered request with an
// empty list while an unfiltered request for the same window and cluster returns that very
// backup set (observed: the set was visible unfiltered while the filtered lookup answered
// empty for 30s of retries). Selection stays exact because every entry is matched in Go by BackupId and,
// when the API reports it, DBClusterId.
func backupQueryBackupSets(client *connectivity.AliyunClient, action, dbClusterId, backupId, startTime, endTime string) ([]map[string]interface{}, bool, error) {
	matched := []map[string]interface{}{}
	sawEntries := false
	for page := 1; page <= backupMaxPages; page++ {
		request := map[string]interface{}{
			"DBClusterId": dbClusterId,
			"StartTime":   startTime,
			"EndTime":     endTime,
			"PageNumber":  page,
			"PageSize":    backupPageSize,
			"RegionId":    client.RegionId,
		}
		response, err := retryRpcPost(client, "polardb", "2017-08-01", action, request, 1*time.Minute)
		if err != nil {
			return nil, sawEntries, err
		}
		entries := backupBackupsList(response)
		if len(entries) > 0 {
			sawEntries = true
		}
		for _, entry := range entries {
			if backupId == "" || backupEntryMatches(entry, dbClusterId, backupId) {
				matched = append(matched, entry)
			}
		}
		if len(entries) < backupPageSize {
			return matched, sawEntries, nil
		}
	}
	if len(matched) > 0 {
		return matched, sawEntries, nil
	}
	// An incomplete scan with no match must never be mistaken for "the backup set does not
	// exist" (identical semantics to a NotFound), so surface it as an error instead.
	return nil, sawEntries, fmt.Errorf("%s did not finish paging within %d pages for backup set %q on cluster %s", action, backupMaxPages, backupId, dbClusterId)
}

// backupSourceUnavailable reports whether an error means a list source could not answer
// (throttling, 5xx, timeout) rather than answering "no match". retryRpcPost wraps the SDK error
// in errs.ComplexError, which does not unwrap, so the transport markers are matched on the
// rendered error as well.
func backupSourceUnavailable(err error) bool {
	if err == nil {
		return false
	}
	if helper.NeedRetry(err) {
		return true
	}
	text := err.Error()
	for _, marker := range []string{"code: 500", "code: 502", "code: 503", "code: 504", "Throttling", "Client.Timeout", "ServiceUnavailable"} {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return false
}

// errBackupClusterNotReleased reports that DescribeDetachedBackups refused the query because
// the cluster still exists. A live cluster has no detached backup sets, so callers treat it as
// "no match" without walking the remaining time windows.
var errBackupClusterNotReleased = errors.New("cluster is not released and has no detached backup sets")

// backupFindDetachedBackupSet resolves a backup set from DescribeDetachedBackups, the API that
// reports the backups retained for a released PolarDB cluster.
func backupFindDetachedBackupSet(client *connectivity.AliyunClient, dbClusterId, backupId, startTime, endTime string) (map[string]interface{}, bool, error) {
	detached, sawEntries, err := backupQueryBackupSets(client, "DescribeDetachedBackups", dbClusterId, backupId, startTime, endTime)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) {
			return nil, false, errBackupClusterNotReleased
		}
		return nil, sawEntries, err
	}
	if len(detached) > 1 {
		return nil, sawEntries, fmt.Errorf("DescribeDetachedBackups returned %d entries for backup set %s on cluster %s", len(detached), backupId, dbClusterId)
	}
	if len(detached) == 1 {
		return detached[0], sawEntries, nil
	}
	return nil, sawEntries, nil
}

// backupFindDetachedBackupSetByID resolves a retained backup set without a live cluster. A
// known backup start time narrows the query to one window; otherwise the lookup walks the
// same backward chunks as backupFindLiveBackupSetByID. errBackupClusterNotReleased is returned
// unchanged so callers can tell "this source cannot describe the cluster" from "no match".
//
// Unlike DescribeBackups, DescribeDetachedBackups answers a window that extends into the future
// with a server-side 500, so the window end is clamped to now: a released cluster's backup sets
// always started in the past.
//
// The second result reports whether the list answered with any backup set at all. A populated
// list that does not contain this id is stronger evidence than an empty one: the cluster's
// retained backups are enumerable and this set is not among them.
func backupFindDetachedBackupSetByID(client *connectivity.AliyunClient, dbClusterId, backupId, knownStart string) (map[string]interface{}, bool, error) {
	const layout = "2006-01-02T15:04Z"
	now := time.Now().UTC()
	if knownStart != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", knownStart); err == nil {
			end := t.Add(72 * time.Hour)
			if end.After(now) {
				end = now
			}
			return backupFindDetachedBackupSet(client, dbClusterId, backupId, t.Add(-24*time.Hour).Format(layout), end.Format(layout))
		}
	}
	sawEntries := false
	for i := 0; i < backupIdentityScanChunks; i++ {
		end := now.Add(-time.Duration(i) * backupIdentityChunk)
		start := end.Add(-backupIdentityChunk)
		set, saw, err := backupFindDetachedBackupSet(client, dbClusterId, backupId, start.Format(layout), end.Format(layout))
		sawEntries = sawEntries || saw
		if errors.Is(err, errBackupClusterNotReleased) {
			// The cluster exists, so it has no detached backup sets at all.
			return nil, sawEntries, err
		}
		if err != nil {
			return nil, sawEntries, err
		}
		if set != nil {
			return set, sawEntries, nil
		}
	}
	return nil, sawEntries, fmt.Errorf("%w: backup set %s on released cluster %s was not found while scanning %s back in %s windows; absence is not proven",
		errBackupIdentityScanExhausted, backupId, dbClusterId, time.Duration(backupIdentityScanChunks)*backupIdentityChunk, backupIdentityChunk)
}

// backupFindBackupTask returns the DescribeBackupTasks entry for the exact task id.
func backupFindBackupTask(client *connectivity.AliyunClient, dbClusterId, backupJobId string) (map[string]interface{}, error) {
	if backupJobId == "" {
		return nil, nil
	}
	request := map[string]interface{}{
		"DBClusterId": dbClusterId,
		"BackupJobId": backupJobId,
		"RegionId":    client.RegionId,
	}
	response, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeBackupTasks", request, 1*time.Minute)
	if err != nil {
		return nil, err
	}
	for _, task := range backupTaskList(response) {
		if fmt.Sprint(task["BackupJobId"]) == backupJobId {
			return task, nil
		}
	}
	return nil, nil
}

// backupWaitForBackupTask waits for the CreateBackup task to reach Finished. Observed
// responses return a task id in the documented "backup set id" field, and the backup set
// only materializes after the task completes, so the task is the authoritative completion
// signal before any set lookup.
func backupWaitForBackupTask(client *connectivity.AliyunClient, dbClusterId, backupJobId string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		task, err := backupFindBackupTask(client, dbClusterId, backupJobId)
		if err != nil {
			return err
		}
		if task != nil && backupStringField(task, "BackupProgressStatus") == "Finished" {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("backup task %s on cluster %s did not reach Finished within %s", backupJobId, dbClusterId, timeout)
		}
		time.Sleep(5 * time.Second)
	}
}

// backupResolveCreatedBackupSet binds the backup set produced by the CreateBackup task named by
// backupJobId and returns its BackupId.
//
// CreateBackup answers with a task id (the documented BackupJobId field) that DescribeBackups
// never reports: the backup set it produces is listed only under its own BackupId. The task
// start time is the single server-provided value shared by both records, so a set is attributed
// to this task only when it is manual AND its BackupStartTime equals the task StartTime exactly.
// No match before the deadline, more than one match, or an unusable task record fails the
// create loudly instead of binding a backup that another actor created on the same cluster.
func backupResolveCreatedBackupSet(client *connectivity.AliyunClient, dbClusterId, backupJobId string, timeout time.Duration) (string, map[string]interface{}, error) {
	const secondsLayout = "2006-01-02T15:04:05Z"
	const minuteLayout = "2006-01-02T15:04Z"
	deadline := time.Now().Add(timeout)
	interval := backupPollDelay(timeout)
	for {
		task, err := backupFindBackupTask(client, dbClusterId, backupJobId)
		if err != nil {
			return "", nil, err
		}
		if task == nil {
			if time.Now().After(deadline) {
				return "", nil, fmt.Errorf("backup task %s on cluster %s is no longer reported by DescribeBackupTasks; refusing to bind any backup set", backupJobId, dbClusterId)
			}
			time.Sleep(interval)
			continue
		}
		startedAt := backupStringField(task, "StartTime")
		taskStart, parseErr := time.Parse(secondsLayout, startedAt)
		if parseErr != nil {
			return "", nil, fmt.Errorf("backup task %s on cluster %s reports an unusable start time %q: %v", backupJobId, dbClusterId, startedAt, parseErr)
		}
		windowStart := taskStart.Add(-1 * time.Minute)
		windowEnd := taskStart.Add(30 * time.Minute)
		if now := time.Now().UTC(); windowEnd.After(now) {
			windowEnd = now
		}
		// The server filters by minute; a zero-width window is rejected, so keep at least
		// one full minute of range around the task start.
		if minimum := windowStart.Add(time.Minute); windowEnd.Before(minimum) {
			windowEnd = minimum
		}
		entries, _, err := backupQueryBackupSets(client, "DescribeBackups", dbClusterId, "", windowStart.Format(minuteLayout), windowEnd.Format(minuteLayout))
		if err != nil {
			return "", nil, err
		}
		matches := make([]map[string]interface{}, 0, 1)
		for _, entry := range entries {
			if backupStringField(entry, "BackupMode") != "Manual" {
				continue
			}
			if backupStringField(entry, "BackupStartTime") == startedAt {
				matches = append(matches, entry)
			}
		}
		if len(matches) > 1 {
			return "", nil, fmt.Errorf("backup task %s on cluster %s started at %s, but DescribeBackups reports %d manual backup sets starting at that second; refusing to bind one of them", backupJobId, dbClusterId, startedAt, len(matches))
		}
		if len(matches) == 1 {
			backupId := backupStringField(matches[0], "BackupId")
			if backupId == "" {
				return "", nil, fmt.Errorf("the backup set produced by task %s on cluster %s has no BackupId", backupJobId, dbClusterId)
			}
			return backupId, matches[0], nil
		}
		if time.Now().After(deadline) {
			return "", nil, fmt.Errorf("backup task %s on cluster %s started at %s, but no manual backup set starting at that second appeared in DescribeBackups within %s; refusing to bind any other backup set", backupJobId, dbClusterId, startedAt, timeout)
		}
		time.Sleep(interval)
	}
}

// A required exact request can identify a singleton configuration whose API omits its ID.
func fwSelectScopedReadObject(raw interface{}, request map[string]interface{}, identities []fwScopedReadIdentity, payloadPaths []string) (map[string]interface{}, error) {
	if len(identities) == 0 {
		return nil, fmt.Errorf("no exact request identity was declared")
	}
	for _, identity := range identities {
		value, exists := request[identity.Target].(string)
		if !exists || identity.Value == "" || value != identity.Value {
			return nil, fmt.Errorf("singleton request does not match the resource identity")
		}
	}
	object, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("singleton response is not an object")
	}
	for _, identity := range identities {
		values := fwWalkMappedPath(object, fwMappingPath(identity.Source), nil)
		if len(values) > 1 {
			return nil, fmt.Errorf("singleton response has ambiguous identity")
		}
		if len(values) == 1 && (values[0].Value == nil || fmt.Sprint(values[0].Value) != identity.Value) {
			return nil, fmt.Errorf("singleton response identity does not match the request")
		}
	}
	for _, path := range payloadPaths {
		values := fwWalkMappedPath(object, fwMappingPath(path), nil)
		for _, value := range values {
			if value.Value != nil {
				return object, nil
			}
		}
	}
	return nil, fmt.Errorf("singleton response has no declared configuration payload")
}

// Optional update parameters only participate when their mapped value changed.
func fwChangedUpdateMappings(bindings []fwMapping, desired, prior map[string]interface{}, requiredSources []string) []fwMapping {
	selected := make([]fwMapping, 0, len(bindings))
	for _, binding := range bindings {
		required := false
		for _, source := range requiredSources {
			if binding.Source == source || strings.HasPrefix(binding.Source, source+".") || strings.HasPrefix(binding.Source, source+"[*]") {
				required = true
				break
			}
		}
		if required || !fwValuesEqual(fwWalkMappedPath(desired, fwMappingPath(binding.Source), nil), fwWalkMappedPath(prior, fwMappingPath(binding.Source), nil)) {
			selected = append(selected, binding)
		}
	}
	return selected
}

// fwStr navigates a dotted key in m and returns a types.String (missing -> StringNull).
// Used by nested-list ObjectValues helpers to read scalar sub-fields from each list item;
// the top-level State closures in resource.go capture the Get response, so they cannot read
// per-item sub-maps.
func fwStr(m map[string]interface{}, key string) types.String {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.StringNull()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.StringNull()
		}
	}
	if v, ok := cur.(string); ok {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func fwBool(m map[string]interface{}, key string) types.Bool {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.BoolNull()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.BoolNull()
		}
	}
	if v, ok := cur.(bool); ok {
		return types.BoolValue(v)
	}
	return types.BoolNull()
}

// backupAdvancedDataPoliciesToTypeList wraps the API object array into a framework
// types.List of nested objects (types.Object). nil/empty/non-array -> empty list
// (.#=0) to match v1 d.Set(..., []), so TestCheckResourceAttrSet("x.#") passes and
// UseStateForUnknown unknown->0 stays consistent (ListNull leaves .# unset, failing
// the .# check).
func backupAdvancedDataPoliciesToTypeList(raw interface{}) types.List {
	objType := backupAdvancedDataPoliciesObjType()
	elems := make([]attr.Value, 0)
	// Handle both array (nested list) and single object (MaxItems:1 object type).
	// When the API returns AdvancedDataPolicies:{AdvancedDataPolicy:[]}, the leaf mappings
	// produce zero entries and raw arrives as the wrapper map. Unwrap the inner
	// "AdvancedDataPolicy" container so an empty array stays an empty list instead of being
	// wrapped into a single all-null block (state item count 0 -> 1).
	if raw != nil {
		var arr []interface{}
		switch v := raw.(type) {
		case []interface{}:
			arr = v
		case map[string]interface{}:
			if inner, ok := v["AdvancedDataPolicy"].([]interface{}); ok {
				arr = inner
			} else {
				arr = []interface{}{v}
			}
		}
		for _, item := range arr {
			m, _ := item.(map[string]interface{})
			// DescribeBackupPolicy documents AutoCreated as "true: system-generated
			// policy, false: user-defined policy", and enabling advanced backup makes the service
			// create one. A system-generated entry is not part of the configuration, so it must
			// not enter state: it would change the block count between apply and import and hide
			// the readable configuration fields behind an import ignore.
			if strings.EqualFold(fmt.Sprint(m["AutoCreated"]), "true") {
				continue
			}
			obj, _ := types.ObjectValue(backupAdvancedDataPoliciesAttrTypes(), backupAdvancedDataPoliciesObjectValues(m))
			elems = append(elems, obj)
		}
	}
	l, _ := types.ListValue(objType, elems)
	return l
}

// backupAdvancedDataPoliciesAttrTypes returns the framework attr types for each nested field.
func backupAdvancedDataPoliciesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action_type":                 types.StringType,
		"only_preserve_one_each_hour": types.BoolType,
		"storage_class":               types.StringType,
		"retention_value":             types.StringType,
		"bak_type":                    types.StringType,
		"filter_key":                  types.StringType,
		"filter_value":                types.StringType,
		"dest_type":                   types.StringType,
		"filter_type":                 types.StringType,
		"src_region":                  types.StringType,
		"auto_created":                types.BoolType,
		"only_preserve_one_each_day":  types.BoolType,
		"src_type":                    types.StringType,
		"dump_action":                 types.StringType,
		"dest_region":                 types.StringType,
		"policy_id":                   types.StringType,
		"retention_type":              types.StringType,
	}
}

// backupAdvancedDataPoliciesObjType is the element ObjectType for the nested list.
func backupAdvancedDataPoliciesObjType() types.ObjectType {
	return types.ObjectType{AttrTypes: backupAdvancedDataPoliciesAttrTypes()}
}

// backupAdvancedDataPoliciesObjectValues maps a single API object (map) to framework attr
// values. Nested sub-lists recurse via their own ToTypeList helpers; scalars use fwStr/fwBool.
func backupAdvancedDataPoliciesObjectValues(m map[string]interface{}) map[string]attr.Value {
	out := map[string]attr.Value{}
	out["action_type"] = types.StringNull()
	out["only_preserve_one_each_hour"] = fwBool(m, "OnlyPreserveOneEachHour")
	out["storage_class"] = fwStr(m, "StorageClass")
	out["retention_value"] = fwStr(m, "RetentionValue")
	out["bak_type"] = fwStr(m, "BakType")
	out["filter_key"] = fwStr(m, "FilterKey")
	out["filter_value"] = fwStr(m, "FilterValue")
	out["dest_type"] = fwStr(m, "DestType")
	out["filter_type"] = fwStr(m, "FilterType")
	out["src_region"] = fwStr(m, "SrcRegion")
	out["auto_created"] = fwBool(m, "AutoCreated")
	out["only_preserve_one_each_day"] = fwBool(m, "OnlyPreserveOneEachDay")
	out["src_type"] = fwStr(m, "SrcType")
	out["dump_action"] = fwStr(m, "DumpAction")
	out["dest_region"] = fwStr(m, "DestRegion")
	out["policy_id"] = fwStr(m, "PolicyId")
	out["retention_type"] = fwStr(m, "RetentionType")
	return out
}

// fwWireField preserves API/source names independently of Terraform schema names.

// backupAdvancedPoliciesChanged compares the configured advanced policies with the state.
// An absent block and an empty list both mean "no policies configured": the framework plans
// an omitted block as null while the read side yields an empty list, so a plain comparison
// would always report a change and send a needless ModifyBackupPolicy for backup-only
// configurations (the cluster policy is shared with the sibling policy resource).
func backupAdvancedPoliciesChanged(plan, state types.List) bool {
	if plan.IsUnknown() {
		return false
	}
	planEmpty := plan.IsNull() || len(plan.Elements()) == 0
	stateEmpty := state.IsUnknown() || state.IsNull() || len(state.Elements()) == 0
	if planEmpty && stateEmpty {
		return false
	}
	if state.IsUnknown() {
		return false
	}
	return fwPlannedChanged(plan, state)
}

// backupAdvancedPolicyIdentity returns the stable identity of one advanced policy entry: its
// PolicyId when the read-back reported one, otherwise the structural chain that identifies the
// policy for an ActionType=CREATE/UPDATE/DELETE request.
func backupAdvancedPolicyIdentity(policy types.Object) string {
	if policy.IsNull() || policy.IsUnknown() {
		return ""
	}
	attributes := policy.Attributes()
	if value, ok := attributes["policy_id"].(types.String); ok && !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
		return "id:" + value.ValueString()
	}
	parts := make([]string, 0, len(backupAdvancedPolicyChainKeys))
	for _, key := range backupAdvancedPolicyChainKeys {
		value, ok := attributes[key].(types.String)
		if !ok || value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
			continue
		}
		parts = append(parts, key+"="+value.ValueString())
	}
	if len(parts) == 0 {
		return ""
	}
	return "chain:" + strings.Join(parts, "|")
}

// backupAdvancedPoliciesRemoved reports that entries present in state have no counterpart in
// the plan. Dropping an entry is not a protocol-level delete here: ModifyBackupPolicy removes a
// policy only through an explicit ActionType=DELETE on that entry, so an implicit removal would
// leave the cloud policy unmanaged. The comparison is by stable identity, never by count: a
// same-sized replacement (A removed, B added) is also a removal of A.
func backupAdvancedPoliciesRemoved(plan, state types.List) bool {
	if state.IsNull() || state.IsUnknown() || len(state.Elements()) == 0 {
		return false
	}
	if plan.IsUnknown() {
		return false
	}
	planned := map[string]struct{}{}
	if !plan.IsNull() {
		for _, element := range plan.Elements() {
			policy, ok := element.(types.Object)
			if !ok {
				continue
			}
			if identity := backupAdvancedPolicyIdentity(policy); identity != "" {
				planned[identity] = struct{}{}
			}
		}
	}
	for _, element := range state.Elements() {
		policy, ok := element.(types.Object)
		if !ok {
			continue
		}
		if action, ok := policy.Attributes()["action_type"].(types.String); ok && !action.IsNull() && !action.IsUnknown() && strings.EqualFold(action.ValueString(), "DELETE") {
			// The entry already carries its own deletion operation: once applied, dropping it
			// from the configuration is the intended retirement, not an implicit removal.
			continue
		}
		identity := backupAdvancedPolicyIdentity(policy)
		if identity == "" {
			continue
		}
		if _, kept := planned[identity]; !kept {
			return true
		}
	}
	return false
}

// backupWeekdaySetChanged reports a real weekday add/remove. An unknown or null planned value
// (an omitted Optional+Computed field) is not a change: comparing it as an empty string would
// always differ and trigger a ModifyBackupPolicy that resets the shared cluster policy.
func backupWeekdaySetChanged(plan, state preferredBackupPeriodStringValue) bool {
	if plan.IsUnknown() || plan.IsNull() || state.IsUnknown() || state.IsNull() {
		return false
	}
	return !backupWeekdaySetEqual(plan.ValueString(), state.ValueString())
}

// fwPlannedChanged ignores computed unknown leaves while comparing configured nested values.
func fwPlannedChanged(plan, state attr.Value) bool {
	if plan == nil || plan.IsUnknown() {
		return false
	}
	if state == nil || state.IsNull() || state.IsUnknown() || plan.IsNull() {
		return state == nil || !plan.Equal(state)
	}
	switch desired := plan.(type) {
	case types.Object:
		previous, ok := state.(types.Object)
		if !ok {
			return true
		}
		oldItems := previous.Attributes()
		for key, value := range desired.Attributes() {
			if fwPlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	case types.List:
		previous, ok := state.(types.List)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for index, value := range items {
			if fwPlannedChanged(value, oldItems[index]) {
				return true
			}
		}
		return false
	case types.Set:
		previous, ok := state.(types.Set)
		if !ok {
			return true
		}
		return !fwSetMatches(desired.Elements(), previous.Elements(), func(a, b attr.Value) bool { return !fwPlannedChanged(a, b) })
	case types.Map:
		previous, ok := state.(types.Map)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for key, value := range items {
			if fwPlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	}
	return !plan.Equal(state)
}

// fwSetMatches pairs unordered members without depending on response or plan order.
func fwSetMatches(left, right []attr.Value, compatible func(attr.Value, attr.Value) bool) bool {
	if len(left) != len(right) {
		return false
	}
	owners := make([]int, len(right))
	for index := range owners {
		owners[index] = -1
	}
	var assign func(int, []bool) bool
	assign = func(index int, seen []bool) bool {
		for candidate := range right {
			if seen[candidate] || !compatible(left[index], right[candidate]) {
				continue
			}
			seen[candidate] = true
			if owners[candidate] == -1 || assign(owners[candidate], seen) {
				owners[candidate] = index
				return true
			}
		}
		return false
	}
	for index := range left {
		if !assign(index, make([]bool, len(right))) {
			return false
		}
	}
	return true
}

func fwKnownValue(value attr.Value) attr.Value {
	switch v := value.(type) {
	case preferredBackupPeriodStringValue:
		if v.IsUnknown() {
			return preferredBackupPeriodStringValue{types.StringNull()}
		}
	case types.String:
		if v.IsUnknown() {
			return types.StringNull()
		}
	case types.Bool:
		if v.IsUnknown() {
			return types.BoolNull()
		}
	case types.Int64:
		if v.IsUnknown() {
			return types.Int64Null()
		}
	case types.Float64:
		if v.IsUnknown() {
			return types.Float64Null()
		}
	case types.List:
		if v.IsUnknown() {
			return types.ListNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, fwKnownValue(item))
			}
			return types.ListValueMust(v.ElementType(context.Background()), items)
		}
	case types.Set:
		if v.IsUnknown() {
			return types.SetNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, fwKnownValue(item))
			}
			return types.SetValueMust(v.ElementType(context.Background()), items)
		}
	case types.Map:
		if v.IsUnknown() {
			return types.MapNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Elements() {
				items[key] = fwKnownValue(item)
			}
			return types.MapValueMust(v.ElementType(context.Background()), items)
		}
	case types.Object:
		if v.IsUnknown() {
			return types.ObjectNull(v.AttributeTypes(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Attributes() {
				items[key] = fwKnownValue(item)
			}
			return types.ObjectValueMust(v.AttributeTypes(context.Background()), items)
		}
	}
	return value
}

func fwValueToInterface(value attr.Value, fields map[string]fwWireField, single bool) interface{} {
	if value == nil || value.IsNull() || value.IsUnknown() {
		return nil
	}
	switch v := value.(type) {
	case interface{ ValueString() string }:
		return v.ValueString()
	case types.Bool:
		return v.ValueBool()
	case types.Int64:
		return v.ValueInt64()
	case types.Float64:
		return v.ValueFloat64()
	case types.Object:
		out := map[string]interface{}{}
		for key, child := range v.Attributes() {
			field, mapped := fields[key]
			if !mapped {
				field.Name = key
			}
			if raw := fwValueToInterface(child, field.Fields, field.Single); raw != nil {
				out[field.Name] = raw
			}
		}
		return out
	case types.Map:
		out := map[string]interface{}{}
		for key, child := range v.Elements() {
			if raw := fwValueToInterface(child, nil, false); raw != nil {
				out[key] = raw
			}
		}
		return out
	case types.List:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := fwValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		if single {
			if len(out) == 0 {
				return nil
			}
			return out[0]
		}
		return out
	case types.Set:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := fwValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		return out
	}
	return nil
}

func fwMappingPath(path string) []string {
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	path = strings.ReplaceAll(path, "[*]", ".*")
	if path == "" {
		return nil
	}
	return strings.Split(strings.TrimPrefix(path, "."), ".")
}

func fwWalkMappedPath(value interface{}, path []string, indexes []int) []fwMappedValue {
	if len(path) == 0 {
		return []fwMappedValue{{Value: value, Indexes: indexes}}
	}
	if path[0] == "*" {
		var out []fwMappedValue
		if items, ok := value.([]interface{}); ok {
			for i, item := range items {
				next := append(append([]int{}, indexes...), i)
				out = append(out, fwWalkMappedPath(item, path[1:], next)...)
			}
		}
		return out
	}
	if m, ok := value.(map[string]interface{}); ok {
		if child, exists := m[path[0]]; exists {
			return fwWalkMappedPath(child, path[1:], indexes)
		}
	}
	return nil
}

func fwSetMappedPath(value interface{}, path []string, indexes []int, raw interface{}) interface{} {
	if len(path) == 0 {
		return raw
	}
	if path[0] == "*" {
		if len(indexes) == 0 {
			return value
		}
		items, _ := value.([]interface{})
		for len(items) <= indexes[0] {
			items = append(items, nil)
		}
		items[indexes[0]] = fwSetMappedPath(items[indexes[0]], path[1:], indexes[1:], raw)
		return items
	}
	m, _ := value.(map[string]interface{})
	if m == nil {
		m = map[string]interface{}{}
	}
	m[path[0]] = fwSetMappedPath(m[path[0]], path[1:], indexes, raw)
	return m
}

func fwApplyMappings(request, values map[string]interface{}, bindings []fwMapping, roa bool) error {
	mapped := map[string]interface{}{}
	repeated := map[string]bool{}
	for _, binding := range bindings {
		source, target := binding.Source, binding.Target
		if binding.Encoding == "comma" || binding.Encoding == "first" || binding.Encoding == "wrap-json" || binding.Encoding == "wrap-array" {
			source = strings.TrimSuffix(source, "[*]")
			target = strings.TrimSuffix(target, "[*]")
		}
		entries := fwWalkMappedPath(values, fwMappingPath(source), nil)
		if binding.Encoding == "scoped-const" {
			parent := source[:strings.LastIndex(source, ".")]
			entries = fwWalkMappedPath(values, fwMappingPath(parent), nil)
		}
		if strings.HasPrefix(binding.Encoding, "concat|") {
			parts := strings.Split(binding.Encoding, "|")
			fields := []string{}
			complete := true
			for _, propertyPath := range parts[2:] {
				found := fwWalkMappedPath(values, fwMappingPath(propertyPath), nil)
				if len(found) == 0 || found[0].Value == nil {
					complete = false
					break
				}
				fields = append(fields, fmt.Sprint(found[0].Value))
			}
			entries = nil
			if complete {
				entries = []fwMappedValue{{Value: strings.Join(fields, parts[1])}}
			}
		}
		if strings.Contains(source, "[*]") && !strings.Contains(target, "[*]") && binding.Encoding == "" {
			items := []interface{}{}
			for _, entry := range entries {
				if entry.Value != nil {
					items = append(items, entry.Value)
				}
			}
			if len(entries) > 0 {
				entries = []fwMappedValue{{Value: items}}
			}
		}
		concreteEntries := entries[:0]
		for _, entry := range entries {
			if entry.Value != nil {
				concreteEntries = append(concreteEntries, entry)
			}
		}
		entries = concreteEntries
		if binding.Constant != "" && !(binding.Encoding == "default" && len(entries) > 0) {
			if (binding.Encoding == "key" || binding.Encoding == "scoped-const") && len(entries) == 0 {
				continue
			}
			var raw interface{}
			if err := json.Unmarshal([]byte(binding.Constant), &raw); err != nil {
				return err
			}
			if len(entries) == 0 {
				entries = []fwMappedValue{{Value: raw}}
			} else {
				for index := range entries {
					entries[index].Value = raw
				}
			}
		}
		if strings.HasPrefix(binding.Encoding, "uniform-scalar|") && len(entries) > 0 {
			scalarType := strings.TrimPrefix(binding.Encoding, "uniform-scalar|")
			for _, entry := range entries {
				valid := false
				switch scalarType {
				case "boolean":
					_, valid = entry.Value.(bool)
				case "integer":
					switch value := entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
						valid = true
					case float64:
						valid = value == float64(int64(value))
					case json.Number:
						_, err := value.Int64()
						valid = err == nil
					}
				case "number":
					switch entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, json.Number:
						valid = true
					}
				}
				if !valid {
					return fmt.Errorf("configured value does not match scalar API parameter %s type %s", target, scalarType)
				}
			}
			first, err := json.Marshal(entries[0].Value)
			if err != nil {
				return err
			}
			for _, entry := range entries[1:] {
				next, err := json.Marshal(entry.Value)
				if err != nil {
					return err
				}
				if string(next) != string(first) {
					return fmt.Errorf("conflicting configured values for scalar API parameter %s", target)
				}
			}
			entries = entries[:1]
		}
		for _, entry := range entries {
			raw := entry.Value
			if raw == nil {
				continue
			}
			if mappedValue, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mappedValue
			}
			if strings.HasPrefix(binding.Encoding, "suffix|") {
				raw = fmt.Sprint(raw) + strings.TrimPrefix(binding.Encoding, "suffix|")
			}
			switch binding.Encoding {
			case "first":
				if list, ok := raw.([]interface{}); ok {
					if len(list) == 0 {
						continue
					}
					raw = list[0]
				}
			case "comma":
				if list, ok := raw.([]interface{}); ok {
					parts := make([]string, len(list))
					for i, item := range list {
						parts[i] = fmt.Sprint(item)
					}
					raw = strings.Join(parts, ",")
				}
			case "wrap-array":
				raw = []interface{}{raw}
			case "wrap-json":
				encoded, err := json.Marshal([]interface{}{raw})
				if err != nil {
					return err
				}
				raw = string(encoded)
			}
			indexes := append([]int(nil), entry.Indexes...)
			if len(binding.IndexSources) > 0 {
				if len(binding.IndexSources) != strings.Count(target, "[*]") {
					return fmt.Errorf("request wildcard alignment has an unexpected target dimension count")
				}
				indexes = make([]int, len(binding.IndexSources))
				for position, sourceIndex := range binding.IndexSources {
					if sourceIndex == -1 {
						continue
					}
					if sourceIndex < 0 || sourceIndex >= len(entry.Indexes) {
						return fmt.Errorf("request wildcard alignment is missing a source dimension")
					}
					indexes[position] = entry.Indexes[sourceIndex]
				}
			} else {
				for len(indexes) < strings.Count(target, "[*]") {
					indexes = append(indexes, 0)
				}
			}
			mapped = fwSetMappedPath(mapped, fwMappingPath(target), indexes, raw).(map[string]interface{})
		}
		if binding.Repeat {
			if parts := fwMappingPath(target); len(parts) > 0 {
				repeated[parts[0]] = true
			}
		}
	}
	if err := fwCompactKvPairs(mapped, bindings); err != nil {
		return err
	}
	for key, raw := range mapped {
		if roa {
			request[key] = raw
			continue
		}
		switch raw.(type) {
		case map[string]interface{}, []interface{}:
			if repeated[key] {
				fwFlattenRequest(request, key, raw)
			} else {
				encoded, err := json.Marshal(raw)
				if err != nil {
					return err
				}
				request[key] = string(encoded)
			}
		default:
			request[key] = raw
		}
	}
	return nil
}

// Only explicitly declared kvPairs groups are compacted; ordinary indexed lists retain shape.
func fwCompactKvPairs(mapped map[string]interface{}, bindings []fwMapping) error {
	groups := map[string][2]string{}
	for _, binding := range bindings {
		if binding.Group == "" {
			continue
		}
		fields := [2]string{binding.GroupKey, binding.GroupValue}
		if fields[0] == "" || fields[1] == "" || fields[0] == fields[1] {
			return fmt.Errorf("Invalid kvPairs member declaration")
		}
		if previous, exists := groups[binding.Group]; exists && previous != fields {
			return fmt.Errorf("Conflicting kvPairs member declaration")
		}
		groups[binding.Group] = fields
	}
	for group, fields := range groups {
		raw, exists := mapped[group]
		if !exists {
			continue
		}
		members, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("kvPairs group is not an indexed object")
		}
		indexes := []int{}
		for key := range members {
			index, err := strconv.Atoi(key)
			if err != nil || index < 1 || strconv.Itoa(index) != key {
				return fmt.Errorf("Invalid kvPairs member index")
			}
			indexes = append(indexes, index)
		}
		sort.Ints(indexes)
		compact := map[string]interface{}{}
		codes := map[string]bool{}
		for position, index := range indexes {
			pair, valid := members[strconv.Itoa(index)].(map[string]interface{})
			if !valid || len(pair) != 2 || pair[fields[0]] == nil || pair[fields[1]] == nil {
				return fmt.Errorf("Incomplete kvPairs member")
			}
			code, valid := pair[fields[0]].(string)
			if !valid || code == "" || codes[code] {
				return fmt.Errorf("Invalid or duplicate kvPairs code")
			}
			codes[code] = true
			compact[strconv.Itoa(position+1)] = pair
		}
		mapped[group] = compact
	}
	return nil
}

func fwFlattenRequest(request map[string]interface{}, key string, raw interface{}) {
	switch v := raw.(type) {
	case []interface{}:
		for i, item := range v {
			fwFlattenRequest(request, fmt.Sprintf("%s.%d", key, i+1), item)
		}
	case map[string]interface{}:
		for name, item := range v {
			fwFlattenRequest(request, key+"."+name, item)
		}
	default:
		if v != nil {
			request[key] = v
		}
	}
}

func fwNormalizeResponse(object map[string]interface{}, bindings []fwMapping) map[string]interface{} {
	// Whole-object mappings must precede their leaf aliases, regardless of AM order.
	// Stable insertion order preserves same-depth precedence and source array indexes.
	bindings = append([]fwMapping(nil), bindings...)
	for index := 1; index < len(bindings); index++ {
		for current := index; current > 0 && len(fwMappingPath(bindings[current].Target)) < len(fwMappingPath(bindings[current-1].Target)); current-- {
			bindings[current], bindings[current-1] = bindings[current-1], bindings[current]
		}
	}
	out := map[string]interface{}{}
	for key, value := range object {
		out[key] = value
	}
	for _, binding := range bindings {
		if strings.HasSuffix(binding.Source, "[*]") && strings.HasSuffix(binding.Target, "[*]") {
			source := strings.TrimSuffix(binding.Source, "[*]")
			target := strings.TrimSuffix(binding.Target, "[*]")
			for _, container := range fwWalkMappedPath(object, fwMappingPath(source), nil) {
				if values, ok := container.Value.([]interface{}); ok && len(values) == 0 {
					out = fwSetMappedPath(out, fwMappingPath(target), container.Indexes, []interface{}{}).(map[string]interface{})
				}
			}
		}
		for _, entry := range fwWalkMappedPath(object, fwMappingPath(binding.Source), nil) {
			raw := entry.Value
			target := binding.Target
			if binding.Encoding == "split-comma" {
				target = strings.TrimSuffix(target, "[*]")
				items := []interface{}{}
				if text, ok := raw.(string); ok && text != "" {
					for _, part := range strings.Split(text, ",") {
						items = append(items, strings.TrimSpace(part))
					}
				}
				raw = items
			}
			if binding.Encoding == "scalar-string" && raw != nil {
				switch raw.(type) {
				case string:
				case int, int32, int64, float32, float64, json.Number:
					raw = fmt.Sprint(raw)
				}
			}
			if mapped, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mapped
			}
			out = fwSetMappedPath(out, fwMappingPath(target), entry.Indexes, raw).(map[string]interface{})
		}
	}
	return out
}

func fwValuesEqual(a, b interface{}) bool {
	aa, errA := json.Marshal(a)
	bb, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aa) == string(bb)
}
