package storagedir

import (
	"testing"

	"github.com/anothermemory/mock"
	"github.com/anothermemory/unit"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
)

func TestUnit_Save_NilUnit(t *testing.T) {
	assert.Error(t, newPersistentUnit(nil, location{}, nil).save())
}

func TestPersistentUnit_save_MkdirError(t *testing.T) {
	s := new(mockPersistentUnitStorage)
	s.On("mkdirAll", m.Anything, m.Anything).Return(errors.New("Mkdir error"))

	assert.Error(t, newPersistentUnit(unit.NewUnit(), location{}, s).save())
	s.AssertExpectations(t)
}

func TestPersistentUnit_save_marshalUnitError(t *testing.T) {
	u := new(mock.UnitUnit)
	u.On("MarshalJSON").Return(nil, errors.New("MarshalJSON error"))
	u.On("Type").Return(unit.TypeUnit)

	s := new(mockPersistentUnitStorage)
	s.On("mkdirAll", m.Anything, m.Anything).Return(nil)

	assert.Error(t, newPersistentUnit(u, location{}, s).save())

	s.AssertExpectations(t)
	u.AssertExpectations(t)
}

func TestPersistentUnit_marshalUnit(t *testing.T) {
	u := new(mock.UnitUnit)
	u.On("MarshalJSON").Return(nil, errors.New("MarshalJSON error"))

	bytes, err := newPersistentUnit(u, location{}, nil).marshalUnit(u)
	assert.Nil(t, bytes)
	assert.Error(t, err)

	u.AssertExpectations(t)
}
