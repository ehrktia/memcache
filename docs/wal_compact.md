### wal cache sync approach

- data is persisted in wal file

- at end of persistance, trigger update to cache



### compaction approach

- wal file exists

- every operation runs through wal

- check wal file size on regular basis

- have a default size configured

- if current size exceeds configured size

- start compaction

- create a new wal file 

- redirect operations to newly created file

- create a tar file

- add wal file with data to tar

- remove the old file

#### considerations

- how to handle operations to wal file during compaction?
  
  create a new wal file, redirect operations to new file

- how to handle any error or failure during compaction?
  
  remove the tar file created
  
  raise a warn level message to indicate failure

**Todo**

- [ ] find a way to recover from failure and deal with multiple wal files during recovery cycle

- [ ] what happens when drive used to store archive is full

- [ ] dependency on `.local` folder in `home` location 
