package repository

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	_ "github.com/mattn/go-sqlite3"
)

func inMemoryRepo() *applicationDataRepository {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		panic("faild to make repository")
	}

	query := `
create table if not exists account_team (
  id integer primary key autoincrement,
  team_id integer not null,
  account_id integer not null,
  created_at text default (datetime('now', 'localtime')),
	unique(team_id, account_id)
);
create index if not exists user_index on account_team(team_id, id, account_id);

create table if not exists account (
  id integer primary key autoincrement,
  name text not null unique,
  password_hash text not null,
  created_at text default (datetime('now', 'localtime'))
);

create table if not exists access_token(
  id integer primary key autoincrement,
  account_id integer not null unique,
  session_token text not null unique,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists token_index on access_token(session_token);

create table if not exists team (
  id integer primary key autoincrement,
  name text not null unique,
  password_hash text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists team_index on team(id);

create table if not exists borrowing(
  id integer primary key autoincrement,
	account_id integer not null,
	team_id integer not null,
	hashed_id text not null unique,
  name text not null,
  memo text not null,
  has_return text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists borrowing_index on borrowing(account_id, team_id, hashed_id, id);

create table if not exists history(
  id integer primary key autoincrement,
  team_id integer,
	account_id integer,
	text text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists history_index on history(team_id);
	`

	db.Exec(query)

	return &applicationDataRepository{
		_db: db,
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

	account, _ := r.CheckAuth(AuthCheckRequest{
		AccessToken: token,
	})

	if account.Name != "hoge" {
		t.Fatalf("can not get account form access token")
	}
}

func TestCheckAuthFailRequest(t *testing.T) {
	r := inMemoryRepo()

	_, err := r.CheckAuth(AuthCheckRequest{
		AccessToken: "wrong",
	})

	if invalidError(apperror.ErrorInvalidAccessToken, err) {
		t.Fail()
	}
}

func TestGetAccount(t *testing.T) {
	r := inMemoryRepo()
	createDummyAccount(r)

	account, _ := r.GetAccount(GetAccountRequest{
		AccountName: dummyAccountName,
	})

	if account.Name != dummyAccountName {
		t.Fail()
	}
}

func TestGetNotExistAccount(t *testing.T) {
	r := inMemoryRepo()
	createDummyAccount(r)

	_, err := r.GetAccount(GetAccountRequest{
		AccountName: "noman",
	})

	if invalidError(apperror.ErrorDataNotFount, err) {
		t.Fail()
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

func TestUpdateAccountId(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	newName := "newName"

	r.UpdateAccountId(UpdateAccountIdRequest{
		OldAccountName: dummyAccountName,
		NewAccountName: newName,
	})

	account, _ := r.GetAccount(GetAccountRequest{
		AccountName: newName,
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

	account, _ := r.GetAccount(GetAccountRequest{
		AccountName: dummyAccountName,
	})

	if account.PasswordHash != newPassword {
		t.Fatalf("fail to update password")
	}
}

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

func TestCreateAccountBorrowing(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.CreateBorrowing(CreateBorrowingRequest{
		BorrowingId:      "borrowinghoge",
		BorrwoedTeam:     dummyTeamName,
		BorrwoingAccount: dummyAccountName,
		ItemName:         "any thing",
	})

	assertCountEqual(t, r.db(), "borrowing", 1)
}

func TestGetAccountBorrowing(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.CreateBorrowing(CreateBorrowingRequest{
		BorrowingId:      "borrowinghoge",
		BorrwoedTeam:     dummyTeamName,
		BorrwoingAccount: dummyAccountName,
		ItemName:         "any thing",
	})
	r.CreateBorrowing(CreateBorrowingRequest{
		BorrowingId:      "fugafuga",
		BorrwoedTeam:     dummyTeamName,
		BorrwoingAccount: dummyAccountName,
		ItemName:         "any thing",
	})

	result, _ := r.GetAccountBorrowing(GetAccountBorrowingRequset{
		AccountName: dummyAccountName,
	})

	if len(result) != 2 {
		t.Fatalf("%d should be %d", len(result), 2)
	}
}

func TestGetTeamBorrowing(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	r.CreateBorrowing(CreateBorrowingRequest{
		BorrowingId:      "borrowinghoge",
		BorrwoedTeam:     dummyTeamName,
		BorrwoingAccount: dummyAccountName,
		ItemName:         "any thing",
	})
	r.CreateBorrowing(CreateBorrowingRequest{
		BorrowingId:      "fugafuga",
		BorrwoedTeam:     dummyTeamName,
		BorrwoingAccount: dummyAccountName,
		ItemName:         "any thing",
	})

	result, _ := r.GetTeamBorrowing(GetTeamBorrowingRequest{
		TeamName: dummyTeamName,
	})

	if len(result) != 2 {
		t.Fatalf("%d should be %d", len(result), 2)
	}

}

func TestReturnItem(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	borrowingId := "hashsed_id"

	r.CreateBorrowing(CreateBorrowingRequest{
		BorrowingId:      borrowingId,
		BorrwoedTeam:     dummyTeamName,
		BorrwoingAccount: dummyAccountName,
		ItemName:         "any thing",
	})

	r.ReturnBorrowing(ReturnBorrowingRequest{
		BorrowingId: borrowingId,
	})

	result, _ := r.GetAccountBorrowing(GetAccountBorrowingRequset{
		AccountName: dummyAccountName,
	})

	if len(result) != 0 {
		t.Fail()
	}
}

// ===== STUBMING =====

const dummyAccountName = "testUser"
const dummyTeamAccounHistoryCount = 3
const dummyAccounHistoryCount = 5
const dummyTeamHistoryCount = 2

func assertCountEqual(t *testing.T, db queryExecter, tableName string, expectCount int) {
	count := checkCount(tableName, db)
	if count != expectCount {
		t.Fatalf("count rows of %s should be %d but got %d", tableName, expectCount, count)
	}
}

func createDummyTeamAcountWithHistory(r DataRepository) {
	dummyAccountJoinToDummyTeam(r)

	for i := 0; ; {
		r.CreateTeamHistory(CreateTeamHistoryRequest{
			TeamName: dummyTeamName,
			History:  "test",
		})

		i += 1
		if i >= dummyTeamHistoryCount {
			break
		}
	}

	for i := 0; ; {
		r.CreateTeamAccountHistory(CreateTeamAccountHistoryRequest{
			AccountName: dummyAccountName,
			TeamName:    dummyTeamName,
			History:     "test",
		})

		i += 1
		if i >= dummyTeamAccounHistoryCount {
			break
		}
	}

	for i := 0; ; {
		r.CreateAccountHistory(CreateAccountHistoryRequest{
			AccountName: dummyAccountName,
			History:     "test",
		})

		i += 1
		if i >= dummyAccounHistoryCount {
			break
		}
	}
}

func dummyAccountJoinToDummyTeam(r DataRepository) {
	createDummyTeam(r)
	createDummyAccount(r)

	account := getDummyAccount(r)
	team := getDummyTeam(r)

	r.CreateTeamAccountReleation(CreateTeamAccountReleationRequest{
		AccountName: account.Name,
		TeamName:    team.Name,
	})
}

func createDummyAccount(r DataRepository) error {
	return r.InsertAccount(InsertAccountRequest{
		Id:                dummyAccountName,
		EncryptedPassword: "pass",
	})
}

const dummyTeamName = "testTeam"

func createDummyTeam(r DataRepository) error {
	err := r.CreateTeam(CreateTeamRequest{
		Name:              dummyTeamName,
		EncryptedPassword: "hoge",
	})
	return err
}

func getDummyTeam(r DataRepository) *model.Team {
	team, _ := r.GetTeam(GetTeamRequest{
		TeamName: dummyTeamName,
	})

	return team
}

func getDummyAccount(r DataRepository) *model.Account {
	account, _ := r.GetAccount(GetAccountRequest{
		AccountName: dummyAccountName,
	})

	return account
}

func invalidError(code int, err error) bool {
	e, ok := err.(apperror.ApplicationError)
	if !ok {
		return true
	}

	return e.Code != code
}

func checkCount(table string, db queryExecter) int {
	row := db.QueryRow(fmt.Sprintf("select count(1) from %s", table))
	var count int
	err := row.Scan(&count)

	if err != nil {
		print(err.Error())
		panic("fail to read table")
	}

	return count
}
