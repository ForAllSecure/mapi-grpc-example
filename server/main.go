// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	pb "github.com/forallsecure/mapi-grpc-example/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
	db   = createDatabase()
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedUserServiceServer
}

func createDatabase() *sql.DB {
	const initDb = `CREATE TABLE IF NOT EXISTS "users" (
		"uid" INTEGER PRIMARY KEY AUTOINCREMENT,
		"username" VARCHAR(64) NULL,
		"email" VARCHAR(256) NULL,
		"created" DATE NULL
	);`

	db, err := sql.Open("sqlite3", "./data.db")
	checkErr(err)

	// insert
	stmt, err := db.Prepare(initDb)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	return db
}

//
// Add a new user to the database
//
func (s *server) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.UserResult, error) {
	currentTime := time.Now()
	created := currentTime.Format("2006-01-02")

	// !! Should use a prepared statement here!
	insertStatement := fmt.Sprintf("INSERT INTO users(username, email, created) values('%s', '%s', '%s');",
		in.Username,
		in.Email,
		created)

	res, err := db.Exec(insertStatement)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to add user - %s", err.Error()))
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to add user - %s", err.Error()))
	}

	return &pb.UserResult{Id: id, Username: in.Username, Email: in.Email, Created: created}, nil
}

//
// Return a list of all users that match the given filter
//
func (s *server) GetUsers(ctx context.Context, in *pb.GetUsersRequest) (*pb.UsersResult, error) {
	// !! Should use a prepared statement here!
	selectStatement := fmt.Sprintf("SELECT * FROM users WHERE username like '%%%s%%'", in.GetFilter())
	rows, err := db.Query(selectStatement)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to read users - %s", err.Error()))
	}
	defer rows.Close()

	var users []*pb.UserResult

	for rows.Next() {
		var user pb.UserResult
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Created); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to read users - %s", err.Error()))
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to read users - %s", err.Error()))
	}
	return &pb.UsersResult{Users: users}, nil
}

//
// Delete a User by ID
//
func (s *server) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResult, error) {
	// !! Should use a prepared statement here!
	deleteStatement := fmt.Sprintf("DELETE FROM users WHERE uid=%d;", in.GetId())

	res, err := db.Exec(deleteStatement)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to delete user - %s", err.Error()))
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to delete user - %s", err.Error()))
	}
	if count == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Could not find user with ID - %d", in.GetId()))
	}

	return &pb.DeleteUserResult{Count: count}, nil
}

//
// Check the list or reserved names to see if there is a match
//
func (s *server) CheckReservedName(ctx context.Context, in *pb.CheckReservedNameRequest) (*pb.CheckReservedNameResult, error) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cat reserved-names.txt | grep %s || true ", in.GetName()))

	out, err := cmd.Output()
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to check reserved name - %s", err.Error()))
	}

	return &pb.CheckReservedNameResult{Reserved: fmt.Sprintf("%s", out)}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
