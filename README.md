# Tiny-Url-Software

This Git repository contains a Tiny URL software implemented in the Go programming language. The software provides a web-based service that converts long URLs into short, more manageable URLs. Users can input a long URL, and the software generates a unique short URL that redirects to the original long URL.

## Table of Contents

- [Installation](#installation)
- [User Authentication](#User-Authentication)
- [URL Shortening](#URL-Shortening)
- [Contributing](#Contributing)
- [License](#license)

## Installation


1. Clone the repository: git clone https://github.com/your-username/your-project.git
2. Change to the project directory:
```bash
$ cd URL_SHOTNER_APPLICATION/
````
3. Install dependencies:
```bash
$ go get -d -v $(cat dependencies.txt)
```
5. Set up the environment variables:
- `DATABASE_URL`: Connection string to your PostgreSQL database
- `Port`: Add your desire PORT
6. Migrate the models ``` go run migrate/migrate.go```
7. Get the build of project:
```bash
$ go build -o URl_SHORTNER_APPLICATION 
```
8. Run the project:
```bash
$ go run main.go
```
## User Authentication

1. Sign up: Send a POST request to `/signup` with the following JSON body:
   `{
   "username": "your-username",
   "password": "your-password"
   }`
2. Log in: Send a POST request to `/login` with the same JSON body. You'll receive a JWT token in the response.
## URL Shortening

1. Authenticate: Include the JWT token in the `Authorization` header for subsequent requests.

2. Create a short URL: Send a POST request to `/url` with the following JSON body:
   Get all URLs: Send a GET request to `/all` to retrieve a list of all created URLs.

4. Redirection: Access a short URL by navigating to `http://localhost:8080/your-shortcode`. You'll be redirected to the original long URL.

## Contributing

1. Fork the repository.

2. Create a new branch: `git checkout -b feature-new-feature`
3. Make your changes and commit: `git commit -m "Add new feature`
4. Push to the branch: `git push origin feature-new-feature`
5. Submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
Testing
