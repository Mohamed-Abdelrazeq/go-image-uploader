package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Session *s3.S3
var bucket string

func init() {
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(os.Getenv("AWS_REGION")),
			},
		),
	)
	s3Session = s3.New(sess)
	bucket = os.Getenv("AWS_BUCKET_NAME")
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return uploadImage(request)
	case "GET":
		if request.PathParameters["filename"] != "" {
			return showImage(request)
		}
		return health()
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method not allowed",
		}, nil
	}
}

func uploadImage(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	contentType := request.Headers["Content-Type"]
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid content type",
		}, nil
	}

	body := strings.NewReader(request.Body)
	reader := multipart.NewReader(body, contentType[30:])

	part, err := reader.NextPart()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid file",
		}, nil
	}
	defer part.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(part)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to read file",
		}, nil
	}

	filePath := filepath.Join("uploads", part.FileName())
	_, err = s3Session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to upload image: %s", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Image uploaded successfully",
	}, nil
}

func showImage(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileName := request.PathParameters["filename"]
	filePath := filepath.Join("uploads", fileName)

	resp, err := s3Session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to retrieve image: %s", err.Error()),
		}, nil
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       imageBase64,
		Headers: map[string]string{
			"Content-Type": "image/jpeg", // Adjust the content type as needed
		},
	}, nil
}

func health() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Healthy",
	}, nil
}

func main() {
	lambda.Start(handler)
}
