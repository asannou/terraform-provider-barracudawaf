package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFFormSpam() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFFormSpamCreate,
		Read:   resourceCudaWAFFormSpamRead,
		Update: resourceCudaWAFFormSpamUpdate,
		Delete: resourceCudaWAFFormSpamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid ID specified. Supposed to be <service_name>/<form_name>")
				}
				d.Set("parent", []string{parts[0]})
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name":                    {Type: schema.TypeString, Required: true, Description: "Form Name"},
			"action_url":              {Type: schema.TypeString, Required: true, Description: "Action URL"},
			"minimum_form_fill_time":  {Type: schema.TypeString, Optional: true, Description: "Minimum Form Fill Time"},
			"mode":                    {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"parameter_class":         {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Parameter Class"},
			"parameter_name":          {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Parameter Name"},
			"status":                  {Type: schema.TypeString, Optional: true, Description: "Status"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_form_spam` manages `Form Spam` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFFormSpamCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/form-spam-forms"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFFormSpamResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFFormSpamRead(d, m)
}

func resourceCudaWAFFormSpamRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/form-spam-forms"
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
	d.Set("action_url", dataItems["action-url"])
	d.Set("minimum_form_fill_time", fmt.Sprintf("%v", dataItems["minimum-form-fill-time"]))
	d.Set("mode", dataItems["mode"])
	d.Set("parameter_class", dataItems["parameter-class"])
	d.Set("parameter_name", dataItems["parameter-name"])
	d.Set("status", dataItems["status"])

	return nil
}

func resourceCudaWAFFormSpamUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/form-spam-forms"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFFormSpamResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFFormSpamRead(d, m)
}

func resourceCudaWAFFormSpamDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/form-spam-forms"
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

func hydrateBarracudaWAFFormSpamResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                   d.Get("name").(string),
		"action-url":             d.Get("action_url").(string),
		"minimum-form-fill-time": d.Get("minimum_form_fill_time").(string),
		"mode":                   d.Get("mode").(string),
		"parameter-class":        d.Get("parameter_class"),
		"parameter-name":         d.Get("parameter_name"),
		"status":                 d.Get("status").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"name", "action-url"}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if v, ok := val.(string); ok && len(v) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
