package domain

import (
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/repository"
)

type TeamPasswordProvider interface {
	Password() string
}

type TeamCreateRequester interface {
	TeamPasswordProvider
	TeamNameProvider
}

type AccountIdUpdateRequester interface {
	AccountIdProvider
}

type AccountPasswordUpdateRequester interface {
	AccountPasswordProvider
}

type JoinTeamRequester interface {
	TeamNameProvider
	TeamPasswordProvider
}

type LeaveTeamRequestor interface {
	TeamNameProvider
}

type AccountDomain interface {
	CreateTeam(TeamCreateRequester) (*model.Team, error)
	JoinTeam(JoinTeamRequester) (*model.Team, error)
	LeaveTeam(LeaveTeamRequestor) error
	GetHistory() ([]model.Hisotry, error)
	GetBorrowings() ([]model.Borrowing, error)
	UpdateAccountPassword(AccountPasswordProvider) (*model.Me, error)
	UpdateAccountId(AccountIdUpdateRequester) (*model.Me, error)
	GetTeams() ([]model.Team, error)
	GetAccount() *model.Me

	RetrunBorrowingWithHash(string) (*model.Borrowing, error)
}

func createAccountApplicatoinDomain(account model.Me, repo repository.DataRepository) AccountDomain {
	domain := applicationAccountDomain{
		account:    account,
		repository: repo,
	}

	return &domain
}

type applicationAccountDomain struct {
	account    model.Me
	repository repository.DataRepository
}

func (self *applicationAccountDomain) CreateTeam(req TeamCreateRequester) (*model.Team, error) {
	r := self.repository
	r.BeginTransaction()
	defer r.CancelTransaction()

	err := r.CreateTeam(repository.CreateTeamRequest{
		Name:              req.TeamName(),
		EncryptedPassword: hashPassword(req.Password()),
	})

	if err != nil {
		return nil, err
	}

	team, err := r.GetTeam(repository.GetTeamRequest{
		TeamName: req.TeamName(),
	})

	if err != nil {
		return nil, err
	}

	err = r.CreateTeamAccountReleation(repository.CreateTeamAccountReleationRequest{
		TeamName:    team.Name,
		AccountName: self.account.Name,
	})

	if err != nil {
		return nil, err
	}

	r.EndTransaction()
	return team, nil
}

func (self *applicationAccountDomain) RetrunBorrowingWithHash(hash string) (*model.Borrowing, error) {
	r := self.repository

	b, err := r.GetBorrowing(repository.GetBorrowingRequest{
		Id: hash,
	})

	if err != nil {
		return nil, err
	}

	err = r.ReturnBorrowing(repository.ReturnBorrowingRequest{
		BorrowingId: hash,
	})

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (self *applicationAccountDomain) JoinTeam(req JoinTeamRequester) (*model.Team, error) {
	r := self.repository

	println(req.TeamName())
	println(req.Password())

	targetHash, err := r.GetTeamPasswordHash(repository.GetTeamPasswordHashRequest{
		TeamName: req.TeamName(),
	})
	if err != nil {
		return nil, err
	}

	err = checkPasswordHash(req.Password(), targetHash)
	if err != nil {
		return nil, err
	}

	err = r.CreateTeamAccountReleation(repository.CreateTeamAccountReleationRequest{
		TeamName:    req.TeamName(),
		AccountName: self.account.Name,
	})

	if err != nil {
		return nil, err
	}

	team, err := r.GetTeam(repository.GetTeamRequest{
		TeamName: req.TeamName(),
	})

	return team, err
}

func (self *applicationAccountDomain) LeaveTeam(req LeaveTeamRequestor) error {
	r := self.repository

	err := r.DeleteTeamAccountReleation(repository.DeleteTeamAccountReleationRequest{
		TeamName:    req.TeamName(),
		AccountName: self.account.Name,
	})

	return err
}

func (self *applicationAccountDomain) GetHistory() ([]model.Hisotry, error) {
	r := self.repository

	hisotries, err := r.GetAccountHistory(repository.GetAccountHistoryRequest{
		AccountName: self.account.Name,
	})

	return hisotries, err
}

func (self *applicationAccountDomain) GetBorrowings() ([]model.Borrowing, error) {
	r := self.repository

	borrowings, err := r.GetAccountBorrowing(repository.GetAccountBorrowingRequset{
		AccountName: self.account.Name,
	})

	return borrowings, err
}
func (self *applicationAccountDomain) UpdateAccountPassword(req AccountPasswordProvider) (*model.Me, error) {
	r := self.repository
	r.BeginTransaction()
	defer r.CancelTransaction()

	err := r.UpdateAccountPassword(repository.UpdateAccountPasswordRequest{
		AccountName:    self.account.Name,
		HashedPassword: hashPassword(req.AccountPassword()),
	})

	if err != nil {
		return nil, err
	}

	account, err := r.GetAccountWithSecretInfo(repository.GetAccountRequest{
		Token: self.account.Token,
	})

	if err != nil {
		return nil, err
	}

	token := newToken()

	r.UpdateOrReplaceAccessToken(repository.UpdateOrReqlaceAccessTokenRequest{
		AccountName: self.account.Name,
		NewToken:    token,
	})

	self.account = *account

	r.EndTransaction()
	return account, nil
}

func (self *applicationAccountDomain) UpdateAccountId(req AccountIdUpdateRequester) (*model.Me, error) {
	r := self.repository

	err := r.UpdateAccountId(repository.UpdateAccountIdRequest{
		OldAccountName: self.account.Name,
		NewAccountName: req.AccountId(),
	})
	if err != nil {
		return nil, err
	}

	account, err := r.GetAccountWithSecretInfo(repository.GetAccountRequest{
		Token: self.account.Token,
	})

	if err != nil {
		return nil, err
	}

	self.account = *account

	return account, nil
}

func (self *applicationAccountDomain) GetAccount() *model.Me {
	return &self.account
}

func (self *applicationAccountDomain) GetTeams() ([]model.Team, error) {
	r := self.repository

	teams, err := r.GetTeams(repository.GetTeamsRequest{
		TeamName: self.account.Name,
	})

	return teams, err
}
