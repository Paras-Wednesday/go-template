package daos

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go-template/models"
)

func CreatePostTx(ctx context.Context, post models.Post, tx *sql.Tx) (models.Post, error) {
	contextExecutor := GetContextExecutor(tx)

	err := post.Insert(ctx, contextExecutor, boil.Infer())
	return post, err
}

func CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	return CreatePostTx(ctx, post, nil)
}

func UpdatePostTx(ctx context.Context, post models.Post, tx *sql.Tx) (models.Post, error) {
	contextExecutor := GetContextExecutor(tx)

	_, err := post.Update(ctx, contextExecutor, boil.Infer())
	return post, err
}

func UpdatePost(ctx context.Context, post models.Post) (models.Post, error) {
	return UpdatePostTx(ctx, post, nil)
}

func DeletePostTx(ctx context.Context, post models.Post, tx *sql.Tx) (int64, error) {
	contextExecutor := GetContextExecutor(tx)

	return post.Delete(ctx, contextExecutor)
}

func DeletePost(ctx context.Context, post models.Post) (int64, error) {
	return DeletePostTx(ctx, post, nil)
}

func FindPostForAuthorByID(ctx context.Context, authorID int, postID int) (*models.Post, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.Posts(
		models.PostWhere.AuthorID.EQ(authorID), models.PostWhere.ID.EQ(postID)).
		One(ctx, contextExecutor)
}

func FindAllPostBylAuthorWithCount(ctx context.Context, authorID int, queries ...qm.QueryMod) (models.PostSlice, int64, error) {
	contextExecutor := GetContextExecutor(nil)

	authQueryMode := models.PostWhere.AuthorID.EQ(authorID)
	count, err := models.Posts(authQueryMode).Count(ctx, contextExecutor)
	if err != nil {
		return models.PostSlice{}, 0, err
	}

	// accumulate all the queries including the author ID filter
	queries = append(queries, authQueryMode)

	posts, err := models.Posts(queries...).All(ctx, contextExecutor)
	return posts, count, err
}
