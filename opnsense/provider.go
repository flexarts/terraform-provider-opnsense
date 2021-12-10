package opnsense

import (
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// ProviderConfiguration struct for opnsense-provider
type ProviderConfiguration struct {
	DHCP  *DHCPSession
	Mutex *sync.Mutex
	Cond  *sync.Cond
}

// Provider libvirt
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OPNSENSE_URI", nil),
				ValidateFunc: validation.IsURLWithHTTPS,
				Description:  "OPNsense platform URI",
			},
			"user": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OPNSENSE_USER_ID", nil),
				ValidateFunc: validation.All(validation.StringIsNotEmpty),
				Description:  "OPNsense platform user ID",
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OPNSENSE_USER_PASSWORD", nil),
				ValidateFunc: validation.All(validation.StringIsNotEmpty),
				Description:  "OPNsense platform user password",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"opnsense_dhcp_static_map": resourceOpnDHCPStaticMap(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	// check for mandatory requirements
	uri := d.Get("uri").(string)
	user := d.Get("user").(string)
	password := d.Get("password").(string)

	if uri == "" || user == "" || password == "" {
		return nil, fmt.Errorf("The opnsense provider needs proper initialization parameters")
	}

	var mut sync.Mutex
	var dhcp DHCPSession
	var provider = ProviderConfiguration{
		DHCP:  &dhcp,
		Mutex: &mut,
		Cond:  sync.NewCond(&mut),
	}
	err := provider.DHCP.Authenticate(uri, user, password)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to OPNSense")
	}

	return &provider, nil
}
