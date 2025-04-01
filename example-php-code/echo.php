<?php
// This is test for php one-line comments ignoring
/*
* This is the test for php multi-line comments ignoring
*/
$var = 3;
// for ($i = 0;$i < 5;$i++) {
//     echo $i . "\n";
//     if ($i == 3) {
//         continue;
//     }
// }
// --
// while(true) {
//     echo $var++;
//     if ($var == 5) {
//         break;
//     }
// }
// --
switch ($var) {
    case 1:
        echo "its 1\n";
        break;
    case 2:
        echo "its 2\n";
        break;
    default:
        echo "its something odd\n";
}
echo "its here!" . "\n";
