# Security Scanner in Docker

## Run locally
- Build: `./build.sh nopush`
- Create `input/input.txt` and `output`.
- Run: `docker run --rm -v $(pwd)/input:/input:ro -v $(pwd)/output:/output registry.gitlab.com/security.surf/scanners/subjack:dev /input/input.txt`

## Parameters

```
Usage of /subjack:
  -a	Find those hidden gems by sending requests to every URL. (Default: Requests are only sent to URLs with identified CNAMEs).
  -c string
    	Path to configuration file. (default "/fingerprints.json")
  -m	Flag the presence of a dead record, but valid CNAME entry.
  -o string
    	Output results to file (Subjack will write JSON if file ends with '.json'). (default "/output/output.txt")
  -ssl
    	Force HTTPS connections (May increase accuracy (Default: http://).
  -t int
    	Number of concurrent threads (Default: 10). (default 10)
  -timeout int
    	Seconds to wait before connection timeout (Default: 10). (default 10)
  -v	Display more information per each request.
```

## References

* https://github.com/m4ll0k/takeover
* https://github.com/DDuarte/aquatone/commit/8306559987fc93162dec15696935dfcac53db2a3
* https://github.com/EdOverflow/can-i-take-over-xyz 
* https://github.com/EdOverflow/can-i-take-over-xyz/issues/29 - cloudfront edge case
