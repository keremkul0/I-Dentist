package helpers

import (
	"dental-clinic-system/models"
	"testing"
)

func TestContainsRole(t *testing.T) {
	type args struct {
		user     models.UserGetModel
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
				user: models.UserGetModel{
					Roles: []*models.Role{{Name: "admin"}, {Name: "user"}},
				},
				roleName: "admin",
			},
			want: true,
		},
		{
			name: "Role does not exist",
			args: args{
				user: models.UserGetModel{
					Roles: []*models.Role{{Name: "user"}},
				},
				roleName: "admin",
			},
			want: false,
		},
		{
			name: "Case insensitive match",
			args: args{
				user: models.UserGetModel{
					Roles: []*models.Role{{Name: "Admin"}},
				},
				roleName: "admin",
			},
			want: true,
		},
		{
			name: "Empty roles",
			args: args{
				user: models.UserGetModel{
					Roles: []*models.Role{},
				},
				roleName: "admin",
			},
			want: false,
		},
		{
			name: "Empty role name",
			args: args{
				user: models.UserGetModel{
					Roles: []*models.Role{{Name: "admin"}},
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
