// Code generated by ffjson <https://github.com/pquerna/ffjson>. DO NOT EDIT.
// source: ./objects/apple/apple.go

package apple

import (
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

// MarshalJSON marshal bytes to json - template
func (j *apple) MarshalJSON() ([]byte, error) {
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
func (j *apple) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"id":`)
	fflib.FormatBits2(buf, uint64(j.ID), 10, false)
	buf.WriteString(`,"dot":`)

	{

		obj, err = j.Dot.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(obj)

	}
	buf.WriteString(`,"type":`)
	fflib.WriteJsonString(buf, string(j.Type))
	buf.WriteByte('}')
	return nil
}
