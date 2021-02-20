package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type MethodDiff struct {
	Params    `json:"parameters,omitempty"`
}

type Params struct {
	AddedParams    ParamNamesByLocation `json:"added,omitempty"`
	DeletedParams  ParamNamesByLocation `json:"deleted,omitempty"`
	ModifiedParams ParamDiffByLocation  `json:"modified,omitempty"`
}

// ParamNamesByLocation maps param location (path, query, header or cookie) to the params in this location
type ParamNamesByLocation map[string]ParamNames

// ParamDiffByLocation maps param location (path, query, header or cookie) to param diffs in this location
type ParamDiffByLocation map[string]ParamDiffs

func newParams() Params {
	return Params{
		AddedParams:    ParamNamesByLocation{},
		DeletedParams:  ParamNamesByLocation{},
		ModifiedParams: ParamDiffByLocation{},
	}
}

// ParamNames is a set of parameter names
type ParamNames map[string]struct{}

// ParamDiffs is map of parameter names to their respective diffs
type ParamDiffs map[string]ParamDiff

func newMethodDiff() *MethodDiff {
	return &MethodDiff{
		Params: newParams(),
	}
}

func (methodDiff *MethodDiff) empty() bool {
	return len(methodDiff.AddedParams) == 0 &&
		len(methodDiff.DeletedParams) == 0 &&
		len(methodDiff.ModifiedParams) == 0
}

func (methodDiff *MethodDiff) addAddedParam(param *openapi3.Parameter) {

	if paramNames, ok := methodDiff.AddedParams[param.In]; ok {
		paramNames[param.Name] = struct{}{}
	} else {
		methodDiff.AddedParams[param.In] = ParamNames{param.Name: struct{}{}}
	}
}

func (methodDiff *MethodDiff) addDeletedParam(param *openapi3.Parameter) {

	if paramNames, ok := methodDiff.DeletedParams[param.In]; ok {
		paramNames[param.Name] = struct{}{}
	} else {
		methodDiff.DeletedParams[param.In] = ParamNames{param.Name: struct{}{}}
	}
}

func (methodDiff *MethodDiff) addModifiedParam(param *openapi3.Parameter, diff ParamDiff) {

	if paramDiffs, ok := methodDiff.ModifiedParams[param.In]; ok {
		paramDiffs[param.Name] = diff
	} else {
		methodDiff.ModifiedParams[param.In] = ParamDiffs{param.Name: diff}
	}
}

func diffParameters(params1 openapi3.Parameters, params2 openapi3.Parameters) *MethodDiff {

	result := newMethodDiff()

	for _, paramRef1 := range params1 {
		if paramRef1 != nil && paramRef1.Value != nil {
			if paramValue2, ok := findParam(paramRef1.Value, params2); ok {
				if diff := diffParamValues(paramRef1.Value, paramValue2); !diff.empty() {
					result.addModifiedParam(paramRef1.Value, diff)
				}
			} else {
				result.addDeletedParam(paramRef1.Value)
			}
		}
	}

	for _, paramRef2 := range params2 {
		if paramRef2 != nil && paramRef2.Value != nil {
			if _, ok := findParam(paramRef2.Value, params1); !ok {
				result.addAddedParam(paramRef2.Value)
			}
		}
	}

	return result
}

func findParam(param1 *openapi3.Parameter, params2 openapi3.Parameters) (*openapi3.Parameter, bool) {
	// TODO: optimize with a map
	for _, paramRef2 := range params2 {
		if paramRef2 == nil || paramRef2.Value == nil {
			continue
		}

		if equalParams(param1, paramRef2.Value) {
			return paramRef2.Value, true
		}
	}

	return nil, false
}

func equalParams(param1 *openapi3.Parameter, param2 *openapi3.Parameter) bool {
	return param1.Name == param2.Name && param1.In == param2.In
}
