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
	"fmt"
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
		searchContacts(os.Args[2])
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

// load reads contactsFile and decodes it into a slice of Contact. A missing
// file is treated as an empty contact list rather than an error, so the
// first "add" works without any setup step.
func load() []Contact {
	data, err := os.ReadFile(contactsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Contact{}
		}
		fmt.Printf("Error reading %s: %v\n", contactsFile, err)
		os.Exit(1)
	}

	var contacts []Contact
	if err := json.Unmarshal(data, &contacts); err != nil {
		fmt.Printf("Error parsing %s: %v\n", contactsFile, err)
		os.Exit(1)
	}
	return contacts
}

// save encodes contacts as indented JSON and writes it back to
// contactsFile, overwriting whatever was there before.
func save(contacts []Contact) {
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding contacts: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(contactsFile, data, 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", contactsFile, err)
		os.Exit(1)
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
	for _, c := range contacts {
		fmt.Printf("%s | %s | %s\n", c.Name, c.Email, c.Phone)
	}
}

// searchContacts prints every contact whose name or email contains query,
// matched case-insensitively.
func searchContacts(query string) {
	contacts := load()
	query = strings.ToLower(query)

	matches := make([]Contact, 0)
	for _, c := range contacts {
		if strings.Contains(strings.ToLower(c.Name), query) ||
			strings.Contains(strings.ToLower(c.Email), query) {
			matches = append(matches, c)
		}
	}

	if len(matches) == 0 {
		fmt.Println("No matching contacts.")
		return
	}
	for _, c := range matches {
		fmt.Printf("%s | %s | %s\n", c.Name, c.Email, c.Phone)
	}
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
