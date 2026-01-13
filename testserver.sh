#!/bin/bash

# curl -X GET http://localhost:8080/api/users
# curl -v GET http://localhost:8080/api/users/10

# SHOULD FAIL!
# curl -v GET http://localhost:8080/api/users/12
curl -v http://localhost:8080/api/users -d '{"first_name": "MARCO", "last_name": "BOY", "biography": "NICE GUY"}'
curl -v http://localhost:8080/api/users/11
