package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

/*
	问题：我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
	答案：应该将错误Wrap给上层，这样才符合错误只处理一次，避免多处判断if != nil和打印日志，在这里GetUserRecord函数调用QueryRow方法查询记录，如果没有返回行时，
		将错误暂存给*Row结构体，同时结构体内有error字段，这个自动将内置错误sql.ErrNoRows；在上层main里，调用GetUserRecord，进行错误判断，判断是否等于sql.ErrNoRows，
		这样才符合错误只处理一次，GetUserRecord里Wrap是将错误进行重新包装，同时保留原有属性，在顶层可以用%+v打印调用栈信息。
*/

type Dao struct {
	db *sql.DB
}

func (d *Dao) GetUserRecord(name string) error {
	query := "select * from users where name = ?"
	e := d.db.QueryRow(query, name).Scan(&name)
	err := errors.Wrap(e, "Dao: GetUserRecord")
	return err
}

func main() {
	conn := "root:root@tcp(localhost:3306)/go_geekang?charset=utf8"
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Panic(err)
	}
	log.Println("数据库连接成功")
	name := "king1"
	dao := Dao{db}
	err = dao.GetUserRecord(name)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("users表没有找到[%s]记录, 错误信息: %+v\n", name, err)
		return
	}

	log.Printf("users表没存在[%s]记录\n", name)
}
