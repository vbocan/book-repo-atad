// Command contact-manager is a reference solution for Exercise 5.1. It stores
// contacts as a JSON array in a file and supports adding, listing, searching,
// and deleting them from the command line.
//
// Usage:
//
//	contact-manager add <name> <email> <phone>
//	contact-manager list
//	contact-manager search <query>
//	contact-manager delete <name>
//
// Contacts persist in contacts.json in the current working directory,
// created automatically on the first "add".
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// contactsFile is where the contact list is persisted, as a JSON array.
const contactsFile = "contacts.json"

// Contact is one address-book entry. The json tags keep the JSON keys
// lowercase while the Go fields stay idiomatic PascalCase.
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 5 {
			usage()
		}
		addContact(os.Args[2], os.Args[3], os.Args[4])
	case "list":
		listContacts()
	case "search":
		if len(os.Args) != 3 {
			usage()
		}
		findContacts(os.Args[2])
	case "delete":
		if len(os.Args) != 3 {
			usage()
		}
		deleteContact(os.Args[2])
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: contact-manager <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  add <name> <email> <phone>   Add a new contact")
	fmt.Println("  list                         List all contacts")
	fmt.Println("  search <query>                Search by name or email")
	fmt.Println("  delete <name>                Delete a contact by name")
	os.Exit(1)
}

// loadContacts reads the file at path and decodes it into a slice of
// Contact. A missing file is not an error: it simply means there are no
// contacts yet, so loadContacts returns a nil slice and a nil error, and the
// first "add" works without any setup step.
func loadContacts(path string) ([]Contact, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var contacts []Contact
	if err := json.Unmarshal(data, &contacts); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return contacts, nil
}

// saveContacts encodes contacts as indented JSON and writes it to path,
// overwriting whatever was there before.
func saveContacts(path string, contacts []Contact) error {
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding contacts: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// searchContacts returns every contact whose name or email contains query,
// matched case-insensitively. It does no I/O, which keeps it easy to test.
func searchContacts(contacts []Contact, query string) []Contact {
	query = strings.ToLower(query)
	matches := make([]Contact, 0)
	for _, c := range contacts {
		if strings.Contains(strings.ToLower(c.Name), query) ||
			strings.Contains(strings.ToLower(c.Email), query) {
			matches = append(matches, c)
		}
	}
	return matches
}

// load and save wrap loadContacts and saveContacts for the command-line
// handlers: any error is reported and ends the program.
func load() []Contact {
	contacts, err := loadContacts(contactsFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	return contacts
}

func save(contacts []Contact) {
	if err := saveContacts(contactsFile, contacts); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printContacts(contacts []Contact) {
	for _, c := range contacts {
		fmt.Printf("%s | %s | %s\n", c.Name, c.Email, c.Phone)
	}
}

func addContact(name, email, phone string) {
	contacts := load()
	contacts = append(contacts, Contact{Name: name, Email: email, Phone: phone})
	save(contacts)
	fmt.Printf("Added %s\n", name)
}

func listContacts() {
	contacts := load()
	if len(contacts) == 0 {
		fmt.Println("No contacts.")
		return
	}
	printContacts(contacts)
}

// findContacts prints the result of searchContacts for the "search" command.
func findContacts(query string) {
	matches := searchContacts(load(), query)
	if len(matches) == 0 {
		fmt.Println("No matching contacts.")
		return
	}
	printContacts(matches)
}

// deleteContact removes every contact with an exact name match. Using
// make(..., 0, ...) instead of a nil slice ensures that deleting every
// contact writes "[]" to contactsFile rather than "null".
func deleteContact(name string) {
	contacts := load()

	kept := make([]Contact, 0, len(contacts))
	deleted := false
	for _, c := range contacts {
		if c.Name == name {
			deleted = true
			continue
		}
		kept = append(kept, c)
	}

	if !deleted {
		fmt.Printf("No contact named %s.\n", name)
		return
	}
	save(kept)
	fmt.Printf("Deleted %s\n", name)
}
