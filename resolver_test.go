package pathresolver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TryGetValueByPath(t *testing.T) {
	type NestedStruct struct {
		FloatField    float64
		floatPtrField *float64
	}

	type TestStruct struct {
		IntField           int
		stringField        string
		stringNilField     *string
		floatPtrToPtrValue **float64
		MapField           map[string]string
		mapFieldOfPointers *map[string]*[]*string
		nestedStruct       NestedStruct
		nestedPtrStruct    *NestedStruct
	}

	floatValue := 300.1
	floatPtrValue := &floatValue

	strVal1 := "hello"
	strVal2 := "world"
	var strVal3 string
	var strVal4 *string
	nestedSlice := make([]*string, 0)
	nestedSlice = append(nestedSlice, &strVal1)
	nestedSlice = append(nestedSlice, nil)
	nestedSlice = append(nestedSlice, &strVal3)
	nestedSlice = append(nestedSlice, strVal4)
	nestedSlice = append(nestedSlice, &strVal2)

	mapFieldOfPointers := make(map[string]*[]*string)
	mapFieldOfPointers["key1"] = &nestedSlice

	testInstance := TestStruct{
		IntField:           599,
		stringField:        "hello",
		MapField:           make(map[string]string),
		mapFieldOfPointers: &mapFieldOfPointers,
		nestedStruct: NestedStruct{
			FloatField:    1.5,
			floatPtrField: &floatValue,
		},
		floatPtrToPtrValue: &floatPtrValue,
	}

	testInstance.MapField["mapFieldKey1"] = "bye"

	t.Run("should get error", func(t *testing.T) {
		actual, err := TryGetValueByPath("IntField", 300)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "data arg must be a struct, map, array or slice")
	})

	t.Run("should get int value", func(t *testing.T) {
		actual, err := TryGetValueByPath("IntField", testInstance)
		assert.Nil(t, err)
		assert.Equal(t, 599, actual)
	})

	t.Run("should get int value by ref", func(t *testing.T) {
		actual, err := TryGetValueByPath("IntField", &testInstance)
		assert.Nil(t, err)
		assert.Equal(t, 599, actual)
	})

	t.Run("should get string field", func(t *testing.T) {
		actual, err := TryGetValueByPath("stringField", &testInstance)
		assert.Nil(t, err)
		assert.Equal(t, "hello", actual)
	})

	t.Run("should get string nil field", func(t *testing.T) {
		actual, err := TryGetValueByPath("stringNilField", &testInstance)
		assert.Nil(t, err)
		assert.Equal(t, nil, actual)
	})

	t.Run("should get map field", func(t *testing.T) {
		actual, err := TryGetValueByPath("MapField.mapFieldKey1", &testInstance)
		assert.Nil(t, err)
		assert.Equal(t, "bye", actual)
	})

	t.Run("should get map field", func(t *testing.T) {
		actual, err := TryGetValueByPath("MapField.mapFieldKey2", &testInstance)
		assert.Nil(t, err)
		assert.Equal(t, nil, actual)
	})

	t.Run("should get nested struct ref field", func(t *testing.T) {
		actual, err := TryGetValueByPath("nestedStruct.floatPtrField", &testInstance)
		assert.Nil(t, err)
		assert.Equal(t, 300.1, actual)
	})

	t.Run("should get ptr to prt float field value", func(t *testing.T) {
		actual, err := TryGetValueByPath("floatPtrToPtrValue", &testInstance)
		assert.Nil(t, err)
		assert.Equal(t, 300.1, actual)
	})

	t.Run("should get value from map", func(t *testing.T) {
		actual1, err1 := TryGetValueByPath("key1[0]", &mapFieldOfPointers)
		actual2, err2 := TryGetValueByPath("key1[4]", &mapFieldOfPointers)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, "hello", actual1)
		assert.Equal(t, "world", actual2)
	})

	t.Run("should get empty or nil values subarray of map", func(t *testing.T) {
		actual1, err1 := TryGetValueByPath("key1[1]", &mapFieldOfPointers)
		actual2, err2 := TryGetValueByPath("key1[2]", &mapFieldOfPointers)
		actual3, err3 := TryGetValueByPath("key1[3]", &mapFieldOfPointers)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Nil(t, err3)
		assert.Equal(t, nil, actual1)
		assert.Equal(t, "", actual2)
		assert.Equal(t, nil, actual3)
	})

	t.Run("should get values from slice", func(t *testing.T) {
		actual0, err0 := TryGetValueByPath("[0]", &nestedSlice)
		assert.Nil(t, err0)
		assert.Equal(t, "hello", actual0)

		actual1, err1 := TryGetValueByPath("[1]", &nestedSlice)
		assert.Nil(t, err1)
		assert.Equal(t, nil, actual1)

		actual2, err2 := TryGetValueByPath("[2]", &nestedSlice)
		assert.Nil(t, err2)
		assert.Equal(t, "", actual2)

		actual3, err3 := TryGetValueByPath("[3]", &nestedSlice)
		assert.Nil(t, err3)
		assert.Equal(t, nil, actual3)

		actual4, err4 := TryGetValueByPath("[4]", &nestedSlice)
		assert.Nil(t, err4)
		assert.Equal(t, "world", actual4)
	})

	t.Run("should get nil from slice non valid index", func(t *testing.T) {
		actual, err := TryGetValueByPath("[10]", &nestedSlice)
		assert.Nil(t, err)
		assert.Equal(t, nil, actual)
	})
}
