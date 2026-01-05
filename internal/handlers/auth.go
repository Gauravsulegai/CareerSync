package handlers

import (
	"net/http"
	"strings"
	"time"
	"fmt"

	"github.com/Gauravsulegai/careersync/internal/database"
	"github.com/Gauravsulegai/careersync/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// SignupInput defines what JSON we expect from the frontend
type SignupInput struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"` // Personal Email
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"` // "student" or "employee"
	
	// Employee Specific
	CompanyName string `json:"company_name"` // Text input from search bar
	WorkEmail   string `json:"work_email"`   // Must match company domain
	Position    string `json:"position"`
}

func Signup(c *gin.Context) {
	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 2. Prepare User Object
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	// 3. EMPLOYEE LOGIC
	if input.Role == "employee" {
		if input.CompanyName == "" || input.WorkEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Employees must provide Company Name and Work Email"})
			return
		}

		// A. Parse Domain from Work Email (e.g., "gaurav@google.com" -> "google.com")
		parts := strings.Split(input.WorkEmail, "@")
		if len(parts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Work Email format"})
			return
		}
		userDomain := parts[1]

		// B. Check if Company Exists
		var company models.Company
		result := database.DB.Where("name = ?", input.CompanyName).First(&company)

		if result.Error == nil {
			// COMPANY EXISTS -> Domain Check
			if company.Domain != userDomain {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Domain Mismatch! You must use an email ending in @" + company.Domain,
				})
				return
			}
			// Link User
			user.CompanyID = &company.ID
			user.Company = &company

		} else {
			// COMPANY DOES NOT EXIST -> Create It (First Settler)
			newCompany := models.Company{
				Name:       input.CompanyName,
				Domain:     userDomain, // Lock this domain for future users!
				FormConfig: []byte(`{"require_resume": true, "require_job_id": true}`), // Default JSON
			}
			database.DB.Create(&newCompany)
			
			// Link User
			user.CompanyID = &newCompany.ID
			user.Company = &newCompany
		}

		user.WorkEmail = input.WorkEmail
		user.Position = input.Position
		user.IsVerified = false // Needs verification
	}

	// 4. Save User to DB
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists or Database error"})
		return
	}

	// 5. SIMULATE VERIFICATION EMAIL (Print to Terminal)
	if input.Role == "employee" {
		fmt.Println("\n===========================================")
		fmt.Println("📧  [SIMULATION] Verification Email Sent!")
		fmt.Println("To verify " + input.Name + ", click this link:")
		fmt.Printf("http://localhost:8080/verify/%d \n", user.ID)
		fmt.Println("===========================================\n")
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!", "user": user})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// We use Preload("Company") so the login response includes company details!
	if err := database.DB.Preload("Company").Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Make sure this key matches require_auth.go!
	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user})
}