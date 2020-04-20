Each elevator has it own "state table" which is a 2D-array containing information about the elevators ID, network online/offline status, motor functionality status,
last floor visited, motor movement direction and whether any of the buttons have been pressed. Each elevator has a map containing the state tables of all the elevators.
The elevators broadcast their state tables repeatedly, so that this map is synchronized. The local elevator calculates its next move by using this map. 


