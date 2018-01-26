package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// MembershipRole - role utils
type MembershipRole struct {
	sync.RWMutex
	tx       *sql.Tx
	RoleID   int    `sql:"role_id"`
	Rolename string `sql:"role"`
}

var membershipRoleLock sync.RWMutex

// Exists - role exists
func (r *MembershipRole) Exists() (bool, error) {
	found := false

	pq := dbUtils.PQuery(`
	    SELECT CASE WHEN EXISTS (
	        SELECT 1
	          FROM role
	         WHERE loweredrole = lower(?)
	    ) THEN 1 ELSE 0 END
	    FROM dual
	`, r.Rolename)

	err := r.tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

// GetByName - get role by name
func (r *MembershipRole) GetByName(role string) error {
	r.Lock()
	defer r.Unlock()

	pq := dbUtils.PQuery(`
	    SELECT role_id,
	           role
	      FROM role
	     WHERE loweredrole = lower(?)
	`, role)

	err := dbUtils.RunQueryTx(r.tx, pq, r)

	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("Role \"%s\" not found", role)
	case err != nil:
		return err
	}

	return nil
}

// GetByID - get role by ID
func (r *MembershipRole) GetByID(roleID int) error {
	r.Lock()
	defer r.Unlock()

	pq := dbUtils.PQuery(`
	    SELECT role_id,
	           role
	      FROM role
	     WHERE role_id = ?
	`, roleID)

	err := dbUtils.RunQueryTx(r.tx, pq, r)

	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("role not found")
	case err != nil:
		return err
	}

	return nil
}

func (r *MembershipRole) testSave() error {
	if len(r.Rolename) == 0 {
		return fmt.Errorf("unknown role \"%s\"", r.Rolename)
	}

	var found bool

	pq := dbUtils.PQuery(`
	    SELECT CASE WHEN EXISTS (
	        SELECT 1
	          FROM role
	         WHERE loweredrole = lower(?)
	           AND role_id <> ?
	    ) THEN 1 ELSE 0 END
	    FROM dual
	`, r.Rolename,
		r.RoleID)

	err := r.tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return err
	}

	if found {
		return fmt.Errorf("duplicate role \"%s\"", r.Rolename)
	}

	return nil
}

// Save - save role details
func (r *MembershipRole) Save() error {
	membershipRoleLock.Lock()
	defer membershipRoleLock.Unlock()

	err := r.testSave()
	if err != nil {
		return err
	}

	if r.RoleID <= 0 {
		pq := dbUtils.PQuery(`
		    INSERT INTO role (role, loweredrole) VALUES (?, lower(?))
		`, r.Rolename, r.Rolename)

		_, err = dbUtils.ExecTx(r.tx, pq)
		if err != nil {
			return err
		}

		pq = dbUtils.PQuery(`
			SELECT role_id FROM role WHERE loweredrole = lower(?)
		`, r.Rolename)

		err = r.tx.QueryRow(pq.Query, pq.Args...).Scan(&r.RoleID)

		switch {
		case err == sql.ErrNoRows:
			r.RoleID = -1
		case err != nil:
			return err
		}

		audit.Log(nil, "add-role", "Add new role.", "new", r)
	} else {
		old := MembershipRole{tx: r.tx}
		err = old.GetByID(r.RoleID)
		if err != nil {
			return err
		}

		pq := dbUtils.PQuery(`
			UPDATE role
			   SET role = ?,
		           loweredrole = lower(?)
		     WHERE role_id = ?
		`, r.Rolename,
			r.Rolename,
			r.RoleID)

		_, err = dbUtils.ExecTx(r.tx, pq)
		if err != nil {
			return err
		}

		audit.Log(nil, "update-role", "Update role.", "old", &old, "new", r)
	}

	return nil
}

// HasMember - role has member
func (r *MembershipRole) HasMember(user string) (bool, error) {
	r.RLock()
	defer r.RUnlock()

	found := false

	dt := time.Now().UTC()

	pq := dbUtils.PQuery(`
	    SELECT CASE WHEN EXISTS (
	        SELECT 1
	          FROM user_role ur
	          JOIN "user" u ON (ur.user_id = u.user_id)
	         WHERE u.loweredusername =  lower(?)
	           AND ur.role_id        =  ?
	           AND ur.valid_from     <= ?
	           AND (ur.valid_until is null OR ur.valid_until > ?)
	    ) THEN 1 ELSE 0 END
	    FROM dual
	`, user,
		r.RoleID,
		dt,
		dt)

	err := r.tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

// HasMemberID - has member ID
func (r *MembershipRole) HasMemberID(userID int) (bool, error) {
	r.RLock()
	defer r.RUnlock()

	found := false

	dt := time.Now().UTC()

	pq := dbUtils.PQuery(`
	    SELECT CASE WHEN EXISTS (
	        SELECT 1
	          FROM user_role ur
	         WHERE ur.user_id =  ?
	           AND ur.role_id =  ?
	           AND ur.valid_from <= ?
	           AND (ur.valid_until is null OR ur.valid_until > ?)
	    ) THEN 1 ELSE 0 END
	    FROM dual
	`, userID,
		r.RoleID,
		dt,
		dt)

	err := r.tx.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

// IsUserInRole - Is user in role
func IsUserInRole(user string, role string) (bool, error) {
	found := false

	dt := time.Now().UTC()

	pq := dbUtils.PQuery(`
	    SELECT CASE WHEN EXISTS (
	        SELECT 1
	          FROM user_role ur
	          JOIN "user" u ON (ur.user_id = u.user_id)
	          JOIN role r ON (ur.role_id = r.role_id)
	         WHERE u.loweredusername =  lower(?)
	           AND lower(r.role)     =  lower(?)
	           AND ur.valid_from     <= ?
	           AND (ur.valid_until is null OR ur.valid_until > ?)
	    ) THEN 1 ELSE 0 END
	    FROM dual
	`, user,
		role,
		dt,
		dt)

	err := db.QueryRow(pq.Query, pq.Args...).Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}
