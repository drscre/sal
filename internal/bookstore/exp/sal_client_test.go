package exp

import (
	"context"
	"testing"
	"time"

	"github.com/go-gad/sal/internal/bookstore"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSalStoreClient_CreateAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	client := NewStoreClient(db)

	req := bookstore.CreateAuthorReq{Name: "foo", Desc: "Bar"}

	rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now().Truncate(time.Millisecond))
	mock.ExpectQuery(`INSERT INTO authors .+`).WithArgs(req.Name, req.Desc).WillReturnRows(rows)

	resp, err := client.CreateAuthor(context.Background(), &req)
	t.Logf("resp %#v err: %s", resp, err)
}