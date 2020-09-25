// Package api Docs
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package api

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
	goboiler "github.com/wednesday-solutions/go-boiler"
	"os"

	"github.com/joho/godotenv"
	"github.com/volatiletech/sqlboiler/boil"
	graphql "github.com/wednesday-solutions/go-boiler/graphql_models"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/postgres"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/zlog"

	"github.com/wednesday-solutions/go-boiler/pkg/api/auth"
	al "github.com/wednesday-solutions/go-boiler/pkg/api/auth/logging"
	at "github.com/wednesday-solutions/go-boiler/pkg/api/auth/transport"
	"github.com/wednesday-solutions/go-boiler/pkg/api/password"
	pl "github.com/wednesday-solutions/go-boiler/pkg/api/password/logging"
	pt "github.com/wednesday-solutions/go-boiler/pkg/api/password/transport"
	"github.com/wednesday-solutions/go-boiler/pkg/api/user"
	ul "github.com/wednesday-solutions/go-boiler/pkg/api/user/logging"
	ut "github.com/wednesday-solutions/go-boiler/pkg/api/user/transport"

	_ "github.com/lib/pq" // here
	"github.com/wednesday-solutions/go-boiler/pkg/utl/config"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/jwt"
	authMw "github.com/wednesday-solutions/go-boiler/pkg/utl/middleware/auth"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/secure"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {

	err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
	if err != nil {
		fmt.Print("Error loading .env file")
		return err
	}
	db, err := postgres.Connect()
	if err != nil {
		return err
	}
	boil.SetDB(db)

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	authMiddleware := authMw.Middleware(jwt)

	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec), log), e, authMiddleware)

	v1 := e.Group("/v1")
	v1.Use(authMiddleware)

	ut.NewHTTP(ul.New(user.Initialize(db, sec), log), v1)
	pt.NewHTTP(pl.New(password.Initialize(db, sec), log), v1)

	graphqlHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &goboiler.Resolver{}}))
	playgroundHandler := playground.Handler("GraphQL playground", "/graphql")

	// graphql apis
	e.POST("/graphql", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		graphqlHandler.ServeHTTP(res, req)
		return nil
	}, authMiddleware)

	// graphql playground
	e.GET("/playground", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		playgroundHandler.ServeHTTP(res, req)
		return nil
	})
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})
	return nil
}

// CustomContext ...
type CustomContext struct {
	echo.Context
	ctx context.Context
}

// Foo ...
func (c *CustomContext) Foo() {
	println("foo")
}

// Bar ...
func (c *CustomContext) Bar() {
	println("bar")
}

type key string

const (
	echoContextKey key = "EchoContextKey"
	// ...
)

// Process is the middleware function.
func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), echoContextKey, c)
		c.SetRequest(c.Request().WithContext(ctx))

		cc := &CustomContext{c, ctx}

		return next(cc)
	}
}
