# Station Manager: logging application

This is a desktop application for **Linux** (currently) and is used for logging QSOs.

# Exit Codes

This section lists the exit codes of the application.

- 0 = OK
- 100 = Failed to determine the working directory (the same directory in which the application is located).
- 101 = Failed to initialize the Ioc/Di container.

# Development Environment
This section describes the development environment for the application.

The main item here is the `.env` file. This file contains the environment `SM_WORKING_DIR` which is the absolute path
to the directory in which the application is located. This directory all contains the subdirectory `db` where the sqlite
database file is located.

