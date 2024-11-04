package main

import (
  "reflect"
  "errors"
  "fmt"
  "github.com/a-h/templ"
  // "github.com/gookit/goutil/dump"
)

type ComponentFnT struct {
  Name string
  F interface{}
}

type ComponentFnsT map[string]*ComponentFnT
var ComponentFns ComponentFnsT 


func RegisterComponent(name string, f interface{}) {
    fmt.Println("Registering component:", name);
    ComponentFns[name] = &ComponentFnT{Name: name, F: f}
}

func RegisterAllComponents() {
    if (len(ComponentFns) == 0) {
        ComponentFns=make(ComponentFnsT)
        RegisterComponent("ExampleUIComponent", ExampleUIComponent);
        RegisterComponent("AnotherUIComponent", AnotherUIComponent);
        RegisterComponent("HeroBlock", HeroBlock);
    }
}



func GetComponentFn(name string) (*ComponentFnT, error) {
   val, ok := ComponentFns[name];
   if(ok) {
      return val, nil;
   } else {
       return nil, errors.New("No component found with name:" + name);
   }
}

func CallComponent(name string, params []interface{}) (interface{}, error) {

    // dump.P("Params=", params);
    RegisterAllComponents();
    componentFn, err := GetComponentFn(name);
    f := reflect.ValueOf(componentFn.F);
    if (err != nil) {
        return f, errors.New("No component with name: " + name);
    }

    /*
    // Trying if we can directly see if component is defined.
    ff := reflect.TypeOf(name)
    if ( ff.Kind() == 0 ) {
          fmt.Println("No component with name: ", name);
          return ff, errors.New("No component with name: " + name);
    }
    f := reflect.ValueOf(name).Interface().(ComponentFnT).F;
    */

    if len(params) != f.Type().NumIn() {
        err := errors.New("CallComponent: The number of params is not sufficient.")
        return nil, err
    }

    in := make([]reflect.Value, len(params))
    for k, p := range params {
        in[k] = reflect.ValueOf(p)
    }
    var r []reflect.Value
    // fmt.Println("Printing inputs");
    // dump.P(in);

    r = f.Call(in)
    result := r[0].Interface() 

    return result, nil
}

func CallComponent_DirectParams(name string, params ... interface{}) (interface{}, error) {

    RegisterAllComponents();


    componentFn, err := GetComponentFn(name);
    if (err != nil) {
          return nil, err;
    }

    f := reflect.ValueOf(componentFn.F);

    if len(params) != f.Type().NumIn() {
        err = errors.New("CallComponent: The number of params is not sufficient.")
        return nil, err
    }

    in := make([]reflect.Value, len(params))
    for k, p := range params {
        in[k] = reflect.ValueOf(p)
    }
    var r []reflect.Value
    r = f.Call(in)
    result := r[0].Interface() 

    return result, nil
}

type export_t string

func (c export_t)GetComponent(name string, params []interface{}) (templ.Component, error) {
    // fmt.Println("Now in GetComponentHTML. Params=", params);
    r, err := CallComponent(name, params);
    if (err != nil){
        fmt.Println("CallComponent has errors:", err);
    }
    // dump.P(r);
    return r.(templ.Component), err;
}

var Exported export_t


