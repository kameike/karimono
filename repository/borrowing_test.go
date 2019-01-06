package repository

import "testing"

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

func TestGetBorrowing(t *testing.T) {
	r := inMemoryRepo()
	dummyAccountJoinToDummyTeam(r)

	id := "borrowingLognHash"

	r.CreateBorrowing(CreateBorrowingRequest{
		BorrowingId:      id,
		BorrwoedTeam:     dummyTeamName,
		BorrwoingAccount: dummyAccountName,
		ItemName:         "any thing",
	})

	result, _ := r.GetBorrowing(GetBorrowingRequest{
		Id: id,
	})

	if result.Uuid != id {
		t.Fatalf("'%s' is not '%s'", result.Uuid, id)
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
