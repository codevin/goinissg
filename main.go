package main

import (
    "grssg/blocks"
    "grssg/html"
    "grssg/plugins"
    // "gopkg.in/ini.v1"
    // "grssg/templcomponents"

    // "strings"
    // "reflect"
    // "context"
    "fmt"
    // "os"
    // "html/template"
    // "github.com/a-h/templ"
	"net/http"
    "log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
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

func gen_site() {
    html.GenerateSiteHTML("site-layout.ini", "./output");
}

func test_component() {
    // output, err  := plugins.CallComponent("AnotherUIComponent", "param 1", "param 2")
    // fmt.Println("output=", output, "\nError=", err);


    // In tmpl file, this component calls another component. 
    block:=blocks.FindBlock("HeroBlock");
    output, err := plugins.CallComponent("ExampleUIComponent", block.Data());
    fmt.Println("output=", output, "\nError=", err);

}

func ui() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Form Widget")

	entry := widget.NewEntry()
	textArea := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Form text Entry", Widget: entry},
        },
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			log.Println("multiline:", textArea.Text)
			myWindow.Close()
		},
	}

	// we can also append items
	form.Append("Textarea Text", textArea)

	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}



func main() {
    plugins.InitComponentCaller();
    root := "./content"
    blocks.Init(root);

    // test_component();
    // os.Exit(0);

    gen_site();

    ui();

	fs := http.FileServer(http.Dir("./assets"))
    // This needs to be used in this way only. Note the trailing / on both handles.
	http.Handle("/sbadmin/", http.StripPrefix("/sbadmin/", fs))
    // Reason: fs needs portion after /static/. It has no awareness of URL.

	fs2 := http.FileServer(http.Dir("./output"))
    // This needs to be used in this way only. Note the trailing / on both handles.
	http.Handle("/", http.StripPrefix("/", fs2))
    // Reason: fs needs portion after /static/. It has no awareness of URL.

   /*
	log.Print("Listening on :3000...")
    err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
    */
}


