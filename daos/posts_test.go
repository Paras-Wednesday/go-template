package daos_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go-template/daos"
	"go-template/models"
	"go-template/testutls"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name    string
		input   models.Post
		wantErr bool
		init    func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Should create a post",
			input: models.Post{
				AuthorID: 1,
				Content:  "This is my first post",
			},
			wantErr: false,
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					models.PostColumns.ID,
					models.PostColumns.DeletedAt,
				}).AddRow(1, null.Time{})
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "posts"`,
				)).WithArgs().
					WillReturnRows(rows)
			},
		},
		{
			name: "Should not create a post",
			input: models.Post{
				AuthorID: 1,
				Content:  "This is my first post",
			},
			wantErr: true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "posts"`,
				)).WithArgs().
					WillReturnError(fmt.Errorf("error inserting in opst"))
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		t.Run(test.name, func(t *testing.T) {
			_, err := daos.CreatePost(context.Background(), test.input)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err,
			)
		})
	}
}

func TestUpdatePost(t *testing.T) {
	tests := []struct {
		name    string
		input   models.Post
		wantErr bool
		init    func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "Should update the Post",
			input:   models.Post{},
			wantErr: false,
			init: func(mock sqlmock.Sqlmock) {
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts"`)).
					WillReturnResult(result)
			},
		},
		{
			name:    "Should not update the Post",
			input:   models.Post{},
			wantErr: true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts"`)).
					WillReturnError(fmt.Errorf("couldn't update the psot"))
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		t.Run(test.name, func(t *testing.T) {
			_, err := daos.UpdatePost(context.Background(), test.input)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err,
			)
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name    string
		input   models.Post
		wantErr bool
		init    func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "Should delete the Post",
			input:   models.Post{},
			wantErr: false,
			init: func(mock sqlmock.Sqlmock) {
				result := driver.Result(driver.RowsAffected(1))
				mock.ExpectExec(regexp.QuoteMeta(
					`DELETE FROM "posts" WHERE "id"=$1`,
				)).
					WillReturnResult(result)
			},
		},
		{
			name:    "Should not delete the Post",
			input:   models.Post{},
			wantErr: true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE "posts"`)).
					WillReturnError(fmt.Errorf("couldn't delete the post"))
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		t.Run(test.name, func(t *testing.T) {
			_, err := daos.DeletePost(context.Background(), test.input)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err,
			)
		})
	}
}

func TestFindPostById(t *testing.T) {
	tests := []struct {
		name     string
		postID   int
		authorID int
		wantErr  bool
		init     func(mock sqlmock.Sqlmock)
	}{
		{
			name:     "Should return post with id 1",
			postID:   1,
			authorID: 1,
			wantErr:  false,
			init: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					models.AuthorColumns.ID,
				}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "posts".* FROM "posts" 
                    WHERE ("posts"."author_id" = $1) AND ("posts"."id" = $2) 
                    LIMIT 1;`,
				)).WillReturnRows(rows)
			},
		},
		{
			name:     "Should return error on no such post",
			postID:   1,
			authorID: 0,
			wantErr:  true,
			init: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "posts".* FROM "posts" 
                    WHERE ("posts"."author_id" = $1) AND ("posts"."id" = $2) 
                    LIMIT 1;`,
				)).WillReturnError(fmt.Errorf("no such post"))
			},
		},
	}
	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		test.init(mock)
		t.Run(test.name, func(t *testing.T) {
			post, err := daos.FindPostForAuthorByID(context.Background(), test.authorID, test.postID)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err,
			)
			if post != nil {
				assert.Equal(t, test.postID, post.ID)
			}
		})
	}
}

// nolint: funlen
func TestFindAllPostsByAuthorWithCount(t *testing.T) {
	tests := []struct {
		name      string
		authorID  int
		qms       []qm.QueryMod
		wantErr   bool
		wantCount int64
		queries   []testutls.QueryData
	}{
		{
			name:     "Retrieve 5 posts by author",
			authorID: 1,
			qms: []qm.QueryMod{
				qm.Limit(5),
				qm.Offset(1),
			},
			wantErr:   false,
			wantCount: 5,
			queries: []testutls.QueryData{
				{
					Query:      `SELECT COUNT(*) FROM "posts" WHERE ("posts"."author_id" = $1);`,
					DbResponse: sqlmock.NewRows([]string{"count"}).AddRow(5),
				},
				{
					Query: `SELECT "posts".* FROM "posts"
					 WHERE ("posts"."author_id" = $1) LIMIT 5 OFFSET 1;`,
					DbResponse: sqlmock.NewRows([]string{"id", "author_id"}).
						AddRow(1, 1).AddRow(1, 2).AddRow(1, 3).AddRow(1, 4).AddRow(1, 5),
				},
			},
		},
		{
			name:     "Retrieve 0 post by author",
			authorID: 90,
			qms: []qm.QueryMod{
				qm.Limit(5),
				qm.Offset(50),
			},
			wantErr:   false,
			wantCount: 5,
			queries: []testutls.QueryData{
				{
					Query:      `SELECT COUNT(*) FROM "posts" WHERE ("posts"."author_id" = $1);`,
					DbResponse: sqlmock.NewRows([]string{"count"}).AddRow(5),
				},
				{
					Query: `SELECT "posts".* FROM "posts"
					 WHERE ("posts"."author_id" = $1) LIMIT 5 OFFSET 50;`,
					DbResponse: sqlmock.NewRows([]string{"id", "author_id"}),
				},
			},
		},
	}

	mock, cleanup, _ := testutls.SetupMockDB(t)
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.wantErr {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT "posts".* FROM "posts"
					 WHERE "author_id"=$1;`,
				)).
					WithArgs().
					WillReturnError(fmt.Errorf("this is some error"))
			}

			for _, dbQuery := range test.queries {
				mock.ExpectQuery(regexp.QuoteMeta(dbQuery.Query)).
					WithArgs().
					WillReturnRows(dbQuery.DbResponse)
			}
			t.Run(test.name, func(t *testing.T) {
				posts, count, err := daos.FindAllPostBylAuthorWithCount(
					context.Background(), test.authorID, test.qms...,
				)
				t.Logf("posts: %+v", posts)

				assert.Equal(t, test.wantErr, err != nil,
					"wantErr: %t, got: %v", test.wantErr, err)
				if err == nil {
					assert.Equal(t, test.wantCount, count)
				}
			})
		})
	}
}
