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
  token text not null unique,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists token_index on access_token(token);

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

const dummyAccountName = "testUser"
const dummyTeamAccounHistoryCount = 3
const dummyAccounHistoryCount = 5
const dummyTeamHistoryCount = 2
const dummyTeamName = "testTeam"
const dummyAccountToken = "tokentoken"

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

func createDummyAccount(r DataRepository) {
	r.InsertAccount(InsertAccountRequest{
		Id:                dummyAccountName,
		EncryptedPassword: "pass",
	})

	r.UpdateOrReplaceAccessToken(UpdateOrReqlaceAccessTokenRequest{
		AccountName: dummyAccountName,
		NewToken:    dummyAccountToken,
	})
}

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

func getDummyAccount(r DataRepository) *model.Me {
	account, _ := r.GetAccountWithSecretInfo(GetAccountRequest{
		Token: dummyAccountToken,
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
