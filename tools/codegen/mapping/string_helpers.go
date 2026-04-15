package mapping

import (
	"strings"

	"github.com/ompgo-dev/ompgo/tools/codegen/model"
)

// isStringOutputCandidate reports whether a function has exactly one string
// output param (CAPIStringView* or CAPIStringBuffer*) and no other
// non-string pointer-out params (except void* and char*).
func IsStringOutputCandidate(fn model.APIFunction) bool {
	count, _ := StringOutputParamInfo(fn)
	return count == 1
}

// stringOutputParamType returns the string output param type if present.
// Returns an empty string if the function is not a string-output candidate.
func StringOutputParamType(fn model.APIFunction) string {
	count, paramType := StringOutputParamInfo(fn)
	if count != 1 {
		return ""
	}
	return paramType
}

// stringOutputParamInfo returns the count of string output params and the
// last seen string output param type.
func StringOutputParamInfo(fn model.APIFunction) (int, string) {
	stringParams := 0
	stringParamType := ""

	for _, p := range fn.Params {
		if p.Type == "CAPIStringView*" || p.Type == "CAPIStringBuffer*" {
			stringParams++
			stringParamType = p.Type
			continue
		}

		// Disqualify if there are other pointer-out params (except void* and char*).
		if strings.HasSuffix(p.Type, "*") &&
			p.Type != "void*" &&
			p.Type != "const char*" &&
			p.Type != "char*" {
			return 0, ""
		}
	}

	return stringParams, stringParamType
}

// isStringParam reports whether a C type is a string input parameter.
func IsStringParam(cType string) bool {
	return cType == "const char*" || cType == "char*"
}

// hasStringOutput reports whether a function has a single string output param.
func HasStringOutput(fn model.APIFunction) bool {
	return IsStringOutputCandidate(fn)
}
