# Fixator

HTTP Service to control the speed limit of vehicles.
## Overview
Service has a rest api with 3 routes to fixate, select sample for the day and return min/max velocities for the day.
The main idea is to use a file system instead of a database to optimize the performance of service.
 
### File System structure
![Filesystem structure](img/Fs.png?raw=true "FS Structure")
* Each day has its own folder.
* Each folder contains files in the format "velocity range"-"hour of the day".
    * What is velocity range?
    Assume that our maximum velocity is 250. Let's split our range of possible velocities.
    The simplest split is ranges from 0-10, 10-20, ..., 240-250, but assuming there are fewer cars with extreme 
    velocities during the day I split our range into the first chunk from 0 to 40, the last from 150 to 250, and 
    in between 40-50, 50-60 ... - as an attempt to make the uniform distribution of entries between the files. 
* Based on the velocity and hour of fixation we will append the entry to the specific file.
For ex. 18.08.2020 14:11:12 PKIT-87 KK 67,3 goes to the file '70-12'.
* Such structure is useful for fast search and save of entries. 

### Service work
#### Save fixation
![Fix logic](img/Fix.png?raw=true "Fix logic")

* Fix method just puts the fixation in buffer. The buffer size should be chosen, depending on the load of the service.
 Now default size is 50.
* Worker goroutine check the directory and defines what goroutine should be used to save the file, 
depending on the velocity range of the fixation.
This structure helps prevent the collisions between goroutines while working with the filesystem. 
#### Select MinMax
![MinMax logic](img/MinMax.png?raw=true "MinMax")
* Min and Max are looking for independently. The main idea is to find min in the range, starting from 40.
If there is at least one file for this range we can guarantee that there will be a minimum. It means that 
there are no need to check other ranges. The same applies to max.
This idea helps to optimize the time for search. 
#### Select Range
![Select logic](img/Select.png?raw=true "Select")
* First we run goroutine to select entries for different hours for the day. It means that if we return sorted arrays 
from each of this goroutine we will need just to concat them to return the array sorted by time.
* For each goroutine we run new goroutines to return entries sorted by time for each velocity range starting from the 
range that contains the boundary. 
For ex. if we want to return cars that exceeds 65,5 velocity we will check only files after 70 range.
The sort happens using the heap. First range (ex. 70) selects only the entries that exceeds the boundary, others 
select all the entries from the file.
* When the result is ready "hour" goroutine will merge sorted arrays using the heap. 
## Usage
1) Service can be set up using config.yaml file.
Use --yaml flag to define the path, otherwise default is used.

2) Service start and end hour in .yaml defines the time when select requests can be handled.
Out of this range service returns the error. 
## ToDo:
Router:
* add default config
* add time period in a message when service unavailable
* add validation for query string params

Save:
* Choose best buffer size depending on the service load.
* We can add some postback url to approve that a fixation is in the system.
* What to do with the limit num of available goroutines 8k?
* Add an interface to work with os file system.

Tests:
* add