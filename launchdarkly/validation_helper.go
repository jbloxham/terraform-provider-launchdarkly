package launchdarkly

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func validateKey() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`),
		"Must contain only letters, numbers, '.', '-', or '_' and must start with an alphanumeric",
	)
}

func validateOp() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("Expected op to be string"))
			return
		}
		switch v {
		case
			"in",
			"endsWith",
			"startsWith",
			"matches",
			"contains",
			"lessThan",
			"lessThanOrEqual",
			"greaterThanOrEqual",
			"before",
			"after",
			"segmentMatch",
			"semVerEqual",
			"semVerLessThan",
			"semVerGreaterThan":
			return
		default:
			es = append(es, fmt.Errorf("%s is an invalid value for Clause argument Op", v))
			return
		}
	}
}
