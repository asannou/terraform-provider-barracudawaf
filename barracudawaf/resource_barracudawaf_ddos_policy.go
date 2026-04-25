package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFDDoSPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFDDoSPolicyCreate,
		Read:   resourceCudaWAFDDoSPolicyRead,
		Update: resourceCudaWAFDDoSPolicyUpdate,
		Delete: resourceCudaWAFDDoSPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DDos Policy Name",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host Match",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL Match",
			},
			"extended_match": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extended Match",
			},
			"extended_match_sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extended Match Sequence",
			},
			"evaluate_clients": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable URL ACL",
			},
			"mouse_check": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Detect Mouse Event",
			},
			"enforce_captcha": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enforce CAPTCHA",
			},
			"num_captcha_tries": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max CAPTCHA Attempts",
			},
			"num_unanswered_captcha": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Unanswered CAPTCHA",
			},
			"expiry_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiry time",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comments",
			},
		},

		Description: "`barracudawaf_ddos_policy` manages `DDoS Policy` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFDDoSPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/ddos-policies"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFDDoSPolicyResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFDDoSPolicyRead(d, m)
}

func resourceCudaWAFDDoSPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/ddos-policies"
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
	d.Set("host", dataItems["host"])
	d.Set("url", dataItems["url"])
	d.Set("extended_match", dataItems["extended-match"])
	d.Set("extended_match_sequence", dataItems["extended-match-sequence"])
	d.Set("evaluate_clients", dataItems["evaluate-clients"])
	d.Set("mouse_check", dataItems["mouse-check"])
	d.Set("enforce_captcha", dataItems["enforce-captcha"])
	d.Set("num_captcha_tries", dataItems["num-captcha-tries"])
	d.Set("num_unanswered_captcha", dataItems["num-unanswered-captcha"])
	d.Set("expiry_time", dataItems["expiry-time"])
	d.Set("comments", dataItems["comments"])
	return nil
}

func resourceCudaWAFDDoSPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/ddos-policies"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFDDoSPolicyResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFDDoSPolicyRead(d, m)
}

func resourceCudaWAFDDoSPolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/ddos-policies"
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

func hydrateBarracudaWAFDDoSPolicyResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"host":                    d.Get("host").(string),
		"url":                     d.Get("url").(string),
		"extended-match":          d.Get("extended_match").(string),
		"extended-match-sequence": d.Get("extended_match_sequence").(string),
		"evaluate-clients":        d.Get("evaluate_clients").(string),
		"mouse-check":             d.Get("mouse_check").(string),
		"enforce-captcha":         d.Get("enforce_captcha").(string),
		"num-captcha-tries":       d.Get("num_captcha_tries").(string),
		"num-unanswered-captcha":  d.Get("num_unanswered_captcha").(string),
		"expiry-time":             d.Get("expiry_time").(string),
		"comments":                d.Get("comments").(string),
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
