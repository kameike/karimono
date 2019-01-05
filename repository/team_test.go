package repository

import (
	"testing"

	apperror "github.com/kameike/karimono/error"
)

func TestInvalidJoin(t *testing.T) {
	r := inMemoryRepo()

	err := r.CreateTeamAccountReleation(CreateTeamAccountReleationRequest{
		AccountName: "kameike",
		TeamName:    "pixiv",
	})

	count := checkCount("account_team", r.db())
	if count != 0 {
		t.Fatalf("count should be 0 but got %d", count)
	}

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}
}

func TestPreventDoubleJoinTeam(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.CreateTeam(CreateTeamRequest{
		Name:              "newTeam",
		EncryptedPassword: "hoge",
	})

	r.CreateTeamAccountReleation(CreateTeamAccountReleationRequest{
		AccountName: dummyAccountName,
		TeamName:    "newTeam",
	})

	err := r.CreateTeamAccountReleation(CreateTeamAccountReleationRequest{
		AccountName: dummyAccountName,
		TeamName:    "newTeam",
	})

	assertCountEqual(t, r.db(), "account_team", 2)

	if invalidError(apperror.ErrorAlreadyJoin, err) {
		t.Fail()
	}
}

func TestJoinTeam(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	assertCountEqual(t, r.db(), "account_team", 1)
}

func TestLeaveTeam(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.DeleteTeamAccountReleation(DeleteTeamAccountReleationRequest{
		AccountName: dummyAccountName,
		TeamName:    dummyTeamName,
	})

	assertCountEqual(t, r.db(), "account_team", 0)
}

func TestGetTeam(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	teams, _ := r.GetTeams(GetTeamsRequest{
		TeamName: dummyAccountName,
	})

	if len(teams) != 1 {
		t.Fail()
	}
}

func TestGetTeamPassowrd(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	password, _ := r.GetTeamPasswordHash(GetTeamPasswordHashRequest{
		TeamName: dummyTeamName,
	})

	if password == "" {
		t.Fail()
	}
}
