// Code generated by SalGen. DO NOT EDIT.
package repo

import (
	"context"
	"database/sql"

	"github.com/go-gad/sal"
	"github.com/go-gad/sal/examples/bookstore"
	"github.com/pkg/errors"
)

type SalStore struct {
	handler  sal.QueryHandler
	ctrl     *sal.Controller
	txOpened bool
}

func NewStore(h sal.QueryHandler, options ...sal.ClientOption) *SalStore {
	s := &SalStore{
		handler:  h,
		ctrl:     sal.NewController(options...),
		txOpened: false,
	}

	return s
}

func (s *SalStore) BeginTx(ctx context.Context, opts *sql.TxOptions) (bookstore.Store, error) {
	dbConn, ok := s.handler.(sal.TransactionBegin)
	if !ok {
		return nil, errors.New("oops")
	}

	ctx = context.WithValue(ctx, sal.ContextKeyTxOpened, s.txOpened)
	ctx = context.WithValue(ctx, sal.ContextKeyOperationType, "Begin")

	// todo middleware
	tx, err := dbConn.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start tx")
	}

	newClient := &SalStore{
		handler:  tx,
		ctrl:     s.ctrl,
		txOpened: true,
	}

	return newClient, nil
}

func (s *SalStore) Tx() sal.TxHandler {
	if tx, ok := s.handler.(sal.TxHandler); ok {
		return tx
	}
	return nil
}

func (s *SalStore) CreateAuthor(ctx context.Context, req bookstore.CreateAuthorReq) (*bookstore.CreateAuthorResp, error) {
	var (
		err      error
		rawQuery = req.Query()
		reqMap   = make(sal.RowMap)
	)
	reqMap["Name"] = &req.Name
	reqMap["Desc"] = &req.Desc

	ctx = context.WithValue(ctx, sal.ContextKeyTxOpened, s.txOpened)
	ctx = context.WithValue(ctx, sal.ContextKeyOperationType, "QueryRow")

	pgQuery, args := sal.ProcessQueryAndArgs(rawQuery, reqMap)

	stmt, err := s.ctrl.PrepareStmt(ctx, s.handler, pgQuery)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, fn := range s.ctrl.BeforeQuery {
		var fnz sal.FinalizerFunc
		ctx, fnz = fn(ctx, rawQuery, req)
		if fnz != nil {
			defer func() { fnz(ctx, err) }()
		}
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute Query")
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch columns")
	}

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, errors.Wrap(err, "rows error")
		}
		return nil, sql.ErrNoRows
	}

	var resp bookstore.CreateAuthorResp
	var respMap = make(sal.RowMap)
	respMap["ID"] = &resp.ID
	respMap["CreatedAt"] = &resp.CreatedAt

	var dest = make([]interface{}, 0, len(respMap))
	for _, v := range cols {
		if intr, ok := respMap[v]; ok {
			dest = append(dest, intr)
		}
	}

	if err = rows.Scan(dest...); err != nil {
		return nil, errors.Wrap(err, "failed to scan row")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "something failed during iteration")
	}

	return &resp, nil
}

func (s *SalStore) GetAuthors(ctx context.Context, req bookstore.GetAuthorsReq) ([]*bookstore.GetAuthorsResp, error) {
	var (
		err      error
		rawQuery = req.Query()
		reqMap   = make(sal.RowMap)
	)
	reqMap["id"] = &req.ID
	reqMap["tags"] = &req.Tags

	req.ProcessRow(reqMap)

	ctx = context.WithValue(ctx, sal.ContextKeyTxOpened, s.txOpened)
	ctx = context.WithValue(ctx, sal.ContextKeyOperationType, "Query")

	pgQuery, args := sal.ProcessQueryAndArgs(rawQuery, reqMap)

	stmt, err := s.ctrl.PrepareStmt(ctx, s.handler, pgQuery)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, fn := range s.ctrl.BeforeQuery {
		var fnz sal.FinalizerFunc
		ctx, fnz = fn(ctx, rawQuery, req)
		if fnz != nil {
			defer func() { fnz(ctx, err) }()
		}
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute Query")
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch columns")
	}

	var list = make([]*bookstore.GetAuthorsResp, 0)

	for rows.Next() {
		var resp bookstore.GetAuthorsResp
		var respMap = make(sal.RowMap)
		respMap["id"] = &resp.ID
		respMap["created_at"] = &resp.CreatedAt
		respMap["name"] = &resp.Name
		respMap["desc"] = &resp.Desc
		respMap["tags"] = &resp.Tags

		resp.ProcessRow(respMap)

		var dest = make([]interface{}, 0, len(respMap))
		for _, v := range cols {
			if intr, ok := respMap[v]; ok {
				dest = append(dest, intr)
			}
		}

		if err = rows.Scan(dest...); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		list = append(list, &resp)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "something failed during iteration")
	}

	return list, nil
}

func (s *SalStore) UpdateAuthor(ctx context.Context, req *bookstore.UpdateAuthorReq) error {
	var (
		err      error
		rawQuery = req.Query()
		reqMap   = make(sal.RowMap)
	)
	reqMap["ID"] = &req.ID
	reqMap["Name"] = &req.Name
	reqMap["Desc"] = &req.Desc

	ctx = context.WithValue(ctx, sal.ContextKeyTxOpened, s.txOpened)
	ctx = context.WithValue(ctx, sal.ContextKeyOperationType, "Exec")

	pgQuery, args := sal.ProcessQueryAndArgs(rawQuery, reqMap)

	stmt, err := s.ctrl.PrepareStmt(ctx, s.handler, pgQuery)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, fn := range s.ctrl.BeforeQuery {
		var fnz sal.FinalizerFunc
		ctx, fnz = fn(ctx, rawQuery, req)
		if fnz != nil {
			defer func() { fnz(ctx, err) }()
		}
	}

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute Exec")
	}

	return nil
}
