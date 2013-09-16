package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "os"
import "os/exec"
import "flag"
import "path/filepath"
import "strings"

var install_prefix = ""
var scratch_path string

func compile (code string, options string) (success bool, buildOutput string, listing string, err error) {
    if file, err := ioutil.TempFile ("", "tryout-"); err == nil {
        defer file.Close()
        defer os.Remove(file.Name())
        if _, err = file.WriteString (code); err == nil {
            output_file := scratch_path + "/" + filepath.Base(file.Name()) + ".out"

            split_options := []string {}
            if len(options) != 0 {
                split_options = strings.Split(options, " ")
            }
            split_options = append (split_options, "-xc")
            split_options = append (split_options, file.Name())
            split_options = append (split_options, "-c")
            split_options = append (split_options, "-o")
            split_options = append (split_options, output_file)
            cmd := exec.Command (install_prefix + "/avr-gcc", split_options...)

            if output, err := cmd.CombinedOutput (); err == nil {
                defer os.Remove(output_file)
                buildOutput = string(output)
                success = cmd.ProcessState.Success()
                cmd := exec.Command (install_prefix + "/avr-objdump", "-S", output_file)
                if output, err = cmd.CombinedOutput (); err == nil {
                    listing = string(output)
                }
            }
        }
    }
    return
}    

func buildHandler (w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic (err)
    }
    type JsonData struct {
        Code string
        Options string
    }
    type JsonRetData struct {
        BuildOutput string
        Success bool
        Listing string 
    }
    
    var jsonData JsonData
    json.Unmarshal (body, &jsonData)
    
    defer r.Body.Close()
    
    var ret JsonRetData
    ret.Success, ret.BuildOutput, ret.Listing, err = compile (jsonData.Code, jsonData.Options)

    if err != nil {
       panic (err)
    }
    
    w.Header ().Add ("Content-Type", "application/json")
    data, err := json.Marshal (ret)

    fmt.Fprintf (w, "%s", data)
}

func initWebServer() {
    http.HandleFunc ("/compile", buildHandler)
    http.Handle ("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("/home/saaadhu/code/git/tryout/src/tryout/www"))))
    http.ListenAndServe(":8082", nil)
}

func main() {
    prefix := flag.String("prefix", "", "Provide prefix path (including bin) for avr-gcc")
    scratch := flag.String("scratch", "", "Path to temporarily store compiled output")
    flag.Parse()
    install_prefix = *prefix;
    scratch_path = *scratch;
    initWebServer()
}

