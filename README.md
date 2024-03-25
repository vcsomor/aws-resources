# aws-resources

Extract your AWS resources into files.

Currently supported resources are:
- S3
- RDS

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
# set up your ${HOME}/.aws/credentials and ${HOME}/.aws/config file or 
# setup the the AWS_XYZ environemntal veraibales according to your need
$ ./bin/aws-resource list --regions us-west-2 --resources s3,rds --threads 8 --output file,stdout
$ ./bin/aws-resource help # for more options
```

### Testing

``make test``
