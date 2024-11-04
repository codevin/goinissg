package plugins 

import (
    "plugin"
    // "reflect"
    "context"
    "fmt"
    "os"
    "github.com/a-h/templ"
    // "html/template" Indirectly used.
)

/*
func try() {
    // Step 2: Render using templ 
    buf := new(strings.Builder)
    // component := templates.Greeting("Vinod") 

    html := "<h1>This is template.</h1>"

    // funcs := template.FuncMap{"join": strings.Join}
    t, err := template.New("ttt").Parse(string(html));
    if (err == nil) {
       component := templ.FromGoHTML(t, data)
	   component.Render(context.Background(), buf) 
       fmt.Println("Component:", buf)
   } else {
     fmt.Println("Component errr:", err);
   }
}
*/

type ComponentT struct {
  Name string
  F interface{}
}

type GetComponent_T interface {
  GetComponent(name string, params []interface{}) (templ.Component, error)
}


var component_caller GetComponent_T;

func InitComponentCaller() {
    mod := "./plugins/templcomponents/components.so"
    plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    symbol, err:= plug.Lookup("Exported");
    if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    // 3. Assert that loaded symbol is of a desired type
	// in this case interface type Greeter (defined above)
	// var caller GetComponent_T 
    var ok bool; 
	component_caller, ok = symbol.(GetComponent_T)
	if !ok {
        fmt.Println("Error: unexpected type from module symbol")
		os.Exit(1)
	}
}


// result, err1 := CallComponent("AnotherUIComponent", "param 1", "param 2")
func CallComponent(name string, params ... interface{}) ( string, error ) {
    result, err1 := component_caller.GetComponent(name, params);
    if (err1 != nil) {
        fmt.Println("Error Calling component AnotherUIComponent =", err1);
        os.Exit(0);
    }
    component:=result.(templ.Component);
    s, err := templ.ToGoHTML(context.Background(), component);
    // Or: component.Render(context.Background(), os.Stdout);
    if (err == nil) {
       return string(s), nil; 
    }
    return "", err;
}

