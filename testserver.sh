#!/bin/bash

echo "LIST ALL"
curl http://localhost:8080/api/users
echo
echo

echo "SHOULD RETURN JOAO"
curl http://localhost:8080/api/users/10
echo
echo

echo "SHOULD FAIL!\n"
curl http://localhost:8080/api/users/12
echo
echo

echo "ADD USER"
curl -v http://localhost:8080/api/users -d '{"first_name": "MARCO", "last_name": "BOY", "biography": "NICE GUY"}'
echo
echo

echo "SHOW NEW USER"
curl http://localhost:8080/api/users/11
echo 
echo

echo "LIST ALL"
curl http://localhost:8080/api/users
echo
echo

echo "UPDATE USER"
curl -v -X PUT http://localhost:8080/api/users/11 -d '{"first_name": "MARCO", "last_name": "BOY", "biography": "BAD BOY"}'
echo
echo

echo "LIST ALL"
curl http://localhost:8080/api/users
echo
echo

echo "DELETE TEST\n"
curl -v -X DELETE http://localhost:8080/api/users/11
echo
echo

echo "LIST ALL"
curl http://localhost:8080/api/users
