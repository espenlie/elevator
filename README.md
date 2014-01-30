elevator
========

Project in Real Time Programming course NTNU

Check for new elevators, or new init?

initFile containing:
* Num elevators?
* Port
* IP of elevators

elevator object ( 
* connectonMap 
* statusArr 
    ..* Errorhandling
* queueArr 
* floor?
)

sharedVariables:
* New order
* New staus


Thoughts/questions:
* An elevator dies, and comes back to life. Needs to be updated.  How to set
* the order if the elevator will stop in several floors on the
  way.
* Handling of shared variables. Can we secure that the queue will be
  handled properly?
* The queue array should "always" be empty.  
* During init, should we have a
 period 1-3min searching and establishing connections between the elevators)
* Do we need dummy connections? Could it be some channel in connection in the map?
 Need to figure out what is most appropriate.
* What about refused connection? First elevator who detects sets 
status on the unreachable elevator. That means it needs to 
be updated when connected again.



