# taskmanagementsystem

Task Management System in Go 


## Problem breakdown and design decisions

A Task Management System with CRUD (Create, Read, Update, Delete) capabilities for tasks. Each task can be assigned to multiple users and includes a status field for filtering. The API is designed for a single client, with pagination set to 20 for GET requests. Tasks are only visible to the users assigned to them.

### Design Decisions for Task Management Application

#### Microservices Architecture: 

I have opted for a microservices architecture for the Task Management System to ensure a clear separation of concerns and improve the modularity of the application. Microservices allow for independent deployment of each service, ensuring that different components (e.g., Task Management) can evolve without tightly coupling them with other parts of the system. This architecture also makes it easier to scale individual services based on demand.

#### Relational Database (PostgreSQL): 

For data storage, I have chosen PostgreSQL as the relational database. Given the structured nature of the data and the relationships between different entities (e.g., tasks, users), PostgreSQL is well-suited to handle these requirements while ensuring consistency, atomicity, and data integrity. The relational model also supports complex queries and transactions effectively. In the future, if load increases, PostgreSQL can be scaled horizontally using read replicas and sharding to ensure better performance.

#### Containerization with Docker: 

Both the PostgreSQL database and the Task Management service will be containerized using Docker. This approach ensures that the services can be consistently deployed across various environments, such as development, staging, and production. Containerization also enables the application to run seamlessly on different machines, improving the portability and scalability of the system.

#### RESTful API for Communication: 

To facilitate communication between the client and the service, and between services, I am utilizing REST APIs. This choice allows for simple, synchronous communication between the client and service, which fits the current requirements of the system. REST APIs are easy to implement and maintain, and they provide a clear and consistent interface. As the system grows, asynchronous communication options like gRPC or message queues (e.g., RabbitMQ or Kafka) can be considered to improve scalability and handle more complex use cases.

#### Task Management Features: 

The Task Management System will support basic CRUD operations (Create, Read, Update, Delete) for tasks. Additionally, pagination will be implemented for the GET /tasks endpoint to manage large datasets efficiently. Filtering by task status (e.g., GET /tasks?status=Completed) will also be supported to enable users to easily find tasks based on their current state.

#### Conclusion:



This design follows best practices of microservices architecture, ensuring that each component has a clear responsibility . The use of PostgreSQL ensures data consistency, and containerization with Docker facilitates smooth deployment and scalability. RESTful APIs provide a straightforward way to handle synchronous communication between services, with the potential to scale asynchronously in the future.

The application adheres to Object-Oriented Programming (OOP) principles to promote cleaner, maintainable, and scalable code:

DRY (Don't Repeat Yourself): The codebase avoids redundancy by reusing components and functions where applicable. This reduces the risk of errors and improves maintainability.

KISS (Keep It Simple, Stupid): The system design emphasizes simplicity in the API design and service structure. Each service focuses on a specific responsibility, ensuring that it is easy to understand, extend, and maintain.

The system is designed to handle growth, both in terms of functionality and scalability, making it flexible for future enhancements.

## Instructions to run the service

#### Install the following dependencies 

```
* go 1.24
* Docker 
* postgresql 
* Install go direct dependencies using the following command -> go mod tidy
* After running these commands, Go will add the necessary dependencies to your go.mod and go.sum files.
```

#### To remove any old images and containers if rerunning 

```
docker ps -a 
docker stop <container_id> 
docker rm <container_id> 

docker images
docker rmi <image_id>  

```

#### To build and run the docker containers:

```
docker build --no-cache -t task-management-system .

docker-compose down -v  -> To stop all the docker containers 
docker compose up -d   

```

#### To run the database container 

```
Docker ps 

docker exec -it {postgres container id} psql -U postgres 

```

#### Login to the db and check the table values 

```
\c taskdb - connect to task db 

taskdb=# select * from users;

 id | name | email | password | created_at 
----+------+-------+----------+------------
(0 rows)

taskdb=# select * from tasks

taskdb-# ;
 id | title | description | status | created_at 
----+-------+-------------+--------+------------
(0 rows)

taskdb=# select * from tasks_users;

 id | user_id | task_id | assigned_at 
----+---------+---------+-------------
(0 rows)

```


## 🚀 API Documentation (endpoints, request/response examples)


### 🌍 Base URL

```
{{host}}:{{port}}
```

Main Service : Main Service (As of now, it doesn't serve any APIs)

* host = "http://localhost/"
* port = 8080

Auth Service: Responsible for serving Authorization APIs like register and login

* host = "http://localhost/"
* port = 8000 

Tasks Service: Responsible for serving Tasks CRUD APIs 

* host = "http://localhost/"
* port = 8001

#### 1. Register User 

Endpoint: 

```
POST {{host}}:{{port}}/auth/register
```

Description: 

```
This API endpoint is used to register a new user. It requires the user's details (Name, Email, and Password)
 to create a new account. Upon successful registration, the user's data is stored in the database, and the 
 user is provided with a token which they can use while using the Tasks APIs.
```

Request: 

```
curl -X POST "http://localhost:8000/auth/register"  -H "Content-Type: application/json" -d '{
           "Name": "tejaswini",
           "Email": "tejaswini@example.com",
           "Password": "securepassword"
         }'
```

```
curl -X POST "http://localhost:8000/auth/register"  -H "Content-Type: application/json" -d '{
           "Name": "aruna",
           "Email": "arunav@example.com",
           "Password": "secure2password"
         }'  
```

Response: 

```
{"message":"User registered successfully","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMzYzOTYsInVzZXJfaWQiOjF9.--VxoU1EHTXAa_JZ5ZKIVaKyV-Huos90nfELSuPxacY","user_id":1}
```

```
{"message":"User registered successfully","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMzY0MDIsInVzZXJfaWQiOjJ9.LoM8bVUXNB78L4Twb4K5T1eZ8F-4ktjF54TYWh1x3sI","user_id":2}  
```

Database Output:

``` 
taskdb=# select * from users;
 id |   name    |         email         |                           password                           |         created_at         
----+-----------+-----------------------+--------------------------------------------------------------+----------------------------
  1 | tejaswini | tejaswini@example.com | $2a$14$6oIEf8XnhdU/oEwcqJEnBeLbT2VY9P55LW0WFNMLUR08oA/2aVkPm | 2025-03-26 00:46:36.361225
  2 | aruna     | arunav@example.com    | $2a$14$fXOF9rXUpISibOJKO2Qce.2VsYAtMWGDflj4M8uSWJrT1wZa5iulW | 2025-03-26 00:46:42.169147
```


#### 2. Login as User 

Endpoint:

```
POST {{host}}:{{port}}/auth/login
```

Description: 

```

This API endpoint is used for user authentication. It requires the user’s email and password. If the credentials are valid, the server generates a JSON Web Token (JWT) and returns it to the user for subsequent requests. The JWT can be used for authenticating the user in further API calls.

```

Request: 


```

curl -X POST "http://localhost:8000/auth/login"  -H "Content-Type: application/json" -d '{
           "Name": "tejaswini",
           "Email": "tejaswini@example.com",
           "Password": "securepassword"
         }'                                                    
```
```
curl -X POST "http://localhost:8000/auth/login"  -H "Content-Type: application/json" -d '{
           "Name": "aruna",
           "Email": "arunav@example.com",
           "Password": "secure2password"
         }'

```
Response: 


```
{"message":"Login successful","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMzY0MDgsInVzZXJfaWQiOjF9.K4lM8YWBjefNlPjKjSNbfSbOukgkyUeksJSvLIU3gXY","user_id":1}   
```

```
{"message":"Login successful","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMzY0MTQsInVzZXJfaWQiOjJ9.UHspiDItQZV2Lxt1RecXlGI2Ny402OxL549V9piHfAs","user_id":2}
```                                                        


#### 3. Create a Task 

Endpoint:

```
POST {{host}}:{{port}}/tasks
```

Description: 

```

This API endpoint is used to create a new task. You must provide the task details, including the title, description, status, and the users assigned to the task.
The task is saved in the database and assigned to the provided users.


title (string): The title of the task.

description (string): A brief description of the task.

status (string): The current status of the task (e.g., "To Do", "In Progress", "Completed").

user_ids (array of integers): A list of user IDs who are assigned to the task.


```


Request 1: 

```

curl -X POST http://localhost:8001/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMDA5MjksInVzZXJfaWQiOjF9.Uy79rzdHPpQI6xR3_RYUD9B1qJN51Fbz61hst-biCKg" \
     -d '{
           "title": "Complete API Integration",
           "description": "Implement API calls and authentication",
           "status": "To Do",  
           "user_ids": [1, 3]
         }'


```

Response 1: 

```

{"error":"Some user IDs do not exist","invalid_users":[3]}%   

```

Request 2: 


```
 curl -X POST http://localhost:8001/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMDA5MjksInVzZXJfaWQiOjF9.Uy79rzdHPpQI6xR3_RYUD9B1qJN51Fbz61hst-biCKg" \
     -d '{
           "title": "Complete API Integration",
           "description": "Implement API calls and authentication",
           "status": "To Do",
           "user_ids": [1, 2]
         }'

```

```
curl -X POST http://localhost:8001/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMzY0MTQsInVzZXJfaWQiOjJ9.UHspiDItQZV2Lxt1RecXlGI2Ny402OxL549V9piHfAs" \
     -d '{
           "title": "Fix Bugs in API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'
```

```
curl -X POST http://localhost:8001/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMzY0MDgsInVzZXJfaWQiOjF9.K4lM8YWBjefNlPjKjSNbfSbOukgkyUeksJSvLIU3gXY" \
     -d '{
           "title": "Complete Frontend Integration",
           "description": "Implement frontend integration with backend I",
           "status": "Completed",
           "user_ids": [1]
         }'                                                                                                  
```

Response 2: 


```
{"message":"Task created successfully","task":{"id":1,"title":"Complete API Integration","description":"Implement API calls and authentication","status":"To Do","created_by":1,"created_at":"0001-01-01T00:00:00Z"}}
```

```
{"message":"Task created successfully","task":{"id":2,"title":"Fix Bugs in API","description":"Fix critical bugs in the API","status":"In Progress","created_by":2,"created_at":"0001-01-01T00:00:00Z"}}
```

```
{"message":"Task created successfully","task":{"id":3,"title":"Complete Frontend Integration","description":"Implement frontend integration with backend I","status":"Completed","created_by":1,"created_at":"0001-01-01T00:00:00Z"}}
```

Database Output: 

```
taskdb=# select * from tasks;
 id |             title             |                  description                  |   status    | created_by |         created_at         
----+-------------------------------+-----------------------------------------------+-------------+------------+----------------------------
  1 | Complete API Integration      | Implement API calls and authentication        | To Do       |          1 | 2025-03-26 00:47:18.534694
  2 | Fix Bugs in API               | Fix critical bugs in the API                  | In Progress |          2 | 2025-03-26 01:14:47.298008
  3 | Complete Frontend Integration | Implement frontend integration with backend I | Completed   |          1 | 2025-03-26 01:17:48.939929
(3 rows)

taskdb=# select * from tasks_users;
 id | user_id | task_id |        assigned_at         
----+---------+---------+----------------------------
  1 |       1 |       1 | 2025-03-26 00:47:18.534694
  2 |       2 |       1 | 2025-03-26 00:47:18.534694
  3 |       2 |       2 | 2025-03-26 01:14:47.298008
  4 |       1 |       3 | 2025-03-26 01:17:48.939929
(4 rows)

```

Request 3: 


```

curl -X POST http://localhost:8001/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMDA5MjksInVzZXJfaWQiOjF9.Uy79rzdHPpQI6xR3_RYUD9B1qJN51Fbz61hst-biCKg" \
     -d '{
           "title": "Complete API Integration",
           "description": "Implement API calls and authentication",
           "status": "pending",
           "user_ids": [1, 3]
         }'


 ```

Response 3: 

```

{"error":"Invalid status. Allowed values: To Do, In Progress, Completed"}

```

Request 4: 


```

curl -X POST http://localhost:8001/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Complete API Integration",
           "description": "Implement API calls and authentication",
           "status": "pending",
           "user_ids": [1, 2]
         }'

```

Response 4: 

```
{"error":"Invalid token"}
2025-03-26 06:37:32 Error: token has invalid claims: token is expired
```

db output:

```

taskdb=# select * from tasks;
 id |             title             |                 description                 |   status    | created_by |         created_at         
----+-------------------------------+---------------------------------------------+-------------+------------+----------------------------
  1 | Complete API Integration      | Implement API calls and authentication      | To Do       |          1 | 2025-03-23 17:58:55.451846
  2 | Complete Frontend Integration | Implement frontend integration with backend | Completed   |          1 | 2025-03-23 17:59:44.037258
  3 | Fix Bugs in API               | Fix critical bugs in the API                | In Progress |          2 | 2025-03-23 18:09:00.198216

```

#### 4. Get all Tasks 

Endpoint:

```
GET {{host}}:{{port}}/tasks
```

Description: 

```
Description:
This API endpoint retrieves a list of all tasks stored in the system. The tasks may include various details such as task names, descriptions, statuses, due dates, or other relevant information depending on the implementation. This endpoint is typically used to fetch and display all tasks to a user or an admin in a task management application.
```

Request1: 

```
curl -X GET "http://localhost:8001/tasks?page=1&limit=5" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwNzAwMjIsInVzZXJfaWQiOjJ9.WZpTaZ2njjbSm1i_3eO6RAGjsRe_uUCwEjEwmfQP3eQ"
```

Database O/P:

```
taskdb=# select * from tasks_users;
 id | user_id | task_id |        assigned_at         
----+---------+---------+----------------------------
  1 |       1 |       1 | 2025-03-26 10:16:54.443524
  2 |       1 |       2 | 2025-03-26 10:17:00.12955
  3 |       2 |       2 | 2025-03-26 10:17:00.12955
  4 |       2 |       3 | 2025-03-26 10:17:22.588586
  5 |       2 |       4 | 2025-03-26 10:19:10.39979
  6 |       2 |       5 | 2025-03-26 10:19:18.598814
  7 |       2 |       6 | 2025-03-26 10:19:27.546983
  8 |       2 |       7 | 2025-03-26 10:19:34.771187
  9 |       2 |       8 | 2025-03-26 10:19:42.61439
```
Response1:

```

{
  "limit": 5,
  "page": 1,
  "tasks": [
    {
      "id": 2,
      "title": "Complete API Integration",
      "description": "Implement API calls and authentication",
      "status": "To Do",
      "created_by": 0,
      "created_at": "2025-03-26T10:17:00.12955Z"
    },
    {
      "id": 3,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:17:22.588586Z"
    },
    {
      "id": 4,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 2 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:10.39979Z"
    },
    {
      "id": 5,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 3 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:18.598814Z"
    },
    {
      "id": 6,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 4 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:27.546983Z"
    }
  ],
  "total_pages": 2,
  "total_tasks": 7
}

```

Request 2: 

```

curl -X GET http://localhost:8001/tasks \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY 
"
{
  "user_id": 3,
  "exp": 1742833386
}

```

* Here user id is 3 (after decoding the bearer token). Since the user is not present in db. He is not authorized. 

Response 2: 

```

{"error":"Invalid token"}% 

```

Request 3 : 

```

 curl -X GET http://localhost:8001/tasks  

```

Response 3: 

```

{
    "error":"Authorization header required"

}

```                   


#### 5. Get a task by taskId 

Endpoint: 


```

GET {{host}}:{{port}}/tasks/:id

```

Description: 

```

This API endpoint retrieves the details of a specific task based on the task ID. It returns the task's attributes such as title, description, status, and the users assigned to the task.
Parameters:

id (int): The unique ID of the task.

```

Request1: 

```

curl -X GET http://localhost:8001/tasks/1 \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMDA4NzEsInVzZXJfaWQiOjF9.-Hq181pG0hTAfSRjIByfZe9m_V7VdwvbPb58kr6dLhg"

```

Response1: 


```
{
  "task": {
    "id": 1,
    "title": "Complete Frontend Integration",
    "description": "Implement frontend integration with backend",
    "status": "Completed",
    "created_by": 1,
    "created_at": "2025-03-26T10:16:54.443524Z"
  }
}
```

Request2: 

```

curl -X GET http://localhost:8001/tasks/1 \ 
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY"

 ```

Response2: 

```

{"error":"Invalid token"}

```

Request3: 

```

 curl -X GET http://localhost:8001/tasks/11 \              
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMDA4NzEsInVzZXJfaWQiOjF9.-Hq181pG0hTAfSRjIByfZe9m_V7VdwvbPb58kr6dLhg"

```

Response3: 

```

{"error":"Task not found"}

```

#### 6. Get Tasks filtered by Status

Endpoint:

```

GET {{host}}:{{port}}/tasks/status/:status

```

Description: 

```

This API endpoint retrieves a list of tasks filtered by their status (e.g., "To Do", "In Progress", "Completed"). It allows you to fetch tasks that match the specified status.
Parameters:


status (string): The status of the task (e.g., "To Do", "In Progress", "Completed").

```

```
Database 
 id |             title             |                 description                 |   status    | created_by |         created_at         
----+-------------------------------+---------------------------------------------+-------------+------------+----------------------------
  1 | Complete Frontend Integration | Implement frontend integration with backend | Completed   |          1 | 2025-03-26 10:16:54.443524
  2 | Complete API Integration      | Implement API calls and authentication      | To Do       |          1 | 2025-03-26 10:17:00.12955
  3 | Fix Bugs in API               | Fix critical bugs in the API                | In Progress |          2 | 2025-03-26 10:17:22.588586
  4 | Fix Bugs in API               | Fix critical bugs 2 in the API              | In Progress |          2 | 2025-03-26 10:19:10.39979
  5 | Fix Bugs in API               | Fix critical bugs 3 in the API              | In Progress |          2 | 2025-03-26 10:19:18.598814
  6 | Fix Bugs in API               | Fix critical bugs 4 in the API              | In Progress |          2 | 2025-03-26 10:19:27.546983
  7 | Fix Bugs in API               | Fix critical bugs 5 in the API              | In Progress |          2 | 2025-03-26 10:19:34.771187
  8 | Fix Bugs in API               | Fix critical bugs 6 in the API              | In Progress |          2 | 2025-03-26 10:19:42.61439
```

Request1: 

``` 

curl -X GET http://localhost:8001/tasks/status/In%20Progress \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMDA4NzEsInVzZXJfaWQiOjF9.-Hq181pG0hTAfSRjIByfZe9m_V7VdwvbPb58kr6dLhg"


```

Response1: 

```

{
  "tasks": [
    {
      "id": 3,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:17:22.588586Z"
    },
    {
      "id": 4,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 2 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:10.39979Z"
    },
    {
      "id": 5,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 3 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:18.598814Z"
    },
    {
      "id": 6,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 4 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:27.546983Z"
    },
    {
      "id": 7,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 5 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:34.771187Z"
    },
    {
      "id": 8,
      "title": "Fix Bugs in API",
      "description": "Fix critical bugs 6 in the API",
      "status": "In Progress",
      "created_by": 0,
      "created_at": "2025-03-26T10:19:42.61439Z"
    }
  ]
}

```
Db output: 
```
taskdb=# select * from users
taskdb-# ;
 id |   name    |         email         |                           password                           |         created_at         
----+-----------+-----------------------+--------------------------------------------------------------+----------------------------
  1 | tejaswini | tejaswini@example.com | $2a$14$qf62mcmnYrG.tTUxDuwvEOnVFC2VQCRqUeC6zYz2pfeXjS35e8lxG | 2025-03-26 12:25:15.311715
  2 | aruna     | arunav@example.com    | $2a$14$ImvSCZBtKJmyU/ttR6728ukaRy/gJNJgMD7RkU40ZXswg8PFDARLe | 2025-03-26 12:25:19.527617
(2 rows)

taskdb=# select * from tasks;
 id |             title             |                 description                 |   status    | created_by |         created_at         
----+-------------------------------+---------------------------------------------+-------------+------------+----------------------------
  1 | Complete Frontend Integration | Implement frontend integration with backend | Completed   |          1 | 2025-03-26 12:25:33.592678
  2 | Complete API Integration      | Implement API calls and authentication      | To Do       |          1 | 2025-03-26 12:25:38.971815
  3 | Fix Bugs in API               | Fix critical bugs in the API                | In Progress |          2 | 2025-03-26 12:25:56.100886
(3 rows)

taskdb=# select * from tasks_users;
 id | user_id | task_id |        assigned_at         
----+---------+---------+----------------------------
  1 |       1 |       1 | 2025-03-26 12:25:33.592678
  2 |       1 |       2 | 2025-03-26 12:25:38.971815
  3 |       2 |       2 | 2025-03-26 12:25:38.971815
  4 |       2 |       3 | 2025-03-26 12:25:56.100886

```

Request 2: 

```

curl -X GET http://localhost:8001/tasks/status/To%20Do \      
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwNTMzNjUsInVzZXJfaWQiOjF9.cH0BOA5X3NR203ReJejZX7QYjl6YfKyZvRZmDe5j7oM" 

```

Response 2: 

```

{
  "tasks": [
    {
      "id": 2,
      "title": "Complete API Integration",
      "description": "Implement API calls and authentication",
      "status": "To Do",
      "created_by": 0,
      "created_at": "2025-03-26T12:25:38.971815Z"
    }
  ]
}         

```

Request 3: 

```
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X GET http://localhost:8001/tasks/status/In%20Progress \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwNTMzNjUsInVzZXJfaWQiOjF9.cH0BOA5X3NR203ReJejZX7QYjl6YfKyZvRZmDe5j7oM"
```

Response 3: 
```
{"tasks":[]} // there are no tasks for user1 in progress and bearer token belongs to user 1
```
#### 7. Update Task 

Endpoint:

```

PUT {{host}}:{{port}}/tasks/:id

```

Description:

```

This API endpoint updates the details of an existing task. You need to provide the task ID and the new values for the task's title, description, and status. Only the owner of the task (user who created it) is allowed to update it.
Parameters:

id (int): The ID of the task you want to update.
Request Body:

title (string): Updated title of the task.

description (string): Updated description of the task.

status (string): Updated status of the task (e.g., "To Do", "In Progress", "Completed").

```


Request1:

```

curl -X PUT http://localhost:8080/tasks/1 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY" \
     -d '{
           "title": "Updated Task Title",
           "description": "Updated Task Description",
           "status": "In Progress"
         }'

```

Response1:

```

{"id":1,"title":"Updated Task Title","description":"Updated Task Description","status":"In Progress","created_by":0,"created_at":"2025-03-23T16:24:52.835104Z"}

```

Output from db: 

```

taskdb=# select * from tasks;
 id |       title        |       description        |   status    | created_by |         created_at         
----+--------------------+--------------------------+-------------+------------+----------------------------
  1 | Updated Task Title | Updated Task Description | In Progress |          1 | 2025-03-23 16:24:52.835104 

```

Request2 : 

  ``` 
  curl -X PUT http://localhost:8001/tasks/4 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'

```

  Response 2: 

```

{"error":"task not found"}%  

```

  Request 3: 

```
 curl -X PUT http://localhost:8001/tasks/3 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'
```

    * This token belongs to user id 1 
    eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc

Response 3: 

  ```
  {"error":"you are not authorized to access or modify this task"}
  ```

#### 8. Delete Task 

Endpoint:

```
DELETE {{host}}:{{port}}/tasks/:id
```

Description: 

```
This API endpoint deletes a task by its ID. Before deletion, it checks whether the task exists and whether the requesting user is the owner of the task. If the task exists and the user has ownership, the task is deleted along with the associated user-task mappings.
Parameters:

id (int): The ID of the task you want to delete.
```

Request: 

```
 curl -X DELETE http://localhost:8001/tasks/4 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'

```

Response: 

```
{"error":"task not found"} 
```

Request 2: 

```
 curl -X DELETE http://localhost:8001/tasks/3 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
         }'
```

Response 2:

```
{"error":"you are not authorized to access or modify this task"}%   
```

Request 3: 

```
curl -X DELETE http://localhost:8001/tasks/1 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwNTMzNjUsInVzZXJfaWQiOjF9.cH0BOA5X3NR203ReJejZX7QYjl6YfKyZvRZmDe5j7oM" 

 ```

Response 3:

```
{"message":"Task deleted successfully"}
```
