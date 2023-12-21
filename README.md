This is a GoLang capstone project which uses gofiber and gorm frameworks and mysql database.
In this I have created a backend rest api and used microservices architecture to serve the movie tickets booking purpose.
Project contains 4 independent microservices - User,Movie,Theatre,Booking each service can be independently deployed to the server.
Some microservices communicate internally each other using http end points.
User authentication and authorization enabled using Jwt token.
Error Handling with user-friendly messages.
Encryptions and decryptions of needed entities.
Each service having configuration file for reference (with dummy values) which serve the purpose of giving runtime values (no hardcoded values in the code).
