package echo

import "github.com/arashi5/echo/internal/repository/echo"

type echoRepo echo.Echo

func (e echoRepo) Dto() Echo {
	return Echo{
		Id:       e.ID,
		Title:    e.Title,
		Reminder: e.Reminder,
	}
}

func (e Echo) Dto() echo.Echo {
	return echo.Echo{
		ID:       e.Id,
		Title:    e.Title,
		Reminder: e.Reminder,
	}
}
