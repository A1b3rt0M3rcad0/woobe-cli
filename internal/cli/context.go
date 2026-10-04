package cli
import (
 "github.com/spf13/cobra"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/config"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)
func(a *App) contextCommands(){
 g:=&cobra.Command{Use:"context",Short:"Local context profiles"};a.Root.AddCommand(g)
 for _,action:=range []string{"list","show","create","update","use","delete","set"}{action:=action;cmd:=&cobra.Command{Use:action,RunE:func(_ *cobra.Command,args []string)error{
 c,e:=config.Load(a.ConfigPath);if e!=nil{return output.New(2,"cannot load config")};name:=c.Current;if a.ContextName!=""{name=a.ContextName};if len(args)>1{return output.New(2,"expected at most one context name")};if len(args)==1{name=args[0]}
 switch action {
 case "list":return a.emit(c)
 case "show":v,e:=a.resolve();if e!=nil{return e};return a.emit(v)
 case "create":if len(args)!=1{return output.New(2,"context name required")};if _,ok:=c.Contexts[name];ok{return output.New(6,"context exists")};if a.APIURL==""{return output.New(2,"--api-url required")};c.Contexts[name]=config.Context{APIURL:a.APIURL,Workspace:a.Workspace,Project:a.Project}
 default:v,ok:=c.Contexts[name];if !ok{return output.New(2,"unknown context")};switch action{case "use":c.Current=name;case "delete":delete(c.Contexts,name);if c.Current==name{c.Current=""};case "update","set":if a.APIURL!=""{v.APIURL=a.APIURL};if a.Workspace!=""&&a.Workspace!=v.Workspace{v.Workspace=a.Workspace;v.Project=""};if a.Project!=""{v.Project=a.Project};c.Contexts[name]=v}
 }
 if e=config.Save(a.ConfigPath,c);e!=nil{return e};return a.emit(map[string]string{"context":name,"action":action})
 }};g.AddCommand(cmd)}
 cr:=&cobra.Command{Use:"credential"};g.AddCommand(cr)
 for _,action:=range []string{"attach","detach"}{action:=action;cr.AddCommand(&cobra.Command{Use:action+" <context>",Args:cobra.ExactArgs(1),RunE:func(_ *cobra.Command,args []string)error{c,e:=config.Load(a.ConfigPath);if e!=nil{return e};v,ok:=c.Contexts[args[0]];if !ok{return output.New(2,"unknown context")};if action=="attach"{if a.Credential==""{return output.New(2,"--credential required")};if _,e=a.store().Get(a.Credential);e!=nil{return e};v.Credential=a.Credential}else{v.Credential=""};c.Contexts[args[0]]=v;if e=config.Save(a.ConfigPath,c);e!=nil{return e};return a.emit(map[string]string{"context":args[0],"action":action})}})}
}
