# grpc-gateway-demo


## Protocol Buffer 官方文档
[官方文档](https://developers.google.com/protocol-buffers)

[语法指南](https://developers.google.com/protocol-buffers/docs/proto3#nested)

[protoc下载地址](https://github.com/protocolbuffers/protobuf/releases/tag/v21.12)
## REST API
[谷歌文档](https://cloud.google.com/apis/design/resources)

## grpc-gateway 官方文档
[官方文档](https://grpc-ecosystem.github.io/grpc-gateway/)

[github readme](https://github.com/grpc-ecosystem/grpc-gateway#readme)

## Buf 文档
[官方文档](https://buf.build/)

## swagger-ui
[安装文档](https://github.com/swagger-api/swagger-ui/blob/master/docs/usage/installation.md)


## protoc-gen-validate
[github](https://github.com/bufbuild/protoc-gen-validate)




# Protocol Buffers的基本使用

## 什么是Protocol buffers

Protocol Buffer是一个由 Google 开发的协议，允许结构化数据的序列化和反序列化协议缓冲。谷歌开发它的目的是提供一种比 XML 更好的方式来使系统进行通信。因此，他们致力于使其比 XML 更简单、更小、更快、更易于维护。这个协议甚至超过了 JSON，具有更好的性能、更好的可维护性和更小的体积。

## 使用Protocol Buffers的好处

Protocol buffer 可以支持多语言，以及跨平台

* 解析快
* 传输速度快，体积小
* 支持多语言

支持的语言包括

* [C++](https://developers.google.com/protocol-buffers/docs/reference/cpp-generated#invocation)
* [C#](https://developers.google.com/protocol-buffers/docs/reference/csharp-generated#invocation)
* [Java](https://developers.google.com/protocol-buffers/docs/reference/java-generated#invocation)
* [Kotlin](https://developers.google.com/protocol-buffers/docs/reference/kotlin-generated#invocation)
* [Objective-C](https://developers.google.com/protocol-buffers/docs/reference/objective-c-generated#invocation)
* [PHP](https://developers.google.com/protocol-buffers/docs/reference/php-generated#invocation)
* [Python](https://developers.google.com/protocol-buffers/docs/reference/python-generated#invocation)
* [Ruby](https://developers.google.com/protocol-buffers/docs/reference/ruby-generated#invocation)
* [Go](https://github.com/protocolbuffers/protobuf-go)
* [Dart](https://github.com/google/protobuf.dart)

## 下载安装Protocol Buffers 编译器

通过[github ](https://github.com/protocolbuffers/protobuf/releases) release页面，下载protoc

![image-20230115110519332](https://img.zhouwanderder.xyz/uPic/image-20230115110519332.png)

打开release页面，我们可以直接找到protoc的安装包，将其下载并且配置即可

安装 Go protocol buffers 生成插件

```go   
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

编译产生的插件`protoc-gen-go`会安装在`$GOBIN`下，默认就是`$GOPATH/bin`路径下

## Protocol Buffer基本语法

首先我们需要创建一个proto文件，在proto文件中，我们需要通过`syntax`关键字来指定我们需要用到的protocol buffers的版本。除此之外，我们还需要添加`option go_package`来指定生成的包地址，`package`字段使用于proto之间的依赖指定的路径，同时proto文件中的每一行都需要使用`;`结尾

### 定义消息类型

我们这里定义一个Person类型的消息

```protobuf
// 指定proto语言版本
syntax = "proto3";

// 生成go文件的包路径
option go_package = "/pb";

// proto文件包的路径
package grpc.gateway.demo.proto.examplepb;


message Person {
    int32 id = 1;
    string name = 2;
    string email = 3;
    float salary = 4;
    bool sex = 5;
}
```

通过protoc生成go文件

`protoc -I=. --go_out=. ./proto/examplepb/person.proto `

通过上面编写的简单proto文件，我们可以发现，proto文件中定义message与我们创建一个go中的结构体类似，包括内部对元素的类型定义，我们也是非常熟悉的。但是我们也还有一些出去基本类型之外的复杂类型，例如:枚举、map、数组等，我们需要如何去定义

#### repeated

repeated关键字的作用是用来定义数组，使用方式是`repeated 数组类型 属性名称 =  字段编号;`

```protobuf
message Person {
	repeated string name = 1;
}
```

#### map

map类型的定义方式是`map <键类型，值类型> 属性名称 = 字段编号;` ，这里需要注意对于map的键类型，只能定义为基本数据类型，但是值的类型可以是任何支持的类型

```protobuf 
message Person {
	map <string, Pet> pets =1;
}

message Pet {
	string name = 1;
}
```

#### enum

对于枚举的定义，我们需要用到enum关键字。

```protobuf 
enum PhoneType {
    PHONE_TYPE_UNSPECIFIED = 0;
    PHONE_TYPE_WORK = 1;
    PHONE_TYPE_HOME = 2;
}
```

在枚举的定义中，我们需要指定一个零值`    PHONE_TYPE_UNSPECIFIED = 0;`

## 序列化与反序列化

```go   
package proto

import (
   "fmt"
   "log"
   "os"
   "testing"

   "google.golang.org/protobuf/proto"

   "grpc-gateway-demo/pb"
)

var fName = "person.txt"

func TestPerson(t *testing.T) {
   person := &pb.Person{}
   phone := &pb.Phone{}
   pet := &pb.Pet{
      Name: "Leo",
      Age:  1,
   }

   phone.PhoneNumber = "1234567789"
   phone.PhoneType = pb.PhoneType_PHONE_TYPE_WORK
   // phone.PhoneType = 1

   person.Name = "Ethan"
   person.Email = "xxxx@yeah.net"
   person.Phone = phone
   person.Pets = map[string]*pb.Pet{}
   person.Pets[pet.Name] = pet

   b, err := proto.Marshal(person)
   if err != nil {
      log.Fatalf("proto marshal failed: %v", err)
      return
   }

   f, err := os.Create(fName)
   defer f.Close()
   if err != nil {
      log.Fatalf("os create failed: %v", err)
      return
   }

   _, err = f.Write(b)
   if err != nil {
      log.Fatalf("write file failed: %v", err)
      return
   }

}

func TestReadPerson(t *testing.T) {
   b, err := os.ReadFile(fName)
   if err != nil {
      log.Fatalf("read file failed: %v", err)
      return
   }
   p := &pb.Person{}
   err = proto.Unmarshal(b, p)
   if err != nil {
      log.Fatalf("proto unmarshal failed: %v", err)
      return
   }

   fmt.Println(p)
}
```


#  grpc-gateway的基本使用

## 简介

grpc-gateway是protoc的以一个插件，它会读取proto文件中的grpc服务的定义并将其生成RestFul JSON API,并将其反向代理到我们的grpc服务中

<img src="https://grpc-ecosystem.github.io/grpc-gateway/assets/images/architecture_introduction_diagram.svg" style="zoom:50%;" />

grpc-gatway是在grpc上做的一个拓展。但是grpc并不能很好的支持客户端，以及传统的RESTful API。因此grpc-gateway诞生了，该项目可以为我们的grpc服务提供HTTP+JSON接口。

## 插件安装

首先，我们在项目中去创建一个tools的文件，然后运行`go mod tidy`

```go
package tools

import (
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

最后我们运行`go install`将其安装在我们的`$GOBIN`目录下

```sh
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### grpc 插件安装

[grpc](https://grpc.io/docs/languages/go/quickstart/)文档

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@book.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@book.2
```

## 使用

### 一、定义proto文件

```protobuf
syntax = "proto3";

option go_package = "/pb;";

package grpc.gateway.demo.proto.examplepb;

import "google/api/annotations.proto";

// 导入api注释依赖

// 定义服务
service BookService {
    // 创建书籍
    rpc CreateBook (CreateBookRequest) returns (CreateBookResponse) {
        option (google.api.http) = {
            post: "/book/books"
            body: "*"
        };
    };

    // 获取书籍
    rpc GetBook (GetBookRequest) returns (GetBookResponse) {
        option (google.api.http) = {
            get: "/book/books/{id}"
        };
    }


//    rpc UpdateBook (UpdateBookRequest) returns (UpdateBookResponse) {
//
//    };
//
//
//    rpc DeleteBook (DeleteBookRequest) returns (DeleteBookResponse) {
//
//    };
//
//    rpc ListBook (ListBookRequest) returns (ListBookResponse) {};


}

message Book {
    // int64 会自动转换为string 进行返回
    int32 id = 1;
    string name = 2;

}

// 定义接收参数
message CreateBookRequest {
    string name = 1;

}


message CreateBookResponse {
    string code = 1;
    string message = 2;
    Book data = 3;
}


message GetBookRequest {
    int32 id = 1;
}


message GetBookResponse {
    string code = 1;
    string message = 2;
    Book data = 3;
}
```

我们在使用grpc-gateway这个框架的时候，需要使用proto文件，将我们的grpc服务进行http+json的一个暴露，以此来对外达到一个提供http+json服务的接口

### 二、生成go文件

在编写好proto文件之后，我们需要使用插件将proto文件生成对应的go文件。

由于我们依赖了google的proto文件，所以我们在使用protoc生成go文件的时候，需要将依赖的proto文件复制到我们的项目中，依赖的proto文件仓库[google/api](https://github.com/googleapis/googleapis/tree/master/google/api)

```te
google/api/annotations.proto
google/api/field_behavior.proto
google/api/http.proto
google/api/httpbody.proto

```

我们还需要依赖的[google/protobuf](https://github.com/protocolbuffers/protobuf/tree/main/src/google/protobuf/compiler)

```tex
google/protobuf/descriptor.proto
```



```she
protoc -I ./proto --grpc-gateway_out=../gen \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --go_out=../gen \
    --go_opt paths=source_relative \
    --go-grpc_out=../gen \
    --go-grpc_opt paths=source_relative \
    ./proto/examplepb/book.proto
```

将需要依赖的文件放入项目中，我们就可以使用protoc生成go文件，这里需要用到三个插件

* go
* grpc
* grpc-gateway

### 三、启动服务

#### 实现BookService

```go
package service

import (
	"context"
	"errors"
	"sync"

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

```

#### 启动gRPC

```go
package server

import (
   "fmt"
   "io"
   "net"
   "os"

   "google.golang.org/grpc"
   "google.golang.org/grpc/grpclog"

   "grpc-gateway-demo/example/config"
   "grpc-gateway-demo/example/service"
   pb "grpc-gateway-demo/pkg/proto/book"
)

// Run 启动rpc服务
func Run() error {
   log := grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard)
   grpclog.SetLoggerV2(log)
   grpcAddr := config.GetRpcAddr()
   l, err := net.Listen("tcp", grpcAddr)
   if err != nil {
      grpclog.Errorf("tcp listen failed: %v", err)
      return err
   }
   defer func() {
      if err = l.Close(); err != nil {

         fmt.Fprintln(os.Stderr, err)
      }
   }()

   s := grpc.NewServer()

   //     注册服务
   registerServer(s)
   log.Infof("Serving gRPC on %s", l.Addr())
   return s.Serve(l)
}

func registerServer(s *grpc.Server) {
   pb.RegisterBookServiceServer(s, service.NewBookService())
}
```

#### 启动gateway

```go
package gateway

import (
   "context"
   "io"
   "net/http"
   "os"

   "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
   "google.golang.org/grpc"
   "google.golang.org/grpc/credentials/insecure"
   "google.golang.org/grpc/grpclog"

   "grpc-gateway-demo/example/config"
   pb "grpc-gateway-demo/pkg/proto/book"
)

func Run() error {
   log := grpclog.NewLoggerV2(os.Stdout, io.Discard, io.Discard)
   grpclog.SetLoggerV2(log)
   ctx := context.Background()
   // 创建grpc连接
   conn, err := grpc.DialContext(ctx, config.GetRpcAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
   if err != nil {
      log.Error("dial failed: %v", err)
      return err
   }
   mux := runtime.NewServeMux()
   err = newGateway(ctx, conn, mux)
   if err != nil {
      log.Error("register handler failed: %v", err)
      return err
   }
   server := http.Server{
      Addr:    config.GetHttpAddr(),
      Handler: mux,
   }
   log.Infof("Serving Http on %s", server.Addr)
   err = server.ListenAndServe()
   if err != nil {
      return err
   }
   return nil
}

func newGateway(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) error {

   err := pb.RegisterBookServiceHandler(ctx, mux, conn)
   if err != nil {
      return err
   }
   return nil
}
```

#### 启动main

```go
package main

import (
   "fmt"
   "os"

   "grpc-gateway-demo/example/gateway"
   "grpc-gateway-demo/example/server"
)

func main() {
   go func() {
      err := server.Run()
      if err != nil {
         fmt.Fprintln(os.Stderr, err)
         os.Exit(1)
      }
   }()

   err := gateway.Run()
   fmt.Fprintln(os.Stderr, err)
   os.Exit(1)

}
```

## 自定义路由

有些时候，rpc的服务不能满足我们的需求，比如文件上传下载，使用proto文件定义api以及实现是无法实现的，这个时候需要我们额外的添加上自定义路由来完成相关操作。



## 使用中间件

