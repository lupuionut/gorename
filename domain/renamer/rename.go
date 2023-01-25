package renamer

import (
    "sync"
    "os"
    "path"
    "strings"
)

func Iterate(files []string, rules map[string]string) {
    var wg sync.WaitGroup
    for _, file := range files {
        wg.Add(1)
        go ProcessFile(file, rules, &wg)
    }
    wg.Wait()
}

func ProcessFile(file string, rules map[string]string, wg *sync.WaitGroup) {
    defer wg.Done()
    directory, filename :=  path.Split(file)
    new_filename, err := Translate(filename, rules)
    if err != nil {
        return
    }
    if filename != new_filename {
        os.Rename(directory + filename, directory + new_filename)
    }
}

func Translate(title string, rules map[string]string) (string, error) {
    if len(rules) == 0 {
        return title, nil
    }
    for key, value := range rules {
       title = strings.ReplaceAll(title, key, value)
    }
    return title, nil
}
