package credentials
import (
 "os"
 "path/filepath"
 "regexp"
 "strings"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/config"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)
var valid=regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,100}$`)
type Store struct {Dir string}
func(s Store) path(name string)(string,error){if !valid.MatchString(name){return "",output.New(2,"invalid credential reference")};return filepath.Join(s.Dir,name+".json"),nil}
func(s Store) Get(name string)(string,error){p,e:=s.path(name);if e!=nil{return "",e};b,e:=os.ReadFile(p);if e!=nil{return "",output.New(3,"credential unavailable")};return strings.TrimSpace(string(b)),nil}
func(s Store) Put(name,value string)error{p,e:=s.path(name);if e!=nil{return e};if strings.TrimSpace(value)==""{return output.New(2,"empty credential")};return config.AtomicWrite(p,[]byte(value),0600)}
func(s Store) Remove(name string)error{p,e:=s.path(name);if e!=nil{return e};return os.Remove(p)}
func(s Store) List()([]string,error){es,e:=os.ReadDir(s.Dir);if os.IsNotExist(e){return []string{},nil};if e!=nil{return nil,e};out:=[]string{};for _,v:=range es{if !v.IsDir()&&strings.HasSuffix(v.Name(),".json"){out=append(out,strings.TrimSuffix(v.Name(),".json"))}};return out,nil}
