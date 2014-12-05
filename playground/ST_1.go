/*

OUTPUT STANDARD 'ST_1'

1. Dots

1.1. Output pattern of dot is `X&Y` where 'X' and 'Y' are dot
coordinates:

	1&2
	2&1
	31&10

1.2. Dot set is printing as `dot1;dot2;dot3` where dotN is packed dot:

	12&12;12&13;12&14;14&15

2. Directions

2.1. Direction packs as exclamation mark and letter of a direction:

	!n
	!s
	!e
	!w

2.2 Used directions

n - north
s - sounth
e - east
w - west

3. Objects

3.1 Output pattern of objects is `'TYPE'ID[OBJECT_DATA]` where
'TYPE' is a type of object (in single quotes), 'ID' is object
identifier on playground, 'OBJECT_DATA' is a dataset which object
contains (string returned Pack method of object). For example:

	'apple'2[22&22]
	'apple'3[1&7]
	'corpse'5[1&7;6&6;5&7;9&8]
	'snake'10[key%3%1&2;21&2;21&3]

'OBJECT_DATA' cannot contain chars '[', ']', ','

3.3 In a objectsets objects are separated by comma:

	'apple'3[1&7],'snake'10[3%1&2;21&2;21&3]

Objects of different types cannot have equal type name markers.

4. Cached objects

4.1. If object state was not updated it are created in cache list.
In cache list objects are represented like numbers according with it's
identifiers on playground separated by comma:

	1,2,3,4,5,6,7

4.2. Identifiers of cached objects are printed in objectset:

	1,2,'apple'3[1&7],4,5,6,7,8,9,'snake'10[3%1&2;21&2;21&3]

4.3 Objects cannot be created in object set twice: as identifier
(like cached object) and as updated packed object

*/
package playground
