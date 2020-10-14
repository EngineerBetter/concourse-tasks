# Concourse Task Tester

Integration-tests Concourse tasks using `fly execute`, using a YAML test definition format.

## Usage

```terminal
$ ginkgo -- --specs some_spec.yml
```

Or alternatively

```terminal
$ ginkgo build .
$ ./wibble-test --specs some_spec.yml
```

## Spec Format

See example `*_spec.yml` files.

