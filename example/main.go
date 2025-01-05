package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/danielMensah/onetable-go"
)

type User struct {
	UserId string `dynamodbav:"userId" onetable:"pk"`
	Name   string `dynamodbav:"userName"`
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	o := onetable.New(dynamodb.NewFromConfig(cfg), "onetable")

	// Create item
	output, _ := o.CreateItem(ctx, User{UserId: "1234", Name: "abigail"})
	fmt.Printf("output: %v", output)

	// Get item
	u := &User{}
	_ = o.Find(ctx, u, "userId = $1", "1234")
	fmt.Printf("user: %v", u)

	// Get items
	uu := &[]User{}
	_ = o.Find(ctx, uu, "name = $1", "abigail")
	fmt.Printf("users: %v", uu)

	// Update item
	_ = o.UpdateItem(ctx, User{UserId: "1234"}, map[string]interface{}{"name": "abigail"})
	fmt.Printf("successfully updated user")
}
