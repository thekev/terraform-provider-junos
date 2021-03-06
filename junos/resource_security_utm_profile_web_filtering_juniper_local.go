package junos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type utmProfileWebFilteringLocalOptions struct {
	timeout            int
	name               string
	defaultAction      string
	customBlockMessage string
	fallbackSettings   []map[string]interface{}
}

func resourceSecurityUtmProfileWebFilteringLocal() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityUtmProfileWebFilteringLocalCreate,
		Read:   resourceSecurityUtmProfileWebFilteringLocalRead,
		Update: resourceSecurityUtmProfileWebFilteringLocalUpdate,
		Delete: resourceSecurityUtmProfileWebFilteringLocalDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityUtmProfileWebFilteringLocalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"custom_block_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (s []string, es []error) {
					value := v.(string)
					if !stringInSlice(value, []string{"block", "log-and-permit", permitWord}) {
						es = append(es, fmt.Errorf(
							"%q %q invalid action", value, k))
					}

					return
				},
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
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntRange(1, 1800),
			},
		},
	}
}

func resourceSecurityUtmProfileWebFilteringLocalCreate(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return err
	}
	defer sess.closeSession(jnprSess)
	if !checkCompatibilitySecurity(jnprSess) {
		return fmt.Errorf("security utm feature-profile web-filtering juniper-local "+
			"not compatible with Junos device %s", jnprSess.Platform[0].Model)
	}
	err = sess.configLock(jnprSess)
	if err != nil {
		return err
	}
	utmProfileWebFLocalExists, err := checkUtmProfileWebFLocalExists(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	if utmProfileWebFLocalExists {
		sess.configClear(jnprSess)

		return fmt.Errorf("security utm feature-profile web-filtering juniper-local "+
			"%v already exists", d.Get("name").(string))
	}

	err = setUtmProfileWebFLocal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("create resource junos_security_utm_profile_web_filtering_juniper_local", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	mutex.Lock()
	utmProfileWebFLocalExists, err = checkUtmProfileWebFLocalExists(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if utmProfileWebFLocalExists {
		d.SetId(d.Get("name").(string))
	} else {
		return fmt.Errorf("security utm feature-profile web-filtering juniper-local %v "+
			"not exists after commit => check your config", d.Get("name").(string))
	}

	return resourceSecurityUtmProfileWebFilteringLocalRead(d, m)
}
func resourceSecurityUtmProfileWebFilteringLocalRead(d *schema.ResourceData, m interface{}) error {
	sess := m.(*Session)
	mutex.Lock()
	jnprSess, err := sess.startNewSession()
	if err != nil {
		mutex.Unlock()

		return err
	}
	defer sess.closeSession(jnprSess)
	utmProfileWebFLocalOptions, err := readUtmProfileWebFLocal(d.Get("name").(string), m, jnprSess)
	mutex.Unlock()
	if err != nil {
		return err
	}
	if utmProfileWebFLocalOptions.name == "" {
		d.SetId("")
	} else {
		fillUtmProfileWebFLocalData(d, utmProfileWebFLocalOptions)
	}

	return nil
}
func resourceSecurityUtmProfileWebFilteringLocalUpdate(d *schema.ResourceData, m interface{}) error {
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
	err = delUtmProfileWebFLocal(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = setUtmProfileWebFLocal(d, m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("update resource junos_security_utm_profile_web_filtering_juniper_local", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	d.Partial(false)

	return resourceSecurityUtmProfileWebFilteringLocalRead(d, m)
}
func resourceSecurityUtmProfileWebFilteringLocalDelete(d *schema.ResourceData, m interface{}) error {
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
	err = delUtmProfileWebFLocal(d.Get("name").(string), m, jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}
	err = sess.commitConf("delete resource junos_security_utm_profile_web_filtering_juniper_local", jnprSess)
	if err != nil {
		sess.configClear(jnprSess)

		return err
	}

	return nil
}
func resourceSecurityUtmProfileWebFilteringLocalImport(
	d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	sess := m.(*Session)
	jnprSess, err := sess.startNewSession()
	if err != nil {
		return nil, err
	}
	defer sess.closeSession(jnprSess)
	result := make([]*schema.ResourceData, 1)
	utmProfileWebFLocalExists, err := checkUtmProfileWebFLocalExists(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	if !utmProfileWebFLocalExists {
		return nil, fmt.Errorf("don't find security utm feature-profile web-filtering juniper-local with id "+
			"'%v' (id must be <name>)", d.Id())
	}
	utmProfileWebFLocalOptions, err := readUtmProfileWebFLocal(d.Id(), m, jnprSess)
	if err != nil {
		return nil, err
	}
	fillUtmProfileWebFLocalData(d, utmProfileWebFLocalOptions)

	result[0] = d

	return result, nil
}

func checkUtmProfileWebFLocalExists(profile string, m interface{}, jnprSess *NetconfObject) (bool, error) {
	sess := m.(*Session)
	profileConfig, err := sess.command("show configuration security utm feature-profile "+
		"web-filtering juniper-local profile \""+profile+"\" | display set", jnprSess)
	if err != nil {
		return false, err
	}
	if profileConfig == emptyWord {
		return false, nil
	}

	return true, nil
}
func setUtmProfileWebFLocal(d *schema.ResourceData, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0)

	setPrefix := "set security utm feature-profile web-filtering juniper-local " +
		"profile \"" + d.Get("name").(string) + "\" "
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
	if d.Get("timeout").(int) != 0 {
		configSet = append(configSet, setPrefix+"timeout "+strconv.Itoa(d.Get("timeout").(int)))
	}

	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}
func readUtmProfileWebFLocal(profile string, m interface{}, jnprSess *NetconfObject) (
	utmProfileWebFilteringLocalOptions, error) {
	sess := m.(*Session)
	var confRead utmProfileWebFilteringLocalOptions

	profileConfig, err := sess.command("show configuration security utm feature-profile web-filtering "+
		"juniper-local profile \""+profile+"\" | display set relative", jnprSess)
	if err != nil {
		return confRead, err
	}
	if profileConfig != emptyWord {
		confRead.name = profile
		for _, item := range strings.Split(profileConfig, "\n") {
			if strings.Contains(item, "<configuration-output>") {
				continue
			}
			if strings.Contains(item, "</configuration-output>") {
				break
			}
			itemTrim := strings.TrimPrefix(item, setLineStart)
			switch {
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
			case strings.HasPrefix(itemTrim, "timeout "):
				var err error
				confRead.timeout, err = strconv.Atoi(strings.TrimPrefix(itemTrim, "timeout "))
				if err != nil {
					return confRead, err
				}
			}
		}
	} else {
		confRead.name = ""

		return confRead, nil
	}

	return confRead, nil
}

func delUtmProfileWebFLocal(profile string, m interface{}, jnprSess *NetconfObject) error {
	sess := m.(*Session)
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security utm feature-profile web-filtering juniper-local "+
		"profile \""+profile+"\"")
	err := sess.configSet(configSet, jnprSess)
	if err != nil {
		return err
	}

	return nil
}

func fillUtmProfileWebFLocalData(d *schema.ResourceData,
	utmProfileWebFLocalOptions utmProfileWebFilteringLocalOptions) {
	tfErr := d.Set("name", utmProfileWebFLocalOptions.name)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("custom_block_message", utmProfileWebFLocalOptions.customBlockMessage)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("default_action", utmProfileWebFLocalOptions.defaultAction)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("fallback_settings", utmProfileWebFLocalOptions.fallbackSettings)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("timeout", utmProfileWebFLocalOptions.timeout)
	if tfErr != nil {
		panic(tfErr)
	}
}
