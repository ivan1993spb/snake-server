// Code generated by ffjson <https://github.com/pquerna/ffjson>. DO NOT EDIT.
// source: objects/watermelon/watermelon.go

package watermelon

import (
	fflib "github.com/pquerna/ffjson/fflib/v1"
)

// MarshalJSON marshal bytes to json - template
func (j *watermelon) MarshalJSON() ([]byte, error) {
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
func (j *watermelon) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if j == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{"uuid":`)
	fflib.WriteJsonString(buf, string(j.UUID))
	buf.WriteByte(',')
	if len(j.Dots) != 0 {
		buf.WriteString(`"dots":`)
		if j.Dots != nil {
			buf.WriteString(`[`)
			for i, v := range j.Dots {
				if i != 0 {
					buf.WriteString(`,`)
				}

				{

					obj, err = v.MarshalJSON()
					if err != nil {
						return err
					}
					buf.Write(obj)

				}
			}
			buf.WriteString(`]`)
		} else {
			buf.WriteString(`null`)
		}
		buf.WriteByte(',')
	}
	buf.WriteString(`"type":`)
	fflib.WriteJsonString(buf, string(j.Type))
	buf.WriteByte('}')
	return nil
}
