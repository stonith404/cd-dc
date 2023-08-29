# CD-DC

CD-DC (Continuous Deployment Docker Compose) is a simple service to continuously deploy docker containers in a docker compose stack. The tool is especially useful in a GitHub Actions workflow.

## Usage

Send a post request to the `/upgrade/<service>` endpoint:

```bash
curl <host-of-ci-cd>/upgrade/<service> -X POST -H <api-key> --fail-with-body
```

## Production Installation

### Install as Systemd service

1. Copy the binary to `/opt/cd-dc`
2. Copy the `config.yml` to `/opt/cd-dc` and modify it accordingly
3. Create a systemd service file in `/etc/systemd/system/cd-dc.service`

   ```
   [Unit]
   Description=CD-DC
   After=network.target

   [Service]
   ExecStart=/opt/cd-dc/cd-dc
   WorkingDirectory=/opt/cd-dc

   [Install]
   WantedBy=default.target
   ```

4. Enable the service
   ```
   sudo systemctl enable cd-dc
   ```
5. Start the service
   ```
   sudo systemctl start cd-dc
   ```
6. Check the status
   ```
   sudo systemctl status cd-dc
   ```
7. Check the logs
   ```
   sudo journalctl -u cd-dc
   ```

#### Update

1. Copy the new binary to `/opt/cd-dc`
2. Reload the systemd daemon
   ```
   sudo systemctl daemon-reload
   ```
3. Restart the service
   ```
   systemctl restart cd-dc
   ```

## Development

### Run

```bash
go run ./cmd
```

### Build

```bash
env GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o build/cd-dc ./cmd
```
