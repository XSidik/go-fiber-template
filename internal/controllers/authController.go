package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/XSidik/go-fiber-template/internal/config"
	"github.com/XSidik/go-fiber-template/internal/database"
	"github.com/XSidik/go-fiber-template/internal/models"
	"github.com/XSidik/go-fiber-template/internal/models/dto"
	"github.com/XSidik/go-fiber-template/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

// Register handles the user registration process.
// @Summary Register a new user
// @Description Register a new user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.Register true "User Data"
// @Success 200 {object} dto.APIResponse "User created successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request"
// @Router /api/v1/auth/register [post]
func Register(c *fiber.Ctx) error {
	validate := validator.New()
	var data dto.Register

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	if err := validate.Struct(data); err != nil {
		formatted := utils.FormatValidationError(err)
		response := dto.APIResponse{
			Status:  false,
			Message: "Validation failed",
			Errors:  formatted,
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validate the data
	data.Password = strings.TrimSpace(data.Password)
	if err := validate.Struct(data); err != nil {
		// Validation failed
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, len(validationErrors))

		for i, v := range validationErrors {
			errorMessages[i] = v.Field() + " failed on the '" + v.Tag() + "' tag"
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.APIResponse{
			Status:  false,
			Message: "Validation failed",
			Errors:  errorMessages,
		})
	}

	// check the username to make it unique
	var userExist models.UserModel
	database.DB.Where("user_name = ?", data.UserName).First(&userExist)
	if userExist.ID != 0 {
		return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
			Status:  false,
			Message: "User already exist",
			Errors:  "Username " + data.UserName + "already exist, use another name",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	user := models.UserModel{
		UserName: data.UserName,
		Password: string(password),
	}

	database.DB.Create(&user)

	return c.JSON(dto.APIResponse{
		Status:  true,
		Message: "User created successfully",
		Data:    user,
	})
}

// Login handles the user login process.
// @Summary Login a user
// @Description Login a user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.Login true "User Data"
// @Success 200 {object} dto.APIResponse "Login successful"
// @Failure 400 {object} dto.APIResponse "Invalid request"
// @Failure 404 {object} dto.APIResponse "User not found"
// @Failure 500 {object} dto.APIResponse "Error generating token"
// @Router /api/v1/auth/login [post]
func Login(c *fiber.Ctx) error {
	var data dto.Login

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	// // Validate the data
	// if err := validate.Struct(data); err != nil {
	// 	return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
	// 		Status:  false,
	// 		Message: "Validation failed",
	// 		Errors:  err.Error(),
	// 	})
	// }

	var user models.UserModel
	database.DB.Where("user_name = ?", data.UserName).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  false,
			Message: "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  false,
			Message: "Invalid credentials",
		})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	config := config.GetConfig()
	accessTokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Error generating access token",
		})
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Error generating refresh token",
		})
	}

	redisDb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		DB:   config.RedisDB,
	})

	redisDb.Set(ctx, "token_"+strconv.Itoa(int(user.ID)), accessTokenString, time.Hour*24*7)
	redisDb.Set(ctx, "refresh_token"+strconv.Itoa(int(user.ID)), refreshTokenString, time.Hour*24*7)

	return c.JSON(dto.APIResponse{
		Status:  true,
		Message: "Login successful",
		Data:    user,
		Meta: map[string]string{
			"access_token":  accessTokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

// RefreshToken handles the user refresh token process.
// @Summary Refresh user token
// @Description Refresh the access and refresh tokens for an authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse "Token refreshed successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request"
// @Failure 500 {object} dto.APIResponse "Error generating token"
// @Router /api/v1/auth/refresh-token [get]
func RefreshToken(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserModel)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	config := config.GetConfig()
	accessTokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Error generating access token",
		})
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  false,
			Message: "Error generating refresh token",
		})
	}

	redisDb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		DB:   config.RedisDB,
	})

	redisDb.Set(ctx, "token_"+strconv.Itoa(int(user.ID)), accessTokenString, time.Hour*24*7)
	redisDb.Set(ctx, "refresh_token"+strconv.Itoa(int(user.ID)), refreshTokenString, time.Hour*24*7)

	return c.JSON(dto.APIResponse{
		Status:  true,
		Message: "Token refreshed",
		Meta: map[string]string{
			"access_token":  accessTokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

// Logout handles the user logout process.
// @Summary Logout a user
// @Description Logout a user and delete their tokens from Redis
// @Tags auth
// @Produce json
// @Success 200 {object} dto.APIResponse "Logout successful"
// @Failure 500 {object} dto.APIResponse "Error during logout"
// @Router /api/v1/auth/logout [get]
func Logout(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserModel)

	redisDb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetConfig().RedisHost, config.GetConfig().RedisPort),
		DB:   config.GetConfig().RedisDB,
	})

	redisDb.Del(ctx, "token_"+strconv.Itoa(int(user.ID)))
	redisDb.Del(ctx, "refresh_token"+strconv.Itoa(int(user.ID)))

	return c.JSON(dto.APIResponse{
		Status:  true,
		Message: "Logout successful",
	})
}
