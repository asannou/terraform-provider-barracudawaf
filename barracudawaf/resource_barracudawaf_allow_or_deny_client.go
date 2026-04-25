package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFAllowOrDenyClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAllowOrDenyClientCreate,
		Read:   resourceCudaWAFAllowOrDenyClientRead,
		Update: resourceCudaWAFAllowOrDenyClientUpdate,
		Delete: resourceCudaWAFAllowOrDenyClientDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule Name",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Description: "Parent format: [service_name]",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status",
			},
			"sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sequence",
			},
			"certificate_serial": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate Serial Number",
			},
			"common_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Common Name",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Organization",
			},
			"organizational_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Organizational Unit",
			},
			"locality": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Locality",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country",
			},
		},

		Description: "`barracudawaf_allow_or_deny_client` manages `Allow or Deny Client` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAllowOrDenyClientCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/allow-deny-clients", serviceName)
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFAllowOrDenyClientResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAllowOrDenyClientRead(d, m)
}

func resourceCudaWAFAllowOrDenyClientRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/allow-deny-clients", serviceName)
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
		log.Printf("[WARN] Barracuda WAF resource (%s) not found on the system", name)
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("action", dataItems["action"])
	d.Set("status", dataItems["status"])
	d.Set("sequence", dataItems["sequence"])
	d.Set("certificate_serial", dataItems["certificate-serial"])
	d.Set("common_name", dataItems["common-name"])
	d.Set("organization", dataItems["organization"])
	d.Set("organizational_unit", dataItems["organizational-unit"])
	d.Set("locality", dataItems["locality"])
	d.Set("state", dataItems["state"])
	d.Set("country", dataItems["country"])
	return nil
}

func resourceCudaWAFAllowOrDenyClientUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/allow-deny-clients", serviceName)
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFAllowOrDenyClientResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAllowOrDenyClientRead(d, m)
}

func resourceCudaWAFAllowOrDenyClientDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/allow-deny-clients", serviceName)
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

func hydrateBarracudaWAFAllowOrDenyClientResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                d.Get("name").(string),
		"action":              d.Get("action").(string),
		"status":              d.Get("status").(string),
		"sequence":            d.Get("sequence").(string),
		"certificate-serial": d.Get("certificate_serial").(string),
		"common-name":        d.Get("common_name").(string),
		"organization":        d.Get("organization").(string),
		"organizational-unit": d.Get("organizational_unit").(string),
		"locality":            d.Get("locality").(string),
		"state":               d.Get("state").(string),
		"country":             d.Get("country").(string),
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
		if reflect.ValueOf(val).Kind() == reflect.String && reflect.ValueOf(val).Len() == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
