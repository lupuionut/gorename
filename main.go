package main

import (
    "os"
    "fmt"
    "strings"
    "github.com/lupuionut/gorename/domain/cli"
    "github.com/lupuionut/gorename/domain/finder"
    "github.com/lupuionut/gorename/domain/rules"
)

func main() {
    args := os.Args
    defaults := make(map[string]string)
    defaults["target"] = ""
    defaults["rules"] = ""

    commands := &cli.Instance{ Commands: defaults }
    commands.ParseArgs(args)

    err, ok := commands.Error.(*cli.CommandParseError)
    if ok {
        fmt.Println(err)
        if err.Level == cli.Fatal {
            os.Exit(1)
        }
    }

    if commands.Commands["target"] == "" {
        fmt.Println("You must specify the path to the folder that contains the files you want to rename or directly to the file, e.g. -folder=/path/to/file.txt; -target=/path/to/folder")
        os.Exit(1)
    }
    if commands.Commands["rules"] == "" {
	fmt.Println("You must specify the file that contains the renaming rules. e.g. -rules=/path/to/rules.rule")
    	os.Exit(1)
    }

    // if rules cannot be read, exit
    rulesContent, errr := os.ReadFile(commands.Commands["rules"])
    if errr != nil {
        fmt.Println(errr)
        os.Exit(1)
    }

    if !strings.HasSuffix(commands.Commands["target"], "/") {
        commands.Commands["target"] += "/"
    }
    searcher := &finder.Finder{
        Cli: commands,
        Path: commands.Commands["target"],
        Recursive: false,
    }
    files, errf := searcher.FindItems();
    if errf != nil {
        fmt.Println(errf)
    }
    fmt.Println(files)

    rulesLines := strings.Split(string(rulesContent), "\n")
    parser := rules.Parser {
        Content: rulesLines,
    }

    // if an error occurs in parsing the rules, exit
    errp := parser.Parse()
    if errp != nil {
        fmt.Println(errp)
        os.Exit(1)
    }

    for _, t := range(parser.Tokens[0]) {
        fmt.Printf("%#v \n", t)
    }
    replacements := make(map[string]string)
    for _, line := range(parser.Tokens) {
        var k string

        if rules.IsValid(line) {
            for _, t := range(line) {
                if t.Type == rules.TokenTagValue {
                    if len(k) == 0 {
                        k = t.Value
                        replacements[k] = ""
                    } else{
                       replacements[k] = t.Value
                    }
                }
            }
        }
    }

    fmt.Println(replacements)
}
