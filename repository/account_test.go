package repository

import (
	"github.com/kameike/karimono/error"
	"testing"
)

func TestUpdateAccountId(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	newName := "newName"

	r.UpdateAccountId(UpdateAccountIdRequest{
		OldAccountName: dummyAccountName,
		NewAccountName: newName,
	})

	account, _ := r.GetAccountWithSecretInfo(GetAccountRequest{
		Token: dummyAccountToken,
	})

	if account.Name != "newName" {
		t.Fatalf("acount name should be %s but get %s", newName, account.Name)
	}
}

func TestUpdateAccountPassword(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	newPassword := "newpassword"

	r.UpdateAccountPassword(UpdateAccountPasswordRequest{
		HashedPassword: newPassword,
		AccountName:    dummyAccountName,
	})

	account, _ := r.GetAccountWithSecretInfo(GetAccountRequest{
		Token: dummyAccountToken,
	})

	if account.PasswordHash != newPassword {
		t.Fatalf("fail to update password")
	}
}

func TestCreateOrReqlaceAccessToken(t *testing.T) {
	r := inMemoryRepo()

	r.InsertAccount(InsertAccountRequest{
		Id:                "hoge",
		EncryptedPassword: "fuga",
	})

	r.UpdateOrReplaceAccessToken(UpdateOrReqlaceAccessTokenRequest{
		AccountName: "hoge",
		NewToken:    "verylongtoken",
	})

	count := checkCount("access_token", r.db())
	if count != 1 {
		t.Fatalf("count should be 1 but result is %d", count)
	}
}

func TestCreateTeam(t *testing.T) {
	r := inMemoryRepo()

	createDummyTeam(r)

	count := checkCount("team", r.db())
	if count != 1 {
		t.Fatalf("count should be 1 but %d", count)
	}
}

func TestTeamNameShouldBeUniqe(t *testing.T) {
	r := inMemoryRepo()

	createDummyTeam(r)
	err := createDummyTeam(r)

	count := checkCount("team", r.db())
	if count != 1 {
		t.Fatalf("count should be 1 but %d", count)
	}

	if invalidError(apperror.ErrorTeamNameAlreadyTaken, err) {
		t.Fail()
	}
}

func TestTeamNotFound(t *testing.T) {
	r := inMemoryRepo()
	createDummyTeam(r)

	_, err := r.GetTeam(GetTeamRequest{
		TeamName: "empty team",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}
}

func TestFindTeam(t *testing.T) {
	r := inMemoryRepo()
	createDummyTeam(r)

	team, _ := r.GetTeam(GetTeamRequest{
		TeamName: dummyTeamName,
	})

	if team.Name != dummyTeamName {
		t.Fail()
	}
}

func TestUpdateNotExistUserPassword(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	newPassword := "newpassword"

	err := r.UpdateAccountPassword(UpdateAccountPasswordRequest{
		HashedPassword: newPassword,
		AccountName:    "noman",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}
}

func TestUpdateAccountIdWhichAlreadyTaken(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	newName := "newName"

	r.InsertAccount(InsertAccountRequest{
		EncryptedPassword: "password",
		Id:                newName,
	})

	err := r.UpdateAccountId(UpdateAccountIdRequest{
		OldAccountName: dummyAccountName,
		NewAccountName: newName,
	})

	if invalidError(apperror.ErrorAccountNameAlreadyTaken, err) {
		t.Fail()
	}
}
