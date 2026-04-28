package barracudawaf

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceSecurityPoliciesParams = map[string][]string{
		"cookie_security": {
			"allow_unrecognized_cookies", "cookie_max_age", "cookie_replay_protection_type", "cookies_exempted", "custom_headers", "days_allowed", "http_only", "same_site", "secure_cookie", "tamper_proof_mode",
		},
		"url_protection": {
			"allowed_content_types", "allowed_methods", "blocked_attack_types", "csrf_prevention", "custom_blocked_attack_types", "enable", "exception_patterns", "max_content_length", "max_parameters", "maximum_parameter_name_length", "maximum_upload_files",
		},
		"parameter_protection": {
			"allowed_file_upload_type", "base64_decode_parameter_value", "blocked_attack_types", "custom_blocked_attack_types", "denied_metacharacters", "enable", "exception_patterns", "file_upload_extensions", "file_upload_mime_types", "ignore_parameters", "maximum_instances", "maximum_parameter_value_length", "maximum_upload_file_size", "validate_parameter_name",
		},
		"cloaking": {
			"filter_response_header", "headers_to_filter", "return_codes_to_exempt", "suppress_return_code",
		},
		"url_normalization": {
			"apply_double_decoding", "default_charset", "detect_response_charset", "normalize_special_chars", "parameter_separators",
		},
		"request_limits": {
			"enable", "max_cookie_name_length", "max_cookie_value_length", "max_header_name_length", "max_header_value_length", "max_number_of_cookies", "max_number_of_headers", "max_query_length", "max_request_length", "max_request_line_length", "max_url_length",
		},
		"client_profile": {
			"client_profile", "exception_client_fingerprints", "high_risk_score", "medium_risk_score",
		},
		"tarpit_profile": {
			"backlog_requests_limit", "tarpit_delay_interval", "tarpit_inactivity_timeout",
		},
		"action_policies": {
			"action", "deny_response", "follow_up_action", "follow_up_action_time", "name", "redirect_url", "response_page", "risk_score",
		},
		"global_acls": {
			"action", "comments", "deny_response", "enable", "extended_match", "extended_match_sequence", "follow_up_action", "follow_up_action_time", "name", "redirect_url", "response_page", "url",
		},
		"protected_data_types": {
			"action", "custom_identity_theft_type", "enable", "identity_theft_type", "initial_characters_to_keep", "name", "trailing_characters_to_keep",
		},
	}
)

func resourceCudaWAFSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSecurityPoliciesCreate,
		Read:   resourceCudaWAFSecurityPoliciesRead,
		Update: resourceCudaWAFSecurityPoliciesUpdate,
		Delete: resourceCudaWAFSecurityPoliciesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"based_on": {Type: schema.TypeString, Optional: true},
			"name":     {Type: schema.TypeString, Required: true, Description: "Policy Name"},
			"cookie_security": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_unrecognized_cookies":    {Type: schema.TypeString, Optional: true},
						"cookie_max_age":                {Type: schema.TypeString, Optional: true},
						"cookie_replay_protection_type": {Type: schema.TypeString, Optional: true},
						"cookies_exempted":              {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"custom_headers":                {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"days_allowed":                  {Type: schema.TypeString, Optional: true},
						"http_only":                     {Type: schema.TypeString, Optional: true},
						"same_site":                     {Type: schema.TypeString, Optional: true},
						"secure_cookie":                 {Type: schema.TypeString, Optional: true},
						"tamper_proof_mode":             {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"url_protection": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_content_types":         {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"allowed_methods":               {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"blocked_attack_types":          {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"csrf_prevention":               {Type: schema.TypeString, Optional: true},
						"custom_blocked_attack_types":   {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"enable":                        {Type: schema.TypeString, Optional: true},
						"exception_patterns":            {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"max_content_length":            {Type: schema.TypeString, Optional: true},
						"max_parameters":                {Type: schema.TypeString, Optional: true},
						"maximum_parameter_name_length": {Type: schema.TypeString, Optional: true},
						"maximum_upload_files":          {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"parameter_protection": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_file_upload_type":       {Type: schema.TypeString, Optional: true},
						"base64_decode_parameter_value":  {Type: schema.TypeString, Optional: true},
						"blocked_attack_types":           {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"custom_blocked_attack_types":    {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"denied_metacharacters":          {Type: schema.TypeString, Optional: true},
						"enable":                         {Type: schema.TypeString, Optional: true},
						"exception_patterns":             {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"file_upload_extensions":         {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"file_upload_mime_types":         {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"ignore_parameters":              {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"maximum_instances":              {Type: schema.TypeString, Optional: true},
						"maximum_parameter_value_length": {Type: schema.TypeString, Optional: true},
						"maximum_upload_file_size":       {Type: schema.TypeString, Optional: true},
						"validate_parameter_name":        {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"cloaking": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_response_header": {Type: schema.TypeString, Optional: true},
						"headers_to_filter":      {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"return_codes_to_exempt": {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"suppress_return_code":   {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"url_normalization": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apply_double_decoding":  {Type: schema.TypeString, Optional: true},
						"default_charset":        {Type: schema.TypeString, Optional: true},
						"detect_response_charset": {Type: schema.TypeString, Optional: true},
						"normalize_special_chars": {Type: schema.TypeString, Optional: true},
						"parameter_separators":   {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"request_limits": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable":                  {Type: schema.TypeString, Optional: true},
						"max_cookie_name_length":  {Type: schema.TypeString, Optional: true},
						"max_cookie_value_length": {Type: schema.TypeString, Optional: true},
						"max_header_name_length":  {Type: schema.TypeString, Optional: true},
						"max_header_value_length": {Type: schema.TypeString, Optional: true},
						"max_number_of_cookies":   {Type: schema.TypeString, Optional: true},
						"max_number_of_headers":   {Type: schema.TypeString, Optional: true},
						"max_query_length":        {Type: schema.TypeString, Optional: true},
						"max_request_length":      {Type: schema.TypeString, Optional: true},
						"max_request_line_length": {Type: schema.TypeString, Optional: true},
						"max_url_length":          {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"client_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_profile":               {Type: schema.TypeString, Optional: true},
						"exception_client_fingerprints": {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"high_risk_score":              {Type: schema.TypeString, Optional: true},
						"medium_risk_score":            {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"tarpit_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backlog_requests_limit":   {Type: schema.TypeString, Optional: true},
						"tarpit_delay_interval":    {Type: schema.TypeString, Optional: true},
						"tarpit_inactivity_timeout": {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"action_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action":                {Type: schema.TypeString, Optional: true},
						"deny_response":         {Type: schema.TypeString, Optional: true},
						"follow_up_action":      {Type: schema.TypeString, Optional: true},
						"follow_up_action_time": {Type: schema.TypeString, Optional: true},
						"name":                  {Type: schema.TypeString, Required: true},
						"redirect_url":          {Type: schema.TypeString, Optional: true},
						"response_page":         {Type: schema.TypeString, Optional: true},
						"risk_score":            {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"global_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action":                {Type: schema.TypeString, Optional: true},
						"comments":              {Type: schema.TypeString, Optional: true},
						"deny_response":         {Type: schema.TypeString, Optional: true},
						"enable":                {Type: schema.TypeString, Optional: true},
						"extended_match":        {Type: schema.TypeString, Optional: true},
						"extended_match_sequence": {Type: schema.TypeString, Optional: true},
						"follow_up_action":      {Type: schema.TypeString, Optional: true},
						"follow_up_action_time": {Type: schema.TypeString, Optional: true},
						"name":                  {Type: schema.TypeString, Required: true},
						"redirect_url":          {Type: schema.TypeString, Optional: true},
						"response_page":         {Type: schema.TypeString, Optional: true},
						"url":                   {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"protected_data_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action":                     {Type: schema.TypeString, Optional: true},
						"custom_identity_theft_type": {Type: schema.TypeString, Optional: true},
						"enable":                     {Type: schema.TypeString, Optional: true},
						"identity_theft_type":        {Type: schema.TypeString, Optional: true},
						"initial_characters_to_keep": {Type: schema.TypeString, Optional: true},
						"name":                       {Type: schema.TypeString, Required: true},
						"trailing_characters_to_keep": {Type: schema.TypeString, Optional: true},
					},
				},
			},
		},

		Description: "`barracudawaf_security_policies` manages `Security Policies` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSecurityPoliciesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFSecurityPoliciesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSecurityPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFSecurityPoliciesRead(d, m)
}

func resourceCudaWAFSecurityPoliciesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
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
		"based_on": "based-on",
	}

	for tfKey, apiKey := range payload {
		if val, ok := dataItems[apiKey]; ok && val != nil {
			d.Set(tfKey, fmt.Sprintf("%v", val))
		}
	}

	// Read sub-resources
	for subResource, subResourceParams := range subResourceSecurityPoliciesParams {
		subResourceURL := strings.Replace(subResource, "_", "-", -1)
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
						subMap[param] = sortFileList(val.([]interface{}), "")
					} else {
						subMap[param] = fmt.Sprintf("%v", val)
					}
				}
			}
			subResourceList = append(subResourceList, subMap)
		}

		sortKey := ""
		for _, param := range subResourceParams {
			if param == "name" {
				sortKey = "name"
				break
			}
		}

		if len(subResourceList) > 0 {
			d.Set(subResource, sortFileList(subResourceList, sortKey))
		} else {
			d.Set(subResource, nil)
		}
	}

	return nil
}

func resourceCudaWAFSecurityPoliciesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFSecurityPoliciesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSecurityPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSecurityPoliciesRead(d, m)
}

func resourceCudaWAFSecurityPoliciesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
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

func hydrateBarracudaWAFSecurityPoliciesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{"based-on": d.Get("based_on").(string), "name": d.Get("name").(string)}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"based-on"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFSecurityPoliciesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceSecurityPoliciesParams {
		subResourceParamsLength := d.Get(subResource + ".#").(int)

		log.Printf("[INFO] Updating Barracuda WAF sub resource (%s) (%s)", name, subResource)

		for i := 0; i < subResourceParamsLength; i++ {
			subResourcePayload := make(map[string]interface{})
			suffix := fmt.Sprintf(".%d", i)

			for _, param := range subResourceParams {
				paramSuffix := fmt.Sprintf(".%s", param)
				paramValue := d.Get(subResource + suffix + paramSuffix)

				if reflect.ValueOf(paramValue).Kind() == reflect.String {
					paramValue = paramValue.(string)
				}

				if reflect.ValueOf(paramValue).Len() > 0 {
					paramKey := strings.Replace(param, "_", "-", -1)
					subResourcePayload[paramKey] = paramValue
				}
			}

			subResourceURL := strings.Replace(subResource, "_", "-", -1)
			if subResource == "global_acls" || subResource == "action_policies" || subResource == "protected_data_types" {
				itemName, ok := subResourcePayload["name"].(string)
				if ok && len(itemName) > 0 {
					subResourceURL = subResourceURL + "/" + itemName
				}
			}

			err := b.UpdateBarracudaWAFSubResource(name, endpoint, &APIRequest{
				URL:  subResourceURL,
				Body: subResourcePayload,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
