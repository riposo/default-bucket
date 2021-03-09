# Default Bucket Plugin

Default bucket plugin for [Riposo](https://github.com/riposo/riposo).

This core plugin adds support for default buckets which is implicitly created for
each user account on on their first use as outlined in the original
[Kinto](https://docs.kinto-storage.org/en/latest/api/1.x/buckets.html#personal-bucket-default)
documentation.

## Configuration

The following environment variable can be used to configure Riposo:

| Variable                            | Description                                                              | Default           |
| ----------------------------------- | ------------------------------------------------------------------------ | ----------------- |
| `RIPOSO_DEFAULT_BUCKET_HASH_SECRET` | Hex-encoded secret used for internal signatures, see (Secrets)(#secrets) | _none_ (required) |

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
