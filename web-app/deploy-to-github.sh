#!/bin/bash

# 🚀 IPO Pilot - GitHub Deployment Script
# This script helps you push your code to GitHub and deploy

set -e  # Exit on error

echo "🚀 IPO Pilot - Deploying to GitHub"
echo "=================================="
echo ""

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed. Please install git first."
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "❌ Please run this script from the web-app directory"
    exit 1
fi

# Initialize git if not already
if [ ! -d ".git" ]; then
    echo "📦 Initializing git repository..."
    git init
    echo "✓ Git initialized"
fi

# Add all files
echo "📝 Adding files to git..."
git add .

# Create commit
echo "💾 Creating commit..."
read -p "Enter commit message (or press Enter for default): " commit_msg
if [ -z "$commit_msg" ]; then
    commit_msg="Update IPO Pilot Web Platform"
fi
git commit -m "$commit_msg" || echo "ℹ️  No changes to commit"

# Check if remote exists
if ! git remote | grep -q "origin"; then
    echo ""
    echo "📡 No remote repository configured."
    echo ""
    echo "Please choose an option:"
    echo "1. Use GitHub CLI (gh) - Recommended"
    echo "2. Enter GitHub repository URL manually"
    read -p "Your choice (1 or 2): " choice
    
    if [ "$choice" == "1" ]; then
        # Check if gh is installed
        if ! command -v gh &> /dev/null; then
            echo "❌ GitHub CLI (gh) is not installed."
            echo "Install it from: https://cli.github.com/"
            echo "Or choose option 2 to enter URL manually."
            exit 1
        fi
        
        echo "🔐 Authenticating with GitHub..."
        gh auth status || gh auth login
        
        read -p "Enter repository name (e.g., ipo-pilot-web): " repo_name
        if [ -z "$repo_name" ]; then
            repo_name="ipo-pilot-web"
        fi
        
        echo "📦 Creating GitHub repository: $repo_name"
        gh repo create "$repo_name" --public --source=. --remote=origin --push
        
        echo ""
        echo "✅ Repository created and code pushed!"
        echo "🌐 Repository URL: https://github.com/$(gh api user -q .login)/$repo_name"
        
    else
        echo ""
        echo "Please create a repository on GitHub:"
        echo "1. Go to: https://github.com/new"
        echo "2. Create repository (e.g., ipo-pilot-web)"
        echo "3. Copy the repository URL"
        echo ""
        read -p "Enter GitHub repository URL (e.g., https://github.com/username/repo.git): " repo_url
        
        if [ -z "$repo_url" ]; then
            echo "❌ No URL provided. Exiting."
            exit 1
        fi
        
        git remote add origin "$repo_url"
        git branch -M main
        echo "🚀 Pushing to GitHub..."
        git push -u origin main
        
        echo ""
        echo "✅ Code pushed to GitHub!"
        echo "🌐 Repository URL: ${repo_url%.git}"
    fi
else
    echo "🚀 Pushing to GitHub..."
    git push
    echo "✅ Code pushed to GitHub!"
fi

echo ""
echo "=================================="
echo "✅ DEPLOYMENT TO GITHUB COMPLETE!"
echo "=================================="
echo ""
echo "📋 NEXT STEPS - Deploy to Cloud Platform:"
echo ""
echo "🔹 Option 1: Railway (Recommended - Easiest)"
echo "   1. Visit: https://railway.app"
echo "   2. Click 'Start a New Project'"
echo "   3. Select 'Deploy from GitHub repo'"
echo "   4. Choose your repository"
echo "   5. Click 'Deploy'"
echo "   → Live in 2-3 minutes!"
echo ""
echo "🔹 Option 2: Render (100% Free)"
echo "   1. Visit: https://render.com"
echo "   2. Click 'New +' → 'Web Service'"
echo "   3. Connect your GitHub repository"
echo "   4. Click 'Create Web Service'"
echo "   → Live in 5 minutes!"
echo ""
echo "🔹 Option 3: Fly.io (Great for Go)"
echo "   1. Install: curl -L https://fly.io/install.sh | sh"
echo "   2. Run: fly auth login"
echo "   3. Run: fly launch"
echo "   → Live in 3 minutes!"
echo ""
echo "📖 Full deployment guide: GITHUB_DEPLOYMENT.md"
echo ""
