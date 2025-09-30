package models

type Role int

const (
	USER Role = iota
	ADMIN
)

const (
	userRoleString  = "User"
	adminRoleString = "Admin"
	unknownRole     = "Unknown"
)

func (r Role) RoleString() string {
	switch r {
	case USER:
		return userRoleString
	case ADMIN:
		return adminRoleString
	default:
		return unknownRole
	}
}

func StringToRole(roleStr string) Role {
	switch roleStr {
	case adminRoleString:
		return ADMIN
	case userRoleString:
		return USER
	default:
		return USER
	}
}

type User struct {
	Username string
	Password string
	Email    string
	ID       int32
	Role     Role
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponseDTO struct {
	Token string `json:"token"`
}
