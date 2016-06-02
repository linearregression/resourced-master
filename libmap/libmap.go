package libmap

import (
	"encoding/json"
	"strings"
	"sync"
)

// NewTSafeMapBytes creates an instance of TSafeMapBytes
func NewTSafeMapBytes(data map[string][]byte) *TSafeMapBytes {
	mp := &TSafeMapBytes{}
	if data == nil {
		mp.data = make(map[string][]byte)
	} else {
		mp.data = data
	}
	return mp
}

// TSafeMapBytes is concurrency-safe map of bytes.
type TSafeMapBytes struct {
	data map[string][]byte
	sync.RWMutex
}

// Set value.
func (mp *TSafeMapBytes) Set(key string, value []byte) {
	mp.Lock()
	defer mp.Unlock()
	mp.data[key] = value
}

// Get value.
func (mp *TSafeMapBytes) Get(key string) []byte {
	mp.Lock()
	defer mp.Unlock()

	return mp.data[key]
}

// All returns all values.
func (mp *TSafeMapBytes) All() map[string][]byte {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string][]byte)
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

// ToJson returns the JSON encoded map.
func (mp *TSafeMapBytes) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}

// NewTSafeMapString creates an instance of TSafeMapString
func NewTSafeMapString(data map[string]string) *TSafeMapString {
	mp := &TSafeMapString{}
	if data == nil {
		mp.data = make(map[string]string)
	} else {
		mp.data = data
	}
	return mp
}

// TSafeMapString is concurrency-safe map of string.
type TSafeMapString struct {
	data map[string]string
	sync.RWMutex
}

// Set value.
func (mp *TSafeMapString) Set(key string, value string) {
	mp.Lock()
	defer mp.Unlock()
	mp.data[key] = value
}

// Get value.
func (mp *TSafeMapString) Get(key string) string {
	mp.Lock()
	defer mp.Unlock()

	return mp.data[key]
}

// Delete by key.
func (mp *TSafeMapString) Delete(key string) string {
	var value string

	mp.Lock()
	defer mp.Unlock()

	value = mp.data[key]
	delete(mp.data, key)

	return value
}

// All returns all values in map.
func (mp *TSafeMapString) All() map[string]string {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string]string)
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

// ToJson returns the JSON encoded values.
func (mp *TSafeMapString) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}

// NewTSafeMapStrings returns instance of TSafeMapStrings
func NewTSafeMapStrings(data map[string][]string) *TSafeMapStrings {
	mp := &TSafeMapStrings{}
	if data == nil {
		mp.data = make(map[string][]string)
	} else {
		mp.data = data
	}
	return mp
}

// TSafeMapStrings is concurrency-safe map of []string.
type TSafeMapStrings struct {
	data map[string][]string
	sync.RWMutex
}

// Set value.
func (mp *TSafeMapStrings) Set(key string, value []string) {
	mp.Lock()
	defer mp.Unlock()
	mp.data[key] = value
}

// Append to slice by key.
func (mp *TSafeMapStrings) Append(key string, value string) {
	mp.Lock()
	defer mp.Unlock()
	mp.data[key] = append(mp.data[key], value)
}

// Get a slice of value.
func (mp *TSafeMapStrings) Get(key string) []string {
	mp.Lock()
	defer mp.Unlock()

	original, ok := mp.data[key]
	if !ok {
		return make([]string, 0)
	}

	copydata := make([]string, len(original))
	for i, value := range original {
		copydata[i] = value
	}

	return copydata
}

// Reset wipes all values.
func (mp *TSafeMapStrings) Reset(key string) {
	mp.Lock()
	defer mp.Unlock()

	mp.data[key] = make([]string, 0)
}

// All returns all values.
func (mp *TSafeMapStrings) All() map[string][]string {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string][]string)
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

// ToJson returns JSON encoded values.
func (mp *TSafeMapStrings) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}

// NewTSafeMapCounter creates an instance of TSafeMapCounter
func NewTSafeMapCounter(data map[string]int) *TSafeMapCounter {
	mp := &TSafeMapCounter{}
	if data == nil {
		mp.data = make(map[string]int)
	} else {
		mp.data = data
	}
	return mp
}

// TSafeMapCounter is concurrency-safe map of counter.
type TSafeMapCounter struct {
	data map[string]int
	sync.RWMutex
}

// Incr increments value by X.
func (mp *TSafeMapCounter) Incr(key string, value int) {
	mp.Lock()
	mp.data[key] = mp.data[key] + value
	mp.Unlock()
}

// Get value.
func (mp *TSafeMapCounter) Get(key string) int {
	mp.RLock()
	defer mp.RUnlock()

	data, ok := mp.data[key]
	if !ok {
		data = 0
	}

	return data
}

// Reset wipes count data to 0.
func (mp *TSafeMapCounter) Reset(key string) {
	mp.Lock()
	mp.data[key] = 0
	mp.Unlock()
}

// All returns all count values.
func (mp *TSafeMapCounter) All() map[string]int {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string]int)
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

// ToJson returns JSON encoded values.
func (mp *TSafeMapCounter) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}

// NewTSafeNestedMapInterface creates an instance of TSafeNestedMapInterface
func NewTSafeNestedMapInterface(data map[string]interface{}) *TSafeNestedMapInterface {
	mp := &TSafeNestedMapInterface{}
	if data == nil {
		mp.data = make(map[string]interface{})
	} else {
		mp.data = data
	}
	return mp
}

// TSafeNestedMapInterface is concurrency-safe map of interface.
type TSafeNestedMapInterface struct {
	data map[string]interface{}
	sync.RWMutex
}

func (mp *TSafeNestedMapInterface) initNestedMap(key string) {
	// Split key by dot, loop deeper into the nesting & create the maps
	keyParts := strings.Split(key, ".")

	mp.Lock()
	m := mp.data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		_, ok := m[keyPart]
		if !ok {
			m[keyPart] = make(map[string]interface{})
		}

		m = m[keyPart].(map[string]interface{})
	}
	mp.Unlock()
}

// Set value.
func (mp *TSafeNestedMapInterface) Set(key string, value interface{}) {
	mp.initNestedMap(key)

	keyParts := strings.Split(key, ".")
	lastPart := keyParts[len(keyParts)-1]

	mp.Lock()
	m := mp.data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		m = m[keyPart].(map[string]interface{})
	}

	m[lastPart] = value
	mp.Unlock()
}

// Get value.
func (mp *TSafeNestedMapInterface) Get(key string) interface{} {
	var data interface{}

	// Split key by dot and loop deeper into the nesting
	keyParts := strings.Split(key, ".")
	lastPart := keyParts[len(keyParts)-1]

	mp.RLock()
	m := mp.data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		m = m[keyPart].(map[string]interface{})
	}

	data = m[lastPart]
	mp.RUnlock()

	return data
}

// All returns all values.
func (mp *TSafeNestedMapInterface) All() map[string]interface{} {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string]interface{})
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

// ToJson returns JSON encoded values.
func (mp *TSafeNestedMapInterface) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}
