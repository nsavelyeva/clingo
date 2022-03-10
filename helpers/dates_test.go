package helpers

import (
	"testing"
	"time"
)

// Verify the returned value of GetMonthDay()
// if zero, positive and negative offset is given to the provided date,
// confirm that leading zeros are displayed when needed.
func TestGetDayMonth(t *testing.T) {
	type args struct {
		tm time.Time
		in int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Custom date with leading zeros",
			args: args{tm: time.Date(2022, time.March, 8, 23, 12, 5, 3, time.UTC), in: 0},
			want: "03-08",
		},
		{
			name: "Custom date without leading zeros",
			args: args{tm: time.Date(2022, time.December, 15, 23, 12, 5, 3, time.UTC), in: 0},
			want: "12-15",
		},
		{
			name: "Future date with leading zeros",
			args: args{tm: time.Date(2022, time.March, 8, 23, 12, 5, 3, time.UTC), in: 30},
			want: "04-07",
		},
		{
			name: "Future date without leading zeros",
			args: args{tm: time.Date(2022, time.November, 28, 23, 12, 5, 3, time.UTC), in: 18},
			want: "12-16",
		},
		{
			name: "Past date with leading zeros (offset should be ignored)",
			args: args{tm: time.Date(2022, time.March, 8, 23, 12, 5, 3, time.UTC), in: -30},
			want: "03-08",
		},
		{
			name: "Past date without leading zeros (offset should be ignored)",
			args: args{tm: time.Date(2022, time.November, 28, 23, 12, 5, 3, time.UTC), in: -18},
			want: "11-28",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMonthDay(tt.args.tm, tt.args.in); got != tt.want {
				t.Errorf("GetMonthDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
