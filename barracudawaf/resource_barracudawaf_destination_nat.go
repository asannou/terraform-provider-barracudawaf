package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFDestinationNAT() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFDestinationNATCreate,
		Read:   resourceCudaWAFDestinationNATRead,
		Update: resourceCudaWAFDestinationNATUpdate,
		Delete: resourceCudaWAFDestinationNATDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"pre_destination_address":  {Type: schema.TypeString, Required: true, Description: "Pre-DNAT Destination"},
			"incoming_interface":       {Type: schema.TypeString, Required: true, Description: "Incoming Interface"},
			"post_destination_address": {Type: schema.TypeString, Required: true, Description: "Post-DNAT Destination"},
			"pre_destination_netmask":  {Type: schema.TypeString, Required: true, Description: "Pre-DNAT Destination Mask"},
			"protocol":                 {Type: schema.TypeString, Required: true, Description: "Protocol"},
			"vsite":                    {Type: schema.TypeString, Required: true, Description: "Network Group"},
			"pre_destination_port":     {Type: schema.TypeString, Optional: true, Description: "Destination Port"},
			"comments":                 {Type: schema.TypeString, Optional: true, Description: "Comments"},
		},

		Description: "`barracudawaf_destination_nat` manages `Destination NAT` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFDestinationNATCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("pre_destination_address").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/destination-nats"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFDestinationNATResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFDestinationNATRead(d, m)
}

func resourceCudaWAFDestinationNATRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/destination-nats"
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
		if dataItems["pre-destination-address"] == name {
			break
		}
	}

	if dataItems["pre-destination-address"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("pre_destination_address", name)
	d.Set("incoming_interface", dataItems["incoming-interface"])
	d.Set("post_destination_address", dataItems["post-destination-address"])
	d.Set("pre_destination_netmask", dataItems["pre-destination-netmask"])
	d.Set("protocol", dataItems["protocol"])
	d.Set("vsite", dataItems["vsite"])
	d.Set("pre_destination_port", fmt.Sprintf("%v", dataItems["pre-destination-port"]))
	d.Set("comments", dataItems["comments"])

	return nil
}

func resourceCudaWAFDestinationNATUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/destination-nats"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFDestinationNATResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFDestinationNATRead(d, m)
}

func resourceCudaWAFDestinationNATDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/destination-nats"
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

func hydrateBarracudaWAFDestinationNATResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"pre-destination-address":  d.Get("pre_destination_address").(string),
		"incoming-interface":       d.Get("incoming_interface").(string),
		"post-destination-address": d.Get("post_destination_address").(string),
		"pre-destination-netmask":  d.Get("pre_destination_netmask").(string),
		"protocol":                 d.Get("protocol").(string),
		"vsite":                    d.Get("vsite").(string),
		"pre-destination-port":     d.Get("pre_destination_port").(string),
		"comments":                 d.Get("comments").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"pre-destination-address"}
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
