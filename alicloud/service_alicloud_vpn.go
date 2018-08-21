package alicloud

import (
	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

const (
	Ssl_Cert_Expiring = Status("expiring-soon")
	Ssl_Cert_Normal   = Status("normal")
	Ssl_Cert_Expired  = Status("expired")
)

func (client *AliyunClient) DescribeVpn(vpnId string) (v vpc.DescribeVpnGatewayResponse, err error) {
	request := vpc.CreateDescribeVpnGatewayRequest()
	request.VpnGatewayId = vpnId

	resp, err := client.vpcconn.DescribeVpnGateway(request)
	if err != nil {
		if IsExceptedError(err, VpnForbidden) || IsExceptedError(err, VpnNotFound) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", vpnId))
		}
		return
	}
	if resp == nil || resp.VpnGatewayId != vpnId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", vpnId))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeCustomerGateway(cgwId string) (v vpc.DescribeCustomerGatewayResponse, err error) {
	request := vpc.CreateDescribeCustomerGatewayRequest()
	request.CustomerGatewayId = cgwId

	resp, err := client.vpcconn.DescribeCustomerGateway(request)
	if err != nil {
		if IsExceptedError(err, VpnForbidden) || IsExceptedError(err, VpnNotFound) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", cgwId))
		}
		return
	}
	if resp == nil || resp.CustomerGatewayId != cgwId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", cgwId))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeVpnConnection(id string) (v vpc.DescribeVpnConnectionResponse, err error) {
	request := vpc.CreateDescribeVpnConnectionRequest()
	request.VpnConnectionId = id

	resp, err := client.vpcconn.DescribeVpnConnection(request)
	if err != nil {
		if IsExceptedError(err, VpnForbidden) || IsExceptedError(err, VpnNotFound) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
		}
		return
	}
	if resp == nil || resp.VpnConnectionId != id {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeSslVpnServers(vpnId string, sslId string) (v vpc.DescribeSslVpnServersResponse, err error) {
	request := vpc.CreateDescribeSslVpnServersRequest()
	if sslId != "" {
		request.SslVpnServerId = sslId
	}

	if vpnId != "" {
		request.VpnGatewayId = vpnId
	}
	resp, err := client.vpcconn.DescribeSslVpnServers(request)
	if err != nil {
		if IsExceptedError(err, VpnForbidden) || IsExceptedError(err, VpnNotFound) || IsExceptedError(err, SslVpnServerNotFound) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", sslId))
		}
		return
	}

	if resp == nil || 0 == len(resp.SslVpnServers.SslVpnServer) {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", sslId))
	}

	if sslId != "" && sslId != resp.SslVpnServers.SslVpnServer[0].SslVpnServerId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", sslId))
	}

	return *resp, nil
}

func (client *AliyunClient) DescribeSslVpnClientCert(id string) (v vpc.DescribeSslVpnClientCertResponse, err error) {
	request := vpc.CreateDescribeSslVpnClientCertRequest()
	request.SslVpnClientCertId = id

	resp, err := client.vpcconn.DescribeSslVpnClientCert(request)
	if err != nil {
		if IsExceptedError(err, VpnForbidden) || IsExceptedError(err, VpnNotFound) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
		}
		return
	}
	if resp == nil || resp.SslVpnClientCertId != id {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
	}
	return *resp, nil
}

func (client *AliyunClient) WaitForVpn(vpnId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		//wait the order effective
		vpn, err := client.DescribeVpn(vpnId)
		if err != nil {
			return err
		}
		if strings.ToLower(vpn.Status) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (client *AliyunClient) WaitForCustomerGateway(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		_, err := client.DescribeCustomerGateway(id)
		if err != nil {
			return err
		} else {
			break
		}
	}
	return nil
}

func (client *AliyunClient) WaitForSslVpnClientCert(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		resp, err := client.DescribeSslVpnClientCert(id)
		if err != nil {
			return err
		}

		if strings.ToLower(resp.Status) == strings.ToLower(string(status)) {
			break
		}
	}
	return nil
}
