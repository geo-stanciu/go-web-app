package main

import (
	"database/sql"

	"github.com/geo-stanciu/go-utils/utils"
)

const (
	allOtherRequests string = "--all-rest--"
)

type menuName struct {
	language string
	name     string
}

type menu struct {
	requestURL string
	name       []menuName
	roles      []userRole
}

func addGetRequestsAccessRules(tx *sql.Tx) (bool, error) {
	menus := []*menu{
		{"stop-process",
			[]menuName{{"EN", "Stop Process"}},
			[]userRole{{"All"}},
		},
		{"index",
			[]menuName{{"EN", "Index"}},
			[]userRole{{"Member"}},
		},
		{"users",
			[]menuName{{"EN", "Users"}},
			[]userRole{{"Administrator"}},
		},
		{"about",
			[]menuName{{"EN", "About"}},
			[]userRole{{"Member"}},
		},
		{"login",
			[]menuName{{"EN", "Login"}},
			[]userRole{{"All"}},
		},
		{"register",
			[]menuName{{"EN", "Register"}},
			[]userRole{{"All"}},
		},
		{"change-password",
			[]menuName{{"EN", "Change Password"}},
			[]userRole{{"Member"}},
		},
	}

	foundNew, err := setAccessRules(tx, "GET", menus)
	return foundNew, err
}

func addPostRequestsAccessRules(tx *sql.Tx) (bool, error) {
	menus := []*menu{
		{"login",
			[]menuName{{"EN", "Login"}},
			[]userRole{{"All"}},
		},
		{"logout",
			[]menuName{{"EN", "Logout"}},
			[]userRole{{"All"}},
		},
		{"register",
			[]menuName{{"EN", "Register"}},
			[]userRole{{"All"}},
		},
	}

	foundNew, err := setAccessRules(tx, "POST", menus)
	return foundNew, err
}

func addGeneralMemberRequestsAccessRules(tx *sql.Tx) (bool, error) {
	menus := []*menu{
		{allOtherRequests,
			[]menuName{},
			[]userRole{{"Member"}},
		},
	}

	foundNew, err := setAccessRules(tx, "All", menus)
	return foundNew, err
}

type accessRule struct {
	ParentID int32  `sql:"parent_id"`
	ChildID  int32  `sql:"child_id"`
	RoleID   int32  `sql:"role_id"`
	Role     string `sql:"role"`
}

func addChildRequestAccessRules(tx *sql.Tx) (bool, error) {
	foundNew := false
	var err error

	pq := dbutl.PQuery(`
		WITH child AS (
			SELECT request_id, parent_id
			  FROM request
			 WHERE parent_id IS NOT NULL
		)
		SELECT c.parent_id,
			   c.request_id AS child_id,
			   rr.role_id,
			   r.role
		  FROM request_role rr
		  JOIN role r ON (rr.role_id = r.role_id)
		  JOIN child c ON (rr.request_id = c.parent_id)
		 ORDER BY rr.request_id
	`)

	var access []*accessRule
	err = dbutl.ForEachRowTx(tx, pq, func(row *sql.Rows, sc *utils.SQLScan) error {
		arule := accessRule{}
		err = sc.Scan(dbutl, row, &arule)
		if err != nil {
			return err
		}

		access = append(access, &arule)
		return nil
	})
	if err != nil {
		return false, err
	}

	for _, arule := range access {
		foundNew, err = addRequest2Role(tx, arule.ChildID, arule.Role)
		if err != nil {
			return false, err
		}
	}

	return foundNew, nil
}

func setAccessRules(tx *sql.Tx, reqType string, menus []*menu) (bool, error) {
	found := 0
	var requestID int32
	var err error

	var pq *utils.PreparedQuery
	foundNew := false

	for _, m := range menus {
		if m.requestURL == allOtherRequests {
			// MySQL does not support except or minus queries at this time
			pq = dbutl.PQuery(`
				SELECT r.request_id
				  FROM request r
				  LEFT OUTER JOIN request_role rr ON (r.request_id = rr.request_id)
				 WHERE rr.request_id is null
			`)

			// on PostgreSQL you need to close your active fetch
			// on a transaction
			// before you can work again on it
			// workaround: on PostgreSQL you can use declared cursors
			// for now: read id's in memory and run the next query on array
			// references: https://github.com/lib/pq/issues/81
			//             https://github.com/lib/pq/issues/635
			var reqID []int32
			err = dbutl.ForEachRowTx(tx, pq, func(row *sql.Rows, sc *utils.SQLScan) error {
				err = row.Scan(&requestID)
				if err != nil {
					return err
				}

				reqID = append(reqID, requestID)
				return nil
			})
			if err != nil {
				return false, err
			}

			for _, req := range reqID {
				for _, r := range m.roles {
					foundNew, err = addRequest2Role(tx, req, r.role)
					if err != nil {
						return false, err
					}
				}
			}
		} else {
			pq = dbutl.PQuery(`
				SELECT request_id
				FROM request
				WHERE request_url = ?
				AND request_type = ?
			`, m.requestURL,
				reqType)

			err = tx.QueryRow(pq.Query, pq.Args...).Scan(&requestID)
			if err != nil {
				return false, err
			}

			for _, n := range m.name {
				pq := dbutl.PQuery(`
					SELECT CASE WHEN EXISTS (
						SELECT 1
						FROM request_name
						WHERE request_id = ?
						AND language = ?
					) THEN 1 ELSE 0 END FROM dual
				`, requestID,
					n.language)

				err := tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
				if err != nil {
					return false, err
				}

				if found == 1 {
					continue
				}

				foundNew = true

				pq = dbutl.PQuery(`
					INSERT INTO request_name (
						request_id,
						language,
						name
					)
					VALUES (?, ?, ?)
				`, requestID,
					n.language,
					n.name)

				_, err = dbutl.ExecTx(tx, pq)
				if err != nil {
					return false, err
				}
			}

			for _, r := range m.roles {
				foundNew, err = addRequest2Role(tx, requestID, r.role)
				if err != nil {
					return false, err
				}
			}
		}
	}

	return foundNew, nil
}

func addRequest2Role(tx *sql.Tx, requestID int32, role string) (bool, error) {
	var roleID int32
	found := 0

	pq := dbutl.PQuery(`
		SELECT role_id
		FROM role
		WHERE loweredrole = lower(?)
	`, role)

	err := tx.QueryRow(pq.Query, pq.Args...).Scan(&roleID)
	if err != nil {
		return false, err
	}

	pq = dbutl.PQuery(`
		SELECT CASE WHEN EXISTS (
			SELECT 1
			FROM request_role
			WHERE role_id = ?
			AND request_id = ?
		) THEN 1 ELSE 0 END FROM dual
	`, roleID,
		requestID)

	err = tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return false, err
	}

	if found == 1 {
		return false, nil
	}

	pq = dbutl.PQuery(`
		INSERT INTO request_role (
			role_id,
			request_id
		)
		VALUES (?, ?)
	`, roleID,
		requestID)

	_, err = dbutl.ExecTx(tx, pq)
	if err != nil {
		return false, err
	}

	return true, nil
}
