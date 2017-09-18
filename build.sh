rm -rf release
mkdir -p release/{conf,i18n,runtime}
GOOS=linux GOARCH=amd64 go build -o starstore
mv starstore release
cp conf/configuration.yaml release/conf/
cp conf/production.yaml release/conf/
cp i18n/*yaml release/i18n/
cp supervisor.conf release
