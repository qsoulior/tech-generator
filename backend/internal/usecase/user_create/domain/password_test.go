package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       Password
		wantErr error
	}{
		{
			name:    "Valid_MinLength_HasPunct",
			p:       "aB1-----",
			wantErr: nil,
		},
		{
			name:    "Valid_MaxLength_HasSymbol",
			p:       Password("aB1" + strings.Repeat("$", 69)),
			wantErr: nil,
		},
		{
			name:    "Invalid_TooShort",
			p:       "aB1----",
			wantErr: ErrPasswordTooShort,
		},
		{
			name:    "Invalid_TooLong",
			p:       Password("aB1" + strings.Repeat("$", 70)),
			wantErr: ErrPasswordTooLong,
		},
		{
			name:    "Invalid_NoDigits",
			p:       "aB------",
			wantErr: ErrPasswordNoDigits,
		},
		{
			name:    "Invalid_NoLettersUppercase",
			p:       "a-1------",
			wantErr: ErrPasswordNoLettersUppercase,
		},
		{
			name:    "Invalid_NoLettersLowercase",
			p:       "-B1------",
			wantErr: ErrPasswordNoLettersLowercase,
		},
		{
			name:    "Invalid_NoSpecial",
			p:       "aB100000",
			wantErr: ErrPasswordNoSpecial,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.Validate()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
