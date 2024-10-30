package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/Magowtham/go_file_streaming_server/proto/filestream"
	"google.golang.org/grpc"
)

type server struct {
	filestream.UnimplementedFileStreamServiceServer
}

func (s *server) DownloadFile(req *filestream.FileRequest, stream filestream.FileStreamService_DownloadFileServer) error {

	fileName := fmt.Sprintf("server/%v", req.Filename)
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Printf("error occurred -> %v", err.Error())
		return err
	}

	defer file.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("error occurred -> %v", err.Error())
			return err
		}

		fileChunk := &filestream.FileChunk{
			Chunk:     buffer,
			ChunkSize: int32(n),
		}

		err = stream.Send(fileChunk)

		if err != nil {
			fmt.Printf("error occurred -> %v", err.Error())
			return err
		}
	}

	return nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8000")

	if err != nil {
		fmt.Printf("error occurred -> %v", err.Error())
		return
	}

	grpcServer := grpc.NewServer()

	filestream.RegisterFileStreamServiceServer(grpcServer, &server{})

	fmt.Printf("grpc server is running on 0.0.0.0:8000\n")

	err = grpcServer.Serve(listener)

	if err != nil {
		fmt.Printf("error occurred -> %v", err.Error())
		return
	}
}
