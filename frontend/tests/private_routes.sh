#!/bin/bash

BASE_URL="http://localhost:5173"

# Protected routes
routes=(
  "/profile"
  "/settings"
  "/dashboard"
  "/create_job"
  "/library"
)

echo "Testing protected routes that should NOT be accessible without authentication: "

for route in "${routes[@]}"; do
  final_url=$(curl -s -o /dev/null -w "%{url_effective}" "$BASE_URL$route")

    if [[ "$final_url" == *"/login"* ]]; then
        echo "PASS: $route redirects to login"
    else
        echo "FAIL: $route does not redirect properly"
    fi
done
