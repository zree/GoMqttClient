package example

import (
	"reflect"
	"fmt"
	"strings"
)

func PVSexplainer(pvsString string){

}

func PrintStruct(t reflect.Type, v reflect.Value, pc int){
	fmt.Println("")
	for i := 0; i<t.NumField(); i++ {
		fmt.Print(strings.Repeat(" ", pc), t.Field(i).Name, ":")
		value := v.Field(i)
		PrintVar(value.Interface(), pc+2)
		fmt.Println("")
	}
}

func PrintVar(i interface{}, ident int){

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr{
		v = reflect.ValueOf(i).Elem()
	}
	switch  v.Kind() {
	case reflect.Array:
		PrintArraySlice(v,ident)
	case reflect.Chan:
		fmt.Println("Chan")
	case reflect.Func:
		fmt.Println("Func")
	case reflect.Interface:
		fmt.Println("Interface")
	case reflect.Map:
		PrintMap(v,ident)
	case reflect.UnsafePointer:
		fmt.Println("UNsafePointer")
	default:
		fmt.Print(strings.Repeat(" ",ident),v.Interface())
	}
}

func PrintArraySlice(v reflect.Value,pc int){
	for j:=0;j<v.Len();j++{
		PrintVar(v.Index(j).Interface(),pc+2)
	}
}

func PrintMap(v reflect.Value,pc int){
	for _, k :=range v.MapKeys(){
		PrintVar(k.Interface(),pc)
		PrintVar(v.MapIndex(k).Interface(),pc)
	}
}