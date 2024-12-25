# Go Image Uploader

This project is a serverless image uploader built with Go, AWS Lambda, and S3. It allows users to upload images via a POST request and retrieve them via a GET request.

## Project Structure


- `.air.toml`: Configuration file for Air, a live reload tool for Go applications.
- `.env`: Environment variables file (not included in version control).
- `.gitignore`: Specifies files and directories to be ignored by Git.
- `bin/`: Directory for the compiled binary.
- `cmd/main.go`: Main application code.
- `go.mod`: Go module file.
- `go.sum`: Go dependencies file.
- `logs.json`: Log file (ignored by Git).
- `migrations/0001_initial.sql`: SQL migration file for creating the `images` table.
- `tmp/build-errors.log`: Log file for build errors.

## Setup

1. **Clone the repository:**
  
    ```sh
    git clone https://github.com/Mohamed-Abdelrazeq/go-image-uploader.git
    cd go-image-uploader
    ```

2. **Install dependencies:**

    ```sh
    go mod tidy
    ```

3. **Set up environment variables: Create a .env file with the following content:**

    ```sh
    AWS_REGION=your-aws-region
    AWS_BUCKET_NAME=your-s3-bucket-name
    ```

4. **Run the application locally:**

    ```sh
    air
    ```

Note: This command only works in local development without the Lambda formatted project.

## Usage
Upload an Image
Send a POST request to the endpoint with the image file in the form-data.

## License
This project is licensed under the MIT License.