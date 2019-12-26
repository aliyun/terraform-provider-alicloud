package alicloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type PolarDBService struct {
	client *connectivity.AliyunClient
}

func (s *PolarDBService) DescribePolarDBCluster(id string) (instance *polardb.DBCluster, err error) {
	request := polardb.CreateDescribeDBClustersRequest()
	request.RegionId = s.client.RegionId
	dbClusterIds := []string{}
	dbClusterIds = append(dbClusterIds, id)
	request.DBClusterIds = id
	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusters(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDBClustersResponse)
	if len(response.Items.DBCluster) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items.DBCluster[0], nil
}

func (s *PolarDBService) DescribePolarDBClusterAttribute(id string) (instance *polardb.DescribeDBClusterAttributeResponse, err error) {
	request := polardb.CreateDescribeDBClusterAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAttribute(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDBClusterAttributeResponse)
	if len(response.DBClusterId) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return response, nil
}

func (s *PolarDBService) DescribePolarDBAutoRenewAttribute(id string) (instance *polardb.AutoRenewAttribute, err error) {
	request := polardb.CreateDescribeAutoRenewAttributeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterIds = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeAutoRenewAttribute(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeAutoRenewAttributeResponse)
	if len(response.Items.AutoRenewAttribute) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.Items.AutoRenewAttribute[0], nil
}

func (s *PolarDBService) DescribeParameters(id string) (ds *polardb.DescribeDBClusterParametersResponse, err error) {
	request := polardb.CreateDescribeDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterParameters(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*polardb.DescribeDBClusterParametersResponse)
	return response, err
}

func (s *PolarDBService) GrantPolarDBAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateGrantAccountPrivilegeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName
	request.AccountPrivilege = parts[2]

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.GrantAccountPrivilege(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := s.WaitForPolarDBAccountPrivilege(id, dbName, Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *PolarDBService) RevokePolarDBAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateRevokeAccountPrivilegeRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.RevokeAccountPrivilege(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if err := s.WaitForPolarDBAccountPrivilegeRevoked(id, dbName, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *PolarDBService) WaitForPolarDBAccountPrivilegeRevoked(id, dbName string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, ProviderERROR)
		}

	}
	return nil
}

func (s *PolarDBService) WaitForPolarDBAccountPrivilege(id, dbName string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccountPrivilege(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) DescribePolarDBAccountPrivilege(id string) (account *polardb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := polardb.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *polardb.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeAccounts(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*polardb.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBAccountPrivilege", id)), NotFoundMsg, ProviderERROR)
	}
	return &response.Accounts[0], nil
}

func (s *PolarDBService) WaitForPolarDBConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBConnection(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if status != Deleted && object != nil && object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *PolarDBService) DescribePolarDBConnection(id string) (*polardb.Address, error) {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(DefaultIntervalLong) * time.Second)
	for {
		object, err := s.DescribePolarDBInstanceNetInfo(parts[0])

		if err != nil {
			if NotFoundError(err) {
				return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return nil, WrapError(err)
		}

		if object != nil {
			for _, o := range object {
				if o.DBEndpointId == parts[1] {
					for _, p := range o.AddressItems {
						if p.NetType == parts[2] {
							return &p, nil
						}
					}
				}
			}
		}
		time.Sleep(DefaultIntervalMini * time.Second)
		if time.Now().After(deadline) {
			break
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("DBConnection", id)), NotFoundMsg, ProviderERROR)
}

func (s *PolarDBService) DescribePolarDBInstanceNetInfo(id string) ([]polardb.DBEndpoint, error) {

	request := polardb.CreateDescribeDBClusterEndpointsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = id
	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterEndpoints(request)
	})

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*polardb.DescribeDBClusterEndpointsResponse)
	if len(response.Items) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBInstanceNetInfo", id)), NotFoundMsg, ProviderERROR)
	}

	return response.Items, nil
}

func (s *PolarDBService) DescribePolarDBDatabase(id string) (ds *polardb.Database, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	dbName := parts[1]
	request := polardb.CreateDescribeDatabasesRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.DBName = dbName

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeDatabases(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDatabasesResponse)
	if len(response.Databases.Database) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBDatabase", dbName)), NotFoundMsg, ProviderERROR)
	}
	ds = &response.Databases.Database[0]
	return ds, nil
}

func (s *PolarDBService) WaitForPolarDBDatabase(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBDatabase(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}
		if status != Deleted && object != nil && object.DBName == parts[1] {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBName, parts[1], ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) WaitForPolarDBAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBAccount(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) DescribePolarDBAccount(id string) (ds *polardb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := polardb.CreateDescribeAccountsRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var response *polardb.DescribeAccountsResponse
	if err := invoker.Run(func() error {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response, _ = raw.(*polardb.DescribeAccountsResponse)
		return nil
	}); err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	if len(response.Accounts) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("DBAccount", id)), NotFoundMsg, ProviderERROR)
	}
	return &response.Accounts[0], nil
}

// WaitForInstance waits for instance to given status
func (s *PolarDBService) WaitForPolarDBInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBCluster(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) RefreshParameters(d *schema.ResourceData, attribute string) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk(attribute)
	if !ok {
		d.Set(attribute, param)
		return nil
	}
	object, err := s.DescribeParameters(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var parameters = make(map[string]interface{})
	for _, i := range object.RunningParameters.Parameter {
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
	if err := d.Set(attribute, param); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *PolarDBService) ModifyParameters(d *schema.ResourceData, attribute string) error {
	request := polardb.CreateModifyDBClusterParametersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = d.Id()
	config := make(map[string]string)
	documented := d.Get(attribute).(*schema.Set).List()
	if len(documented) > 0 {
		for _, i := range documented {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request.Parameters = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForCluster(d.Id(), Running, DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterParameters(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance parameter expect after modifying
		if err := s.WaitForPolarDBParameter(d.Id(), DefaultTimeoutMedium, config); err != nil {
			return WrapError(err)
		}
	}
	d.SetPartial(attribute)
	return nil
}

// WaitForCluster waits for cluster to given status
func (s *PolarDBService) WaitForCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PolarDBService) DescribeDBSecurityIps(clusterId string) (ips []string, err error) {

	request := polardb.CreateDescribeDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeDBClusterAccessWhitelist(request)
	})
	if err != nil {
		return ips, WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	resp, _ := raw.(*polardb.DescribeDBClusterAccessWhitelistResponse)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range resp.Items.DBClusterIPArray {
		if ip.DBClusterIPArrayAttribute != "hidden" {
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

func (s *PolarDBService) ModifyDBSecurityIps(clusterId, ips string) error {

	request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.SecurityIps = ips

	raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.ModifyDBClusterAccessWhitelist(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *PolarDBService) DescribeBackupPolicy(id string) (policy *polardb.DescribeBackupPolicyResponse, err error) {

	request := polardb.CreateDescribeBackupPolicyRequest()
	request.DBClusterId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.DescribeBackupPolicy(request)
	})

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return policy, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return raw.(*polardb.DescribeBackupPolicyResponse), nil
}

func (s *PolarDBService) ModifyDBBackupPolicy(clusterId, backupTime, backupPeriod string) error {

	request := polardb.CreateModifyBackupPolicyRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterId = clusterId
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime

	raw, err := s.client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
		return polardbClient.ModifyBackupPolicy(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, clusterId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return nil
}

func (s *PolarDBService) PolarDBClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePolarDBClusterAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBClusterStatus == failState {
				return object, object.DBClusterStatus, WrapError(Error(FailedToReachTargetStatus, object.DBClusterStatus))
			}
		}
		return object, object.DBClusterStatus, nil
	}
}

// WaitForDBParameter waits for instance parameter to given value.
// Status of DB instance is Running after ModifyParameters API was
// call, so we can not just wait for instance status become
// Running, we should wait until parameters have expected values.
func (s *PolarDBService) WaitForPolarDBParameter(clusterId string, timeout int, expects map[string]string) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeParameters(clusterId)
		if err != nil {
			return WrapError(err)
		}

		var actuals = make(map[string]string)
		for _, i := range object.RunningParameters.Parameter {
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

		time.Sleep(DefaultIntervalShort * time.Second)

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, clusterId, GetFunc(1), timeout, got_value, expected_value, ProviderERROR)
		}
	}
	return nil
}
