package finder

import (
    "os"
    "strings"
    "errors"
    "github.com/lupuionut/gorename/domain/cli"
)

type Finder struct {
    Cli *cli.Instance
    Path string
    Recursive bool
}

func (ren *Finder) FindItems() ([]string, error) {
    var files []string
    var errText string
    stat, err := os.Stat(ren.Path)
    if err != nil {
        return nil, err
    }
    if !stat.IsDir() {
        return []string{ren.Path}, nil
    }

    items, err := os.ReadDir(ren.Path)
    if err != nil {
        return nil, err
    }

    for _, item := range(items) {
        location := ren.Path + item.Name()
        stat, err = os.Stat(location)
        if err != nil {
            errText += err.Error() + "\n"
            continue
        }
        if !stat.IsDir() {
           files = append(files, location)
        } else {
            if ren.Recursive == true {
                if !strings.HasSuffix(location, "/") {
                    location += "/"
                }
                NewFinder := &Finder {
                    Cli: ren.Cli,
                    Path: location,
                    Recursive: true,
                }
                newfiles, err := NewFinder.FindItems()
                if err != nil {
                    errText += err.Error() + "\n"
                    continue
                }
                files = append(files, newfiles...)
            }
        }
    }

    if errText == "" {
        return files, nil
    }
    return files, errors.New(errText)
}
