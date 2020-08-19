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
#### Select MinMax
![MinMax logic](img/MinMax.png?raw=true "MinMax")
#### Select Range
![Select logic](img/Select.png?raw=true "Select")
## Usage
## ToDo:
Router:
* add default config
* add time period in a message when service unavailable
* add validation for query string params

Save:
* Choose best buffer size depending on the service load.
* We can add some postback url to approve that a fixation is in the system.
* What to do with limit num of available goroutines 8k?