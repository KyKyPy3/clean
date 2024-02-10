package common_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/stretchr/testify/assert"
)

func TestEmail_NewEmail(t *testing.T) {
	t.Parallel()

	// Build our needed testcase data struct
	type testCase struct {
		test        string
		email       string
		expectedErr error
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:        "Empty Email validation",
			email:       "",
			expectedErr: common.ErrEmptyEmail,
		}, {
			test:        "Valid Email",
			email:       "peter.parker@email.com",
			expectedErr: nil,
		}, {
			test:        "Invalid Email",
			email:       "peter.parker",
			expectedErr: common.ErrBadFormat,
		},
	}

	for _, tc := range testCases {
		email := tc.email
		expectedErr := tc.expectedErr

		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			t.Parallel()

			// Create a new email
			_, err := common.NewEmail(email)
			// Check if the error matches the expected error
			if !errors.Is(err, expectedErr) {
				t.Errorf("Expected error %v, got %v", expectedErr, err)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, _ := common.NewEmail("peter.parker@gmail.com")

	assert.Equal(t, "peter.parker@gmail.com", email.String())
}

func TestEmail_IsEmpty(t *testing.T) {
	email := common.Email{}

	assert.True(t, email.IsEmpty())
}

func TestEmail_MarshalText(t *testing.T) {
	email, _ := common.NewEmail("peter.parker@gmail.com")
	marshalled, err := email.MarshalText()

	require.NoError(t, err)
	assert.Equal(t, marshalled, []byte("peter.parker@gmail.com"))
}
