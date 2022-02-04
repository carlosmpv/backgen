# Backgen

A simple golang API generator that uses key/value based databases.

It does not provide the database itself, only uses a interface to access set/get functions.


## Installing

```sh
go install github.com/carlosmpv/backgen
```

## Usage

```sh
backgen <Model name> <field type> <field name> <field type> <field name> ...
```

It will generate a struct with the first argument as name, and add fields with json tags given its types.
The generated struct will also have getters and setters that updates their states into the database.

An api with handler for creating, updating and deleting data based on fiber will also be generated.
