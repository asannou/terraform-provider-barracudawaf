package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFNetworkACLCreate,
		Read:   resourceCudaWAFNetworkACLRead,
		Update: resourceCudaWAFNetworkACLUpdate,
		Delete: resourceCudaWAFNetworkACLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name":                      {Type: schema.TypeString, Required: true, Description: "Name"},
			"priority":                  {Type: schema.TypeString, Required: true, Description: "Priority"},
			"action":                    {Type: schema.TypeString, Optional: true, Description: "Action"},
			"comments":                  {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"destination_address":       {Type: schema.TypeString, Optional: true, Description: "Destination IP Address"},
			"destination_netmask":       {Type: schema.TypeString, Optional: true, Description: "Destination Netmask"},
			"destination_port":          {Type: schema.TypeString, Optional: true, Description: "Destination Port Range"},
			"enable_logging":            {Type: schema.TypeString, Optional: true, Description: "Log Status"},
			"icmp_response":             {Type: schema.TypeString, Optional: true, Description: "ICMP Response"},
			"interface":                 {Type: schema.TypeString, Optional: true, Description: "Interface"},
			"ip_version":                {Type: schema.TypeString, Optional: true, Description: "IP Protocol Version"},
			"ipv6_destination_address":  {Type: schema.TypeString, Optional: true, Description: "Destination IP Address (IPv6)"},
			"ipv6_destination_netmask":  {Type: schema.TypeString, Optional: true, Description: "Destination Netmask (IPv6)"},
			"ipv6_source_address":       {Type: schema.TypeString, Optional: true, Description: "Source IP Address (IPv6)"},
			"ipv6_source_netmask":       {Type: schema.TypeString, Optional: true, Description: "Source Netmask (IPv6)"},
			"max_connections":           {Type: schema.TypeString, Optional: true, Description: "Max Number of Connections"},
			"max_half_open_connections": {Type: schema.TypeString, Optional: true, Description: "Max Connection Rate"},
			"protocol":                  {Type: schema.TypeString, Optional: true, Description: "Protocol"},
			"source_address":            {Type: schema.TypeString, Optional: true, Description: "Source IP Address"},
			"source_netmask":            {Type: schema.TypeString, Optional: true, Description: "Source Netmask"},
			"source_port":               {Type: schema.TypeString, Optional: true, Description: "Source Port Range"},
			"status":                    {Type: schema.TypeString, Optional: true, Description: "Enabled"},
			"vsite":                     {Type: schema.TypeString, Optional: true, Description: "Network Group"},
		},

		Description: "`barracudawaf_network_acl` manages `Network ACL` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFNetworkACLCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/network-acls"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFNetworkACLResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFNetworkACLRead(d, m)
}

func resourceCudaWAFNetworkACLRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/network-acls"
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
	d.Set("priority", fmt.Sprintf("%v", dataItems["priority"]))
	d.Set("action", dataItems["action"])
	d.Set("comments", dataItems["comments"])
	d.Set("destination_address", dataItems["destination-address"])
	d.Set("destination_netmask", dataItems["destination-netmask"])
	d.Set("destination_port", dataItems["destination-port"])
	d.Set("enable_logging", dataItems["enable-logging"])
	d.Set("icmp_response", dataItems["icmp-response"])
	d.Set("interface", dataItems["interface"])
	d.Set("ip_version", dataItems["ip-version"])
	d.Set("ipv6_destination_address", dataItems["ipv6-destination-address"])
	d.Set("ipv6_destination_netmask", dataItems["ipv6-destination-netmask"])
	d.Set("ipv6_source_address", dataItems["ipv6-source-address"])
	d.Set("ipv6_source_netmask", dataItems["ipv6-source-netmask"])
	d.Set("max_connections", fmt.Sprintf("%v", dataItems["max-connections"]))
	d.Set("max_half_open_connections", fmt.Sprintf("%v", dataItems["max-half-open-connections"]))
	d.Set("protocol", dataItems["protocol"])
	d.Set("source_address", dataItems["source-address"])
	d.Set("source_netmask", dataItems["source-netmask"])
	d.Set("source_port", dataItems["source-port"])
	d.Set("status", dataItems["status"])
	d.Set("vsite", dataItems["vsite"])

	return nil
}

func resourceCudaWAFNetworkACLUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/network-acls"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFNetworkACLResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFNetworkACLRead(d, m)
}

func resourceCudaWAFNetworkACLDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/network-acls"
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

func hydrateBarracudaWAFNetworkACLResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                      d.Get("name").(string),
		"priority":                  d.Get("priority").(string),
		"action":                    d.Get("action").(string),
		"comments":                  d.Get("comments").(string),
		"destination-address":       d.Get("destination_address").(string),
		"destination-netmask":       d.Get("destination_netmask").(string),
		"destination-port":          d.Get("destination_port").(string),
		"enable-logging":            d.Get("enable_logging").(string),
		"icmp-response":             d.Get("icmp_response").(string),
		"interface":                 d.Get("interface").(string),
		"ip-version":                d.Get("ip_version").(string),
		"ipv6-destination-address":  d.Get("ipv6_destination_address").(string),
		"ipv6-destination-netmask":  d.Get("ipv6_destination_netmask").(string),
		"ipv6-source-address":       d.Get("ipv6_source_address").(string),
		"ipv6-source-netmask":       d.Get("ipv6_source_netmask").(string),
		"max-connections":           d.Get("max_connections").(string),
		"max-half-open-connections": d.Get("max_half_open_connections").(string),
		"protocol":                  d.Get("protocol").(string),
		"source-address":            d.Get("source_address").(string),
		"source-netmask":            d.Get("source_netmask").(string),
		"source-port":               d.Get("source_port").(string),
		"status":                    d.Get("status").(string),
		"vsite":                     d.Get("vsite").(string),
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
