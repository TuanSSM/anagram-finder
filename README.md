# Anagram Finder Project

## Description

Task description can be found [here](Platform%20Engineer%20Case%20Description-Anagram.pdf).

### Aim

Aim is to build a kubernetes compliant Anagram Finder Microservice, providing multiple Algorithms to draw contrast between their efficencies.

## Project

### Part 1

API takes a JSON Request object for creating new datasource item on database.

`Request` validation achieved with `fiber`, `Ctx`, and `BodyParser` functions.

- JSON Request object should contain `rawUrl` of the remote file

### Part 2

API takes a JSON Request object for finding anagrams.

- Iterates over datasource lines, with given Strategy

#### Strategy

##### Encoding 
`string`s are encoded into a 27 boolean bits &  27 integer weights array

`AlphabetBools`: 27 boolean bits for each letter in alphabet including `'`, characters with higher frequency are assigned to more significant bitweights array for each letter

Input lines and combined lines with same boolean bits are stored in the same `<uint32>.json` files. Anagrams are appended to corresponding `weights`.

### Part 3

json files in the `data/<datasource-slug-id>` directory is read concurrently and entries having anagrams are inserted to database in batches.

## Building

### Running with Makefile

```bash
make run
```

### Running with Docker

```bash
docker-compose up -d
```

### Running with Kubernetes

```bash
kubectl apply -f kubernetes/
```