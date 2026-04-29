package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSystemUpdate,
		Read:   resourceCudaWAFSystemRead,
		Update: resourceCudaWAFSystemUpdate,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Default Domain",
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "" {
						return
					}
					validLocales := []string{
						"한국어", "Dansk", "日本語", "Italiano", "Español", "Deutsch", "Magyar", "Français", "Русский", "English", "Česky", "繁體中文", "Polski", "Catalan", "íslenska", "简体中文", "Português (BR)", "Nederlands",
					}
					found := false
					for _, locale := range validLocales {
						if v == locale {
							found = true
							break
						}
					}
					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validLocales))
					}
					return
				},
				Description: "Default Language and Encoding",
			},
			"interface_for_system_services": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "WAN" && v != "Management" && v != "" {
						errs = append(errs, fmt.Errorf("%q must be either WAN, Management or empty", key))
					}
					return
				},
				Description: "Interface for System Services",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Time Zone",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default Hostname",
			},
			"enable_ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Yes" && v != "No" && v != "" {
						errs = append(errs, fmt.Errorf("%q must be either Yes, No or empty", key))
					}
					return
				},
				Description: "Enable IPv6",
			},
			"device_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Firewall Name",
			},
			"operation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Bridge All Traffic" && v != "Proxy" && v != "" {
						errs = append(errs, fmt.Errorf("%q must be either 'Bridge All Traffic', 'Proxy' or empty", key))
					}
					return
				},
				Description: "Mode of Operation",
			},
		},

		Description: "`barracudawaf_system` manages `System` settings on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSystemRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/system"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	resources, err := client.GetBarracudaWAFResource(name, request)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	if resources.Data == nil {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var dataItems map[string]interface{}
	dataItems = resources.Data["system"]
	if dataItems == nil {
		for _, v := range resources.Data {
			dataItems = v
			break
		}
	}

	if dataItems == nil {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("domain", dataItems["domain"])
	d.Set("locale", dataItems["locale"])
	d.Set("interface_for_system_services", dataItems["interface-for-system-services"])
	d.Set("time_zone", dataItems["time-zone"])
	d.Set("hostname", dataItems["hostname"])
	d.Set("enable_ipv6", dataItems["enable-ipv6"])
	d.Set("device_name", dataItems["device-name"])
	d.Set("operation_mode", dataItems["operation-mode"])

	return nil
}

func resourceCudaWAFSystemUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := "system"
	if d.Id() == "" {
		d.SetId(name)
	}

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/system"
	request := hydrateBarracudaWAFSystemResource(d, resourceEndpoint)
	_, err := client.putReq(request.Body, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSystemRead(d, m)
}

func hydrateBarracudaWAFSystemResource(d *schema.ResourceData, endpoint string) *APIRequest {

	resourcePayload := map[string]string{
		"domain":                        d.Get("domain").(string),
		"locale":                        d.Get("locale").(string),
		"interface-for-system-services": d.Get("interface_for_system_services").(string),
		"time-zone":                     d.Get("time_zone").(string),
		"hostname":                      d.Get("hostname").(string),
		"enable-ipv6":                   d.Get("enable_ipv6").(string),
		"device-name":                   d.Get("device_name").(string),
		"operation-mode":                d.Get("operation_mode").(string),
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if len(val) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
