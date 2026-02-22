<?php

/*
    var db_path - /index

    fun num_to_base62 - /utils/encode
*/

function get_db_file(
    string $db_name,
    string $user_name,
    string $file_name
): string {
    global $db_path;

    $today = num_to_base62(strtotime('today'));
    $dir = $db_path . '/' . $db_name . '/' . $user_name . '/' . $today;
    if (!is_dir($dir)) {
        mkdir($dir, 0777, true);
    }

    return $dir . '/' . $file_name . '.json';
}