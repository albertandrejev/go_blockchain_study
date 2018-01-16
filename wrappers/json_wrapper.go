package wrappers

import "encoding/json"

//IJSONWrapper wrapper interface
type IJSONWrapper interface {
	Encode(v interface{}) ([]byte, error)
}

//JSONWrapper Json serilizier wrapper
type JSONWrapper struct {
}

//Encode object to json byte array
func (t JSONWrapper) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
