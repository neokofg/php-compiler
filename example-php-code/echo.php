<?php
// This is test for php one-line comments ignoring
/*
* This is the test for php multi-line comments ignoring
*/
$var = 1;
// for ($i = 0;$i < 5;$i++) {
//     echo $i . "\n";
//     if ($i == 3) {
//         continue;
//     }
// }
while(true) {
    echo $var++;
    if ($var == 5) {
        break;
    }
}
echo "its here!" . "\n";
