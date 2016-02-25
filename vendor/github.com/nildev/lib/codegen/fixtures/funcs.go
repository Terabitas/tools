package fixtures

import (
	"github.com/nildev/lib/codegen/fixtures/sub"
	alias "github.com/nildev/lib/codegen/fixtures/suby"
)

type (
	MyStruct struct {
	}

	MyOtherStruct struct {
		F MyStruct
	}
)

// CheckPrimitive to check primitive types
// @method GET
func CheckPrimitive(
	msg string,
	num int,
	num8 int8,
	num16 int16,
	num32 int32,
	num64 int64,
	flt32 float32,
	flt64 float64,
) (rmsg string, rnum int, rnum8 int8, rnum16 int16, rnum32 int32, rnum64 int64, rflt32 float32, rflt64 float64) {

	return rmsg, rnum, rnum8, rnum16, rnum32, rnum64, rflt32, rflt64
}

// CheckNestedMap to check primitive types
// @method DELETE
func CheckNestedMap(
	mapSS map[string]string,
	mapSMI map[string]map[string]int,
	mapASAMSAS map[int][]map[string][]string,
	mapPSSI map[*string]sub.MyInterface,
	mapMPSI map[*MyStruct]int,
	mapMSI map[MyStruct]int,
	mapMSII map[sub.MyInterface]int,
	mapMSIISI map[struct{}]interface{},
	x ...map[sub.MyInterface]int,
) (rmapSS map[string]string, rmapSMI map[string]map[string]int, rmapASAMSAS map[int][]map[string][]string, rmapPSSI map[*string]sub.MyInterface, rmapMPSI map[*MyStruct]int, rmapMSI map[MyStruct]int, rmapMSII map[sub.MyInterface]int, rmapMSIISI map[struct{}]interface{}) {

	return rmapSS, rmapSMI, rmapASAMSAS, rmapPSSI, rmapMPSI, rmapMSI, rmapMSII, rmapMSIISI
}

// CheckVariadicParam to check primitive types
// @method PUT
func CheckVariadicParam(
	x ...map[sub.MyInterface]int,
) {

}

// CheckComplex to check primitive types
// @method POST
func CheckComplex(
	str1 MyStruct,
	strPtr *MyStruct,
	strSubPtr *sub.MyStruct,
	str2 MyOtherStruct,
	inteM sub.MyInterface,
	inte interface{},
	st struct{},
	ai alias.MyInterface,
) (rstr1 MyStruct, rstrPtr *MyStruct, rstrSubPtr *sub.MyStruct, rstr2 MyOtherStruct, rinteM sub.MyInterface, rinte interface{}, rst struct{}, rai alias.MyInterface) {

	return rstr1, rstrPtr, rstrSubPtr, rstr2, rinteM, rinte, rst, rai
}

// CheckGetter getter
// All signature params will have to be passed in query string or path
// Params which have type of pointer will be optional and wont be included in query
// If no validation rule is required than you do not need to add to @query it will be taken from signature
// @method GET
// @path /my-path/{varOne}/{xxx}/sub-resource/{second:[a-z]+}
// @query {mineRegex:[A-Z]+} {mineRegex2:[0-9]{4}[a-z]{5}} {mineRegex26:[0-9]{4}[a-z]{5}} {notOptional}
func CheckGetter(varOne int, xxx string, second string, mineRegex string, mineRegex2 string, mineRegex26 string, notOptional string, optionalParam *string) (rez string) {
	return rez
}

// CheckPoster post
// All signature params will have to be included in POST body
// Only params from @path and @query will be passed to function from path and query string
// @method POST
// @path /my-path/{varOne}/{xxx}/sub-resource/{second:[a-z]+}
// @query {mineRegex:[A-Z]+} {notOptional}
func CheckPoster(varOne int, xxx string, second string, mineRegex string, mineRegex2 string, mineRegex26 string, notOptional string, optionalParam *string) (rez string) {
	return rez
}

// CheckNone to check primitive types
// @protected
func CheckNone(boo string) (rez string) {
	return rez
}
