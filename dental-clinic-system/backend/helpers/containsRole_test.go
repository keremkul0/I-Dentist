package helpers

import (
	"dental-clinic-system/models/user"
	"testing"
)

func TestContainsRole(t *testing.T) {
	type args struct {
		user     user.UserGetModel
		roleName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Role exists",
			args: args{
				user: user.UserGetModel{
					Roles: []*user.Role{{Name: "admin"}, {Name: "user"}},
				},
				roleName: "admin",
			},
			want: true,
		},
		{
			name: "Role does not exist",
			args: args{
				user: user.UserGetModel{
					Roles: []*user.Role{{Name: "user"}},
				},
				roleName: "admin",
			},
			want: false,
		},
		{
			name: "Case insensitive match",
			args: args{
				user: user.UserGetModel{
					Roles: []*user.Role{{Name: "Admin"}},
				},
				roleName: "admin",
			},
			want: true,
		},
		{
			name: "Empty roles",
			args: args{
				user: user.UserGetModel{
					Roles: []*user.Role{},
				},
				roleName: "admin",
			},
			want: false,
		},
		{
			name: "Empty role name",
			args: args{
				user: user.UserGetModel{
					Roles: []*user.Role{{Name: "admin"}},
				},
				roleName: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsRole(tt.args.user, tt.args.roleName); got != tt.want {
				t.Errorf("ContainsRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
