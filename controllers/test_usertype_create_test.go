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
	m "github.com/hail2skins/splattastic/middlewares"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestUserTypeCreate(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(m.AuthenticateUser())

	r.LoadHTMLGlob("../templates/**/**")

	r.POST("/login", Login)

	admin := r.Group("/admin")
	{
		usertypes := admin.Group("/usertypes")
		{
			usertypes.GET("/new", UserTypeNew)
			usertypes.POST("/", UserTypeCreate)
		}
	}

	// Create a test user type
	testUserType := models.UserType{Name: "TestType"}
	db.Database.Create(&testUserType)

	// Create a test user with hashed password and set its user type to admin
	testUser, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		"TestType",
	)
	assert.NoError(t, err)
	assert.NotNil(t, testUser)

	// Create a request with POST method and form data
	form := url.Values{}
	form.Add("name", "Test UserType")
	req, err := http.NewRequest(http.MethodPost, "/admin/usertypes", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the UserTypeCreate function directly with the request and response recorder
	UserTypeCreate(ginContext(req, w))

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the user type was created in the database
	var userType models.UserType
	db.Database.Where("name = ?", "Test UserType").First(&userType)
	assert.Equal(t, "Test UserType", userType.Name)

	// Cleanup
	db.Database.Unscoped().Delete(&userType)
	db.Database.Unscoped().Delete(testUser)
	db.Database.Unscoped().Delete(&testUserType)
}

func ginContext(req *http.Request, w *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}
