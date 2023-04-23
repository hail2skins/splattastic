package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	// ... setup code from the previous test ...
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("../templates/**/**")
	r.POST("/logout", Logout)
	r.POST("/login", Login)

	// Create a test user
	userType := models.UserType{Name: "Test User Type"}
	require.NoError(t, db.Database.Create(&userType).Error)
	user, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		userType.Name,
	)
	require.NoError(t, err)
	require.NotNil(t, user)

	// Login the test user
	loginReq := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(
		url.Values{
			"email":    {user.Email},
			"password": {"testpassword"},
		}.Encode(),
	))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginResp := httptest.NewRecorder()
	r.ServeHTTP(loginResp, loginReq)
	require.Equal(t, http.StatusMovedPermanently, loginResp.Code)

	// Logout the test user
	logoutReq := httptest.NewRequest(http.MethodPost, "/logout", nil)
	logoutReq.Header.Set("Cookie", loginResp.Header().Get("Set-Cookie"))
	logoutResp := httptest.NewRecorder()
	r.ServeHTTP(logoutResp, logoutReq)
	require.Equal(t, http.StatusOK, logoutResp.Code)
	assert.Contains(t, logoutResp.Body.String(), "Successfully logged out")

	// Clean up the test setup
	require.NoError(t, db.Database.Unscoped().Delete(user).Error)
	require.NoError(t, db.Database.Unscoped().Delete(&userType).Error)
}
