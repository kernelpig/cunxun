package script

func InitScript() error {
	if err := CreateDatabase(); err != nil {
		return err
	}
	user, err := CreateSuperAdmin()
	if err != nil {
		return err
	}
	_, err = CreateColumns(user)
	if err != nil {
		return err
	}
	return nil
}
