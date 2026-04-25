package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFInterfaceRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFInterfaceRouteCreate,
		Read:   resourceCudaWAFInterfaceRouteRead,
		Update: resourceCudaWAFInterfaceRouteUpdate,
		Delete: resourceCudaWAFInterfaceRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {Type: schema.TypeString, Required: true, Description: "IP/Network Address"},
			"interface":  {Type: schema.TypeString, Required: true, Description: "Network Interface"},
			"netmask":    {Type: schema.TypeString, Required: true, Description: "Netmask"},
			"vsite":      {Type: schema.TypeString, Required: true, Description: "Network Group"},
			"ip_version": {Type: schema.TypeString, Optional: true, Description: "IP Version"},
			"comments":   {Type: schema.TypeString, Optional: true, Description: "Comments"},
		},

		Description: "`barracudawaf_interface_route` manages `Interface Route` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFInterfaceRouteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("ip_address").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/interface-routes"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFInterfaceRouteResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFInterfaceRouteRead(d, m)
}

func resourceCudaWAFInterfaceRouteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/interface-routes"
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
		if dataItems["ip-address"] == name {
			break
		}
	}

	if dataItems["ip-address"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("ip_address", name)
	d.Set("interface", dataItems["interface"])
	d.Set("netmask", dataItems["netmask"])
	d.Set("vsite", dataItems["vsite"])
	d.Set("ip_version", dataItems["ip-version"])
	d.Set("comments", dataItems["comments"])

	return nil
}

func resourceCudaWAFInterfaceRouteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/interface-routes"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFInterfaceRouteResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFInterfaceRouteRead(d, m)
}

func resourceCudaWAFInterfaceRouteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/interface-routes"
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

func hydrateBarracudaWAFInterfaceRouteResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"ip-address": d.Get("ip_address").(string),
		"interface":  d.Get("interface").(string),
		"netmask":    d.Get("netmask").(string),
		"vsite":      d.Get("vsite").(string),
		"ip-version": d.Get("ip_version").(string),
		"comments":   d.Get("comments").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"ip-address"}
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
