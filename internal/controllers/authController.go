package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/XSidik/go-fiber-template/internal/config"
	"github.com/XSidik/go-fiber-template/internal/database"
	"github.com/XSidik/go-fiber-template/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.UserModel{
		UserName: data["user_name"],
		Password: string(password),
	}

	database.DB.Create(&user)

	return c.JSON(models.APIResponse{
		Status:  true,
		Message: "User created successfully",
		Data:    user,
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
	}

	var user models.UserModel
	database.DB.Where("user_name = ?", data["user_name"]).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  false,
			Message: "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
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
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
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
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
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

	return c.JSON(models.APIResponse{
		Status:  true,
		Message: "Login successful",
		Data:    user,
		Meta: map[string]string{
			"access_token":  accessTokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

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
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
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
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
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

	return c.JSON(models.APIResponse{
		Status:  true,
		Message: "Token refreshed",
		Meta: map[string]string{
			"access_token":  accessTokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserModel)

	redisDb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetConfig().RedisHost, config.GetConfig().RedisPort),
		DB:   config.GetConfig().RedisDB,
	})

	redisDb.Del(ctx, "token_"+strconv.Itoa(int(user.ID)))
	redisDb.Del(ctx, "refresh_token"+strconv.Itoa(int(user.ID)))

	return c.JSON(models.APIResponse{
		Status:  true,
		Message: "Logout successful",
	})
}
