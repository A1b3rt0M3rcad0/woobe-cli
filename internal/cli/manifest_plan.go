package cli
import (
 "github.com/spf13/cobra"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/manifest"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)
func(a *App) manifestPlan(d manifest.Document)([]map[string]any,error){if e:=a.validateManifest(d);e!=nil{return nil,e};steps,e:=d.Order();if e!=nil{return nil,e};out:=[]map[string]any{};for _,s:=range steps{for _,op:=range a.Registry{if op.Command==s.Command{out=append(out,map[string]any{"step":s,"operation":op,"authorization":"not_evaluated","status":"planned"})}}};return out,nil}
func(a *App) manifestPlanCommands(g *cobra.Command){g.AddCommand(&cobra.Command{Use:"plan",Args:cobra.NoArgs,RunE:func(*cobra.Command,[]string)error{d,e:=a.readManifest();if e!=nil{return e};plan,e:=a.manifestPlan(d);if e!=nil{return e};return a.emit(map[string]any{"manifest_hash":d.Hash(),"workspace_id":d.Workspace,"project_id":d.Project,"operations":plan,"complete":true,"atomic":false})}})
 g.AddCommand(&cobra.Command{Use:"diff",Args:cobra.NoArgs,RunE:func(*cobra.Command,[]string)error{return output.New(9,"semantic remote diff requires resource schemas and server revision contracts; use plan for explicit operations")}})
 a.manifestApplyCommand(g)
}
