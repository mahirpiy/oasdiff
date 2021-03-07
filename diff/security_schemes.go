package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// SecuritySchemesDiff is a diff between two sets of security scheme objects: https://swagger.io/specification/#security-scheme-object
type SecuritySchemesDiff struct {
	Added    StringList              `json:"added,omitempty"`
	Deleted  StringList              `json:"deleted,omitempty"`
	Modified ModifiedSecuritySchemes `json:"modified,omitempty"`
}

func (diff *SecuritySchemesDiff) empty() bool {
	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

// ModifiedSecuritySchemes is map of security schemes to their respective diffs
type ModifiedSecuritySchemes map[string]SecuritySchemeDiff

func newSecuritySchemesDiff() *SecuritySchemesDiff {
	return &SecuritySchemesDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedSecuritySchemes{},
	}
}

func getSecuritySchemesDiff(config *Config, securitySchemes1, securitySchemes2 openapi3.SecuritySchemes) *SecuritySchemesDiff {

	result := newSecuritySchemesDiff()

	for value1, ref1 := range securitySchemes1 {
		if ref1 != nil && ref1.Value != nil {
			if value2, ok := securitySchemes2[value1]; ok {
				if diff := getSecuritySchemeDiff(config, ref1.Value, value2.Value); !diff.empty() {
					result.Modified[value1] = diff
				}
			} else {
				result.Deleted = append(result.Deleted, value1)
			}
		}
	}

	for value2, ref2 := range securitySchemes2 {
		if ref2 != nil && ref2.Value != nil {
			if _, ok := securitySchemes1[value2]; !ok {
				result.Added = append(result.Added, value2)
			}
		}
	}

	return result

}

func (diff *SecuritySchemesDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(diff.Added),
		Deleted:  len(diff.Deleted),
		Modified: len(diff.Modified),
	}
}