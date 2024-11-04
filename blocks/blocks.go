package blocks

import (
    "fmt"
    "gopkg.in/ini.v1"
    "log"

    "os"
	"path/filepath"

    "github.com/microcosm-cc/bluemonday"
    "github.com/gomarkdown/markdown"
)

type BlockT struct {
  name string; 
  filepath string;
  section *ini.Section;
  data map[string]interface{};
}


var ContentBlocks []*BlockT;


func (block *BlockT) Name() string {
   return block.name;
}

func CreateBlock(name string, filepath string, section *ini.Section) (block *BlockT) {
    new_block := &BlockT{name: name, filepath: filepath, section: section, data: make(map[string]interface{})};

    key_names:=section.KeyStrings()
    for _, key_name  := range key_names {
      new_block.data[key_name]=section.Key(key_name).Value();
    }
    return new_block;
}


func (block *BlockT) Section() *ini.Section {
    return block.section;
}


func (block *BlockT) Data() map[string]interface{} {
    return block.data;
}

func (block *BlockT) Lookup(keyname string) string {
    v, ok := block.data[keyname]
    if (!ok) {
        fmt.Println("Block=", block.Name(), " Didn't find key:", keyname);
        return "";
    }
    return v.(string);
}


func (block *BlockT) Generate_block_html() string {
    section := block.Section()
    // TODO: Generate html.
    fmt.Println("Generating text for block:", block.Name());

    mdtext:=section.Key("Text").String();

    maybeUnsafeHTML := markdown.ToHTML([]byte(mdtext), nil, nil)
    html := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)
    return string(html)
}


func (block *BlockT) AddBlock(path string) {
    fmt.Println("Read block: ", block.Name());
}

func ini_file_walker(root string, process_block func(string, *BlockT)) {
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Println(err)
            return nil
        }
        if !info.IsDir() && filepath.Ext(path) == ".ini" {
            cfg, err := ini.Load(path);
            if err != nil {
                fmt.Printf("Fail to read file: %v", err)
                os.Exit(1)
            }
            for _, section := range cfg.Sections() {
                block:=CreateBlock(section.Name(), path, section); 
                ContentBlocks = append(ContentBlocks, block);
                process_block(path, block);
            }
        }
        return nil
    })

    if err != nil {
       log.Fatal(err)
    }
}


func FindBlock(name string) *BlockT {

    var matched_blocks []*BlockT;

    for _, block := range ContentBlocks {
        if (block.name == name) {
           matched_blocks = append(matched_blocks, block);
        }
    }
    if len(matched_blocks) == 0 {
       fmt.Println("No content blocks matches for referenced block in site_layout.html: ", name);
       return nil; 

    } else if len(matched_blocks) > 1 {
       fmt.Println("Multiple blocks match for block: ", name, " Use FilePath.block to refer to such blocks. Using first matching block.");
       return nil;
    }

    block := matched_blocks[0];
    return block;
}

func VisitBlock(filepath string, block *BlockT) {
   fmt.Println("Path=", filepath, " Created Block:", block.Name());
}

func Init(rootpath string) {
    ContentBlocks = []*BlockT{};
    ini_file_walker(rootpath, VisitBlock);
}


