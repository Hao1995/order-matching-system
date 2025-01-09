// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package events

import (
	"errors"
	"fmt"
)

const (
	// OrderEventTypeCREATE is a OrderEventType of type CREATE.
	OrderEventTypeCREATE OrderEventType = "CREATE"
	// OrderEventTypeCANCEL is a OrderEventType of type CANCEL.
	OrderEventTypeCANCEL OrderEventType = "CANCEL"
)

var ErrInvalidOrderEventType = errors.New("not a valid OrderEventType")

// String implements the Stringer interface.
func (x OrderEventType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x OrderEventType) IsValid() bool {
	_, err := ParseOrderEventType(string(x))
	return err == nil
}

var _OrderEventTypeValue = map[string]OrderEventType{
	"CREATE": OrderEventTypeCREATE,
	"CANCEL": OrderEventTypeCANCEL,
}

// ParseOrderEventType attempts to convert a string to a OrderEventType.
func ParseOrderEventType(name string) (OrderEventType, error) {
	if x, ok := _OrderEventTypeValue[name]; ok {
		return x, nil
	}
	return OrderEventType(""), fmt.Errorf("%s is %w", name, ErrInvalidOrderEventType)
}

// MarshalText implements the text marshaller method.
func (x OrderEventType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *OrderEventType) UnmarshalText(text []byte) error {
	tmp, err := ParseOrderEventType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

const (
	// SideBUY is a Side of type BUY.
	SideBUY Side = iota
	// SideSELL is a Side of type SELL.
	SideSELL
)

var ErrInvalidSide = errors.New("not a valid Side")

const _SideName = "BUYSELL"

var _SideMap = map[Side]string{
	SideBUY:  _SideName[0:3],
	SideSELL: _SideName[3:7],
}

// String implements the Stringer interface.
func (x Side) String() string {
	if str, ok := _SideMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Side(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Side) IsValid() bool {
	_, ok := _SideMap[x]
	return ok
}

var _SideValue = map[string]Side{
	_SideName[0:3]: SideBUY,
	_SideName[3:7]: SideSELL,
}

// ParseSide attempts to convert a string to a Side.
func ParseSide(name string) (Side, error) {
	if x, ok := _SideValue[name]; ok {
		return x, nil
	}
	return Side(0), fmt.Errorf("%s is %w", name, ErrInvalidSide)
}

// MarshalText implements the text marshaller method.
func (x Side) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Side) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseSide(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
