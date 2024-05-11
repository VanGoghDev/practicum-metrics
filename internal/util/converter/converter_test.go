package converter

import "testing"

func TestStr(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "int",
			args: args{
				v: 12,
			},
			want:    "12",
			wantErr: false,
		},
		{
			name: "signed int",
			args: args{
				v: -12,
			},
			want:    "-12",
			wantErr: false,
		},
		{
			name: "min int32 val",
			args: args{
				v: -2147483647,
			},
			want:    "-2147483647",
			wantErr: false,
		},
		{
			name: "max int32 val",
			args: args{
				v: 2147483647,
			},
			want:    "2147483647",
			wantErr: false,
		},
		{
			name: "int64 val",
			args: args{
				v: 2147483648,
			},
			want:    "2147483648",
			wantErr: false,
		},
		{
			name: "min int64 val",
			args: args{
				v: -9223372036854775808,
			},
			want:    "-9223372036854775808",
			wantErr: false,
		},
		{
			name: "max int64 val",
			args: args{
				v: 9223372036854775807,
			},
			want:    "9223372036854775807",
			wantErr: false,
		},
		{
			name: "uint",
			args: args{
				v: 255,
			},
			want:    "255",
			wantErr: false,
		},
		{
			name: "not a number",
			args: args{
				v: "sss",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "not a number",
			args: args{
				v: args{
					v: "ss",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "float",
			args: args{
				v: 3216543.53416,
			},
			want:    "3216543.53416",
			wantErr: false,
		},
		{
			name: "signed float",
			args: args{
				v: -3216543.53416,
			},
			want:    "-3216543.53416",
			wantErr: false,
		},
		{
			name: "zero float",
			args: args{
				v: -0.53416,
			},
			want:    "-0.53416",
			wantErr: false,
		},
		{
			name: "zero.zero float",
			args: args{
				v: 0.0,
			},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Str(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Str() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Str() = %v, want %v", got, tt.want)
			}
		})
	}
}
