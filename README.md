# Asql
An systems programming project that explores the capabilities of compiler learning for use in the `Systems Programming` course.

## Project Stages

- [x] Lexer
Developed a tokenizer and lexer rules to differentiate tokens that belong to the syntactic table from those that do not.
- [ ] Parser
Build a parser and an AST (Abstract Syntax Tree) to determine the correct token sequence and verify that the grammar is correct.
- [ ] Semantics
Implementation of a semantic checker to detect compilation errors and type verification.
- [ ] Optimization
Final stage of the project where language ambiguities will be addressed and overall performance will be optimized.

## Usage

By default, `make` runs the standard rule and reads the source file `code.asql`
```sh
make
```
Compile the source code
```
make build
```
Once the source code is compiled, it can be integrated as a shell script inside a file and executed with the following example code.
```sql
#!./asql
SELECT ANOMBRE
FROM ALUMNOS,INSCRITOS,CARRERAS
WHERE ALUMNOS.A#=INSCRITOS.A# AND ALUMNOS.C#=CARRERAS.C#
AND INSCRITOS.SEMESTRE='2010I'
AND CARRERAS.CNOMBRE='ISC'
AND ALUMNOS.GENERACION='2010'
```
and run as
```sh
./code.asql
```
**expected output**
| IDENTIFIER | VALUE | LINE    |
|------------|-------|---------|
| ANOMBRE    | 401   | 1       |
| ALUMNOS    | 402   | 2,3,3,6 |
| INSCRITOS  | 403   | 2,3,4   |
| CARRERAS   | 404   | 2,3,5   |
| A#         | 405   | 3,3     |
| C#         | 406   | 3,3     |
| SEMESTRE   | 407   | 4       |
| CNOMBRE    | 408   | 5       |
| GENERACION | 409   | 6       |

| CONSTANT | VALUE | LINE |
|----------|-------|------|
| '2010I'  | 600   | 4    |
| 'ISC'    | 601   | 5    |
| '2010'   | 602   | 6    |

| NO. | LINE          | TOKEN      | TYPE | CODE |
|-----|---------------|------------|------|------|
| 1   | 1             | SELECT     | 1    | 10   |
| 2   | 2             | FROM       | 1    | 11   |
| 3   | 3             | WHERE      | 1    | 12   |
| 4   | 3,4,5,6       | AND        | 1    | 14   |
| 5   | 2,2           | ,          | 5    | 50   |
| 6   | 3,3,3,3,4,5,6 | .          | 5    | 51   |
| 7   | 3,3,4,5,6     | =          | 8    | 83   |
| 8   | 1             | ANOMBRE    | 4    | 401  |
| 9   | 2,3,3,6       | ALUMNOS    | 4    | 402  |
| 10  | 2,3,4         | INSCRITOS  | 4    | 403  |
| 11  | 2,3,5         | CARRERAS   | 4    | 404  |
| 12  | 3,3           | A#         | 4    | 405  |
| 13  | 3,3           | C#         | 4    | 406  |
| 14  | 4             | SEMESTRE   | 4    | 407  |
| 15  | 5             | CNOMBRE    | 4    | 408  |
| 16  | 6             | GENERACION | 4    | 409  |
| 17  | 4             | '2010I'    | 6    | 600  |
| 18  | 5             | 'ISC'      | 6    | 601  |
| 19  | 6             | '2010'     | 6    | 602  |

_More technical implementation details will be recorded in this same file for future reference or to showcase the work done throughout this course._
