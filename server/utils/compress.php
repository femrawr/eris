<?php

/*
    fun num_to_base62 - /utils/encode
    fun base62_to_num - /utils/encode
*/

function compress(array $data, array $schema): array {
    $result = [];

    foreach ($schema as $long => $what) {
        $short = is_array($what) ? $what[0] : $what;
        $type = is_array($what) ? $what[1] : null;

        if (!isset($data[$long])) {
            continue;
        }

        $value = $data[$long];

        if ($type === 'bool') {
            $result[$short] = $value ? 1 : 0;
            continue;
        }

        if ($type === 'int') {
            $result[$short] = num_to_base62((int) $value);
            continue;
        }

        $result[$short] = $value;
    }

    return $result;
}

function decompress(array $data, array $schema): array {
    $result = [];

    foreach ($schema as $long => $what) {
        $short = is_array($what) ? $what[0] : $what;
        $type = is_array($what) ? $what[1] : null;

        if (!isset($data[$short])) {
            continue;
        }

        $value = $data[$short];

        if ($type === 'bool') {
            $result[$long] = $value === 1 || $value === '1';
            continue;
        }

        if ($type === 'int') {
            $result[$long] = base62_to_num($value);
            continue;
        }

        $result[$long] = $value;
    }

    return $result;
}