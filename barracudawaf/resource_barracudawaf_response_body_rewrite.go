package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFResponseBodyRewrite() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFResponseBodyRewriteCreate,
		Read:   resourceCudaWAFResponseBodyRewriteRead,
		Update: resourceCudaWAFResponseBodyRewriteUpdate,
		Delete: resourceCudaWAFResponseBodyRewriteDelete,
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
			"name":            {Type: schema.TypeString, Required: true, Description: "Rule Name"},
			"comments":        {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"host":            {Type: schema.TypeString, Required: true, Description: "Host Match"},
			"replace_string":  {Type: schema.TypeString, Optional: true, Description: "Replace String"},
			"search_string":   {Type: schema.TypeString, Required: true, Description: "Search String"},
			"sequence_number": {Type: schema.TypeString, Required: true, Description: "Sequence number"},
			"url":             {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_response_body_rewrite` manages `Response Body Rewrite` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFResponseBodyRewriteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/response-body-rewrite-rules"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFResponseBodyRewriteResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFResponseBodyRewriteRead(d, m)
}

func resourceCudaWAFResponseBodyRewriteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/response-body-rewrite-rules"
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
	d.Set("comments", dataItems["comments"])
	d.Set("host", dataItems["host"])
	d.Set("replace_string", dataItems["replace-string"])
	d.Set("search_string", dataItems["search-string"])
	d.Set("sequence_number", fmt.Sprintf("%v", dataItems["sequence-number"]))
	d.Set("url", dataItems["url"])

	return nil
}

func resourceCudaWAFResponseBodyRewriteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/response-body-rewrite-rules"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFResponseBodyRewriteResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFResponseBodyRewriteRead(d, m)
}

func resourceCudaWAFResponseBodyRewriteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/response-body-rewrite-rules"
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

func hydrateBarracudaWAFResponseBodyRewriteResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":            d.Get("name").(string),
		"comments":        d.Get("comments").(string),
		"host":            d.Get("host").(string),
		"replace-string":  d.Get("replace_string").(string),
		"search-string":   d.Get("search_string").(string),
		"sequence-number": d.Get("sequence_number").(string),
		"url":             d.Get("url").(string),
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
