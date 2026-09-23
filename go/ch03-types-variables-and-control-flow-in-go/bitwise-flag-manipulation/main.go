package main

import (
	"fmt"
	"strings"
)

// Permission represents a set of access rights as a bitmask.
type Permission uint8

const (
	PermRead Permission = 1 << iota
	PermWrite
	PermExecute
)

// Set returns the permission set with flag added.
func (p Permission) Set(flag Permission) Permission {
	return p | flag
}

// Clear returns the permission set with flag removed.
func (p Permission) Clear(flag Permission) Permission {
	return p &^ flag
}

// Has reports whether flag is present in the permission set.
func (p Permission) Has(flag Permission) bool {
	return p&flag != 0
}

// String returns the pipe-joined names of the set flags, e.g. "Read | Write".
func (p Permission) String() string {
	var names []string
	if p.Has(PermRead) {
		names = append(names, "Read")
	}
	if p.Has(PermWrite) {
		names = append(names, "Write")
	}
	if p.Has(PermExecute) {
		names = append(names, "Execute")
	}
	if len(names) == 0 {
		return "None"
	}
	return strings.Join(names, " | ")
}

func main() {
	perms := PermRead.Set(PermWrite)
	fmt.Printf("Permissions: %s (0b%03b)\n", perms, perms)

	perms = perms.Set(PermExecute)
	fmt.Printf("After adding Execute: %s (0b%03b)\n", perms, perms)

	perms = perms.Clear(PermWrite)
	fmt.Printf("After removing Write: %s (0b%03b)\n", perms, perms)

	fmt.Printf("Has Read? %t\n", perms.Has(PermRead))
	fmt.Printf("Has Write? %t\n", perms.Has(PermWrite))
}
