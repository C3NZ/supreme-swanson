// Getting started with echo

package main

import (
    "net/http"
    "io/ioutil"
    "log"

    "github.com/labstack/echo"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// The swanson quote to be stored
type SwansonQuote struct {
    gorm.Model
    Quote string `json:"quote"`
}

// Init the database and return it 
func initDatabase() *gorm.DB {
    db, err := gorm.Open("sqlite3", "test.db")

    if err != nil {
        panic("Failed to init databse")
    }

    return db
}

func main() {
    e := echo.New()
    
    // Init the database and migrate the struct as a schema 
    db := initDatabase()
    db.AutoMigrate(&SwansonQuote{})


    // Get all Ron Swanson quotes
    e.GET("/swanson", func(c echo.Context) error {
        // Send a get request to our swanson API 
        // Grabs both the response and err (if any)
        resp, err := http.Get("http://ron-swanson-quotes.herokuapp.com/v2/quotes");

        // Check if an error occurred
        if err != nil {
            log.Fatal("Request could not be made")
        }
       
        // Read the response data into body as a byte slice
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal("Could not parse the body of the response")
        }
        //
        quote := string(body)
        db.Create(&SwansonQuote{Quote: quote})

        // Return the quote to the user 
        return c.String(http.StatusCreated, quote)
    })

    e.GET("/swansonquotes", func(c echo.Context) error {
        quotes := []SwansonQuote{}
        db.Find(&quotes) 
        return c.JSON(http.StatusCreated, quotes)
    })
    e.GET("/bitcoin", func(c echo.Context) error {
        resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")

        if err != nil {
            log.Fatal("Request could not be made")
        }

        body, err := ioutil.ReadAll(resp.Body)

        return c.JSONPretty(http.StatusCreated, body, "   ")
    })

    e.Logger.Fatal(e.Start(":1323"))

}
