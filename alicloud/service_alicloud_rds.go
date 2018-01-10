package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
)

//
//       _______________                      _______________                       _______________
//       |              | ______param______\  |              |  _____request_____\  |              |
//       |   Business   |                     |    Service   |                      |    SDK/API   |
//       |              | __________________  |              |  __________________  |              |
//       |______________| \    (obj, err)     |______________|  \ (status, cont)    |______________|
//                           |                                    |
//                           |A. {instance, nil}                  |a. {200, content}
//                           |B. {nil, error}                     |b. {200, nil}
//                      					  |c. {4xx, nil}
//
// The API return 200 for resource not found.
// When getInstance is empty, then throw InstanceNotfound error.
// That the business layer only need to check error.
func (client *AliyunClient) DescribeDBInstanceById(id string) (instance *rds.DBInstanceAttribute, err error) {
	arrtArgs := rds.DescribeDBInstancesArgs{
		DBInstanceId: id,
	}
	resp, err := client.rdsconn.DescribeDBInstanceAttribute(&arrtArgs)
	if err != nil {
		return nil, err
	}

	attr := resp.Items.DBInstanceAttribute

	if len(attr) <= 0 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s is not found.", id))
	}

	return &attr[0], nil
}

func (client *AliyunClient) DescribeDatabaseAccount(instanceId, accountName string) (ds *rds.DBInstanceAccount, err error) {
	conn := client.rdsconn
	resp, err := conn.DescribeAccounts(&rds.DescribeAccountsArgs{
		DBInstanceId: instanceId,
		AccountName:  accountName,
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Accounts.DBInstanceAccount) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("Data account %s is not found in the instance %s.", accountName, instanceId))
	}
	return &resp.Accounts.DBInstanceAccount[0], nil
}

func (client *AliyunClient) DescribeDatabaseByName(instanceId, dbName string) (ds *rds.Database, err error) {

	resp, err := client.rdsconn.DescribeDatabases(&rds.DescribeDatabasesArgs{
		DBInstanceId: instanceId,
		DBName:       dbName,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Databases.Database) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("Database %s is not found in the instance %s.", dbName, instanceId))
	}
	return &resp.Databases.Database[0], nil
}

func (client *AliyunClient) AllocateDBPublicConnection(instanceId, prefix, port string) error {
	conn := client.rdsconn
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		if _, err := conn.AllocateInstancePublicConnection(&rds.AllocateInstancePublicConnectionArgs{
			DBInstanceId:           instanceId,
			ConnectionStringPrefix: prefix,
			Port: port,
		}); err != nil {
			if IsExceptedError(err, NetTypeExists) {
				connection, err := client.DescribeDBInstanceNetInfoByIpType(instanceId, rds.Public)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return resource.NonRetryableError(fmt.Errorf("The connection string with specified prefix %s has already existed. "+
					"Please import it using ID '%s:%s' or specify a new 'connection_prefix' and try again.", prefix, instanceId, connection.ConnectionString))
			} else if IsExceptedError(err, OperationDeniedDBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Allocate db connection got an error: %#v.", err))
			}

			return resource.NonRetryableError(fmt.Errorf("Allocate db connection got an error: %#v.", err))
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := conn.WaitForDBConnection(instanceId, rds.Public, 300); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) DescribeDBInstanceNetInfos(instanceId string) ([]rds.DBInstanceNetInfo, error) {

	resp, err := client.rdsconn.DescribeDBInstanceNetInfo(&rds.DescribeDBInstanceNetInfoArgs{
		DBInstanceId: instanceId,
	})

	if err != nil {
		return nil, err
	}

	if len(resp.DBInstanceNetInfos.DBInstanceNetInfo) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have any connection.", instanceId))
	}

	return resp.DBInstanceNetInfos.DBInstanceNetInfo, nil
}

func (client *AliyunClient) DescribeDBInstanceNetInfoByIpType(instanceId string, ipType rds.IPType) (*rds.DBInstanceNetInfo, error) {

	resps, err := client.DescribeDBInstanceNetInfos(instanceId)

	if err != nil {
		return nil, err
	}

	if len(resps) < 1 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have any connection.", instanceId))
	}

	for _, conn := range resps {
		if conn.IPType == ipType {
			return &conn, nil
		}
	}

	return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s does not have specified type %s connection.", instanceId, ipType))
}

func (client *AliyunClient) GrantAccountPrivilege(instanceId, account, dbName, privilege string) error {
	args := rds.GrantAccountPrivilegeArgs{
		DBInstanceId:     instanceId,
		AccountName:      account,
		DBName:           dbName,
		AccountPrivilege: rds.AccountPrivilege(privilege),
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ag := args
		if _, err := client.rdsconn.GrantAccountPrivilege(&ag); err != nil {
			if IsExceptedError(err, OperationDeniedDBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Grant DB %s account %s privilege got an error: %#v.", dbName, account, err))
			}
			return resource.NonRetryableError(fmt.Errorf("Grant DB %s account %s privilege got an error: %#v.", dbName, account, err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := client.rdsconn.WaitForAccountPrivilege(instanceId, account, dbName, rds.AccountPrivilege(privilege), 200); err != nil {
		return fmt.Errorf("Wait for grantting DB %s account %s privilege got an error: %#v.", dbName, account, err)
	}

	return nil
}

func (client *AliyunClient) RevokeAccountPrivilege(instanceId, account, dbName string) error {
	args := rds.RevokeAccountPrivilegeArgs{
		DBInstanceId: instanceId,
		AccountName:  account,
		DBName:       dbName,
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		ag := args
		if err := client.rdsconn.RevokeAccountPrivilege(&ag); err != nil {
			if IsExceptedError(err, OperationDeniedDBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("Revoke DB %s account %s privilege got an error: %#v.", dbName, account, err))
			}
			return resource.NonRetryableError(fmt.Errorf("Revoke DB %s account %s privilege got an error: %#v.", dbName, account, err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := client.rdsconn.WaitForAccountPrivilegeRevoked(instanceId, account, dbName, 200); err != nil {
		return fmt.Errorf("Wait for revoking DB %s account %s privilege got an error: %#v.", dbName, account, err)
	}

	return nil
}

func (client *AliyunClient) ReleaseDBPublicConnection(instanceId, connection string) error {
	conn := client.rdsconn

	if err := conn.ReleaseInstancePublicConnection(&rds.ReleaseInstancePublicConnectionArgs{
		DBInstanceId:            instanceId,
		CurrentConnectionString: connection,
	}); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) SwitchDBInstanceNetType(instanceId, prefix string, port int, vswitchId string) error {
	conn := client.rdsconn

	if err := conn.SwitchDBInstanceNetType(&rds.SwitchDBInstanceNetTypeArgs{
		DBInstanceId:           instanceId,
		ConnectionStringPrefix: prefix,
		Port: port,
	}); err != nil {
		return err
	}

	ipType := rds.Private
	if vswitchId == "" {
		ipType = rds.Inner
	}

	if err := conn.WaitForDBConnection(instanceId, ipType, 300); err != nil {
		return err
	}

	return nil
}

func (client *AliyunClient) ConfigDBBackup(instanceId, backupTime, backupPeriod string, retentionPeriod int) error {
	bargs := rds.BackupPolicy{
		PreferredBackupTime:   backupTime,
		PreferredBackupPeriod: backupPeriod,
		BackupRetentionPeriod: retentionPeriod,
	}
	args := rds.ModifyBackupPolicyArgs{
		DBInstanceId: instanceId,
		BackupPolicy: bargs,
	}

	if _, err := client.rdsconn.ModifyBackupPolicy(&args); err != nil {
		return err
	}

	if err := client.rdsconn.WaitForInstance(instanceId, rds.Running, 600); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) ModifyDBSecurityIps(instanceId, ips string) error {
	args := rds.ModifySecurityIpsArgs{
		DBInstanceId: instanceId,
		SecurityIps:  ips,
	}

	if _, err := client.rdsconn.ModifySecurityIps(&args); err != nil {
		return err
	}

	if err := client.rdsconn.WaitForInstance(instanceId, rds.Running, 600); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) DescribeDBSecurityIps(instanceId string) (ips []rds.DBInstanceIPList, err error) {
	args := rds.DescribeDBInstanceIPsArgs{
		DBInstanceId: instanceId,
	}

	resp, err := client.rdsconn.DescribeDBInstanceIPs(&args)
	if err != nil {
		return nil, err
	}
	return resp.Items.DBInstanceIPArray, nil
}

func (client *AliyunClient) GetSecurityIps(instanceId string, securityIps interface{}) ([]string, error) {
	arr, err := client.DescribeDBSecurityIps(instanceId)
	if err != nil {
		return nil, err
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range arr {
		ips += separator + ip.SecurityIPList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ips, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	// Sort security ips according to security_ips's order
	var finalIps []string
	if securityIps != nil {
		ipList := expandStringList(securityIps.([]interface{}))
		for _, ip := range ipList {
			if _, ok := ipsMap[ip]; ok {
				finalIps = append(finalIps, ip)
				delete(ipsMap, ip)
				continue
			}
			finalIps = append(finalIps, "")
		}
	}

	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (client *AliyunClient) ModifyDBClassStorage(instanceId, class, storage string) error {
	conn := client.rdsconn
	args := rds.ModifyDBInstanceSpecArgs{
		DBInstanceId:      instanceId,
		PayType:           rds.Postpaid,
		DBInstanceClass:   class,
		DBInstanceStorage: storage,
	}

	if _, err := conn.ModifyDBInstanceSpec(&args); err != nil {
		return err
	}

	if err := conn.WaitForInstance(instanceId, rds.Running, 600); err != nil {
		return err
	}
	return nil
}

// turn period to TimeType
func TransformPeriod2Time(period int, chargeType string) (ut int, tt common.TimeType) {
	if chargeType == string(rds.Postpaid) {
		return 1, common.Day
	}

	if period >= 1 && period <= 9 {
		return period, common.Month
	}

	if period == 12 {
		return 1, common.Year
	}

	if period == 24 {
		return 2, common.Year
	}
	return 0, common.Day

}

// turn TimeType to Period
func TransformTime2Period(ut int, tt common.TimeType) (period int) {
	if tt == common.Year {
		return 12 * ut
	}

	return ut

}

// Flattens an array of databases into a []map[string]interface{}
func flattenDatabaseMappings(list []rds.Database) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"db_name":            i.DBName,
			"character_set_name": i.CharacterSetName,
			"db_description":     i.DBDescription,
		}
		result = append(result, l)
	}
	return result
}

func flattenDBBackup(list []rds.BackupPolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"preferred_backup_period": i.PreferredBackupPeriod,
			"preferred_backup_time":   i.PreferredBackupTime,
			"backup_retention_period": i.LogBackupRetentionPeriod,
		}
		result = append(result, l)
	}
	return result
}

func flattenDBSecurityIPs(list []rds.DBInstanceIPList) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"security_ips": i.SecurityIPList,
		}
		result = append(result, l)
	}
	return result
}

// Flattens an array of databases connection into a []map[string]interface{}
func flattenDBConnections(list []rds.DBInstanceNetInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"connection_string": i.ConnectionString,
			"ip_type":           i.IPType,
			"ip_address":        i.IPAddress,
		}
		result = append(result, l)
	}
	return result
}
