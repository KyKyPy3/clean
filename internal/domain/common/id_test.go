package common_test

import (
	"errors"
	"testing"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEntityID_NewEntityID(t *testing.T) {
	entity := common.NewUID()

	assert.NotNil(t, entity.GetID())
}

func TestEntityID_NewEntityWithSpecifiedID(t *testing.T) {
	id := uuid.New()
	entity := common.NewWithSpecifiedID(id)

	assert.Equal(t, entity.GetID(), id)
}

func TestEntityID_ParseEntityID(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test        string
		id          string
		expectedErr error
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:        "Empty entity id validation",
			id:          "",
			expectedErr: common.ErrEmptyUID,
		}, {
			test:        "Valid entity id",
			id:          "2b0c8791-2136-46b6-bc38-b33038ca2e80",
			expectedErr: nil,
		}, {
			test:        "Invalid entity id",
			id:          "123",
			expectedErr: common.ErrUIDBadFormat,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			// Create a new entity id
			_, err := common.ParseUID(tc.id)
			// Check if the error matches the expected error
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
