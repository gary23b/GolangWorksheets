/*
Master the built int package: encoding/json/v2
https://pkg.go.dev/encoding/json/v2
*/
package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"time"
)

// To work with JSON data, we can use the encoding/json/v2 package.
// It provides functions to marshal and unmarshal JSON data.

// Make a type called data1 that has a string field called Name
// Make sure you add the JSON struct tag for the Name field: `json:"name"`

type data1 struct {
	// !->
	Name string `json:"name"`
	// <-!
}

// Now we create a function to marshal the data1 struct into JSON.
func MarshalData1(d data1) ([]byte, error) {
	// !->
	return json.Marshal(d)
	// <-!
}

// We can also create a function to unmarshal JSON data into the data1 struct.
func UnmarshalData1(b []byte) (data1, error) {
	// !->
	var ret data1
	err := json.Unmarshal(b, &ret)
	return ret, err
	// <-!
}

// We can also marshal and unmarshal to io.Writer and io.Reader using json.NewEncoder and json.NewDecoder.
func MarshalData1ToWriter(d data1, w io.Writer) error {
	// !->
	return json.MarshalWrite(w, d)
	// <-!
}

func UnmarshalData1FromReader(r io.Reader) (data1, error) {
	// !->
	var ret data1
	err := json.UnmarshalRead(r, &ret)
	return ret, err
	// <-!
}

// json/v2 has the following struct tags

// omitempty: omits fields that are empty according to the JSON specification.
// Such as an empty string, null, or an empty array.

// omitzero: omits things that are zero from a golang point of view.
// such as false for a bool.

// case:ignore or case:strict (default): specifies how the JSON keys should be matched with the struct fields.
// case:ignore: ignores the case when matching JSON keys with struct fields.
// case:strict (default): requires an exact match of the case between JSON keys and struct fields.

// string: causes fields that would be encoded as JSON numbers to be encoded as JSON strings instead.

// embed: promote a named struct field to be treated as if its fields are part of the outer struct.
// this won't get used much...

// Fill in the data struct below with the correct JSON struct tags for each field.
/*
ID    string    `json:"id"`     // always omit
Name  string    `json:"name"`   // omit if ""
Name2 *string   `json:"name2"`  // omit if "" or nil
Name3 *string   `json:"name3"`  // omit if nil
List  []string  `json:"list"`   // omit if nil
List2 []string  `json:"list2"`  // omit if JSON will be {}
Bool1 bool      `json:"bool1"`  // omit if false
Sub   subStruct `json:"sub"`    // omit if nothing filled in
Case1 string    `json:"case1"`  // omit if "", make it also match to CASE1
*/

type subStruct struct {
	// !->
	SubName string `json:"subName,omitempty"`
	// <-!
}

type data2 struct {
	// !->
	ID    string    `json:"-"`                           // always omit
	Name  string    `json:"name,omitempty"`              // omit if ""
	Name2 *string   `json:"name2,omitempty"`             // omit if "" or nil
	Name3 *string   `json:"name3,omitzero"`              // omit if nil
	List  []string  `json:"list,omitzero"`               // omit if nil
	List2 []string  `json:"list2,omitempty"`             // omit if JSON will be {}
	Bool1 bool      `json:"bool1,omitzero"`              // omit if false
	Sub   subStruct `json:"sub,omitzero"`                // omit if nothing filled in
	Case1 string    `json:"case1,omitempty,case:ignore"` // omit if "", make it also match to CASE1
	// <-!
}

// The final thing that is used a lot is custom marshaling and unmarshaling.
// There are a couple ways to implement custom marshaling and unmarshaling with v2,
// but one of the faster options is to implement the json.MarshalerTo json.UnmarshalerFrom interfaces.

// the interfaces look like this:
/*
type MarshalerTo interface {
    MarshalJSONTo(*jsontext.Encoder) error
}

type UnmarshalerFrom interface {
    UnmarshalJSONFrom(*jsontext.Decoder) error
}
*/

// so for example. Salesforce dates are represented as strings in JSON, but you might want to parse them into time.Time in Go.
// You can implement custom marshaling and unmarshaling for such fields using the interfaces above.
// They give you the date in the format "2006-01-02"
// https://go.dev/doc/tutorial/json#custom-marshal-and-unmarshal

type SalesforceDate struct {
	time.Time
}

// !->
func (s SalesforceDate) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, s.Format("2006-01-02"))
}

func (s *SalesforceDate) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var dateStr string
	if err := json.UnmarshalDecode(dec, &dateStr); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	s.Time = t
	return nil
}

// <-!

// Or we could have private data fields and implement custom marshaling and unmarshaling for them as well.
type Version struct {
	major int64
	minor int64
	patch int64
}

// !->
func (v Version) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch))
}

func (v *Version) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	var s string
	if err := json.UnmarshalDecode(dec, &s); err != nil {
		return err
	}

	_, err := fmt.Sscanf(s, "%d.%d.%d", &v.major, &v.minor, &v.patch)
	return err
}

// <-!

// jsontext.Value is useful for discriminated unions or dynamic message types.
// A common pattern: implement custom unmarshaling with UnmarshalJSONFrom to handle
// type detection and route to the correct struct based on a "type" field.
// This is especially useful with APIs like Salesforce that declare the message
// type in a field before providing the actual data.

type DynamicMessage struct {
	Type     string           `json:"type"`
	Order    *OrderMessage    `json:"-"`
	Customer *CustomerMessage `json:"-"`
}

type OrderMessage struct {
	OrderID string  `json:"orderId"`
	Amount  float64 `json:"amount"`
}

type CustomerMessage struct {
	CustomerID string `json:"customerId"`
	Name       string `json:"name"`
}

func (v *DynamicMessage) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	type messageEnvelope struct {
		Type string         `json:"type"`
		Data jsontext.Value `json:"data"` // Defer parsing until we know the type
	}

	var envelope messageEnvelope
	if err := json.UnmarshalDecode(dec, &envelope); err != nil {
		return err
	}

	v.Type = envelope.Type

	switch envelope.Type {
	// !->
	case "order":
		var msg OrderMessage
		if err := json.Unmarshal(envelope.Data, &msg); err != nil {
			return err
		}
		v.Order = &msg
		return nil
	case "customer":
		var msg CustomerMessage
		if err := json.Unmarshal(envelope.Data, &msg); err != nil {
			return err
		}
		v.Customer = &msg
		return nil
	// <-!
	default:
		return fmt.Errorf("unknown message type: %s", envelope.Type)
	}
}
