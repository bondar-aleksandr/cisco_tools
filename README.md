## cisco_tools

This web app is used to host CLI executables for various network tools, as well as UI representation of  [**bondar-aleksandr/cisco_parser**](https://github.com/bondar-aleksandr/cisco_parser)
___
## Usage
Typical folder structure for app to run is below

```
.
└── app_root_dir/
    ├── app_executable
    ├── config/
    │   └── config.yml
    ├── temp/
    ├── downloads/
    │   └── additional tools tar-files for download
    └── ui/
        ├── html
        └── static
```
App reads config parameters from **./config/config.yml** file. Config parameters are as follow:


| Option | Description |
| ------ | ----------- |
| server, host   | Ip address for webapp to listen to, default - all interfaces |
| server, port   | webapp HTTP port |
| server, readTimeout | webapp server read timeout, seconds |
| server, writeTimeout | webapp server write timeout, seconds |
| server, idleTimeout | webapp server idle timeout, seconds |
| server, maxUpload | webapp maximum request size, Megabytes |
| server, uploadMIMETypes | list of allowed for upload MIME types (files with configuration) |

**config.yml** example:
```yaml
server:
  host: "localhost"
  port: "4000"
  readTimeout: 5      # seconds
  writeTimeout: 10    # seconds
  idleTimeout: 60     # seconds
  maxUpload: 1        # Mbytes
  uploadMIMETypes:
    - "text/plain"
```
