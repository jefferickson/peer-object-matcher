## Peer Object Matcher

#### Author: Jeff Erickson `<jeff@erick.so>`
#### Date: 2016-03-25

### Overview

Use to match objects (peer objects) based on categorical and continuous data. First, the objects are matched exactly based on their categorical data, then they are matched within each categorical group based on the Euclidean distance of their continuous data.

The input should be of the following form:

`object_id, categorical_data, no_match_groups, cont_point_1, cont_point_2, ..., cont_point_n`

with one object per line. Leaving `no_match_group` as a blank field will cause all objects to be compared within the categorical group. A sample input file, `test_data/sample_input_data.csv`, is provided for testing.

The output will be:

`object_id, peer_object_id_1, peer_object_id_2, ..., peer_object_id_m`

### Categorical Data

The categorical data can be anything. Grouping will be done on the unique values of this field. This can be a group label (`group1`, `group2`, etc.), concatenated categorical fields (`x:y` where `x` and `y` are different categorical flags), etc. While any number of data can be used, they must be concatenated into one field.

### Continuous Data

Any number of continuous dimensions can be used, as long as each object has the same number of dimensions. Objects will be matched based within their categorical groups based on the shortest distance between objects. Euclidean distance in _n_ dimensions is used here, but other distance algorithms could be substituted.

### Lag Peers

It is possible to specify a separate list of objects from which to peer. The lag file must have the same format as the input file. The objects in the lag file will be used as peers, but will not be peered themselves.

### Installation

After installing [Go](https://golang.org/dl/), follow these steps:
```
go get github.com/jefferickson/peer-object-matcher
go install github.com/jefferickson/peer-object-matcher
```

### Usage

```
$GOPATH/bin/peer-object-matcher --input /path/to/input.csv --output /path/to/output.csv
```
For a listing of all config flags, type:
```
$GOPATH/bin/peer-object-matcher --help
```

### Python Version

See [peer-object-matching](https://github.com/jefferickson/peer-object-matching) for the Python prototype of this peering algorithm.
