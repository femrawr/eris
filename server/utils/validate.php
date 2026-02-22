<?php

/*
    var body - /index
*/

function validate_fields(array $schema): void {
    global $body;

    $missing = [];

    foreach ($schema as $long => $short) {
        if (!isset($body[$long])) {
            $missing[] = $long;
            continue;
        }

        if ($body[$long] === '') {
            $missing[] = $long;
            continue;
        }

        if ($body[$long] === null) {
            $missing[] = $long;
            continue;
        }
    }

    if (!empty($missing)) {
        $missing_str = implode(', ', $missing);

        http_response_code(400);
        echo json_encode([
            'ok' => false,
            'err' => "missing field(s) - $missing_str"
        ]);

        exit;
    }
}