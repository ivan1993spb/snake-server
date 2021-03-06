// Code generated by ffjson <https://github.com/pquerna/ffjson>. DO NOT EDIT.
// source: ./player/message.go

package player

import (
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

// MarshalJSON marshal bytes to json - template
func (j *Message) MarshalJSON() ([]byte, error) {
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
func (j *Message) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"type":`)

	{

		obj, err = j.Type.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(obj)

	}
	buf.WriteString(`,"payload":`)
	/* Interface types must use runtime reflection. type=interface {} kind=interface */
	err = buf.Encode(j.Payload)
	if err != nil {
		return err
	}
	buf.WriteByte('}')
	return nil
}

// MarshalJSON marshal bytes to json - template
func (j *MessageSize) MarshalJSON() ([]byte, error) {
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
func (j *MessageSize) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"width":`)
	fflib.FormatBits2(buf, uint64(j.Width), 10, false)
	buf.WriteString(`,"height":`)
	fflib.FormatBits2(buf, uint64(j.Height), 10, false)
	buf.WriteByte('}')
	return nil
}
