# Plesk Mass Account Checker

## Overview
This tool allows you to perform bulk login checks for Plesk accounts using credentials from a `list.txt` file. It supports multi-threading to speed up the process and saves results in organized files.

## Features
- **Multi-threaded checking**: Allows users to specify the number of threads for faster execution.
- **Credential validation**: Detects successful and failed logins.
- **Organized results**: Saves results in `result/Success.txt` and `result/Failed.txt`.
- **Simple CLI interface**: Easy-to-use command-line tool.

## Installation
### Prerequisites
- Go 1.18 or later installed
- Internet connection
- A `list.txt` file containing login credentials in `username:password:host` format

### Clone the Repository
```sh
git clone https://github.com/yourusername/plesk-checker.git
cd plesk-checker
```

### Install Dependencies
No external dependencies are required beyond the Go standard library.

### Build the Project
```sh
go build -o plesk-checker
```

## Usage
### Running the Tool
```sh
./plesk-checker
```

1. Enter the number of threads when prompted.
2. The tool will read credentials from `list.txt`.
3. Results will be saved in:
   - `result/Success.txt` for valid logins.
   - `result/Failed.txt` for invalid logins.

### Example of `list.txt`
```
admin:password123:example.com
user:testpass:plesk.hosting.com
```

### Expected Output
```sh
Enter number of threads: 10
[SUCCESS] - user:testpass:plesk.hosting.com
[FAILED] - admin:password123:example.com
Done.
```

## Contributing
Feel free to open issues or submit pull requests to improve the tool.

## License
This project is licensed under the MIT License.
