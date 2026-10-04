package cli
import (
 "context"
 "encoding/json"
 "errors"
 "os"
 "github.com/spf13/cobra"
 sdk "github.com/A1b3rt0M3rcad0/woobe-sdk-go"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)
type runtimeTarget interface {
 Run(context.Context,string,*sdk.RunOptions)(*sdk.RunResult,error)
 Chat(string,*sdk.ChatOptions)(*sdk.Chat,error)
 ObserveRun(context.Context,string)(*sdk.Stream,error)
 ActiveRun(context.Context,string)(*sdk.ActiveRun,error)
 CancelRun(context.Context,string)error
}
func runtimeError(e error)error{if e==nil{return nil};var r *sdk.RequestError;if errors.As(e,&r){code:=7;switch r.StatusCode{case 400,422:code=2;case 401:code=3;case 403:code=4;case 404:code=5;case 409:code=6};return &output.Error{Code:code,Message:"runtime request failed",Status:r.StatusCode}};return output.New(7,"runtime transport or protocol failed")}
func(a *App) runtimeCommands(){
 for _,op:=range []Operation{
 {Command:"runtime run list",Method:"GET",Path:"/runtime/runs",QueryScope:"project_id"},{Command:"runtime run get",Method:"GET",Path:"/runtime/runs/{run_id}"},{Command:"runtime run cancel",Method:"POST",Path:"/runtime/runs/{run_id}/cancel",Effect:"execution"},{Command:"runtime run rerun",Method:"POST",Path:"/runtime/runs/{run_id}/rerun",Effect:"execution",Body:true},{Command:"runtime run invocations",Method:"GET",Path:"/runtime/runs/{run_id}/model-invocations"},
 {Command:"project trace list",Method:"GET",Path:"/observability/traces",QueryScope:"project_id"},{Command:"project trace get",Method:"GET",Path:"/observability/traces/{trace_id}"},{Command:"project usage summary",Method:"GET",Path:"/metering/projects/{project_id}/usage"},{Command:"project usage daily",Method:"GET",Path:"/metering/projects/{project_id}/daily"},{Command:"project agent usage daily",Method:"GET",Path:"/metering/agents/{agent_id}/daily"},
 {Command:"project network execution diagnostics",Method:"GET",Path:"/network/executions/{execution_id}/diagnostics"},{Command:"project network execution snapshot",Method:"GET",Path:"/network/executions/{execution_id}/snapshot"},{Command:"project network execution events",Method:"GET",Path:"/network/executions/{execution_id}/events"},{Command:"project network execution cancel",Method:"POST",Path:"/network/executions/{execution_id}/cancel",Effect:"execution"},
 {Command:"runtime session messages",Method:"GET",Path:"/runtime/sessions/{session_id}/messages"},{Command:"runtime session create",Method:"POST",Path:"/runtime/sessions",Body:true},
 }{a.register(op)}
 g:=a.group("runtime target");var kind string
 for _,action:=range []string{"run","stream","observe","active","cancel"}{action:=action;n:=1;if action=="observe"||action=="active"||action=="cancel"{n=2};c:=&cobra.Command{Use:action+" <alias> [run-or-session-id]",Args:cobra.ExactArgs(n),RunE:func(cmd *cobra.Command,args []string)error{
 if (action=="stream"||action=="observe")&&a.Mode!="jsonl"{return output.New(2,"stream and observe require --output jsonl")};v,e:=a.resolve();if e!=nil{return e};key:=os.Getenv("WOOBE_RUNTIME_KEY");if v.RuntimeCredential!=""{key,e=a.store().Get(v.RuntimeCredential);if e!=nil{return e}};if key==""{return output.New(3,"explicit runtime credential required")};if kind!="agent"&&kind!="network"{return output.New(2,"--target-kind must be agent or network")}
 cp,e:=a.client();if e!=nil{return e};client,e:=sdk.New(sdk.WithBaseURL(v.APIURL),sdk.WithHTTPClient(cp.HTTP),sdk.WithReconnect(0,0));if e!=nil{return runtimeError(e)};defer client.CloseIdleConnections();var t runtimeTarget
 if kind=="agent"{t,e=client.Connect.Agent(args[0],key)}else{t,e=client.Connect.Network(args[0],key)};if e!=nil{return runtimeError(e)}
 ctx,cancel:=context.WithTimeout(cmd.Context(),a.Timeout);defer cancel()
 switch action{
 case "active":v,e:=t.ActiveRun(ctx,args[1]);if e!=nil{return runtimeError(e)};return a.emit(v)
 case "cancel":if !a.Yes{return output.New(2,"--yes required")};if e=t.CancelRun(ctx,args[1]);e!=nil{return runtimeError(e)};return a.emit(map[string]bool{"cancel_requested":true})
 }
 var stream *sdk.Stream
 if action=="observe"{stream,e=t.ObserveRun(ctx,args[1])}else{
 if !a.Yes{return output.New(2,"execution requires --yes")};b,e:=a.body(true);if e!=nil{return e};var input struct{Input string `json:"input"`;Options sdk.RunOptions `json:"options"`};if e=json.Unmarshal(b,&input);e!=nil||input.Input==""{return output.New(2,"input JSON requires nonempty input")}
 if action=="run"{result,e:=t.Run(ctx,input.Input,&input.Options);if e!=nil{return runtimeError(e)};return a.emit(result)}
 chat,err:=t.Chat(input.Input,&input.Options);if err!=nil{return runtimeError(err)};stream,e=chat.Stream(ctx)
 }
 if e!=nil{return runtimeError(e)};if a.Mode!="jsonl"{return output.New(2,"stream and observe require --output jsonl")};enc:=json.NewEncoder(a.Out);for event:=range stream.Events{if e=enc.Encode(event);e!=nil{return e}};if e=stream.Err();e!=nil{if ctx.Err()!=nil{return ctx.Err()};return runtimeError(e)};return nil
 }};c.Flags().StringVar(&kind,"target-kind","agent","agent or network");g.AddCommand(c)}
}
