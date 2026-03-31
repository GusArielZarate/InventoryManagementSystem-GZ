package users

type User struct {
	Name string `json:"name"`
	Rol  string `json:"rol"`
}

func NewUser(nombre, rol string) User {
	return User{Name: nombre, Rol: rol}
}

func (u *User) GetNombre() string       { return u.Name }
func (u *User) GetRol() string          { return u.Rol }
func (u *User) SetNombre(nombre string) { u.Name = nombre }
func (u *User) SetRol(rol string)       { u.Rol = rol }
