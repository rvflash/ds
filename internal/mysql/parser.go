// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mysql

import (
	"io"

	"github.com/xwb1989/sqlparser"
)

// Parse parses the given SQL statements as MySQL queries.
// It tries to convert it as a Storage.
func Parse(r io.Reader) (Storage, error) {
	var (
		cur = DefaultDatabaseName
		dbs = Storage{}
		tkz = sqlparser.NewTokenizer(r)
	)
	for {
		stmt, err := sqlparser.ParseNext(tkz)
		if err == io.EOF {
			break
		}
		switch stmt := stmt.(type) {
		case *sqlparser.DBDDL:
			switch stmt.Action {
			case sqlparser.CreateStr:
				dbs, cur = dbs.addDatabase(stmt.DBName, stmt.Charset)
			case sqlparser.DropStr:
				dbs = dbs.dropDatabase(stmt.DBName)
			}
		case *sqlparser.DDL:
			// By default, if no database are specified, we use a default one to wrap any tables.
			if cur == DefaultDatabaseName {
				dbs, cur = dbs.addDatabase(cur, DefaultCharset)
			}
			switch stmt.Action {
			case sqlparser.CreateStr:
				err = dbs.createTable(cur, stmt)
			case sqlparser.AlterStr:
				err = dbs.AlterTable(cur, stmt)
			case sqlparser.DropStr:
				err = dbs.DropTable(cur, stmt)
			case sqlparser.RenameStr:
				err = dbs.RenameTable(cur, stmt)
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return dbs, nil
}
