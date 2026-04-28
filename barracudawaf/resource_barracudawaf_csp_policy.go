package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFCSPPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFCSPPolicyCreate,
		Read:   resourceCudaWAFCSPPolicyRead,
		Update: resourceCudaWAFCSPPolicyUpdate,
		Delete: resourceCudaWAFCSPPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"csp_policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CSP Policy Name",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Description: "Parent format: [service_name, rule_name]",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status",
			},
			"csp_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mode",
			},
			"default_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Default Source",
			},
			"script_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Script Source",
			},
			"style_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Style Source",
			},
			"img_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Image Source",
			},
			"connect_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Connect Source",
			},
			"font_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Font Source",
			},
			"object_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Object Source",
			},
			"media_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Media Source",
			},
			"frame_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Frame Source",
			},
			"child_src": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Child Source",
			},
			"form_action": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Forms Action",
			},
			"frames_ancestors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Frames Ancestors",
			},
			"base_uri": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Base URI",
			},
			"report_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Report URI",
			},
			"upgrade_unsecure_req": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Upgrade Insecure Requests",
			},
			"block_all_mixed_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block All Mixed Content",
			},
		},

		Description: "`barracudawaf_csp_policy` manages `CSP Policy` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFCSPPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("csp_policy_name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	ruleName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/client-side-protection/%s/csp-policies", serviceName, ruleName)
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFCSPPolicyResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFCSPPolicyRead(d, m)
}

func resourceCudaWAFCSPPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	ruleName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/client-side-protection/%s/csp-policies", serviceName, ruleName)
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
		if dataItems["csp-policy-name"] == name {
			break
		}
	}

	if dataItems["csp-policy-name"] != name {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found on the system", name)
		d.SetId("")
		return nil
	}

	d.Set("csp_policy_name", name)
	d.Set("status", dataItems["status"])
	d.Set("csp_mode", dataItems["csp-mode"])
	if val, ok := dataItems["default-src"]; ok && val != nil {
		d.Set("default_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("default_src", nil)
	}
	if val, ok := dataItems["script-src"]; ok && val != nil {
		d.Set("script_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("script_src", nil)
	}
	if val, ok := dataItems["style-src"]; ok && val != nil {
		d.Set("style_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("style_src", nil)
	}
	if val, ok := dataItems["img-src"]; ok && val != nil {
		d.Set("img_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("img_src", nil)
	}
	if val, ok := dataItems["connect-src"]; ok && val != nil {
		d.Set("connect_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("connect_src", nil)
	}
	if val, ok := dataItems["font-src"]; ok && val != nil {
		d.Set("font_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("font_src", nil)
	}
	if val, ok := dataItems["object-src"]; ok && val != nil {
		d.Set("object_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("object_src", nil)
	}
	if val, ok := dataItems["media-src"]; ok && val != nil {
		d.Set("media_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("media_src", nil)
	}
	if val, ok := dataItems["frame-src"]; ok && val != nil {
		d.Set("frame_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("frame_src", nil)
	}
	if val, ok := dataItems["child-src"]; ok && val != nil {
		d.Set("child_src", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("child_src", nil)
	}
	if val, ok := dataItems["form-action"]; ok && val != nil {
		d.Set("form_action", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("form_action", nil)
	}
	if val, ok := dataItems["frames-ancestors"]; ok && val != nil {
		d.Set("frames_ancestors", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("frames_ancestors", nil)
	}
	if val, ok := dataItems["base-uri"]; ok && val != nil {
		d.Set("base_uri", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("base_uri", nil)
	}
	d.Set("report_uri", dataItems["report-uri"])
	d.Set("upgrade_unsecure_req", dataItems["upgrade-unsecure-req"])
	d.Set("block_all_mixed_content", dataItems["mixed-box-content"])
	return nil
}

func resourceCudaWAFCSPPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	ruleName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/client-side-protection/%s/csp-policies", serviceName, ruleName)
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFCSPPolicyResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFCSPPolicyRead(d, m)
}

func resourceCudaWAFCSPPolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	ruleName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/client-side-protection/%s/csp-policies", serviceName, ruleName)
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

func hydrateBarracudaWAFCSPPolicyResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"csp-policy-name":         d.Get("csp_policy_name").(string),
		"status":                  d.Get("status").(string),
		"csp-mode":                d.Get("csp_mode").(string),
		"default-src":             d.Get("default_src"),
		"script-src":              d.Get("script_src"),
		"style-src":               d.Get("style_src"),
		"img-src":                 d.Get("img_src"),
		"connect-src":             d.Get("connect_src"),
		"font-src":                d.Get("font_src"),
		"object-src":              d.Get("object_src"),
		"media-src":               d.Get("media_src"),
		"frame-src":               d.Get("frame_src"),
		"child-src":               d.Get("child_src"),
		"form-action":             d.Get("form_action"),
		"frames-ancestors":        d.Get("frames_ancestors"),
		"base-uri":                d.Get("base_uri"),
		"report-uri":              d.Get("report_uri").(string),
		"upgrade-unsecure-req":    d.Get("upgrade_unsecure_req").(string),
		"mixed-box-content":       d.Get("block_all_mixed_content").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"csp-policy-name"}
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
