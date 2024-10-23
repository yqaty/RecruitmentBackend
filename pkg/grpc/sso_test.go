package grpc

import "testing"

func TestGetUserInfoByUID(t *testing.T) {
	userInfo, err := GetUserInfoByUID("afb6e834-3615-4ebb-9d9d-825af333a3ca")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("Get UserInfo Success")
	t.Logf("%#v", userInfo)
}

func TestGetRolesByUID(t *testing.T) {
	userRoles, err := GetRolesByUID("c4fb1c23-e9de-40a6-b1d4-b4bc2df0a625")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("Get UserRoles Success")
	t.Logf("%#v", userRoles)
}

func TestGetUsers(t *testing.T) {
	users, err := GetUsers([]string{
		"ffb6e834-3615-4ebb-9d9d-825af333a3ca",
		"afb6e834-3615-4ebb-9d9d-825af333a3ca",
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("Get Users Success")
	t.Logf("%#v", users)
}

func TestGetGroupsDetail(t *testing.T) {
	groupsDetail, err := GetGroupsDetail()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("Get GroupsDetail Success")
	t.Logf("%#v", groupsDetail)
}
