package brooklyn

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jittakal/terraform-provider-brooklyn/brooklyn/application"
	"github.com/jittakal/terraform-provider-brooklyn/brooklyn/config"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,

		Schema: map[string]*schema.Schema{
			"application_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"location_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"catalog_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"catalog_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"wait_for_running": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"application_configuration": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			// "brooklyn_tags": &schema.Schema{
			//	Type:     schema.TypeList,
			//	Optional:  false,
			// },

		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] resource application create")
	waitForRunning := d.Get("wait_for_running").(bool)

	appln := &application.Application{
		Client:   m.(config.Client),
		Name:     d.Get("application_name").(string),
		Location: d.Get("location_name").(string),
		Type:     fmt.Sprintf("%s:%s", d.Get("catalog_id").(string), d.Get("catalog_version").(string)),
	}

	if v, ok := d.GetOk("application_configuration"); ok {
		configurations := configurationsFromMap(v.(map[string]interface{}))
		appln.Configurations = configurations
	}

	// TODO: preCreateValudation ()
	// Application name already exists
	// Location name not found
	// catalog does not exits
	// catalog exists but disabled

	err := appln.Create()
	if err != nil {
		return err
	}
	d.SetId(appln.ID)

	if waitForRunning {
		// Wait for instance starting
		err = appln.WaitForApplicationStarting()
		if err != nil {
			return err
		}

		// Wait for instance running
		err = appln.WaitForApplicationRunning()
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceApplicationRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceApplicationUpdate(d *schema.ResourceData, m interface{}) error {

	if d.HasChange("application_name") {
		appln := &application.Application{
			Client:   m.(config.Client),
			ID:       d.Id(),
			Name:     d.Get("application_name").(string),
			Location: d.Get("location_name").(string),
			Type:     fmt.Sprintf("%s:%s", d.Get("catalog_id").(string), d.Get("catalog_version").(string)),
		}

		return appln.Rename()
	}

	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, m interface{}) error {
	appln := &application.Application{
		Client: m.(config.Client),
		ID:     d.Id(),
	}

	err := appln.Delete()
	if err == nil {
		d.SetId("")
	}

	return err
}

// configurationsFromMap returns the configurations for the given map of data.
func configurationsFromMap(m map[string]interface{}) []application.Configuration {
	result := make([]application.Configuration, 0, len(m))
	for k, v := range m {
		t := application.Configuration{
			Key:   k,
			Value: v.(string),
		}
		result = append(result, t)
	}

	return result
}
