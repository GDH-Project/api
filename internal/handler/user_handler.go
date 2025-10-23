package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/middleware"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type userResponse struct {
	Status int
	Body   struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
}

type tokenResponse struct {
	Status int
	Body   domain.Token
}

// RegisterAuthHandler 인증 및 유저 관련 Handler
func RegisterAuthHandler(api huma.API, log *zap.Logger, authUseCase domain.AuthUseCase, userUseCase domain.UserUseCase, m middleware.Middleware) {
	v1 := huma.NewGroup(api, "/api/v1")

	// 회원 가입
	huma.Register(v1, huma.Operation{
		OperationID:   "v1AuthSignUp",
		Method:        http.MethodPost,
		Path:          "/auth/sign-up",
		Summary:       "회원 가입",
		Description:   "회원 가입 API 입니다.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *struct {
		Body struct {
			Email    string `json:"email,required" format:"email" doc:"사용자 이메일 입니다."`
			Name     string `json:"name,required" minLength:"3" doc:"사용자 닉네임 입니다" example:"tester"`
			Role     string `json:"role" enum:"user,device" default:"user" doc:"사용자의 권한 입니다."`
			Password string `json:"password,required" minLength:"8" format:"password"`
		}
	}) (*struct{}, error) {
		err := userUseCase.CreateUser(ctx, &domain.User{
			Name:     i.Body.Name,
			Email:    i.Body.Email,
			Password: i.Body.Password,
			Role:     domain.ParseStringRoleToUserRole(i.Body.Role),
		})

		if err != nil {
			return nil, huma.Error400BadRequest("사용자 생성에 실패했습니다.", err)
		}

		return nil, nil
	})

	// 로그인
	huma.Register(v1, huma.Operation{
		OperationID:   "v1AuthSignIn",
		Method:        http.MethodPost,
		Path:          "/auth/sign-in",
		Summary:       "로그인",
		Description:   "로그인 API 입니다.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		Type string `query:"type" enum:"password" default:"password" doc:"로그인 구분 입니다."`
		Body struct {
			Email    string `json:"email,required" format:"email"`
			Password string `json:"password,required" minLength:"8" format:"password"`
		}
	}) (*tokenResponse, error) {
		var resp tokenResponse
		// ?type=password 인 경우
		if strings.EqualFold(i.Type, "password") {
			token, err := authUseCase.Login(ctx, i.Body.Email, i.Body.Password)
			if err != nil {
				return nil, huma.Error400BadRequest("id 혹은 패스워드를 확인해주세요")
			}
			resp.Body = *token

			return &resp, nil
		}

		// type 이 잘못된 경우
		return nil, huma.Error400BadRequest("잘못된 접근 입니다.")
	})

	// 토큰 재발급
	huma.Register(v1, huma.Operation{
		OperationID:   "v1AuthRefresh",
		Method:        http.MethodPost,
		Path:          "/auth/refresh",
		Summary:       "토큰 재발급",
		Description:   "토큰 재발급 API 입니다.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *struct {
		Body struct {
			RefreshToken string `json:"refresh_token"`
		}
	}) (*tokenResponse, error) {
		var resp tokenResponse
		token, err := authUseCase.RefreshToken(ctx, i.Body.RefreshToken)
		if err != nil {
			return nil, huma.Error401Unauthorized("토큰이 유효하지 않습니다.")
		}

		resp.Body = *token
		return &resp, nil
	})

	// 사용자 정보 조회
	huma.Register(v1, m.WithAuth(huma.Operation{
		OperationID:   "v1GetUserInfo",
		Method:        http.MethodGet,
		Path:          "/user",
		Summary:       "사용자 정보 조회",
		Description:   "사용자 정보 조회 API 입니다.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusOK,
	}), func(ctx context.Context, i *struct{}) (*userResponse, error) {
		var resp userResponse
		userId, _ := ctx.Value("user_id").(string)

		u, err := userUseCase.GetUserInfoByUserID(ctx, userId)
		if err != nil {
			return nil, huma.Error404NotFound("존재하지 않는 사용자 입니다.")
		}

		resp.Body.Name = u.Name
		resp.Body.Email = u.Email
		resp.Body.Role = string(u.Role)

		return &resp, nil
	})

	// 사용자 정보 업데이트
	huma.Register(v1, m.WithAuth(huma.Operation{
		OperationID:   "v1UpdateUserInfo",
		Method:        http.MethodPut,
		Path:          "/user",
		Summary:       "사용자 정보 업데이트",
		Description:   "사용자 정보 업데이트 API 입니다.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusOK,
	}), func(ctx context.Context, i *struct {
		Body struct {
			Name     string `json:"name,omitempty"`
			Password string `json:"password,omitempty" minLength:"8" format:"password"`
		}
	}) (*userResponse, error) {
		var resp userResponse
		userID, _ := ctx.Value("user_id").(string)

		u, err := userUseCase.UpdateUser(ctx, &domain.User{
			ID:       userID,
			Name:     i.Body.Name,
			Password: i.Body.Password,
		})
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		resp.Body.Name = u.Name
		resp.Body.Email = u.Email
		resp.Body.Role = string(u.Role)

		return &resp, nil
	})

	// 회원 탈퇴
	huma.Register(v1, m.WithAuth(huma.Operation{
		OperationID:   "v1DeleteUser",
		Method:        http.MethodPost,
		Path:          "/user/delete",
		Summary:       "회원 탈퇴",
		Description:   "회원 탈퇴 API 입니다.",
		Tags:          []string{"Auth"},
		DefaultStatus: http.StatusOK,
	}), func(ctx context.Context, i *struct {
		Body struct {
			Password string `json:"password" minLength:"8" format:"password"`
		}
	}) (*struct{}, error) {
		userID, _ := ctx.Value("user_id").(string)

		err := userUseCase.DeleteUser(ctx, userID, i.Body.Password)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	})

}
