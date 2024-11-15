package string_utils

import (
	"bytes"
	"github.com/saichler/shared/go/src/interfaces"
)

// String is a wrapper over buff.bytes to make it seamless to concatenate strings
type String struct {
	buff               *bytes.Buffer
	TypesPrefix        bool
	AddSpaceWhenAdding bool
}

// New construct a new String instance and initialize the buff with the input string
func New(anys ...interface{}) *String {
	s := &String{}
	s.init()
	if anys != nil {
		for _, any := range anys {
			s.Add(s.StringOf(any))
		}
	}
	return s
}

// init initialize the buff if needed
func (s *String) init() {
	if s.buff == nil {
		s.buff = &bytes.Buffer{}
	}
}

// Add concatenate a string to this String instance
func (s *String) Add(strs ...string) *String {
	s.init()
	if s.AddSpaceWhenAdding && len(s.Bytes()) > 0 {
		s.buff.WriteString(" ")
	}
	if strs != nil {
		for _, str := range strs {
			s.buff.WriteString(str)
		}
	}
	return s
}

// Join concatenate a String instance to this String instance
func (s *String) Join(other *String) *String {
	s.init()
	s.buff.Write(other.buff.Bytes())
	return s
}

// String convert the String instance buffer to primitive string
func (s *String) String() string {
	s.init()
	return s.buff.String()
}

// IsBlank return is this String instance is blank
func (s *String) IsBlank() bool {
	s.init()
	return s.buff.Len() == 0
}

// Len return the length of the current string
func (s *String) Len() int {
	return s.buff.Len()
}

// Bytes return the bytes of the string
func (s *String) Bytes() []byte {
	return s.buff.Bytes()
}

func (s *String) AddBytes(bytes []byte) {
	s.buff.Write(bytes)
}

func (s *String) LogError() {
	interfaces.Logger().Error(s.buff.String())
}
