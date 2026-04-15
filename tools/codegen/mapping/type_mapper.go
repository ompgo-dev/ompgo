package mapping

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

// TypeInfo describes a parsed C type with convenience flags.
type TypeInfo struct {
	CType              string
	Base               string
	IsPointer          bool
	IsConst            bool
	IsVoidPtr          bool
	IsCharPtr          bool
	IsCAPIStringView   bool
	IsCAPIStringBuffer bool
	IsComponentVersion bool
}

// ParseCType normalizes and parses a C type string into a TypeInfo.
func ParseCType(cType string) TypeInfo {
	trimmed := strings.TrimSpace(cType)
	info := TypeInfo{CType: trimmed}

	if strings.HasPrefix(trimmed, "const ") {
		info.IsConst = true
		trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "const "))
	}

	if strings.HasSuffix(trimmed, "*") {
		info.IsPointer = true
		info.Base = strings.TrimSpace(strings.TrimSuffix(trimmed, "*"))
	} else {
		info.Base = trimmed
	}

	info.IsVoidPtr = info.IsPointer && info.Base == "void"
	info.IsCharPtr = info.IsPointer && info.Base == "char"
	info.IsCAPIStringView = info.Base == "CAPIStringView"
	info.IsCAPIStringBuffer = info.Base == "CAPIStringBuffer"
	info.IsComponentVersion = info.Base == "ComponentVersion"

	return info
}

// GoTypeSafe maps C types to Go types for safe APIs (string-friendly).
func GoTypeSafe(cType string) string {
	info := ParseCType(cType)

	if info.IsPointer {
		if info.IsVoidPtr {
			return "unsafe.Pointer"
		}
		if info.IsCharPtr || cType == "const char*" || cType == "char*" {
			return "string"
		}
		if info.IsCAPIStringView {
			return "string"
		}
		if info.IsCAPIStringBuffer || info.IsComponentVersion {
			return "unsafe.Pointer"
		}
		return "*" + GoTypeSafe(info.Base)
	}

	switch info.Base {
	case "void":
		return ""
	case "bool":
		return "bool"
	case "int", "int32_t":
		return "int32"
	case "unsigned int", "uint", "uint32_t":
		return "uint32"
	case "int8_t":
		return "int8"
	case "uint8_t", "unsigned char", "uchar":
		return "uint8"
	case "int16_t":
		return "int16"
	case "uint16_t":
		return "uint16"
	case "int64_t":
		return "int64"
	case "uint64_t":
		return "uint64"
	case "size_t":
		return "uint"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "char":
		return "byte"
	case "CAPIStringView":
		return "string"
	case "CAPIStringBuffer", "ComponentVersion":
		return "unsafe.Pointer"
	default:
		return "interface{}"
	}
}

// GoTypeUnsafe maps C types to Go types for low-level bindings.
func GoTypeUnsafe(cType string) string {
	info := ParseCType(cType)

	if info.IsPointer {
		if info.IsVoidPtr {
			return "unsafe.Pointer"
		}
		if info.IsCharPtr || cType == "const char*" || cType == "char*" {
			return "*byte"
		}
		if info.IsCAPIStringView || info.IsCAPIStringBuffer || info.IsComponentVersion {
			return "unsafe.Pointer"
		}
		return "*" + GoTypeUnsafe(info.Base)
	}

	switch info.Base {
	case "void":
		return ""
	case "bool":
		return "bool"
	case "int", "int32_t":
		return "int32"
	case "unsigned int", "uint", "uint32_t":
		return "uint32"
	case "int8_t":
		return "int8"
	case "uint8_t", "unsigned char", "uchar":
		return "uint8"
	case "int16_t":
		return "int16"
	case "uint16_t":
		return "uint16"
	case "int64_t":
		return "int64"
	case "uint64_t":
		return "uint64"
	case "size_t":
		return "uint"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "char":
		return "byte"
	default:
		return "interface{}"
	}
}

// CTypeForHeader maps C types to explicit C header types.
func CTypeForHeader(cType string) string {
	info := ParseCType(cType)

	if info.IsPointer {
		switch info.Base {
		case "CAPIStringView":
			return "struct CAPIStringView*"
		case "CAPIStringBuffer":
			return "struct CAPIStringBuffer*"
		case "ComponentVersion":
			return "struct ComponentVersion*"
		}
		if info.IsVoidPtr {
			return "void*"
		}
		return cType
	}

	switch info.Base {
	case "void":
		return "void"
	case "bool":
		return "bool"
	case "int":
		return "int"
	case "int8_t":
		return "int8_t"
	case "int16_t":
		return "int16_t"
	case "int32_t":
		return "int32_t"
	case "int64_t":
		return "int64_t"
	case "unsigned int", "uint":
		return "uint32_t"
	case "uint8_t", "unsigned char", "uchar":
		return "uint8_t"
	case "uint16_t":
		return "uint16_t"
	case "uint32_t":
		return "uint32_t"
	case "uint64_t":
		return "uint64_t"
	case "size_t":
		return "size_t"
	case "float":
		return "float"
	case "double":
		return "double"
	case "char":
		return "char"
	case "CAPIStringView":
		return "struct CAPIStringView"
	case "CAPIStringBuffer":
		return "struct CAPIStringBuffer"
	case "ComponentVersion":
		return "struct ComponentVersion"
	case "void*":
		return "void*"
	case "char*":
		return "const char*"
	default:
		return cType
	}
}

// GoDefaultValue returns a zero/default value literal for a Go type.
func GoDefaultValue(goType string) string {
	switch goType {
	case "":
		return ""
	case "bool":
		return "false"
	case "string":
		return "\"\""
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "byte":
		return "0"
	case "unsafe.Pointer":
		return "nil"
	default:
		if strings.HasPrefix(goType, "*") {
			return "nil"
		}
		return "0"
	}
}

// CDefaultValue returns a default literal for a C type.
func CDefaultValue(cType string) string {
	info := ParseCType(cType)

	if info.IsPointer || info.Base == "void*" || info.Base == "char*" {
		return "NULL"
	}

	switch info.Base {
	case "bool":
		return "false"
	case "float", "double":
		return "0.0"
	case "int", "int8_t", "int16_t", "int32_t", "int64_t",
		"unsigned int", "uint", "uint8_t", "uint16_t", "uint32_t", "uint64_t",
		"char", "unsigned char", "uchar", "size_t":
		return "0"
	default:
		return "0"
	}
}

// CGoParamType maps a C type to a cgo parameter type.
func CGoParamType(cType string) string {
	info := ParseCType(cType)

	if info.IsPointer {
		if info.IsVoidPtr {
			return "unsafe.Pointer"
		}
		if info.IsCharPtr || cType == "const char*" || cType == "char*" {
			return "*C.char"
		}
		if info.IsCAPIStringView {
			return "C.struct_CAPIStringView"
		}
		if info.IsCAPIStringBuffer {
			return "C.struct_CAPIStringBuffer"
		}
		return "unsafe.Pointer"
	}

	switch info.Base {
	case "bool":
		return "C.bool"
	case "int", "int32_t":
		return "C.int"
	case "unsigned int", "uint", "uint32_t":
		return "C.uint32_t"
	case "int8_t":
		return "C.int8_t"
	case "uint8_t", "unsigned char", "uchar":
		return "C.uint8_t"
	case "int16_t":
		return "C.int16_t"
	case "uint16_t":
		return "C.uint16_t"
	case "int64_t":
		return "C.int64_t"
	case "uint64_t":
		return "C.uint64_t"
	case "size_t":
		return "C.size_t"
	case "float":
		return "C.float"
	case "double":
		return "C.double"
	case "CAPIStringView":
		return "C.struct_CAPIStringView"
	case "CAPIStringBuffer":
		return "C.struct_CAPIStringBuffer"
	case "ComponentVersion":
		return "C.struct_ComponentVersion"
	default:
		return "C.int"
	}
}

// CastFromC builds a Go cast expression from a C value name and type.
func CastFromC(name, cType string) string {
	switch cType {
	case "void*":
		return name
	case "const char*", "char*":
		return fmt.Sprintf("C.GoString(%s)", name)
	case "CAPIStringView":
		return fmt.Sprintf("stringFromCAPIStringView(%s)", name)
	case "bool":
		return fmt.Sprintf("bool(%s)", name)
	case "int8_t":
		return fmt.Sprintf("int8(%s)", name)
	case "uint8_t", "unsigned char", "uchar":
		return fmt.Sprintf("uint8(%s)", name)
	case "int16_t":
		return fmt.Sprintf("int16(%s)", name)
	case "uint16_t":
		return fmt.Sprintf("uint16(%s)", name)
	case "int64_t":
		return fmt.Sprintf("int64(%s)", name)
	case "uint64_t":
		return fmt.Sprintf("uint64(%s)", name)
	case "uint32_t":
		return fmt.Sprintf("uint32(%s)", name)
	case "unsigned int", "uint":
		return fmt.Sprintf("uint(%s)", name)
	case "size_t":
		return fmt.Sprintf("uint(%s)", name)
	case "float":
		return fmt.Sprintf("float32(%s)", name)
	case "double":
		return fmt.Sprintf("float64(%s)", name)
	case "int", "int32_t":
		return fmt.Sprintf("int32(%s)", name)
	default:
		return fmt.Sprintf("int32(%s)", name)
	}
}

// GoEventParamCast builds Go call expressions for event parameters.
func GoEventParamCast(name, cType string) string {
	return CastFromC(name, cType)
}

// MapCTypeToGo maps C types to Go types as strings for code generation
func MapCTypeToGo(cType string) string {
	switch cType {
	case "void*":
		return "unsafe.Pointer"
	case "const char*", "char*", "CAPIStringView":
		return "string"
	case "bool":
		return "bool"
	case "int", "int32_t":
		return "int32"
	case "unsigned int", "uint", "uint32_t":
		return "uint32"
	case "int8_t":
		return "int8"
	case "uint8_t", "unsigned char", "uchar":
		return "uint8"
	case "int16_t":
		return "int16"
	case "uint16_t":
		return "uint16"
	case "int64_t":
		return "int64"
	case "uint64_t":
		return "uint64"
	case "size_t":
		return "uint"
	case "float":
		return "float32"
	case "double":
		return "float64"
	default:
		return "int32"
	}
}

// CTypeToGoUnsafe maps C types to jen.Code for unsafe bindings.
func CTypeToGoUnsafe(cType string) jen.Code {
	if strings.HasSuffix(cType, "*") {
		baseType := strings.TrimSuffix(cType, "*")

		// void* -> unsafe.Pointer
		if baseType == "void" {
			return jen.Qual("unsafe", "Pointer")
		}

		// const char* -> *byte (or we could use string, but keeping it low-level)
		if cType == "const char*" || cType == "char*" {
			return jen.Op("*").Byte()
		}
		if baseType == "char" {
			return jen.Op("*").Byte()
		}

		// Struct pointers -> unsafe.Pointer (unknown struct types)
		if baseType == "CAPIStringView" || baseType == "CAPIStringBuffer" || baseType == "ComponentVersion" {
			return jen.Qual("unsafe", "Pointer")
		}

		// For other pointer types, make pointer to the base type
		return jen.Op("*").Add(CTypeToGoUnsafe(baseType))
	}

	switch cType {
	case "void":
		return jen.Empty()
	case "bool":
		return jen.Bool()
	case "int", "int32_t":
		return jen.Int32()
	case "unsigned int", "uint", "uint32_t":
		return jen.Uint32()
	case "int8_t":
		return jen.Int8()
	case "uint8_t", "unsigned char", "uchar":
		return jen.Uint8()
	case "int16_t":
		return jen.Int16()
	case "uint16_t":
		return jen.Uint16()
	case "int64_t":
		return jen.Int64()
	case "uint64_t":
		return jen.Uint64()
	case "size_t":
		return jen.Uint()
	case "float":
		return jen.Float32()
	case "double":
		return jen.Float64()
	case "char":
		return jen.Byte()
	default:
		return jen.Interface()
	}
}

// CTypeToGoSafe maps C types to jen.Code for safe bindings.
func CTypeToGoSafe(cType string) jen.Code {
	return CTypeToGoSafeWithName(cType, "")
}

// CTypeToGoSafeWithName maps C types to jen.Code and considers parameter naming.
func CTypeToGoSafeWithName(cType, paramName string) jen.Code {
	if strings.HasSuffix(cType, "*") {
		baseType := strings.TrimSuffix(cType, "*")

		// const char* -> string
		if cType == "const char*" || cType == "char*" {
			return jen.String()
		}

		// Struct pointers from open.mp C API
		if baseType == "CAPIStringView" || baseType == "CAPIStringBuffer" || baseType == "ComponentVersion" {
			return jen.Qual("unsafe", "Pointer")
		}

		// int* -> *int
		if baseType == "int" || baseType == "unsigned int" || baseType == "uint" {
			goType := CTypeToGoSafe(baseType)
			return jen.Op("*").Add(goType)
		}

		if baseType == "void" {
			if kind, ok := ResolveVoidPtrKind("", paramName); ok {
				return KindToEventGoType(kind)
			}
			return jen.Interface()
		}
		if baseType == "char" {
			return jen.String()
		}
		return jen.Op("*").Add(CTypeToGoSafeWithName(baseType, paramName))
	}

	switch cType {
	case "void":
		return jen.Empty()
	case "bool":
		return jen.Bool()
	case "int", "int32_t":
		return jen.Int32()
	case "unsigned int", "uint", "uint32_t":
		return jen.Uint32()
	case "int8_t":
		return jen.Int8()
	case "uint8_t", "unsigned char", "uchar":
		return jen.Uint8()
	case "int16_t":
		return jen.Int16()
	case "uint16_t":
		return jen.Uint16()
	case "int64_t":
		return jen.Int64()
	case "uint64_t":
		return jen.Uint64()
	case "size_t":
		return jen.Uint()
	case "float":
		return jen.Float32()
	case "double":
		return jen.Float64()
	case "char":
		return jen.Byte()
	// Struct types from open.mp C API
	case "CAPIStringView":
		return jen.Id("BorrowedStringView")
	case "CAPIStringBuffer", "ComponentVersion":
		return jen.Qual("unsafe", "Pointer")
	default:
		return jen.Interface()
	}
}

// CTypeToGoSafeWithEvent maps C types to jen.Code with event context.
func CTypeToGoSafeWithEvent(cType, paramName, eventName string) jen.Code {
	if cType == "void*" {
		if kind, ok := ResolveVoidPtrKind(eventName, paramName); ok {
			return KindToEventGoType(kind)
		}
		return jen.Interface()
	}
	return CTypeToGoSafeWithName(cType, paramName)
}

// CTypeToC maps C types to C header types.
func CTypeToC(cType string) string {
	return CTypeForHeader(cType)
}

// cTypeToC maps C types to C header types.
func cTypeToC(cType string) string {
	return CTypeForHeader(cType)
}

// DefaultReturnValue returns a default C return value literal.
func DefaultReturnValue(cType string) string {
	return CDefaultValue(cType)
}

// getDefaultReturnValue returns a default C return value literal.
func getDefaultReturnValue(cType string) string {
	return CDefaultValue(cType)
}

// GoTypeForCType maps C types to Go types for cgo wrappers.
func GoTypeForCType(cType string) string {
	return goTypeForCType(cType)
}

// goTypeForCType maps C types to Go types for cgo wrappers.
func goTypeForCType(cType string) string {
	info := ParseCType(cType)
	if info.IsPointer {
		switch info.Base {
		case "void":
			return "handle.Handle"
		case "CAPIStringView":
			return "*C.struct_CAPIStringView"
		case "CAPIStringBuffer":
			return "*C.struct_CAPIStringBuffer"
		case "ComponentVersion":
			return "*C.struct_ComponentVersion"
		}
		if info.IsCharPtr || cType == "const char*" || cType == "char*" {
			return "string"
		}
		return "*" + GoTypeSafe(info.Base)
	}

	switch info.Base {
	case "void":
		return ""
	case "CAPIStringView":
		return "C.struct_CAPIStringView"
	case "CAPIStringBuffer":
		return "C.struct_CAPIStringBuffer"
	case "ComponentVersion":
		return "C.struct_ComponentVersion"
	default:
		return GoTypeSafe(info.Base)
	}
}

// GoArgForCType builds cgo argument expressions for a given C type.
func GoArgForCType(argName, cType string) string {
	return goArgForCType(argName, cType)
}

// goArgForCType builds cgo argument expressions for a given C type.
func goArgForCType(argName, cType string) string {
	safeName := SafeGoName(argName)
	info := ParseCType(cType)

	if info.IsPointer {
		switch info.Base {
		case "void":
			// handle.Handle -> unsafe.Pointer for cgo calls (runtime package provides ptrFromHandle)
			return "ptrFromHandle(" + safeName + ")"
		case "CAPIStringView", "CAPIStringBuffer", "ComponentVersion":
			return safeName
		case "char":
			if cType == "const char*" || cType == "char*" {
				// string -> *C.char setup is done in the wrapper
				return "c_" + safeName
			}
			return "(*C.char)(unsafe.Pointer(" + safeName + "))"
		case "bool":
			return "(*C.bool)(unsafe.Pointer(" + safeName + "))"
		case "int", "int32_t":
			return "(*C.int)(unsafe.Pointer(" + safeName + "))"
		case "unsigned int", "uint", "uint32_t":
			return "(*C.uint)(unsafe.Pointer(" + safeName + "))"
		case "int8_t":
			return "(*C.int8_t)(unsafe.Pointer(" + safeName + "))"
		case "uint8_t", "unsigned char", "uchar":
			return "(*C.uint8_t)(unsafe.Pointer(" + safeName + "))"
		case "int16_t":
			return "(*C.int16_t)(unsafe.Pointer(" + safeName + "))"
		case "uint16_t":
			return "(*C.uint16_t)(unsafe.Pointer(" + safeName + "))"
		case "int64_t":
			return "(*C.int64_t)(unsafe.Pointer(" + safeName + "))"
		case "uint64_t":
			return "(*C.uint64_t)(unsafe.Pointer(" + safeName + "))"
		case "size_t":
			return "(*C.size_t)(unsafe.Pointer(" + safeName + "))"
		case "float":
			return "(*C.float)(unsafe.Pointer(" + safeName + "))"
		case "double":
			return "(*C.double)(unsafe.Pointer(" + safeName + "))"
		default:
			return safeName
		}
	}

	switch info.Base {
	case "bool":
		return fmt.Sprintf("C.bool(%s)", safeName)
	case "int", "int32_t":
		return fmt.Sprintf("C.int(%s)", safeName)
	case "unsigned int", "uint":
		return fmt.Sprintf("C.uint(%s)", safeName)
	case "uint32_t":
		return fmt.Sprintf("C.uint32_t(%s)", safeName)
	case "int8_t":
		return fmt.Sprintf("C.int8_t(%s)", safeName)
	case "uint8_t", "unsigned char", "uchar":
		return fmt.Sprintf("C.uint8_t(%s)", safeName)
	case "int16_t":
		return fmt.Sprintf("C.int16_t(%s)", safeName)
	case "uint16_t":
		return fmt.Sprintf("C.uint16_t(%s)", safeName)
	case "int64_t":
		return fmt.Sprintf("C.int64_t(%s)", safeName)
	case "uint64_t":
		return fmt.Sprintf("C.uint64_t(%s)", safeName)
	case "size_t":
		return fmt.Sprintf("C.size_t(%s)", safeName)
	case "float":
		return fmt.Sprintf("C.float(%s)", safeName)
	case "double":
		return fmt.Sprintf("C.double(%s)", safeName)
	case "CAPIStringView", "CAPIStringBuffer", "ComponentVersion":
		return safeName
	case "char":
		// Not expected (char params should be pointers), but keep behavior safe.
		return safeName
	default:
		return safeName
	}
}

// DefaultGoReturn returns a default Go return value for a C type.
func DefaultGoReturn(ret string) string {
	return defaultGoReturn(ret)
}

// defaultGoReturn returns a default Go return value for a C type.
func defaultGoReturn(ret string) string {
	goType := GoTypeSafe(ret)
	return GoDefaultValue(goType)
}

// CEventParamType maps event C types to cgo parameter types.
func CEventParamType(cType string) string {
	return cEventParamType(cType)
}

// cEventParamType maps event C types to cgo parameter types.
func cEventParamType(cType string) string {
	return CGoParamType(cType)
}

// MapGoTypeToCType maps Go types to C types for cgo.
func MapGoTypeToCType(goType string) string {
	return mapGoTypeToCType(goType)
}

// mapGoTypeToCType maps Go types to C types for cgo.
func mapGoTypeToCType(goType string) string {
	switch goType {
	case "int32", "int":
		return "C.int"
	case "float32":
		return "C.float"
	case "bool":
		return "C.bool"
	case "string":
		return "*C.char"
	case "unsafe.Pointer":
		return "unsafe.Pointer"
	default:
		return "C.int"
	}
}
