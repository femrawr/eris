<?php

/*
    var body - /index
*/

require __DIR__ . '/../index.php';

$type = $body['type'] ?? null;
if (!$type) {
    http_response_code(400);
    echo json_encode([
        'ok' => false,
        'err' => 'no type'
    ]);

    exit;
}

$handler = __DIR__ . '/../funcs/reddit/' . $type . '.php';
if (!file_exists($handler)) {
    http_response_code(400);
    echo json_encode([
        'ok' => false,
        'err' => "bad type - $type"
    ]);

    exit;
}

require $handler;