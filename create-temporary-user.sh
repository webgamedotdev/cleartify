#!/bin/bash

# Function to create a temporary user account
function create_temporary_user() {
  local username=$1
  local expire_days=$2

  # Create a new user
  useradd "$username"

  # Set the account to expire after the specified number of days
  chage -E $(date -d "+$expire_days days" +%F) "$username"

  echo "Temporary user $username created and set to expire after $expire_days days."
}

# Usage
# Change the USERNAME and EXPIRE_DAYS variables as needed
USERNAME="tempuser"
EXPIRE_DAYS=7

# Call the function with the username and number of days to expire
create_temporary_user "$USERNAME" "$EXPIRE_DAYS"

