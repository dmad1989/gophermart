// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jsonobject

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject(in *jlexer.Lexer, out *Withdrawls) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Withdrawls, 0, 0)
			} else {
				*out = Withdrawls{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Withdraw
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject(out *jwriter.Writer, in Withdrawls) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Withdrawls) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdrawls) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdrawls) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdrawls) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject(l, v)
}
func easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject1(in *jlexer.Lexer, out *Withdraw) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "order":
			out.Order = string(in.String())
		case "sum":
			out.Sum = float64(in.Float64())
		case "processed_at":
			out.ProcessedDate = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject1(out *jwriter.Writer, in Withdraw) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"order\":"
		out.RawString(prefix[1:])
		out.String(string(in.Order))
	}
	{
		const prefix string = ",\"sum\":"
		out.RawString(prefix)
		out.Float64(float64(in.Sum))
	}
	if in.ProcessedDate != "" {
		const prefix string = ",\"processed_at\":"
		out.RawString(prefix)
		out.String(string(in.ProcessedDate))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Withdraw) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdraw) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdraw) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdraw) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject1(l, v)
}
func easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject2(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject2(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject2(l, v)
}
func easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject3(in *jlexer.Lexer, out *Orders) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Orders, 0, 0)
			} else {
				*out = Orders{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 Order
			(v4).UnmarshalEasyJSON(in)
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject3(out *jwriter.Writer, in Orders) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			(v6).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Orders) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Orders) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Orders) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Orders) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject3(l, v)
}
func easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject4(in *jlexer.Lexer, out *Order) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "number":
			out.Number = string(in.String())
		case "status":
			out.Status = string(in.String())
		case "accrual":
			out.Accrual = string(in.String())
		case "uploaded_at":
			out.UploadDate = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject4(out *jwriter.Writer, in Order) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"number\":"
		out.RawString(prefix[1:])
		out.String(string(in.Number))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	if in.Accrual != "" {
		const prefix string = ",\"accrual\":"
		out.RawString(prefix)
		out.String(string(in.Accrual))
	}
	{
		const prefix string = ",\"uploaded_at\":"
		out.RawString(prefix)
		out.String(string(in.UploadDate))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Order) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Order) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Order) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Order) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject4(l, v)
}
func easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject5(in *jlexer.Lexer, out *Balance) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "current":
			out.AccrualCurrent = float64(in.Float64())
		case "withdrawn":
			out.Withdrawn = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject5(out *jwriter.Writer, in Balance) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"current\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.AccrualCurrent))
	}
	{
		const prefix string = ",\"withdrawn\":"
		out.RawString(prefix)
		out.Float64(float64(in.Withdrawn))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Balance) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Balance) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfc1bcb3EncodeGithubComDmad1989GophermartInternalJsonobject5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Balance) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Balance) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfc1bcb3DecodeGithubComDmad1989GophermartInternalJsonobject5(l, v)
}
