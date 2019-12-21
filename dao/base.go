package dao

type DaoBase interface {
	InsertSelf() (res bool, err error)
}

func Insert(d DaoBase) (res bool, err error) {
	res, err = d.InsertSelf()
	return
}
