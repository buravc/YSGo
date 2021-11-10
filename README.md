# YSGo
YS Go Task I was assigned to.

## Description
This is an Api library with a Redis like in-memory key-value (NoSQL) storage implementation called as Godis.

## Why Godis? What does it mean?
This name was inspired from the consistency of library naming convention of Go Community.

## Installation
Godis uses two Environment Variable for configuration. One for backup interval `GodisBackupInterval`, and the other for the address to listen `GodisApiEndpoint`
`GodisBackupInterval` Supports hour, minute and second configs (E.g. `1h`, `12h`, `30m`, `11m`, `10s`)
`GodisApiEndpoint` is formatted as follows `<ip-address>:<port>`
Backup file path is `<os-default-temp-dir>/Latest-GodisBackup.json` (A symbolic link to the latest backup file)

#### Containerized:
Simply download the `Dockerfile.production` 
```bash
$ curl -O https://raw.githubusercontent.com/callduckk/YSGo/main/Dockerfile.production
```
and run this command within the directory of the downloaded file
```bash
$ docker build . -f Dockerfile.production
```
Then run the created container image with
```bash
$ docker run -p <host-port>:<container-port> <image-id>
```
You can change the default configurations for Environment Variables within the `Dockerfile.production`

```Dockerfile
ENV GodisApiEndpoint=0.0.0.0:8090
ENV GodisBackupInterval=1h
```
`docker-compose.yml` file was prepared with a [reverse proxy](https://github.com/callduckk/YSGo-nginxReverseProxy) container in mind.
#### Non-Containerized
Clone the git repository with
```bash
$ git clone https://github.com/callduckk/YSGo
```
then build the module in the repository directory using
```bash
$ go build -o <output-file-path> ./cmd 
```
Configure environment variables
```bash
$ export GodisApiEndpoint=0.0.0.0:8090
$ export GodisBackupInterval=1h
```
Finally start the Api by simply executing the build output from the previous command.
```bash
$ <output-file-path>
```
## Usage
Three endpoints are currently supported.
* `/get?key=<key>` 
* `/set` 
* `/flush`

#### /get
This endpoint supports GET verb only. Requires one parameter `key`. Queries the Godis with the provided key and returns the corresponding value for the queried key.
Returns
```json
{
  "success": true/false,
  "value": <value>
}
```
Example request
```
/get?key=firstName
```
#### /set
This endpoint supports POST verb with `Content-Type: application/json`. Allows user to store new key-value pair or change existing ones. Required parameters are as follows
```json
{
	"key": <key>,
	"value": <value>
}
```
Returns
```json
{
  "success": true/false
}
```
### /flush
This endpoint supports GET verb only. Commands Godis to clear all key-value pairs.
Returns
```json
{
  "success": true/false
}
```

## Questions
#### Why not use Redis?
Task description stated that I shall only use golang stdlib. This is the same reason why I didn't use `cron` utility, instead implemented my own.

#### Why are backups saved to `/tmp`?
Task description included a sample backup directory with `/tmp`. There is no other reason.

#### Is Godis thread safe?
Godis is developed with thread safety in mind. Therefore it uses `sync/Map` library to achieve atomic transactions.
The implemented singleton pattern uses mutex as a lock mechanism to prevent data races from occurring.

#### Which style guide was followed?
[Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) was followed.

## Limitations
Godis only supports string:string key-value pairs. So it can't store collections neither as key nor as value.
