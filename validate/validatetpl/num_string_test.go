package validatetpl

import "testing"

func TestValidateNumString(t *testing.T) {
	type args struct {
		v    interface{}
		ok   bool
		min  int
		max  int
		name string
	}

	testCases := []args{
		{
			name: "not a num string",
			v:    "31qadads",
			min:  0,
			max:  100,
			ok:   false,
		},
		{
			name: "with white space",
			v:    " 123456 ",
			min:  0,
			max:  100,
			ok:   false,
		},
		{
			name: "min len num string",
			v:    "1",
			min:  1,
			max:  3,
			ok:   true,
		},
		{
			name: "max len num string",
			min:  1,
			max:  3,
			v:    "123",
			ok:   true,
		},
		{
			name: "too small num string",
			min:  10,
			max:  30,
			v:    "123",
			ok:   false,
		},
		{
			name: "too big num string",
			min:  1,
			max:  4,
			v:    "123456",
			ok:   false,
		},
		{
			name: "valid num string",
			min:  1,
			max:  10,
			v:    "12345",
			ok:   true,
		},
	}

	for _, tc := range testCases {
		if ok, _ := NewValidateNumStringLength(tc.min, tc.max)(tc.v); ok != tc.ok {
			t.Errorf("ValidateBankCard input:%s name:%s got:%v want:%v", tc.v, tc.name, ok, tc.ok)
		}
	}
}
