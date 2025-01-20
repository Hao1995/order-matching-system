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
	// OrderTypeBuy is a OrderType of type Buy.
	OrderTypeBuy OrderType = iota
	// OrderTypeSell is a OrderType of type Sell.
	OrderTypeSell
)

var ErrInvalidOrderType = errors.New("not a valid OrderType")

const _OrderTypeName = "BuySell"

var _OrderTypeMap = map[OrderType]string{
	OrderTypeBuy:  _OrderTypeName[0:3],
	OrderTypeSell: _OrderTypeName[3:7],
}

// String implements the Stringer interface.
func (x OrderType) String() string {
	if str, ok := _OrderTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("OrderType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x OrderType) IsValid() bool {
	_, ok := _OrderTypeMap[x]
	return ok
}

var _OrderTypeValue = map[string]OrderType{
	_OrderTypeName[0:3]: OrderTypeBuy,
	_OrderTypeName[3:7]: OrderTypeSell,
}

// ParseOrderType attempts to convert a string to a OrderType.
func ParseOrderType(name string) (OrderType, error) {
	if x, ok := _OrderTypeValue[name]; ok {
		return x, nil
	}
	return OrderType(0), fmt.Errorf("%s is %w", name, ErrInvalidOrderType)
}

// MarshalText implements the text marshaller method.
func (x OrderType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *OrderType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseOrderType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
