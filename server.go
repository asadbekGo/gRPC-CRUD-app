package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"app/database"
	"app/models"
	"app/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedUserServiceServer
}

func (*server) AllUser(c context.Context, req *proto.GetAllUserRequest) (*proto.GetAllUserResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	users := []models.User{}

	rows := db.Table(req.TableName).Find(&users)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	fmt.Println(users)

	var result []*proto.User

	for _, data := range users {

		res := &proto.User{
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Username:  data.Username,
			Age:       data.Age,
		}

		result = append(result, res)
	}

	return &proto.GetAllUserResponse{
		User: result,
	}, nil
}

func (*server) GetUser(c context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Print(err)
	}

	user := models.User{}

	rows := db.Table("users").
		Where("id = ?", req.GetId()).
		Find(&user)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	res := &proto.GetUserResponse{
		Id: user.Id,
		User: &proto.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Age:       user.Age,
		},
	}

	return res, nil
}

func (*server) PostUser(c context.Context, req *proto.PostUserRequest) (*proto.PostUserResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Print(err)
	}

	var id int32

	rows := db.Table("users").
		Create(&req.User).
		Select("id").
		Where("username = ?", req.User.Username).
		Find(&id)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	res := &proto.PostUserResponse{
		Id: id,
		User: &proto.User{
			FirstName: req.User.FirstName,
			LastName:  req.User.LastName,
			Username:  req.User.Username,
			Age:       req.User.Age,
		},
	}

	return res, nil
}

func (*server) DeleteUser(c context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	user := models.PostUser{}
	fmt.Println(req.Id)

	rows := db.Table("users").
		Select("first_name, last_name, username, age").
		Where("id = ?", req.Id).
		Find(&user)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	del := db.Table("users").
		Delete(&models.User{}, req.Id)

	if database.IsNotFound(del) {
		log.Fatal(del.Error)
	}

	res := &proto.DeleteUserResponse{
		User: &proto.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Age:       user.Age,
		}}

	return res, nil
}

func (*server) UpdateUser(c context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	db, err := database.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	userTemp := models.PostUser{}

	findUser := db.Table("users").
		Select("first_name, last_name, username, age").
		Where("id = ?", req.Id).
		Find(&userTemp)

	if database.IsNotFound(findUser) {
		log.Fatal(findUser.Error)
	}

	if req.User.FirstName == "" {
		req.User.FirstName = userTemp.FirstName
	}

	if req.User.LastName == "" {
		req.User.LastName = userTemp.LastName
	}

	if req.User.Username == "" {
		req.User.Username = userTemp.Username
	}

	if req.User.Age == 0 {
		req.User.Age = userTemp.Age
	}

	rows := db.Table("users").
		Where("id = ?", req.Id).
		Updates(
			map[string]interface{}{
				"first_name": req.User.FirstName,
				"last_name":  req.User.LastName,
				"username":   req.User.Username,
				"age":        req.User.Age,
			},
		)

	if database.IsNotFound(rows) {
		log.Fatal(rows.Error)
	}

	return &proto.UpdateUserResponse{}, nil
}

func main() {

	fmt.Println("Connection Server RPC...")

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	proto.RegisterUserServiceServer(s, &server{})

	err = s.Serve(lis)

	if err != nil {
		log.Fatal(err)
	}
}
