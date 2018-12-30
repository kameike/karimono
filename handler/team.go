package handler

import (
	"database/sql"

	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	"github.com/labstack/echo"
)

var GetTeams = injectDbConn(checkAuth(getTeams))

func getTeams(u model.Account, db *sql.DB, c echo.Context) error {
	teams := readTeams(u, db)
	c.JSON(200, teams)
	return nil
}

func readTeams(u model.Account, db *sql.DB) []model.Team {
	smit, err := db.Prepare(`
		select team.id, team.name from user
		join team on team.id = user.team_id
		where user.account_id = ?
	`)
	util.CheckInternalFatalError(err)

	rows, err := smit.Query(u.Id)
	util.CheckInternalFatalError(err)

	var teams []model.Team
	for rows.Next() {
		var team model.Team
		rows.Scan(&team.Id, &team.Name)
		teams = append(teams, team)
	}

	return teams
}
