package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFWebScrapingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFWebScrapingPolicyCreate,
		Read:   resourceCudaWAFWebScrapingPolicyRead,
		Update: resourceCudaWAFWebScrapingPolicyUpdate,
		Delete: resourceCudaWAFWebScrapingPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Web Scraping Policy Name",
			},
			"detect_mouse_event": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Detect Mouse Event",
			},
			"insert_delay": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Insert Delay in Robots.txt",
			},
			"whitelisted_bots": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Allow-listed Bots",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment",
			},
			"blacklisted_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Blocklisted Categories",
			},
			"delay_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Delay Time",
			},
			"insert_disallowed_urls": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Insert Disallowed URLs in Robots.txt",
			},
			"insert_hidden_links": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Insert Hidden Links in Response",
			},
			"insert_javascript_in_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Insert JavaScript in Response",
			},
		},

		Description: "`barracudawaf_web_scraping_policy` manages `Web Scraping Policy` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFWebScrapingPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/web-scraping-policies"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFWebScrapingPolicyResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFWebScrapingPolicyRead(d, m)
}

func resourceCudaWAFWebScrapingPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/web-scraping-policies"
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

	found := false
	for _, dataItems = range resources.Data {
		if dataItems["name"] == name {
			found = true
			break
		}
	}

	if !found {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found on the system", name)
		d.SetId("")
		return nil
	}

	d.Set("name", name)

	payload := map[string]string{
		"detect_mouse_event":            "detect-mouse-event",
		"insert_delay":                  "insert-delay",
		"comments":                      "comments",
		"delay_time":                    "delay-time",
		"insert_disallowed_urls":        "insert-disallowed-urls",
		"insert_hidden_links":           "insert-hidden-links",
		"insert_javascript_in_response": "insert-javascript-in-response",
	}

	for tfKey, apiKey := range payload {
		if val, ok := dataItems[apiKey]; ok && val != nil {
			d.Set(tfKey, fmt.Sprintf("%v", val))
		} else {
			d.Set(tfKey, nil)
		}
	}

	if val, ok := dataItems["whitelisted-bots"]; ok && val != nil {
		if sliceVal, ok := val.([]interface{}); ok {
			d.Set("whitelisted_bots", sortFileList(sliceVal, ""))
		} else {
			d.Set("whitelisted_bots", nil)
		}
	} else {
		d.Set("whitelisted_bots", nil)
	}

	if val, ok := dataItems["blacklisted-categories"]; ok && val != nil {
		if sliceVal, ok := val.([]interface{}); ok {
			d.Set("blacklisted_categories", sortFileList(sliceVal, ""))
		} else {
			d.Set("blacklisted_categories", nil)
		}
	} else {
		d.Set("blacklisted_categories", nil)
	}

	return nil
}

func resourceCudaWAFWebScrapingPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/web-scraping-policies"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFWebScrapingPolicyResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFWebScrapingPolicyRead(d, m)
}

func resourceCudaWAFWebScrapingPolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/web-scraping-policies"
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

func hydrateBarracudaWAFWebScrapingPolicyResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	resourcePayload := map[string]interface{}{
		"name":                          d.Get("name").(string),
		"detect-mouse-event":            d.Get("detect_mouse_event").(string),
		"insert-delay":                  d.Get("insert_delay").(string),
		"whitelisted-bots":              d.Get("whitelisted_bots").([]interface{}),
		"comments":                      d.Get("comments").(string),
		"blacklisted-categories":        d.Get("blacklisted_categories").([]interface{}),
		"delay-time":                    d.Get("delay_time").(string),
		"insert-disallowed-urls":        d.Get("insert_disallowed_urls").(string),
		"insert-hidden-links":           d.Get("insert_hidden_links").(string),
		"insert-javascript-in-response": d.Get("insert_javascript_in_response").(string),
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
		if reflect.ValueOf(val).Kind() == reflect.Slice && reflect.ValueOf(val).Len() == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
