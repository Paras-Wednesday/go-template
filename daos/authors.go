package daos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go-template/models"
)

func CreateAuthorTx(ctx context.Context, author models.Author, tx *sql.Tx) (models.Author, error) {
	contextExecutor := GetContextExecutor(tx)

	err := author.Insert(ctx, contextExecutor, boil.Infer())
	return author, err
}

func CreateAuthor(ctx context.Context, author models.Author) (models.Author, error) {
	return CreateAuthorTx(ctx, author, nil)
}

func UpdateAuthor(ctx context.Context, author models.Author) (models.Author, error) {
	contextExecutor := GetContextExecutor(nil)

	_, err := author.Update(ctx, contextExecutor, boil.Infer())
	return author, err
}

func DeleteAuthor(ctx context.Context, author models.Author) (int64, error) {
	contextExecutor := GetContextExecutor(nil)

	return author.Delete(ctx, contextExecutor)
}

func FindAuthorByID(ctx context.Context, id int) (*models.Author, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.FindAuthor(ctx, contextExecutor, id)
}

func FindAuthorByFirstName(ctx context.Context, fname string) (*models.Author, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.Authors(qm.Where(fmt.Sprintf("%s=$1", models.AuthorColumns.FirstName), fname)).
		One(ctx, contextExecutor)
}

func FindAuthorByLastName(ctx context.Context, lname string) (*models.Author, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.Authors(qm.Where(fmt.Sprintf("%s=$1", models.AuthorColumns.LastName), lname)).
		One(ctx, contextExecutor)
}

func FindAuthorByEmail(ctx context.Context, email string) (*models.Author, error) {
	contextExecutor := GetContextExecutor(nil)
	return models.Authors(models.AuthorWhere.Email.EQ(email)).One(ctx, contextExecutor)
}

func GetAllAuthorsWithCount(ctx context.Context, queries ...qm.QueryMod) (models.AuthorSlice, int64, error) {
	contextExecutor := GetContextExecutor(nil)

	count, err := models.Authors().Count(ctx, contextExecutor)
	if err != nil {
		return models.AuthorSlice{}, 0, err
	}
	authors, err := models.Authors(queries...).All(ctx, contextExecutor)
	return authors, count, err
}
