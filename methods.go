package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type pantryValue struct {
	UnitID string `json:"unit_id" binding:"required"`
	Values []int  `json:"items" binding:"required"`
}

type createDeviceT struct {
	ApiKey string `json:"api_key"`
}

func GenerateRandomStringURLSafe(n int) (string, error) {
	/* Creates a length n string of random characters that are
	   allowed in a URL address */

	b := make([]byte, n)
	_, err := rand.Read(b)
	return base64.URLEncoding.EncodeToString(b), err
}

func createDeviceEntry(db *sql.DB) gin.HandlerFunc {
	/* Create a new database entry for an account.
	   Returns API key for user.
	*/
	return func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
			"items (" +
			"api_key VARCHAR(255) UNIQUE NOT NULL, " +
			"items TEXT )"); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
			return
		}

		apikey, err := GenerateRandomStringURLSafe(32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Server error.")
			return
		}

		if _, err := db.Exec("INSERT INTO items(api_key) VALUES ($1);", apikey); err != nil {
			c.JSON(http.StatusBadRequest, "Error")
			return
		}
		c.JSON(http.StatusAccepted, createDeviceT{apikey})
	}
}

func postPantry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.Request.Header.Get("Authorization")
		api_key_parts := strings.Split(header, "Bearer ")
		if api_key_parts == nil {
			c.JSON(http.StatusInternalServerError, "No bearer token found.")
			return
		}

		api_key := api_key_parts[1]

		var pantryStatus pantryValue

		if err := c.BindJSON(&pantryStatus); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error %s", err))
			return
		}

		row := db.QueryRow("SELECT exists (SELECT 1 FROM items WHERE api_key = $1 LIMIT 1);", api_key)

		var exists *bool

		if err := row.Scan(&exists); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error fetching data: %s", err))
			return
		}

		if !*exists {
			c.JSON(http.StatusNotFound, "Not found")
			return
		}

		data := fmt.Sprint(pantryStatus.Values)
		data = strings.Replace(data, " ", ",", -1)

		old_data, err := getData(db, &api_key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Something failed.")
		}

		old_data[pantryStatus.UnitID] = pantryStatus.Values

		err = saveData(db, &api_key, &old_data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
		}

		c.JSON(http.StatusAccepted, fmt.Sprintf("Data added: %s", data))
	}
}

func saveData(db *sql.DB, api_key *string, data *map[string][]int) error {
	/* Save map to database in JSON format. */

	json_str, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO items(api_key, items) VALUES ($1, $2)"+
		"ON CONFLICT (api_key) DO UPDATE SET items = $2;", api_key, json_str)

	return err
}

func getData(db *sql.DB, api_key *string) (map[string][]int, error) {
	/* Get value of items field in database. */
	row := db.QueryRow("SELECT items FROM items WHERE api_key = $1;", api_key)
	datamap := make(map[string][]int)

	var data sql.NullString
	err := row.Scan(&data)
	if err != nil {
		return datamap, err
	}

	err = json.Unmarshal([]byte(data.String), &datamap)

	return datamap, err
}

func getPantry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.Request.Header.Get("Authorization")
		api_key_parts := strings.Split(header, "Bearer ")
		if api_key_parts == nil {
			c.JSON(http.StatusInternalServerError, "No bearer token found.")
			return
		}

		api_key := api_key_parts[1]

		data, err := getData(db, &api_key)

		switch err {
		case sql.ErrNoRows:
			c.IndentedJSON(http.StatusNotFound, "No data found!")
		case nil:
		default:
			c.IndentedJSON(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", err))
			return
		}
		c.JSON(http.StatusAccepted, data)
	}
}
