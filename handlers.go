package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type dbRow struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Sex        string `json:"sex"`
	Country    string `json:"country"`
}

type idJson struct {
	Id int `json:"id"`
}

type ageJson struct {
	Age int `json:"age"`
}

type sexJson struct {
	Sex string `json:"gender"`
}

type countryJson struct {
	Country     string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type countriesJson struct {
	Countries []countryJson `json:"country"`
}

func getListOfTasks(c *gin.Context) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s;", "people"))
	if err != nil {
		log.Printf("Error while getting the list of people: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jRows := make([]dbRow, 0)
	for rows.Next() {
		jRows = append(jRows, dbRow{})
		i := len(jRows) - 1
		err := rows.Scan(
			&jRows[i].Id,
			&jRows[i].Name,
			&jRows[i].Surname,
			&jRows[i].Patronymic,
			&jRows[i].Age,
			&jRows[i].Sex,
			&jRows[i].Country)
		if err != nil {
			log.Printf("Error while getting the list of people: %+v\n", err)
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}
	resp, err := json.Marshal(jRows)
	log.Println(jRows)
	log.Println(string(resp))
	if err != nil {
		log.Printf("Error while getting the list of people: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jRowsFiltered := make([]dbRow, 0)
	id, idErr := strconv.Atoi(c.Request.URL.Query().Get("id"))
	name := c.Request.URL.Query().Get("name")
	surname := c.Request.URL.Query().Get("surname")
	patronymic := c.Request.URL.Query().Get("patronymic")
	age, ageErr := strconv.Atoi(c.Request.URL.Query().Get("age"))
	sex := c.Request.URL.Query().Get("sex")
	country := c.Request.URL.Query().Get("country")
	for _, el := range jRows {
		if el.Id != id && idErr == nil {
			continue
		}
		if el.Name != name && name != "" {
			continue
		}
		if el.Surname != surname && surname != "" {
			continue
		}
		if el.Patronymic != patronymic && patronymic != "" {
			continue
		}
		if el.Age != age && ageErr == nil {
			continue
		}
		if el.Sex != sex && sex != "" {
			continue
		}
		if el.Country != country && country != "" {
			continue
		}
		jRowsFiltered = append(jRowsFiltered, el)
	}
	log.Println("Successfully got the list of people\n")
	c.JSON(http.StatusOK, jRowsFiltered)
}

func deleteTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Error while deleting a human from data base: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = %d;", "people", id))
	if err != nil {
		log.Printf("Error while deleting a human from data base: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	log.Println("Successfully deleted the row from the list of people\n")
	c.JSON(http.StatusNoContent, gin.H{})
}

func updateTaskById(c *gin.Context) {
	var row dbRow
	if err := c.BindJSON(&row); err != nil {
		log.Printf("Error while updating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := db.Exec(fmt.Sprintf("UPDATE %s SET name = %s, surname = %s, patronymic = %s, age = %s, sex = %s, country = %s WHERE id = %d;",
		cfg.DbName, row.Id, row.Name, row.Surname, row.Patronymic, row.Age, row.Sex, row.Country, row.Id))
	count, _ := res.RowsAffected()
	if err != nil {
		log.Printf("Error while updating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		createTask(c)
		return
	}
	resp, err := json.Marshal(row)
	if err != nil {
		log.Printf("Error while updating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Successfully updated the row\n")
	c.JSON(http.StatusOK, resp)
}

func createTask(c *gin.Context) {
	var row dbRow
	if err := c.BindJSON(&row); err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jRow, err := json.Marshal(row)
	log.Printf("row: %s\n", jRow)
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	//getting age via open API
	bodyReaderAge := bytes.NewBuffer(jRow)
	reqAge, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.agify.io/?name=%s", row.Name), bodyReaderAge)
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resAge, err := client.Do(reqAge)
	defer resAge.Body.Close()
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	get := make([]byte, 1000)
	_, err = resAge.Body.Read(get)
	var age ageJson
	json.Unmarshal(rmZeroes(get), &age)
	row.Age = age.Age
	log.Printf("age: %s\n", get)

	//getting sex via open API
	bodyReaderSex := bytes.NewBuffer(jRow)
	reqSex, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.genderize.io/?name=%s", row.Name), bodyReaderSex)
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resSex, err := client.Do(reqSex)
	defer resSex.Body.Close()
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	get = make([]byte, 1000)
	_, err = resSex.Body.Read(get)
	var sex sexJson
	json.Unmarshal(rmZeroes(get), &sex)
	row.Sex = sex.Sex
	log.Printf("sex: %s\n", get)

	//getting country via open API
	bodyReaderCountry := bytes.NewBuffer(jRow)
	reqCountry, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.nationalize.io/?name=%s", row.Name), bodyReaderCountry)
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resCountry, err := client.Do(reqCountry)
	defer resCountry.Body.Close()
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	get = make([]byte, 1000)
	_, err = resCountry.Body.Read(get)
	var prbs countriesJson
	json.Unmarshal(rmZeroes(get), &prbs)
	row.Country = getMostProbableCountry(prbs.Countries)
	log.Printf("country: %s\n", get)

	executed, _ := db.Exec(fmt.Sprintf("INSERT INTO %s (name, surname, patronymic, age, sex, country)\nVALUES ('%s', '%s', '%s', %d, '%s', '%s');", "people",
		row.Name, row.Surname, row.Patronymic, row.Age, row.Sex, row.Country))
	id64, _ := executed.LastInsertId()
	row.Id = int(id64)
	if err != nil {
		log.Printf("Error while creating the task: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("Successfully created the row\n")
	c.JSON(http.StatusCreated, row)
}
