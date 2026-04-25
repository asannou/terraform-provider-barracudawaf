package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFRateControlPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFRateControlPoolCreate,
		Read:   resourceCudaWAFRateControlPoolRead,
		Update: resourceCudaWAFRateControlPoolUpdate,
		Delete: resourceCudaWAFRateControlPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name":                       {Type: schema.TypeString, Required: true, Description: "Rate Control Pool Name"},
			"max_active_requests":        {Type: schema.TypeString, Optional: true, Description: "Maximum Active Requests"},
			"max_per_client_backlog":     {Type: schema.TypeString, Optional: true, Description: "Maximum Per Client Backlog"},
			"max_unconfigured_clients":   {Type: schema.TypeString, Optional: true, Description: "Maximum Unconfigured Clients"},
		},

		Description: "`barracudawaf_rate_control_pool` manages `Rate Control Pool` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFRateControlPoolCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/rate-control-pools"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFRateControlPoolResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFRateControlPoolRead(d, m)
}

func resourceCudaWAFRateControlPoolRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/rate-control-pools"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	var dataItems map[string]interface{}
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

	for _, dataItems = range resources.Data {
		if dataItems["name"] == name {
			break
		}
	}

	if dataItems["name"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("name", name)
	d.Set("max_active_requests", fmt.Sprintf("%v", dataItems["max-active-requests"]))
	d.Set("max_per_client_backlog", fmt.Sprintf("%v", dataItems["max-per-client-backlog"]))
	d.Set("max_unconfigured_clients", fmt.Sprintf("%v", dataItems["max-unconfigured-clients"]))

	return nil
}

func resourceCudaWAFRateControlPoolUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/rate-control-pools"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFRateControlPoolResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFRateControlPoolRead(d, m)
}

func resourceCudaWAFRateControlPoolDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/rate-control-pools"
	request := &APIRequest{
		Method: "delete",
		URL:    resourceEndpoint,
	}

	err := client.DeleteBarracudaWAFResource(name, request)

	if err != nil {
		return fmt.Errorf("Unable to delete the Barracuda WAF resource (%s) (%v)", name, err)
	}

	return nil
}

func hydrateBarracudaWAFRateControlPoolResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                       d.Get("name").(string),
		"max-active-requests":        d.Get("max_active_requests").(string),
		"max-per-client-backlog":     d.Get("max_per_client_backlog").(string),
		"max-unconfigured-clients":   d.Get("max_unconfigured_clients").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"name"}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
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
