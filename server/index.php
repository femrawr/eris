<?php

header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

$method = $_SERVER['REQUEST_METHOD'];

if ($method === 'OPTIONS') {
    http_response_code(204);
    exit;
}

if ($method !== 'POST') {
    http_response_code(405);
    echo json_encode([
        'ok' => false,
        'err' => 'bad method'
    ]);

    exit;
}

require_once __DIR__ . '/schemas/reddit.php';

require_once __DIR__ . '/utils/compress.php';
require_once __DIR__ . '/utils/database.php';
require_once __DIR__ . '/utils/encode.php';
require_once __DIR__ . '/utils/validate.php';

header('Content-Type: application/json');

$db_path = __DIR__ . '/../_database';
if (!is_dir($db_path)) {
    mkdir($db_path, 0777, true);
}

$db_path = realpath($db_path);

$body = json_decode(file_get_contents('php://input'), true) ?? [];