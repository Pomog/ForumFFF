#!/bin/bash
# Use chmod +x getGitUdates.sh

echo "Getting Updates from GitHub"

# Switch to the main branch
if git checkout main; then
    echo "Switched to main branch successfully"
else
    echo "Error: Failed to switch to main branch"
    exit 1
fi

# Pull the latest changes from the remote main branch
if git pull origin main; then
    echo "Pulled changes from main successfully"
else
    echo "Error: Failed to pull changes from main"
    exit 1
fi

# Switch back to your branch (den_dev)
if git checkout den_dev; then
    echo "Switched back to den_dev branch successfully"
else
    echo "Error: Failed to switch back to den_dev branch"
    exit 1
fi

# Rebase your branch onto the latest changes from the local main branch
if git rebase main; then
    echo "Rebased den_dev onto main successfully"
else
    echo "Error: Failed to rebase den_dev onto main"
    exit 1
fi

echo "Updated Successfully"
