package controls

import "fmt"

// ControlError implements GetType() wich is used in tests to test the validation.
type ControlError interface {
	Error() string
	GetType() string
}

// InvalidCategoryError is an error type for invalid categories
type InvalidCategoryError struct {
	Category string
}

// GetType returns a string containing the error Type
func (e *InvalidCategoryError) GetType() string {
	return "InvalidCategoryError"
}

func (e *InvalidCategoryError) Error() string {
	return fmt.Sprintf("Invalid category: %s", e.Category)
}

func newInvalidCategoryError(Category string) *InvalidCategoryError {
	return &InvalidCategoryError{
		Category: Category,
	}
}

// InvalidControlNameError is an error type for invalid commands
type InvalidControlNameError struct {
	Name string
}

// GetType returns a string containing the error Type
func (e *InvalidControlNameError) GetType() string {
	return "InvalidControlNameError"
}

func (e *InvalidControlNameError) Error() string {
	return fmt.Sprintf("Invalid name for control: %s Only ASCII letters (a-z), underscore (_) and hyphen (-) allowed.", e.Name)
}

func newInvalidControlNameError(name string) *InvalidControlNameError {
	return &InvalidControlNameError{
		Name: name,
	}
}

// InvalidCommandError is an error type for invalid commands
type InvalidCommandError struct {
	Category string
	Command  string
}

// GetType returns a string containing the error Type
func (e *InvalidCommandError) GetType() string {
	return "InvalidCommandError"
}

func (e *InvalidCommandError) Error() string {
	return fmt.Sprintf("Invalid comand for type %s: %s", e.Category, e.Command)
}

func newInvalidCommandError(category, command string) *InvalidCommandError {
	return &InvalidCommandError{
		Category: category,
		Command:  command,
	}
}

// InvalidAuthKeyError is an error type for invalid authKeys
type InvalidAuthKeyError struct {
	Name string
}

// GetType returns a string containing the error Type
func (e *InvalidAuthKeyError) GetType() string {
	return "InvalidAuthKeyError"
}

func (e *InvalidAuthKeyError) Error() string {
	return fmt.Sprintf("Invalid authKey: %s", e.Name)
}

func newInvalidAuthKeyError(name string) *InvalidAuthKeyError {
	return &InvalidAuthKeyError{
		Name: name,
	}
}

// NoAuthKeysError is an error type for invalid authKeys
type NoAuthKeysError struct {
	Name string
}

// GetType returns a string containing the error Type
func (e *NoAuthKeysError) GetType() string {
	return "NoAuthKeysError"
}

func (e *NoAuthKeysError) Error() string {
	return "No authKeys defined"
}

func newNoAuthKeysError() *NoAuthKeysError {
	return &NoAuthKeysError{
		Name: "all",
	}
}
