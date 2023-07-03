# Anagram Finder Project

## Description

Task description can be found [here](Platform Engineer Case Description-Anagram.pdf).

## Project

### Part 1

API takes a JSON Request object for fetching new datasource.

`Request` validation achieved with `fiber`, `Ctx`, and `BodyParser` functions.

- JSON Request object should contain `rawUrl` of the remote file

### Part 2

API takes a JSON Request object for finding anagrams.

- By iterating over content...

#### Valid Edge Cases

##### Empty

## TODO

- [] App
 + [] Strategies
   - [x] `PrimeMultiplication`
   - [] `BitwiseMatching`
   - [] ~~LettersSorted~~
 + [x] API
   - [x] `DataSource` handlers
   - [x] `FindAnagrams` handler
 + [] Unit Test
- [] Kubernetes
  + [x] Minimal Docker image
  + [x] Kubernetes configuration
  + [] Helm chart
- [] Further Improvements
  + [] Update README.md
  + [] Linting
  + [] Swagger
  + [] `Bombardier` benchmark
