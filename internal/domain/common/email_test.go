package common

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmail_NewEmail(t *testing.T) {
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
			expectedErr: ErrEmptyEmail,
		}, {
			test:        "Valid Email",
			email:       "peter.parker@email.com",
			expectedErr: nil,
		}, {
			test:        "Invalid Email",
			email:       "peter.parker",
			expectedErr: ErrBadFormat,
		},
	}

	for _, tc := range testCases {
		email := tc.email
		expectedErr := tc.expectedErr

		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			t.Parallel()

			// Create a new email
			_, err := NewEmail(email)
			// Check if the error matches the expected error
			if !errors.Is(err, expectedErr) {
				t.Errorf("Expected error %v, got %v", expectedErr, err)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, _ := NewEmail("peter.parker@gmail.com")

	assert.Equal(t, email.String(), "peter.parker@gmail.com")
}

func TestEmail_IsEmpty(t *testing.T) {
	email := Email{}

	assert.Equal(t, email.IsEmpty(), true)
}

func TestEmail_MarshalText(t *testing.T) {
	email, _ := NewEmail("peter.parker@gmail.com")
	marshalled, err := email.MarshalText()

	assert.Nil(t, err)
	assert.Equal(t, marshalled, []byte("peter.parker@gmail.com"))
}
