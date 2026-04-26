package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFBond() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFBondCreate,
		Read:   resourceCudaWAFBondRead,
		Update: resourceCudaWAFBondUpdate,
		Delete: resourceCudaWAFBondDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name":       {Type: schema.TypeString, Required: true, Description: "Name"},
			"bond_ports": {Type: schema.TypeString, Required: true, Description: "Ports"},
			"min_link":   {Type: schema.TypeString, Optional: true, Description: "Minimum Links"},
			"mode":       {Type: schema.TypeString, Optional: true, Description: "Bonding Mode"},
			"duplexity":  {Type: schema.TypeString, Optional: true, Description: "Duplexity"},
			"mtu":        {Type: schema.TypeString, Optional: true, Description: "MTU"},
			"speed":      {Type: schema.TypeString, Optional: true, Description: "Speed"},
		},

		Description: "`barracudawaf_bond` manages `Bond` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFBondCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFBondResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFBondRead(d, m)
}

func resourceCudaWAFBondRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
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
	d.Set("bond_ports", dataItems["bond-ports"])
	d.Set("min_link", fmt.Sprintf("%v", dataItems["min-link"]))
	d.Set("mode", dataItems["mode"])
	d.Set("duplexity", dataItems["duplexity"])
	d.Set("mtu", fmt.Sprintf("%v", dataItems["mtu"]))
	d.Set("speed", dataItems["speed"])

	return nil
}

func resourceCudaWAFBondUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFBondResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFBondRead(d, m)
}

func resourceCudaWAFBondDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
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

func hydrateBarracudaWAFBondResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":       d.Get("name").(string),
		"bond-ports": d.Get("bond_ports").(string),
		"min-link":   d.Get("min_link").(string),
		"mode":       d.Get("mode").(string),
		"duplexity":  d.Get("duplexity").(string),
		"mtu":        d.Get("mtu").(string),
		"speed":      d.Get("speed").(string),
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
