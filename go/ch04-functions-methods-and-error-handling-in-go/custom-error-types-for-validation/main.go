package main

import (
	"errors"
	"fmt"
	"strings"
)

// Person is the record we validate.
type Person struct {
	Name  string
	Email string
	Age   int
}

// ValidationError describes a single failed validation rule.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors collects every rule a record failed, not just the first.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	messages := make([]string, len(ve))
	for i, e := range ve {
		messages[i] = e.Error()
	}
	return strings.Join(messages, "; ")
}

// validate checks every rule against p and returns all the failures
// together, instead of stopping at the first one.
func validate(p Person) error {
	var errs ValidationErrors

	if p.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "cannot be empty"})
	}
	if !strings.Contains(p.Email, "@") {
		errs = append(errs, ValidationError{Field: "email", Message: "must contain @"})
	}
	if p.Age < 0 || p.Age > 150 {
		errs = append(errs, ValidationError{
			Field:   "age",
			Message: fmt.Sprintf("%d is out of valid range (0-150)", p.Age),
		})
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func main() {
	people := []Person{
		{Name: "Alice", Email: "alice@example.com", Age: 30},
		{Name: "", Email: "invalid-email", Age: -5},
		{Name: "Bob", Email: "bob@test.com", Age: 200},
	}

	for i, p := range people {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("Validating %+v:\n", p)

		err := validate(p)
		if err == nil {
			fmt.Println("  Valid")
			continue
		}

		var verrs ValidationErrors
		errors.As(err, &verrs)
		fmt.Println("  INVALID:")
		for _, v := range verrs {
			fmt.Printf("    - %s\n", v.Error())
		}
	}
}
