# IPO Pilot

Automated IPO Application System for MeroShare

Professional solution for monitoring and automatically applying to IPOs

## Activation

This software uses a one-time password-based activation system:

### First-Time Activation
1. On first run, you'll be prompted to enter an activation password
2. The activation is device-specific and never expires
3. You only need to activate once per device

### Password Information
- For activation password, contact the developer: **Mahesh**
- Password is case-sensitive
- Once activated, no renewal is needed

### For Developers: Password Management
If you need to change the activation password, modify the `correctPassword` variable in the source code:

```bash
go run license_manager.go
```

The license manager allows you to:
- Generate valid license keys for customers
- Track all issued licenses
- Revoke licenses when necessary
- Validate license key formats

See `LICENSE_MANAGER_README.md` for detailed instructions.

### Important Notes
- Licenses are hardware-bound and cannot be transferred
- One license per machine only
- Software will not function without a valid license

## Installation and Running

### Prerequisites
- Go programming language installed (version 1.19 or higher)
- Windows operating system

### Running on Windows

#### Method 1: Run directly with Go
1. Open Command Prompt or PowerShell
2. Navigate to the IPO~Master directory:
   ```
   cd "c:\path\to\IPO~Master"
   ```
3. Run the program:
   ```
   go run main.go
   ```
4. The menu will appear for profile management and IPO monitoring

#### Method 2: Build and run executable
1. Open Command Prompt or PowerShell
2. Navigate to the IPO~Master directory:
   ```
   cd "c:\path\to\IPO~Master"
   ```
3. Build the executable:
   ```
   go build main.go
   ```
4. Run the executable:
   ```
   .\main.exe
   ```
5. The menu will appear for profile management and IPO monitoring

### Running in Terminal (PowerShell/Command Prompt)

1. **Open Terminal**: Press `Win + R`, type `cmd` or `powershell`, press Enter
2. **Navigate to Directory**: 
   ```
   cd "c:\1 VS Code\1\IPO~Master"
   ```
3. **Run the Program**:
   ```
   go run main.go
   ```
4. **Follow Menu Prompts**:
   - Choose option 1 to create a new profile
   - Choose option 2 to run continuous monitoring
   - Choose option 3 to run all profiles once
   - Choose option 4 to exit

### Notes
- Keep the terminal window open while running continuous monitoring
- The program will check for IPOs every 5 minutes in continuous mode
- All data is stored locally and encrypted for security

## Menu Options:
1. **Create a new profile** - Set up a new named account profile with credentials and settings
2. **Run selected profiles continuously** - Choose multiple profiles to monitor and apply to IPOs automatically every 5 minutes
3. **Run all profiles once** - Execute all existing profiles once to apply to current open IPOs
4. **Manage profiles (edit/delete)** - Edit profile names or delete unwanted profiles
5. **License information** - View current license status, expiry date, and days remaining
6. **Exit** - Quit the application

## Features:
- Multi-account support with named profiles
- Continuous monitoring mode with automatic retries
- Secure credential encryption (passwords, PINs, CRN stored encrypted locally)
- Automatic IPO application with error recovery
- Profile management (create, edit, delete)
- Colored console output for better user experience
- Input validation and sanitization
- Configurable kittas per profile
- Profiles store all information locally for automatic reuse