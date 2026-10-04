package output
import "strings"
func Sensitive(k string) bool {
 k=strings.ToLower(strings.ReplaceAll(k,"-","_"))
 return k=="value" || k=="key" || k=="api_key" || k=="secret" || k=="password" || k=="csrf" || k=="token" || strings.HasSuffix(k,"_secret") || strings.HasSuffix(k,"_token") || k=="authorization" || k=="cookie"
}
