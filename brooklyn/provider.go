package brooklyn

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jittakal/terraform-provider-brooklyn/brooklyn/config"
)

// Provider returns schema.Provider
func Provider() *schema.Provider {
	log.Printf("[Info] Provider for brooklyn")
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["secret_key"],
			},
			"endpoint_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["endpoint_url"],
			},
			"skip_ssl_checks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_ssl_checks"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"brooklyn_application": resourceApplication(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key":      "Access key",
		"secret_key":      "Secret Key",
		"endpoint_url":    "Endpoint URL",
		"skip_ssl_checks": "Skip SSL Checks",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := config.Config{
		AccessKey:     d.Get("access_key").(string),
		SecretKey:     d.Get("secret_key").(string),
		EndpointURL:   d.Get("endpoint_url").(string),
		SkipSslChecks: d.Get("skip_ssl_checks").(bool),
	}

	return config.Client()
}
