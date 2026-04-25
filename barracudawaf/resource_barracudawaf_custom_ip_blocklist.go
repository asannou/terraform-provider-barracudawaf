package barracudawaf

import (
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceCustomIPBlocklistParams = map[string][]string{}
)

func resourceCudaWAFCustomIPBlocklist() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFCustomIPBlocklistCreate,
		Read:   resourceCudaWAFCustomIPBlocklistRead,
		Update: resourceCudaWAFCustomIPBlocklistUpdate,
		Delete: resourceCudaWAFCustomIPBlocklistDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"custom_ip_list": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom IP List",
			},
			"blacklisted_ips": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Upload File Of BlackListed IPs",
			},
			"download_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Download URL",
			},
			"validate_server_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Validate Server Certificate",
			},
			"trusted_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Trusted Certificate",
			},
		},

		Description: "`barracudawaf_custom_ip_blocklist` manages `Custom IP Blocklist` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFCustomIPBlocklistCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("custom_ip_list").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/custom-ip-blocklist"
	err := client.UpdateBarracudaWAFResource("", &APIRequest{
		URL:  resourceEndpoint,
		Body: hydrateBarracudaWAFCustomIPBlocklistResource(d, "put", resourceEndpoint).Body,
	})

	if err != nil {
		log.Printf("[ERROR] Unable to create/update Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId("custom-ip-blocklist")
	return resourceCudaWAFCustomIPBlocklistRead(d, m)
}

func resourceCudaWAFCustomIPBlocklistRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	log.Println("[INFO] Fetching Barracuda WAF resource custom-ip-blocklist")

	resourceEndpoint := "/custom-ip-blocklist"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	resources, err := client.GetBarracudaWAFResource("", request)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Barracuda WAF resource (%v) ", err)
		return err
	}

	if resources.Data == nil {
		log.Printf("[WARN] Barracuda WAF resource not found, removing from state")
		d.SetId("")
		return nil
	}

	dataItems := resources.Data["custom-ip-blocklist"]
	if dataItems == nil {
		// Try another way to match
		for _, v := range resources.Data {
			dataItems = v
			break
		}
	}

	if dataItems != nil {
		d.Set("custom_ip_list", dataItems["custom-ip-list"])
		d.Set("blacklisted_ips", dataItems["blacklisted-ips"])
		d.Set("download_url", dataItems["download-url"])
		d.Set("validate_server_certificate", dataItems["validate-server-certificate"])
		d.Set("trusted_certificate", dataItems["trusted-certificate"])
	}

	return nil
}

func resourceCudaWAFCustomIPBlocklistUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceCudaWAFCustomIPBlocklistCreate(d, m)
}

func resourceCudaWAFCustomIPBlocklistDelete(d *schema.ResourceData, m interface{}) error {
	// Custom IP Blocklist is a singleton configuration, usually we don't "delete" it but might reset it.
	// WAF API doesn't show a DELETE method for this.
	d.SetId("")
	return nil
}

func hydrateBarracudaWAFCustomIPBlocklistResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"custom-ip-list":              d.Get("custom_ip_list").(string),
		"blacklisted-ips":             d.Get("blacklisted_ips").(string),
		"download-url":                d.Get("download_url").(string),
		"validate-server-certificate": d.Get("validate_server_certificate").(string),
		"trusted-certificate":         d.Get("trusted_certificate").(string),
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if reflect.ValueOf(val).Kind() == reflect.String && reflect.ValueOf(val).Len() == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
