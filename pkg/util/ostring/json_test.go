package ostring

import "testing"

func TestInterfaceIsNil(t *testing.T) {
	type args struct {
		a interface{}
	}
	var nilStruct *args
	s := &args{}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil struct is nil",
			args: args{a: nilStruct},
			want: true,
		},
		{
			name: "nil is nil",
			args: args{a: nil},
			want: true,
		},
		{
			name: "struct is not nil",
			args: args{a: s},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceIsNil(tt.args.a); got != tt.want {
				t.Errorf("InterfaceIsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}
