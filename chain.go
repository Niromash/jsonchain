package main

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type JsonChain[T string, K any] map[T]K
type JsonOutput []byte
type JsonChainKeyNotExist struct {
	error
}
type JsonChainKeyAlreadyExist struct {
	error
}

func (j *JsonChainKeyNotExist) Error() string {
	return "key not exist"
}

func (j *JsonChainKeyAlreadyExist) Error() string {
	return "key already exist"
}

// NewJsonChain returns a new JsonChain with given types T and K.
func NewJsonChain[T string, K any]() JsonChain[T, K] {
	return make(JsonChain[T, K])
}

// Set sets the value of the key to the value and send back the JsonChain object.
func (j JsonChain[T, K]) Set(key T, value K) JsonChain[T, K] {
	j[key] = value
	return j
}

// SetWithError sets the value of the key to the value and throws an error if the key already exists.
func (j JsonChain[T, K]) SetWithError(key T, value K) error {
	if _, ok := j[key]; ok {
		return &JsonChainKeyAlreadyExist{}
	}
	j.Set(key, value)
	return nil
}

// GetWithError gets the value of the key and throws an error if the key does not exist.
func (j JsonChain[T, K]) GetWithError(key T) (K, error) {
	value := j.Get(key)
	if reflect.ValueOf(value).IsZero() {
		return reflect.ValueOf(*new(K)).Interface().(K), &JsonChainKeyNotExist{}
	}
	return value, nil
}

// Get gets the value of the key.
func (j JsonChain[T, K]) Get(key T) K {
	return j[key]
}

// Clear clears the JsonChain object.
func (j *JsonChain[T, K]) Clear() JsonChain[T, K] {
	*j = map[T]K{}
	return *j
}

// Load loads the JsonChain object from the given other JsonChain given in parameter, and overwrite all existing data.
func (j *JsonChain[T, K]) Load(otherChain JsonChain[T, K]) {
	j.Clear()
	*j = otherChain
}

// LoadFromBytes loads the JsonChain object from the given json string, and overwrite all existing data.
func (j *JsonChain[T, K]) LoadFromBytes(data []byte) error {
	j.Clear()
	return j.AppendFromBytes(data)
}

func (j *JsonChain[T, K]) Each(fun func(key T, value K)) {
	for key, value := range *j {
		fun(key, value)
	}
}

// Copy copies the given JsonChain object to the current one without overwrite data.
func (j *JsonChain[T, K]) Copy(otherChain JsonChain[T, K]) {
	otherChain.Each(func(key T, value K) {
		j.Set(key, value)
	})
	return
}

// Append appends the given JsonChain object in parameter to current JsonChain object.
func (j JsonChain[T, K]) Append(otherChain JsonChain[T, K]) {
	otherChain.Each(func(key T, value K) {
		_ = j.SetWithError(key, value)
	})
}

// AppendFromBytes appends the given json string to JsonChain object.
func (j JsonChain[T, K]) AppendFromBytes(data []byte) error {
	return json.NewDecoder(bytes.NewBuffer(data)).Decode(&j)
}

// Clone clones the current JsonChain object.
func (j JsonChain[T, K]) Clone() JsonChain[T, K] {
	clone := NewJsonChain[T, K]()
	clone.Copy(j)
	return clone
}

// Equal returns true if the current JsonChain object is equal to the given one.
func (j JsonChain[T, K]) Equal(otherChain JsonChain[T, K]) bool {
	return reflect.DeepEqual(j, otherChain)
}

// ToJson converts the JsonChain object to json string.
func (j JsonChain[T, K]) ToJson() (JsonOutput, error) {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(j); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Pretty prints the JsonChain object in a pretty format.
func (o *JsonOutput) Pretty() string {
	var buffer bytes.Buffer
	if err := json.Indent(&buffer, *o, "", "  "); err != nil {
		return ""
	}
	return buffer.String()
}
