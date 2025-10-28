### co-ordinator

manages the sync and routing 
- [ ] publish ip to app layer so app can send register request   
- [ ] build a lookup routing table using a hashmap/map data structure 
- [ ] persist the lookup routing table  
- [ ] config option to dynamically manage number of instances  
- [ ] have a prepare req for each commit which supplies a uuid to caller who requires a txn  
- [ ] keeps txn log which records status of txn  
