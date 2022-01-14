# Default Bucket Plugin

Default bucket plugin for [Riposo](https://github.com/riposo/riposo).

This core plugin adds support for default buckets which is implicitly created for
each user account on on their first use as outlined in the original
[Kinto](https://docs.kinto-storage.org/en/latest/api/1.x/buckets.html#personal-bucket-default)
documentation.

## Configuration

The following additional configuration options can be used:

| Option                  | Type     | Description                                                               | Default           |
| ----------------------- | -------- | ------------------------------------------------------------------------- | ----------------- |
| `default_bucket.secret` | `string` | Hex-encoded secret used for internal signatures, see (Secrets)(#secrets)` | _none_ (required) |

### Secrets

You must generate and configure a 32-byte cryptographically strong hex-encoded secret. To generate a good random secret use:

```shell
openssl rand -hex 32
```

Alternatively:

```shell
xxd -c 32 -l 32 -p /dev/urandom
tr -dc 'a-f0-9' < /dev/urandom | head -c64; echo
```

## License

Copyright 2021-2022 Black Square Media Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this material except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
