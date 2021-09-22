# Image Repository
## Shopify Winter 2022 Backend Developer Intern Challenge

Hi! This is my submission for Shopify's backend challenge. This repo contains the code and instructions to run and test an image repository API.

This API allows users to sign up and upload, delete and search for images in bulk, all while being secured and authenticated by JWT tokens. This project is written in Go and uses PostgreSQL as its primary database. Image artifacts are stored in AWS S3 buckets.

### Features

- Secure bulk uploads 
- Secure bulk deletes, cannot delete another user's images
- Image search by user IDs
- Authentication middleware using JWT tokens

### Usage

Docker and `docker-compose` are required. Please clone this repository and execute the following command in the `/backend` directory to run this project: 
```
docker-compose up --build
```
The image repository API can now be accessed at `localhost:8000` 

To stop the containers, press `ctrl + c` and run `docker-compose down` to remove the containers. 

### Operations
Please visit [this Postman doc](https://documenter.getpostman.com/view/13042235/UUxuj9xn) to see detailed documentation on all operations and instructions on how to test the API using Postman, cURL or other tools.  
- `POST /signup`: Create an account
- `POST /login`: Log into an account
- `PUT /images`: Upload images
- `DELETE /images`: Deletes images uploaded by user
- `GET /search`: Search for images by user IDs
- `GET /image/{id}`: Fetches the requested image

