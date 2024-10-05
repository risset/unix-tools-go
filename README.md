# unix-utils-go

Simple Go implementations of classic Unix shell utilities.

## Implemented so far (work in progress)
### cat
- just plain concatenation

### find
- regex-based matching only
- has `-0` flag for usage with xargs (separate matches with null-separator)

### xargs
- has `-0` flag for usage with find (split input at null-separator)
