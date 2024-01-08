package common

import (
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
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			// Create a new email
			_, err := NewEmail(tc.email)
			// Check if the error matches the expected error
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

		})
	}
}
