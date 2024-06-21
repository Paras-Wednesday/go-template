package auth

import (
	"context"
	"reflect"

	graphql2 "github.com/99designs/gqlgen/graphql"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/vektah/gqlparser/v2/ast"

	"go-template/daos"
	"go-template/models"
	resultwrapper "go-template/pkg/utl/resultwrapper"
)

type key string

const (
	authorization key = "Authorization"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}

// CustomContext ...
type CustomContext struct {
	echo.Context
	ctx context.Context
}

var (
	UserCtxKey   = &ContextKey{"user"}
	AuthorCtxKey = &ContextKey{"author"}
)

type ContextKey struct {
	Name string
}

// FromContext finds the user from the context. REQUIRES Middleware to have run.
func FromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(UserCtxKey).(*models.User)
	return user
}

// UserIDFromContext ...
func UserIDFromContext(ctx context.Context) int {
	user := FromContext(ctx)
	if user != nil {
		return user.ID
	}
	return 0
}

func AuthorIDFromContext(ctx context.Context) int {
	author, ok := ctx.Value(AuthorCtxKey).(*models.Author)
	if !ok || author == nil {
		return 0
	}
	return author.ID
}

// GqlMiddleware ...
func GqlMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(
				c.Request().Context(),
				authorization,
				c.Request().Header.Get(string(authorization)),
			)
			c.SetRequest(c.Request().WithContext(ctx))
			cc := &CustomContext{c, ctx}
			return next(cc)
		}
	}
}

// WhiteListedOperations...
var WhiteListedOperations = map[string][]string{
	"query":        {"__schema", "introspectionquery", "userNotification"},
	"mutation":     {"login", "authorLogin", "createAuthor"},
	"subscription": {"userNotification"},
}

// AdminOperations...
var AdminOperations = map[string][]string{
	"query": {"users"},
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getAccessNeeds(operation *ast.OperationDefinition) (needsAuthAccess bool, needsSuperAdminAccess bool) {
	operationName := string(operation.Operation)
	for _, selectionSet := range operation.SelectionSet {
		selection := reflect.ValueOf(selectionSet).Elem()
		if !contains(WhiteListedOperations[operationName], selection.FieldByName("Name").Interface().(string)) {
			needsAuthAccess = true
		}
		if contains(AdminOperations[operationName], selection.FieldByName("Name").Interface().(string)) {
			needsSuperAdminAccess = true
		}
	}
	return needsAuthAccess, needsSuperAdminAccess
}

// GraphQLMiddleware ...
func GraphQLMiddleware(
	ctx context.Context,
	tokenParser TokenParser,
	next graphql2.OperationHandler,
) graphql2.ResponseHandler {
	if !needsAuthOrSuperAdminAccess(ctx) {
		return next(ctx)
	}

	// strip token
	tokenStr := ctx.Value(authorization).(string)
	if len(tokenStr) == 0 {
		return resultwrapper.HandleGraphQLError("Authorization header is missing")
	}
	token, err := tokenParser.ParseToken(tokenStr)
	if err != nil || !token.Valid {
		return resultwrapper.HandleGraphQLError("Invalid authorization token")
	}
	claims := token.Claims.(jwt.MapClaims)

	if resp := verifySuperAdminRole(ctx, claims); resp != nil {
		return resp
	}

	email := claims["e"].(string)
	entityType := claims["type"].(string)
	if !isValidEntityType(entityType) {
		return resultwrapper.HandleGraphQLError("Invalid authorization token")
	}
	if entityType == "user" {
		ctx, resp := ctxWithUser(ctx, email)
		if resp != nil {
			return resp
		}
		return next(ctx)
	}
	ctx, resp := ctxWithAuthor(ctx, email)
	if resp != nil {
		return resp
	}
	return next(ctx)
}

// needsAuthOrSuperAdminAccess helper to check if auth or any super admin access
// is needed to perform the given operation
func needsAuthOrSuperAdminAccess(ctx context.Context) bool {
	operation := graphql2.GetOperationContext(ctx).Operation
	needsAuthAccess, needsSuperAdminAccess := getAccessNeeds(operation)
	return needsAuthAccess || needsSuperAdminAccess
}

func verifySuperAdminRole(ctx context.Context, claims jwt.MapClaims) graphql2.ResponseHandler {
	operation := graphql2.GetOperationContext(ctx).Operation
	_, needsSuperAdminAccess := getAccessNeeds(operation)
	role, ok := claims["role"].(string)
	if !ok {
		return resultwrapper.HandleGraphQLError(
			"Unauthorized! \n Only admins are authorized to make this request.",
		)
	}
	if needsSuperAdminAccess && role != "SUPER_ADMIN" {
		return resultwrapper.HandleGraphQLError(
			"Unauthorized! \n Only admins are authorized to make this request.",
		)
	}
	return nil
}

// isValidEntityType will check if given entityType is author or user
func isValidEntityType(entityType string) bool {
	return entityType == "author" || entityType == "user"
}

// ctxWithUser the helper function to GQlMiddleware to set context with user
func ctxWithUser(ctx context.Context, email string) (context.Context, graphql2.ResponseHandler) {
	user, err := daos.FindUserByEmail(email, ctx)
	if err != nil {
		return nil, resultwrapper.HandleGraphQLError("no user found with this email")
	}
	ctx = context.WithValue(ctx, UserCtxKey, user)
	return ctx, nil
}

// ctxWithAuthor the helper function to GQlMiddleware to set context with author
func ctxWithAuthor(ctx context.Context, email string) (context.Context, graphql2.ResponseHandler) {
	author, err := daos.FindAuthorByEmail(ctx, email)
	if err != nil {
		return nil, resultwrapper.HandleGraphQLError("no author found with this email")
	}
	ctx = context.WithValue(ctx, AuthorCtxKey, author)
	return ctx, nil
}
