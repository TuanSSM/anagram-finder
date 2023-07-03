# Anagram Finder Project

## Description

Task description can be found [here](Platform%20Engineer%20Case%20Description-Anagram.pdf).

### Aim

Aim is to build a kubernetes compliant Anagram Finder Microservice, providing multiple Algorithms to draw contrast between their efficencies.

## Project

### Part 1

API takes a JSON Request object for fetching new datasource.

`Request` validation achieved with `fiber`, `Ctx`, and `BodyParser` functions.

- JSON Request object should contain `rawUrl` of the remote file

### Part 2

API takes a JSON Request object for finding anagrams.

- Iterates over datasource lines, with given Strategy

#### Prime Multiplication Strategy

By making use of Fundamental Theorem of Algebra,

+ Assign a prime number for each letter in the alphabet with considering dictionary frequencies
+ Multiply each letter with it's corresponding prime number, in order to obtain an unique number for equivalent anagrams

+ Pros
  - Does not need a sorting
+ Cons
  - Unique numbers get exponentially larger as anagrams spread into words

#### Bit Encoded Matching Strategy

##### Encoding 
`string`s are encoded into a 27 boolean bits & weights list

`AlphabetBools`: 27 boolean bits for each letter in alphabet including `'`, characters with higher frequency are assigned to more significant bitweights array for each letter

##### Mathcing Anagrams
1. Check equivalence of `AlphabetBools`
2. If True check equivalence of `Weights`
3. Both True => Anagram found

+ Pros
  - Low memory usage

### Part 3

1. Write singleton anagrams (only one word) to a file with unique anagram identifier, one word per line.
2. Generate 2 word anagrams, repeat same write procedure with a combined anagram per line in a directory with name number of words
3. By making use of last generated n word anagrams and singleton anagrams, generate n+1 word anagrams.
4. Squash files with same name as lines with a delimeter `,` to a single result file

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

```

## TODO

- [ ] App
 + [ ] Strategies
   - [x] `PrimeMultiplication`
   - [ ] `BitwiseMatching`
   - [ ] ~~LettersSorted~~
 + [x] API
   - [x] `DataSource` handlers
   - [x] `FindAnagrams` handler
 + [ ] Unit Test
 + [ ] Bounded concurrency, semaphores
- [ ] Kubernetes
  + [x] Minimal Docker image
  + [x] Kubernetes configuration
  + [ ] Helm chart
- [ ] Further Improvements
  + [ ] README.md
  + [ ] Linting
  + [ ] Swagger
  + [ ] `Bombardier` benchmark