package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type Dao struct {
	db *sql.DB
}

func NewDao(db *sql.DB) *Dao {
	return &Dao{
		db: db,
	}
}

func (d *Dao) GetUserData(account string) (user UserData, err error) {
	query := "select user_name,birth_day,gender from t_user where account=?"
	e := d.db.QueryRow(query, account).Scan(&user.UserName, &user.Birthday, &user.Gender)
	err = errors.Wrap(e, "dao.GetUserData")
	return
}
