package entity

type User struct {
	BaseEntity
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}
