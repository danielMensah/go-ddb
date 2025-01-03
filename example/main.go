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

	CreateUser(o, ctx)
	GetUser(o, ctx)
	GetUsers(o, ctx)
}

func CreateUser(o onetable.Client, ctx context.Context) {
	output, _ := o.CreateItem(ctx, User{UserId: "1234", Name: "abigail"})
	fmt.Printf("output: %v", output)
}

func GetUser(o onetable.Client, ctx context.Context) {
	u := &User{}
	_ = o.Find(ctx, u, "userId = $1", "1234")

	fmt.Printf("user: %v", u)
}

func GetUsers(o onetable.Client, ctx context.Context) {
	u := &[]User{}
	_ = o.Find(ctx, u, "name = $1, age = $2", "abigail", 25)

	fmt.Printf("users: %v", u)
}
