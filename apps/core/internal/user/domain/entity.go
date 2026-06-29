package domain

import "golang.org/x/crypto/bcrypt"

type Role string

const (
	AdminRole         Role = "admin"
	MemberRole        Role = "member"
	OrgAdminRole      Role = "org_admin"
	TeamAdminRole     Role = "team_admin"
	PlatformAdminRole Role = "platform_admin"
)

type Password string

func (p Password) IsValid() bool {
	return len(p) >= 8
}

type Email string

func (e Email) IsValid() bool {
	return len(e) > 0
}

type User struct {
	ID       string
	Username string
	Role     Role
	Email    Email
	Password string
}

func NewUser(id, username, rawEmail, rawPassword string) (*User, error) {
	password := Password(rawPassword)

	if !password.IsValid() {
		return nil, ErrInvalidPassword
	}

	email := Email(rawEmail)

	if !email.IsValid() {
		return nil, ErrInvalidEmail
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     MemberRole,
	}, nil
}

func (u *User) VerifyPassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword))
	return err == nil
}
