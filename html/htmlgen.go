package html

import (
    "fmt"
    "gopkg.in/ini.v1"
    "log"
    "regexp"
    "bytes"

    // "strings"
    // "context"
    "html/template"
    // "github.com/a-h/templ"

    "os"
    // "path/filepath"

    "github.com/microcosm-cc/bluemonday"
    "github.com/gomarkdown/markdown"

    "grssg/blocks"
    "grssg/plugins"
)


func MarkdownContent(mdtext string) string {
    maybeUnsafeHTML := markdown.ToHTML([]byte(mdtext), nil, nil)
    t := string(bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML))
    return t;
}

func TestComponent() {
    block := blocks.FindBlock("HeroBlock");
    output, err := plugins.CallComponent("ExampleUIComponent", block.Data());
    fmt.Println("output=", output, "\nError=", err);
}

func ComponentContent(block *blocks.BlockT) string {
    component := block.Lookup("Component");
    output, err := plugins.CallComponent(component, block.Data());
    if (err == nil) {
         return output;
    }
    fmt.Println("Error in component execution of section: ", block.Name(), " Error=", err);
    return "Error parsing this content:" + block.Name()
}



func HTMLTemplateContent(block *blocks.BlockT) string {
    
     text := block.Lookup("Text");

     funcMap := template.FuncMap{
         "c": func (component_name string) {
                ComponentContent(block);
            },
     }
     t := template.Must(template.New(block.Name()).Funcs(funcMap).Parse(text))

     var buff bytes.Buffer
     if err := t.Execute(&buff, block.Data()); err != nil {
         fmt.Println("Error parsing template from:", block.Name(), " Error=", err); 
         return "Error parsing this template. Name=" + block.Name();
     }

     return buff.String();
}


func GenerateBlockHTML(block *blocks.BlockT) string {

    fmt.Println("Generating text for block:", block.Name()); 

    text:=block.Lookup("Text");
    html:=""

    if component := block.Lookup("Component"); component != ""  {
        html =  ComponentContent(block);

    } else if  HTMLTemplate := block.Lookup("HTMLTemplate"); HTMLTemplate != ""  {
        html =  HTMLTemplateContent(block);

    } else if  HTML := block.Lookup("HTML"); HTML != ""  {
        html = HTML; 

    } else if Markdown := block.Lookup("Markdown"); Markdown != "" {
           html = MarkdownContent(text);

    } else {
           html = MarkdownContent(text);
    }

    return html;
}


func GeneratePageHTML(page string, section *ini.Section, output_dir string) {

    page_output := ""
    if ! section.HasKey("rows") {
        log.Fatal("There is no 'rows=' key in this section. section=", section.Name());
    }

    rows := section.Key("rows").Strings(",");
    for _, blockname := range rows {
           matched_block := blocks.FindBlock(blockname);
           if matched_block != nil {
               page_output += GenerateBlockHTML(matched_block);
           }
    }
   
    file, err2 := os.Create(output_dir + "/" + page); 
    if (err2 != nil) {
       log.Fatal("Error creating file:", page);
    }
    defer file.Close();
    _, err3 := file.WriteString(page_output);
    if (err3 != nil) {
       log.Fatal("Error writing to file:", page);
    }
    fmt.Println("Created html file:", page);
}


func GenerateSiteHTML(site_ini string, output_dir string) {
    cfg, err := ini.Load(site_ini);
    if (err != nil) {
        log.Fatal("Unable to find file site_layout.ini");
        os.Exit(-1);
    }
    for _, section := range cfg.Sections() {
        page:=section.Name();
        if page == "DEFAULT" {
           continue;    
        }
        _, err := regexp.Match(`[a-zA-Z1-9]+.html`, []byte(page));
        if (err != nil) {
           log.Fatal("site-layout.ini: Page should be simple name and end with .html: ", page);
        }
        GeneratePageHTML(page, section, output_dir); 
    }
}

