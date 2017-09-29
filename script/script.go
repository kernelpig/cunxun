package script

func InitScript() {
	if err := CreateDatabase(); err != nil {
		panic(err.Error())
	}
	user, err := CreateSuperAdmin()
	if err != nil {
		panic(err.Error())
	}
	_, err = CreateColumns(user)
	if err != nil {
		panic(err.Error())
	}
}
