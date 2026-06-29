package application

type TokenService interface {
	GenerateToken(userID string, role string) (string, error)
	ValidateToken(token string) (string, error)
}
