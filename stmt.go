// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package alticli

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/bulenttokuzlu/alticli/api"
)

type Stmt struct {
	c     *Conn
	query string
	os    *ODBCStmt
	mu    sync.Mutex
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	fmt.Println("-------------------------Prepare------------------------------", query)

	//--------
	connJson, _ := json.Marshal(c)
	fmt.Println("connJson = ", string(connJson))
	//--------

	if c.bad {
		return nil, driver.ErrBadConn
	}
	os, err := c.PrepareODBCStmt(query)

	//--------
	osJson, _ := json.Marshal(os)
	fmt.Println("osJson = ", string(osJson))
	//--------

	if err != nil {
		return nil, err
	}

	//--------
	stmtJson, _ := json.Marshal(&Stmt{c: c, os: os, query: query})
	fmt.Println("stmtJson = ", string(stmtJson))
	//--------

	return &Stmt{c: c, os: os, query: query}, nil
}

func (s *Stmt) NumInput() int {
	fmt.Println("-------------------------NumInput------------------------------")
	return -1
	/*if s.os == nil {
		return -1
	}
	return len(s.os.Parameters)*/
}

func (s *Stmt) Close() error {
	fmt.Println("-------------------------Close------------------------------")
	if s.os == nil {
		return errors.New("Stmt is already closed")
	}
	ret := s.os.closeByStmt()
	s.os = nil
	return ret
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	fmt.Println("-------------------------Exec------------------------------")
	if s.os == nil {
		return nil, errors.New("Stmt is closed")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.os.usedByRows {
		s.os.closeByStmt()
		s.os = nil
		os, err := s.c.PrepareODBCStmt(s.query)
		if err != nil {
			return nil, err
		}
		s.os = os
	}
	err := s.os.Exec(args)
	if err != nil {
		return nil, err
	}
	var c api.SQLLEN
	ret := api.SQLRowCount(s.os.h, &c)
	if IsError(ret) {
		return nil, NewError("SQLRowCount", s.os.h)
	}
	return &Result{rowCount: int64(c)}, nil
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	fmt.Println("-------------------------Query------------------------------")
	if s.os == nil {
		return nil, errors.New("Stmt is closed")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.os.usedByRows {
		s.os.closeByStmt()
		s.os = nil
		os, err := s.c.PrepareODBCStmt(s.query)
		if err != nil {
			return nil, err
		}
		s.os = os
	}
	err := s.os.Exec(args)
	if err != nil {
		return nil, err
	}
	err = s.os.BindColumns()
	if err != nil {
		return nil, err
	}
	s.os.usedByRows = true // now both Stmt and Rows refer to it
	return &Rows{os: s.os}, nil
}
