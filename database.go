package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Database *sql.DB
)

func DatabaseInit(host string, port int, user, password, name string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, name))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	Database = db
	return nil
}

func DatabaseVideoExists(id int) (bool, error) {
	res, err := Database.Query("SELECT EXISTS(SELECT id FROM dv_video WHERE id = ?)", id)
	if err != nil {
		return false, err
	}

	defer res.Close()

	if res.Next() {
		var e bool

		err = res.Scan(&e)
		if err != nil {
			return false, err
		}

		return e, nil
	}

	return false, fmt.Errorf("Row not found")
}

func DatabaseInsertConversion(c *Conversion) error {
	res, err := Database.Exec("INSERT INTO dv_conversion (video_id, format_id, status_id) VALUES (?, ?, ?)", c.VideoID, c.FormatID, c.StatusID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = int(id)

	return nil
}

func DatabaseGetConversion(id int) (Conversion, error) {
	var c Conversion

	res, err := Database.Query("SELECT * FROM dv_conversion WHERE id = ? LIMIT 1", id)
	if err != nil {
		return c, err
	}

	defer res.Close()

	if res.Next() {
		err = res.Scan(&c.ID, &c.VideoID, &c.FormatID, &c.StatusID)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

func DatabaseGetVideoConversions(vid int) ([]Conversion, error) {
	res, err := Database.Query("SELECT * FROM dv_conversion WHERE video_id = ?", vid)
	if err != nil {
		return nil, err
	}

	cs := make([]Conversion, 0)
	for res.Next() {
		var c Conversion

		err = res.Scan(&c.ID, &c.VideoID, &c.FormatID, &c.StatusID)
		if err != nil {
			return nil, err
		}

		cs = append(cs, c)
	}

	return cs, nil
}

func DatabaseUpdateConversion(c *Conversion) error {
	_, err := Database.Exec("REPLACE INTO dv_conversion VALUES (?, ?, ?, ?)", c.ID, c.VideoID, c.FormatID, c.StatusID)
	if err != nil {
		return err
	}

	return nil
}
