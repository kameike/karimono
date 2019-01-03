package repository

import (
	"testing"

	apperror "github.com/kameike/karimono/error"
)

func TestCreateTeamAccountHistoryRequest(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.CreateTeamAccountHistory(CreateTeamAccountHistoryRequest{
		AccountName: dummyAccountName,
		TeamName:    dummyTeamName,
		History:     "borrow xxx",
	})

	assertCountEqual(t, r.db(), "history", 1)
}

func TestCreateInvalidTeamAccountHistoryRequest(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	err := r.CreateTeamAccountHistory(CreateTeamAccountHistoryRequest{
		AccountName: dummyAccountName,
		TeamName:    "invalid team",
		History:     "borrow xxx",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}

	err = r.CreateTeamAccountHistory(CreateTeamAccountHistoryRequest{
		AccountName: "invalid account",
		TeamName:    dummyTeamName,
		History:     "borrow xxx",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}

	assertCountEqual(t, r.db(), "history", 0)
}

func TestCreateTeamHistory(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.CreateTeamHistory(CreateTeamHistoryRequest{
		History:  "hoge",
		TeamName: dummyTeamName,
	})

	assertCountEqual(t, r.db(), "history", 1)
}

func TestCreateInvalidTeamHistory(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	err := r.CreateTeamHistory(CreateTeamHistoryRequest{
		History:  "message",
		TeamName: "invalid team",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}

	assertCountEqual(t, r.db(), "history", 0)
}

func TestCreateInvalidAccountHistory(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	err := r.CreateAccountHistory(CreateAccountHistoryRequest{
		History:     "message",
		AccountName: "invalid team",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}

	assertCountEqual(t, r.db(), "history", 0)
}

func TestCreateAccountHistory(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.CreateAccountHistory(CreateAccountHistoryRequest{
		AccountName: dummyAccountName,
		History:     "history",
	})

	assertCountEqual(t, r.db(), "history", 1)
}

func TestGetTeamHisotry(t *testing.T) {
	r := inMemoryRepo()
	createDummyTeamAcountWithHistory(r)

	result, _ := r.GetTeamHistory(GetTeamHistoryRequest{
		TeamName: dummyTeamName,
	})

	size := len(result)

	if size != dummyTeamHistoryCount+dummyTeamAccounHistoryCount {
		t.Fatalf("size should be %d, but got %d", dummyTeamHistoryCount+dummyTeamAccounHistoryCount, size)
	}
}

func TestGetAccountHistory(t *testing.T) {
	r := inMemoryRepo()
	createDummyTeamAcountWithHistory(r)

	result, _ := r.GetAccountHistory(GetAccountHistoryRequest{
		AccountName: dummyAccountName,
	})

	size := len(result)
	expected := dummyAccounHistoryCount + dummyTeamAccounHistoryCount

	if size != expected {
		t.Fatalf("count data should be %d but got %d", expected, size)
	}
}
