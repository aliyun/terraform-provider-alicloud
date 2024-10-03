package providers

import (
	"fmt"
	"os"
	"strings"
)

type DefaultCredentialsProvider struct {
	providerChain    []CredentialsProvider
	lastUsedProvider CredentialsProvider
}

func NewDefaultCredentialsProvider() (provider *DefaultCredentialsProvider) {
	providers := []CredentialsProvider{}

	// Add static ak or sts credentials provider
	envProvider, err := NewEnvironmentVariableCredentialsProviderBuilder().Build()
	if err == nil {
		providers = append(providers, envProvider)
	}

	// oidc check
	oidcProvider, err := NewOIDCCredentialsProviderBuilder().Build()
	if err == nil {
		providers = append(providers, oidcProvider)
	}

	// cli credentials provider
	cliProfileProvider, err := NewCLIProfileCredentialsProviderBuilder().Build()
	if err == nil {
		providers = append(providers, cliProfileProvider)
	}

	// profile credentials provider
	profileProvider, err := NewProfileCredentialsProviderBuilder().Build()
	if err == nil {
		providers = append(providers, profileProvider)
	}

	// Add IMDS
	if os.Getenv("ALIBABA_CLOUD_ECS_METADATA") != "" {
		ecsRamRoleProvider, err := NewECSRAMRoleCredentialsProviderBuilder().WithRoleName(os.Getenv("ALIBABA_CLOUD_ECS_METADATA")).Build()
		if err == nil {
			providers = append(providers, ecsRamRoleProvider)
		}
	}

	// TODO: ALIBABA_CLOUD_CREDENTIALS_URI check

	return &DefaultCredentialsProvider{
		providerChain: providers,
	}
}

func (provider *DefaultCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	if provider.lastUsedProvider != nil {
		inner, err1 := provider.lastUsedProvider.GetCredentials()
		if err1 != nil {
			return
		}

		cc = &Credentials{
			AccessKeyId:     inner.AccessKeyId,
			AccessKeySecret: inner.AccessKeySecret,
			SecurityToken:   inner.SecurityToken,
			ProviderName:    fmt.Sprintf("%s/%s", provider.GetProviderName(), provider.lastUsedProvider.GetProviderName()),
		}
		return
	}

	errors := []string{}
	for _, p := range provider.providerChain {
		provider.lastUsedProvider = p
		inner, errInLoop := p.GetCredentials()
		if errInLoop != nil {
			errors = append(errors, errInLoop.Error())
			// 如果有错误，进入下一个获取过程
			continue
		}

		if inner != nil {
			cc = &Credentials{
				AccessKeyId:     inner.AccessKeyId,
				AccessKeySecret: inner.AccessKeySecret,
				SecurityToken:   inner.SecurityToken,
				ProviderName:    fmt.Sprintf("%s/%s", provider.GetProviderName(), p.GetProviderName()),
			}
			return
		}
	}

	err = fmt.Errorf("unable to get credentials from any of the providers in the chain: %s", strings.Join(errors, ", "))
	return
}

func (provider *DefaultCredentialsProvider) GetProviderName() string {
	return "default"
}
