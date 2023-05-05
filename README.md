# About

Dumper is a lib to dump bitbucket / github repositories.

## Usage examples

Some examples how to use library can be found in [examples folder](https://github.com/Updiver/dumper/tree/dev/examples/dumper) folder.

### Run examples

Clone all repos from your account (your personal account + teams you are member of):
```
make build-all && ./build/dumper dump github -u username -d ~/destination_folder -t token_or_user_app_password
```
## Tests 

```
make run-tests
```