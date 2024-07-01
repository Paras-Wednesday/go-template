package resolver_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/internal/config"
	"go-template/internal/jwt"
	"go-template/internal/service"
	"go-template/models"
	"go-template/pkg/utl/resultwrapper"
	"go-template/pkg/utl/secure"
	"go-template/resolver"
	"go-template/testutls"
)

const (
	UserRoleName               = "UserRole"
	SuperAdminRoleName         = "SuperAdminRole"
	ErrorFromRedisCache        = "RedisCache Error"
	ErrorFromGetRole           = "RedisCache GetRole Error"
	ErrorUnauthorizedUser      = "Unauthorized User"
	ErrorFromCreateRole        = "CreateRole Error"
	ErrorPasswordValidation    = "Fail on PasswordValidation"
	ErrorActiveStatus          = "Fail on ActiveStatus"
	ErrorInsecurePassword      = "Insecure password"
	ErrorInvalidToken          = "Fail on FindByToken"
	ErrorUpdateUser            = "User Update Error"
	ErrorDeleteUser            = "User Delete Error"
	ErrorFromConfig            = "Config Error"
	ErrorFromBool              = "Boolean Error"
	ErrorMsgFromConfig         = "error in loading config"
	ErrorMsginvalidToken       = "error from FindByToken"
	ErrorMsgFindingUser        = "error in finding the user"
	ErrorMsgFromJwt            = "error in creating auth service"
	ErrorMsgfromUpdateUser     = "error while updating user"
	ErrorMsgPasswordValidation = "username or password does not exist "
	TestPasswordHash           = "$2a$10$dS5vK8hHmG5"
	OldPasswordHash            = "$2a$10$dS5vK8hHmG5gzwV8f7TK5.WHviMBqmYQLYp30a3XvqhCW9Wvl2tOS"
	SuccessCase                = "Success"
	ErrorFindingUser           = "Fail on finding user"
	ErrorFromCreateUser        = "Fail on Create User"
	ErrorFromThrottleCheck     = "Throttle error"
	ErrorFromJwt               = "Jwt Error"
	ErrorFromGenerateToken     = "Token Error"
	OldPassword                = "adminuser"
	NewPassword                = "adminuser!A9@"
	TestPassword               = "pass123"
	TestUsername               = "wednesday"
	TestToken                  = "refreshToken"
	ReqToken                   = "refresh_token"
)

type loginArgs struct {
	UserName string
	Password string
}

type loginType struct {
	name     string
	req      loginArgs
	wantResp *fm.LoginResponse
	wantErr  bool
	err      error
	init     func() *gomonkey.Patches
}

func errorFindingUserCase() loginType {
	return loginType{
		name: ErrorFindingUser,
		req: loginArgs{
			UserName: TestUsername,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFindingUser),
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf(ErrorMsgFindingUser)
				})
		},
	}
}

func errorPasswordValidationCase() loginType {
	return loginType{
		name: ErrorPasswordValidation,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: TestPassword,
		},
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgPasswordValidation),
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					return &models.User{
						ID:       testutls.MockID,
						Username: null.StringFrom(testutls.MockEmail),
						Password: null.StringFrom(""),
						Active:   null.BoolFrom(true),
					}, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return false
				})
		},
	}
}

func errorActiveStatusCase() loginType {
	return loginType{
		name: ErrorActiveStatus,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					return &models.User{
						ID:       testutls.MockID,
						Username: null.StringFrom(testutls.MockEmail),
						Password: null.StringFrom(""),
						Active:   null.BoolFrom(false),
					}, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return true
				})
		},
	}
}

// func errorFromConfigCase() loginType {
// 	return loginType{
// 		name: ErrorFromConfig,
// 		req: loginArgs{
// 			UserName: testutls.MockEmail,
// 			Password: OldPassword,
// 		},
// 		wantErr: true,
// 		err:     fmt.Errorf(ErrorMsgFromConfig),
// 		init: func() *gomonkey.Patches {
// 			return gomonkey.ApplyFunc(daos.FindUserByUserName,
// 				func(username string, ctx context.Context) (*models.User, error) {
// 					user := testutls.MockUser()
// 					user.Password = null.StringFrom(OldPasswordHash)
// 					user.Active = null.BoolFrom(false)
// 					return user, nil
// 				}).
// 				ApplyFunc(config.Load, func() (*config.Configuration, error) {
// 					return nil, fmt.Errorf("error in loading config")
// 				})
// 		},
// 	}
// }

func errorWhileGeneratingToken() loginType {
	return loginType{
		name: ErrorFromJwt,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     resultwrapper.ErrUnauthorized,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					return &models.User{
						ID:       testutls.MockID,
						Username: null.StringFrom(testutls.MockEmail),
						Password: null.StringFrom(""),
						Active:   null.BoolFrom(true),
					}, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return true
				}).ApplyMethod(reflect.TypeOf(jwt.Service{}),
				"GenerateToken",
				func(svc jwt.Service, u *models.User) (string, error) {
					return "", fmt.Errorf(ErrorMsgFromJwt)
				})
		},
	}
}

func errorUpdateUserCase() loginType {
	err := fmt.Errorf(ErrorMsgfromUpdateUser)
	return loginType{
		name: ErrorUpdateUser,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantErr: true,
		err:     err,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					return &models.User{
						ID:       testutls.MockID,
						Username: null.StringFrom(testutls.MockEmail),
						Password: null.StringFrom(""),
						Active:   null.BoolFrom(true),
					}, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return true
				}).ApplyMethod(reflect.TypeOf(jwt.Service{}),
				"GenerateToken",
				func(svc jwt.Service, u *models.User) (string, error) {
					return "token", nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"Token",
				func(sec secure.Service, str string) string {
					return "refresh token"
				}).ApplyFunc(daos.UpdateUser,
				func(user models.User, ctx context.Context) (models.User, error) {
					return models.User{}, fmt.Errorf(ErrorMsgfromUpdateUser)
				},
			)
		},
	}
}

func loginSuccessCase() loginType {
	jwtToken := "jwttokenstring"
	return loginType{
		name: SuccessCase,
		req: loginArgs{
			UserName: testutls.MockEmail,
			Password: OldPassword,
		},
		wantResp: &fm.LoginResponse{
			Token:        jwtToken,
			RefreshToken: TestToken,
		},
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			sec := secure.Service{}
			return gomonkey.ApplyFunc(daos.FindUserByUserName,
				func(username string, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(true)
					return user, nil
				}).
				ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
					return tg, nil
				}).
				ApplyFunc(service.Secure, func(cfg *config.Configuration) secure.Service {
					return sec
				}).
				ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
					func(jwt.Service, *models.User) (string, error) {
						return jwtToken, nil
					}).
				ApplyMethod(reflect.TypeOf(sec), "Token",
					func(secure.Service, string) string {
						return TestToken
					}).
				ApplyFunc(daos.UpdateUser,
					func(u models.User, ctx context.Context) (models.User, error) {
						return *testutls.MockUser(), nil
					}).ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return nil, nil
			})
		},
	}
}

// func errorWhileCreatingJWTService() loginType {
// 	err := fmt.Errorf("error in creating auth service")
// 	return loginType{
// 		name: "Error while creating a JWT Service",
// 		req: loginArgs{
// 			UserName: testutls.MockEmail,
// 			Password: OldPassword,
// 		},
// 		wantErr: true,
// 		err:     err,
// 		init: func() *gomonkey.Patches {
// 			sec := secure.Service{}
// 			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
// 				return nil, nil
// 			}).ApplyFunc(service.Secure, func(cfg *config.Configuration) secure.Service {
// 				return sec
// 			}).ApplyFunc(daos.FindUserByUserName,
// 				func(username string, ctx context.Context) (*models.User, error) {
// 					user := testutls.MockUser()
// 					user.Password = null.StringFrom(OldPasswordHash)
// 					user.Active = null.BoolFrom(false)
// 					return user, nil
// 				}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
// 				// mock service.JWT
// 				return jwt.Service{}, err
// 			})
// 		},
// 	}
// }

func loadLoginTestCases() []loginType {
	return []loginType{
		errorWhileGeneratingToken(),
		errorFindingUserCase(),
		// errorFromConfigCase(),
		errorPasswordValidationCase(),
		errorActiveStatusCase(),
		// errorWhileCreatingJWTService(),
		errorUpdateUserCase(),
		loginSuccessCase(),
	}
}

func TestLogin(
	t *testing.T,
) {
	cases := loadLoginTestCases()
	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				patch := tt.init()
				c := context.Background()
				// Call the login mutation with the given arguments and check the response and error against the expected values
				response, err := resolver1.Mutation().Login(c, tt.req.UserName, tt.req.Password)
				if tt.wantResp != nil &&
					response != nil {
					// Assert that the expected response matches the actual response
					assert.Equal(t, tt.wantResp, response)
				} else {
					// Assert that the expected error value matches the actual error value
					assert.Equal(t, true, strings.Contains(err.Error(), tt.err.Error()),
						"wantErr: %t, with message: %s got: %v", tt.wantErr, tt.err.Error(), err)
					assert.Equal(t, tt.wantErr, err != nil)
				}
				if patch != nil {
					patch.Reset()
				}
			},
		)
	}
}

type changeReq struct {
	OldPassword string
	NewPassword string
}

type changePasswordType struct {
	name     string
	req      changeReq
	wantResp *fm.ChangePasswordResponse
	wantErr  bool
	init     func() *gomonkey.Patches
}

func changePasswordErrorFindingUserCase() changePasswordType {
	return changePasswordType{
		name: ErrorFindingUser,
		req: changeReq{
			OldPassword: TestPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf(ErrorMsgFindingUser)
				})
		},
	}
}

func changePasswordErrorPasswordValidationcase() changePasswordType {
	return changePasswordType{
		name: ErrorPasswordValidation,
		req: changeReq{
			OldPassword: TestPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf(ErrorMsgFindingUser)
				}).ApplyMethod(reflect.TypeOf(secure.Service{}), "HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return false
				})
		},
	}
}

func changePasswordErrorInsecurePasswordCase() changePasswordType {
	return changePasswordType{
		name: ErrorInsecurePassword,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: testutls.MockEmail,
		},
		wantErr: true,
		init: func() *gomonkey.Patches {
			sec := secure.Service{}
			// mock FindUserByUserName with the proper password, and active state
			return gomonkey.ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, nil
				}).ApplyMethod(reflect.TypeOf(sec), "HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return true
				}).ApplyMethod(reflect.TypeOf(sec), "Password",
				func(sec secure.Service, pass string, inputs ...string) bool {
					return false
				})
		},
	}
}

func changePasswordErrorUpdateUserCase() changePasswordType {
	return changePasswordType{
		name: ErrorUpdateUser,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: NewPassword,
		},
		wantErr: true,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, fmt.Errorf(ErrorInsecurePassword)
				}).ApplyMethod(reflect.TypeOf(secure.Service{}), "HashMatchesPassword", func(secure.Service, string, string) bool {
				return true
			}).ApplyMethod(reflect.TypeOf(secure.Service{}), "Password", func(secure.Service, string, ...string) bool {
				return true
			}).ApplyFunc(daos.UpdateUser,
				func(user models.User, ctx context.Context) (*models.User, error) {
					return nil, fmt.Errorf(ErrorUpdateUser)
				})
		},
	}
}

// func changePasswordErrorFromConfigCase() changePasswordType {
// 	return changePasswordType{
// 		name: ErrorFromConfig,
// 		req: changeReq{
// 			OldPassword: OldPassword,
// 			NewPassword: testutls.MockEmail,
// 		},
// 		wantErr: true,
// 		init: func() *gomonkey.Patches {
// 			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
// 				return nil, fmt.Errorf("error in loading config")
// 			}).ApplyFunc(daos.FindUserByID,
// 				func(userID int, ctx context.Context) (*models.User, error) {
// 					user := testutls.MockUser()
// 					user.Password = null.StringFrom(OldPasswordHash)
// 					user.Active = null.BoolFrom(false)
// 					return user, nil
// 				})
// 		},
// 	}
// }

func changePasswordSuccessCase() changePasswordType {
	return changePasswordType{
		name: SuccessCase,
		req: changeReq{
			OldPassword: OldPassword,
			NewPassword: NewPassword,
		},
		wantResp: &fm.ChangePasswordResponse{
			Ok: true,
		},
		wantErr: false,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByID,
				func(userID int, ctx context.Context) (*models.User, error) {
					user := testutls.MockUser()
					user.Password = null.StringFrom(OldPasswordHash)
					user.Active = null.BoolFrom(false)
					return user, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash, password string) bool {
					return true
				}).ApplyFunc(daos.UpdateUser,
				func(user models.User, ctx context.Context) (models.User, error) {
					return *testutls.MockUser(), nil
				})
		},
	}
}

func loadChangePasswordTestCases() []changePasswordType {
	return []changePasswordType{
		changePasswordErrorFindingUserCase(),
		changePasswordErrorPasswordValidationcase(),
		changePasswordErrorInsecurePasswordCase(),
		changePasswordErrorUpdateUserCase(),
		// changePasswordErrorFromConfigCase(),
		changePasswordSuccessCase(),
	}
}

func TestChangePassword(
	t *testing.T,
) {
	// Define a struct to represent the change password request
	cases := loadChangePasswordTestCases()
	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				// Handle the case where there is an error while loading the configuration
				patches := tt.init()
				// Set up the context with the mock user
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())

				// Call the ChangePassword mutation and check the response and error against the expected values
				response, err := resolver1.Mutation().ChangePassword(ctx, tt.req.OldPassword, tt.req.NewPassword)
				if tt.wantResp != nil {
					// Assert that the expected response matches the actual response
					assert.Equal(t, tt.wantResp, response)
				}
				// Assert that the expected error value matches the actual error value
				assert.Equal(t, tt.wantErr, err != nil)
				if patches != nil {
					patches.Reset()
				}
			},
		)
	}
}

type refereshTokenType struct {
	name     string
	req      string
	wantResp *fm.RefreshTokenResponse
	wantErr  bool
	err      error
	init     func() *gomonkey.Patches
}

func refreshTokenInvalidCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorInvalidToken,
		req:     TestToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsginvalidToken),
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, nil
			}).ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
				func(jwt.Service, *models.User) (string, error) {
					return "", fmt.Errorf(ErrorMsginvalidToken)
				}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
		},
	}
}

// func refreshTokenErrorFromConfigCase() refereshTokenType {
// 	return refereshTokenType{
// 		name:    ErrorFromConfig,
// 		req:     ReqToken,
// 		wantErr: true,
// 		err:     fmt.Errorf(ErrorMsgFromConfig),
// 		init: func() *gomonkey.Patches {
// 			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
// 				return nil, fmt.Errorf(ErrorFromConfig)
// 			}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
// 				return testutls.MockUser(), nil
// 			})
// 		},
// 	}
// }

func refereshTokenerrorWhileGeneratingToken() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromJwt,
		req:     ReqToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorMsgFromJwt),
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			}).ApplyMethod(jwt.Service{},
				"GenerateToken",
				func(tg jwt.Service, u *models.User) (string, error) {
					return "", fmt.Errorf(ErrorMsgFromJwt)
				})
		},
	}
}

func refereshTokenErrorFromGenerateTokenCase() refereshTokenType {
	return refereshTokenType{
		name:    ErrorFromGenerateToken,
		req:     ReqToken,
		wantErr: true,
		err:     fmt.Errorf(ErrorFromGenerateToken),
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return &config.Configuration{}, nil
			}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, nil
			}).ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
				func(jwt.Service, *models.User) (string, error) {
					return "", fmt.Errorf(ErrorFromGenerateToken)
				}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
		},
	}
}

func refreshTokenSuccessCase() refereshTokenType {
	return refereshTokenType{
		name: SuccessCase,
		req:  ReqToken,
		wantResp: &fm.RefreshTokenResponse{
			Token: testutls.MockToken,
		},
		wantErr: false,
		init: func() *gomonkey.Patches {
			tg := jwt.Service{}
			return gomonkey.ApplyFunc(config.Load, func() (*config.Configuration, error) {
				return &config.Configuration{}, nil
			}).ApplyFunc(service.JWT, func(cfg *config.Configuration) (jwt.Service, error) {
				return tg, nil
			}).ApplyMethod(reflect.TypeOf(tg), "GenerateToken",
				func(jwt.Service, *models.User) (string, error) {
					return "", nil
				}).ApplyFunc(daos.FindUserByToken, func(token string, ctx context.Context) (*models.User, error) {
				return testutls.MockUser(), nil
			})
		},
	}
}

func loadRefereshTokenCases() []refereshTokenType {
	return []refereshTokenType{
		refreshTokenInvalidCase(),
		// refreshTokenErrorFromConfigCase(),
		refereshTokenerrorWhileGeneratingToken(),
		refereshTokenErrorFromGenerateTokenCase(),
		refreshTokenSuccessCase(),
	}
}

func TestRefreshToken(t *testing.T) {
	cases := loadRefereshTokenCases()
	// Create a new instance of the resolver
	resolver1 := resolver.Resolver{}
	for _, tt := range cases {
		t.Run(
			tt.name,
			func(t *testing.T) {
				// Handle the case where authentication token is invalid
				patches := tt.init()
				// Set up the context with the mock user
				c := context.Background()
				ctx := context.WithValue(c, testutls.UserKey, testutls.MockUser())
				// Call the refresh token mutation with the given arguments and check the response and error against the expected values
				response, err := resolver1.Mutation().
					RefreshToken(ctx, tt.req)
				if tt.wantResp != nil &&
					response != nil {
					tt.wantResp.Token = response.Token
					// Assert that the expected response matches the actual response
					assert.Equal(t, tt.wantResp, response)
				} else {
					// Assert that the expected error value matches the actual error value
					assert.Equal(t, true, strings.Contains(err.Error(), tt.err.Error()))
				}
				if patches != nil {
					patches.Reset()
				}
			},
		)
	}
}

type authorLoginInput struct {
	email, password string
}

type authorTestInput struct {
	name     string
	input    authorLoginInput
	wantErr  bool
	wantResp *fm.LoginResponse
	init     func() *gomonkey.Patches
}

func TestAuthorLogin(t *testing.T) {
	tests := []authorTestInput{
		authorLoginSuccess(),
		authorInvalidEmailOrPassword(),
		authorInvalidEmail(),
		authorErrorGeneratingToken(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			patch := test.init()
			defer func() {
				if patch != nil {
					patch.Reset()
				}
			}()
			resolver := resolver.Resolver{}
			resp, err := resolver.Mutation().AuthorLogin(
				context.Background(), test.input.email, test.input.password,
			)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err)
			if err == nil {
				assert.Equal(t, test.wantResp, resp)
			}
		})
	}
}

func authorLoginSuccess() authorTestInput {
	return authorTestInput{
		name:    "Should login the author",
		input:   authorLoginInput{email: "user.email@domain.com", password: "simple-password"},
		wantErr: false,
		wantResp: &fm.LoginResponse{
			Token:        "token",
			RefreshToken: "refresh token",
		},
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindAuthorByEmail,
				func(ctx context.Context, email string) (*models.Author, error) {
					return &models.Author{
						ID:        3,
						FirstName: "User",
						LastName:  null.StringFrom("Demo"),
						Email:     "user.email@domain.com",
						Password:  "simple-password",
					}, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return true
				}).ApplyMethod(reflect.TypeOf(jwt.Service{}),
				"GenerateTokenForAuthor",
				func(svc jwt.Service, a *models.Author) (string, error) {
					return "token", nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"Token",
				func(sec secure.Service, str string) string {
					return "refresh token"
				})
		},
	}
}

func authorInvalidEmailOrPassword() authorTestInput {
	return authorTestInput{
		name:     "Should return with message invalid email or password",
		input:    authorLoginInput{email: "user.email@domain.com", password: "simple-password"},
		wantErr:  true,
		wantResp: nil,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindAuthorByEmail,
				func(ctx context.Context, email string) (*models.Author, error) {
					return &models.Author{
						ID:        3,
						FirstName: "User",
						LastName:  null.StringFrom("Demo"),
						Email:     "user.email@domain.com",
						Password:  "different password then input",
					}, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return false
				})
		},
	}
}

func authorInvalidEmail() authorTestInput {
	return authorTestInput{
		name:     "Should return with message invalid email or password",
		input:    authorLoginInput{email: "user.email@domain.com", password: "simple-password"},
		wantErr:  true,
		wantResp: nil,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindAuthorByEmail,
				func(ctx context.Context, email string) (*models.Author, error) {
					return nil, fmt.Errorf("no such user")
				})
		},
	}
}

func authorErrorGeneratingToken() authorTestInput {
	return authorTestInput{
		name:     "Should login the author",
		input:    authorLoginInput{email: "user.email@domain.com", password: "simple-password"},
		wantErr:  true,
		wantResp: nil,
		init: func() *gomonkey.Patches {
			return gomonkey.ApplyFunc(
				daos.FindAuthorByEmail,
				func(ctx context.Context, email string) (*models.Author, error) {
					return &models.Author{
						ID:        3,
						FirstName: "User",
						LastName:  null.StringFrom("Demo"),
						Email:     "user.email@domain.com",
						Password:  "simple-password",
					}, nil
				}).ApplyMethod(reflect.TypeOf(secure.Service{}),
				"HashMatchesPassword",
				func(sec secure.Service, hash string, password string) bool {
					return true
				}).ApplyMethod(reflect.TypeOf(jwt.Service{}),
				"GenerateTokenForAuthor",
				func(svc jwt.Service, a *models.Author) (string, error) {
					return "", fmt.Errorf("couldn't generate token")
				})
		},
	}
}
