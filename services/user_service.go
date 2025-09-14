package services

import (
	"context"
	"errors"
	"math"
	"time"

	"go-boilerplate/logger"
	"go-boilerplate/models"
	"go-boilerplate/models/request"
	"go-boilerplate/models/response"
	repoInterfaces "go-boilerplate/repository/interfaces"
	serviceInterfaces "go-boilerplate/services/interfaces"
	"go-boilerplate/utilities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userService struct {
	userRepo     repoInterfaces.UserRepository
	authService  serviceInterfaces.AuthService
	redisService serviceInterfaces.RedisService
}

func NewUserService(userRepo repoInterfaces.UserRepository, authService serviceInterfaces.AuthService, redisService serviceInterfaces.RedisService) serviceInterfaces.UserService {
	return &userService{
		userRepo:     userRepo,
		authService:  authService,
		redisService: redisService,
	}
}

func (s *userService) CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error) {
	ctx := context.Background()
	logger.Info(ctx, "UserService.CreateUser start", map[string]any{"email": req.Email})
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		logger.Warn(ctx, "CreateUser: email already exists", map[string]any{"email": req.Email})
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "CreateUser: password hash failed", map[string]any{"email": req.Email, "error": err.Error()})
		return nil, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error(ctx, "CreateUser: user repo create failed", map[string]any{"email": req.Email, "error": err.Error()})
		return nil, err
	}

	userResponse := utilities.ToUserResponse(user)

	// Cache the user
	if err := s.redisService.SetJSON(ctx, utilities.UserCacheKey(user.ID), userResponse, 30*time.Minute); err != nil {
		logger.Warn(ctx, "CreateUser: cache set failed", map[string]any{"user_id": user.ID, "error": err.Error()})
	}

	logger.Info(ctx, "UserService.CreateUser success", map[string]any{"user_id": user.ID, "email": req.Email})
	return userResponse, nil
}

func (s *userService) GetUserByID(id uint) (*response.UserResponse, error) {
	ctx := context.Background()
	logger.Debug(ctx, "UserService.GetUserByID start", map[string]any{"user_id": id})
	// Try to get from cache first
	var cachedUser response.UserResponse
	if err := s.redisService.GetJSON(ctx, utilities.UserCacheKey(id), &cachedUser); err == nil {
		logger.Info(ctx, "GetUserByID: cache hit", map[string]any{"user_id": id})
		return &cachedUser, nil
	} else {
		logger.Debug(ctx, "GetUserByID: cache miss", map[string]any{"user_id": id, "error": err.Error()})
	}

	// If not in cache, get from database
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(ctx, "GetUserByID: not found", map[string]any{"user_id": id})
			return nil, errors.New("user not found")
		}
		logger.Error(ctx, "GetUserByID: repo error", map[string]any{"user_id": id, "error": err.Error()})
		return nil, err
	}

	userResponse := utilities.ToUserResponse(user)

	// Cache the user for future requests
	if err := s.redisService.SetJSON(ctx, utilities.UserCacheKey(user.ID), userResponse, 30*time.Minute); err != nil {
		logger.Warn(ctx, "GetUserByID: cache set failed", map[string]any{"user_id": user.ID, "error": err.Error()})
	}

	logger.Info(ctx, "UserService.GetUserByID success", map[string]any{"user_id": id})
	return userResponse, nil
}

func (s *userService) GetUsers(page, perPage int) (*response.PaginationResponse, error) {
	ctx := context.Background()
	logger.Debug(ctx, "UserService.GetUsers start", map[string]any{"page": page, "per_page": perPage})
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	users, total, err := s.userRepo.GetAll(offset, perPage)
	if err != nil {
		logger.Error(ctx, "GetUsers: repo error", map[string]any{"page": page, "per_page": perPage, "error": err.Error()})
		return nil, err
	}

	userResponses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = utilities.ToUserResponse(user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	logger.Info(ctx, "UserService.GetUsers success", map[string]any{"count": len(userResponses), "total": total, "page": page, "per_page": perPage, "total_pages": totalPages})
	return &response.PaginationResponse{
		Data:       userResponses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) UpdateUser(id uint, req *request.UpdateUserRequest) (*response.UserResponse, error) {
	ctx := context.Background()
	logger.Info(ctx, "UserService.UpdateUser start", map[string]any{"user_id": id})
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(ctx, "UpdateUser: not found", map[string]any{"user_id": id})
			return nil, errors.New("user not found")
		}
		logger.Error(ctx, "UpdateUser: repo get failed", map[string]any{"user_id": id, "error": err.Error()})
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		logger.Error(ctx, "UpdateUser: repo update failed", map[string]any{"user_id": id, "error": err.Error()})
		return nil, err
	}

	userResponse := utilities.ToUserResponse(user)

	// Update cache
	if err := s.redisService.SetJSON(ctx, utilities.UserCacheKey(user.ID), userResponse, 30*time.Minute); err != nil {
		logger.Warn(ctx, "UpdateUser: cache set failed", map[string]any{"user_id": user.ID, "error": err.Error()})
	}

	logger.Info(ctx, "UserService.UpdateUser success", map[string]any{"user_id": id})
	return userResponse, nil
}

func (s *userService) DeleteUser(id uint) error {
	ctx := context.Background()
	logger.Info(ctx, "UserService.DeleteUser start", map[string]any{"user_id": id})
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(ctx, "DeleteUser: not found", map[string]any{"user_id": id})
			return errors.New("user not found")
		}
		logger.Error(ctx, "DeleteUser: repo get failed", map[string]any{"user_id": id, "error": err.Error()})
		return err
	}

	// Delete from database
	if err := s.userRepo.Delete(id); err != nil {
		logger.Error(ctx, "DeleteUser: repo delete failed", map[string]any{"user_id": id, "error": err.Error()})
		return err
	}

	// Remove from cache
	if err := s.redisService.Delete(ctx, utilities.UserCacheKey(id)); err != nil {
		logger.Warn(ctx, "DeleteUser: cache delete failed", map[string]any{"user_id": id, "error": err.Error()})
	}

	logger.Info(ctx, "UserService.DeleteUser success", map[string]any{"user_id": id})
	return nil
}

func (s *userService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	ctx := context.Background()
	logger.Info(ctx, "UserService.Login start", map[string]any{"email": req.Email})
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		logger.Warn(ctx, "Login: user not found", map[string]any{"email": req.Email})
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn(ctx, "Login: password mismatch", map[string]any{"email": req.Email})
		return nil, errors.New("invalid email or password")
	}

	token, err := s.authService.GenerateToken(user.ID)
	if err != nil {
		logger.Error(ctx, "Login: token generation failed", map[string]any{"user_id": user.ID, "error": err.Error()})
		return nil, err
	}

	// Cache user session
	if err := s.redisService.CacheUserSession(ctx, user.ID, token, 24*time.Hour); err != nil {
		logger.Warn(ctx, "Login: cache session failed", map[string]any{"user_id": user.ID, "error": err.Error()})
	}

	logger.Info(ctx, "UserService.Login success", map[string]any{"user_id": user.ID})
	return &response.LoginResponse{
		Token: token,
		User:  *utilities.ToUserResponse(user),
	}, nil
}

func (s *userService) Logout(userID uint) error {
	ctx := context.Background()
	logger.Info(ctx, "UserService.Logout start", map[string]any{"user_id": userID})
	err := s.redisService.DeleteUserSession(ctx, userID)
	if err != nil {
		logger.Error(ctx, "Logout: delete session failed", map[string]any{"user_id": userID, "error": err.Error()})
		return err
	}
	logger.Info(ctx, "UserService.Logout success", map[string]any{"user_id": userID})
	return nil
}
