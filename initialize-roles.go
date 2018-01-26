package main

import "database/sql"

type userRole struct {
	role string
}

func addRoles(tx *sql.Tx) (bool, error) {
	roles := []userRole{
		{"Administrator"},
		{"Member"},
		{"All"},
	}

	foundNew := false

	for _, r := range roles {
		mrole := MembershipRole{tx: tx}
		mrole.Rolename = r.role

		found, err := mrole.Exists()
		if err != nil {
			return false, err
		}

		if found {
			continue
		}

		foundNew = true

		err = mrole.Save()
		if err != nil {
			return false, err
		}
	}

	return foundNew, nil
}
