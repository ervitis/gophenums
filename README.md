# Gophenums

Generate enum types with type-safe.

## Getting Started

### Prerequisites

Requirements for the software and other tools to build, test and push

- [Go1.21](https://go.dev/doc/install)

### Installing

You can use the generator like this:

```bash
go get github.com/ervitis/gophenums
```

Then let's create a file with the types and constants. And add the comment associated to the type.

```go
// gophenum:generate color
type color string

const (
	Blue   color = "blue"
	Yellow color = "yellow"
	Red    color = "red"
)
```

> [!WARNING]  
> The type and the name of the type in the comment has to be the same or it won't work

## Running the tests

```bash
make tests
```

## Built With

- Go1.21

## Contributing

Please read [CONTRIBUTING.md](./.github/CONTRIBUTING.md) for details on our code
of conduct, and the process for submitting pull requests to us.

## Versioning

We use [Semantic Versioning](http://semver.org/) for versioning. For the versions
available, see the [tags on this
repository](https://github.com/PurpleBooth/a-good-readme-template/tags).

## Authors

- [@ervitis](https://github.com/ervitis)

## License

This project is licensed under the [Apache 2.0](LICENSE)
