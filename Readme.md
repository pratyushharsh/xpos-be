Legacy

Layer will be exported using aws sam

While building the zip for image compression use the following commabg to bulf sharp dependency for linux.
rm -rf node_modules/sharp
SHARP_IGNORE_GLOBAL_LIBVIPS=1 npm install --arch=x64 --platform=linux sharp


## To create new lambda function in go

Create a new directory in **go-lambda-func**.\
cd [new directory] && go mod init [new directory]