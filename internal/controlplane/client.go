package controlplane
import (
 "bytes"
 "context"
 "encoding/json"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"
 "github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)
type Client struct {Base string; Token string; HTTP *http.Client; Headers http.Header; ResponseHook func(*http.Response)error}
func New(base,token string,timeout time.Duration)(*Client,error){
 u,e:=url.Parse(base);if e!=nil||u.Host==""||(u.Scheme!="http"&&u.Scheme!="https")||u.User!=nil||u.RawQuery!=""||u.Fragment!="" {return nil,output.New(2,"api URL must be an HTTP(S) origin without credentials, query or fragment")}
 h:=&http.Client{Timeout:timeout,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 return &Client{Base:strings.TrimRight(base,"/"),Token:token,HTTP:h,Headers:make(http.Header)},nil
}
func(c *Client) Request(ctx context.Context,method,path string,q url.Values,body []byte)(any,http.Header,error){
 if !strings.HasPrefix(path,"/")||strings.HasPrefix(path,"//")||strings.ContainsAny(path,"?#\\") {return nil,nil,output.New(2,"request path must be an origin-relative path; use --query for query parameters")}
 raw:=c.Base+path;if len(q)>0 {raw+="?"+q.Encode()}
 req,e:=http.NewRequestWithContext(ctx,method,raw,bytes.NewReader(body));if e!=nil{return nil,nil,output.New(2,"invalid request")}
 req.Header=c.Headers.Clone();req.Header.Set("Accept","application/json");if c.Token!=""{req.Header.Set("Authorization","Bearer "+c.Token)};if body!=nil{req.Header.Set("Content-Type","application/json")}
 resp,e:=c.HTTP.Do(req);if e!=nil {err:=output.Normalize(e);if err.Code==1{err.Code=7;err.Message="HTTP connection failed"};if method!="GET"&&method!="HEAD"{err.Outcome="unknown"};return nil,nil,err};defer resp.Body.Close();if c.ResponseHook!=nil {if e=c.ResponseHook(resp);e!=nil{return nil,resp.Header,&output.Error{Code:10,Message:"server responded but session persistence failed",Outcome:"unknown"}}}
 b,e:=io.ReadAll(io.LimitReader(resp.Body,32<<20+1));if e!=nil{return nil,resp.Header,output.New(7,"response could not be read")};if len(b)>32<<20{return nil,resp.Header,output.New(9,"response exceeds limit")}
 if resp.StatusCode<200||resp.StatusCode>=300{
  code:=7;switch resp.StatusCode{case 400,422:code=2;case 401:code=3;case 403:code=4;case 404:code=5;case 409,412:code=6;case 408,504:code=8;case 405,501:code=9}
  outcome:="rejected";if resp.StatusCode>=500 {outcome="unknown"}
  return nil,resp.Header,&output.Error{Code:code,Message:"server rejected request",Status:resp.StatusCode,RequestID:resp.Header.Get("X-Request-ID"),Outcome:outcome}
 }
 if len(bytes.TrimSpace(b))==0{return nil,resp.Header,nil};var v any;d:=json.NewDecoder(bytes.NewReader(b));d.UseNumber();if e=d.Decode(&v);e!=nil{return nil,resp.Header,output.New(9,"server returned non-JSON response")};return v,resp.Header,nil
}
