### Golang Send Mail From Excel File  
This is the simple example of sending mail from excel file using golang.

#### Environment Variables
```go  
host := getEnvOrFatal("MAIL_HOST")
port := getEnvOrFatal("MAIL_PORT")
username := getEnvOrFatal("MAIL_USERNAME")
password := getEnvOrFatal("MAIL_PASSWORD")
fromName := getEnvOrFatal("MAIL_FROM_NAME")
fromMail := getEnvOrFatal("MAIL_FROM_MAIL")
mailSubject := getEnvOrFatal("MAIL_SUBJECT")
```

#### How to run  

```bash
cp ./templates/email.html.sample ./templates/email.html # copy email sample
go build -o send-mail . # build
./send-mail test.xlsx # run
```