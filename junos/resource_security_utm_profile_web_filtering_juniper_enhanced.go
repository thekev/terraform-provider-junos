package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type utmProfileWebFilteringEnhancedOptions struct {
	noSafeSearch            bool
	timeout                 int
	name                    string
	defaultAction           string
	customBlockMessage      string
	quarantineCustomMessage string
	blockMessage            []map[string]interface{}
	category                []map[string]interface{}
	fallbackSettings        []map[string]interface{}
	quarantineMessage       []map[string]interface{}
	siteReputationAction    []map[string]interface{}
}

func resourceSecurityUtmProfileWebFilteringEnhanced() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityUtmProfileWebFilteringEnhancedCreate,
		Read:   resourceSecurityUtmProfileWebFilteringEnhancedRead,
		Update: resourceSecurityUtmProfileWebFilteringEnhancedUpdate,
		Delete: resourceSecurityUtmProfileWebFilteringEnhancedDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityUtmProfileWebFilteringEnhancedImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"block_message": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type_custom_redirect_url": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"category": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateNameObjectJunos(),
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateUtmProfileAction(),
						},
						"reputation_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"site_reputation": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											if !stringInSlice(value, []string{"fairly-safe", "harmful", "moderately-safe", "suspicious", "very-safe"}) {
												errors = append(errors, fmt.Errorf(
													"%q %q invalid reputation", value, k))
											}

											return
										},
									},
									"action": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateUtmProfileAction(),
									},
								},
							},
						},
					},
				},
			},
			"custom_block_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateUtmProfileAction(),
			},
			"fallback_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"block", "log-and-permit"}) {
									errors = append(errors, fmt.Errorf(
										"%q %q invalid action", value, k))
								}

								return
							},
						},
						"server_connectivity": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"block", "log-and-permit"}) {
									errors = append(errors, fmt.Errorf(
										"%q %q invalid action", value, k))
								}

								return
							},
						},
						"timeout": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"block", "log-and-permit"}) {
									errors = append(errors, fmt.Errorf(
										"%q %q invalid action", value, k))
								}

								return
							},
						},
						"too_many_requests": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"block", "log-and-permit"}) {
									errors = append(errors, fmt.Errorf(
										"%q %q invalid action", value, k))
								}

								return
							},
						},
					},
				},
			},
			"no_safe_search": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"quarantine_custom_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quarantine_message": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type_custom_redirect_url": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"site_reputation_action": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_reputation": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if !stringInSlice(value, []string{"fairly-safe", "harmful", "moderately-safe", "suspicious", "very-safe"}) {
									errors = append(errors, fmt.Errorf(
										"%q %q invalid reputation", value, k))
								}

								return
							},
						},
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateUtmProfileAction(),
						},
					},
				},
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 1800),
			},
		},
	}
}

func validateUtmProfileAction() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (s []string, es []error) {
		value := v.(string)
		if !stringInSlice(value, []string{"block", "log-and-permit", permitWord, "quarantine"}) {
			es = append(es, fmt.Errorf(
				"%q %q invalid action", value, k))
		}

		return
	}
}

func resourceSecurityUtmProfileWebFilteringEnhancedCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security utm feature-profile web-filtering juniper-enhanced "+
			"not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	utmProfileWebFEnhancedExists, err := checkUtmProfileWebFEnhancedExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if utmProfileWebFEnhancedExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("security utm feature-profile web-filtering juniper-enhanced "+
			"%v already exists", d.Get("name").(string))
	}

	err = setUtmProfileWebFEnhanced(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_security_utm_profile_web_filtering_juniper_enhanced", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	mutex.Lock()
	utmProfileWebFEnhancedExists, err = checkUtmProfileWebFEnhancedExists(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if utmProfileWebFEnhancedExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security utm feature-profile web-filtering juniper-enhanced %v "+
			"not exists after commit => check your config", d.Get("name").(string))
	}

	return resourceSecurityUtmProfileWebFilteringEnhancedRead(d, m)
}
func resourceSecurityUtmProfileWebFilteringEnhancedRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	utmProfileWebFEnhancedOptions, err := readUtmProfileWebFEnhanced(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if utmProfileWebFEnhancedOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmProfileWebFEnhancedData(d, utmProfileWebFEnhancedOptions)
	}

	return nil
}
func resourceSecurityUtmProfileWebFilteringEnhancedUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delUtmProfileWebFEnhanced(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setUtmProfileWebFEnhanced(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_security_utm_profile_web_filtering_juniper_enhanced", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceSecurityUtmProfileWebFilteringEnhancedRead(d, m)
}
func resourceSecurityUtmProfileWebFilteringEnhancedDelete(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	err = delUtmProfileWebFEnhanced(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_security_utm_profile_web_filtering_juniper_enhanced", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
}
func resourceSecurityUtmProfileWebFilteringEnhancedImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	utmProfileWebFEnhancedExists, err := checkUtmProfileWebFEnhancedExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !utmProfileWebFEnhancedExists {
		return nil, fmt.Errorf("don't find security utm feature-profile web-filtering juniper-enhanced with id "+
			"'%v' (id must be <name>)", d.Id())
	}
	utmProfileWebFEnhancedOptions, err := readUtmProfileWebFEnhanced(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillUtmProfileWebFEnhancedData(d, utmProfileWebFEnhancedOptions)

	result[0] = d

	return result, nil
}

func checkUtmProfileWebFEnhancedExists(profile string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	profileConfig, err := sess.command("show configuration security utm feature-profile "+
		"web-filtering juniper-enhanced profile \""+profile+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if profileConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setUtmProfileWebFEnhanced(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security utm feature-profile web-filtering juniper-enhanced " +
		"profile \"" + d.Get("name").(string) + "\" "
	for _, v := range d.Get("block_message").([]interface{}) {
		if v != nil {
			message := v.(map[string]interface{})
			if message["url"].(string) != "" {
				configSet = append(configSet, setPrefix+"block-message url \""+message["url"].(string)+"\"")
			}
			if message["type_custom_redirect_url"].(bool) {
				configSet = append(configSet, setPrefix+"block-message type custom-redirect-url")
			}
		} else {
			configSet = append(configSet, setPrefix+"block-message")
		}
	}
	for _, v := range d.Get("category").([]interface{}) {
		category := v.(map[string]interface{})
		setPrefixCategory := setPrefix + "category \"" + category["name"].(string) + "\" "
		configSet = append(configSet, setPrefixCategory+"action "+category["action"].(string))
		for _, r := range category["reputation_action"].([]interface{}) {
			reputation := r.(map[string]interface{})
			configSet = append(configSet, setPrefixCategory+"reputation-action "+
				reputation["site_reputation"].(string)+" "+reputation["action"].(string))
		}
	}
	if d.Get("custom_block_message").(string) != "" {
		configSet = append(configSet, setPrefix+"custom-block-message \""+d.Get("custom_block_message").(string)+"\"")
	}
	if d.Get("default_action").(string) != "" {
		configSet = append(configSet, setPrefix+"default "+d.Get("default_action").(string))
	}
	for _, v := range d.Get("fallback_settings").([]interface{}) {
		if v != nil {
			fSettings := v.(map[string]interface{})
			if fSettings["default"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings default "+
					fSettings["default"].(string))
			}
			if fSettings["server_connectivity"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings server-connectivity "+
					fSettings["server_connectivity"].(string))
			}
			if fSettings["timeout"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings timeout "+
					fSettings["timeout"].(string))
			}
			if fSettings["too_many_requests"].(string) != "" {
				configSet = append(configSet, setPrefix+"fallback-settings too-many-requests "+
					fSettings["too_many_requests"].(string))
			}
		} else {
			configSet = append(configSet, setPrefix+"fallback-settings")
		}
	}
	if d.Get("no_safe_search").(bool) {
		configSet = append(configSet, setPrefix+"no-safe-search")
	}
	if d.Get("quarantine_custom_message").(string) != "" {
		configSet = append(configSet,
			setPrefix+"quarantine-custom-message \""+d.Get("quarantine_custom_message").(string)+"\"")
	}
	for _, v := range d.Get("quarantine_message").([]interface{}) {
		if v != nil {
			message := v.(map[string]interface{})
			if message["url"].(string) != "" {
				configSet = append(configSet, setPrefix+"quarantine-message url \""+message["url"].(string)+"\"")
			}
			if message["type_custom_redirect_url"].(bool) {
				configSet = append(configSet, setPrefix+"quarantine-message type custom-redirect-url")
			}
		} else {
			configSet = append(configSet, setPrefix+"quarantine-message")
		}
	}
	for _, v := range d.Get("site_reputation_action").([]interface{}) {
		siteReputation := v.(map[string]interface{})
		configSet = append(configSet, setPrefix+"site-reputation-action "+
			siteReputation["site_reputation"].(string)+" "+siteReputation["action"].(string))
	}
	if d.Get("timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+"timeout "+strconv.Itoa(d.Get("timeout").(int)))
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func readUtmProfileWebFEnhanced(profile string, m interface{}, jnprSess *NetconfObject) (
	utmProfileWebFilteringEnhancedOptions, error) {
	sess := m.(*Session)
	var confRead utmProfileWebFilteringEnhancedOptions

	profileConfig, err := sess.command("show configuration security utm feature-profile web-filtering "+
		"juniper-enhanced profile \""+profile+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if profileConfig != emptyWord {
		confRead.name = profile
		categoryList := make([]map[string]interface{}, 0)
		for _, item := range strings.Split(profileConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
			case strings.HasPrefix(itemTrim, "block-message"):
				if len(confRead.blockMessage) == 0 {
					confRead.blockMessage = append(confRead.blockMessage, map[string]interface{}{
						"url":                      "",
						"type_custom_redirect_url": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "block-message type custom-redirect-url"):
					confRead.blockMessage[0]["type_custom_redirect_url"] = true
				case strings.HasPrefix(itemTrim, "block-message url "):
					confRead.blockMessage[0]["url"] = strings.Trim(strings.TrimPrefix(itemTrim, "block-message url "), "\"")
				}
			case strings.HasPrefix(itemTrim, "category "):
				catergoryLineCut := strings.Split(itemTrim, " ")
				c := map[string]interface{}{
					"name":              catergoryLineCut[1],
					"action":            "",
					"reputation_action": make([]map[string]interface{}, 0),
				}
				c, categoryList = copyAndRemoveItemMapList("name", false, c, categoryList)
				itemTrimCategory := strings.TrimPrefix(itemTrim, "category "+catergoryLineCut[1]+" ")
				switch {
				case strings.HasPrefix(itemTrimCategory, "action "):
					c["action"] = strings.TrimPrefix(itemTrimCategory, "action ")
				case strings.HasPrefix(itemTrimCategory, "reputation-action "):
					cutReputationAction := strings.Split(strings.TrimPrefix(itemTrimCategory, "reputation-action "), " ")
					c["reputation_action"] = append(c["reputation_action"].([]map[string]interface{}), map[string]interface{}{
						"site_reputation": cutReputationAction[0],
						"action":          cutReputationAction[1],
					})
				}
				categoryList = append(categoryList, c)
			case strings.HasPrefix(itemTrim, "custom-block-message "):
				confRead.customBlockMessage = strings.Trim(strings.TrimPrefix(itemTrim, "custom-block-message "), "\"")
			case strings.HasPrefix(itemTrim, "default "):
				confRead.defaultAction = strings.TrimPrefix(itemTrim, "default ")
			case strings.HasPrefix(itemTrim, "fallback-settings"):
				if len(confRead.fallbackSettings) == 0 {
					confRead.fallbackSettings = append(confRead.fallbackSettings, map[string]interface{}{
						"default":             "",
						"server_connectivity": "",
						"timeout":             "",
						"too_many_requests":   "",
					})
				}
				itemTrimFallback := strings.TrimPrefix(itemTrim, "fallback-settings ")
				switch {
				case strings.HasPrefix(itemTrimFallback, "default "):
					confRead.fallbackSettings[0]["default"] = strings.TrimPrefix(itemTrimFallback, "default ")
				case strings.HasPrefix(itemTrimFallback, "server-connectivity "):
					confRead.fallbackSettings[0]["server_connectivity"] = strings.TrimPrefix(itemTrimFallback, "server-connectivity ")
				case strings.HasPrefix(itemTrimFallback, "timeout "):
					confRead.fallbackSettings[0]["timeout"] = strings.TrimPrefix(itemTrimFallback, "timeout ")
				case strings.HasPrefix(itemTrimFallback, "too-many-requests "):
					confRead.fallbackSettings[0]["too_many_requests"] = strings.TrimPrefix(itemTrimFallback, "too-many-requests ")
				}
			case strings.HasSuffix(itemTrim, "no-safe-search"):
				confRead.noSafeSearch = true
			case strings.HasPrefix(itemTrim, "quarantine-custom-message "):
				confRead.quarantineCustomMessage = strings.Trim(strings.TrimPrefix(itemTrim, "quarantine-custom-message "), "\"")
			case strings.HasPrefix(itemTrim, "quarantine-message"):
				if len(confRead.quarantineMessage) == 0 {
					confRead.quarantineMessage = append(confRead.quarantineMessage, map[string]interface{}{
						"url":                      "",
						"type_custom_redirect_url": false,
					})
				}
				switch {
				case strings.HasPrefix(itemTrim, "quarantine-message type custom-redirect-url"):
					confRead.quarantineMessage[0]["type_custom_redirect_url"] = true
				case strings.HasPrefix(itemTrim, "quarantine-message url "):
					confRead.quarantineMessage[0]["url"] = strings.Trim(strings.TrimPrefix(itemTrim, "quarantine-message url "), "\"")
				}
			case strings.HasPrefix(itemTrim, "site-reputation-action "):
				itemTrimSiteReput := strings.Split(strings.TrimPrefix(itemTrim, "site-reputation-action "), " ")
				confRead.siteReputationAction = append(confRead.siteReputationAction, map[string]interface{}{
					"site_reputation": itemTrimSiteReput[0],
					"action":          itemTrimSiteReput[1],
				})
			case strings.HasPrefix(itemTrim, "timeout "):
				var err error
				confRead.timeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "timeout "))
				if err != nil {
					return confRead, err
				}
			}
		}
		confRead.category = categoryList
	} else {
		confRead.name = ""

		return confRead, nil
	}

	return confRead, nil
}

func delUtmProfileWebFEnhanced(profile string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm feature-profile web-filtering juniper-enhanced "+
		"profile \""+profile+"\"")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}

func fillUtmProfileWebFEnhancedData(d *schema.ResourceData,
	utmProfileWebFEnhancedOptions utmProfileWebFilteringEnhancedOptions) {
	tfErr := d.Set("name", utmProfileWebFEnhancedOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("block_message", utmProfileWebFEnhancedOptions.blockMessage)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("category", utmProfileWebFEnhancedOptions.category)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("custom_block_message", utmProfileWebFEnhancedOptions.customBlockMessage)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("default_action", utmProfileWebFEnhancedOptions.defaultAction)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("fallback_settings", utmProfileWebFEnhancedOptions.fallbackSettings)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("no_safe_search", utmProfileWebFEnhancedOptions.noSafeSearch)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("quarantine_custom_message", utmProfileWebFEnhancedOptions.quarantineCustomMessage)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("quarantine_message", utmProfileWebFEnhancedOptions.quarantineMessage)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("site_reputation_action", utmProfileWebFEnhancedOptions.siteReputationAction)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("timeout", utmProfileWebFEnhancedOptions.timeout)
	if tfErr != nil {
		panic(tfErr)
	}
}
