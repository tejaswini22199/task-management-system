# taskmanagementsystem
Task Management System in Go 

## Problem breakdown and design decisions

A Task Management System with CRUD (Create, Read, Update, Delete) capabilities for tasks. Each task can be assigned to multiple users and includes a status field for filtering. The API is designed for a single client, with pagination set to 20 for GET requests. Tasks are only visible to the users assigned to them.

Design Decisions: 
1. I'm using micro services architecture because
2. I'm using relational database postgresql database because 
3. I'm containerising the database and task service apis in seperate docker containers to ensure seamless running of these services in different environments. 
4. I'm using REST API for http calls because 



## Instructions to run the service

#### To remove any old images and containers if rerunning 

docker ps -a 
docker stop <container_id> 
docker rm <container_id>   

docker images
docker rmi <image_id>  


#### To build and run : 
docker build --no-cache -t task-management-system .

docker-compose down -v  -> To stop all the docker containers 
docker compose up -d   

#### To run the db docker container 

Docker ps 

docker exec -it {postgres container id} psql -U postgres 

#### Login to the db and check the table values 
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

## 🚀 API Documentation (endpoints, request/response examples)

### 🌍 Base URL

1. Register User 
Endpoint:

Description: 

Request: 
```
curl -X POST "http://localhost:8080/register"  -H "Content-Type: application/json" -d '{
           "Name": "tejaswini",
           "Email": "tejaswini@example.com",
           "Password": "securepassword"
         }'
         ```
Response: 

```{"message":"User registered successfully","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4Mjg2MTd9.Mj8vug-iWGSHo7PzYSxnFqlCo3rQPR8vu0J6ah0u46s","user_id":1}%   
```
2. Login as User 

Endpoint:

Description: 

Request: 

```
curl -X POST "http://localhost:8080/login"  -H "Content-Type: application/json" -d '{
           "Email": "tejaswini@example.com",
           "Password": "securepassword"
         }'
```

Response: 

```
{"message":"Login successful","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4Mjg3MDB9.1jdxFBT7kNZ4NxlUAvtMIXCgzd0oCzkfVoBtQkp_YQg","user_id":1}
```


3. Create a Task 

Endpoint:

Description: 

Request 1: 
curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Complete API Integration",
           "description": "Implement API calls and authentication",
           "status": "pending",
           "user_ids": [1, 3]
         }'
Response 1: 
{"error":"Some user IDs do not exist","invalid_users":[3]}

Request 2: 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Complete API Integration",
           "description": "Implement API calls and authentication",
           "status": "To Do",
           "user_ids": [1, 2]
         }'


Response 2: 
{"assigned_users":[1,2],"message":"Task created successfully","task_id":1}%     

Request 3: 
{"error":"Invalid status. Allowed values: To Do, In Progress, Completed"}%                                                                                      
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % 
curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Complete API Integration",
           "description": "Implement API calls and authentication",
           "status": "pending",
           "user_ids": [1, 2]
         }'
 
Response 3: 
{"error":"Invalid status. Allowed values: To Do, In Progress, Completed"}%   

Request 4: 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY" \
     -d '{
           "title": "Fix Bugs in API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'

Response 4: 
{"error":"Invalid token"}%  

taskdb=# select * from tasks;
 id |             title             |                 description                 |   status    | created_by |         created_at         
----+-------------------------------+---------------------------------------------+-------------+------------+----------------------------
  1 | Complete API Integration      | Implement API calls and authentication      | To Do       |          1 | 2025-03-23 17:58:55.451846
  2 | Complete Frontend Integration | Implement frontend integration with backend | Completed   |          1 | 2025-03-23 17:59:44.037258
  3 | Fix Bugs in API               | Fix critical bugs in the API                | In Progress |          2 | 2025-03-23 18:09:00.198216

5. Get all Tasks 

Endpoint:

Description: 

Request1: 

tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X GET http://localhost:8080/tasks \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY"

     Here user id is 2 (after decoding the bearer token). So, the task was fetched since user 2 is assigned to it. 
    {
  "user_id": 2,
  "exp": 1742833386
}


Response1:
{"limit":20,"page":1,"tasks":[{"created_at":"2025-03-23T16:24:52.835104Z","description":"Implement API calls and authentication","id":1,"status":"To Do","title":"Complete API Integration","user_ids":[1,2]}],"total_pages":1,"total_tasks":1}%  

Request 2: 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X GET http://localhost:8080/tasks \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY 
"
{
  "user_id": 3,
  "exp": 1742833386
}

Here user id is 3 (after decoding the bearer token). Since the user is not present in db. He is not authorized. 



Response 2: 

{"error":"Invalid token"}%   

Request 3 : 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X GET http://localhost:8080/tasks  

Response 3: 

{"error":"Authorization header required"}%                           


4. Get a task by taskId 
Endpoint:

Description: 

Request1: 
curl -X GET http://localhost:8080/tasks/1 \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY"

Response1: 
{"task":{"created_at":"2025-03-23T16:24:52.835104Z","description":"Implement API calls and authentication","id":1,"status":"To Do","title":"Complete API Integration","user_ids":[1,2]}}%    

Request2: 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X GET http://localhost:8080/tasks/1 \ 
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY
"
Response2: 
{"error":"Invalid token"}%  

Request3: 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X GET http://localhost:8080/tasks/2 \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY"

Response3: 
{"error":"Task not found"}%   
6. Get Tasks filtered by Status

Endpoint:

Description: 

Request1: 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X GET http://localhost:8080/tasks/status/In%20Progress \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY"


Response1: 
{"tasks":null}%  


Request 2:  curl -X GET http://localhost:8080/tasks/status/To%20Do \      
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY"

Response 2: 
{"tasks":[{"created_at":"2025-03-23T16:24:52.835104Z","description":"Implement API calls and authentication","id":1,"status":"To Do","title":"Complete API Integration"}]}%    

7. Update Task 

Endpoint:

Description:

Request1:

curl -X PUT http://localhost:8080/tasks/1 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY" \
     -d '{
           "title": "Updated Task Title",
           "description": "Updated Task Description",
           "status": "In Progress"
         }'

Response1:

{"id":1,"title":"Updated Task Title","description":"Updated Task Description","status":"In Progress","created_by":0,"created_at":"2025-03-23T16:24:52.835104Z"}%   

from db:
taskdb=# select * from tasks;
 id |       title        |       description        |   status    | created_by |         created_at         
----+--------------------+--------------------------+-------------+------------+----------------------------
  1 | Updated Task Title | Updated Task Description | In Progress |          1 | 2025-03-23 16:24:52.835104 

  Request2 : 
  tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X PUT http://localhost:8080/tasks/4 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'

  Response 2: 

{"error":"task not found"}%  

  Request 3: 

  tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X PUT http://localhost:8080/tasks/3 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'


    This token belongs to user id 1 
    eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc

  Response 3: 

  {"error":"you are not authorized to access or modify this task"}%    

8. Delete Task 
Endpoint:

Description: 

Request: 
tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X DELETE http://localhost:8080/tasks/4 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'


Response: 
{"error":"task not found"}%   

Request 2: 

tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X DELETE http://localhost:8080/tasks/3 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3NDI4MTIzNTJ9.PchfuUCqcaCR2JhUfWm7nkcMcDhvMJDsYGwxcS2Jygc" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'

Response 2:

{"error":"you are not authorized to access or modify this task"}%   

Request 3: 

tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X DELETE http://localhost:8080/tasks/3 \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJleHAiOjE3NDI4MzMzODZ9.VczoPaa2E28eFrio0OasqqyakIHR2OxPKPqs_AHJfWY" \
     -d '{
           "title": "Fix Bugs in Rest API",
           "description": "Fix critical bugs in the API",
           "status": "In Progress",
           "user_ids": [2]
         }'
Response 3:
{"message":"Task deleted successfully"}%    

## Explanation of how the service demonstrates microservices concepts



## Tasks Pending 

4. if someone login without registering, it gives followinf reposnse tejaswinivakkalagaddi@Tejaswinis-MacBook-Air task-management-system % curl -X POST "http://localhost:8080/login"  -H "Content-Type: application/json" -d '{
           "Email": “arunav@example.com",
           "Password": "password"
         }'
{"error":"Invalid input"}%    
make it more clear
7. rephrase code to ensure maintainability and cleanliness , extendability and srp 

## Tasks Completed 

1. Test CreateTask API
2. check the auth bearer token thing 
4. Test Login API 
5. Test Register API
3. Test GetTaskById API  
1. Test GetTasksByStatus API 

2. Test updateTask API 
3. Test deleteTask API