// Code generated by go-enum DO NOT EDIT.
// Version: 0.6.0
// Revision: 919e61c0174b91303753ee3898569a01abb32c97
// Build Date: 2023-12-18T15:54:43Z
// Built By: goreleaser

package events

import (
	"errors"
	"fmt"
)

const (
	// MatchingEventTypeCreate is a MatchingEventType of type Create.
	MatchingEventTypeCreate MatchingEventType = "Create"
	// MatchingEventTypeCancel is a MatchingEventType of type Cancel.
	MatchingEventTypeCancel MatchingEventType = "Cancel"
)

var ErrInvalidMatchingEventType = errors.New("not a valid MatchingEventType")

// String implements the Stringer interface.
func (x MatchingEventType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x MatchingEventType) IsValid() bool {
	_, err := ParseMatchingEventType(string(x))
	return err == nil
}

var _MatchingEventTypeValue = map[string]MatchingEventType{
	"Create": MatchingEventTypeCreate,
	"Cancel": MatchingEventTypeCancel,
}

// ParseMatchingEventType attempts to convert a string to a MatchingEventType.
func ParseMatchingEventType(name string) (MatchingEventType, error) {
	if x, ok := _MatchingEventTypeValue[name]; ok {
		return x, nil
	}
	return MatchingEventType(""), fmt.Errorf("%s is %w", name, ErrInvalidMatchingEventType)
}

// MarshalText implements the text marshaller method.
func (x MatchingEventType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *MatchingEventType) UnmarshalText(text []byte) error {
	tmp, err := ParseMatchingEventType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
