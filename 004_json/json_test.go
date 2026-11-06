package main

import (
	"bytes"
	"encoding/json/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMarshalData1(t *testing.T) {
	in := data1{Name: ""}
	want, err := MarshalData1(in)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"name":""}`), want)

	out, err := UnmarshalData1(want)
	require.NoError(t, err)
	require.Equal(t, in, out)

	in = data1{Name: "test"}
	want, err = MarshalData1(in)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"name":"test"}`), want)

	out, err = UnmarshalData1(want)
	require.NoError(t, err)
	require.Equal(t, in, out)
}

func TestMarshalData1ToWriter(t *testing.T) {
	in := data1{Name: ""}
	wr := &bytes.Buffer{}
	err := MarshalData1ToWriter(in, wr)
	require.NoError(t, err)
	require.Equal(t, `{"name":""}`, wr.String())

	rd := bytes.Buffer{}
	rd.Write(wr.Bytes())
	out, err := UnmarshalData1FromReader(&rd)
	require.NoError(t, err)
	require.Equal(t, in, out)

	in = data1{Name: "test"}
	wr = &bytes.Buffer{}
	err = MarshalData1ToWriter(in, wr)
	require.NoError(t, err)
	require.Equal(t, `{"name":"test"}`, wr.String())

	rd = bytes.Buffer{}
	rd.Write(wr.Bytes())
	out, err = UnmarshalData1FromReader(&rd)
	require.NoError(t, err)
	require.Equal(t, in, out)
}

func TestStructTags(t *testing.T) {
	in := data2{Name: ""}
	got, err := json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(got))

	in = data2{ID: "123"}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(got))

	in = data2{Name: "test"}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"name":"test"}`, string(got))

	in = data2{Name2: new("")}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(got))

	in = data2{Name2: new("abc")}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"name2":"abc"}`, string(got))

	in = data2{Name3: new("")}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"name3":""}`, string(got))

	in = data2{List: []string{}}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"list":[]}`, string(got))

	in = data2{List: []string{"abc", "123"}}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"list":["abc","123"]}`, string(got))

	in = data2{List2: []string{}}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(got))

	in = data2{List2: []string{"abc", "123"}}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"list2":["abc","123"]}`, string(got))

	in = data2{Bool1: true}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"bool1":true}`, string(got))

	in = data2{Sub: subStruct{SubName: "test"}}
	got, err = json.Marshal(in)
	require.NoError(t, err)
	require.Equal(t, `{"sub":{"subName":"test"}}`, string(got))

	// Test case-insensitive un-marshalling for Case1 field
	jsonStr := `{"case1":"test"}`
	in = data2{}
	err = json.Unmarshal([]byte(jsonStr), &in)
	require.NoError(t, err)
	require.Equal(t, in, data2{Case1: "test"})

	jsonStr = `{"CASE1":"test2"}`
	in = data2{Case1: "test2"}
	err = json.Unmarshal([]byte(jsonStr), &in)
	require.NoError(t, err)
	require.Equal(t, in, data2{Case1: "test2"})
}

func Test_SalesforceDate(t *testing.T) {
	type wrapper struct {
		Date    SalesforceDate  `json:"date"`
		DatePtr *SalesforceDate `json:"datePtr"`
	}
	w := wrapper{
		Date:    SalesforceDate{Time: time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)},
		DatePtr: &SalesforceDate{Time: time.Date(2000, 1, 3, 3, 4, 5, 6, time.UTC)},
	}

	got, err := json.Marshal(w)
	require.NoError(t, err)
	require.Equal(t, `{"date":"2000-01-02","datePtr":"2000-01-03"}`, string(got))

	var w2 wrapper
	err = json.Unmarshal(got, &w2)
	require.NoError(t, err)

	want := wrapper{
		Date:    SalesforceDate{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
		DatePtr: &SalesforceDate{Time: time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC)},
	}
	require.Equal(t, want, w2)
}

func Test_Version(t *testing.T) {
	type wrapper struct {
		Version    Version  `json:"version"`
		VersionPtr *Version `json:"versionPtr"`
	}
	w := wrapper{
		Version:    Version{major: 1, minor: 2, patch: 3},
		VersionPtr: &Version{major: 4, minor: 5, patch: 6},
	}

	got, err := json.Marshal(w)
	require.NoError(t, err)
	require.Equal(t, `{"version":"1.2.3","versionPtr":"4.5.6"}`, string(got))

	var w2 wrapper
	err = json.Unmarshal(got, &w2)
	require.NoError(t, err)

	want := wrapper{
		Version:    Version{major: 1, minor: 2, patch: 3},
		VersionPtr: &Version{major: 4, minor: 5, patch: 6},
	}
	require.Equal(t, want, w2)
}

func Test_DynamicMessage_DiscriminatedUnion(t *testing.T) {
	// Test order message
	orderJSON := []byte(`{"type":"order","data":{"orderId":"12345","amount":99.99}}`)
	var msg DynamicMessage
	err := json.Unmarshal(orderJSON, &msg)
	require.NoError(t, err)
	require.Equal(t, "order", msg.Type)
	require.NotNil(t, msg.Order)
	require.Equal(t, OrderMessage{OrderID: "12345", Amount: 99.99}, *msg.Order)
	require.Nil(t, msg.Customer)

	// Test customer message
	customerJSON := []byte(`{"type":"customer","data":{"customerId":"cust-001","name":"John Doe"}}`)
	msg = DynamicMessage{}
	err = json.Unmarshal(customerJSON, &msg)
	require.NoError(t, err)
	require.Equal(t, "customer", msg.Type)
	require.Nil(t, msg.Order)
	require.NotNil(t, msg.Customer)
	require.Equal(t, CustomerMessage{CustomerID: "cust-001", Name: "John Doe"}, *msg.Customer)

	// Test unknown type
	unknownJSON := []byte(`{"type":"unknown","data":{}}`)
	msg = DynamicMessage{}
	err = json.Unmarshal(unknownJSON, &msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown message type")
}

func Test_RawMessage_DiscriminatedUnion2(t *testing.T) {
	// Test order message
	orderJSON := []byte(`{"type":"order","data":{"orderId":"12345","amount":99.99}}`)

	var w2 DynamicMessage
	err := json.Unmarshal(orderJSON, &w2)
	require.NoError(t, err)
	require.Equal(t, DynamicMessage{Type: "order", Order: &OrderMessage{OrderID: "12345", Amount: 99.99}}, w2)

	// Test customer message
	customerJSON := []byte(`{"type":"customer","data":{"customerId":"cust-001","name":"John Doe"}}`)
	w2 = DynamicMessage{}
	err = json.Unmarshal(customerJSON, &w2)
	require.NoError(t, err)
	require.Equal(t, DynamicMessage{Type: "customer", Customer: &CustomerMessage{CustomerID: "cust-001", Name: "John Doe"}}, w2)

	// Test unknown type
	unknownJSON := []byte(`{"type":"unknown","data":{}}`)
	w2 = DynamicMessage{}
	err = json.Unmarshal(unknownJSON, &w2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown message type")
}
