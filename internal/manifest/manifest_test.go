package manifest
import "testing"
func TestOrderAndCycles(t *testing.T){d:=Document{Steps:[]Step{{ID:"network",Command:"x",DependsOn:[]string{"agent"}},{ID:"agent",Command:"x"}}};steps,e:=d.Order();if e!=nil||steps[0].ID!="agent"{t.Fatal(steps,e)};d.Steps[1].DependsOn=[]string{"network"};if _,e=d.Order();e==nil{t.Fatal("cycle accepted")};d.Steps[1].DependsOn=[]string{"missing"};if _,e=d.Order();e==nil{t.Fatal("missing dependency accepted")}}
func TestStrictDocument(t *testing.T){for _,b:=range []string{`{"schema_version":"2","steps":[]}`,`{"schema_version":"1","steps":[],"typo":true}`,`{} {}`,`{"schema_version":"1","steps":[{"id":"a","command":"x"},{"id":"a","command":"y"}]}`}{if _,e:=Parse([]byte(b));e==nil{t.Fatal(b)}}}
