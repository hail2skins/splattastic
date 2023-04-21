package testhelpers

import "github.com/hail2skins/splattastic/models"

func CreateTestUser() *models.User {
	user, _ := models.UserCreate(
		"test@example.com",
		"testpassword",
		"Test",
		"User",
		"testuser",
		models.UserType("Athlete"),
	)
	return user
}
