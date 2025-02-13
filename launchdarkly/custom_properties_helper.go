package launchdarkly

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ldapi "github.com/launchdarkly/api-client-go"
)

// https://docs.launchdarkly.com/docs/custom-properties
const CUSTOM_PROPERTY_CHAR_LIMIT = 64
const CUSTOM_PROPERTY_ITEM_LIMIT = 64

func customPropertiesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Set:      customPropertyHash,
		MaxItems: CUSTOM_PROPERTY_ITEM_LIMIT,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				key: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringLenBetween(1, CUSTOM_PROPERTY_CHAR_LIMIT),
				},
				name: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringLenBetween(1, CUSTOM_PROPERTY_CHAR_LIMIT),
				},
				value: {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: CUSTOM_PROPERTY_ITEM_LIMIT,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringLenBetween(1, CUSTOM_PROPERTY_CHAR_LIMIT),
					},
				},
			},
		},
	}
}

func customPropertiesFromResourceData(d *schema.ResourceData) map[string]ldapi.CustomProperty {
	customPropertiesRaw := d.Get(custom_properties)
	schemaCustomProperties := customPropertiesRaw.(*schema.Set)
	customProperties := make(map[string]ldapi.CustomProperty)
	for _, cpRaw := range schemaCustomProperties.List() {
		key, cp := customPropertyFromResourceData(cpRaw)
		customProperties[key] = cp
	}
	return customProperties
}

func customPropertyFromResourceData(val interface{}) (string, ldapi.CustomProperty) {
	customPropertyMap := val.(map[string]interface{})

	var values []string
	for _, v := range customPropertyMap[value].([]interface{}) {
		values = append(values, v.(string))
	}
	sort.Strings(values)

	cp := ldapi.CustomProperty{
		Name:  customPropertyMap[name].(string),
		Value: values,
	}

	return customPropertyMap[key].(string), cp
}

func customPropertiesToResourceData(customProperties map[string]ldapi.CustomProperty) []interface{} {
	transformed := make([]interface{}, 0)

	for k, cp := range customProperties {
		var values []interface{}
		for _, v := range cp.Value {
			values = append(values, v)
		}
		cpRaw := map[string]interface{}{
			key:   k,
			name:  cp.Name,
			value: values,
		}
		transformed = append(transformed, cpRaw)
	}
	return transformed
}

// https://godoc.org/github.com/hashicorp/terraform/helper/schema#SchemaSetFunc
func customPropertyHash(val interface{}) int {
	customPropertyMap := val.(map[string]interface{})
	return hashcode.String(fmt.Sprintf("%v", customPropertyMap[key]))
}
