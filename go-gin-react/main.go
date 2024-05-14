package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Channel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID        int    `json:"id"`
	ChannelID int    `json:"channel_id"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
}

func main() {
	// Get the working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Print the working directory
	fmt.Println("Working directory:", wd)

	// Open the SQLite database file
	db, err := sql.Open("sqlite3", wd+"/database.db")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// Create the channels table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS channels (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
)`)

	if err != nil {
		log.Fatal(err)
	}

	// Create the messages table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY,
    channel_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Gin router
	r := gin.Default()

	r.Use(cors.Default())

	if err != nil {
		log.Fatal(err)
	}

	// creation endpoints
	r.POST("/users", func(c *gin.Context) { createUser(c, db) })
	r.POST("/channels", func(c *gin.Context) { createChannel(c, db) })
	r.POST("/messages", func(c *gin.Context) { createMessage(c, db) })

	// Listing endpoints
	r.GET("/channels", func(c *gin.Context) { listChannels(c, db) })
	r.GET("/messages", func(c *gin.Context) { listMessages(c, db) })

	// Login endpoint
	r.POST("/login", func(c *gin.Context) { login(c, db) })

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

// User creation endpoint
func createUser(c *gin.Context, db *sql.DB) {
	// parse JSON request body into User struct
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert user into database
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?,?)", user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get ID of newly inserted user
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return ID of newly inserted user
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "User created successfully"})
}

// Login endpoint
func login(c *gin.Context, db *sql.DB) {
	// parse JSON request body into User struct
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query database for user
	row := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", user.Username, user.Password)

	// Get ID of user
	var id int
	err := row.Scan(&id)
	if err != nil {
		// Check if user was not found
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		// Return error if other error occurred
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Return ID of user
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "User Login Successfully"})

}

// Channel creation endpoint
func createChannel(c *gin.Context, db *sql.DB) {
	// Parse JSON request body into Channel struct
	var channel Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO channels (name) VALUES (?)")
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	id, err := stmt.Exec(channel.Name)
	if err != nil {
		fmt.Println("Error executing SQL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lastInsertedID, err := id.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last inserted ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return ID of newly inserted channel
	c.JSON(http.StatusOK, gin.H{"id": lastInsertedID, "message": "Channel created successfully"})

}

// Channel listing endpoint
func listChannels(c *gin.Context, db *sql.DB) {
	// Query database for channels
	rows, err := db.Query("SELECT id, name FROM channels")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create slice of channels
	var channels []Channel

	// Iterate over rows
	for rows.Next() {
		// Create new channel
		var channel Channel

		// Scan row into channel
		err := rows.Scan(&channel.ID, &channel.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Append channel to slice
		channels = append(channels, channel)
	}

	// Return slice of channels
	c.JSON(http.StatusOK, channels)
}

// Message creation endpoint
func createMessage(c *gin.Context, db *sql.DB) {
	// pare JSON request body into Message struct
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert message into database
	result, err := db.Exec("INSERT INTO messages (channel_id, user_id, message) VALUES (?, ?, ?)", message.ChannelID, message.UserID, message.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get ID of newly inserted message
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return ID of newly inserted message
	c.JSON(http.StatusOK, gin.H{"id": id})

}

func listMessages(c *gin.Context, db *sql.DB) {
	// Parse channel ID from URL
	channelIDStr := c.Query("channelID")
	if channelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channelID is required"})
		return
	}

	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse optional limit query parameter from URL
	limitParam := c.Query("limit")
	limit := 100 // default limit
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	// Parse last message ID query parameter from URL. This is used to get messages after a certain message.
	lastMessageIDParam := c.Query("lastMessageID")
	lastMessageID := 0 // default lastMessageID
	if lastMessageIDParam != "" {
		lastMessageID, err = strconv.Atoi(lastMessageIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lastMessageID parameter"})
			return
		}
	}

	// Query database for messages
	rows, err := db.Query("SELECT m.id, channel_id, user_id, u.username AS user_name, message FROM messages m LEFT JOIN users u ON u.id = m.user_id WHERE channel_id = ? AND m.id > ? ORDER BY m.id ASC LIMIT ?", channelID, lastMessageID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create slice of messages
	var messages []Message

	// Iterate over rows
	for rows.Next() {
		// Create new message
		var message Message
		var userNameNullString sql.NullString // Using sql.NullString to handle NULL values

		// Scan row into message
		err := rows.Scan(&message.ID, &message.ChannelID, &message.UserID, &userNameNullString, &message.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Assign user_name value to message.UserName
		if userNameNullString.Valid {
			message.UserName = userNameNullString.String
		} else {
			message.UserName = "" // or any other default value you prefer
		}

		// Append message to slice
		messages = append(messages, message)
	}

	// Return slice of messages
	c.JSON(http.StatusOK, messages)
}
