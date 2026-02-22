<?php

const CHARS = '1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM';

function num_to_base62(int $num): string {
    $result = '';

    while ($num > 0) {
        $result = CHARS[$num % 62] . $result;
        $num = intdiv($num, 62);
    }

    return $result;
}

function base62_to_num(string $base62): int {
    $result = 0;

    for ($i = 0; $i < strlen($base62); $i++) {
        $result = $result * 62 + strpos(CHARS, $base62[$i]);
    }

    return $result;
}