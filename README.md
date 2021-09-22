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


Bengali Sweets! As a Bengali, I can comfortably say the demand for Bengali sweets in Canada is very high amongst the Bengali diaspora, so much so that my friends and family try to squeeze in as many treats in their suitcases on their return flights to Toronto from Kolkata! And the demand is definitely not unwarranted - the softness of a Rasgulla, a round treat soaked in a sugar solution and the creaminess of a Chom-Chom are absolutely to die for. It is not uncommon to be served a Rasgulla alongside a cup of tea when you are a guest at someone's home in Kolkata or Dhaka. 

I'd love to share this gift of Bengal with Canadians. Brown folk and other minorities are generally considered as monoliths and there's very little interest in learning more about different cultures, languages and foods. Hopefully, a warm Chom-Chom delivered straight to the customer's door through a Shopify store would allow us to shed light on the beautiful things present in every culture. 