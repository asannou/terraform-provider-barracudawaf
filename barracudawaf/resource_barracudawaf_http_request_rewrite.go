package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFHTTPRequestRewrite() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFHTTPRequestRewriteCreate,
		Read:   resourceCudaWAFHTTPRequestRewriteRead,
		Update: resourceCudaWAFHTTPRequestRewriteUpdate,
		Delete: resourceCudaWAFHTTPRequestRewriteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid ID specified. Supposed to be <service_name>/<rule_name>")
				}
				d.Set("parent", []string{parts[0]})
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name":                {Type: schema.TypeString, Required: true, Description: "Rule Name"},
			"action":              {Type: schema.TypeString, Optional: true, Description: "Action"},
			"comments":            {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"condition":           {Type: schema.TypeString, Optional: true, Description: "Rewrite Condition"},
			"continue_processing": {Type: schema.TypeString, Optional: true, Description: "Continue Processing"},
			"header":              {Type: schema.TypeString, Optional: true, Description: "Header Name"},
			"old_value":           {Type: schema.TypeString, Optional: true, Description: "Old Value"},
			"rewrite_value":       {Type: schema.TypeString, Optional: true, Description: "Rewrite Value"},
			"sequence_number":     {Type: schema.TypeString, Required: true, Description: "Sequence Number"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_http_request_rewrite` manages `HTTP Request Rewrite` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFHTTPRequestRewriteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/http-request-rewrite-rules"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFHTTPRequestRewriteResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFHTTPRequestRewriteRead(d, m)
}

func resourceCudaWAFHTTPRequestRewriteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/http-request-rewrite-rules"
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
	d.Set("action", dataItems["action"])
	d.Set("comments", dataItems["comments"])
	d.Set("condition", dataItems["condition"])
	d.Set("continue_processing", dataItems["continue-processing"])
	d.Set("header", dataItems["header"])
	d.Set("old_value", dataItems["old-value"])
	d.Set("rewrite_value", dataItems["rewrite-value"])
	d.Set("sequence_number", fmt.Sprintf("%v", dataItems["sequence-number"]))

	return nil
}

func resourceCudaWAFHTTPRequestRewriteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/http-request-rewrite-rules"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFHTTPRequestRewriteResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFHTTPRequestRewriteRead(d, m)
}

func resourceCudaWAFHTTPRequestRewriteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/http-request-rewrite-rules"
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

func hydrateBarracudaWAFHTTPRequestRewriteResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                d.Get("name").(string),
		"action":              d.Get("action").(string),
		"comments":            d.Get("comments").(string),
		"condition":           d.Get("condition").(string),
		"continue-processing": d.Get("continue_processing").(string),
		"header":              d.Get("header").(string),
		"old-value":           d.Get("old_value").(string),
		"rewrite-value":       d.Get("rewrite_value").(string),
		"sequence-number":     d.Get("sequence_number").(string),
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
