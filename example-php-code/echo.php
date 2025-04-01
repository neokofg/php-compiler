<?php
// This is test for php one-line comments ignoring
/*
* This is the test for php multi-line comments ignoring
*/
$var = 1;
while ($var <= 10) {
    echo $var++ . "\n";
    if ($var === "5") {
        echo "var is 5" . "\n";
    }
}