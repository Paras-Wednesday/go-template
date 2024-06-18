package daos_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"go-template/daos"
	"go-template/models"
	"go-template/testutls"
)

func TestCreateAuthor(t *testing.T) {
	tests := []struct {
		name    string
		author  models.Author
		wantErr bool
	}{
		{
			name: "should create a author",
			author: models.Author{
				FirstName: testutls.MockAuthor().FirstName,
				LastName:  testutls.MockAuthor().LastName,
			},
			wantErr: false,
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	// As the insertion will return the "id", and "deleted_at"
	// column which are non default column
	rows := sqlmock.NewRows([]string{
		models.AuthorColumns.ID,
		models.AuthorColumns.DeletedAt,
	}).AddRow(
		testutls.MockAuthor().ID,
		testutls.MockAuthor().DeletedAt,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "authors"`)).
		WithArgs().
		WillReturnRows(rows)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotAuthor, err := daos.CreateAuthor(context.Background(), test.author)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err,
			)
			assert.Equal(t, test.author.FirstName, gotAuthor.FirstName)
			assert.Equal(t, test.author.LastName, gotAuthor.LastName)
		})
	}
}

func TestFindAuthorByID(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		wantErr bool
		init    func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "Should return author with id 1",
			userID:  1,
			wantErr: false,
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{models.AuthorColumns.ID}).AddRow(testutls.MockAuthor().ID)
				mock.ExpectQuery(regexp.QuoteMeta(
					`select * from "authors" where "id"=$1`,
				)).WithArgs().WillReturnRows(rows)
			},
		},
		{
			name:    "Should return error on finding user",
			userID:  1,
			wantErr: true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`select * from "authors" where "id"=$1`,
				)).WithArgs().WillReturnError(
					fmt.Errorf("couldn't find user"),
				)
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		t.Run(test.name, func(t *testing.T) {
			author, err := daos.FindAuthorByID(context.Background(), test.userID)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t got: %v", test.wantErr, err)
			if author != nil {
				assert.Equal(t, test.userID, author.ID)
			}
		})
	}
}

func TestFindAuthorByFirstName(t *testing.T) {
	tests := []struct {
		testName  string
		firstName string
		wantErr   bool
		init      func(mock sqlmock.Sqlmock)
	}{
		{
			testName:  "Should return error in finding author",
			firstName: "abc",
			wantErr:   true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "authors".* FROM "authors"
					 WHERE (first_name=$1) LIMIT 1;`)).
					WithArgs().
					WillReturnError(fmt.Errorf("No such user"))
			},
		},
		{
			testName:  "Should return author with first_name \"abc\"",
			firstName: "abc",
			wantErr:   false,
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					models.AuthorColumns.ID,
					models.AuthorColumns.FirstName,
				}).AddRow(
					testutls.MockAuthor().ID,
					testutls.MockAuthor().FirstName,
				)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "authors".* FROM "authors"
					 WHERE (first_name=$1) LIMIT 1;`,
				)).
					WithArgs().
					WillReturnRows(rows)
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		t.Run(test.testName, func(t *testing.T) {
			author, err := daos.FindAuthorByFirstName(context.Background(), test.firstName)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err)
			t.Logf("author: %+v", author)
		})
	}
}

func TestFindAuthorByLastName(t *testing.T) {
	tests := []struct {
		testName string
		lastName string
		wantErr  bool
		init     func(mock sqlmock.Sqlmock)
	}{
		{
			testName: "Should return error in finding author",
			lastName: "abc",
			wantErr:  true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "authors".* FROM "authors"
					 WHERE (last_name=$1) LIMIT 1;`,
				)).
					WithArgs().
					WillReturnError(fmt.Errorf("No such user"))
			},
		},
		{
			testName: "Should return author with last_name \"abc\"",
			lastName: "abc",
			wantErr:  false,
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					models.AuthorColumns.ID,
					models.AuthorColumns.LastName,
				}).AddRow(
					testutls.MockAuthor().ID,
					testutls.MockAuthor().LastName,
				)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "authors".* FROM "authors" WHERE (last_name=$1) LIMIT 1;`,
				)).
					WithArgs().
					WillReturnRows(rows)
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		t.Run(test.testName, func(t *testing.T) {
			author, err := daos.FindAuthorByLastName(context.Background(), test.lastName)
			t.Logf("author: %+v", author)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err)
		})
	}
}

func TestUpdateAuthor(t *testing.T) {
	cases := []struct {
		name    string
		req     models.Author
		wantErr bool
	}{
		{
			name:    "Should update Author",
			req:     models.Author{},
			wantErr: false,
		},
	}
	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()
	for _, test := range cases {
		result := driver.Result(driver.RowsAffected(1))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "authors"`)).
			WillReturnResult(result)

		t.Run(test.name, func(t *testing.T) {
			_, err := daos.UpdateAuthor(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err)
		})
	}
}

func TestDeleteAuthor(t *testing.T) {
	cases := []struct {
		name    string
		req     models.Author
		wantErr bool
	}{
		{
			name:    "Should delete author",
			req:     models.Author{},
			wantErr: false,
		},
	}
	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()
	for _, test := range cases {
		result := driver.Result(driver.RowsAffected(1))
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM "authors" WHERE "id"=$1`,
		)).
			WillReturnResult(result)
		t.Run(test.name, func(t *testing.T) {
			_, err := daos.DeleteAuthor(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err)
		})
	}
}

//nolint:funlen
func TestGetAllAuthors(t *testing.T) {
	tests := []struct {
		name      string
		wantErr   bool
		wantCount int
		init      func(mock sqlmock.Sqlmock)
	}{
		{
			name:      "Should retrieve 2 authors",
			wantErr:   false,
			wantCount: 2,
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					models.AuthorColumns.ID,
					models.AuthorColumns.FirstName,
					models.AuthorColumns.LastName,
				}).AddRow(
					testutls.MockAuthor().ID,
					testutls.MockAuthor().FirstName,
					testutls.MockAuthor().LastName,
				).AddRow(
					testutls.MockAuthor().ID,
					testutls.MockAuthor().FirstName,
					testutls.MockAuthor().LastName,
				)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "authors".* FROM "authors";`,
				)).WithArgs().
					WillReturnRows(rows)
			},
		},
		{
			name:      "Should retrieve only 1 author",
			wantErr:   false,
			wantCount: 1,
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					models.AuthorColumns.ID,
					models.AuthorColumns.FirstName,
					models.AuthorColumns.LastName,
				}).AddRow(
					testutls.MockAuthor().ID,
					testutls.MockAuthor().FirstName,
					testutls.MockAuthor().LastName,
				)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "authors".* FROM "authors";`,
				)).WithArgs().
					WillReturnRows(rows)
			},
		},
		{
			name:      "Should retrieve no authors",
			wantErr:   true,
			wantCount: 0,
			init: func(mock sqlmock.Sqlmock) {
				sqlmock.NewRows([]string{}).
					AddRow()
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "authors".* FROM "authors";`,
				)).WithArgs().
					WillReturnError(fmt.Errorf("error in fetching "))
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		authors, err := daos.GetAllAuthors(context.Background())
		assert.Equal(t, test.wantErr, err != nil,
			"wantErr: %t, got: %v", test.wantErr, err)
		assert.Equal(t, test.wantCount, len(authors))
	}
}
