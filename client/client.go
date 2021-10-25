package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"app/models"
	"app/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var c proto.UserServiceClient

func AllUser(ctx *gin.Context) {

	req := proto.GetAllUserRequest{
		TableName: "users",
	}

	res, err := c.AllUser(context.Background(), &req)

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, res)
}

func getUser(ctx *gin.Context) {

	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Print(id)
	}

	req := proto.GetUserRequest{
		Id: int32(id),
	}

	res, err := c.GetUser(context.Background(), &req)

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(200, res)
}

func postUser(ctx *gin.Context) {

	pUser := models.PostUser{}

	ctx.BindJSON(&pUser)

	req := proto.PostUserRequest{
		User: &proto.User{
			FirstName: pUser.FirstName,
			LastName:  pUser.LastName,
			Username:  pUser.Username,
			Age:       pUser.Age,
		},
	}

	res, err := c.PostUser(context.Background(), &req)

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(201, res)
}

func delUser(ctx *gin.Context) {

	id := models.Id{}

	ctx.BindJSON(&id)

	req := proto.DeleteUserRequest{
		Id: id.Id,
	}

	res, err := c.DeleteUser(context.Background(), &req)

	if err != nil {
		log.Print(err)
	}

	ctx.JSON(202, res)
}

func updateUser(ctx *gin.Context) {

	user := models.User{}

	ctx.BindJSON(&user)

	req := proto.UpdateUserRequest{
		Id: user.Id,
		User: &proto.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Age:       user.Age,
		},
	}

	res, err := c.UpdateUser(context.Background(), &req)

	if err != nil {
		log.Print(err)
	}

	ctx.JSON(202, res)
}

func main() {

	fmt.Println("Connection Client RPC...")

	httpRouter := gin.Default()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c = proto.NewUserServiceClient(conn)

	httpRouter.GET("/getUser", AllUser)
	httpRouter.GET("/getUser/:id", getUser)
	httpRouter.POST("/postUser", postUser)
	httpRouter.DELETE("/delUser", delUser)
	httpRouter.PUT("/putUser", updateUser)

	httpRouter.Run()
}
