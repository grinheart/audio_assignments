package user

import ("database/sql"
_ "github.com/go-sql-driver/mysql"
"errors"
"strconv"
"os"
"log")

type User struct {
	id int
	Name string
	Email string
	Pwd string
	name string
	email string
	pwd string
	db *sql.DB
}

func (u *User) SetDB(db *sql.DB) {
	u.db = db
}

func (u *User) checkDBSet() {
	if (u.db == nil) {
		panic(errors.New("DB not set for User"))
	}
}

func (u *User) GetId() int {
	return u.id;
}

func (u *User) retrieveIdWithQuery(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := u.db.Prepare(query)
	res, err := stmt.Exec(args...)
	if (err != nil) {
		log.Println(err)
		return res, err
	}
	id64, err := res.LastInsertId()
	u.id = int(id64)
	return res, err
}

func (u *User) Reg() (int) {
	u.save()
	res, err := u.db.Query("SELECT * from users WHERE email = ?", u.email)
	if (err != nil) {
		log.Println(err)
	}
	if (res.Next()) {
		return 1
	}
	_, err = u.retrieveIdWithQuery("INSERT INTO users(name, email, pwd) VALUES(?, ?, ?);", u.name, u.email, u.pwd)
	if (err != nil) {
		log.Println(err)
		return -1
	}
	path := "./audio/" + strconv.Itoa(u.id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	if (err != nil) {
		log.Println(err)
		return -1
	}
	return 0
}

func (u *User) save() {
	//todo: validation
	u.name = u.Name
	u.email = u.Email
	u.pwd = u.Pwd
}

func (u *User) Auth() (int) {
	u.save()
	u.checkDBSet()
	log.Print("u.id u.name u.email ", u.id, " ", u.name, " ", u.email)

	res, err := u.db.Query("SELECT id, name, email FROM users WHERE email=? AND pwd=?", u.email, u.pwd)
	if (err != nil) {
		return -1
	} else if res.Next() {
		res.Scan(&u.id, &u.Name, &u.Email)
		u.save()
		log.Print(u.id, " ", u.name, " ", u.email)
		return 0
	}
	return 1
}

func (u *User) Load(id int) (bool, error) {
	res, err := u.db.Query("SELECT name, email, pwd FROM users WHERE id=?;", id)
	log.Println("id in Load", id)
	if (err != nil) {
		log.Println("error loading with ", id)
		return false, err
	}
	if (!res.Next()) {
		log.Println("empty set with ", id)
		u.id = 0
		return false, nil
	}
	u.id = id
	err = res.Scan(&u.Name, &u.Email, &u.Pwd)
	if (err != nil) {
		return false, err
	}
	u.save()
	return true, nil
}