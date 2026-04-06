package utility_test

import (
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUuidFromString(t *testing.T) {

	testCases := []struct {
		Name          string
		Value         string
		ExpectedError bool
	}{
		{Name: "Valid UUID", Value: "6cdab6d0-db92-4ed4-8cd1-55b811591cee", ExpectedError: false},
		{Name: "Invalid UUID", Value: "1234", ExpectedError: true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			val, err := utility.UuidFromString(tc.Value)

			if tc.ExpectedError {
				assert.Error(t, err)
				assert.False(t, val.Valid)
			} else {
				assert.NoError(t, err)
				assert.True(t, val.Valid)
			}
		})
	}
}
