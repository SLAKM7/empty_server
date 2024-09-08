package book

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc/metadata"

	pb "grpc-gateway-demo/pkg/proto/book"
)

type BookService struct {
	pb.UnimplementedBookServiceServer
}

func NewBookService() *BookService {
	l := localStorage{
		Count: 0,
		DB:    make(map[int32]*pb.Book),
	}
	db = &l
	return &BookService{}
}

var db *localStorage

type localStorage struct {
	Count int32
	DB    map[int32]*pb.Book
	mux   sync.Mutex
}

func (l *localStorage) getId() int32 {
	l.Count = l.Count + 1
	return l.Count
}

func (l *localStorage) Store(d *pb.Book) error {
	if d == nil {
		return errors.New("data is nil")
	}

	if d.Id <= 0 {
		return errors.New("illegal id")
	}
	l.DB[d.Id] = d
	return nil
}

func (l *localStorage) Load(id int32) (*pb.Book, error) {
	if id <= 0 {
		return nil, errors.New("illegal id")
	}
	book := l.DB[id]
	return book, nil
}

func (b *BookService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	resp := &pb.CreateBookResponse{}
	db.mux.Lock()
	defer db.mux.Unlock()
	id := db.getId()
	book := pb.Book{
		Name: req.GetName(),
		Id:   id,
	}

	err := db.Store(&book)
	if err != nil {
		return resp, err
	}

	fmt.Printf("user %s create a book", getUserId(ctx))
	resp.Data = &book

	return resp, nil
}

func (b *BookService) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	resp := &pb.GetBookResponse{}
	db.mux.Lock()
	defer db.mux.Unlock()

	book, err := db.Load(req.GetId())
	if err != nil {
		return resp, err
	}

	resp.Data = book
	return resp, nil
}

func getUserId(ctx context.Context) string {
	userId := ""
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md["x-user-id"]) == 0 {
		return userId
	}

	userId = strings.Join(md["x-user-id"], ",")

	return userId
}
