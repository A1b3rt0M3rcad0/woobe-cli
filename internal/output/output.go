package output

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "io"
)

type Error struct {
 Code int `json:"exit_code"`
 Message string `json:"message"`
 Status int `json:"http_status,omitempty"`
 RequestID string `json:"request_id,omitempty"`
 Outcome string `json:"write_outcome,omitempty"`
 Cause error `json:"-"`
}
func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Cause }
func New(code int, message string) *Error { return &Error{Code:code,Message:message} }
func Normalize(err error) *Error {
 var e *Error
 if errors.As(err,&e) {return e}
 if errors.Is(err,context.Canceled) {return New(130,"local operation interrupted")}
 if errors.Is(err,context.DeadlineExceeded) {return New(8,"deadline exceeded")}
 return &Error{Code:1,Message:err.Error(),Cause:err}
}
type Envelope struct {
 SchemaVersion string `json:"schema_version"`
 Success bool `json:"success"`
 Context map[string]string `json:"context,omitempty"`
 Data any `json:"data,omitempty"`
 Error *Error `json:"error,omitempty"`
 Meta map[string]any `json:"meta"`
}
func Write(w io.Writer, mode string, data any, scope map[string]string, err error) error {
 e:=Envelope{SchemaVersion:"1",Success:err==nil,Context:scope,Data:data,Meta:map[string]any{"complete":err==nil}}
 if err!=nil {e.Error=Normalize(err)}
 enc:=json.NewEncoder(w)
 if mode=="table" {enc.SetIndent("","  ")}
 if mode!="table" && mode!="json" && mode!="jsonl" {return New(2,fmt.Sprintf("unsupported output %q",mode))}
 return enc.Encode(e)
}
// Redact recursively protects secret-bearing fields even in generic responses.
func Redact(v any) any {
 switch x:=v.(type) {
 case map[string]any:
  out:=make(map[string]any,len(x)); for k,v:=range x {if Sensitive(k) {out[k]="[REDACTED]"} else {out[k]=Redact(v)}}; return out
 case []any: out:=make([]any,len(x)); for i,v:=range x {out[i]=Redact(v)}; return out
 default:return v
 }
}
