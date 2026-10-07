//go:build windows

package packagecheckpoint

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"unsafe"
)

func supportedProtection() error        { return nil }
func privateMode(info os.FileInfo) bool { return info.Mode().IsRegular() }

// Windows privacy is an ACL property, never inferred from Unix mode bits.
// Only the current account and Windows administrative principals may have access.
func checkPrivateFile(file *os.File) error {
	descriptor, err := windows.GetSecurityInfo(windows.Handle(file.Fd()), windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION|windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return output.New(3, "Package operational ACL could not be verified")
	}
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return output.New(3, "Windows account identity is unavailable")
	}
	owner, _, err := descriptor.Owner()
	if err != nil || owner == nil || (owner.String() != user.User.Sid.String() && owner.String() != "S-1-5-32-544") {
		return output.New(3, "Package operational file must belong to the current Windows account")
	}
	acl, _, err := descriptor.DACL()
	if err != nil || acl == nil {
		return output.New(3, "Package operational file requires a private Windows DACL")
	}
	for index := uint32(0); index < uint32(acl.AceCount); index++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err = windows.GetAce(acl, index, &ace); err != nil {
			return output.New(3, "Package operational ACL could not be inspected")
		}
		if ace.Header.AceType == windows.ACCESS_DENIED_ACE_TYPE {
			continue
		}
		if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
			return output.New(3, "Package operational ACL modality is unsupported")
		}
		sid := (*windows.SID)(unsafe.Pointer(&ace.SidStart)).String()
		if ace.Mask != 0 && sid != user.User.Sid.String() && sid != "S-1-5-18" && sid != "S-1-5-32-544" && !(sid == "S-1-3-0" && ace.Header.AceFlags&0x08 != 0) {
			return output.New(3, "Package operational directory must have a private Windows DACL")
		}
	}
	return nil
}
func checkPrivateParent(root *os.Root) error {
	parent, err := root.Open(".")
	if err != nil {
		return err
	}
	defer parent.Close()
	return checkPrivateFile(parent)
}
func syncParent(file *os.File) error { return nil } // Windows has no directory fsync; file data is flushed before atomic rename.

func openConfinedRead(root *os.Root, name string) (*os.File, error) {
	info, err := root.Lstat(name)
	if err != nil || !info.Mode().IsRegular() {
		if err != nil {
			return nil, err
		}
		return nil, output.New(3, "Package operational file must be regular")
	}
	file, err := root.Open(name)
	if err != nil {
		return nil, err
	}
	opened, err := file.Stat()
	if err != nil || !os.SameFile(info, opened) || !opened.Mode().IsRegular() {
		file.Close()
		return nil, output.New(3, "Package operational file changed while opening")
	}
	return file, nil
}
func openPrivateRead(path string) (*os.File, error) {
	root, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return openConfinedRead(root, filepath.Base(path))
}
func lockCheckpoint(root *os.Root, name string) (*os.File, error) {
	if info, err := root.Lstat(name); err == nil && !info.Mode().IsRegular() {
		return nil, output.New(3, "Package lock must be regular")
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	file, err := root.OpenFile(name, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	info, err := root.Lstat(name)
	opened, statErr := file.Stat()
	if err != nil || statErr != nil || !info.Mode().IsRegular() || !os.SameFile(info, opened) {
		file.Close()
		return nil, output.New(3, "Package lock changed while opening")
	}
	if err = checkPrivateFile(file); err != nil {
		file.Close()
		return nil, err
	}
	var overlapped windows.Overlapped
	if err = windows.LockFileEx(windows.Handle(file.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &overlapped); err != nil {
		file.Close()
		return nil, output.New(2, "Package checkpoint is in use by another command")
	}
	return file, nil
}
