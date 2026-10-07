package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWindowsProtectionVaultPersistenceAndIsolation(t *testing.T) {
	if !nativeStore {
		t.Skip("Windows native protection")
	}
	store := Store{Dir: t.TempDir()}
	if e := store.Put("private", "fixture"); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { _ = store.Remove("private") })
	restored := Store{Dir: store.Dir}
	v, e := restored.Get("private")
	if e != nil || v != "fixture" {
		t.Fatal("credential did not persist", e)
	}
	if _, e = os.Stat(filepath.Join(store.Dir, "private.json")); !os.IsNotExist(e) {
		t.Fatal("secret written to disk")
	}
	if _, e = (Store{Dir: t.TempDir()}).Get("private"); e == nil {
		t.Fatal("credential leaked across configurations")
	}
	names, e := restored.List()
	if e != nil || len(names) != 1 || names[0] != "private" {
		t.Fatal(names, e)
	}
	if e = restored.Remove("private"); e != nil {
		t.Fatal(e)
	}
	if _, e = restored.Get("private"); e == nil {
		t.Fatal("credential survived logout")
	}
}
