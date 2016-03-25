library(reshape2)

### Use to compare the output of:
### https://github.com/jefferickson/peer-object-matcher [Go]
### https://github.com/jefferickson/peer-object-matching [Python Prototype]

go.version <- read.csv("/tmp/go_test.csv", header=FALSE)
py.version <- read.csv("/tmp/py_test.csv", header=FALSE)

go.long <- melt(go.version, id.vars="V1")
go.long$variable <- as.character(go.long$variable)

py.long <- melt(py.version, id.vars="V1")
py.long$variable <- as.character(py.long$variable)

comp <- merge(go.long, py.long, by=c("V1", "variable"), all.x=TRUE, all.y=TRUE)
comp$match <- ifelse(comp$value.x == comp$value.y, 1, 0)

print(nrow(comp[comp$match == 0, ]))