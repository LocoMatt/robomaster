package result

import (
	"encoding/json"
	"fmt"

	"github.com/brunoga/unitybridge/unity/key"
)

// Reesult represents a result from an operation on a key. The zero value
// is not valid.
type Result struct {
	key       *key.Key
	tag       uint64
	errorCode int32
	errorDesc string
	value     any
}

type jsonResultValue struct {
	Value any `json:"value"`
}

type jsonResult struct {
	Key   uint32
	Tag   uint64
	Error int32
	Value any
}

func New(key *key.Key, tag uint64, errorCode int32, errorDesc string,
	value any) *Result {
	return &Result{
		key:       key,
		tag:       tag,
		errorCode: errorCode,
		errorDesc: errorDesc,
		value:     value,
	}
}

// NewFromJSON creates a new Result from the given JSON data. Any errors are
// reported in the Result itself and should be handled by anyone that cares
// about it.
func NewFromJSON(jsonData []byte) *Result {
	r := &Result{}

	if len(jsonData) == 0 {
		r.errorCode = -1
		r.errorDesc = "empty or nil json data"
		return r
	}

	jr := jsonResult{}

	err := json.Unmarshal(jsonData, &jr)
	if err != nil {
		r.errorCode = -1
		r.errorDesc = fmt.Sprintf("error unmarshalling json data: %s",
			err.Error())
		return r
	}

	key, err := key.FromSubType(jr.Key)
	if err != nil {
		fmt.Printf("error creating key from sub type %d: %s\n", jr.Key, err.Error())
		r.errorCode = -1
		r.errorDesc = fmt.Sprintf("error creating key from sub type %d: %s",
			jr.Key, err.Error())
		return r
	}

	errorDesc := ""
	if jr.Error != 0 {
		errorDesc = fmt.Sprintf("error %d", jr.Error)
	}

	var value any
	var ok bool
	switch jr.Value.(type) {
	case map[string]interface{}:
		outerValue := jr.Value.(map[string]interface{})
		value, ok = outerValue["value"]
		if !ok {
			r.errorCode = -1
			r.errorDesc = fmt.Sprintf("value field not found: %v",
				outerValue)
			return r
		}
	default:
		value = jr.Value
	}

	r.key = key
	r.tag = jr.Tag
	r.errorCode = jr.Error
	r.errorDesc = errorDesc
	r.value = value

	return r
}

// Key returns the key associated with this result.
func (r *Result) Key() *key.Key {
	return r.key
}

// SetKey sets the key associated with this result.
func (r *Result) SetKey(key *key.Key) {
	r.key = key
}

// Tag returns the tag associated with this result.
func (r *Result) Tag() uint64 {
	return r.tag
}

// ErrorCode returns the error code associated with this result.
func (r *Result) ErrorCode() int32 {
	return r.errorCode
}

// SetErrorCode sets the error code associated with this result.
func (r *Result) SetErrorCode(errorCode int32) {
	r.errorCode = errorCode
}

// ErrorDesc returns the error description associated with this result.
func (r *Result) ErrorDesc() string {
	return r.errorDesc
}

// SetErrorDesc sets the error description associated with this result.
func (r *Result) SetErrorDesc(errorDesc string) {
	r.errorDesc = errorDesc
}

// Value returns the value associated with this result.
func (r *Result) Value() any {
	return r.value
}

// SetValue sets the value associated with this result.
func (r *Result) SetValue(value any) {
	r.value = value
}

// Succeeded returns true if this result represents a successful operation.
func (r *Result) Succeeded() bool {
	return r.errorCode == 0
}

// String returns a string representation of this result.
func (r *Result) String() string {
	return fmt.Sprintf("Result{Key: %s, Tag: %d, ErrorCode: %d, ErrorDesc: "+
		"%s, Value: %v}", r.key, r.tag, r.errorCode, r.errorDesc, r.value)
}

func (r *Result) MarshalJSON() ([]byte, error) {
	jr := jsonResult{
		Key:   r.key.SubType(),
		Tag:   r.tag,
		Error: r.errorCode,
		Value: jsonResultValue{r.value},
	}

	return json.Marshal(jr)
}
