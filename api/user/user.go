package user

import ("database/sql"
_ "github.com/go-sql-driver/mysql"
"errors"
"strconv"
"os")

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

func (u *User) Reg() error {
	u.save()
	stmt, err := u.db.Prepare("INSERT INTO users(name, email, pwd) VALUES(?, ?, ?);")
	res, err := stmt.Exec(u.name, u.email, u.pwd)
	if (err != nil) {
		panic(err)
	}
	id64, err := res.LastInsertId()
	u.id = int(id64)
	path := "./audio/" + strconv.Itoa(u.id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	return err
}

func (u *User) save() {
	//todo: validation
	u.name = u.Name
	u.email = u.Email
	u.pwd = u.Pwd
}

func (u *User) Auth() (bool, error) {
	u.save()
	u.checkDBSet()
	res, err := u.db.Query("SELECT id, name, email FROM users WHERE email='" + u.email + "' AND pwd='" + u.pwd + "'")
	if (err != nil) {
		return false, err
	} else {
		res.Scan(&u.id, &u.Name, &u.Email)
		u.save()
	}
	return res.Next(), err
}

func (u *User) Load(id int, db *sql.DB) (bool, error) {
	res, err := u.db.Query("SELECT name, email, pwd FROM user WHERE id=" + strconv.Itoa(id))
	if (err != nil) {
		return false, err
	}
	if (!res.Next()) {
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