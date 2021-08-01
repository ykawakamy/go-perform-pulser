package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateTopicNameSupplier(t *testing.T) {
	factory := CreateFactory()
	factory.SetDistoribution(2)
	factory.ParseTopicExpression("/test/:1-100")

	assert.Equal(t, 50, len(factory.Build().topics))
	assert.Equal(t, 50, len(factory.Build().topics))
}

func Test_CreateTopicNameSupplier2(t *testing.T) {
	factory := CreateFactory()
	factory.SetDistoribution(3)
	factory.ParseTopicExpression("/test/:1-100")

	tsupplier1 := factory.Build()
	assert.Equal(t, 33, len(tsupplier1.topics))
	tsupplier2 := factory.Build()
	assert.Equal(t, 33, len(tsupplier2.topics))
	tsupplier3 := factory.Build()
	assert.Equal(t, 34, len(tsupplier3.topics))

	assert.NotEqual(t, tsupplier1.topics, tsupplier2.topics)
	assert.NotEqual(t, tsupplier2.topics, tsupplier3.topics)
	assert.NotEqual(t, tsupplier3.topics, tsupplier1.topics)

}

func Test_CreateTopicNameSupplier3(t *testing.T) {
	factory := CreateFactory()
	factory.SetDistoribution(1)
	factory.ParseTopicExpression("/test/:1-100")

	tsupplier1 := factory.Build()
	assert.Equal(t, 100, len(tsupplier1.topics))
	tsupplier2 := factory.Build()
	assert.Equal(t, 100, len(tsupplier2.topics))

	assert.Equal(t, tsupplier1.topics, tsupplier2.topics)
}

func Test_CreateTopicNameSupplier4(t *testing.T) {
	factory := CreateFactory()
	factory.SetDistoribution(2)
	factory.ParseTopicExpression("/test/:1-100")

	tsupplier1 := factory.Build()
	assert.Equal(t, 50, len(tsupplier1.topics))
	tsupplier2 := factory.Build()
	assert.Equal(t, 50, len(tsupplier2.topics))
	tsupplier3 := factory.Build()
	assert.Equal(t, 50, len(tsupplier3.topics))
	tsupplier4 := factory.Build()
	assert.Equal(t, 50, len(tsupplier4.topics))

	assert.Equal(t, tsupplier1.topics, tsupplier3.topics)
	assert.Equal(t, tsupplier2.topics, tsupplier4.topics)

	assert.NotEqual(t, tsupplier1.topics, tsupplier2.topics)
	assert.NotEqual(t, tsupplier3.topics, tsupplier4.topics)

}
