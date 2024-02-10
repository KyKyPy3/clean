package vo_test

import (
	"errors"
	"testing"

	"github.com/KyKyPy3/clean/internal/modules/user/domain/vo"
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
			expectedErr: vo.ErrEmptyFirstName,
		}, {
			test:        "Empty last name validation",
			firstName:   "John",
			lastName:    "",
			middleName:  "Sr",
			expectedErr: nil,
		}, {
			test:        "Valid fullname",
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
			_, err := vo.NewFullName(tc.firstName, tc.lastName, tc.middleName)
			// Check if the error matches the expected error
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
