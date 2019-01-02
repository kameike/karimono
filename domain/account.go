package domain

import (
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/repository"
)

// Interfaces
type TeamIdProvider interface {
	TeamId() string
}

type TeamNameProvider interface {
	TeamName() string
}

type TeamPasswordProvider interface {
	Password() string
}

type TeamCreateRequester interface {
	TeamPasswordProvider
	TeamIdProvider
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
	UpdateAccountPassword(AccountPasswordProvider) (*model.Account, error)
	UpdateAccountId(AccountIdProvider) (*model.Account, error)
	GetAccount() *model.Account
}

func createAccountApplicatoinDomain(account model.Account, repo repository.DataRepository) AccountDomain {
	domain := applicationAccountDomain{
		account:    account,
		repository: repo,
	}

	return &domain
}

type applicationAccountDomain struct {
	account    model.Account
	repository repository.DataRepository
}

func (self *applicationAccountDomain) CreateTeam(req TeamCreateRequester) (*model.Team, error) {
	r := self.repository
	r.BeginTransaction()

	err := r.CreateTeam(repository.CreateTeamRequest{
		Name:              req.TeamId(),
		EncryptedPassword: hashPassword(req.Password()),
	})

	if err != nil {
		return nil, err
	}

	team, err := r.GetTeam(repository.GetTeamRequest{
		TeamName: req.TeamId(),
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

func (self *applicationAccountDomain) JoinTeam(req JoinTeamRequester) (*model.Team, error) {
	r := self.repository

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
func (self *applicationAccountDomain) UpdateAccountPassword(req AccountPasswordProvider) (*model.Account, error) {
	r := self.repository
	r.BeginTransaction()

	account, err := r.UpdateAccountPassword(repository.UpdateAccountPasswordRequest{
		AccountName:    self.account.Name,
		HashedPassword: hashPassword(req.AccountPassword()),
	})

	if err != nil {
		return nil, err
	}

	token := newToken()

	r.UpdateOrReplaceAccessToken(repository.UpdateOrReqlaceAccessTokenRequest{
		AccountName: self.account.Name,
		NewToken:    token,
	})

	account.Token = token
	self.account = *account

	r.EndTransaction()
	return account, nil
}

func (self *applicationAccountDomain) UpdateAccountId(req AccountIdProvider) (*model.Account, error) {
	r := self.repository

	account, err := r.UpdateAccountId(repository.UpdateAccountIdRequest{
		OldAccountName: self.account.Token,
		NewAccountName: req.AccountId(),
	})

	if err != nil {
		return nil, err
	}

	self.account = *account

	return account, nil
}

func (self *applicationAccountDomain) GetAccount() *model.Account {
	return &self.account
}
