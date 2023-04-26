# About

Dumper is a core library to be used by other Updiver.Backup repositories

## How to use

Run this command to get help:
```
./build/backup -h
```

Clone all repos from your account (your personal account + teams you are member of):
```
make build-all && ./build/dumper dump -u username -d ~/destination_folder -t token_or_user_app_password
```