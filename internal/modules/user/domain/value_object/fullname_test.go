package value_object

import (
	"testing"
)

func TestFullName_NewFullName(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		test        string
		firstName   string
		lastName    string
		middleName  string
		expectedErr error
	}

	// Create new test cases
	testCases := []testCase{
		{
			test:        "Empty first name validation",
			firstName:   "",
			lastName:    "Smith",
			middleName:  "Sr",
			expectedErr: ErrEmptyFirstName,
		}, {
			test:        "Empty last name validation",
			firstName:   "John",
			lastName:    "",
			middleName:  "Sr",
			expectedErr: ErrEmptyLastName,
		}, {
			test:        "Valid fullanme",
			firstName:   "John",
			lastName:    "Smith",
			middleName:  "Sr",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			// Create a new email
			_, err := NewFullName(tc.firstName, tc.lastName, tc.middleName)
			// Check if the error matches the expected error
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

		})
	}
}
