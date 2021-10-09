package main

import (
	"context"
	"fmt"

	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id       string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"Email,omitempty" bson:"Email,omitempty"`
	Password string `json:"Password,omitempty" bson:"Password,omitempty"`
}

type Post struct {
	Userid    string    `json:"Userid,omitempty" bson:"Userid,omitempty"`
	Id        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	caption   string    `json:"caption,omitempty" bson:"caption,omitempty"`
	url       string    `json:"url,omitempty" bson:"url,omitempty"`
	timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

//function to produce random strings of different lengths for user name email and password
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("go_search")

	app := fiber.New()

	app.Use(cors.New())

	var savuserid = ""

	//creating an user with all its fields random string
	app.Post("/users", func(c *fiber.Ctx) error {
		collection := db.Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		var id = RandomString(15)
		var name = RandomString(9)
		var email = RandomString(10)
		var password = RandomString(5)
		collection.InsertOne(ctx, User{
			Id:       id,
			Name:     name,
			Email:    email,
			Password: password,
		})

		// The id of the user is now stored in savuserid variable
		savuserid += id
		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	// the user with id saved in savuserid variable will be able to post further having userid same to all his/her post
	app.Post("/posts", func(c *fiber.Ctx) error {
		collection := db.Collection("posts")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		collection.InsertOne(ctx, Post{
			Userid:    savuserid,
			Id:        RandomString(15),
			caption:   RandomString(20),
			url:       fmt.Sprintf("http://lorempixel.com/200/200?%s", rand.Intn(100)), //this site generetes random images provided we give random integer at the ned of url
			timestamp: time.Now(),
		})

		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	//get request to fetch specific user with given userid
	var userid = "kuggasbjka"
	app.Get("/users/%s,userid", func(c *fiber.Ctx) error {
		collection := db.Collection("Posts")
		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

		var users []User

		//Querying Documents from users Collection with a Filter where id is same as we have passed through get request
		cursor, _ := collection.Find(ctx, bson.M{"Id": userid})

		for cursor.Next(ctx) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}
		fmt.Println(c.JSON(users))
		return c.JSON(users)
	})

	//get request to fetch specific post with given postid
	var postid = "aeqynsdfgsd845"
	app.Get("/users/%s,postid", func(c *fiber.Ctx) error {
		collection := db.Collection("posts")
		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

		var posts []Post

		//Querying Documents from posts Collection with a Filter where id is same as we have passed through get request
		cursor, _ := collection.Find(ctx, bson.M{"Id": postid})

		for cursor.Next(ctx) {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
		}
		fmt.Println(c.JSON(posts))
		return c.JSON(posts)
	})

	//get request to fetch all posts with given user id
	var useridforposts = "aeqynsdfgsd845"
	app.Get("/users/%s,useridforposts", func(c *fiber.Ctx) error {
		collection := db.Collection("posts")
		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

		var posts []Post

		//Querying Documents from posts Collection with a Filter where id is same as we have passed through get request
		cursor, _ := collection.Find(ctx, bson.M{"Userid": useridforposts})

		for cursor.Next(ctx) {
			var post Post
			cursor.Decode(&post)
			posts = append(posts, post)
		}

		fmt.Println(c.JSON(posts))
		return c.JSON(posts)
	})

	app.Listen(":8000")
}
