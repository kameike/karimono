package repository

import (
	"github.com/kameike/karimono/error"
	"testing"
)

func TestAccountTeamRelationWhenNotJoined(t *testing.T) {
	r := inMemoryRepo()
	createDummyAccount(r)
	createDummyTeam(r)

	err := r.CheckAccountTeamRelation(CheckAccountTeamRelationRequest{
		AccountName: dummyAccountName,
		TeamName:    dummyTeamName,
	})

	if invalidError(apperror.ErrorNoRelationBetweenUserAndTeam, err) {
		t.Fail()
	}
}

func TestCreateAccountFailWhenTryToRegisterSameName(t *testing.T) {
	r := inMemoryRepo()

	r.InsertAccount(InsertAccountRequest{
		Id:                "test",
		EncryptedPassword: "pass",
	})
	err := r.InsertAccount(InsertAccountRequest{
		Id:                "test",
		EncryptedPassword: "pass",
	})

	if invalidError(apperror.ErrorAccountNameAlreadyTaken, err) {
		t.Fail()
	}
}

func TestInvalidNameAccount(t *testing.T) {
	r := inMemoryRepo()

	err := r.InsertAccount(InsertAccountRequest{
		Id:                "",
		EncryptedPassword: "pass",
	})

	if invalidError(apperror.ErrorInvalidAccountName, err) {
		t.Fatalf("empty name is not allowed")
	}
}

func TestAccountTeamRelation(t *testing.T) {
	r := inMemoryRepo()
	createDummyTeamAcountWithHistory(r)

	err := r.CheckAccountTeamRelation(CheckAccountTeamRelationRequest{
		AccountName: dummyAccountName,
		TeamName:    dummyTeamName,
	})

	if err != nil {
		t.Fail()
	}
}

func TestCreateAccount(t *testing.T) {
	r := inMemoryRepo()

	r.InsertAccount(InsertAccountRequest{
		Id:                "test",
		EncryptedPassword: "pass",
	})

	count := checkCount("account", r.db())
	if count != 1 {
		t.Fatalf("user count should be 1 but %d", count)
	}
}

func TestFindInvalidAccount(t *testing.T) {
	r := inMemoryRepo()
	createDummyAccount(r)

	_, err := r.GetAccountWithSecretInfo(GetAccountRequest{
		Token: "bad token",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
	}
}

func TestCheckAuthWithAccessToken(t *testing.T) {
	r := inMemoryRepo()

	r.InsertAccount(InsertAccountRequest{
		Id:                "hoge",
		EncryptedPassword: "fuga",
	})

	token := "verylongtokentokentoken"

	r.UpdateOrReplaceAccessToken(UpdateOrReqlaceAccessTokenRequest{
		AccountName: "hoge",
		NewToken:    token,
	})

	r.CheckAuth(AuthCheckRequest{
		AccessToken: token,
	})

	account, _ := r.GetAccountWithSecretInfo(GetAccountRequest{
		Token: token,
	})

	if account.Name != "hoge" {
		t.Fatalf("can not get account form access token")
	}
}

func TestCheckAuthFailRequest(t *testing.T) {
	r := inMemoryRepo()

	err := r.CheckAuth(AuthCheckRequest{
		AccessToken: "wrong",
	})

	if invalidError(apperror.ErrorInvalidAccessToken, err) {
		t.Fail()
	}
}
