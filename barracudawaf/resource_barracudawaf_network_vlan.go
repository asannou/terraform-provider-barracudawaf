package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFNetworkVLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFNetworkVLANCreate,
		Read:   resourceCudaWAFNetworkVLANRead,
		Update: resourceCudaWAFNetworkVLANUpdate,
		Delete: resourceCudaWAFNetworkVLANDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name":      {Type: schema.TypeString, Required: true, Description: "VLAN Name"},
			"interface": {Type: schema.TypeString, Required: true, Description: "VLAN Interface"},
			"vlan_id":   {Type: schema.TypeString, Required: true, Description: "VLAN ID"},
			"vsite":     {Type: schema.TypeString, Required: true, Description: "Network Group"},
			"comments":  {Type: schema.TypeString, Optional: true, Description: "Comments"},
		},

		Description: "`barracudawaf_network_vlan` manages `Network VLAN` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFNetworkVLANCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/vlans"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFNetworkVLANResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFNetworkVLANRead(d, m)
}

func resourceCudaWAFNetworkVLANRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/vlans"
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
	d.Set("interface", dataItems["interface"])
	d.Set("vlan_id", fmt.Sprintf("%v", dataItems["vlan-id"]))
	d.Set("vsite", dataItems["vsite"])
	d.Set("comments", dataItems["comments"])

	return nil
}

func resourceCudaWAFNetworkVLANUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/vlans"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFNetworkVLANResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFNetworkVLANRead(d, m)
}

func resourceCudaWAFNetworkVLANDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/vlans"
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

func hydrateBarracudaWAFNetworkVLANResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":      d.Get("name").(string),
		"interface": d.Get("interface").(string),
		"vlan-id":   d.Get("vlan_id").(string),
		"vsite":     d.Get("vsite").(string),
		"comments":  d.Get("comments").(string),
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
