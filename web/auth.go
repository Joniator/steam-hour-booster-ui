package web

type Auth struct {
	Username string
	Password string
}

func (auth *Auth) IsEnabled() bool {
	return auth.Username != "" && auth.Password != ""
}
