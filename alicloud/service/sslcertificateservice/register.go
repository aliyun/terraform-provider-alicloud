package sslcertificateservice

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the SslCertificatesService
// service. It carries the SDKv2 resources and data sources (including the legacy
// alicloud_cas_* aliases) and the framework Action form of
// alicloud_ssl_certificates_service_pca_cert_sync, which coexists with its SDKv2
// resource counterpart under separate registry namespaces.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "SslCertificatesService",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ssl_certificates_service_certificate", Factory: resourceAliCloudSslCertificatesServiceCertificate},
			{TypeName: "alicloud_ssl_certificates_service_certificate_apply", Factory: resourceAliCloudSslCertificatesServiceCertificateApply},
			{TypeName: "alicloud_ssl_certificates_service_certificate_validation", Factory: resourceAliCloudSslCertificatesServiceCertificateValidation},
			{TypeName: "alicloud_ssl_certificates_service_company", Factory: resourceAliCloudSslCertificatesServiceCompany},
			{TypeName: "alicloud_ssl_certificates_service_contact", Factory: resourceAliCloudSslCertificatesServiceContact},
			{TypeName: "alicloud_ssl_certificates_service_instance", Factory: resourceAliCloudSslCertificatesServiceInstance},
			{TypeName: "alicloud_ssl_certificates_service_instance_certificate", Factory: resourceAliCloudSslCertificatesServiceInstanceCertificate},
			{TypeName: "alicloud_ssl_certificates_service_pca_cert", Factory: resourceAliCloudSslCertificatesServicePcaCert},
			{TypeName: "alicloud_ssl_certificates_service_pca_certificate", Factory: resourceAliCloudSslCertificatesServicePcaCertificate},
			{TypeName: "alicloud_ssl_certificates_service_pca_cert_sync", Factory: resourceAliCloudSslCertificatesServicePcaCertSync},
			// Legacy alias kept for backward compatibility; reuses the certificate factory.
			{TypeName: "alicloud_cas_certificate", Factory: resourceAliCloudSslCertificatesServiceCertificate},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_ssl_certificates_service_certificates", Factory: dataSourceAliCloudSslCertificatesServiceCertificates},
			{TypeName: "alicloud_ssl_certificates_service_companies", Factory: dataSourceAliCloudSslCertificatesServiceCompanies},
			{TypeName: "alicloud_ssl_certificates_service_contacts", Factory: dataSourceAliCloudSslCertificatesServiceContacts},
			{TypeName: "alicloud_ssl_certificates_service_instances", Factory: dataSourceAliCloudSslCertificatesServiceInstances},
			{TypeName: "alicloud_ssl_certificates_service_instance_certificates", Factory: dataSourceAliCloudSslCertificatesServiceInstanceCertificates},
			// Legacy alias kept for backward compatibility; reuses the certificates factory.
			{TypeName: "alicloud_cas_certificates", Factory: dataSourceAliCloudSslCertificatesServiceCertificates},
		},
		Actions: []conns.Action{
			{TypeName: "alicloud_ssl_certificates_service_pca_cert_sync", Factory: NewPcaCertSyncAction},
		},
	}
}
