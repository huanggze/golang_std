# persistentConn

空闲超时(tooOld)
IdleAt + IdleConnTimeout < Now 

# dialer

deadline returns the earliest of:
- now+Timeout
- d.Deadline
- the context's deadline
 