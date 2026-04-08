#!/bin/bash

BASE_URL="http://localhost:8080" # FOR NOW!

# Public routes that should be accessible
routes=(
  "/login"
  "/signup"
  "/"
)

echo "Testing public routes accessibility: "

for route in "${routes[@]}"; do
  status=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL$route")

  if [ "$status" -eq 200 ]; then
    echo "PASS: $route is accessible (HTTP $status)"
  else
    echo "FAIL: $route is NOT accessible (HTTP $status)"
  fi
done
