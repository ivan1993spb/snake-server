// Code generated by ffjson <https://github.com/pquerna/ffjson>. DO NOT EDIT.
// source: world/event.go

package world

import (
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

// MarshalJSON marshal bytes to json - template
func (j *Event) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if j == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := j.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONBuf marshal buff to json - template
func (j *Event) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"Type":`)

	{

		obj, err = j.Type.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(obj)

	}
	buf.WriteString(`,"Payload":`)
	/* Interface types must use runtime reflection. type=interface {} kind=interface */
	err = buf.Encode(j.Payload)
	if err != nil {
		return err
	}
	buf.WriteByte('}')
	return nil
}
