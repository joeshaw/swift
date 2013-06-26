package swift

import (
	"bytes"
	"encoding/json"
	"io"
)

// V2 Authentication request
//
// http://docs.openstack.org/developer/keystone/api_curl_examples.html
// http://docs.rackspace.com/servers/api/v2/cs-gettingstarted/content/curl_auth.html
type V2AuthRequest struct {
	Auth struct {
		ApiKeyCredentials struct {
			UserName string `json:"username"`
			ApiKey   string `json:"apiKey"`
		} `json:"RAX-KSKEY:apiKeyCredentials"`
		PasswordCredentials struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		} `json:"passwordCredentials"`
	} `json:"auth"`
}

// V2 Authentication reply
//
// http://docs.openstack.org/developer/keystone/api_curl_examples.html
// http://docs.rackspace.com/servers/api/v2/cs-gettingstarted/content/curl_auth.html
type V2Auth struct {
	Access struct {
		ServiceCatalog []struct {
			Endpoints []struct {
				InternalUrl string
				PublicUrl   string
				Region      string
				TenantId    string
			}
			Name string
			Type string
		}
		Token struct {
			Expires string
			Id      string
			Tenant  struct {
				id   string
				name string
			}
		}
		User struct {
			DefaultRegion string `json:"RAX-AUTH:defaultRegion"`
			Id            string
			Name          string
			Roles         []struct {
				Description string
				Id          string
				Name        string
				TenantId    string
			}
		}
	}
}

// Create a V2 auth request for the body of the connection
//
// Adds both ApiKey and Password auth
func NewV2AuthRequest(UserName, ApiKey, Password string) (io.Reader, error) {
	v2 := V2AuthRequest{}
	v2.Auth.ApiKeyCredentials.UserName = UserName
	v2.Auth.ApiKeyCredentials.ApiKey = ApiKey
	v2.Auth.PasswordCredentials.UserName = UserName
	v2.Auth.PasswordCredentials.Password = Password
	body, err := json.Marshal(v2)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

// Finds the Endpoint Url of "type" from the V2Auth using the Region
// if set in Connection or defaulting to the first one if not
//
// Returns "" if not found
func (c *Connection) V2AuthEndpointUrl(Type string) string {
	if c.V2 {
		for _, catalog := range c.V2Auth.Access.ServiceCatalog {
			if catalog.Type == Type {
				for _, endpoint := range catalog.Endpoints {
					if c.Region == "" || (c.Region == endpoint.Region) {
						// FIXME could use PrivateUrl?
						return endpoint.PublicUrl
					}
				}
			}
		}
	}
	return ""
}
