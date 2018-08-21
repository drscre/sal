// Code generated by SalGen. DO NOT EDIT.
package actsal

import (
	"context"
	"database/sql"
	"github.com/go-gad/sal"
	"github.com/go-gad/sal/examples/bookstore1"
	"github.com/pkg/errors"
)

type SalStoreClient struct {
	DB *sql.DB
}

func NewStoreClient(db *sql.DB) *SalStoreClient {
	return &SalStoreClient{DB: db}
}

func (s *SalStoreClient) CreateAuthor(ctx context.Context, req bookstore1.CreateAuthorReq) (*bookstore1.CreateAuthorResp, error) {
	var reqMap = make(sal.KeysIntf)
	reqMap["Name"] = &req.Name
	reqMap["Desc"] = &req.Desc
	pgQuery, args := sal.ProcessQueryAndArgs(req.Query(), reqMap)

	rows, err := s.DB.Query(pgQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute Query")
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch columns")
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, errors.Wrap(err, "rows error")
		}
		return nil, sql.ErrNoRows
	}

	var resp bookstore1.CreateAuthorResp
	var mm = make(sal.KeysIntf)
	mm["ID"] = &resp.ID
	mm["CreatedAt"] = &resp.CreatedAt
	var dest = make([]interface{}, 0, len(mm))
	for _, v := range cols {
		if intr, ok := mm[v]; ok {
			dest = append(dest, intr)
		}
	}

	if err = rows.Scan(dest...); err != nil {
		return nil, errors.Wrap(err, "failed to scan row")
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "something failed during iteration")
	}

	return &resp, nil
}

func (s *SalStoreClient) GetAuthors(ctx context.Context, req bookstore1.GetAuthorsReq) ([]*bookstore1.GetAuthorsResp, error) {
	var reqMap = make(sal.KeysIntf)
	reqMap["ID"] = &req.ID
	pgQuery, args := sal.ProcessQueryAndArgs(req.Query(), reqMap)

	rows, err := s.DB.Query(pgQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute Query")
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch columns")
	}

	var list = make([]*bookstore1.GetAuthorsResp, 0)

	for rows.Next() {
		var resp bookstore1.GetAuthorsResp
		var mm = make(sal.KeysIntf)
		mm["ID"] = &resp.ID
		mm["CreatedAt"] = &resp.CreatedAt
		mm["Name"] = &resp.Name
		mm["Desc"] = &resp.Desc
		var dest = make([]interface{}, 0, len(mm))
		for _, v := range cols {
			if intr, ok := mm[v]; ok {
				dest = append(dest, intr)
			}
		}

		if err = rows.Scan(dest...); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		list = append(list, &resp)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "something failed during iteration")
	}

	return list, nil
}

func (s *SalStoreClient) UpdateAuthor(ctx context.Context, req *bookstore1.UpdateAuthorReq) error {
	var reqMap = make(sal.KeysIntf)
	reqMap["ID"] = &req.ID
	reqMap["Name"] = &req.Name
	reqMap["Desc"] = &req.Desc
	pgQuery, args := sal.ProcessQueryAndArgs(req.Query(), reqMap)

	_, err := s.DB.Exec(pgQuery, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute Exec")
	}

	return nil
}
