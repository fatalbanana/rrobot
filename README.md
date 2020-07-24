## rrobot

Experimental fast test framework for Rspamd

## Usage

```
./rrobot --help
Usage of ./rrobot:
  -concurrency int
    	How many goroutines (default 8)
  -config string
    	Config file(s) to process
  -url string
    	Rspamd URL (default "http://127.0.0.1:11333/checkv2")
```

## Example config

```
test "MY TEST NAME" {
  assertions = [
    "Result.Action == 'no action' || Result.Action == 'greylist'",
    "not('FAKE_REPLY' in Result.Symbols)",
  ]
  inputs = [
    "/home/adam/emails/kaminski-v/*/*",
  ]
  headers = {
    IP = "127.0.0.1",
  }
}
```

The format of the file is [Hashicorp HCL](https://github.com/hashicorp/hcl).
[antonmedv/expr](https://github.com/antonmedv/expr) is used for assertions.
