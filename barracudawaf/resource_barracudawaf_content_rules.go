package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceContentRulesParams = map[string][]string{
		"caching": {
			"cache_negative_responses",
			"expiry_age",
			"file_extensions",
			"ignore_request_headers",
			"ignore_response_headers",
			"max_size",
			"min_size",
			"status",
		},
		"compression": {
			"compress_unknown_content_types",
			"content_types",
			"min_size",
			"status",
		},
		"load_balancing": {
			"cookie_age",
			"failover_method",
			"header_name",
			"lb_algorithm",
			"parameter_name",
			"persistence_cookie_domain",
			"persistence_cookie_name",
			"persistence_cookie_path",
			"persistence_idle_timeout",
			"persistence_method",
			"source_ip_netmask",
		},
		"captcha_settings": {
			"recaptcha_type",
			"rg_recaptcha_domain",
			"rg_recaptcha_site_key",
			"rg_recaptcha_site_secret",
		},
		"advanced_client_analysis": {
			"advanced_analysis",
		},
	}
)

func resourceCudaWAFContentRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFContentRulesCreate,
		Read:   resourceCudaWAFContentRulesRead,
		Update: resourceCudaWAFContentRulesUpdate,
		Delete: resourceCudaWAFContentRulesDelete,
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
			"access_log":              {Type: schema.TypeString, Optional: true, Description: "Access Log"},
			"app_id":                  {Type: schema.TypeString, Optional: true, Description: "Rule App Id"},
			"comments":                {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"host_match":              {Type: schema.TypeString, Required: true, Description: "Host Match"},
			"name":                    {Type: schema.TypeString, Required: true, Description: "Rule Group Name"},
			"status":                  {Type: schema.TypeString, Optional: true, Description: "Status"},
			"extended_match":          {Type: schema.TypeString, Optional: true, Description: "Extended Match"},
			"extended_match_sequence": {Type: schema.TypeString, Optional: true, Description: "Extended Match Sequence"},
			"mode":                    {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"url_match":               {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"web_firewall_policy":     {Type: schema.TypeString, Optional: true, Description: "Web Firewall Policy"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"caching": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_negative_responses": {Type: schema.TypeString, Optional: true},
						"expiry_age":               {Type: schema.TypeString, Optional: true},
						"file_extensions":          {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"ignore_request_headers":   {Type: schema.TypeString, Optional: true},
						"ignore_response_headers":  {Type: schema.TypeString, Optional: true},
						"max_size":                 {Type: schema.TypeString, Optional: true},
						"min_size":                 {Type: schema.TypeString, Optional: true},
						"status":                   {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"compression": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compress_unknown_content_types": {Type: schema.TypeString, Optional: true},
						"content_types":                  {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"min_size":                       {Type: schema.TypeString, Optional: true},
						"status":                         {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"load_balancing": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_age":                {Type: schema.TypeString, Optional: true},
						"failover_method":           {Type: schema.TypeString, Optional: true},
						"header_name":               {Type: schema.TypeString, Optional: true},
						"lb_algorithm":              {Type: schema.TypeString, Optional: true},
						"parameter_name":            {Type: schema.TypeString, Optional: true},
						"persistence_cookie_domain": {Type: schema.TypeString, Optional: true},
						"persistence_cookie_name":   {Type: schema.TypeString, Optional: true},
						"persistence_cookie_path":   {Type: schema.TypeString, Optional: true},
						"persistence_idle_timeout":  {Type: schema.TypeString, Optional: true},
						"persistence_method":        {Type: schema.TypeString, Optional: true},
						"source_ip_netmask":         {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"captcha_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recaptcha_type":           {Type: schema.TypeString, Optional: true},
						"rg_recaptcha_domain":      {Type: schema.TypeString, Optional: true},
						"rg_recaptcha_site_key":    {Type: schema.TypeString, Optional: true},
						"rg_recaptcha_site_secret": {Type: schema.TypeString, Optional: true, Sensitive: true},
					},
				},
			},
			"advanced_client_analysis": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_analysis": {Type: schema.TypeString, Optional: true},
					},
				},
			},
		},

		Description: "`barracudawaf_content_rules` manages `Content Rules` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFContentRulesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFContentRulesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFContentRulesRead(d, m)
}

func resourceCudaWAFContentRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules"
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
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("name", name)

	payload := map[string]string{
		"access_log":              "access-log",
		"app_id":                  "app-id",
		"comments":                "comments",
		"host_match":              "host-match",
		"status":                  "status",
		"extended_match":          "extended-match",
		"extended_match_sequence": "extended-match-sequence",
		"mode":                    "mode",
		"url_match":               "url-match",
		"web_firewall_policy":     "web-firewall-policy",
	}

	for tfKey, apiKey := range payload {
		if val, ok := dataItems[apiKey]; ok && val != nil {
			d.Set(tfKey, fmt.Sprintf("%v", val))
		}
	}

	// Read sub-resources
	for subResource, subResourceParams := range subResourceContentRulesParams {
		subResourceURL := strings.Replace(subResource, "_", "-", -1)
		if subResource == "captcha_settings" {
			subResourceURL = "rg-captcha-settings"
		}

		subResourceEndpoint := fmt.Sprintf("%s/%s/%s", resourceEndpoint, name, subResourceURL)
		subRequest := &APIRequest{
			Method: "get",
			URL:    subResourceEndpoint,
		}

		subResources, err := client.GetBarracudaWAFResource(name, subRequest)
		if err != nil {
			log.Printf("[ERROR] Unable to Retrieve Barracuda WAF sub-resource (%s) (%v) ", subResource, err)
			continue
		}

		if subResources.Data == nil {
			continue
		}

		var subResourceList []interface{}
		for _, subDataItems := range subResources.Data {
			subMap := make(map[string]interface{})
			for _, param := range subResourceParams {
				apiParam := strings.Replace(param, "_", "-", -1)
				if val, ok := subDataItems[apiParam]; ok && val != nil {
					if reflect.TypeOf(val).Kind() == reflect.Slice {
						subMap[param] = val
					} else {
						subMap[param] = fmt.Sprintf("%v", val)
					}
				}
			}
			subResourceList = append(subResourceList, subMap)
		}

		if len(subResourceList) > 0 {
			d.Set(subResource, subResourceList)
		}
	}

	return nil
}

func resourceCudaWAFContentRulesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFContentRulesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFContentRulesRead(d, m)
}

func resourceCudaWAFContentRulesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules"
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

func hydrateBarracudaWAFContentRulesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"access-log":              d.Get("access_log").(string),
		"app-id":                  d.Get("app_id").(string),
		"comments":                d.Get("comments").(string),
		"host-match":              d.Get("host_match").(string),
		"name":                    d.Get("name").(string),
		"status":                  d.Get("status").(string),
		"extended-match":          d.Get("extended_match").(string),
		"extended-match-sequence": d.Get("extended_match_sequence").(string),
		"mode":                    d.Get("mode").(string),
		"url-match":               d.Get("url_match").(string),
		"web-firewall-policy":     d.Get("web_firewall_policy").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{}
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

func (b *BarracudaWAF) hydrateBarracudaWAFContentRulesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceContentRulesParams {
		subResourceParamsLength := d.Get(subResource + ".#").(int)

		log.Printf("[INFO] Updating Barracuda WAF sub resource (%s) (%s)", name, subResource)

		for i := 0; i < subResourceParamsLength; i++ {
			subResourcePayload := make(map[string]interface{})
			suffix := fmt.Sprintf(".%d", i)

			for _, param := range subResourceParams {
				paramSuffix := fmt.Sprintf(".%s", param)
				paramVaule := d.Get(subResource + suffix + paramSuffix)

				if reflect.ValueOf(paramVaule).Len() > 0 {
					param = strings.Replace(param, "_", "-", -1)
					subResourcePayload[param] = paramVaule
				}
			}

			// Special handling for captcha_settings endpoint
			url := strings.Replace(subResource, "_", "-", -1)
			if subResource == "captcha_settings" {
				url = "rg-captcha-settings"
			}

			err := b.UpdateBarracudaWAFSubResource(name, endpoint, &APIRequest{
				URL:  url,
				Body: subResourcePayload,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
