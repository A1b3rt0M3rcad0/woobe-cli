package cli
import (
 "context"
 "encoding/json"
 "io"
 "net/url"
 "os"
 "path/filepath"
 "strings"
 "time"
 "github.com/spf13/cobra"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/identity"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/config"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/controlplane"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/credentials"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)
var Version="dev"
type App struct {
 Root *cobra.Command
 In io.Reader; Out,Err io.Writer
 ConfigPath,ContextName,APIURL,Workspace,Project,Credential,RuntimeCredential,Mode string
 Timeout time.Duration; Yes,DryRun,NoInput bool
 File,IfMatch,IdempotencyKey,SecretFile string
 Query []string
 Registry []Operation
}
func New(in io.Reader,out,errOut io.Writer)*App {
 a:=&App{In:in,Out:out,Err:errOut}
 r:=&cobra.Command{Use:"woobe",Short:"Woobe public control plane and runtime",SilenceErrors:true,SilenceUsage:true}
 a.Root=r;r.SetIn(in);r.SetOut(out);r.SetErr(errOut)
 f:=r.PersistentFlags();f.StringVar(&a.ConfigPath,"config",config.DefaultPath(),"Private config file");f.StringVar(&a.ContextName,"context","","Context name")
 f.StringVar(&a.APIURL,"api-url","","API origin");f.StringVar(&a.Workspace,"workspace","","Workspace ID");f.StringVar(&a.Project,"project","","Project ID")
 f.StringVar(&a.Credential,"credential","","Administrative credential reference");f.StringVar(&a.RuntimeCredential,"runtime-credential","","Runtime credential reference")
 f.StringVar(&a.Mode,"output","json","json, jsonl or table (indented JSON)");f.DurationVar(&a.Timeout,"timeout",30*time.Second,"HTTP deadline")
 f.BoolVar(&a.Yes,"yes",false,"Accept the specified destructive operation");f.BoolVar(&a.NoInput,"no-input",false,"Deterministic execution (always enabled)");f.BoolVar(&a.DryRun,"dry-run",false,"Render the request without executing")
 f.StringVar(&a.File,"file","","JSON input file, or - for stdin");f.StringArrayVar(&a.Query,"query",nil,"Query name=value (repeatable)")
 f.StringVar(&a.IfMatch,"if-match","","Expected server ETag");f.StringVar(&a.IdempotencyKey,"idempotency-key","","Key, only when supported by server")
 f.StringVar(&a.SecretFile,"secret-file","","Exclusive private destination for issued secret")
 r.PersistentPreRunE=func(*cobra.Command,[]string)error{if a.Mode!="json"&&a.Mode!="jsonl"&&a.Mode!="table"{return output.New(2,"invalid output mode")};if a.Timeout<=0{return output.New(2,"timeout must be positive")};return nil}
 r.AddCommand(&cobra.Command{Use:"version",Args:cobra.NoArgs,RunE:func(*cobra.Command,[]string)error{return a.emit(map[string]any{"version":Version,"schema_version":"1"})}})
 a.contextCommands();a.requestCommands();a.workspaceCommands();a.projectCommands();a.agentCommands();a.releaseCommands();a.networkCommands();a.toolCommands();a.knowledgeCommands();a.catalogCommands();a.surfaceCommands();a.keysCommands();a.permissionCommands();a.authCommands();a.runtimeCommands();a.manifestCommands();a.discoveryCommands();return a
}
func(a *App) emit(v any)error{return output.Write(a.Out,a.Mode,output.Redact(v),map[string]string{"workspace_id":a.Workspace,"project_id":a.Project},nil)}
func(a *App) store()credentials.Store{return credentials.Store{Dir:filepath.Join(filepath.Dir(a.ConfigPath),"credentials")}}
func(a *App) resolve()(config.Context,error){
 c,e:=config.Load(a.ConfigPath);if e!=nil{return config.Context{},output.New(2,"cannot load config")};name:=a.ContextName;if name==""{name=os.Getenv("WOOBE_CONTEXT")};if name==""{name=c.Current}
 v:=config.Context{};if name!=""{var ok bool;v,ok=c.Contexts[name];if !ok{return v,output.New(2,"unknown context")}}
 vals:=[]struct{dst *string;flag,env string}{{&v.APIURL,a.APIURL,"WOOBE_API_URL"},{&v.Workspace,a.Workspace,"WOOBE_WORKSPACE_ID"},{&v.Project,a.Project,"WOOBE_PROJECT_ID"},{&v.Credential,a.Credential,"WOOBE_CREDENTIAL"},{&v.RuntimeCredential,a.RuntimeCredential,"WOOBE_RUNTIME_CREDENTIAL"}}
 for _,x:=range vals{if x.flag!=""{*x.dst=x.flag}else if s:=os.Getenv(x.env);s!=""{*x.dst=s}}
 if v.APIURL==""{v.APIURL="http://localhost:8000"};a.Workspace=v.Workspace;a.Project=v.Project;return v,nil
}
func(a *App) client()( *controlplane.Client,error){v,e:=a.resolve();if e!=nil{return nil,e};key:=os.Getenv("WOOBE_CONTROL_KEY");if v.Credential!=""{key,e=a.store().Get(v.Credential);if e!=nil{return nil,e}};c,e:=controlplane.New(v.APIURL,key,a.Timeout);if e==nil{if key==""{session,jar,err:=identity.Load(filepath.Join(filepath.Dir(a.ConfigPath),"sessions"),v.APIURL);if err!=nil{return nil,output.New(3,"cannot restore session")};c.HTTP.Jar=jar;c.ResponseHook=session.Capture;c.Headers.Set("X-CSRF-Token",session.CSRF)};if a.IfMatch!=""{c.Headers.Set("If-Match",a.IfMatch)};if a.IdempotencyKey!=""{c.Headers.Set("Idempotency-Key",a.IdempotencyKey)}};return c,e}
func(a *App) body(required bool)([]byte,error){if a.File==""{if required{return nil,output.New(2,"--file is required (use - for stdin)")};return nil,nil};var b []byte;var e error;if a.File=="-"{b,e=io.ReadAll(io.LimitReader(a.In,8<<20+1))}else{b,e=os.ReadFile(a.File)};if e!=nil{return nil,output.New(2,"cannot read input")};if len(b)>8<<20||!json.Valid(b){return nil,output.New(2,"input must be JSON within 8 MiB")};return b,nil}
func(a *App) query()(url.Values,error){q:=url.Values{};for _,s:=range a.Query{k,v,ok:=strings.Cut(s,"=");if !ok||k==""{return nil,output.New(2,"query must be name=value")};q.Add(k,v)};return q,nil}
func(a *App) Execute(ctx context.Context,args []string)int{a.Root.SetArgs(args);e:=a.Root.ExecuteContext(ctx);if e==nil{return 0};e=output.Normalize(e);w:=a.Out;if a.Mode=="jsonl"{w=a.Err};_ = output.Write(w,a.Mode,nil,map[string]string{"workspace_id":a.Workspace,"project_id":a.Project},e);return output.Normalize(e).Code}
