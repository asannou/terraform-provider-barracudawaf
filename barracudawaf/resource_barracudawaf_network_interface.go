package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFNetworkInterfaceUpdate,
		Read:   resourceCudaWAFNetworkInterfaceRead,
		Update: resourceCudaWAFNetworkInterfaceUpdate,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name":                    {Type: schema.TypeString, Required: true, Description: "NIC Card Name"},
			"auto_negotiation_status": {Type: schema.TypeString, Optional: true, Description: "Auto-Negotiation Status"},
			"duplexity":               {Type: schema.TypeString, Required: true, Description: "NIC Cards Duplexity"},
			"speed":                   {Type: schema.TypeString, Required: true, Description: "NIC Speed"},
		},

		Description: "`barracudawaf_network_interface` manages `Network Interface` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFNetworkInterfaceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/network-interfaces"
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
	d.Set("auto_negotiation_status", dataItems["auto-negotiation-status"])
	d.Set("duplexity", dataItems["duplexity"])
	d.Set("speed", dataItems["speed"])

	return nil
}

func resourceCudaWAFNetworkInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)
	if d.Id() == "" {
		d.SetId(name)
	}

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/network-interfaces"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFNetworkInterfaceResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFNetworkInterfaceRead(d, m)
}

func hydrateBarracudaWAFNetworkInterfaceResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                    d.Get("name").(string),
		"auto-negotiation-status": d.Get("auto_negotiation_status").(string),
		"duplexity":               d.Get("duplexity").(string),
		"speed":                   d.Get("speed").(string),
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
