package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrivateReferences(t *testing.T) {
	s := Store{Dir: t.TempDir()}
	for _, n := range []string{"../other", "/absolute", "a/b", ""} {
		if s.Put(n, "secret") == nil {
			t.Fatal(n)
		}
	}
	if e := s.Put("fixture", "secret"); e != nil {
		t.Fatal(e)
	}
	info, e := os.Stat(filepath.Join(s.Dir, "fixture.json"))
	if e != nil || info.Mode().Perm() != 0600 {
		t.Fatal(info, e)
	}
	v, e := s.Get("fixture")
	if e != nil || v != "secret" {
		t.Fatal(v, e)
	}
}
func TestSymlinkAndPublicFilesRejected(t *testing.T) {
	dir := t.TempDir()
	s := Store{Dir: dir}
	p := filepath.Join(dir, "public.json")
	_ = os.WriteFile(p, []byte("secret"), 0644)
	if _, e := s.Get("public"); e == nil {
		t.Fatal("public credential accepted")
	}
	_ = os.Symlink(p, filepath.Join(dir, "link.json"))
	if _, e := s.Get("link"); e == nil {
		t.Fatal("symlink accepted")
	}
}
